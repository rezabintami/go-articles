package users

import (
	"context"
	"errors"
	"fmt"
	"go-articles/constants"
	"go-articles/helpers"
	"go-articles/modules/roles"
	"go-articles/server/config"
	"go-articles/server/middleware"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis"
)

const (
	forgotPasswordLatestAction = "latestAction"
	forgotPassword             = "forgotPassword"
	forgotPasswordKey          = "forgotPasswordKey"
	forgotPasswordVerifyKey    = "verify_key"
	forgotPasswordUserLocked   = "user_locked"
	defaultForgotPasswordDur   = 24 * time.Hour
	sessionHashName            = "admin:login"
)

type UserUsecase struct {
	userRepository Repository
	contextTimeout time.Duration
	jwtAuth        *middleware.ConfigJWT
	jwtForgot      *middleware.ConfigForgotJWT
	mail           helpers.MailConnection
	redis          *redis.Client
}

func NewUserUsecase(ur Repository, rr roles.Repository, jwtAuth *middleware.ConfigJWT, jwtForgot *middleware.ConfigForgotJWT, timeout time.Duration, mail helpers.MailConnection, redis *redis.Client) Usecase {
	return &UserUsecase{
		userRepository: ur,
		jwtAuth:        jwtAuth,
		jwtForgot:      jwtForgot,
		contextTimeout: timeout,
		mail:           mail,
		redis:          redis,
	}
}

func (usecase *UserUsecase) Login(ctx context.Context, email, password string) (string, string, error) {
	existedUser, err := usecase.userRepository.GetByEmail(ctx, email)
	if err != nil {
		return "", "", err
	}

	if !helpers.ValidateHash(password, existedUser.Password) {
		return "", "", constants.ErrEmailPasswordNotFound
	}

	accessToken := usecase.jwtAuth.GenerateToken(existedUser.ID, existedUser.Name, existedUser.Role.Name)
	refreshToken := usecase.jwtAuth.GenerateRefreshToken(existedUser.ID)

	return accessToken, refreshToken, nil
}

func (usecase *UserUsecase) GetByID(ctx context.Context, id int) (Domain, error) {
	users, err := usecase.userRepository.GetByID(ctx, id)

	if err != nil {
		return Domain{}, err
	}

	return users, nil
}

func (usecase *UserUsecase) Register(ctx context.Context, userDomain *Domain) error {
	ctx, cancel := context.WithTimeout(ctx, usecase.contextTimeout)
	defer cancel()

	existedUser, err := usecase.userRepository.GetByEmail(ctx, userDomain.Email)
	if err != nil {
		if !strings.Contains(err.Error(), "not found") {
			return err
		}
	}

	if existedUser != (Domain{}) {
		return constants.ErrDuplicateData
	}

	userDomain.Password, _ = helpers.Hash(userDomain.Password)

	err = usecase.userRepository.Register(ctx, userDomain)
	return err
}

func (usecase *UserUsecase) Update(ctx context.Context, userDomain *Domain, id int) error {
	ctx, cancel := context.WithTimeout(ctx, usecase.contextTimeout)
	defer cancel()

	err := usecase.userRepository.Update(ctx, userDomain, helpers.CheckStringNull(userDomain.Password), id)
	if err != nil {
		return err
	}

	return nil
}

func (usecase *UserUsecase) SetForgotPasswordRedis(key string, value string, duration time.Duration) error {
	fmt.Println("key :", key, "value :", value, "duration :", duration)
	err := usecase.redis.Set(key, value, duration).Err()
	return err
}

func (usecase *UserUsecase) SetSessionRedis(sessionName string, id int, time time.Time) error {
	fmt.Println("sessionName :", sessionName, "id :", id, "time :", time)
	idString := strconv.Itoa(id)
	err := usecase.redis.HSet(sessionName, idString, time).Err()
	return err
}

func (usecase *UserUsecase) GetForgotPasswordRedis(key string) (string, error) {
	value, err := usecase.redis.Get(key).Result()
	if err == redis.Nil || value == "" {
		err = errors.New("Invalid Key")
	}

	return value, nil
}

func (usecase *UserUsecase) DeleteForgotPasswordRedis(key string) error {
	return usecase.redis.Del(key).Err()
}

func (usecase *UserUsecase) ForgotPassword(ctx context.Context, email string) error {
	key := helpers.RandLowerAlphanumericString(50)
	users, err := usecase.userRepository.GetByEmail(ctx, email)
	if err != nil {
		return err
	}

	err = usecase.SetForgotPasswordRedis(forgotPassword+string(key), strconv.Itoa(users.ID), defaultForgotPasswordDur)
	if err != nil {
		return err
	}

	err = usecase.SetForgotPasswordRedis(forgotPasswordKey+strconv.Itoa(users.ID), string(key), defaultForgotPasswordDur)
	if err != nil {
		return err
	}

	forgotPasswordLink := config.GetConfiguration("app.development_url") + config.GetConfiguration("feature.forgot_password.path") + "?key=" + string(key)
	stringEmailBody := `<p>Hai, %s.</p>
				<p>Terima kasih telah melakukan permintaan reset kata sandi di articles. Mohon memperbaharui kata sandi Anda melalui <a href="%s" target="_blank">laman</a> ini.</p>
				<p>Abaikan email ini jika Anda tidak melakukan permintaan reset kata sandi.</p>`
	body := fmt.Sprintf(stringEmailBody, users.Name, forgotPasswordLink)
	mailContent := helpers.MailContent{
		From:    config.GetConfiguration("mail.default"),
		To:      email,
		Subject: "Permintaan Reset Kata Sandi",
		Body:    body,
	}

	_, err = usecase.mail.MailSend(mailContent.From, mailContent.To, mailContent.Subject, mailContent.Body)

	return err
}

func (usecase *UserUsecase) VerifyTokenForgotPassword(ctx context.Context, key string) (string, string, error) {
	id, err := usecase.GetForgotPasswordRedis(forgotPassword + key)
	if err != nil {
		return "", "", err
	}

	// Check valid key
	validKey, err := usecase.GetForgotPasswordRedis(forgotPasswordKey + id)
	if err != nil {
		return "", "", err
	}

	if validKey != key {
		return "", "", errors.New("Expired Key")
	}

	newId, _ := strconv.Atoi(id)
	data, err := usecase.userRepository.GetByID(ctx, newId)
	if err != nil {
		return "", "", err
	}

	err = usecase.SetForgotPasswordRedis(forgotPasswordLatestAction+id, forgotPasswordVerifyKey, defaultForgotPasswordDur)
	if err != nil {
		return "", "", err
	}

	err = usecase.DeleteForgotPasswordRedis(forgotPassword + key)
	if err != nil {
		return "", "", err
	}

	expSecretToken, _ := strconv.Atoi(config.GetConfiguration("jwt.exp_secret"))
	expSecretTokenTime := time.Duration(expSecretToken)

	token, tokenExpiredAt, err := usecase.jwtForgot.GenerateTokenResetPassword(data.ID, expSecretTokenTime)
	if err != nil {
		return "", "", errors.New("Invalid generate token")
	}

	err = usecase.SetSessionRedis(sessionHashName, data.ID, time.Now())
	if err != nil {
		return "", "", err
	}

	return token, tokenExpiredAt, err
}

func (usecase *UserUsecase) SetForgotPassword(ctx context.Context, id int, password string) error {
	hashPassword, _ := helpers.Hash(password)
	err := usecase.userRepository.SetPassword(ctx, id, hashPassword)
	return err
}

func (usecase *UserUsecase) Fetch(ctx context.Context, page, perpage int, by, search, sort string) ([]Domain, int, error) {
	if page <= 0 {
		page = 1
	}

	if perpage <= 0 {
		perpage = 25
	}

	if strings.ToUpper(sort) != "ASC" && strings.ToUpper(sort) != "DESC" {
		sort = "ASC"
	}
	switch by {
		case "name":
			by = `users."name"`
		case "roles":
			by = `roles."name"`
		default:
			by = `users."created_at"`
	}

	res, total, err := usecase.userRepository.Fetch(ctx, page, perpage, by, strings.ToLower(search), sort)
	if err != nil {
		return []Domain{}, 0, err
	}

	return res, total, nil
}

func (usecase *UserUsecase) Delete(ctx context.Context, id int) error {
	err := usecase.userRepository.Delete(ctx, id)

	if err != nil {
		return err
	}

	return nil
}

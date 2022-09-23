package users_test

// import (
// 	"context"
// 	"errors"
// 	"go-articles/app/middleware"
// 	"go-articles/helpers"
// 	"go-articles/modules/users"
// 	userMock "go-articles/modules/users/mocks"
// 	"os"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// )

// var (
// 	userRepository userMock.Repository
// 	userUsecase    users.Usecase
// 	jwtAuth        *middleware.ConfigJWT
// )

// func setup() {
// 	jwtAuth = &middleware.ConfigJWT{SecretJWT: "abc123", ExpiresDuration: 2}
// 	userUsecase = users.NewUserUsecase(&userRepository, jwtAuth, 2)
// }

// func TestMain(m *testing.M) {
// 	setup()
// 	os.Exit(m.Run())
// }

// func TestLoginUsers(t *testing.T) {
// 	t.Run("test case 1, valid test", func(t *testing.T) {
// 		email := "agent@gmail.com"
// 		sso := false
// 		pass, _ := helpers.Hash("123123")
// 		usersDomain := users.Domain{
// 			ID:       1,
// 			Password: pass,
// 			Name:     "agent",
// 			Email:    email,
// 		}

// 		userRepository.On("GetByEmail", mock.Anything, mock.AnythingOfType("string")).Return(usersDomain, nil).Once()

// 		accessToken, refreshToken, err := userUsecase.Login(context.Background(), email, "123123", sso)
// 		assert.Nil(t, err)
// 		assert.NotEmpty(t, accessToken)
// 		assert.NotEmpty(t, refreshToken)
// 	})
// 	t.Run("test case 2, password error", func(t *testing.T) {
// 		pass, _ := helpers.Hash("1231231")
// 		usersDomain := users.Domain{
// 			ID:       1,
// 			Password: pass,
// 			Name:     "users",
// 			Email:    "users@gmail.com",
// 		}

// 		userRepository.On("GetByEmail", mock.Anything, mock.AnythingOfType("string")).Return(usersDomain, nil).Once()

// 		_, _, err := userUsecase.Login(context.Background(), "users@gmail.com", "123123", false)
// 		assert.Equal(t, err, helpers.ErrEmailPasswordNotFound)

// 	})

// 	t.Run("test case 3, error record", func(t *testing.T) {

// 		errRepository := errors.New("error record")
// 		userRepository.On("GetByEmail", mock.Anything, mock.AnythingOfType("string")).Return(users.Domain{}, errRepository).Once()

// 		accessToken, refreshToken, err := userUsecase.Login(context.Background(), "logistik@gmail.com", "123123", false)

// 		assert.Equal(t, err, errRepository)
// 		assert.Equal(t, "", accessToken)
// 		assert.Equal(t, "", refreshToken)
// 	})
// }

// func TestRegisterUsers(t *testing.T) {
// 	t.Run("test case 1, valid test", func(t *testing.T) {
// 		domain := users.Domain{
// 			ID:       1,
// 			Password: "asyudasd820aisd",
// 			Name:     "users",
// 			Email:    "users@gmail.com",
// 		}
// 		userRepository.On("GetByEmail", mock.Anything, mock.Anything).Return(users.Domain{}, nil).Once()
// 		userRepository.On("Register", mock.Anything, mock.Anything).Return(nil).Once()

// 		err := userUsecase.Register(context.Background(), &domain, false)

// 		assert.Nil(t, err)
// 	})

// 	t.Run("test case 2, duplicate data", func(t *testing.T) {
// 		domain := users.Domain{
// 			ID:       1,
// 			Password: "asyudasd820aisd",
// 			Name:     "users",
// 			Email:    "users@gmail.com",
// 		}

// 		errRepository := errors.New("duplicate data")
// 		userRepository.On("GetByEmail", mock.Anything, mock.Anything).Return(domain, errRepository).Once()

// 		err := userUsecase.Register(context.Background(), &domain, false)

// 		assert.Equal(t, err, helpers.ErrDuplicateData)
// 	})

// 	t.Run("test case 3, data has exist", func(t *testing.T) {
// 		domain := users.Domain{
// 			ID:       1,
// 			Password: "asyudasd820aisd",
// 			Name:     "users",
// 			Email:    "users@gmail.com",
// 		}

// 		userRepository.On("GetByEmail", mock.Anything, mock.Anything).Return(domain, nil).Once()

// 		err := userUsecase.Register(context.Background(), &domain, false)

// 		assert.Equal(t, err, helpers.ErrDuplicateData)
// 	})

// 	t.Run("test case 4, register failed", func(t *testing.T) {
// 		domain := users.Domain{
// 			ID:    1,
// 			Name:  "users",
// 			Email: "users@gmail.com",
// 		}
// 		errRepository := errors.New("register failed")
// 		userRepository.On("GetByEmail", mock.Anything, mock.Anything).Return(users.Domain{}, nil).Once()
// 		userRepository.On("Register", mock.Anything, mock.Anything).Return(errRepository).Once()

// 		err := userUsecase.Register(context.Background(), &domain, false)

// 		assert.Equal(t, err, errRepository)
// 	})
// }

// func TestGetByIdUsers(t *testing.T) {
// 	t.Run("test case 1, valid test", func(t *testing.T) {
// 		domain := users.Domain{
// 			ID:    1,
// 			Name:  "users",
// 			Email: "users@gmail.com",
// 		}
// 		userRepository.On("GetByID", mock.Anything, mock.AnythingOfType("int")).Return(domain, nil).Once()

// 		result, err := userUsecase.GetByID(context.Background(), 1)

// 		assert.Nil(t, err)
// 		assert.Equal(t, domain.ID, result.ID)
// 	})

// 	t.Run("test case 2, data not found", func(t *testing.T) {
// 		errRepository := errors.New("data not found")
// 		userRepository.On("GetByID", mock.Anything, mock.AnythingOfType("int")).Return(users.Domain{}, errRepository).Once()
// 		result, err := userUsecase.GetByID(context.Background(), -1)
// 		assert.Equal(t, result, users.Domain{})
// 		assert.Equal(t, err, errRepository)
// 	})
// }

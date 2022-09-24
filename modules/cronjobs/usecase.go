package cronjobs

import (
	"context"
	"fmt"
	"go-articles/modules/articles"
	"time"

	"github.com/jasonlvhit/gocron"
)

type CronjobsUsecase struct {
	articleUsecase articles.Usecase
}

func NewCronjobUsecase(au articles.Usecase) *CronjobsUsecase {
	return &CronjobsUsecase{
		articleUsecase: au,
	}
}

func (usecase *CronjobsUsecase) RegisterCronjob() {
	location, _ := time.LoadLocation("Asia/Jakarta")
	gocron.ChangeLoc(location)
	schedule := gocron.NewScheduler()
	schedule.Every(1).Minute().Do(usecase.GetCountArticles)
	<-schedule.Start()
}

// GetCountArticles implements Usecase
func (usecase *CronjobsUsecase) GetCountArticles() {
	count, err := usecase.articleUsecase.Count(context.Background())
	if err != nil {
		fmt.Println("error count : ",err)
	}
	fmt.Println("Count :",count)
}
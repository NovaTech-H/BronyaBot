package service

import (
	"BronyaBot/config"
	"BronyaBot/global"
	"BronyaBot/internal/entity"
	"BronyaBot/internal/service/cx_service"
	"BronyaBot/internal/service/gongxueyun_service"
	"github.com/robfig/cron/v3"
)

type AppService struct {
	users []entity.SignEntity
	cron  *cron.Cron
}

func NewAppService() *AppService {
	return &AppService{
		cron: cron.New(),
	}
}

func (svc *AppService) Init() {
	svc.scheduleTasks()
	svc.cron.Start()
	select {}
}

func (svc *AppService) scheduleTasks() {
	taskCount := len(global.Config.Tasks)
	global.Log.Infof("Total tasks to schedule: %d", taskCount)
	for _, task := range global.Config.Tasks {
		global.Log.Infof("Loading Task: Cron: %v; Description: %v", task.Cron, task.Description)
		svc.addTask(task)
	}
	global.Log.Info("Scheduling tasks...")
}

func (svc *AppService) addTask(task config.CronTask) {
	svc.cron.AddFunc(task.Cron, func() {
		global.Log.Infof("Running task: %v", task.Description)
		svc.StartGongxueYun(task.TaskType)
		global.Log.Info("Task finished!")
	})
}

func (svc *AppService) StartGongxueYun(taskType string) {
	svc.users = gongxueyun_service.LoadUsers()
	global.Log.Info("Starting Gongxueyun module...")
	for _, user := range svc.users {
		svc.createMoguDing(user).Run(taskType)
	}
}

func (svc *AppService) StartTestCX() {
	cxLogic := cx_service.CxLogic{}
	cxLogic.Run()
}

func (svc *AppService) createMoguDing(user entity.SignEntity) *gongxueyun_service.MoguDing {
	return &gongxueyun_service.MoguDing{
		ID:          user.ID,
		PhoneNumber: user.Username,
		Password:    user.Password,
		Email:       user.Email,
		Sign: gongxueyun_service.SignInfo{
			City:      user.City,
			Area:      user.Area,
			Address:   user.Address,
			Country:   user.Country,
			Province:  user.Province,
			Latitude:  user.Latitude,
			Longitude: user.Longitude,
		},
	}
}

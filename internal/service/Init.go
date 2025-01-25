package service

import (
	"BronyaBot/config"
	"BronyaBot/global"
	"BronyaBot/internal/entity"
	"BronyaBot/internal/service/cx_service"
	"BronyaBot/internal/service/gongxueyun_service"
	"github.com/robfig/cron/v3"
	"time"
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

	// 获取任务总数
	taskCount := len(global.Config.Tasks)
	global.Log.Infof("Total tasks to schedule: %d", taskCount)

	// 遍历任务并添加
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

//func (svc *AppService) scheduleTasks() {
//	global.Log.Info("Scheduling tasks...")
//
//	svc.addCronTask("0 8 * * *", "每天早上8点签到", "sign")
//	svc.addCronTask("0 18 * * *", "每天晚上6点签到", "sign")
//	svc.addCronTask("0 10 * * 5", "每周周五早上10点提交周报", "week")
//	svc.cron.AddFunc("0 10 ? * 1L", func() {
//		if isLastWeek(time.Now()) {
//			global.Log.Info("Running task: 每月最后一周的周一早上10点提交月报")
//			svc.StartGongxueYun("month")
//			global.Log.Info("Task finished!")
//		}
//	})
//}
//
//func (svc *AppService) addCronTask(schedule, logMessage, taskType string) {
//	svc.cron.AddFunc(schedule, func() {
//		global.Log.Infof("Running task: %s", logMessage)
//		svc.StartGongxueYun(taskType)
//		global.Log.Info("Task finished!")
//	})
//}

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

func isLastWeek(t time.Time) bool {
	_, week := t.ISOWeek()
	nextMonday := t.AddDate(0, 0, 7-int(t.Weekday()))
	nextMonthWeek, _ := nextMonday.ISOWeek()
	return week != nextMonthWeek
}

package config

type CronTask struct {
	Cron        string `yaml:"cron"`
	Description string `yaml:"description"`
	TaskType    string `yaml:"taskType"`
}

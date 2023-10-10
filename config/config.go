package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		Worker
		Manager
		General
	}
	Worker struct {
		workerQueueSize int `yaml:"worker_queue_size" env:"WORKER_QUEUE_SIZE" env-default:"1000"`
	}
	Manager struct {
		numberOfWorkers int `yaml:"manager_number_of_workers" env:"MANAGER_NUMBER_OF_WORKERS" env-default:"10"`
	}
	General struct {
	}
)

func NewConfig(path string) (*Config, error) {
	var conf = &Config{}
	err := cleanenv.ReadConfig(path, conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

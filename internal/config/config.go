package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type Db struct {
}

type Config struct {
	DbHost     string `envconfig:"DB_HOST"`
	DbPort     string `envconfig:"DB_PORT"`
	DbName     string `envconfig:"DB_NAME"`
	DbUser     string `envconfig:"DB_USER"`
	DbPassword string `envconfig:"DB_PASSWORD"`

	ConnURI          string `envconfig:"CONN_URI" default:"http://bsm.api.iql.ru/ords/bsm/segmentation/get_segmentation"`
	ConnAuthLoginPwd string `envconfig:"CONN_AUTH_LOGIN_PWD" default:"4Dfddf5:jKlljHGH"`
	ConnUserAgent    string `envconfig:"CONN_USER_AGENT " default:"spacecount-test"`
	ConnTimeout      int64  `envconfig:"CONN_TIMEOUT" default:"5"`
	ConnInterval     int64  `envconfig:"CONN_INTERVAL" default:"1500"`

	ImportBatchSize int64 `envconfig:"IMPORT_BATCH_SIZE" default:"50"`

	LogCleanupMaxAge int64  `envconfig:"LOG_CLEANUP_MAX_AGE" default:"7"`
	LogPath          string `envconfig:"LOG_PATH" default:"/log/segmentation_import.log"`
}

func Parse() *Config {
	res := &Config{}
	err := envconfig.Process("", res)
	if err != nil {
		log.Fatalf("failed to parse config: %s", err)
	}
	return res
}

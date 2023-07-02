package econf

import (
	"fmt"
	"os"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	vip *viper.Viper
	err error
)

func load() error {
	// if CONSUL_HTTP_ADDR is not found in env we only load config as filepath
	if os.Getenv("CONSUL_HTTP_ADDR") == "" {
		vip.AddConfigPath(GetEnv("ECONF_FILEPATH", "."))

		return vip.ReadInConfig()
	}

	// retrieve config from consul
	check := func() error {
		val, err := os.ReadFile(os.Getenv("CONSUL_HTTP_TOKEN_FILE"))
		if err != nil {
			logrus.Warnf("unable to read consul token file : %s", err.Error())
		}

		if string(val) == "" && os.Getenv("CONSUL_HTTP_TOKEN") == "" {
			return fmt.Errorf("unable to retrieve token")
		}

		return nil
	}

	notify := func(err error, t time.Duration) {
		logrus.Info(err.Error(), t)
	}

	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = 2 * time.Minute

	err = backoff.RetryNotify(check, b, notify)
	if err != nil {
		logrus.Info("http token can't be retrieved")
		panic(err)
	}

	viper.AddRemoteProvider("consul", os.Getenv("CONSUL_HTTP_ADDR"), os.Getenv("CONSUL_FILENAME"))
	return viper.ReadRemoteConfig()
}

func Configure() {
	if err := godotenv.Load(); err != nil {
		logrus.Warnf("unable to load .env file : %s", err.Error())
	}

	vip = viper.New()
	vip.SetConfigName(GetEnv("ECONF_FILENAME", "config.yaml"))
	vip.SetConfigType(GetEnv("ECONF_FILETYPE", "yaml"))
	vip.AutomaticEnv()

	err = load()

	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			logrus.Errorf("config file is not found : %s", err.Error())
		} else {
			logrus.Fatalf("unknown error : %s", err.Error())
		}
	}
}

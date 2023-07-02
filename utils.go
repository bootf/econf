package econf

import (
	"os"

	"github.com/spf13/viper"
)

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func Config() *viper.Viper {
	return vip
}

func GetString(key string) string {
	return vip.GetString(key)
}

func GetInt(k string) int {
	return vip.GetInt(k)
}

func GetFloat64(k string) float64 {
	return vip.GetFloat64(k)
}

func GetStringSlice(k string) []string {
	return vip.GetStringSlice(k)
}

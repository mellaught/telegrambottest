package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config interface {
	GetString(key string) string
	GetInt(key string) int
	GetBool(key string) bool
	GetFloat64(key string) float64
	Init()
}

type viperConfig struct {
}

func (v *viperConfig) Init() {
	viper.AutomaticEnv()

	viper.AddConfigPath(".")

	replacer := strings.NewReplacer(`.`, `_`)
	viper.SetEnvKeyReplacer(replacer)
	viper.SetConfigType(`json`)
	viper.SetConfigFile(`config.json`)
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}
	//panic
}

func (v *viperConfig) GetString(key string) string {
	return viper.GetString(key)
}

func (v *viperConfig) GetInt(key string) int {
	return viper.GetInt(key)
}

func (v *viperConfig) GetBool(key string) bool {
	return viper.GetBool(key)
}

func (v *viperConfig) GetFloat64(key string) float64 {
	return viper.GetFloat64(key)
}

func NewViperConfig() Config {
	v := &viperConfig{}
	v.Init()
	return v
}

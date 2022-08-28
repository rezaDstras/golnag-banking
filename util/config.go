package util

import "github.com/spf13/viper"

// Config stores all configurations of the application
// The values are read from viper from a config file or envirement variables
type Config struct {
	//unmarshal from viper
	DBDriver   string `mapstructure:"DB_DRIVER"`
	DBSource   string `mapstructure:"DB_SOURCE"`
	ServerAddr string `mapstructure:"SERVER_ADDR"`
}

func LoadConfig(path string) (cinfig Config , err error){
	viper.AddConfigPath(path)
	//name of file config which viper should read from
	viper.SetConfigName("app")
	// type of file config which viper should read from
	viper.SetConfigType("env") //like env , json , ymal , toml ,xaml

	//read automatically from envirement variables if exist
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	//unmarshal
	err = viper.Unmarshal(&cinfig)
	return
}
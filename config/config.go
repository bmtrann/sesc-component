package config

import (
	"time"

	"github.com/spf13/viper"
)

type AuthConfig struct {
	HashSalt string
	SignKey  []byte
	TokenTTL time.Duration
}

type DBConfig struct {
	URI               string
	Name              string
	UserCollection    string
	StudentCollection string
	CourseCollection  string
}

func Init() error {
	viper.AddConfigPath("./config")
	viper.SetConfigName("config")

	return viper.ReadInConfig()
}

func LoadAuthConfig() *AuthConfig {
	return &AuthConfig{
		viper.GetString("auth.hash_salt"),
		[]byte(viper.GetString("auth.sign_key")),
		viper.GetDuration("auth.token_ttl"),
	}
}

func LoadDBConfig() *DBConfig {
	return &DBConfig{
		viper.GetString("mongo.uri"),
		viper.GetString("mongo.db_name"),
		viper.GetString("mongo.user_collection"),
		viper.GetString("mongo.student_collection"),
		viper.GetString("mongo.course_collection"),
	}
}

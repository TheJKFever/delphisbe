package config

import (
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Environment string        `json:"env" mapstructure:"env"`
	DBConfig    DBConfig      `json:"db" mapstructure:"db"`
	SQLDBConfig SQLDBConfig   `json:"sqldb" mapstructure:"sqldb"`
	Twitter     TwitterConfig `json:"twitter" mapstructure:"twitter"`
	Auth        AuthConfig    `json:"authConfig"`
	AWS         AWSConfig     `json:"awsConfig" mapstructure:"aws"`
}

func (c *Config) ReadEnvAndUpdate() {
	viper.AutomaticEnv()
	c.Twitter.ConsumerKey, c.Twitter.ConsumerSecret, c.Auth.HMACSecret = viper.GetString("twitter_consumer_key"), viper.GetString("twitter_consumer_secret"), viper.GetString("auth_hmac_secret")
	c.SQLDBConfig.Username, c.SQLDBConfig.Password = viper.GetString("db_user"), viper.GetString("db_password")
}

type AWSConfig struct {
	Region         string         `json:"region" mapstructure:"region"`
	UseCredentials bool           `json:"useCredentials" mapstructure:"useCredentials"`
	Credentials    AWSCredentials `json:"credentials" mapstructure:"credentials"`
	IsFargate      bool           `json:"isFargate" mapstructure:"isFargate"`
}

type AWSCredentials struct {
	ID     string `json:"id" mapstructure:"id"`
	Secret string `json:"secret" mapstructure:"secret"`
	Token  string `json:"token" mapstructure:"token"`
}

type AuthConfig struct {
	HMACSecret string `json:"hmacSecret"`
	Domain     string `json:"domain" mapstructure:"domain"`
}

type TwitterConfig struct {
	ConsumerKey    string
	ConsumerSecret string
	Callback       string `json:"callback" mapstructure:"callback"`
	Redirect       string `json:"redirect" mapstructure:"redirect"`
}

type DBConfig struct {
	Host         string       `json:"host" mapstructure:"host"`
	Port         int          `json:"port" mapstructure:"port"`
	Region       string       `json:"region" mapstructure:"region"`
	TablesConfig TablesConfig `json:"tables_config" mapstructure:"tables_config"`
}

type SQLDBConfig struct {
	Host     string `json:"host" mapstructure:"host"`
	Port     int    `json:"port" mapstructure:"port"`
	DBName   string `json:"db_name" mapstructure:"db_name"`
	Username string
	Password string
}

type TablesConfig struct {
	Discussions   TableConfig `json:"discussions" mapstructure:"discussions"`
	Participants  TableConfig `json:"participants" mapstructure:"participants"`
	PostBookmarks TableConfig `json:"post_bookmarks" mapstructure:"post_bookmarks"`
	Posts         TableConfig `json:"posts" mapstructure:"posts"`
	Users         TableConfig `json:"users" mapstructure:"users"`
	UserProfiles  TableConfig `json:"user_profiles" mapstructure:"user_profiles"`
	Viewers       TableConfig `json:"viewers" mapstructure:"viewers"`
}

type TableConfig struct {
	TableName string `json:"table_name" mapstructure:"table_name"`
}

func AddConfigDirectory(dir string) {
	viper.AddConfigPath(dir)
}

func clearConfig() {
	viper.Reset()
}

func ReadConfig() (*Config, error) {
	env := os.Getenv("DELPHIS_ENV")
	if env == "" {
		env = "local"
	}
	viper.SetConfigType("json")
	viper.SetConfigName(env)
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	config := Config{}
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	viper.SetEnvPrefix("delphis")
	_ = viper.BindEnv("twitter_consumer_key", "twitter_consumer_secret", "auth_hmac_secret", "db_user", "db_password")
	config.ReadEnvAndUpdate()

	return &config, nil
}

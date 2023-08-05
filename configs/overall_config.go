package configs

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func parseConfig(cfgPath string) (*Config, error) {
	// Read the YAML file
	yamlFile, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		return &Config{}, err
	}

	// Parse the YAML data into a struct
	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return &Config{}, err
	}

	return &config, nil
}

type Config struct {
	MySQL               MySQLConfig               `yaml:"mysql"`
	Redis               RedisConfig               `yaml:"redis"`
	AuthenticateAndPost AuthenticateAndPostConfig `yaml:"authenticate_and_post_config"`
	Newsfeed            NewsfeedConfig            `yaml:"newsfeed_config"`
	NewsfeedPublishing  NewsfeedPublishingConfig  `yaml:"newsfeed_publishing_config"`
	Web                 WebConfig                 `yaml:"web_config"`
}

type MySQLConfig struct {
	DSN                       string `yaml:"dsn"`
	DefaultStringSize         int    `yaml:"defaultStringSize"`
	DisableDatetimePrecision  bool   `yaml:"disableDatetimePrecision"`
	DontSupportRenameIndex    bool   `yaml:"dontSupportRenameIndex"`
	SkipInitializeWithVersion bool   `yaml:"skipInitializeWithVersion"`
}

type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
}

type KafkaConfig struct {
	Topic   string   `yaml:"topic"`
	Brokers []string `yaml:"brokers"`
}

type AuthenticateAndPostConfig struct {
	Port               int         `yaml:"port"`
	MySQL              MySQLConfig `yaml:"mysql"`
	Redis              RedisConfig `yaml:"redis"`
	NewsfeedPublishing HostConfig  `yaml:"newsfeed_publishing"`
}

type NewsfeedConfig struct {
	Port                int         `yaml:"port"`
	MySQL               MySQLConfig `yaml:"mysql"`
	Redis               RedisConfig `yaml:"redis"`
	Kafka               KafkaConfig `yaml:"kafka"`
	AuthenticateAndPost HostConfig  `yaml:"authenticate_and_post"`
}

type WebConfig struct {
	Port                int         `yaml:"port"`
	APIVersions         []string    `yaml:"api_version"`
	AuthenticateAndPost HostConfig  `yaml:"authenticate_and_post"`
	Newsfeed            HostConfig  `yaml:"newsfeed"`
	Redis               RedisConfig `yaml:"redis"`
}

type NewsfeedPublishingConfig struct {
	Port                int         `yaml:"port"`
	Redis               RedisConfig `yaml:"redis"`
	Kafka               KafkaConfig `yaml:"kafka"`
	AuthenticateAndPost HostConfig  `yaml:"authenticate_and_post"`
}

type HostConfig struct {
	Hosts []string `yaml:"hosts"`
}

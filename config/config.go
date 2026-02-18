package config

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	ServiceHTTPPort int           `mapstructure:"API_HTTP_PORT"`
	ServiceGRPCPort int           `mapstructure:"API_GRPC_PORT"`
	IntervalSend    time.Duration `mapstructure:"INTERVAL_SEND"`
	Logger          `mapstructure:",squash"`
	Tracer          `mapstructure:",squash"`
	Meter           `mapstructure:",squash"`
	OtelExporter    string `mapstructure:"OTEL_EXPORTER"`
	Gmail           `mapstructure:",squash"`
	Ses             `mapstructure:",squash"`
	MailerSend      `mapstructure:",squash"`
	SendGrid        `mapstructure:",squash"`
}

type Logger struct {
	ServerMode string `mapstructure:"SERVER_MODE"`
	Encoding   string `mapstructure:"LOG_ENCODING"`
	Level      string `mapstructure:"LOG_LEVEL"`
}

type Tracer struct {
	Name              string `mapstructure:"TRACE_NAME"`
	TracerExporterURL string `mapstructure:"TRACER_EXPORTER_URL"`
}

type Meter struct {
	Name             string `mapstructure:"METER_NAME"`
	MeterExporterURL string `mapstructure:"METER_EXPORTER_URL"`
}

type Ses struct {
}

type Gmail struct {
	Identity string `mapstructure:"GMAIL_IDENTITY"`
	Username string `mapstructure:"GMAIL_USERNAME"`
	Password string `mapstructure:"GMAIL_PASSWORD"`
	Host     string `mapstructure:"GMAIL_HOST"`
	Port     int    `mapstructure:"GMAIL_PORT"`
}

type MailerSend struct {
	ApiKey string `mapstructure:"MAILERSEND_API_KEY"`
}

type SendGrid struct {
	SendGridApiKey string `mapstructrure:"SENDGRID_API_KEY"`
}

func Load(path string) *Config {

	if os.Getenv("PROD") == "true" {
		return loadFromEnv()
	}

	log.Println("env variable PROD=false reading from file")

	return loadFromFile(path)
}

func loadFromFile(path string) *Config {

	var config Config

	viper.SetConfigName("app")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("LoadFile.ReadInConfig: %v", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("LoadFile.Unmarshal: %v", err)
	}

	return &config
}

func recu(v any) {
	e := reflect.TypeOf(v)
	tp := e.Elem()

	for i := 0; i < tp.NumField(); i++ {
		field := tp.Field(i)

		switch tp.Field(i).Type.Kind() {
		case reflect.Int:
			fmt.Println("iinttttttt", field)
		case reflect.String:
			fmt.Println("string", field)

		case reflect.Struct:

			if field.Type == reflect.TypeOf(time.Duration(0)) {
				fmt.Println("duration", field)
			} else {
				recu(field)
			}
		}
	}
}

func loadFromEnv() *Config {

	var c *Config

	recu(c)

	return &Config{}
}

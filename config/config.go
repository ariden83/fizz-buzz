package config

import (
	"context"
	"fmt"
	"github.com/heetch/confita"
	"github.com/heetch/confita/backend"
	"github.com/heetch/confita/backend/env"
	"github.com/heetch/confita/backend/file"
	"github.com/heetch/confita/backend/flags"
	"os"
	"reflect"
	"strings"
	"time"
)

type Logger struct {
	Host     string `config:"logger_host"`
	Port     int    `config:"logger_port"`
	Level    string `config:"logger_level"`
	CLILevel string `config:"cli_level"`
}

type Swagger struct {
	Port     int    `config:"swager_port"`
	Host     string `config:"swager_host"`
	Path     string `config:"swager_path"`
	SpecPath string `config:"swager_specpath"`
}

type Metrics struct {
	Port int    `config:"metrics_port"`
	Host string `config:"metrics_host"`
}

type Config struct {
	Name      string
	Port      int
	Env       string `config:"env"`
	Host      string
	PublicURL string `config:"public_url"`

	CallHTTPTimeout     time.Duration `config:"call_http_timeout"`
	APIReadTimeout      time.Duration `config:"api_read_timeout"`
	APIWriteTimeout     time.Duration `config:"api_write_timeout"`
	HealthzReadTimeout  time.Duration `config:"healthz_read_timeout"`
	HealthzWriteTimeout time.Duration `config:"healthz_write_timeout"`

	MaxNBParametersLimit int `config:"max_nb_parameters_limit"`

	Logger

	Metrics

	Swagger
}

func getDefaultConfig() *Config {
	return &Config{
		Name: "fizz-buzz",
		Host: "127.0.0.1",
		Port: 8080,

		APIReadTimeout:      4,
		APIWriteTimeout:     100,
		HealthzReadTimeout:  10,
		HealthzWriteTimeout: 10,

		MaxNBParametersLimit: 100000,

		PublicURL: "127.0.0.1:8080",

		Swagger: Swagger{
			Port:     8082,
			Host:     "127.0.0.1",
			Path:     "swagger/",
			SpecPath: "swagger/swagger.yaml",
		},

		Logger: Logger{
			CLILevel: "INFO",
			Host:     "127.0.0.1",
			Port:     12201,
			Level:    "INFO",
		},

		Metrics: Metrics{
			Port: 8081,
			Host: "127.0.0.1",
		},
	}
}

// New Load the config
func New() *Config {
	loaders := []backend.Backend{
		env.NewBackend(),
		flags.NewBackend(),
	}

	environment := os.Getenv("CONFIGOR_ENV")
	if environment != "" {
		configFile := findConfigFilePathRecursively(environment, 0)
		loaders = append(loaders, file.NewBackend(configFile))
	}

	loader := confita.NewLoader(loaders...)

	cfg := getDefaultConfig()
	err := loader.Load(context.Background(), cfg)
	if err != nil {
		panic(err)
	}

	fmt.Println(fmt.Sprintf("%+v", cfg))
	return cfg
}

func findConfigFilePathRecursively(environment string, depth int) string {
	char := "../"
	if depth == 0 {
		char = "./"
	}

	filepath := strings.Repeat(char, depth) + "config/config." + environment + ".yaml"
	if _, err := os.Stat(filepath); err == nil {
		return filepath
	}
	depth++

	return findConfigFilePathRecursively(environment, depth)
}

func (c *Config) String() string {
	val := reflect.ValueOf(c).Elem()
	s := "\n-------------------------------\n"
	s += "-  Application configuration  -\n"
	s += "-------------------------------\n"
	for i := 0; i < val.NumField(); i++ {
		v := val.Field(i)
		t := val.Type().Field(i)
		c.applyWithType(&s, "", v, t)
	}
	return s
}

func (c *Config) applyWithType(s *string, parentKey string, v reflect.Value, k reflect.StructField) {
	obfuscate := false

	tag := k.Tag.Get("config")
	if idx := strings.Index(tag, ","); idx != -1 {
		opts := strings.Split(tag[idx+1:], ",")

		for _, opt := range opts {
			if opt == "obfuscate" {
				obfuscate = true
			}
		}
	}
	if !obfuscate {
		if parentKey != "" {
			parentKey += "-"
		}
		switch v.Kind() {
		case reflect.String:
			*s += fmt.Sprintf("%s: \"%v\"\n", parentKey+k.Name, v.Interface())
			return
		case reflect.Bool:
		case reflect.Int:
			*s += fmt.Sprintf("%s: %v\n", parentKey+k.Name, v.Interface())
			return
		case reflect.Struct:
			parentKey += k.Name
			c.DeepStructFields(s, parentKey, v.Interface())
			return
		}
	}
}

func (c *Config) DeepStructFields(s *string, parentKey string, iface interface{}) {
	ifv := reflect.ValueOf(iface)
	ift := reflect.TypeOf(iface)

	for i := 0; i < ift.NumField(); i++ {
		v := ifv.Field(i)
		t := ift.Field(i)
		c.applyWithType(s, parentKey, v, t)
	}
}

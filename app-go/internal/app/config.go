package app

import (
	"context"
	"embed"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/sethvargo/go-envconfig"
)

type EnvConfig struct {
	AppEnv      string `env:"APP_ENVIRONMENT"`
	HttpAddress string `env:"HTTP_ADDRESS"`
	Postgres    struct {
		ConnectionString string `env:"POSTGRES_CONNECTION_STRING"`
	}
	AWS struct {
		Sqs struct {
			ImageTransformQueueUrl string `env:"AWS_SQS_QUEUE_URL_IMAGE_TRANSFORMER"`

			// NOTE: add more queue urls here
			// ...
		}
	}
}

// NewConfig is to parse env either from file embedded and env vars attached on os
func NewConfig(embedFs embed.FS) *EnvConfig {
	appEnv := os.Getenv("APP_ENVIRONMENT")

	// load file from embed
	envs, err := embedFs.ReadFile(fmt.Sprintf("configs/%s.env", appEnv))
	if err != nil {
		log.Fatalf("failed to load .env file from embed.Fs. err=%v", err)
	}

	// load envs to runtime line by line
	lines := strings.Split(string(envs), "\n")
	for _, line := range lines {
		if line != "" {
			splits := strings.SplitN(line, "=", 2)

			// skip env definition from file
			// if its already defined on os
			if os.Getenv(splits[0]) != "" {
				continue
			}
			if err := os.Setenv(splits[0], splits[1]); err != nil {
				log.Fatalf("failed to inject .env values. env=%s. value%s. err=%v", splits[0], splits[1], err)
			}
		}
	}

	// read all required env values
	var cfg EnvConfig
	if err := envconfig.Process(context.Background(), &cfg); err != nil {
		log.Fatalf("failed to read environment variables. err=%v", err)
	}

	return &cfg
}

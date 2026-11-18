package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	config, err := LoadConfiguration("../")
	if err != nil {
		logrus.Fatalf("Error loading configuration: %v", err)
	}
	configBytes, _ := json.MarshalIndent(config, "", "  ")
	logrus.Info(string(configBytes))
}

type Configuration struct {
	NatsUrl        string   `mapstructure:"nats_url" json:"nats_url"`
	Backend        string   `mapstructure:"backend" json:"backend"`
	BackendBaseUrl string   `mapstructure:"backend_base_url" json:"backend_base_url"`
	Model          string   `mapstructure:"model" json:"model"`
	MinIntervalS   int      `mapstructure:"min_interval_s" json:"min_interval_s"`
	RemoteApiKey   string   `mapstructure:"-" json:"remote_api_key"`
	Tags           []string `mapstructure:"tags" json:"tags"`
	Categories     []string `mapstructure:"categories" json:"categories"`

	LlmResponseSchema struct {
		Schema      string `mapstructure:"$schema" json:"$schema"`
		Id          string `mapstructure:"$id" json:"$id"`
		Title       string `mapstructure:"title" json:"title"`
		Description string `mapstructure:"description" json:"description"`
		Type        string `mapstructure:"type" json:"type"`
		Properties  struct {
			Summary  Property `mapstructure:"summary" json:"summary"`
			Tags     Property `mapstructure:"tags" json:"tags"`
			Category Property `mapstructure:"category" json:"category"`
		} `mapstructure:"properties" json:"properties"`
		Required []string `mapstructure:"required" json:"required"`
	} `mapstructure:"llm_response_schema" json:"llm_response_schema"`
}

type Property struct {
	Type     string   `mapstructure:"type" json:"type"`
	Items    *Items   `mapstructure:"items" json:"items,omitempty"`
	Enum     []string `mapstructure:"enum" json:"enum,omitempty"`
	MinItems int      `mapstructure:"minItems" json:"minItems,omitempty"`
	MaxItems int      `mapstructure:"maxItems" json:"maxItems,omitempty"`
}

type Items struct {
	Type string   `mapstructure:"type" json:"type"`
	Enum []string `mapstructure:"enum" json:"enum,omitempty"`
}

// LoadConfiguration
// Reads configuration.json once, at startup. For each key, the environment overrides the file,
// and the file overrides the default: "model" is overridden by ANALYZER_MODEL.
func LoadConfiguration(configurationDirectory string) (Configuration, error) {
	var configuration Configuration
	configurationPath := filepath.Join(configurationDirectory, "configuration.json")

	viper.SetConfigFile(configurationPath)

	// all env variable starting with ANALYZER will be taken into account
	viper.SetEnvPrefix("ANALYZER")
	viper.AutomaticEnv()
	viper.BindEnv("nats_url", "NATS_URL")

	viper.SetDefault("nats_url", "nats://localhost:4222")
	viper.SetDefault("model", "gemini-3.5-flash-lite")
	viper.SetDefault("min_interval_s", 5)

	err := viper.ReadInConfig()
	if err != nil {
		return configuration, fmt.Errorf("cannot read %s: %w", configurationPath, err)
	}
	// put everything in configuration variable, ENV overrides file
	err = viper.Unmarshal(&configuration)
	if err != nil {
		return configuration, fmt.Errorf("%s: %w", configurationPath, err)
	}

	// fill tags and categories into llm response schema
	configuration.LlmResponseSchema.Properties.Tags.Items.Enum = configuration.Tags
	configuration.LlmResponseSchema.Properties.Category.Enum = configuration.Categories

	apiKey, err := ReadConfigurationSecret("REMOTE_API_KEY")
	if err != nil {
		return configuration, err
	}
	configuration.RemoteApiKey = apiKey

	return configuration, err
}

// ReadConfigurationSecret
// Reads a secret from the environment (NAME), or from the file whose path is in NAME_FILE.
// Never from configuration.json: it would be committed.
func ReadConfigurationSecret(environmentName string) (string, error) {
	value := os.Getenv(environmentName)
	path := os.Getenv(environmentName + "_FILE")

	if value != "" && path != "" {
		return "", fmt.Errorf("%s and %s_FILE are both set: use one or the other", environmentName, environmentName)
	}
	if value == "" && path == "" {
		return "", fmt.Errorf("neither %s nor %s_FILE is defined", environmentName, environmentName)
	}
	if path == "" {
		return value, nil
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("cannot read %s_FILE (%s): %w", environmentName, path, err)
	}
	return strings.TrimSpace(string(content)), nil
}

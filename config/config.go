package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type APIConfig struct {
	Secret string `yaml:"secret"`
	Key    string `yaml:"key"`
}

type TokenConfig struct {
	OAuthToken       string `yaml:"oauth_token"`
	OAuthTokenSecret string `yaml:"oauth_token_secret"`
}

func SetApiConfig(secret, key string) error {
	apiConfig := APIConfig{
		Secret: secret,
		Key:    key,
	}

	if err := writeConfigFile(apiConfig, "api"); err != nil {
		return fmt.Errorf("failed to save API config: %w", err)
	}

	return nil
}

func SetTokenConfig(oauthToken, oauthTokenSecret string) error {
	tokenConfig := TokenConfig{
		OAuthToken:       oauthToken,
		OAuthTokenSecret: oauthTokenSecret,
	}

	if err := writeConfigFile(tokenConfig, "token"); err != nil {
		return fmt.Errorf("failed to save Token config: %w", err)
	}

	return nil
}

func GetAPIConfig() (*APIConfig, error) {
	data, err := readConfigFile("api")
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var apiConfig APIConfig
	if err := yaml.Unmarshal(data, &apiConfig); err != nil {
		return nil, fmt.Errorf("failed to unmarshal API config from YAML: %w", err)
	}

	return &apiConfig, nil
}

func GetTokenConfig() (*TokenConfig, error) {
	data, err := readConfigFile("token")
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var tokenConfig TokenConfig
	if err := yaml.Unmarshal(data, &tokenConfig); err != nil {
		return nil, fmt.Errorf("failed to unmarshal Token config from YAML: %w", err)
	}

	return &tokenConfig, nil
}

func readConfigFile(file string) ([]byte, error) {
	configFilePath, err := getConfigFilePath(file)
	if err != nil {
		return nil, fmt.Errorf("failed to get config file: %w", err)
	}

	data, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	return data, nil
}

func writeConfigFile(data interface{}, file string) error {
	configFilePath, err := getConfigFilePath(file)

	if err != nil {
		return fmt.Errorf("failed to get config file: %w", err)
	}

	stringData, err := yaml.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal Token config to YAML: %w", err)
	}

	if err := os.WriteFile(configFilePath, stringData, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	return nil
}

func getConfigFilePath(file string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", errors.New("cannot get user home directory")
	}

	configDir := filepath.Join(homeDir, ".flkctl")
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return "", errors.New("cannot create config directory")
	}

	configFilePath := filepath.Join(configDir, file+".yaml")
	return configFilePath, nil
}

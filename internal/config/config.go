package config

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
)

type Config struct {
	APIEndpoint string
	APIKey      string
	Model       string
}

var (
	instance *Config
	once     sync.Once
	initErr  error
)

func Load() error {
	Get()
	return initErr
}

func Get() *Config {
	once.Do(func() {
		instance = &Config{}
		// Load API config from the environment variables
		api_endpoint := os.Getenv("DESCRIBER_API_ENDPOINT")
		api_key := os.Getenv("DESCRIBER_API_KEY")
		api_model := os.Getenv("DESCRIBER_MODEL")
		// Check Availability
		if api_key == "" {
			log.Println("WARN: Environment variable 'DESCRIBER_API_KEY' is not set.Please make sure the model do not need key to access.")
		}
		if api_endpoint == "" {
			initErr = fmt.Errorf("ERR: Environment variable 'DESCRIBER_API_ENDPOINT' is not set.")
		} else if !IsValidOpenaiEndpointURL(api_endpoint) {
			initErr = fmt.Errorf("ERR: Environment variable 'DESCRIBER_API_ENDPOINT' is not a valid OpenAI compatible endpoint URL.\n It should be an OpenAI compatible endpoint like 'https://your-provider.ai/api/v1/'. The last path segment should be the version, e.g., 'v1'.")
		} else if api_model == "" {
			initErr = fmt.Errorf("ERR: Environment variable 'DESCRIBER_MODEL' is not set")
		}
		instance.APIEndpoint = api_endpoint
		instance.APIKey = api_key
		instance.Model = api_model
	})
	return instance
}

func IsValidOpenaiEndpointURL(raw string) bool {

	// 1. parse URL
	u, err := url.Parse(raw)
	if err != nil {
		return false
	}

	if u.Scheme == "" || u.Host == "" {
		return false
	}

	path := strings.TrimSuffix(u.Path, "/")
	if path == "" {
		return false
	}

	lastPart := path[strings.LastIndex(path, "/")+1:]
	if lastPart == "" {
		return false
	}
	// 2. check if the last path segment matches version pattern like "v1", "v2", etc.
	var versionSegmentRegex = regexp.MustCompile(`^v\d+$`)
	return versionSegmentRegex.MatchString(lastPart)
}

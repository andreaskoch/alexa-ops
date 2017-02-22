package main

import (
	"os"
	"io/ioutil"
	"encoding/json"
)

// Config contains the configuration parameters for alexaops.
type Config struct {
	Skill AlexaSkill `json:"skill"`

	JenkinsAPI  *JenkinsAPI `json:"jenkinsAPI,omitempty"`
	Deployments []DeploymentConfig `json:"deployments"`
}

type AlexaSkill struct {
	// AppID is the application id from your
	// Amazon Dashboard (developer.amazon.com > Alexa > Alexa Skills Kit > Your Skill > ID)
	// Example: amzn1.ask.skill.abc12345-1111-dddd-aaaa-aaaabbbb3333
	AppID string `json:"appID"`
}

type JenkinsAPI struct {
	URL      string `json:"url"`
	Username string `json:"username"`
	APIToken string `json:"apiToken"`
}
type DeploymentConfig struct {
	Type string `json:"type"`
	Name string `json:"name"`

	Jenkins *JenkinsDeployment `json:"jenkins,omitempty"`
}

type JenkinsDeployment struct {
	JobName string `json:"jobName"`
}

func readConfigFromFile(configPath string) (Config, error) {

	configFile, err := os.Open(configPath)
	if err != nil {
		return Config{}, err
	}

	defer configFile.Close()

	configFileData, readErr := ioutil.ReadAll(configFile)
	if readErr != nil {
		return Config{}, readErr
	}

	var config Config
	jsonDecodeErr := json.Unmarshal(configFileData, &config)
	if jsonDecodeErr != nil {
		return Config{}, jsonDecodeErr
	}

	return config, nil
}

func saveConfigToFile(configPath string, config Config) error {
	configFile, err := os.Create(configPath)
	if err != nil {
		return err
	}

	defer configFile.Close()

	jsonData, encodeError := json.MarshalIndent(config, "", "  ")
	if encodeError != nil {
		return encodeError
	}

	_, writeError := configFile.Write(jsonData)
	if writeError != nil {
		return writeError
	}

	return nil
}

// newSampleConfig creates a new sample configuration that can be used
// as a starting point for more elaborate configurations.
func newSampleConfig() Config {
	return Config{
		Skill: AlexaSkill{
			AppID: "Your-Alexa-Skill-ID",
		},
		JenkinsAPI: &JenkinsAPI{
			URL:      "http://jenkins.example.com:8080",
			Username: "alexaops",
			APIToken: "8ebb23329c7a4575077462bd810030c16390dd7d",
		},
		Deployments: []DeploymentConfig{
			DeploymentConfig{
				Type: "Jenkins",
				Name: "wambo",
				Jenkins: &JenkinsDeployment{
					JobName: "wambo-shop-deploy",
				},
			},
		},
	}
}

package main

import (
	"os"
	"io/ioutil"
	"encoding/json"
)

// Config contains the configuration parameters for alexaops.
type Config struct {
	Skill AlexaSkill `json:"skill"`

	Parameters map[string]string `json:"parameters"`
	Projects   []Project `json:"projects"`
}

type AlexaSkill struct {
	// AppID is the application id from your
	// Amazon Dashboard (developer.amazon.com > Alexa > Alexa Skills Kit > Your Skill > ID)
	// Example: amzn1.ask.skill.abc12345-1111-dddd-aaaa-aaaabbbb3333
	AppID string `json:"appID"`
}

type Project struct {
	Name string `json:"name"`

	Parameters map[string]string `json:"parameters"`

	DeploymentIntend DeploymentIntendConfig `json:"deploymentIntend"`
	RestartIntend    RestartIntendConfig `json:"restartIntend"`
	VersionIntend    VersionIntendConfig `json:"versionIntend"`
	StatusIntend     StatusIntendConfig `json:"statusIntend"`
}

type BashCode struct {
	Code string `json:"code"`
}

type DeploymentIntendConfig struct {
	Deploy BashCode `json:"deploy"`
}

type RestartIntendConfig struct {
	Restart BashCode `json:"restart"`
}

type VersionIntendConfig struct {
	GetCurrentVersion   BashCode `json:"getCurrentVersion"`
	GetAvailableVersion BashCode `json:"getAvailableVersion"`
}

type StatusIntendConfig struct {
	GetStatus BashCode `json:"getStatus"`
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
		Parameters: map[string]string{
			"PROJECT_NAME": "SomeValue",
		},
		Projects: []Project{
			Project{
				Name: "ak7.io",
				Parameters: map[string]string{
					"APPLICATION_NAME": "SomeOtherValue",
				},
				DeploymentIntend: DeploymentIntendConfig{
					Deploy: BashCode{
						Code: "curl https://ak7.io/deploy/$PROJECT_NAME/$APPLICATION_NAME",
					},
				},
				RestartIntend: RestartIntendConfig{
					Restart: BashCode{
						Code: "curl https://ak7.io/restart/$PROJECT_NAME/$APPLICATION_NAME",
					},
				},
				VersionIntend: VersionIntendConfig{
					GetCurrentVersion: BashCode{
						Code: "curl https://ak7.io/version/$PROJECT_NAME/$APPLICATION_NAME",
					},
				},
				StatusIntend: StatusIntendConfig{
					GetStatus: BashCode{
						Code: "curl https://ak7.io/status/$PROJECT_NAME/$APPLICATION_NAME",
					},
				},
			},
		},
	}
}

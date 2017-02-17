package main

// Config contains the configuration parameters for alexaops.
type Config struct {
	// AppID is the application id from your
	// Amazon Dashboard (developer.amazon.com > Alexa > Alexa Skills Kit > Your Skill > ID)
	// Example: amzn1.ask.skill.abc12345-1111-dddd-aaaa-aaaabbbb3333
	AppID string `json:"appID"`
}

func readConfigFromFile(configPath string) (Config, error) {
	return Config{
		AppID: "amzn1.ask.skill.dbc35ee7-1d08-4a0d-a865-b2c2a5435690",
	}, nil
}

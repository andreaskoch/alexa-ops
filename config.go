package main

// Config contains the configuration parameters for alexaops.
type Config struct {
	// AppID is the application id from your
	// Amazon Dashboard (developer.amazon.com > Alexa > Alexa Skills Kit > Your Skill > ID)
	// Example: amzn1.ask.skill.abc12345-1111-dddd-aaaa-aaaabbbb3333
	AppID string `json:"appID"`
}

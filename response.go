package main

func createSpeechResponse(text string) ServiceResponse {
	response := ServiceResponse{}
	response.Version = "1.0"

	response.ResponseBody.Card.Type = "Simple"
	response.ResponseBody.Card.Title = "Deploy"
	response.ResponseBody.Card.Content = text

	response.ResponseBody.OutputSpeech.Type = "PlainText"
	response.ResponseBody.OutputSpeech.Text = text
	response.ResponseBody.ShouldEndSession = true
	
	return response
}

type ServiceResponse struct {
	Version string `json:"version"`
	SessionAttributes struct {
	} `json:"sessionAttributes,omitempty"`
	ResponseBody struct {
		OutputSpeech struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"outputSpeech,omitempty"`
		Card struct {
			Type    string `json:"type"`
			Title   string `json:"title"`
			Content string `json:"content"`
		} `json:"card,omitempty"`
		Reprompt         *Reprompt `json:"reprompt,omitempty"`
		ShouldEndSession bool `json:"shouldEndSession"`
	} `json:"response"`
}

type Reprompt struct {
	OutputSpeech *OutputSpeech `json:"outputSpeech"`
}

type OutputSpeech struct {
	Type string `json:"type"`
	Text string `json:"text"`
}


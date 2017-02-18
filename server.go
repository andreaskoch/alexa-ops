package main

import (
	"net/http"
	"fmt"
	"time"
	"encoding/json"
	"io/ioutil"
	"log"
)

func NewServer(listenAddress string, config Config) (Server, error) {
	return Server{
		listenAddress: listenAddress,
		config:        config,
	}, nil
}

type Server struct {
	listenAddress string
	config        Config
}

func (server *Server) Run() error {
	http.HandleFunc("/", server.handleRequest)
	return http.ListenAndServe(server.listenAddress, nil)
}

func (server *Server) handleRequest(w http.ResponseWriter, r *http.Request) {
	server.logRequest(r)

	serviceRequestModel, err := readServiceRequest(r)
	if err != nil {
		server.logError("Failed to decode request: %s", err.Error())
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		return
	}

	if _, err := requestMatchesApplicationID(serviceRequestModel, server.config); err != nil {
		server.logError("Application ID does not match: %s", err.Error())
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		return
	}

	response := createSpeechResponse(fmt.Sprintf("Deploying %s", serviceRequestModel.RequestBody.Intent.Slots.ApplicationName.Value))

	w.Header().Set("Content-Type", "application/json; charset=utf-8");
	if err := writeJSONResponse(w, response); err != nil {
		server.logError("Failed to write JSON response: %s", err.Error())
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (server *Server) logRequest(request *http.Request) {
	log.Printf("[Info] "+"%s %s %q %q", request.RemoteAddr, request.Method, request.URL.String(), request.UserAgent())
}

func (server *Server) logInfo(format string, p ...interface{}) {
	log.Printf("[Info] "+format, p...)
}

func (server *Server) logError(format string, p ...interface{}) {
	log.Printf("[Error] "+format, p...)
}

// requestMatchesApplicationID checks if the given ServiceRequest matches the application ID
// in the supplied configuration. If the application IDs match true and no error is returned.
// If the application IDs don't match false and an error is returned.
func requestMatchesApplicationID(request ServiceRequest, config Config) (bool, error) {
	isMatch := request.Session.Application.ApplicationID == config.AppID
	if isMatch {
		return true, nil
	}

	return false, fmt.Errorf("Given: %s, Required: %s", request.Session.Application.ApplicationID, config.AppID)
}

// readServiceRequest reads reads the ServiceRequest model from the given http request.
// Returns an error if the ServiceRequest could not be read.
func readServiceRequest(httpRequest *http.Request) (ServiceRequest, error) {
	var serviceRequest ServiceRequest
	body, readBodyError := ioutil.ReadAll(httpRequest.Body)
	if readBodyError != nil {
		return ServiceRequest{}, readBodyError
	}

	if unmarshalError := json.Unmarshal(body, &serviceRequest); unmarshalError != nil {
		return ServiceRequest{}, unmarshalError
	}

	return serviceRequest, nil
}

func writeJSONResponse(w http.ResponseWriter, serviceResponse ServiceResponse) error {

	jsonBody, err := json.MarshalIndent(serviceResponse, "", "  ")
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "%s", jsonBody)
	return nil
}

type ServiceRequest struct {
	Session struct {
		SessionID string `json:"sessionId"`
		Application struct {
			ApplicationID string `json:"applicationId"`
		} `json:"application"`
		Attributes struct {
		} `json:"attributes"`
		User struct {
			UserID string `json:"userId"`
		} `json:"user"`
		New bool `json:"new"`
	} `json:"session"`
	RequestBody struct {
		Type      string `json:"type"`
		RequestID string `json:"requestId"`
		Locale    string `json:"locale"`
		Timestamp time.Time `json:"timestamp"`
		Intent struct {
			Name string `json:"name"`
			Slots struct {
				ApplicationName struct {
					Name  string `json:"name"`
					Value string `json:"value"`
				} `json:"ApplicationName"`
			} `json:"slots"`
		} `json:"intent"`
	} `json:"request"`
	Version string `json:"version"`
}

func createSpeechResponse(text string) ServiceResponse {
	response := ServiceResponse{}
	response.Version = "1.0"

	response.ResponseBody.Card.Type = "Simple"
	response.ResponseBody.Card.Title = "Deploy"
	response.ResponseBody.Card.Content = text

	response.ResponseBody.OutputSpeech.Type = "PlainText"
	response.ResponseBody.OutputSpeech.Text = text
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

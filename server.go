package main

import (
	"net/http"
	"fmt"
	"time"
	"encoding/json"
	"io/ioutil"
	"log"
	"github.com/gorilla/mux"
	"github.com/andreaskoch/alexa-ops/files"
	"io"
)

func NewServer(listenAddress string, config Config, intendHandlerProvider intendHandlerProvider) (Server, error) {
	return Server{
		listenAddress:         listenAddress,
		config:                config,
		intendHandlerProvider: intendHandlerProvider,
	}, nil
}

type Server struct {
	listenAddress         string
	config                Config
	intendHandlerProvider intendHandlerProvider
}

func (server *Server) Run() error {
	r := mux.NewRouter()

	r.HandleFunc("/", server.intendHandler).Methods(http.MethodPost)

	r.HandleFunc("/", server.websiteHandler).Methods(http.MethodGet)
	r.HandleFunc("/index.html", server.websiteHandler).Methods(http.MethodGet)

	http.Handle("/", r)
	return http.ListenAndServe(server.listenAddress, nil)
}

func (server *Server) websiteHandler(w http.ResponseWriter, r *http.Request) {
	server.logRequest(r)

	indexHtml, err  := files.Open("website/index.html")
	defer indexHtml.Close()

	if err != nil {
		server.logError("Failed to decode request: %s", err.Error())
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		return
	}

	io.Copy(w, indexHtml)
}

func (server *Server) intendHandler(w http.ResponseWriter, r *http.Request) {
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

	intendHandler, intendHandlerErr := server.intendHandlerProvider.Get(serviceRequestModel.RequestBody.Intent.Name)
	if intendHandlerErr != nil {
		server.logError("No matching intend handler found: %s", intendHandlerErr.Error())
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		return
	}

	response, intendErr := intendHandler.Handle(serviceRequestModel)
	if intendErr != nil {
		server.logError("Failed to execute the %q intend handler: %s", intendHandler.Name(), intendErr.Error())
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
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
	isMatch := request.Session.Application.ApplicationID == config.Skill.AppID
	if isMatch {
		return true, nil
	}

	return false, fmt.Errorf("Given: %s, Required: %s", request.Session.Application.ApplicationID, config.Skill.AppID)
}

// readServiceRequest reads reads the ServiceRequest model from the given http request.
// Returns an error if the ServiceRequest could not be read.
func readServiceRequest(httpRequest *http.Request) (ServiceRequest, error) {
	if httpRequest.Method != http.MethodPost {
		return ServiceRequest{}, fmt.Errorf("Invalid HTTP method: %s", httpRequest.Method)
	}

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

package main

import (
	"net/http"
	"fmt"
	"time"
	"encoding/json"
	"io/ioutil"
)

func NewServer(listenAddress string, config Config) (Server, error) {
	return Server{
		listenAddress: listenAddress,
		appID: config.AppID,
		}, nil
}

type Server struct {
	listenAddress string
	appID string
}

func (server *Server) Run() error {

	http.HandleFunc("/", server.handleRequest)

	return http.ListenAndServe(server.listenAddress, nil)
}

func (server *Server) handleRequest(w http.ResponseWriter, r *http.Request) {
	serviceRequestModel, err := readServiceRequest(r)
	if err != nil {
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		return
	}

	fmt.Fprint(w, "%s", serviceRequestModel.Request.Intent.Slots.ApplicationName.Value)
}


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
	Request struct {
		Type string `json:"type"`
		RequestID string `json:"requestId"`
		Locale string `json:"locale"`
		Timestamp time.Time `json:"timestamp"`
		Intent struct {
			Name string `json:"name"`
			Slots struct {
				ApplicationName struct {
					Name string `json:"name"`
					Value string `json:"value"`
				} `json:"ApplicationName"`
			} `json:"slots"`
		} `json:"intent"`
	} `json:"request"`
	Version string `json:"version"`
}
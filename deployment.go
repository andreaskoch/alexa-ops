package main

import (
	"fmt"
	"github.com/andreaskoch/golang-jenkins"
	"net/url"
)

const DEPLOYMENT_STARTED = "deployment_of_%s_started"
const DEPLOYMENT_FAILED = "deployment_of_%s_failed"

func getJenkinsJobs() error {
	auth := &gojenkins.Auth{
		Username: "alexaops",
		ApiToken: "0395e8ffce4a42ec8912265ec70bffbe",
	}
	jenkins := gojenkins.NewJenkins(auth, "http://jenkins.002.io:50080")
	jobs, err := jenkins.GetJobs()
	if err != nil {
		return err
	}

	for _, job := range jobs {
		fmt.Println()
		fmt.Println(job.Name)
		fmt.Println(job.LastCompletedBuild.ChangeSet.Kind)
	}

	return nil
}

func newDeploymentIntendHandler(config Config) intendHandler {
	return &deploymentIntendHandler{
		config:    config,
		localizer: getDeploymentHandlerLocalizations(),
		handler:   &jenkinsDeploymentHandler{*config.JenkinsAPI, config.Deployments},
	}
}

// deploymentIntendHandler handles deployment intend requests.
type deploymentIntendHandler struct {
	config    Config
	localizer localizer
	handler   deploymentHandler
}

func (deployment *deploymentIntendHandler) Name() string {
	return "Deployment"
}

func (deployment *deploymentIntendHandler) Handle(request ServiceRequest) (ServiceResponse, error) {

	culture := request.RequestBody.Locale
	applicationName := request.RequestBody.Intent.Slots.ApplicationName.Value;
	deploymentError := deployment.handler.Deploy(applicationName)
	if deploymentError != nil {
		localizationDeploymentFailed, localizationError := deployment.localizer.Localize(DEPLOYMENT_FAILED, culture, applicationName)
		if localizationError != nil {
			return ServiceResponse{}, fmt.Errorf("Localization failed: %s", localizationError.Error())
		}

		return createSpeechResponse(localizationDeploymentFailed), nil
	}

	localizationDeploymentStarted, localizationError := deployment.localizer.Localize(DEPLOYMENT_STARTED, culture, applicationName)
	if localizationError != nil {
		return ServiceResponse{}, fmt.Errorf("Localization failed: %s", localizationError.Error())
	}

	return createSpeechResponse(localizationDeploymentStarted), nil
}

// getDeploymentHandlerLocalizations returns localizations for the deployment intend handler.
func getDeploymentHandlerLocalizations() localizer {
	localizations := make(localizations);

	// deployment started
	deploymentStarted := newLocalization(DEPLOYMENT_STARTED, defaultCulture, "The deployment of %s has been started")
	deploymentStarted.Add(cultureGerman, "Das Deployment von %s wurde gestartet")
	localizations[deploymentStarted.Key()] = deploymentStarted

	// deployment failed
	deploymentFailed := newLocalization(DEPLOYMENT_FAILED, defaultCulture, "The deployment of %s has failed")
	deploymentFailed.Add(cultureGerman, "Das Deployment von %s ist fehlgeschlagen")
	localizations[deploymentFailed.Key()] = deploymentFailed

	return &inMemoryLocalizer{localizations}
}

type deploymentHandler interface {
	Deploy(projectName string) error
}

type jenkinsDeploymentHandler struct {
	jenkinsConfig     JenkinsAPI
	deploymentConfigs []DeploymentConfig
}

func (deploymentHandler *jenkinsDeploymentHandler) Deploy(projectName string) error {
	deploymentConfig, err := getMatchingDeploymentConfig(projectName, deploymentHandler.deploymentConfigs)
	if err != nil {
		return err
	}

	auth := &gojenkins.Auth{
		Username: deploymentHandler.jenkinsConfig.Username,
		ApiToken: deploymentHandler.jenkinsConfig.APIToken,
	}
	jenkins := gojenkins.NewJenkins(auth, deploymentHandler.jenkinsConfig.URL)
	job, jobError := jenkins.GetJob(deploymentConfig.Jenkins.JobName)
	if jobError != nil {
		return jobError
	}

	buildResult := jenkins.Build(job, url.Values{})
	if buildResult != nil {
		return buildResult
	}

	return nil
}

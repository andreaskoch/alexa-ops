package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"github.com/satori/go.uuid"
	"os/exec"
)

const DEPLOYMENT_STARTED = "deployment_of_%s_started"
const DEPLOYMENT_FAILED = "deployment_of_%s_failed"

func newDeploymentIntendHandler(config Config) intendHandler {
	return &deploymentIntendHandler{
		config:    config,
		localizer: getDeploymentHandlerLocalizations(),
		handler:   &bashDeploymentHandler{config.Parameters, config.Projects},
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
		log.Println("Error", deploymentError.Error())

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

type bashDeploymentHandler struct {
	parameters map[string]string
	projects   []Project
}

func (deploymentHandler *bashDeploymentHandler) Deploy(projectName string) error {
	project, projectErr := getMatchingProject(projectName, deploymentHandler.projects)
	if projectErr != nil {
		return projectErr
	}

	scriptDirectory := os.TempDir()
	scriptPath := filepath.Join(scriptDirectory, fmt.Sprintf("%s.sh", uuid.NewV4().String()))
	scriptFile, err := os.Create(scriptPath)
	if err != nil {
		return err
	}

	defer func() {
		scriptFile.Close()
		os.Remove(scriptPath)
	}()

	scriptFile.WriteString(project.DeploymentIntend.Deploy.Code)

	deploymentCommand := exec.Command("/bin/bash", scriptPath)
	deploymentCommand = addEnvironmentVariables(deploymentCommand, deploymentHandler.parameters)
	deploymentCommand = addEnvironmentVariables(deploymentCommand, project.Parameters)

	if err := deploymentCommand.Run(); err != nil {
		return err
	}

	return nil
}

func addEnvironmentVariables(command *exec.Cmd, environmentVariables map[string]string) *exec.Cmd {
	for key, value := range environmentVariables {
		command.Env = append(command.Env, fmt.Sprintf("%s=%s", key, value))
	}

	return command
}

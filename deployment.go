package main

import "fmt"

const DEPLOYMENT_STARTED = "deployment_of_%s_started"

func newDeploymentIntendHandler(config Config) intendHandler {
	return &deploymentIntendHandler{
		config:    config,
		localizer: getDeploymentHandlerLocalizations(),
	}
}

// deploymentIntendHandler handles deployment intend requests.
type deploymentIntendHandler struct {
	config    Config
	localizer localizer
}

func (deployment *deploymentIntendHandler) Name() string {
	return "Deployment"
}

func (deployment *deploymentIntendHandler) Handle(request ServiceRequest) (ServiceResponse, error) {

	culture := request.RequestBody.Locale
	applicationName := request.RequestBody.Intent.Slots.ApplicationName.Value;

	deploymentStarted, localizationError := deployment.localizer.Localize(DEPLOYMENT_STARTED, culture, applicationName)
	if localizationError != nil {
		return ServiceResponse{}, fmt.Errorf("Localization failed: %s", localizationError.Error())
	}

	return createSpeechResponse(deploymentStarted), nil
}

// getDeploymentHandlerLocalizations returns localizations for the deployment intend handler.
func getDeploymentHandlerLocalizations() localizer {
	localizations := make(localizations);

	// deployment
	deploymentStarted := newLocalization(DEPLOYMENT_STARTED, defaultCulture, "The deployment of %s has been started")
	deploymentStarted.Add(cultureGerman, "Das Deployment von %s wurde gestartet")
	localizations[deploymentStarted.Key()] = deploymentStarted

	return &inMemoryLocalizer{localizations}
}

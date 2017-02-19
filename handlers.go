package main

import "fmt"

type intendHandler interface {
	// Name returns the name of the intend handler
	Name() string

	// Handle executes the intend with the given request parameters.
	// Returns the indent response if everything went well.
	// Otherwise it will return an error.
	Handle(request ServiceRequest) (ServiceResponse, error)
}

type intendHandlerProvider struct {
	handlers map[string]intendHandler
}

// Get returns the intend handler which matches the given name.
// Returns an error if no matching intend handler was found.
func (intendHandlerFactory *intendHandlerProvider) Get(intendName string) (intendHandler, error) {
	if handler, exists := intendHandlerFactory.handlers[intendName]; exists {
		return handler, nil
	}

	return nil, fmt.Errorf("Handler %q is unknown", intendName)
}

func newIntendHandlerProvider(config Config) intendHandlerProvider {
	handlers := make(map[string]intendHandler)

	// deployment
	deploymentHandler := newDeploymentIntendHandler(config)
	handlers[deploymentHandler.Name()] = deploymentHandler

	return intendHandlerProvider{handlers}
}

type genericIntendHandler struct {
	name    string
	execute func(request ServiceRequest) (ServiceResponse, error)
}

func (intendHandler *genericIntendHandler) Name() string {
	return intendHandler.name
}

func (intendHandler *genericIntendHandler) Execute(request ServiceRequest) (ServiceResponse, error) {
	return intendHandler.execute(request)
}

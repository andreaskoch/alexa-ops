package main

import (
	"fmt"
	"regexp"
	"strings"
)

func getMatchingDeploymentConfig(projectName string, deploymentConfigs []DeploymentConfig) (DeploymentConfig, error) {
	for _, deploymentConfig := range deploymentConfigs {

		if normalizeProjectName(deploymentConfig.Name) == normalizeProjectName(projectName) {
			return deploymentConfig, nil
		}
	}

	return DeploymentConfig{}, fmt.Errorf("No matching deployment config found for %q", projectName)
}

var whitespace = regexp.MustCompile(`[\s\.]`)

func normalizeProjectName(projectName string) string {
	unified := strings.ToLower(projectName)
	trimmed := strings.TrimSpace(unified)
	whitespaceRemoved := whitespace.ReplaceAllString(trimmed, "")
	return whitespaceRemoved

}
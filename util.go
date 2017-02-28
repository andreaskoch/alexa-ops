package main

import (
	"fmt"
	"regexp"
	"strings"
)

func getMatchingProject(projectName string, deploymentConfigs []Project) (Project, error) {
	for _, deploymentConfig := range deploymentConfigs {

		if normalizeProjectName(deploymentConfig.Name) == normalizeProjectName(projectName) {
			return deploymentConfig, nil
		}
	}

	return Project{}, fmt.Errorf("No matching project found for %q", projectName)
}

var whitespace = regexp.MustCompile(`[\s\.]`)

func normalizeProjectName(projectName string) string {
	unified := strings.ToLower(projectName)
	trimmed := strings.TrimSpace(unified)
	whitespaceRemoved := whitespace.ReplaceAllString(trimmed, "")
	return whitespaceRemoved

}
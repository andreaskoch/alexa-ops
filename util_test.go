package main

import "testing"

func Test_normalizeProjectName(t *testing.T) {
	inputs := []struct {
		input          string
		expectedResult string
	}{
		{"Andy K Docs", "andykdocs"},
		{"Andy K. Docs", "andykdocs"},
		{"Andy K Docs ", "andykdocs"},
		{" Andy K Docs", "andykdocs"},
		{" AndyK.Docs", "andykdocs"},
	}

	for _, input := range inputs {
		result := normalizeProjectName(input.input)

		if result != input.expectedResult {
			t.Fail()
			t.Logf("normalizeProjectName(%q) returned %q instead of %q", input.input, result, input.expectedResult)
		}
	}
}

func Test_getMatchingDeploymentConfig(t *testing.T) {
	inputs := []struct {
		input          string
		expectedResult string
	}{
		{"Andy K Docs", "andykdocs"},
		{"ak 7 io", "ak7.io"},
	}

	deployments := []DeploymentConfig{
		DeploymentConfig{
			Name: "andykdocs",
		},
		DeploymentConfig{
			Name: "la di da will not match",
		},
		DeploymentConfig{
			Name: "ak7.io",
		},
	}

	for _, input := range inputs {
		deploymentConfig, err := getMatchingDeploymentConfig(input.input, deployments)

		if err != nil {
			t.Fail()
			t.Logf("getMatchingDeploymentConfig(%q, ...) returned an error: %s", err.Error())
		}

		if deploymentConfig.Name != input.expectedResult {
			t.Fail()
			t.Logf("getMatchingDeploymentConfig(%q, ...) returned %q instead of %q", deploymentConfig.Name, input.expectedResult)
		}
	}
}

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

	deployments := []Project{
		Project{
			Name: "andykdocs",
		},
		Project{
			Name: "la di da will not match",
		},
		Project{
			Name: "ak7.io",
		},
	}

	for _, input := range inputs {
		deploymentConfig, err := getMatchingProject(input.input, deployments)

		if err != nil {
			t.Fail()
			t.Logf("getMatchingProject(%q, ...) returned an error: %s", err.Error())
		}

		if deploymentConfig.Name != input.expectedResult {
			t.Fail()
			t.Logf("getMatchingProject(%q, ...) returned %q instead of %q", deploymentConfig.Name, input.expectedResult)
		}
	}
}

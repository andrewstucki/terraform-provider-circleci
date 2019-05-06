package circleci

import (
	circleciapi "github.com/andrewstucki/terraform-provider-circleci/circleci/client"
	"github.com/cenkalti/backoff"
)

// ProviderClient is a thin commodity wrapper on top of circleciapi
type ProviderClient struct {
	client       *circleciapi.Client
	vcsType      string
	organization string
}

// NewConfig initialize circleci API client and returns a new config object
func NewConfig(token, vscType, organization string) *ProviderClient {
	return &ProviderClient{
		client: &circleciapi.Client{
			Token: token,
		},
		vcsType:      vscType,
		organization: organization,
	}
}

// GetEnvVar get the environment variable with given name
// It returns an empty structure if no environment variable exists with that name
func (pv *ProviderClient) GetEnvVar(projectName, envVarName string) (*circleciapi.EnvVar, error) {
	return pv.client.GetEnvVar(pv.vcsType, pv.organization, projectName, envVarName)
}

// EnvVarExists check if environment variable exists with given name
func (pv *ProviderClient) EnvVarExists(projectName, envVarName string) (bool, error) {
	envVar, err := pv.client.GetEnvVar(pv.vcsType, pv.organization, projectName, envVarName)
	if err != nil {
		return false, err
	}
	return bool(envVar.Name != ""), nil
}

// AddEnvVar create an environment variable with given name and value
func (pv *ProviderClient) AddEnvVar(projectName, envVarName, envVarValue string) (*circleciapi.EnvVar, error) {
	return pv.client.AddEnvVar(pv.vcsType, pv.organization, projectName, envVarName, envVarValue)
}

// DeleteEnvVar delete the environment variable with given name
func (pv *ProviderClient) DeleteEnvVar(projectName, envVarName string) error {
	retry := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 3)

	return backoff.Retry(func() error {
		return pv.client.DeleteEnvVar(pv.vcsType, pv.organization, projectName, envVarName)
	}, retry)
}

// GetProject reads the project with given name
func (pv *ProviderClient) GetProject(projectName string) (*circleciapi.Project, error) {
	retry := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 3)

	var err error
	var project *circleciapi.Project
	err = backoff.Retry(func() error {
		project, err = pv.client.GetProject(pv.organization, projectName)
		return err
	}, retry)
	return project, err
}

// EnableProject enables the project with given name
func (pv *ProviderClient) EnableProject(projectName string) error {
	retry := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 3)

	return backoff.Retry(func() error {
		return pv.client.EnableProject(pv.vcsType, pv.organization, projectName)
	}, retry)
}

// FollowProject follows the project with given name
func (pv *ProviderClient) FollowProject(projectName string) (*circleciapi.Project, error) {
	retry := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 3)

	var err error
	var project *circleciapi.Project
	err = backoff.Retry(func() error {
		project, err = pv.client.FollowProject(pv.vcsType, pv.organization, projectName)
		return err
	}, retry)
	return project, err
}

// DisableProject disables the project with given name
func (pv *ProviderClient) DisableProject(projectName string) error {
	retry := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 3)

	return backoff.Retry(func() error {
		return pv.client.DisableProject(pv.vcsType, pv.organization, projectName)
	}, retry)
}

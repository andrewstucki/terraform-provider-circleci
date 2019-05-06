package circleci

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var testProvider *schema.Provider
var testProviders map[string]terraform.ResourceProvider

func init() {
	testProvider = Provider().(*schema.Provider)
	testProviders = map[string]terraform.ResourceProvider{
		"circleci": testProvider,
	}
}

func testPreCheck(t *testing.T) {
	if v := os.Getenv("CIRCLECI_TOKEN"); v == "" {
		t.Fatal("CIRCLECI_TOKEN must be set for acceptance tests")
	}

	if v := os.Getenv("CIRCLECI_VCS_TYPE"); v == "" {
		t.Fatal("CIRCLECI_VCS_TYPE must be set for acceptance tests")
	}

	if v := os.Getenv("CIRCLECI_ORGANIZATION"); v == "" {
		t.Fatal("CIRCLECI_ORGANIZATION must be set for acceptance tests")
	}

	if v := os.Getenv("CIRCLECI_PROJECT"); v == "" {
		t.Fatal("CIRCLECI_PROJECT must be set for acceptance tests")
	}
}

func testCircleCIEnvironmentVariableConfig(project, name, value string) string {
	return fmt.Sprintf(`
resource "circleci_project" "%[1]s" {
  repo = "%[1]s"
}

resource "circleci_environment_variable" "%[2]s" {
	project = "${circleci_project.%[1]s.id}"
  name    = "%[2]s"
  value   = "%[3]s"
}`, project, name, value)
}

func testCircleCIEnvironmentVariableConfigIdentical(project, name, value string) string {
	return fmt.Sprintf(`
resource "circleci_project" "%[1]s" {
  repo = "%[1]s"
}

resource "circleci_environment_variable" "%[2]s" {
  project = "${circleci_project.%[1]s.id}"
  name    = "%[2]s"
  value   = "%[3]s"
}

resource "circleci_environment_variable" "%[2]s_2" {
	project = "${circleci_project.%[1]s.id}"
  name    = "%[2]s"
  value   = "%[3]s"
}`, project, name, value)
}

func testCircleCICheckDestroy(s *terraform.State) error {
	providerClient := testProvider.Meta().(*ProviderClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "circleci_environment_variable" {
			continue
		}

		project, err := providerClient.GetProject(rs.Primary.Attributes["project"])
		if err != nil {
			return err
		}
		if project != nil {
			return fmt.Errorf("Project should have been destroyed")
		}

		envVar, err := providerClient.GetEnvVar(rs.Primary.Attributes["project"], rs.Primary.Attributes["name"])
		if err != nil {
			return err
		}

		if envVar.Name != "" {
			return errors.New("Environment variable should have been destroyed")
		}
	}

	return nil
}

func TestCircleCICreateThenUpdate(t *testing.T) {
	project := os.Getenv("CIRCLECI_PROJECT")
	envName := "TEST_" + acctest.RandString(8)

	resourceName := "circleci_environment_variable." + envName

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testPreCheck(t)
		},
		Providers:    testProviders,
		CheckDestroy: testCircleCICheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCircleCIEnvironmentVariableConfig(project, envName, "value-for-the-test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "project", project),
					resource.TestCheckResourceAttr(resourceName, "name", envName),
					resource.TestCheckResourceAttr(resourceName, "value", "value-for-the-test"),
				),
			},
			{
				Config: testCircleCIEnvironmentVariableConfig(project, envName, "value-for-the-test-again"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "project", project),
					resource.TestCheckResourceAttr(resourceName, "name", envName),
					resource.TestCheckResourceAttr(resourceName, "value", "value-for-the-test-again"),
				),
			},
		},
	})
}

func TestCircleCICreateAlreadyExists(t *testing.T) {
	project := os.Getenv("CIRCLECI_PROJECT")
	envName := "TEST_" + acctest.RandString(8)
	envValue := acctest.RandString(8)

	resourceName := "circleci_environment_variable." + envName

	resource.Test(t, resource.TestCase{
		Providers: testProviders,
		PreCheck: func() {
			testPreCheck(t)
		},
		CheckDestroy: testCircleCICheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCircleCIEnvironmentVariableConfig(project, envName, envValue),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "project", project),
					resource.TestCheckResourceAttr(resourceName, "name", envName),
					resource.TestCheckResourceAttr(resourceName, "value", envValue),
				),
			},
			{
				Config:      testCircleCIEnvironmentVariableConfigIdentical(project, envName, envValue),
				ExpectError: regexp.MustCompile("already exists"),
			},
		},
	})
}

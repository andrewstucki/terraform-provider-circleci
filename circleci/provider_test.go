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

const (
	testPrivateKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEA4B1zFe/o1y06ZXh8F1Sw8RSs42Zt2icrduRXW6SxRyh3LJ5j
wEHqM26rzDzIfsYEicLno09+OxDCzQ4V3DFN/z4M1ydiZacVTn3tTjX4ir9w08Jc
U8eR+70kFPS2QYMS5m0+XPzStW4JFr3j4bIzWIar5u0iYjp6vJN2SkvFnoRakrg+
jfd2wwB6/HessWeMkELy81fJsL0AiJAiUrZSw+aujWVI6EISquU9Ru4iV1z9NvPI
kPHvNCFYqaKoWJETcN82iOjvmUF2CizzAeXLeyXXGT/vY06nnPdfgMWXNmvKjbdd
DO+9/xLUzhUugvvRClfRtcmMxIv3dtUmDGVnBwIDAQABAoIBAQCWOQdAaByBx05C
X92F4f3syvgMQUdXGDRZMDuiMSWnVed0IAMrBsPOj9tWPlQCLgcytFOCMzGgs95v
hoZ+nwnyEgyXV03sZ2+vulcLur+LeUFOoBQ1ECu5OxHpfhKAnTRZAlbwC4PNmyE6
fjQ2v7UNHNAzLTaD80D8EDgVGu8vDbUCTIKEbHTUs+L1SHRXTAKXeC2MMiLMfc2m
91cooEqgMVjWxKIHSKfLhaVv/XZSQMz+dJuIKcg0ZA3GW32HQQdxo539wanprYjq
8sZYZwOm65i/O5nTDj6C615HzijO9PIKABIx54/WYxDL6oSL5xT50voWZVqAheSV
lMT/z4HRAoGBAPw1KvjKGqjtYqP6sRGHx3lLNbl9dFBaRAMwhtX2zNu6osS4nM5O
bg+IvhF9Gt09Z4zMHkfOYfenBCSs6UyQwWUa0Qi93AcxQKBALCdN0+bcMc4q2duO
Ur/GAjUHhNUGQJ5v2wFLhYPt1AD4YVyI6R1eUBwJm4IFriFLOkBAMbNZAoGBAON8
JMX47D7jf+Bb9gt6pmYbJxTTfL1QGvpp+UKRokSVfqBKznQQFVhHe2FAtWbYQ3+J
1U+MIYJAqRTmhPgc4QyJAW+kiDc0cyg/ILG3l+7awqRoCTD10nHjbe4HZxa9a0mW
SYUFU0GIZZRxQZ7C9rfuuKW25+bZLn7ezp00tIFfAoGAU6WuouUlAnH5DTnQEGhg
GDKBlwus0BmgBQ7LKZu5RgcYhPZVy3bnue84WsSLbGU5OtFYGaixhVm3XhKbLfG0
sru6KJQPrbMAJCYkfsSpSyAsxJwhtVf2yfP6N2xO+fgg5mtiz4MkvSTb85ZtdCtU
ZZEqMKJfGTiZECHLKBQiZ8ECgYAlWHEVCyuFm4WXyKEY+1ar9pMw6RNWZPs41wLz
ucLg7YXvPLit9yH57ypDKgNd0e0q1+7r8z5hCsp3QuzbaqpLi4Zv1JwELBknp01v
v4syzDkeEnJH1mNpDQQ0CoUTB5/AYerJ6rjjTkgW2Y0DSlCEm602j1N843StoVhc
GJX1kwKBgHDVTxhBTdDT4VN0bDndKm6nKjWlhF6PifpeQZ0noEMlu0a/TrULmOHj
7uhyL6+QR4MMeCon9EcdunElEq4gDkVBHkttbl9YnfWlpBb4eLGrd6kIfjPmH9Ks
k8ZO3kUnDGr3/IpYOe61RE0/10f46fm3PQ7+p+o0e4DMOzLQ2Dy7
-----END RSA PRIVATE KEY-----`
)

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

resource "circleci_ssh_key" "%[1]s" {
	project = "${circleci_project.%[1]s.id}"
	hostname = "github.com"
	private_key = <<EOF
%[4]s
EOF
}

resource "circleci_environment_variable" "%[2]s" {
	project = "${circleci_project.%[1]s.id}"
  name    = "%[2]s"
  value   = "%[3]s"
}`, project, name, value, testPrivateKey)
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

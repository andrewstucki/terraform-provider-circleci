package circleci

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceCircleCIProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceCircleCIProjectCreate,
		Read:   resourceCircleCIProjectRead,
		Delete: resourceCircleCIProjectDelete,
		Exists: resourceCircleCIProjectExists,

		Schema: map[string]*schema.Schema{
			"repo": {
				Type:        schema.TypeString,
				Description: "The name of the CircleCI project to enable",
				Required:    true,
				ForceNew:    true,
			},
		},

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceCircleCIProjectCreate(d *schema.ResourceData, m interface{}) error {
	providerClient := m.(*ProviderClient)

	name := d.Get("repo").(string)

	_, err := providerClient.FollowProject(name)
	if err != nil {
		return err
	}

	err = providerClient.EnableProject(name)
	if err != nil {
		return err
	}

	d.SetId(name)

	return nil
}

func resourceCircleCIProjectRead(d *schema.ResourceData, m interface{}) error {
	providerClient := m.(*ProviderClient)

	name := d.Get("repo").(string)

	_, err := providerClient.GetProject(name)
	if err != nil {
		return err
	}

	d.SetId(name)

	return nil
}

func resourceCircleCIProjectDelete(d *schema.ResourceData, m interface{}) error {
	providerClient := m.(*ProviderClient)

	name := d.Get("repo").(string)

	err := providerClient.DisableProject(name)
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}

func resourceCircleCIProjectExists(d *schema.ResourceData, m interface{}) (bool, error) {
	providerClient := m.(*ProviderClient)

	name := d.Get("repo").(string)

	project, err := providerClient.GetProject(name)
	if err != nil {
		return false, err
	}

	return bool(project != nil), nil
}

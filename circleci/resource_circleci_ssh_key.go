package circleci

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"golang.org/x/crypto/ssh"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceCircleCISSHKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceCircleCISSHKeyCreate,
		Read:   resourceCircleCISSHKeyRead,
		Delete: resourceCircleCISSHKeyDelete,

		Schema: map[string]*schema.Schema{
			"project": {
				Type:        schema.TypeString,
				Description: "The name of the CircleCI project to which you want to add the SSH key",
				Required:    true,
				ForceNew:    true,
			},
			"hostname": {
				Type:        schema.TypeString,
				Description: "The hostname where we want to use the SSH key",
				Required:    true,
				ForceNew:    true,
			},
			"private_key": {
				Type:        schema.TypeString,
				Description: "The SSH private key",
				Required:    true,
				ForceNew:    true,
				Sensitive:   true,
			},
			"fingerprint": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceCircleCISSHKeyCreate(d *schema.ResourceData, m interface{}) error {
	providerClient := m.(*ProviderClient)

	name := d.Get("project").(string)
	hostname := d.Get("hostname").(string)
	privateKey := d.Get("private_key").(string)

	block, rest := pem.Decode([]byte(privateKey))
	if len(rest) != 0 {
		return fmt.Errorf("found %d stray bytes in private key", len(rest))
	}

	privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return err
	}

	pubKey, err := ssh.NewPublicKey(&privKey.PublicKey)
	if err != nil {
		return err
	}

	err = providerClient.AddSSHKey(name, hostname, privateKey)
	if err != nil {
		return err
	}

	fingerprint := ssh.FingerprintLegacyMD5(pubKey)

	d.SetId(fmt.Sprintf("%s|%s|%s", name, hostname, fingerprint))
	d.Set("fingerprint", fingerprint)

	return nil
}

func resourceCircleCISSHKeyRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceCircleCISSHKeyDelete(d *schema.ResourceData, m interface{}) error {
	providerClient := m.(*ProviderClient)

	name := d.Get("project").(string)
	hostname := d.Get("hostname").(string)
	fingerprint := d.Get("fingerprint").(string)

	err := providerClient.DeleteSSHKey(name, hostname, fingerprint)
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}

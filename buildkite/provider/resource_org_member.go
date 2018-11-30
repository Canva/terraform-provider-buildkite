package provider

import (
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/pkg/errors"
	"log"

	"github.com/hashicorp/terraform/helper/schema"

	"github.com/saymedia/terraform-buildkite/buildkite/client"
)

var (
	ValidOrganizationMemberRole = []string{client.OrganizationMemberRoleMember, client.OrganizationMemberRoleAdmin}
)

func resourceOrgMember() *schema.Resource {
	return &schema.Resource{
		Create: CreateOrganizationMember,
		Read:   ReadOrganizationMember,
		Update: UpdateOrganizationMember,
		Delete: DeleteOrganizationMember,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"member_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"role": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(ValidOrganizationMemberRole, false),
			},
			"user_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_email": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func CreateOrganizationMember(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[TRACE] CreateOrganizationMember")
	return errors.New("org member cannot be created")
}

func ReadOrganizationMember(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[TRACE] ReadOrganizationMember")

	buildkiteClient := meta.(*client.Client)
	uuid := d.Id()

	orgMember, err := buildkiteClient.GetOrganizationMember(uuid)
	if err != nil {
		if _, ok := err.(*client.NotFound); ok {
			d.SetId("")
			return nil
		}
		return err
	}

	return updateOrgMemberFromAPI(d, orgMember)
}

func UpdateOrganizationMember(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[TRACE] UpdateOrganizationMember")

	buildkiteClient := meta.(*client.Client)

	orgMember := prepareOrgMemberRequestPayload(d)

	res, err := buildkiteClient.UpdateOrganizationMember(orgMember)
	if err != nil {
		return err
	}

	return updateOrgMemberFromAPI(d, res)
}

func DeleteOrganizationMember(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[TRACE] DeleteOrganizationMember")

	buildkiteClient := meta.(*client.Client)
	id := d.Get("member_id").(string)

	return buildkiteClient.DeleteOrganizationMember(id)
}

func updateOrgMemberFromAPI(d *schema.ResourceData, t *client.OrganizationMember) error {
	d.SetId(t.UUID)
	log.Printf("[INFO] buildkite: Pipeline ID: %s", d.Id())

	d.Set("member_id", t.Id)
	d.Set("uuid", t.UUID)
	d.Set("role", t.Role)
	d.Set("created_at", t.CreatedAt)
	d.Set("user_id", t.User.Id)
	d.Set("user_name", t.User.Name)
	d.Set("user_email", t.User.Email)

	return nil
}

func prepareOrgMemberRequestPayload(d *schema.ResourceData) *client.OrganizationMember {
	req := &client.OrganizationMember{}

	req.Id = d.Get("member_id").(string)
	req.UUID = d.Get("uuid").(string)
	req.Role = d.Get("role").(string)
	req.CreatedAt = d.Get("created_at").(string)
	req.User = client.User{
		Id:    d.Get("user_id").(string),
		Name:  d.Get("user_name").(string),
		Email: d.Get("user_email").(string),
	}

	return req
}

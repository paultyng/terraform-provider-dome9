package dome9

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
)

func dataSourceImageAssurancePolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceImageAssurancePolicyRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"target_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ruleset_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"notification_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"admission_control_action": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"admission_control_unscanned_action": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceImageAssurancePolicyRead(d *schema.ResourceData, meta interface{}) error {
	d9Client := meta.(*Client)

	policyID := d.Get("id").(string)
	log.Printf("Getting data for Image Assurance Policy id: %s\n", policyID)

	resp, _, err := d9Client.imageAssurancePolicy.Get(policyID)
	if err != nil {
		return err
	}

	d.SetId(resp.ID)
	_ = d.Set("target_id", resp.TargetId)
	_ = d.Set("target_type", resp.TargetType)
	_ = d.Set("ruleset_id", resp.RulesetId)
	_ = d.Set("admission_control_action", resp.AdmissionControllerAction)
	_ = d.Set("admission_control_unscanned_action", resp.AdmissionControlUnScannedAction)
	if err := d.Set("notification_ids", resp.NotificationIds); err != nil {
		return err
	}

	return nil
}

package main

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/pactflow/terraform/client"
)

func pacticipants() *schema.Resource {
	return &schema.Resource{
		Read: pacticipantsRead,
		Schema: map[string]*schema.Schema{
			"pacticipants": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"repository_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"main_branch": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"display_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func pacticipantsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*client.Client)

	res, err := client.ListPacticipants()
	if err != nil {
		return fmt.Errorf("error reading pacticipants: %w", err)
	}

	items := make([]map[string]interface{}, len(res.Embedded.Items))
	for i, p := range res.Embedded.Items {
		items[i] = map[string]interface{}{
			"name":           p.Name,
			"repository_url": p.RepositoryURL,
			"main_branch":    p.MainBranch,
			"display_name":   p.DisplayName,
		}
	}

	d.SetId("pacticipants")

	return d.Set("pacticipants", items)
}

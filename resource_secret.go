package main

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pactflow/terraform/broker"
	"github.com/pactflow/terraform/client"
)

func secret() *schema.Resource {
	return &schema.Resource{
		CreateContext: secretCreate,
		UpdateContext: secretUpdate,
		ReadContext:   secretRead,
		DeleteContext: secretDelete,
		Importer:      &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "A short name of the secret (alphanumeric characters only)",
				ValidateFunc: validateName,
			},
			"description": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A longer description for the secret",
			},
			"value": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				ExactlyOneOf: []string{"value", "value_wo"},
				Description:  "The actual secret. The value is stored in the Terraform state; prefer `value_wo` with Terraform 1.11+",
			},
			"value_wo": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				WriteOnly:    true,
				ExactlyOneOf: []string{"value", "value_wo"},
				RequiredWith: []string{"value_wo_version"},
				Description:  "Write-only variant of `value`: the secret is sent to the broker but never stored in the Terraform plan or state. Requires Terraform 1.11+. Must be used with `value_wo_version`",
			},
			"value_wo_version": {
				Type:         schema.TypeInt,
				Optional:     true,
				RequiredWith: []string{"value_wo"},
				Description:  "Version of `value_wo`. Terraform cannot detect changes to write-only values, so change (e.g. increment) this number to push a new `value_wo` to the broker",
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The UUID of secret",
			},
			"team": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The team this secret should be associated with (uuid). Leave empty for a non-team secret",
			},
		},
	}
}

func validateName(val interface{}, key string) (warns []string, errs []error) {
	v := val.(string)
	if matched, _ := regexp.MatchString(`[^a-zA-z0-9].*`, v); matched {
		errs = append(errs, fmt.Errorf("%q must be a string containing alphanumeric letters, got: %s", key, v))
	}
	return
}

// parseSecret builds a broker.Secret from the resource configuration.
// The secret value comes from `value`, or, when that is not set, from the
// write-only `value_wo` attribute (only available in the raw config).
func parseSecret(d *schema.ResourceData) (broker.Secret, diag.Diagnostics) {
	secret := broker.Secret{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Value:       d.Get("value").(string),
		TeamUUID:    d.Get("team").(string),
	}

	if secret.Value == "" {
		value, diags := rawConfigString(d, cty.GetAttrPath("value_wo"))
		if diags.HasError() {
			return secret, diags
		}
		secret.Value = value
	}

	// Existing secret?
	if d.Id() != "" {
		secret.UUID = d.Id()
	}

	return secret, nil
}

func secretCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.Client)
	secret, diags := parseSecret(d)
	if diags.HasError() {
		return diags
	}
	log.Println("[DEBUG] creating secret", secret.Name)

	res, err := client.CreateSecret(secret)
	if err != nil {
		return diag.FromErr(err)
	}

	items := strings.Split(res.Links["self"].Href, "/")
	id := items[len(items)-1]
	d.SetId(id)
	secret.UUID = id

	return diag.FromErr(setSecretState(d, secret))
}

func secretUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.Client)
	secret, diags := parseSecret(d)
	if diags.HasError() {
		return diags
	}

	log.Println("[DEBUG] updating secret", secret.UUID)

	if _, err := client.UpdateSecret(secret); err != nil {
		return diag.FromErr(err)
	}

	return diag.FromErr(setSecretState(d, secret))
}

func secretRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	httpClient := meta.(*client.Client)

	secret, err := httpClient.ReadSecret(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// The UUID is not part of the response body (json:"-"), it is the resource ID
	s := secret.Secret
	s.UUID = d.Id()

	return diag.FromErr(setSecretState(d, s))
}

func secretDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.Client)

	log.Println("[DEBUG] deleting secret", d.Id())

	if err := client.DeleteSecret(broker.Secret{UUID: d.Id()}); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

// setSecretState stores the non-sensitive secret attributes in the state.
//
// The broker never returns the secret value, so `value` is intentionally not
// touched here: Terraform keeps the configured value (or the value from the
// previous state), and `value_wo` is write-only and never persisted.
func setSecretState(d *schema.ResourceData, secret broker.Secret) error {
	log.Printf("[DEBUG] setting secret state for %s\n", secret.UUID)

	if err := d.Set("name", secret.Name); err != nil {
		return err
	}
	if err := d.Set("uuid", secret.UUID); err != nil {
		return err
	}
	if err := d.Set("description", secret.Description); err != nil {
		return err
	}
	return d.Set("team", secret.TeamUUID)
}

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"reflect"
	"regexp"
	"sort"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/mapstructure"
	"github.com/pactflow/terraform/broker"
	"github.com/pactflow/terraform/client"
)

var allowedEvents = []string{
	"contract_content_changed",
	"contract_published",
	"provider_verification_failed",
	"provider_verification_published",
	"provider_verification_succeeded",
	"contract_requiring_verification_published",
}

// NOTE: SDK v2 does not support TypeMap with a *schema.Resource Elem. The map
// only ever carries a single "name" key, so model it as a map of strings.
var pacticipantType = &schema.Schema{
	Type:        schema.TypeMap,
	Optional:    true,
	Computed:    true,
	ForceNew:    true,
	Description: "The pacticipant this webhook applies to, e.g. { name = \"my-app\" }",
	Elem: &schema.Schema{
		Type: schema.TypeString,
	},
}

var eventsType = &schema.Schema{
	Type:     schema.TypeSet,
	Optional: true,
	Elem: &schema.Schema{
		Type:         schema.TypeString,
		ValidateFunc: validateEvents,
	},
}

var requestType = &schema.Schema{
	Type:     schema.TypeList, // Terraform hack for complex objects
	MaxItems: 1,
	Required: true,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateURL,
				Description:  "A valid URL to send the webhook request to",
			},
			"method": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateMethod,
				Description:  "The HTTP method to use with the request",
			},
			"username": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"request.0.username_wo"},
				Description:   "An optional (basic auth) username to send with the request",
			},
			"username_wo": {
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				WriteOnly:     true,
				ConflictsWith: []string{"request.0.username"},
				RequiredWith:  []string{"request.0.username_wo_version"},
				Description:   "Write-only variant of `username`: sent to the broker but never stored in the Terraform plan or state. Requires Terraform 1.11+. Must be used with `username_wo_version`",
			},
			"username_wo_version": {
				Type:         schema.TypeInt,
				Optional:     true,
				RequiredWith: []string{"request.0.username_wo"},
				Description:  "Version of `username_wo`. Terraform cannot detect changes to write-only values, so change (e.g. increment) this number to push a new `username_wo` to the broker",
			},
			"password": {
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				ConflictsWith: []string{"request.0.password_wo"},
				Description:   "An optional (basic auth) password to send with the request. The value is stored in the Terraform state; prefer `password_wo` with Terraform 1.11+",
			},
			"password_wo": {
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				WriteOnly:     true,
				ConflictsWith: []string{"request.0.password"},
				RequiredWith:  []string{"request.0.password_wo_version"},
				Description:   "Write-only variant of `password`: sent to the broker but never stored in the Terraform plan or state. Requires Terraform 1.11+. Must be used with `password_wo_version`",
			},
			"password_wo_version": {
				Type:         schema.TypeInt,
				Optional:     true,
				RequiredWith: []string{"request.0.password_wo"},
				Description:  "Version of `password_wo`. Terraform cannot detect changes to write-only values, so change (e.g. increment) this number to push a new `password_wo` to the broker",
			},
			"headers": {
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Request headers to send with the request",
			},
			"headers_wo": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				WriteOnly:    true,
				RequiredWith: []string{"request.0.headers_wo_version"},
				ValidateFunc: validateHeadersJSON,
				Description:  "Write-only request headers, as a JSON object of strings (e.g. `jsonencode({ Authorization = \"Bearer ...\" })`). They are merged with `headers`, sent to the broker but never stored in the Terraform plan or state. Requires Terraform 1.11+. Must be used with `headers_wo_version`",
			},
			"headers_wo_version": {
				Type:         schema.TypeInt,
				Optional:     true,
				RequiredWith: []string{"request.0.headers_wo"},
				Description:  "Version of `headers_wo`. Terraform cannot detect changes to write-only values, so change (e.g. increment) this number to push new `headers_wo` to the broker",
			},
			"body": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "A request body to send with the request",
				DiffSuppressFunc: ignoreJSONFormatting,
			},
		},
	},
}

func stringContains(s []string, searchterm string) bool {
	sort.Strings(s)
	i := sort.SearchStrings(s, searchterm)
	return i < len(s) && s[i] == searchterm
}

func validateEvents(val interface{}, key string) (warns []string, errs []error) {
	v := val.(string)
	if !stringContains(allowedEvents, v) {
		errs = append(errs, fmt.Errorf("%q must be one of the allowed events %v, got %v", key, allowedEvents, v))
	}
	return
}

func validateURL(val interface{}, key string) (warns []string, errs []error) {
	v := val.(string)
	_, err := url.ParseRequestURI(v)

	if err != nil {
		errs = append(errs, fmt.Errorf("%q must be a valid URL, got: %v", key, err))
	}
	return
}

func validateMethod(val interface{}, key string) (warns []string, errs []error) {
	v := val.(string)
	if matched, _ := regexp.MatchString(`^(GET|PUT|PATCH|POST|DELETE)$`, v); !matched {
		errs = append(errs, fmt.Errorf("%q must one of the following HTTP Verbs 'GET, PUT, PATCH, POST, DELETE', got: %s", key, v))
	}
	return
}

func validateHeadersJSON(val interface{}, key string) (warns []string, errs []error) {
	if _, err := parseHeadersJSON(val.(string)); err != nil {
		errs = append(errs, fmt.Errorf("%q %v", key, err))
	}
	return
}

// parseHeadersJSON parses a JSON object of strings (e.g. the output of jsonencode) into headers.
func parseHeadersJSON(s string) (map[string]string, error) {
	headers := map[string]string{}
	if s == "" {
		return headers, nil
	}
	if err := json.Unmarshal([]byte(s), &headers); err != nil {
		return nil, fmt.Errorf("must be a JSON object of strings: %v", err)
	}
	return headers, nil
}

func webhook() *schema.Resource {
	return &schema.Resource{
		CreateContext: webhookCreate,
		UpdateContext: webhookUpdate,
		ReadContext:   webhookRead,
		DeleteContext: webhookDelete,
		Importer:      &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
		Schema: map[string]*schema.Schema{
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"webhook_provider": pacticipantType,
			"webhook_consumer": pacticipantType,
			"request":          requestType,
			"events":           eventsType,
			"enabled": {
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},
			"team": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The team this webhook should be associated with (uuid). Leave empty for a non-team Webhook",
			},
		},
	}
}

func parseWebhook(d *schema.ResourceData, meta interface{}) (broker.Webhook, error) {
	request := new(broker.Request)
	webhook := &broker.Webhook{
		Enabled: true,
		Events:  []broker.WebhookEvent{},
	}

	log.Printf("[DEBUG] create or update webhook with data %+v \n", d)

	webhook.Description = d.Get("description").(string)

	// Team
	if team, ok := d.GetOk("team"); ok {
		if team != "" {
			webhook.TeamUUID = team.(string)
		}
	}

	// Provider
	if rawProvider, ok := d.GetOk("webhook_provider"); ok {
		provider := new(broker.Pacticipant)
		log.Printf("[DEBUG] raw provider %+v \n", rawProvider)
		err := mapstructure.Decode(rawProvider, provider)
		if err != nil {
			log.Println("[ERROR] error decoding webhook config: webhook_provider", err)
			return *webhook, err
		}

		if provider.Name != "" {
			webhook.Provider = provider
		}
	}

	// Consumer
	if rawConsumer, ok := d.GetOk("webhook_consumer"); ok {
		consumer := new(broker.Pacticipant)
		log.Printf("[DEBUG] raw consumer %+v \n", rawConsumer)
		err := mapstructure.Decode(rawConsumer, consumer)
		if err != nil {
			log.Println("[ERROR] error decoding webhook config: webhook_consumer", err)
			return *webhook, err
		}

		if consumer.Name != "" {
			webhook.Consumer = consumer
		}
	}

	// Events
	if eventsRaw, ok := d.GetOk("events"); ok {
		events := eventsRaw.(*schema.Set)
		for _, event := range ExpandStringSet(events) {
			log.Printf("[DEBUG]event item %+v\n", event)
			webhook.Events = append(webhook.Events, broker.WebhookEvent{
				Name: event,
			})
		}
	}

	// Request
	log.Println("[DEBUG] checking request")
	if rawRequest, ok := d.GetOk("request"); ok {
		log.Printf("[DEBUG] have raw request of %+v \n", rawRequest)

		rawRequestList := rawRequest.([]interface{})
		requestMap := rawRequestList[0].(map[string]interface{})
		log.Printf("[DEBUG] have converted request %+v \n", requestMap)

		// Method
		if method, ok := requestMap["method"]; ok {
			request.Method = method.(string)
		}

		// Username (either the regular attribute or the write-only one)
		if username, ok := requestMap["username"]; ok {
			request.Username = username.(string)
		}
		if request.Username == "" {
			username, diags := rawConfigString(d, cty.GetAttrPath("request").IndexInt(0).GetAttr("username_wo"))
			if diags.HasError() {
				return *webhook, fmt.Errorf("unable to read request.username_wo: %v", diags[0].Summary)
			}
			request.Username = username
		}

		// Password (either the regular attribute or the write-only one)
		if password, ok := requestMap["password"]; ok {
			request.Password = password.(string)
		}
		if request.Password == "" {
			password, diags := rawConfigString(d, cty.GetAttrPath("request").IndexInt(0).GetAttr("password_wo"))
			if diags.HasError() {
				return *webhook, fmt.Errorf("unable to read request.password_wo: %v", diags[0].Summary)
			}
			request.Password = password
		}

		// URL
		if url, ok := requestMap["url"]; ok {
			request.URL = url.(string)
		}

		// Convert headers JSON string into map type
		if headers, ok := requestMap["headers"]; ok {
			request.Headers = make(map[string]string)
			if headers, ok := headers.(map[string]interface{}); ok {
				for k, v := range headers {
					log.Println("[DEBUG] header", k, "type", reflect.TypeOf(v))
					request.Headers[k] = v.(string)
				}
			} else {
				err := fmt.Errorf("unable parse request headers into a map[string]interface, got %v", reflect.TypeOf(requestMap["headers"]))
				log.Print("[ERROR] error", err)
				return *webhook, err
			}
		} else {
			log.Printf("[ERROR] 'headers' is a required field")
			return *webhook, fmt.Errorf("headers is a mandatory field")
		}

		// Write-only headers, merged with the regular ones
		rawHeaders, diags := rawConfigString(d, cty.GetAttrPath("request").IndexInt(0).GetAttr("headers_wo"))
		if diags.HasError() {
			return *webhook, fmt.Errorf("unable to read request.headers_wo: %v", diags[0].Summary)
		}
		woHeaders, err := parseHeadersJSON(rawHeaders)
		if err != nil {
			return *webhook, fmt.Errorf("request.headers_wo %v", err)
		}
		for k, v := range woHeaders {
			if _, exists := request.Headers[k]; exists {
				return *webhook, fmt.Errorf("header %q is set in both request.headers and request.headers_wo", k)
			}
			request.Headers[k] = v
		}

		// Body
		if body, ok := requestMap["body"]; ok {
			// parse JSON into an intermediate object if possible, as this will avoid double escaping of the
			// JSON (e.g. quotes) when it's sent over the wire
			var i interface{}
			err := json.Unmarshal([]byte(body.(string)), &i)
			if err != nil {
				log.Println("[DEBUG] unable to parse JSON, default to string")
				request.Body = body.(string)
			} else {
				request.Body = i
			}
		}

		log.Printf("[DEBUG] have fully serialised request %+v \n", redactRequest(*request))

		webhook.Request = *request
	} else {
		log.Println("[ERROR] request attribute not found")
		return *webhook, fmt.Errorf("request is a mandatory field")
	}

	// Existing webhook for update?
	if d.Id() != "" {
		webhook.ID = d.Id()
	}

	return *webhook, nil
}

func setWebhookState(d *schema.ResourceData, webhook broker.Webhook) error {
	log.Printf("[DEBUG] setting webhook state: %+v \n", redactWebhook(webhook))
	if err := d.Set("description", webhook.Description); err != nil {
		log.Println("[ERROR] error setting key 'description'", err)
		return err
	}

	if err := d.Set("enabled", webhook.Enabled); err != nil {
		log.Println("[ERROR] error setting key 'enabled'", err)
		return err
	}

	if err := d.Set("team", webhook.TeamUUID); err != nil {
		log.Println("[ERROR] error setting key 'team'", err)
		return err
	}

	if webhook.Consumer != nil {
		if err := d.Set("webhook_consumer", map[string]interface{}{
			"name": webhook.Consumer.Name,
		}); err != nil {
			log.Println("[ERROR] error setting key 'webhook_consumer'", err)
			return err
		}
	} else {
		d.Set("webhook_consumer", nil)
	}

	if webhook.Provider != nil {
		if err := d.Set("webhook_provider", map[string]interface{}{
			"name": webhook.Provider.Name,
		}); err != nil {
			log.Println("[ERROR] error setting key 'webhook_provider", err)
			return err
		}
	} else {
		d.Set("webhook_provider", nil)
	}

	if err := d.Set("events", flattenEvents(webhook)); err != nil {
		log.Println("[ERROR] error setting key 'events'", err)
		return err
	}

	if err := d.Set("request", flattenRequest(d, webhook.Request)); err != nil {
		log.Println("[ERROR] error setting key 'request'", err)
		return err
	}
	return nil
}

func flattenEvents(w broker.Webhook) []string {
	events := make([]string, len(w.Events), len(w.Events))
	for i, event := range w.Events {
		events[i] = event.Name
	}

	sort.Strings(events)

	return events
}

func flattenRequest(d *schema.ResourceData, r broker.Request) []interface{} {
	// NOTE: the top level structure to set is a map
	m := make(map[string]interface{})
	m["url"] = r.URL
	m["method"] = r.Method

	// Never persist values that may come from the write-only `username_wo`
	// attribute: keep whatever `username` value Terraform already has.
	if version, ok := d.GetOk("request.0.username_wo_version"); ok {
		m["username_wo_version"] = version.(int)
		m["username"] = d.Get("request.0.username").(string)
	} else {
		m["username"] = r.Username
	}

	// The broker obscures the password ("*****"), and it may have been provided
	// via the write-only `password_wo` attribute, which must never be persisted.
	// Keep whatever `password` value Terraform already has (config or prior state).
	if original, ok := d.GetOk("request.0.password"); ok {
		m["password"] = original.(string)
	}
	if version, ok := d.GetOk("request.0.password_wo_version"); ok {
		m["password_wo_version"] = version.(int)
	}

	// When write-only headers are in use, only the headers Terraform already
	// knows about (`headers`) are stored; the others come from `headers_wo`.
	headers := r.Headers
	if version, ok := d.GetOk("request.0.headers_wo_version"); ok {
		m["headers_wo_version"] = version.(int)
		known := d.Get("request.0.headers").(map[string]interface{})
		headers = make(map[string]string, len(known))
		for k, v := range r.Headers {
			if _, ok := known[k]; ok {
				headers[k] = v
			}
		}
	}
	m["headers"] = mapStringStringToMapStringInterface(headers)

	// We want to store the body as a string in the state file
	// Try to parse body into JSON, fallback to a string if not
	if bodyAsStr, ok := r.Body.(string); ok {
		log.Println("[DEBUG] parsed webhook body as string", bodyAsStr)
		m["body"] = bodyAsStr
	} else if bytes, err := json.Marshal(r.Body); err == nil {
		log.Println("[DEBUG] parsed webhook body as JSON", string(bytes))
		m["body"] = string(bytes)
	} else {
		log.Println("[DEBUG] unable to parse the body as a JSON string or a plain string!")
	}

	return []interface{}{m}
}

// Lowercases all keys
func mapStringStringToMapStringInterface(in map[string]string) map[string]interface{} {
	var out = make(map[string]interface{}, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

func webhookCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	httpClient := meta.(*client.Client)
	webhook, err := parseWebhook(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	res, err := httpClient.CreateWebhook(webhook)
	if err != nil {
		log.Println("[ERROR] webhook creation failed", err)
		d.SetId("")
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] response from creating webhook %+v\n", redactWebhook(res.Webhook))

	items := strings.Split(res.Links["self"].Href, "/")
	id := items[len(items)-1]
	d.SetId(id)

	return diag.FromErr(setWebhookState(d, webhook))
}

func webhookUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	httpClient := meta.(*client.Client)
	webhook, err := parseWebhook(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	res, err := httpClient.UpdateWebhook(webhook)
	if err != nil {
		log.Println("[ERROR] webhook update failed", err)
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] response from updating webhook %+v\n", redactWebhook(res.Webhook))

	return diag.FromErr(setWebhookState(d, webhook))
}

func webhookRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	httpClient := meta.(*client.Client)
	res, err := httpClient.ReadWebhook(d.Id())

	if err != nil {
		log.Println("[ERROR] webhook read failed", err)
		d.SetId("")
		return nil
	}
	log.Printf("[DEBUG] response from reading webhook %+v\n", redactWebhook(*res))

	return diag.FromErr(setWebhookState(d, *res))
}

func webhookDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	httpClient := meta.(*client.Client)

	log.Println("[DEBUG] deleting webhook", d.Id())

	if err := httpClient.DeleteWebhook(broker.Webhook{ID: d.Id()}); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

// redactRequest returns a copy of the request that is safe to log.
func redactRequest(r broker.Request) broker.Request {
	if r.Password != "" {
		r.Password = "*****"
	}
	if len(r.Headers) > 0 {
		headers := make(map[string]string, len(r.Headers))
		for k := range r.Headers {
			headers[k] = "*****"
		}
		r.Headers = headers
	}
	return r
}

// redactWebhook returns a copy of the webhook that is safe to log.
func redactWebhook(w broker.Webhook) broker.Webhook {
	w.Request = redactRequest(w.Request)
	return w
}

func tryParseJSONObject(s string) interface{} {
	log.Println("[DEBUG] checking if", s, "is a JSON string")
	var i interface{}
	err := json.Unmarshal([]byte(s), &i)

	if err != nil {
		log.Println("[DEBUG] input body is not JSON")
		return nil
	}

	return i
}

func ignoreJSONFormatting(k, old, new string, d *schema.ResourceData) bool {
	// old = strings.TrimSpace(tryParseJSONString(old))
	// new = strings.TrimSpace(tryParseJSONString(new))
	log.Println("[DEBUG] checking if we should ignore white space and JSON formatting", old, new)

	if tryParseJSONObject(old) != nil && reflect.DeepEqual(tryParseJSONObject(old), tryParseJSONObject(new)) {
		log.Println("[DEBUG] JSON bodies are identical")
		return true
	}

	if strings.TrimSpace(old) == strings.TrimSpace(new) {
		return true
	}

	return false
}

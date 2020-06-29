package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/pact-foundation/terraform/broker"
)

const (
	libraryVersion                      = "0.0.1"
	userAgent                           = "go-pact/" + libraryVersion
	defaultBaseURL                      = "http://localhost"
	webhookReadUpdateDeleteTemplate     = "/webhooks/%s"
	webhookCreateTemplate               = "/webhooks"
	pacticipantReadUpdateDeleteTemplate = "/pacticipants/%s"
	pacticipantCreateTemplate           = "/pacticipants"
	userReadUpdateDeleteTemplate        = "/admin/users/%s"
	userCreateTemplate                  = "/admin/users/invite-user"
	userAdminUpdateTemplate             = "/admin/users/%s/role/admin"
	secretReadUpdateDeleteTemplate      = "/secrets/%s"
	secretCreateTemplate                = "/secrets"
	listTokensTemplate                  = "/settings/tokens"
	tokenRegenerateTemplate             = "/settings/tokens/%s/regenerate"
	metadataTemplate                    = "/"
	// {"_links":{"self":{"href":"https://dius.pact.dius.com.au","title":"Index","templated":false},"pb:publish-pact":{"href":"https://dius.pact.dius.com.au/pacts/provider/{provider}/consumer/{consumer}/version/{consumerApplicationVersion}","title":"Publish a pact","templated":true},"pb:latest-pact-versions":{"href":"https://dius.pact.dius.com.au/pacts/latest","title":"Latest pact versions","templated":false},"pb:tagged-pact-versions":{"href":"https://dius.pact.dius.com.au/pacts/provider/{provider}/consumer/{consumer}/tag/{tag}","title":"All versions of a pact for a given consumer, provider and consumer version tag","templated":false},"pb:pacticipants":{"href":"https://dius.pact.dius.com.au/pacticipants","title":"Pacticipants","templated":false},"pb:pacticipant":{"href":"https://dius.pact.dius.com.au/pacticipants/{pacticipant}","title":"Fetch pacticipant by name","templated":true},"pb:latest-provider-pacts":{"href":"https://dius.pact.dius.com.au/pacts/provider/{provider}/latest","title":"Latest pacts by provider","templated":true},"pb:latest-provider-pacts-with-tag":{"href":"https://dius.pact.dius.com.au/pacts/provider/{provider}/latest/{tag}","title":"Latest pacts for provider with the specified tag","templated":true},"pb:provider-pacts-with-tag":{"href":"https://dius.pact.dius.com.au/pacts/provider/{provider}/tag/{tag}","title":"All pact versions for the provider with the specified consumer version tag","templated":true},"pb:provider-pacts":{"href":"https://dius.pact.dius.com.au/pacts/provider/{provider}","title":"All pact versions for the specified provider","templated":true},"pb:latest-version":{"href":"https://dius.pact.dius.com.au/pacticipants/{pacticipant}/latest-version","title":"Latest pacticipant version","templated":true},"pb:latest-tagged-version":{"href":"https://dius.pact.dius.com.au/pacticipants/{pacticipant}/latest-version/{tag}","title":"Latest pacticipant version with the specified tag","templated":true},"pb:webhooks":{"href":"https://dius.pact.dius.com.au/webhooks","title":"Webhooks","templated":false},"pb:webhook":{"href":"https://dius.pact.dius.com.au/webhooks/{uuid}","title":"Webhook","templated":true},"pb:integrations":{"href":"https://dius.pact.dius.com.au/integrations","title":"Integrations","templated":false},"pb:pacticipant-version-tag":{"href":"https://dius.pact.dius.com.au/pacticipants/{pacticipant}/versions/{version}/tags/{tag}","title":"Get, create or delete a tag for a pacticipant version","templated":true},"pb:metrics":{"href":"https://dius.pact.dius.com.au/metrics","title":"Get Pact Broker metrics"},"pb:can-i-deploy-pacticipant-version-to-tag":{"href":"https://dius.pact.dius.com.au/can-i-deploy?pacticipant={pacticipant}\u0026version={version}\u0026to={tag}","title":"Determine if an application can be safely deployed to an environment identified by the given tag","templated":true},"curies":[{"name":"pb","href":"https://dius.pact.dius.com.au/doc/{rel}?context=index","templated":true},{"name":"beta","href":"https://dius.pact.dius.com.au/doc/{rel}?context=index","templated":true}],"beta:provider-pacts-for-verification":{"href":"https://dius.pact.dius.com.au/pacts/provider/{provider}/for-verification","title":"Pact versions to be verified for the specified provider","templated":true},"pb:api-tokens":{"href":"https://dius.pact.dius.com.au/settings/tokens","title":"API tokens","templated":false},"pb:audit":{"href":"https://dius.pact.dius.com.au/audit","title":"Audit trail","templated":false},"pb:secrets":{"href":"https://dius.pact.dius.com.au/secrets","title":"Secrets","templated":false}}}
)

const (
	readOnlyTokenType  = "read-only"
	readWriteTokenType = "read-write"
)

var tokenTypes = map[string]string{
	readOnlyTokenType:  "Read only token (developer)",
	readWriteTokenType: "Read/write token (CI)",
}

var (
	ErrBadRequest        = errors.New("bad request")
	ErrSystemUnavailable = errors.New("system unavailable")
	ErrUnauthorized      = errors.New("unauthorized")
)

// Config is the primary means to modify the Pact Broker http client
type Config struct {
	AccessToken       string
	BasicAuthUsername string
	BasicAuthPassword string
	BaseURL           *url.URL
	CustomTLSConfig   *tls.Config
}

// Client is the main Broker API interface.
// Use NewClient to get started
type Client struct {
	client    http.Client
	Config    Config
	UserAgent string
}

// NewClient creates a new Broker API client with sensible but overridable defaults
func NewClient(httpClient *http.Client, config Config) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	if config.CustomTLSConfig.InsecureSkipVerify {
		httpClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	client := Client{
		client:    *httpClient,
		Config:    config,
		UserAgent: userAgent,
	}

	return &client
}

// ReadWebhook returns a Webhook or an error for a given ID
func (c *Client) ReadWebhook(id string) (*broker.Webhook, error) {
	res, err := c.doCrud("GET", fmt.Sprintf(webhookReadUpdateDeleteTemplate, id), nil, new(broker.Webhook))
	return res.(*broker.Webhook), err
}

// CreateWebhook creates a new webhook
func (c *Client) CreateWebhook(w broker.Webhook) (*broker.WebhookResponse, error) {
	res, err := c.doCrud("POST", webhookCreateTemplate, w, new(broker.WebhookResponse))
	return res.(*broker.WebhookResponse), err
}

// UpdateWebhook updates an existing webhook. Not all properties are mutable
func (c *Client) UpdateWebhook(w broker.Webhook) (*broker.WebhookResponse, error) {
	res, err := c.doCrud("PUT", fmt.Sprintf(webhookReadUpdateDeleteTemplate, w.ID), w, new(broker.WebhookResponse))
	return res.(*broker.WebhookResponse), err
}

// DeleteWebhook removes an existing webhook
func (c *Client) DeleteWebhook(w broker.Webhook) error {
	_, err := c.doCrud("DELETE", fmt.Sprintf(webhookReadUpdateDeleteTemplate, w.ID), nil, nil)
	return err
}

// ReadPacticipant gets a pacticipant
func (c *Client) ReadPacticipant(name string) (*broker.Pacticipant, error) {
	res, err := c.doCrud("GET", fmt.Sprintf(pacticipantReadUpdateDeleteTemplate, name), nil, new(broker.Pacticipant))
	return res.(*broker.Pacticipant), err
}

// CreatePacticipant creates a new Pacticipant
func (c *Client) CreatePacticipant(p broker.Pacticipant) (*broker.Pacticipant, error) {
	res, err := c.doCrud("POST", pacticipantCreateTemplate, p, new(broker.Pacticipant))
	return res.(*broker.Pacticipant), err
}

// UpdatePacticipant updates an existing Pacticipant
func (c *Client) UpdatePacticipant(p broker.Pacticipant) (*broker.Pacticipant, error) {
	res, err := c.doCrud("PATCH", pacticipantReadUpdateDeleteTemplate, p, new(broker.Pacticipant))
	return res.(*broker.Pacticipant), err
}

// DeletePacticipant removes an existing Pacticipant
func (c *Client) DeletePacticipant(p broker.Pacticipant) error {
	_, err := c.doCrud("DELETE", fmt.Sprintf(pacticipantReadUpdateDeleteTemplate, p.Name), nil, nil)
	return err
}

// ReadUser gets a User
func (c *Client) ReadUser(name string) (*broker.User, error) {
	res, err := c.doCrud("GET", fmt.Sprintf(userReadUpdateDeleteTemplate, name), nil, new(broker.User))
	return res.(*broker.User), err
}

// CreateUsers creates a user
func (c *Client) CreateUser(p broker.User) (*broker.User, error) {
	res, err := c.doCrud("POST", userCreateTemplate, p, new(broker.User))
	return res.(*broker.User), err
}

// UpdateUser updates an existing User
// currently only supports modifying the "active" property
func (c *Client) UpdateUser(p broker.User) (*broker.User, error) {
	res, err := c.doCrud("PUT", fmt.Sprintf(userReadUpdateDeleteTemplate, p.UUID), p, new(broker.User))
	return res.(*broker.User), err
}

// DeleteUser removes an existing User
// TODO: is this even possible?
// Perhaps we _should_ allow deleting a user so that it's even gone from the UI
func (c *Client) DeleteUser(p broker.User) error {
	p.Active = false
	_, err := c.UpdateUser(p)

	return err
}

// AddAdminRoleToUser converts a user to an administrator
func (c *Client) AddAdminRoleToUser(p broker.User) (*broker.User, error) {
	res, err := c.doCrud("PUT", fmt.Sprintf(userAdminUpdateTemplate, p.UUID), p, new(broker.User))
	return res.(*broker.User), err
}

// RemoveAdminRoleToUser removes the administrator role from a user
func (c *Client) RemoveAdminRoleToUser(p broker.User) (*broker.User, error) {
	res, err := c.doCrud("DELETE", fmt.Sprintf(userAdminUpdateTemplate, p.UUID), p, new(broker.User))
	return res.(*broker.User), err
}

// ReadSecret gets the current Secret information (the actual secret is not returned)
func (c *Client) ReadSecret(uuid string) (*broker.SecretResponseCreate, error) {
	res, err := c.doCrud("GET", fmt.Sprintf(secretReadUpdateDeleteTemplate, uuid), nil, new(broker.SecretResponseCreate))
	return res.(*broker.SecretResponseCreate), err
}

// CreateSecret creates a new secret
// TODO: better response message for OSS broker vs Pactflow
func (c *Client) CreateSecret(s broker.Secret) (*broker.SecretResponseCreate, error) {
	res, err := c.doCrud("POST", secretCreateTemplate, s, new(broker.SecretResponseCreate))
	return res.(*broker.SecretResponseCreate), err
}

// UpdateSecret updates an existing secret. All values may be changed
func (c *Client) UpdateSecret(s broker.Secret) (*broker.SecretResponseCreate, error) {
	res, err := c.doCrud("PUT", fmt.Sprintf(secretReadUpdateDeleteTemplate, s.UUID), s, new(broker.SecretResponseCreate))
	return res.(*broker.SecretResponseCreate), err
}

// DeleteSecret removes an existing secret
func (c *Client) DeleteSecret(s broker.Secret) error {
	_, err := c.doCrud("DELETE", fmt.Sprintf(secretReadUpdateDeleteTemplate, s.UUID), nil, nil)
	return err
}

// ListTokens lists all tokens for the given user principal
func (c *Client) ListTokens() (*broker.APITokensResponse, error) {
	res, err := c.doCrud("GET", listTokensTemplate, nil, new(broker.APITokensResponse))
	return res.(*broker.APITokensResponse), err
}

// FindTokenByUUID finds a token given a UUID
func (c *Client) FindTokenByUUID(uuid string) (*broker.APIToken, error) {
	tokens, err := c.ListTokens()
	log.Println("[DEBUG] have tokens", tokens)

	if err != nil {
		return nil, err
	}
	for _, t := range tokens.Embedded.Items {
		log.Println("[DEBUG] have token", t)
		if t.UUID == uuid {
			return &t, nil
		}
	}
	return nil, fmt.Errorf("token with uuid '%s' not found", uuid)
}

// FindTokenByType finds a token given it's s
// NOTE: this API will be deprecated once a full CRUD API is available
func (c *Client) FindTokenByType(tokenType string) (*broker.APIToken, error) {
	if _, ok := tokenTypes[tokenType]; !ok {
		return nil, fmt.Errorf("invalid token type specified, need one of %v, got %s", tokenTypes, tokenType)
	}

	tokens, err := c.ListTokens()
	log.Println("[DEBUG] have tokens", tokens)

	if err != nil {
		return nil, err
	}
	for _, t := range tokens.Embedded.Items {
		log.Println("[DEBUG] have token", t)
		if t.Description == tokenTypes[tokenType] {
			return &t, nil
		}
	}
	return nil, fmt.Errorf("token of type %s not found", tokenType)
}

// RegenerateToken generates a new API Token for the given UUID
func (c *Client) RegenerateToken(t broker.APIToken) (*broker.APITokenResponse, error) {
	res, err := c.doCrud("POST", fmt.Sprintf(tokenRegenerateTemplate, t.UUID), nil, new(broker.APITokenResponse))
	return res.(*broker.APITokenResponse), err
}

func (c *Client) newRequest(method, path string, body interface{}) (*http.Request, error) {
	rel := &url.URL{Path: path}
	u := c.Config.BaseURL.ResolveReference(rel)
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if c.Config.AccessToken != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.Config.AccessToken))
	} else if c.Config.BasicAuthUsername != "" {
		req.SetBasicAuth(c.Config.BasicAuthUsername, c.Config.BasicAuthPassword)
	}

	req.Header.Set("Accept", "application/hal+json")
	req.Header.Set("User-Agent", c.UserAgent)

	log.Println("[DEBUG] creating new request", req)
	return req, nil
}

func (c *Client) do(req *http.Request, v interface{}) (*http.Response, error) {
	log.Println("[DEBUG] sending body for request", req)
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	// Drain and close the body to let the Transport reuse the connection
	// See https://github.com/google/go-github/pull/317/files for more info/background
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()
	log.Println("[DEBUG] response for request:", req, "resp:", resp)

	if resp.StatusCode >= 500 {
		return nil, ErrSystemUnavailable
	}

	if resp.StatusCode >= 401 {
		return nil, ErrUnauthorized
	}

	if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		return nil, ErrBadRequest
	}

	if resp.StatusCode >= 300 && resp.StatusCode < 400 {
		// TODO: deal with redirect
	}

	if v != nil {
		err = json.NewDecoder(resp.Body).Decode(v)
		if err != nil {
			log.Println("[DEBUG] error decoding response for", req.URL.Path, ". Error", err)
			return resp, err
		}
		log.Printf("[DEBUG] Response body: %+v \n", v)
	}

	return resp, err
}

func (c *Client) doCrud(method string, path string, requestEntity interface{}, responseEntity interface{}) (interface{}, error) {
	req, err := c.newRequest(method, path, requestEntity)
	if err != nil {
		return responseEntity, err
	}
	if responseEntity == nil {
		_, err = c.do(req, nil)
	} else {
		_, err = c.do(req, &responseEntity)
	}
	return responseEntity, err
}

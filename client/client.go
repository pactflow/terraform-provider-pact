package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/pactflow/terraform/broker"
	"github.com/pactflow/terraform/version"
)

const (
	userAgent                           = "terraform-pact/" + version.LIBRARY_VERSION
	defaultBaseURL                      = "http://localhost"
	webhookReadUpdateDeleteTemplate     = "/webhooks/%s"
	webhookCreateTemplate               = "/webhooks"
	pacticipantReadUpdateDeleteTemplate = "/pacticipants/%s"
	pacticipantCreateTemplate           = "/pacticipants"
	teamReadUpdateDeleteTemplate        = "/admin/teams/%s"
	teamCreateTemplate                  = "/admin/teams"
	teamAssignmentTemplate              = "/admin/teams/%s/users"
	teamUserTemplate                    = "/admin/teams/%s/users/%s"
	tenantAuthenticationTemplate        = "/admin/tenant/authentication-settings"
	roleCreateTemplate                  = "/admin/roles"
	roleReadUpdateDeleteTemplate        = "/admin/roles/%s"
	userReadUpdateDeleteTemplate        = "/admin/users/%s"
	userRolesUpdateTemplate             = "/admin/users/%s/roles"
	userRolesDeleteAppendTemplate       = "/admin/users/%s/roles/%s"
	userCreateTemplate                  = "/admin/users/invite-user"
	systemAccountCreateTemplate         = "/admin/system-accounts"
	userAdminUpdateTemplate             = "/admin/users/%s/role/admin"
	secretReadUpdateDeleteTemplate      = "/secrets/%s"
	secretCreateTemplate                = "/secrets"
	listTokensTemplate                  = "/settings/tokens"
	tokenRegenerateTemplate             = "/settings/tokens/%s/regenerate"
	metadataTemplate                    = "/"
	environmentCreateTemplate           = "/environments"
	environmentReadUpdateDeleteTemplate = "/environments/%s"
)

const (
	readOnlyTokenType  = "read-only"
	readWriteTokenType = "read-write"
)

var tokenTypes = map[string]string{
	readOnlyTokenType:  "Read only token (developer)",
	readWriteTokenType: "Read/write token (CI)",
}

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

	if config.CustomTLSConfig != nil && config.CustomTLSConfig.InsecureSkipVerify {
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
	res, err := c.doCrud("GET", urlEncodeTemplate(webhookReadUpdateDeleteTemplate, id), nil, new(broker.Webhook))
	return res.(*broker.Webhook), err
}

// CreateWebhook creates a new webhook
func (c *Client) CreateWebhook(w broker.Webhook) (*broker.WebhookResponse, error) {
	res, err := c.doCrud("POST", webhookCreateTemplate, w, new(broker.WebhookResponse))
	return res.(*broker.WebhookResponse), err
}

// UpdateWebhook updates an existing webhook. Not all properties are mutable
func (c *Client) UpdateWebhook(w broker.Webhook) (*broker.WebhookResponse, error) {
	res, err := c.doCrud("PUT", urlEncodeTemplate(webhookReadUpdateDeleteTemplate, w.ID), w, new(broker.WebhookResponse))
	return res.(*broker.WebhookResponse), err
}

// DeleteWebhook removes an existing webhook
func (c *Client) DeleteWebhook(w broker.Webhook) error {
	_, err := c.doCrud("DELETE", urlEncodeTemplate(webhookReadUpdateDeleteTemplate, w.ID), nil, nil)
	return err
}

// ReadPacticipant gets a pacticipant
func (c *Client) ReadPacticipant(name string) (*broker.Pacticipant, error) {
	res, err := c.doCrud("GET", urlEncodeTemplate(pacticipantReadUpdateDeleteTemplate, name), nil, new(broker.Pacticipant))
	return res.(*broker.Pacticipant), err
}

// CreatePacticipant creates a new Pacticipant
func (c *Client) CreatePacticipant(p broker.Pacticipant) (*broker.Pacticipant, error) {
	res, err := c.doCrud("POST", pacticipantCreateTemplate, p, new(broker.Pacticipant))
	return res.(*broker.Pacticipant), err
}

// UpdatePacticipant updates an existing Pacticipant
func (c *Client) UpdatePacticipant(p broker.Pacticipant) (*broker.Pacticipant, error) {
	res, err := c.doCrud("PATCH", urlEncodeTemplate(pacticipantReadUpdateDeleteTemplate, p.Name), p, new(broker.Pacticipant))
	return res.(*broker.Pacticipant), err
}

// DeletePacticipant removes an existing Pacticipant
func (c *Client) DeletePacticipant(p broker.Pacticipant) error {
	_, err := c.doCrud("DELETE", urlEncodeTemplate(pacticipantReadUpdateDeleteTemplate, p.Name), nil, nil)
	return err
}

// ReadTeam gets a Team
func (c *Client) ReadTeam(t broker.Team) (*broker.Team, error) {
	res, err := c.doCrud("GET", urlEncodeTemplate(teamReadUpdateDeleteTemplate, t.UUID), nil, new(broker.Team))
	return res.(*broker.Team), err
}

// CreateTeam creates a Team
func (c *Client) CreateTeam(t broker.TeamCreateOrUpdateRequest) (*broker.Team, error) {
	res, err := c.doCrud("POST", teamCreateTemplate, t, new(broker.Team))
	return res.(*broker.Team), err
}

// ReadTeamAssignments finds all users currently in a team
func (c *Client) ReadTeamAssignments(t broker.Team) (*broker.TeamsAssignmentResponse, error) {
	res, err := c.doCrud("GET", urlEncodeTemplate(teamAssignmentTemplate, t.UUID), t, new(broker.TeamsAssignmentResponse))
	return res.(*broker.TeamsAssignmentResponse), err
}

// UpdateTeamAssignments sets the users for a given team, removing any existing users not in the specified request
func (c *Client) UpdateTeamAssignments(r broker.TeamsAssignmentRequest) (*broker.TeamsAssignmentResponse, error) {
	res, err := c.doCrud("PUT", urlEncodeTemplate(teamAssignmentTemplate, r.UUID), r, new(broker.TeamsAssignmentResponse))

	if err != nil {
		return nil, err
	}

	if len(r.Users) > 0 {
		apiResponse := res.(*broker.TeamsAssignmentResponse)

		return apiResponse, err
	}

	return nil, err
}

// AppendTeamAssignments adds users to an existing Team (does not remove absent ones)
func (c *Client) AppendTeamAssignments(r broker.TeamsAssignmentRequest) (*broker.TeamsAssignmentResponse, error) {
	res, err := c.doCrud("POST", urlEncodeTemplate(teamAssignmentTemplate, r.UUID), r, new(broker.TeamsAssignmentResponse))

	if err != nil {
		return nil, err
	}

	if len(r.Users) > 0 {
		apiResponse := res.(*broker.TeamsAssignmentResponse)

		return apiResponse, err
	}

	return nil, err
}

// DeleteTeamAssignment removes a single user from a team
func (c *Client) DeleteTeamAssignment(t broker.Team, u broker.User) error {
	_, err := c.doCrud("DELETE", urlEncodeTemplate(teamUserTemplate, t.UUID, u.UUID), nil, nil)

	return err
}

// DeleteTeamAssignments removes specified users from the team
func (c *Client) DeleteTeamAssignments(t broker.TeamsAssignmentRequest) error {
	if len(t.Users) > 0 {
		_, err := c.doCrud("DELETE", urlEncodeTemplate(teamAssignmentTemplate, t.UUID), t, nil)
		return err
	}
	return nil
}

// UpdateTeam updates the team
func (c *Client) UpdateTeam(t broker.TeamCreateOrUpdateRequest) (*broker.Team, error) {
	res, err := c.doCrud("PUT", urlEncodeTemplate(teamReadUpdateDeleteTemplate, t.UUID), t, new(broker.Team))
	return res.(*broker.Team), err
}

// DeleteTeam deletes the Team
func (c *Client) DeleteTeam(t broker.Team) error {
	_, err := c.doCrud("DELETE", urlEncodeTemplate(teamReadUpdateDeleteTemplate, t.UUID), nil, nil)

	return err
}

// ReadRole gets a Role
func (c *Client) ReadRole(uuid string) (*broker.Role, error) {
	res, err := c.doCrud("GET", urlEncodeTemplate(roleReadUpdateDeleteTemplate, uuid), nil, new(broker.Role))
	return res.(*broker.Role), err
}

// CreateRole creates a Role
func (c *Client) CreateRole(p broker.Role) (*broker.Role, error) {
	res, err := c.doCrud("POST", roleCreateTemplate, p, new(broker.Role))
	return res.(*broker.Role), err
}

// UpdateRole updates an existing Role
func (c *Client) UpdateRole(p broker.Role) (*broker.Role, error) {
	res, err := c.doCrud("PUT", urlEncodeTemplate(roleReadUpdateDeleteTemplate, p.UUID), p, new(broker.Role))
	return res.(*broker.Role), err
}

// DeleteRole removes a role
func (c *Client) DeleteRole(p broker.Role) error {
	_, err := c.doCrud("DELETE", urlEncodeTemplate(roleReadUpdateDeleteTemplate, p.UUID), nil, nil)

	return err
}

// ReadUser gets a User
func (c *Client) ReadUser(uuid string) (*broker.User, error) {
	res, err := c.doCrud("GET", urlEncodeTemplate(userReadUpdateDeleteTemplate, uuid), nil, new(broker.User))
	return res.(*broker.User), err
}

// CreateUser creates a user or a system account
func (c *Client) CreateUser(u broker.User) (*broker.User, error) {
	template := userCreateTemplate
	if u.Type == broker.SystemAccount {
		return c.CreateSystemAccount(u)
	}
	res, err := c.doCrud("POST", template, u, new(broker.User))
	return res.(*broker.User), err
}

// CreateUser creates a user or a system account
func (c *Client) CreateSystemAccount(u broker.User) (*broker.User, error) {
	res, err := c.doCrud("POST", systemAccountCreateTemplate, u, nil)

	if err != nil {
		return nil, err
	}

	// TODO: Returns a 201
	// e.g. https://tf-acceptance.pactflow.io/admin/system-accounts/f996d7e5-6525-4649-b479-9299793d105e
	// + a list of users
	// Extracting the UUID from response
	parts := strings.Split(res.(string), "/")
	u.UUID = parts[len(parts)-1]

	return &u, err
}

// UpdateUser updates an existing User
// currently only supports modifying the "active" property
func (c *Client) UpdateUser(p broker.User) (*broker.User, error) {
	res, err := c.doCrud("PUT", urlEncodeTemplate(userReadUpdateDeleteTemplate, p.UUID), p, new(broker.User))
	return res.(*broker.User), err
}

// DeleteUser simply de-activates an existing user. Users are global on the platform,
// but can be enabled/disabled at the tenant level
func (c *Client) DeleteUser(p broker.User) error {
	p.Active = false
	_, err := c.UpdateUser(p)

	return err
}

// AddAdminRoleToUser converts a user to an administrator
func (c *Client) AddAdminRoleToUser(p broker.User) (*broker.User, error) {
	res, err := c.doCrud("PUT", urlEncodeTemplate(userAdminUpdateTemplate, p.UUID), p, new(broker.User))
	return res.(*broker.User), err
}

// RemoveAdminRoleFromUser removes the administrator role from a user
func (c *Client) RemoveAdminRoleFromUser(p broker.User) (*broker.User, error) {
	res, err := c.doCrud("DELETE", urlEncodeTemplate(userAdminUpdateTemplate, p.UUID), p, new(broker.User))
	return res.(*broker.User), err
}

// ReadSecret gets the current Secret information (the actual secret is not returned)
func (c *Client) ReadSecret(uuid string) (*broker.SecretResponse, error) {
	res, err := c.doCrud("GET", urlEncodeTemplate(secretReadUpdateDeleteTemplate, uuid), nil, new(broker.SecretResponse))
	return res.(*broker.SecretResponse), err
}

// CreateSecret creates a new secret
// TODO: better response message for OSS broker vs Pactflow
func (c *Client) CreateSecret(s broker.Secret) (*broker.SecretResponse, error) {
	res, err := c.doCrud("POST", secretCreateTemplate, s, new(broker.SecretResponse))
	return res.(*broker.SecretResponse), err
}

// UpdateSecret updates an existing secret. All values may be changed
func (c *Client) UpdateSecret(s broker.Secret) (*broker.SecretResponse, error) {
	res, err := c.doCrud("PUT", urlEncodeTemplate(secretReadUpdateDeleteTemplate, s.UUID), s, new(broker.SecretResponse))
	return res.(*broker.SecretResponse), err
}

// DeleteSecret removes an existing secret
func (c *Client) DeleteSecret(s broker.Secret) error {
	_, err := c.doCrud("DELETE", urlEncodeTemplate(secretReadUpdateDeleteTemplate, s.UUID), nil, nil)
	return err
}

// ReadTokens lists all tokens for the given user principal
func (c *Client) ReadTokens() (*broker.APITokensResponse, error) {
	res, err := c.doCrud("GET", listTokensTemplate, nil, new(broker.APITokensResponse))
	return res.(*broker.APITokensResponse), err
}

// ReadToken finds an API token given a UUID
func (c *Client) ReadToken(uuid string) (*broker.APIToken, error) {
	tokens, err := c.ReadTokens()
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

	tokens, err := c.ReadTokens()
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
	res, err := c.doCrud("POST", urlEncodeTemplate(tokenRegenerateTemplate, t.UUID), nil, new(broker.APITokenResponse))
	return res.(*broker.APITokenResponse), err
}

// SetUserRoles sets the roles for a given user, removing any not given and adding those that were provided
func (c *Client) SetUserRoles(uuid string, r broker.SetUserRolesRequest) error {
	_, err := c.doCrud("PUT", urlEncodeTemplate(userRolesUpdateTemplate, uuid), r, nil)
	return err
}

// ReadTenantAuthenticationSettings configures the authentication settings on a given Pactflow account
func (c *Client) ReadTenantAuthenticationSettings() (*broker.AuthenticationSettings, error) {
	res, err := c.doCrud("GET", tenantAuthenticationTemplate, nil, new(broker.AuthenticationSettings))

	return res.(*broker.AuthenticationSettings), err
}

// SetTenantAuthenticationSettings configures the authentication settings on a given Pactflow account
func (c *Client) SetTenantAuthenticationSettings(r broker.AuthenticationSettings) (*broker.AuthenticationSettings, error) {
	res, err := c.doCrud("PUT", tenantAuthenticationTemplate, r, new(broker.AuthenticationSettings))

	return res.(*broker.AuthenticationSettings), err
}

// ReadEnvironment gets an Environment
func (c *Client) ReadEnvironment(uuid string) (*broker.Environment, error) {
	res, err := c.doCrud("GET", urlEncodeTemplate(environmentReadUpdateDeleteTemplate, uuid), nil, new(broker.Environment))
	return res.(*broker.Environment), err
}

// CreateEnvironment creates an Environment
func (c *Client) CreateEnvironment(p broker.EnvironmentCreateOrUpdateRequest) (*broker.EnvironmentCreateOrUpdateResponse, error) {
	res, err := c.doCrud("POST", environmentCreateTemplate, p, new(broker.EnvironmentCreateOrUpdateResponse))
	return res.(*broker.EnvironmentCreateOrUpdateResponse), err
}

// UpdateEnvironment updates an Environment
func (c *Client) UpdateEnvironment(p broker.EnvironmentCreateOrUpdateRequest) (*broker.EnvironmentCreateOrUpdateResponse, error) {
	res, err := c.doCrud("PUT", urlEncodeTemplate(environmentReadUpdateDeleteTemplate, p.UUID), p, new(broker.EnvironmentCreateOrUpdateResponse))
	return res.(*broker.EnvironmentCreateOrUpdateResponse), err
}

// DeleteEnvironment removes an Environment
func (c *Client) DeleteEnvironment(p broker.Environment) error {
	_, err := c.doCrud("DELETE", urlEncodeTemplate(environmentReadUpdateDeleteTemplate, p.UUID), nil, nil)

	return err
}

func (c *Client) newRequest(method, path string, body interface{}) (*http.Request, error) {
	rel := &url.URL{Path: path}
	u := c.Config.BaseURL.ResolveReference(rel)
	var buf = new(bytes.Buffer)
	if body != nil {
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}

		log.Printf("[INFO] raw body to be sent over wire: '%s'", buf.String())
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

	req.Header.Set("Accept", "application/hal+json, application/json")
	req.Header.Set("User-Agent", c.UserAgent)

	log.Println("[DEBUG] creating new request", req)
	return req, nil
}

func handleError(err error, req *http.Request, resp *http.Response) (*http.Response, error) {
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close() //  must close
	log.Println("[DEBUG] handling error response:", string(bodyBytes))

	// TODO: decode the multiple concrete error types here
	var e error

	e = &apiErrorResponse{
		err: err,
	}
	decodingErr := json.NewDecoder(bytes.NewBuffer(bodyBytes)).Decode(e)
	if decodingErr != nil {
		log.Println("[DEBUG] error decoding APIErrorResponse from response for", req.Method, req.URL.Path, ". Error", decodingErr)

		e = &apiArrayErrorResponse{
			err: err,
		}
		decodingErr = json.NewDecoder(bytes.NewBuffer(bodyBytes)).Decode(e)
		if decodingErr != nil {
			log.Println("[DEBUG] error decoding APIArrayErrorResponse from response for", req.Method, req.URL.Path, ". Error", decodingErr)
		}
	}

	return resp, e
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
		return handleError(ErrSystemUnavailable, req, resp)
	}

	if resp.StatusCode == 401 {
		return handleError(ErrUnauthorized, req, resp)
	}

	if resp.StatusCode == 403 {
		return handleError(ErrForbidden, req, resp)
	}

	if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		return handleError(ErrBadRequest, req, resp)
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
	var resp *http.Response

	if err != nil {
		return responseEntity, err
	}
	if responseEntity == nil {
		resp, err = c.do(req, nil)

		// 201 -> extract the location header if the expectation is a string value
		if resp.StatusCode == 201 {
			log.Println("[DEBUG] have 201, returning Location header", resp.Header)
			return resp.Header.Get("Location"), err
		}
	} else {
		_, err = c.do(req, &responseEntity)
	}

	return responseEntity, err
}

func urlEncodeTemplate(template string, parameters ...string) string {
	encodedParams := make([]interface{}, len(parameters))

	for i, p := range parameters {
		encodedParams[i] = url.PathEscape(p)
	}

	return fmt.Sprintf(template, encodedParams...)
}

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
	secretReadUpdateDeleteTemplate      = "/secrets/%s"
	secretCreateTemplate                = "/secrets"
)

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
	config    Config
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
		config:    config,
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
func (c *Client) CreateWebhook(w broker.Webhook) (*broker.WebhookResponseCreate, error) {
	res, err := c.doCrud("POST", webhookCreateTemplate, w, new(broker.WebhookResponseCreate))
	return res.(*broker.WebhookResponseCreate), err
}

// UpdateWebhook updates an existing webhook. Not all properties are mutable
func (c *Client) UpdateWebhook(w broker.Webhook) (*broker.WebhookResponseCreate, error) {
	res, err := c.doCrud("PUT", fmt.Sprintf(webhookReadUpdateDeleteTemplate, w.ID), w, new(broker.WebhookResponseCreate))
	return res.(*broker.WebhookResponseCreate), err
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

func (c *Client) newRequest(method, path string, body interface{}) (*http.Request, error) {
	rel := &url.URL{Path: path}
	u := c.config.BaseURL.ResolveReference(rel)
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

	if c.config.AccessToken != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.config.AccessToken))
	} else if c.config.BasicAuthUsername != "" {
		req.SetBasicAuth(c.config.BasicAuthUsername, c.config.BasicAuthPassword)
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

	if v != nil {
		err = json.NewDecoder(resp.Body).Decode(v)
		log.Println("[DEBUG] response err?", err)
	}

	if resp.StatusCode >= 500 {
		err = ErrSystemUnavailable
	}

	if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		err = ErrBadRequest
	}

	if resp.StatusCode >= 300 && resp.StatusCode < 400 {
		// TODO: deal with redirect
	}

	return resp, err
}

func (c *Client) doCrud(method string, path string, requestEntity interface{}, responseEntity interface{}) (interface{}, error) {
	req, err := c.newRequest(method, path, requestEntity)
	if err != nil {
		return responseEntity, err
	}
	_, err = c.do(req, &responseEntity)
	return responseEntity, err
}

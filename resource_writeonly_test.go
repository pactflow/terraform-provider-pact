package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/go-cty/cty"
	ctyjson "github.com/hashicorp/go-cty/cty/json"
	"github.com/hashicorp/go-cty/cty/msgpack"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeBroker records the JSON bodies sent to it and answers like the Pactflow API.
type fakeBroker struct {
	mu     sync.Mutex
	bodies map[string][]map[string]interface{} // "METHOD path" -> bodies
}

func newFakeBroker(t *testing.T) (*fakeBroker, *httptest.Server) {
	b := &fakeBroker{bodies: map[string][]map[string]interface{}{}}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		if len(raw) > 0 {
			var body map[string]interface{}
			require.NoError(t, json.Unmarshal(raw, &body))
			b.mu.Lock()
			key := r.Method + " " + r.URL.Path
			b.bodies[key] = append(b.bodies[key], body)
			b.mu.Unlock()
		}

		w.Header().Set("Content-Type", "application/hal+json")
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/secrets":
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"name":"MySecret","description":"desc","_links":{"self":{"href":"http://broker/secrets/secret-uuid"}}}`))
		case r.Method == http.MethodPut && r.URL.Path == "/secrets/secret-uuid":
			_, _ = w.Write([]byte(`{"name":"MySecret","description":"desc","_links":{"self":{"href":"http://broker/secrets/secret-uuid"}}}`))
		case r.Method == http.MethodPost && r.URL.Path == "/webhooks":
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"description":"hook","_links":{"self":{"href":"http://broker/webhooks/webhook-id"}}}`))
		case r.Method == http.MethodPut && r.URL.Path == "/webhooks/webhook-id":
			_, _ = w.Write([]byte(`{"description":"hook","_links":{"self":{"href":"http://broker/webhooks/webhook-id"}}}`))
		case r.Method == http.MethodGet && r.URL.Path == "/webhooks/webhook-id":
			// the broker echoes every header, including the write-only ones
			_, _ = w.Write([]byte(`{"description":"hook","enabled":true,"events":[{"name":"contract_published"}],` +
				`"request":{"url":"https://example.com/hook","method":"POST","username":"wo-user","password":"*****",` +
				`"headers":{"Content-Type":"application/json","Authorization":"Bearer t0k3n"},"body":{"pact":"url"}},` +
				`"_links":{"self":{"href":"http://broker/webhooks/webhook-id"}}}`))
		default:
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	t.Cleanup(srv.Close)
	return b, srv
}

func (b *fakeBroker) last(key string) map[string]interface{} {
	b.mu.Lock()
	defer b.mu.Unlock()
	bodies := b.bodies[key]
	if len(bodies) == 0 {
		return nil
	}
	return bodies[len(bodies)-1]
}

// testProviderServer returns a configured gRPC provider server, talking to the given broker URL.
func testProviderServer(t *testing.T, brokerURL string) (tfprotov5.ProviderServer, *schema.Provider) {
	t.Helper()
	p := Provider()
	srv := schema.NewGRPCProviderServer(p)

	cfgTypes := map[string]cty.Type{}
	for name, s := range p.Schema {
		switch s.Type {
		case schema.TypeBool:
			cfgTypes[name] = cty.Bool
		default:
			cfgTypes[name] = cty.String
		}
	}
	cfgType := cty.Object(cfgTypes)
	cfg := objectFromJSON(t, cfgType, `{"host": "`+brokerURL+`", "access_token": "token"}`)

	resp, err := srv.ConfigureProvider(context.Background(), &tfprotov5.ConfigureProviderRequest{
		TerraformVersion: "1.11.0",
		Config:           dynamicValue(t, cfgType, cfg),
	})
	require.NoError(t, err)
	requireNoDiags(t, resp.Diagnostics)

	return srv, p
}

// objectFromJSON builds a cty object of the given type from a (partial) JSON document; missing attributes are null.
func objectFromJSON(t *testing.T, ty cty.Type, doc string) cty.Value {
	t.Helper()
	var partial map[string]json.RawMessage
	require.NoError(t, json.Unmarshal([]byte(doc), &partial))

	vals := map[string]cty.Value{}
	for name, attrTy := range ty.AttributeTypes() {
		raw, ok := partial[name]
		if !ok {
			vals[name] = cty.NullVal(attrTy)
			continue
		}
		v, err := ctyjson.Unmarshal(raw, attrTy)
		require.NoError(t, err, name)
		vals[name] = v
	}
	return cty.ObjectVal(vals)
}

func dynamicValue(t *testing.T, ty cty.Type, v cty.Value) *tfprotov5.DynamicValue {
	t.Helper()
	b, err := msgpack.Marshal(v, ty)
	require.NoError(t, err)
	return &tfprotov5.DynamicValue{MsgPack: b}
}

func decodeDynamicValue(t *testing.T, ty cty.Type, dv *tfprotov5.DynamicValue) cty.Value {
	t.Helper()
	require.NotNil(t, dv)
	v, err := msgpack.Unmarshal(dv.MsgPack, ty)
	require.NoError(t, err)
	return v
}

func requireNoDiags(t *testing.T, diags []*tfprotov5.Diagnostic) {
	t.Helper()
	for _, d := range diags {
		if d.Severity == tfprotov5.DiagnosticSeverityError {
			t.Fatalf("unexpected error diagnostic: %s: %s", d.Summary, d.Detail)
		}
	}
}

func hasErrorDiag(diags []*tfprotov5.Diagnostic) bool {
	for _, d := range diags {
		if d.Severity == tfprotov5.DiagnosticSeverityError {
			return true
		}
	}
	return false
}

// planAndApply runs PlanResourceChange + ApplyResourceChange the way Terraform does, and returns the new state.
func planAndApply(t *testing.T, srv tfprotov5.ProviderServer, typeName string, ty cty.Type, prior, config cty.Value) cty.Value {
	t.Helper()
	ctx := context.Background()

	validate, err := srv.ValidateResourceTypeConfig(ctx, &tfprotov5.ValidateResourceTypeConfigRequest{
		TypeName:           typeName,
		Config:             dynamicValue(t, ty, config),
		ClientCapabilities: &tfprotov5.ValidateResourceTypeConfigClientCapabilities{WriteOnlyAttributesAllowed: true},
	})
	require.NoError(t, err)
	requireNoDiags(t, validate.Diagnostics)

	// Terraform nullifies write-only attributes in the proposed new state
	proposed := nullWriteOnly(t, config)
	if !prior.IsNull() {
		// approximate Terraform's proposed new state: (optional+)computed
		// attributes not set in the configuration keep their prior value
		vals := proposed.AsValueMap()
		priorVals := prior.AsValueMap()
		for k, v := range vals {
			if v.IsNull() && !strings.HasSuffix(k, "_wo") {
				vals[k] = priorVals[k]
			}
		}
		proposed = cty.ObjectVal(vals)
	}

	plan, err := srv.PlanResourceChange(ctx, &tfprotov5.PlanResourceChangeRequest{
		TypeName:         typeName,
		PriorState:       dynamicValue(t, ty, prior),
		ProposedNewState: dynamicValue(t, ty, proposed),
		Config:           dynamicValue(t, ty, config),
	})
	require.NoError(t, err)
	requireNoDiags(t, plan.Diagnostics)
	planned := decodeDynamicValue(t, ty, plan.PlannedState)
	assertNoWriteOnlyValues(t, planned, "planned state")

	apply, err := srv.ApplyResourceChange(ctx, &tfprotov5.ApplyResourceChangeRequest{
		TypeName:       typeName,
		PriorState:     dynamicValue(t, ty, prior),
		PlannedState:   plan.PlannedState,
		Config:         dynamicValue(t, ty, config),
		PlannedPrivate: plan.PlannedPrivate,
	})
	require.NoError(t, err)
	requireNoDiags(t, apply.Diagnostics)

	newState := decodeDynamicValue(t, ty, apply.NewState)
	assertNoWriteOnlyValues(t, newState, "new state")
	return newState
}

// nullWriteOnly replaces every attribute whose name ends in "_wo" by null.
func nullWriteOnly(t *testing.T, v cty.Value) cty.Value {
	t.Helper()
	out, err := cty.Transform(v, func(p cty.Path, v cty.Value) (cty.Value, error) {
		if len(p) > 0 {
			if step, ok := p[len(p)-1].(cty.GetAttrStep); ok && strings.HasSuffix(step.Name, "_wo") {
				return cty.NullVal(v.Type()), nil
			}
		}
		return v, nil
	})
	require.NoError(t, err)
	return out
}

func assertNoWriteOnlyValues(t *testing.T, v cty.Value, what string) {
	t.Helper()
	_ = cty.Walk(v, func(p cty.Path, v cty.Value) (bool, error) {
		if len(p) > 0 {
			if step, ok := p[len(p)-1].(cty.GetAttrStep); ok && strings.HasSuffix(step.Name, "_wo") {
				assert.True(t, v.IsNull(), "%s: %s must be null, got %#v", what, step.Name, v)
			}
		}
		return true, nil
	})
}

func TestSecret_ValueWriteOnly(t *testing.T) {
	broker, brokerSrv := newFakeBroker(t)
	srv, p := testProviderServer(t, brokerSrv.URL)
	ty := p.ResourcesMap["pact_secret"].CoreConfigSchema().ImpliedType()

	config := objectFromJSON(t, ty, `{"name": "MySecret", "description": "desc", "value_wo": "s3cr3t", "value_wo_version": 1}`)
	state := planAndApply(t, srv, "pact_secret", ty, cty.NullVal(ty), config)

	// The secret is sent to the broker...
	body := broker.last("POST /secrets")
	require.NotNil(t, body)
	assert.Equal(t, "s3cr3t", body["value"])

	// ...but never stored in the state
	assert.True(t, state.GetAttr("value").IsNull() || state.GetAttr("value").AsString() == "", "value must not be stored")
	assert.Equal(t, "secret-uuid", state.GetAttr("id").AsString())
	assert.Equal(t, "secret-uuid", state.GetAttr("uuid").AsString())
	assert.Equal(t, "1", state.GetAttr("value_wo_version").AsBigFloat().String())
	assert.NotContains(t, state.GoString(), "s3cr3t")

	// Rotating the secret: bump the version, the new value is pushed to the broker
	config = objectFromJSON(t, ty, `{"name": "MySecret", "description": "desc", "value_wo": "n3w-s3cr3t", "value_wo_version": 2}`)
	state = planAndApply(t, srv, "pact_secret", ty, state, config)

	body = broker.last("PUT /secrets/secret-uuid")
	require.NotNil(t, body)
	assert.Equal(t, "n3w-s3cr3t", body["value"])
	assert.NotContains(t, state.GoString(), "s3cr3t")
}

func TestSecret_Value(t *testing.T) {
	broker, brokerSrv := newFakeBroker(t)
	srv, p := testProviderServer(t, brokerSrv.URL)
	ty := p.ResourcesMap["pact_secret"].CoreConfigSchema().ImpliedType()

	config := objectFromJSON(t, ty, `{"name": "MySecret", "description": "desc", "value": "s3cr3t"}`)
	state := planAndApply(t, srv, "pact_secret", ty, cty.NullVal(ty), config)

	assert.Equal(t, "s3cr3t", broker.last("POST /secrets")["value"])
	// existing behaviour: the (sensitive) value is kept in the state
	assert.Equal(t, "s3cr3t", state.GetAttr("value").AsString())
}

func TestSecret_Validation(t *testing.T) {
	_, brokerSrv := newFakeBroker(t)
	srv, p := testProviderServer(t, brokerSrv.URL)
	ty := p.ResourcesMap["pact_secret"].CoreConfigSchema().ImpliedType()

	cases := map[string]string{
		"neither value nor value_wo":     `{"name": "MySecret", "description": "desc"}`,
		"both value and value_wo":        `{"name": "MySecret", "description": "desc", "value": "a", "value_wo": "b", "value_wo_version": 1}`,
		"value_wo without version":       `{"name": "MySecret", "description": "desc", "value_wo": "b"}`,
		"value_wo_version without value": `{"name": "MySecret", "description": "desc", "value": "a", "value_wo_version": 1}`,
	}

	for name, doc := range cases {
		t.Run(name, func(t *testing.T) {
			resp, err := srv.ValidateResourceTypeConfig(context.Background(), &tfprotov5.ValidateResourceTypeConfigRequest{
				TypeName:           "pact_secret",
				Config:             dynamicValue(t, ty, objectFromJSON(t, ty, doc)),
				ClientCapabilities: &tfprotov5.ValidateResourceTypeConfigClientCapabilities{WriteOnlyAttributesAllowed: true},
			})
			require.NoError(t, err)
			assert.True(t, hasErrorDiag(resp.Diagnostics), "expected a validation error")
		})
	}
}

func TestSecret_WriteOnlyRequiresSupportingTerraform(t *testing.T) {
	_, brokerSrv := newFakeBroker(t)
	srv, p := testProviderServer(t, brokerSrv.URL)
	ty := p.ResourcesMap["pact_secret"].CoreConfigSchema().ImpliedType()

	resp, err := srv.ValidateResourceTypeConfig(context.Background(), &tfprotov5.ValidateResourceTypeConfigRequest{
		TypeName: "pact_secret",
		Config:   dynamicValue(t, ty, objectFromJSON(t, ty, `{"name": "MySecret", "description": "desc", "value_wo": "b", "value_wo_version": 1}`)),
		// Terraform < 1.11 does not advertise write-only support
		ClientCapabilities: &tfprotov5.ValidateResourceTypeConfigClientCapabilities{WriteOnlyAttributesAllowed: false},
	})
	require.NoError(t, err)
	assert.True(t, hasErrorDiag(resp.Diagnostics), "expected an error for Terraform versions without write-only support")
}

func TestWebhook_PasswordWriteOnly(t *testing.T) {
	broker, brokerSrv := newFakeBroker(t)
	srv, p := testProviderServer(t, brokerSrv.URL)
	ty := p.ResourcesMap["pact_webhook"].CoreConfigSchema().ImpliedType()

	doc := func(password string, version int) string {
		return `{
			"description": "hook",
			"events": ["contract_published"],
			"request": [{
				"url": "https://example.com/hook",
				"method": "POST",
				"username": "user",
				"password_wo": "` + password + `",
				"password_wo_version": ` + strconv.Itoa(version) + `,
				"headers": {"Content-Type": "application/json"},
				"body": "{\"pact\":\"url\"}"
			}]
		}`
	}

	state := planAndApply(t, srv, "pact_webhook", ty, cty.NullVal(ty), objectFromJSON(t, ty, doc("hunter2", 1)))

	body := broker.last("POST /webhooks")
	require.NotNil(t, body)
	request := body["request"].(map[string]interface{})
	assert.Equal(t, "hunter2", request["password"])
	assert.Equal(t, "user", request["username"])

	assert.Equal(t, "webhook-id", state.GetAttr("id").AsString())
	req := state.GetAttr("request").Index(cty.NumberIntVal(0))
	assert.True(t, req.GetAttr("password").IsNull() || req.GetAttr("password").AsString() == "", "password must not be stored")
	assert.Equal(t, "1", req.GetAttr("password_wo_version").AsBigFloat().String())
	assert.NotContains(t, state.GoString(), "hunter2")

	// Rotate the password
	state = planAndApply(t, srv, "pact_webhook", ty, state, objectFromJSON(t, ty, doc("hunter3", 2)))

	body = broker.last("PUT /webhooks/webhook-id")
	require.NotNil(t, body)
	assert.Equal(t, "hunter3", body["request"].(map[string]interface{})["password"])
	assert.NotContains(t, state.GoString(), "hunter3")
}

func TestWebhook_PasswordConflictsWithPasswordWriteOnly(t *testing.T) {
	_, brokerSrv := newFakeBroker(t)
	srv, p := testProviderServer(t, brokerSrv.URL)
	ty := p.ResourcesMap["pact_webhook"].CoreConfigSchema().ImpliedType()

	config := objectFromJSON(t, ty, `{
		"description": "hook",
		"events": ["contract_published"],
		"request": [{
			"url": "https://example.com/hook",
			"method": "POST",
			"password": "a",
			"password_wo": "b",
			"password_wo_version": 1,
			"headers": {}
		}]
	}`)

	resp, err := srv.ValidateResourceTypeConfig(context.Background(), &tfprotov5.ValidateResourceTypeConfigRequest{
		TypeName:           "pact_webhook",
		Config:             dynamicValue(t, ty, config),
		ClientCapabilities: &tfprotov5.ValidateResourceTypeConfigClientCapabilities{WriteOnlyAttributesAllowed: true},
	})
	require.NoError(t, err)
	assert.True(t, hasErrorDiag(resp.Diagnostics), "expected a validation error")
}

func TestWebhook_UsernameAndHeadersWriteOnly(t *testing.T) {
	broker, brokerSrv := newFakeBroker(t)
	srv, p := testProviderServer(t, brokerSrv.URL)
	ty := p.ResourcesMap["pact_webhook"].CoreConfigSchema().ImpliedType()

	doc := func(token string, version int) string {
		return `{
			"description": "hook",
			"events": ["contract_published"],
			"request": [{
				"url": "https://example.com/hook",
				"method": "POST",
				"username_wo": "wo-user",
				"username_wo_version": 1,
				"headers": {"Content-Type": "application/json"},
				"headers_wo": "{\"Authorization\": \"Bearer ` + token + `\"}",
				"headers_wo_version": ` + strconv.Itoa(version) + `,
				"body": "{\"pact\":\"url\"}"
			}]
		}`
	}

	state := planAndApply(t, srv, "pact_webhook", ty, cty.NullVal(ty), objectFromJSON(t, ty, doc("t0k3n", 1)))

	request := broker.last("POST /webhooks")["request"].(map[string]interface{})
	assert.Equal(t, "wo-user", request["username"])
	assert.Equal(t, map[string]interface{}{"Content-Type": "application/json", "Authorization": "Bearer t0k3n"}, request["headers"])

	req := state.GetAttr("request").Index(cty.NumberIntVal(0))
	assert.Equal(t, cty.MapVal(map[string]cty.Value{"Content-Type": cty.StringVal("application/json")}), req.GetAttr("headers"))
	assert.Equal(t, "1", req.GetAttr("username_wo_version").AsBigFloat().String())
	assert.Equal(t, "1", req.GetAttr("headers_wo_version").AsBigFloat().String())
	assert.NotContains(t, state.GoString(), "wo-user")
	assert.NotContains(t, state.GoString(), "t0k3n")

	// Refreshing must not pull the write-only values back from the broker into the state
	read, err := srv.ReadResource(context.Background(), &tfprotov5.ReadResourceRequest{
		TypeName:     "pact_webhook",
		CurrentState: dynamicValue(t, ty, state),
	})
	require.NoError(t, err)
	requireNoDiags(t, read.Diagnostics)
	refreshed := decodeDynamicValue(t, ty, read.NewState)
	assert.NotContains(t, refreshed.GoString(), "wo-user")
	assert.NotContains(t, refreshed.GoString(), "t0k3n")
	assert.Equal(t, req.GetAttr("headers"), refreshed.GetAttr("request").Index(cty.NumberIntVal(0)).GetAttr("headers"))

	// Rotate the token
	state = planAndApply(t, srv, "pact_webhook", ty, state, objectFromJSON(t, ty, doc("n3w-t0k3n", 2)))

	request = broker.last("PUT /webhooks/webhook-id")["request"].(map[string]interface{})
	assert.Equal(t, "Bearer n3w-t0k3n", request["headers"].(map[string]interface{})["Authorization"])
	assert.NotContains(t, state.GoString(), "t0k3n")
}

func TestWebhook_WriteOnlyValidation(t *testing.T) {
	_, brokerSrv := newFakeBroker(t)
	srv, p := testProviderServer(t, brokerSrv.URL)
	ty := p.ResourcesMap["pact_webhook"].CoreConfigSchema().ImpliedType()

	doc := func(request string) string {
		return `{"description": "hook", "events": ["contract_published"], "request": [{
			"url": "https://example.com/hook", "method": "POST", "headers": {}, ` + request + `}]}`
	}

	cases := map[string]string{
		"both username and username_wo":  `"username": "a", "username_wo": "b", "username_wo_version": 1`,
		"username_wo without version":    `"username_wo": "b"`,
		"headers_wo without version":     `"headers_wo": "{}"`,
		"headers_wo_version alone":       `"headers_wo_version": 1`,
		"headers_wo is not JSON":         `"headers_wo": "nope", "headers_wo_version": 1`,
		"headers_wo is not a string map": `"headers_wo": "{\"a\": 1}", "headers_wo_version": 1`,
	}

	for name, request := range cases {
		t.Run(name, func(t *testing.T) {
			resp, err := srv.ValidateResourceTypeConfig(context.Background(), &tfprotov5.ValidateResourceTypeConfigRequest{
				TypeName:           "pact_webhook",
				Config:             dynamicValue(t, ty, objectFromJSON(t, ty, doc(request))),
				ClientCapabilities: &tfprotov5.ValidateResourceTypeConfigClientCapabilities{WriteOnlyAttributesAllowed: true},
			})
			require.NoError(t, err)
			assert.True(t, hasErrorDiag(resp.Diagnostics), "expected a validation error")
		})
	}
}

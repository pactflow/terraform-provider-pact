package main

import (
	"encoding/json"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func TestParseWebhookEnabled(t *testing.T) {
	for _, enabled := range []bool{true, false} {
		raw := map[string]interface{}{
			"description": "test",
			"enabled":     enabled,
			"events":      []interface{}{"contract_content_changed"},
			"request": []interface{}{map[string]interface{}{
				"url":    "https://example.com",
				"method": "POST",
			}},
		}
		d := schema.TestResourceDataRaw(t, webhook().Schema, raw)

		w, err := parseWebhook(d, nil)
		if err != nil {
			t.Fatal(err)
		}
		if w.Enabled != enabled {
			t.Errorf("enabled = %v, want %v", w.Enabled, enabled)
		}

		body, err := json.Marshal(w)
		if err != nil {
			t.Fatal(err)
		}
		if got := strings.Contains(string(body), `"enabled":true`); got != enabled {
			t.Errorf("payload %s: enabled true = %v, want %v", body, got, enabled)
		}
		if !enabled && !strings.Contains(string(body), `"enabled":false`) {
			t.Errorf("payload %s must carry enabled:false", body)
		}
	}
}

func TestValidateEventsConcurrent(t *testing.T) {
	events := append([]string(nil), allowedEvents...) // snapshot before goroutines start
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		for _, e := range events {
			wg.Add(1)
			go func(e string) {
				defer wg.Done()
				if _, errs := validateEvents(e, "events.0"); len(errs) > 0 {
					t.Errorf("valid event %q rejected: %v", e, errs)
				}
			}(e)
		}
	}
	wg.Wait()
}

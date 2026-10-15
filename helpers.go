package main

import (
	"sort"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func arrayInterfaceToArrayString(raw []interface{}) []string {
	items := make([]string, len(raw))
	if len(raw) > 0 {
		for i, s := range raw {
			items[i] = s.(string)
		}
	}

	sort.Strings(items)

	return items
}

func interfaceToStringArray(o interface{}) []string {
	items := o.([]interface{})
	res := make([]string, len(items))
	for i, item := range items {
		res[i] = item.(string)
	}

	sort.Strings(res)

	return res
}

// Finds the items in b that don't exist in a
func diff(a, b []string) []string {
	diff := make([]string, 0)
	m := make(map[string]bool)

	for _, item := range a {
		m[item] = true
	}

	for _, item := range b {
		if _, ok := m[item]; !ok {
			diff = append(diff, item)
		}
	}

	return diff
}

// From: https://github.com/hashicorp/terraform-provider-aws/blob/77cbe287f2805319b1c25aa94d70b7a971165f2e/internal/flex/flex.go

// Takes the result of schema.Set of strings and returns a []*string
func ExpandStringSet(configured *schema.Set) []string {
	return ExpandStringList(configured.List()) // nosemgrep: helper-schema-Set-extraneous-ExpandStringList-with-List
}

// Takes the result of flatmap.Expand for an array of strings
// and returns a []*string
func ExpandStringList(configured []interface{}) []string {
	vs := make([]string, 0, len(configured))
	for _, v := range configured {
		val, ok := v.(string)
		if ok && val != "" {
			vs = append(vs, val)
		}
	}
	return vs
}

// rawConfigString returns the string value found at path in the raw resource
// configuration. This is the only way to read write-only attributes, which are
// never present in the plan or state (and hence not available via d.Get).
//
// An empty string is returned if the raw config is not available (e.g. during
// Read or Delete) or if the value at path is null or unknown.
func rawConfigString(d *schema.ResourceData, path cty.Path) (string, diag.Diagnostics) {
	val := d.GetRawConfig()

	for _, step := range path {
		if val.IsNull() || !val.IsKnown() {
			return "", nil
		}

		switch s := step.(type) {
		case cty.GetAttrStep:
			if !val.Type().IsObjectType() || !val.Type().HasAttribute(s.Name) {
				return "", diag.Errorf("unable to read attribute %q from the configuration", s.Name)
			}
			val = val.GetAttr(s.Name)
		case cty.IndexStep:
			if !val.CanIterateElements() {
				return "", diag.Errorf("unable to index into non-collection configuration value")
			}
			if val.LengthInt() == 0 {
				return "", nil
			}
			if hasIndex := val.HasIndex(s.Key); !hasIndex.IsKnown() || hasIndex.False() {
				return "", nil
			}
			val = val.Index(s.Key)
		default:
			return "", diag.Errorf("unsupported configuration path step %T", step)
		}
	}

	if val.IsNull() || !val.IsKnown() {
		return "", nil
	}

	if !val.Type().Equals(cty.String) {
		return "", diag.Errorf("expected a string in the configuration, got %s", val.Type().FriendlyName())
	}

	return val.AsString(), nil
}

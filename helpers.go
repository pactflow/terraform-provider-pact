package main

import (
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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

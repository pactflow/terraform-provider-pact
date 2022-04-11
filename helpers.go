package main

import "sort"

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

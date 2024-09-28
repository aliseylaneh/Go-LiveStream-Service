package appconfigs

import (
	"flag"
	"fmt"
	"strings"
)

var requiredFields []string

func Int(name string, desc string) *int {
	requiredFields = append(requiredFields, name)
	return flag.Int(name, 0, desc)
}

func String(name string, desc string) *string {
	requiredFields = append(requiredFields, name)
	return flag.String(name, "", desc)
}

func Bool(name string, desc string) *bool {
	requiredFields = append(requiredFields, name)
	return flag.Bool(name, false, desc)
}

func Parse() error {
	var unsets []string
	flag.Parse()

	for _, name := range requiredFields {
		if !isFlagPassed(name) {
			unsets = append(unsets, name)
		}
	}

	if len(unsets) > 0 {
		return fmt.Errorf("some required flags are not set: %s", strings.Join(unsets, ", "))
	}

	return nil
}

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

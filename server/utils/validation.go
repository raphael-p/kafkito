package utils

import (
	"errors"
	"regexp"
)

var nameRegex = regexp.MustCompile("^[a-zA-Z0-9_.-]*$")

func CheckNameFormat(name string, errPrefix string) error {
	if !nameRegex.MatchString(name) {
		return errors.New(errPrefix +
			"may only contain letters, numbers, '_', '.', and '-'",
		)
	}
	return nil
}

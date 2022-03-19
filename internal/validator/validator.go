package validator

import (
	"net/url"
	"regexp"
	"strconv"
)

// Validator type which contains a map of validation errors
type Validator struct {
	Errors map[string]string
}

// New is helper which creates a new Validator instance with an empty errors map.
func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

// Valid returns true if the errors map is empty
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

// AddError adds an error message to the map (so long as no entry already exists
// for the given key).
func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

// Check adds an error message to the map if the validation check is not 'ok'.
func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

// In returns true if a specific value is in the list of strings
func In(value string, list ...string) bool {
	for i := range list {
		if value == list[i] {
			return true
		}
	}
	return false
}

// Matches returns true if a string value matches a specific regexp pattern
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

// Unique returns true if all string values in a slice are unique.
func Unique(values []string) bool {
	uniqueValues := make(map[string]bool)

	for _, value := range values {
		uniqueValues[value] = true
	}

	return len(values) == len(uniqueValues)
}

// IsURL checks that the supplied URL has a valid scheme
func IsURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Host != "" && u.Scheme == "http" || u.Scheme == "https"

}

// ReadInt reads string value from the query string and converts it to an integer
// before returning. If no matching key is found it returns the provided default
// value. If the value couldn't be converted to an int, then we record an error
// message in the provided Validator instance.
func ReadInt(qs url.Values, key string, defaultValue int, v *Validator) int {
	s := qs.Get(key)

	if s == "" {
		return defaultValue
	}

	// Try to convert the value to an int. If this fails, add an error message to
	// the validator instance and return the default value.
	i, err := strconv.Atoi(s)
	if err != nil {
		v.AddError(key, "must be an integer value")
		return defaultValue
	}
	return i
}

// ReadString returns a string value from the query string, or the provided
// default value if no matching key is found.
func ReadString(qs url.Values, key string, defaultValue string) string {
	// Extract the value for a given key from the query string. If no key exists
	// this will return the empty string "".
	s := qs.Get(key)

	// If no key exists (or value is empty) then return the default value.
	if s == "" {
		return defaultValue
	}
	return s
}

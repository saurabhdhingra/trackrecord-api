package validator

import (
	"regexp"
	"time"
)

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type Validator struct {
	Errors map[string]string
}

func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

func In(value string, list ...string) bool {
	for i := range list {
		if value == list[i] {
			return true
		}
	}
	return false
}

func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

func Unique(values []string) bool {
	uniqueValues := make(map[string]bool)
	
	for _, value := range values {
		uniqueValues[value] = true
	}
	
	return len(values) == len(uniqueValues)
}

func MinChars(value string, n int) bool {
	return len(value) >= n
}

func MaxChars(value string, n int) bool {
	return len(value) <= n
}

func Min(value, min int) bool {
	return value >= min
}

func Max(value, max int) bool {
	return value <= max
}

func MinFloat(value, min float64) bool {
	return value >= min
}

func MaxFloat(value, max float64) bool {
	return value <= max
}

func IsDate(value time.Time) bool {
	return !value.IsZero()
}

func IsFuture(value time.Time) bool {
	return value.After(time.Now())
} 
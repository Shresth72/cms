package utils

import "unicode/utf8"

type Validator struct {
	Errors map[string]string
}

func (v *Validator) Check(ok bool, key, val string) {
	if !ok {
		v.AddError(key, val)
	}
}

func (v *Validator) AddError(key, val string) {
	if _, fnd := v.Errors[key]; !fnd {
		v.Errors[key] = val
	}
}

func MinChars(field string, n int) bool {
	return utf8.RuneCountInString(field) >= n
}

func MaxChars(field string, n int) bool {
	return utf8.RuneCountInString(field) <= n
}

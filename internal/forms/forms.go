package forms

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

type Form struct {
	url.Values
	Errors errors
}

func (f *Form) Valid() bool {

	return len(f.Errors) == 0
}

// New initialized a form struct

func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Required checks for required fields
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}
func (f *Form) Has(field string) bool {
	x := f.Get(field)
	if x == "" {
		//f.Errors.Add(field, fmt.Sprintf("missing field %s", field))
		return false
	}
	return true

}

// MinLength checks for string minimum length
func (f *Form) MinLength(field string, length int) bool {
	if len(f.Get(field)) < length {
		f.Errors.Add(field, fmt.Sprintf("this field must be at least %d characters long", length))
		return false
	}
	return true
}

// IsEmail checks for valid email address
func (f *Form) IsEmail(field string) {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field, "Invalid email address.")
	}
}

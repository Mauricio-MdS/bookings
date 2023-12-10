package forms

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

// Form creates a custom form struct, embeds a url.Values object
type Form struct {
	url.Values
	Errors errors
}

// Has checks if form field is in post and not empty
func (f *Form) Has(field string) bool {
	 x := f.Get(field)
	 return x != ""
}

// IsEmail checks for valid email address. Add error if invalid
func (f *Form) IsEmail(field string) {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field, "Invalid email address")
	}
}

// MinLength checks for field minimun length
func (f *Form) MinLength(fieldName string, length int) bool {
	field := f.Get(fieldName)
	if len(field) < length {
		f.Errors.Add(fieldName, fmt.Sprintf("This field must be at least %d characters long", length))
		return false
	}
	return true
}

// New initializes a Form struct
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Required add errors for blank provided fields
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

// Valid returns true if there are no errors, otherwise false
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}
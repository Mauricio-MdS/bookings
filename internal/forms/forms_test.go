package forms

import (
	"net/url"
	"testing"
)

const (
	blankField = "blank field"
	nonExistantField = "invalid field"
	randomString = "whatever"
	smallString = "s"	
	validEmail = "test@testing.com"
)

func TestForm_Has(t *testing.T) {
	data := url.Values{}
	data.Add(validEmail, validEmail)
	data.Add(blankField, "")
	form := New(data)
	if !form.Has(validEmail) {
		t.Error("form has field, but Has function returns false")
	}
	if form.Has(nonExistantField) {
		t.Error("form does not has field, but Has function returns true")
	}
	if form.Has(blankField) {
		t.Error("form field is blank, but Has function returns true")
	}
}

func TestForm_IsEmail(t *testing.T) {
	data := url.Values{}
	data.Add(validEmail, validEmail)
	data.Add(randomString, randomString)
	form := New(data)
	form.IsEmail(validEmail)
	if !form.Valid() {
		t.Error("shows invalid email for a valid email")
	}
	form.IsEmail(randomString)
	if form.Valid() {
		t.Error("shows valid email for a invalid email")
	}
}

func TestForm_MinLength(t *testing.T) {

	data := url.Values{}
	data.Add(smallString, smallString)
	data.Add(randomString, randomString)
	form := New(data)
	if form.MinLength(smallString, 3) {
		t.Error("shows minimun length for a small string")
	}
	if !form.MinLength(randomString, 3) {
		t.Error("shows that does not satisty minimum length for a passing string")
	}
}

func TestForm_Required(t *testing.T) {
	data := url.Values{}
	data.Add(validEmail, validEmail)
	data.Add(randomString, randomString)
	data.Add(blankField, "")
	form := New(data)
	form.Required(validEmail, randomString)
	if !form.Valid() {
		t.Error("form shows invalid when required fields are present")
	}
	form.Required(nonExistantField)
	if form.Valid() {
		t.Error("form shows valid when required fields are missing")
	}
}

func TestForm_Valid(t *testing.T) {
	form := New(url.Values{})
	isValid := form.Valid()
	if !isValid {
		t.Error("got invalid for a form with no errors")
	}
}
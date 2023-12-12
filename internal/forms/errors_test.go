package forms

import "testing"

const(
	errorKey = "error key"
	firstMessage = "first message"
	secondMessage = "second message"
)

func TestErrors_Get(t *testing.T) {
	err := errors{}
	if err.Get(errorKey) != "" {
		t.Error("does not return blank string for empty errors")
	}
	err.Add(errorKey, firstMessage)
	err.Add(errorKey, secondMessage)
	if err.Get(errorKey) != firstMessage {
		t.Error("does not return first error message")
	}
}
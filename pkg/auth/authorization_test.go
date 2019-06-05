package auth

import (
	"errors"
	"testing"
)

func TestAuthorizeAdminUser(t *testing.T) {
	testcases := []struct {
		description string
		session     *Session
		expected    error
	}{
		{
			description: "authorized",
			session:     &Session{IsSuperuser: true},
			expected:    nil,
		},
		{
			description: "not authorized",
			session:     &Session{},
			expected:    errors.New("USER_UNAUTHORIZED"),
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.description, func(t *testing.T) {
			result := AuthorizeAdminUser(testcase.session)
			expected := testcase.expected

			var failed bool
			// check that the type returned is an error
			_, ok := result.(error)

			if ok {
				failed = result.Error() != expected.Error()
			} else {
				failed = result != expected
			}

			if failed {
				t.Errorf("got %#v, expected %#v", result, testcase.expected)
			}
		})
	}
}

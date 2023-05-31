package shortener

import (
	"fmt"
	"testing"
)

func TestValidateURL(t *testing.T) {
	testCases := []struct {
		url  string
		want bool
	}{
		{"https://xxx.com", true},
		{"https://xxx.com/yyy", true},
		{"https://", false},
		{"xxx.com", false},
		{"https://xxx", false},
	}

	shortener := SimpleShortener{}

	for _, test := range testCases {
		t.Run(fmt.Sprintf("validate %s", test.url), func(t *testing.T) {
			got, err := shortener.ValidateURL(test.url)
			if err != nil {
				t.Fatal(err)
			}

			if got != test.want {
				t.Errorf("got %v, want %v", got, test.want)
			}
		})
	}
}

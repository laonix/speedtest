package provider

import (
	"reflect"
	"testing"
)

func TestProviderFactory_GetProvider_typesAssertion(t *testing.T) {
	f := ProviderFactory{}

	tests := []struct {
		name     string
		input    string
		expected reflect.Type
	}{
		{
			name:     "ookla",
			input:    "ookla",
			expected: reflect.TypeOf(&Ookla{}),
		},
		{
			name:     "netflix",
			input:    "netflix",
			expected: reflect.TypeOf(&Netflix{}),
		},
		{
			name:     "unknown provider",
			input:    "some",
			expected: reflect.TypeOf(&NoopProvider{}),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, _ := f.GetProvider(test.input); reflect.TypeOf(got) != test.expected {
				t.Errorf("Wrong provider: got = %s, expected = %s", reflect.TypeOf(got), test.expected)
			}
		})
	}
}

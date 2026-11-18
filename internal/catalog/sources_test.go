package catalog

import (
	"reflect"
	"testing"
)

func TestSources(t *testing.T) {
	tests := []struct {
		input string
		want  []string
	}{
		{"", []string{"https://libgen.bz/index.php"}},
		{" Libgen.bz, books.example.org, LIBGEN.BZ ", []string{"https://libgen.bz/index.php", "https://books.example.org/search.php"}},
		{"https://books.xyz", nil},
		{"invalid-domain", nil},
		{"books.xyz:443", nil},
		{"-books.xyz", nil},
		{"books..xyz", nil},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			sources, err := Sources(test.input)
			if (err != nil) != (test.want == nil) {
				t.Fatalf("Sources(%q) error = %v", test.input, err)
			}
			var got []string
			for _, source := range sources {
				got = append(got, source.String())
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Fatalf("sources = %v, want %v", got, test.want)
			}
		})
	}
}

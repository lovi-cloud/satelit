package team_test

import (
	"testing"

	"github.com/go-test/deep"
	"github.com/whywaita/satelit/pkg/api/team"
)

func TestGetTeamID(t *testing.T) {
	tests := []struct {
		input string
		want  uint
		err   bool
	}{
		{
			input: "10.160.0.101",
			want:  0,
			err:   true,
		},
		{
			input: "10.160.1.101",
			want:  1,
			err:   false,
		},
		{
			input: "10.160.11.104",
			want:  11,
			err:   false,
		},
		{
			input: "10.165.76.103",
			want:  576,
			err:   false,
		},
		{
			input: "10.165.77.101",
			want:  0,
			err:   true,
		},
	}
	for _, test := range tests {
		got, err := team.GetTeamID(test.input)
		if !test.err && err != nil {
			t.Fatalf("should not be error for %+v but: %+v", test.input, err)
		}
		if test.err && err == nil {
			t.Fatalf("should be error for %+v but not:", test.input)
		}
		if diff := deep.Equal(test.want, got); len(diff) != 0 {
			t.Fatalf("want %v, but %v, diff %q:", test.want, got, diff)
		}
	}
}

package join

import (
	"bytes"
	"strings"
	"testing"
)

func TestJoinTwoInputs(t *testing.T) {
	tests := []struct {
		intent, left, right, want string
	}{
		{"empty inputs", "", "", ""},
		{"empty lines", "\n", "\n", ""},
		{"first lines match", "1 a", "1 b", "1 a b\n"},
		{"no match", "1 a", "2 b", ""},
		{"left[1] matches right[0]", "1 a\n2 b", "2 c", "2 b c\n"},
		{"left[0] matches right[1]", "2 a", "1 b\n2 c", "2 a c\n"},
		{"first lines empty, second ones match", "\n1 a", "\n1 b", "1 a b\n"},
		{"matching lines with single word", "1", "1", "1\n"},
		{"all lines match", "1 a\n2 b", "1 c\n2 d", "1 a c\n2 b d\n"},
		{"last line has end-of-line", "1 a\n", "1 b\n", "1 a b\n"},
		{"inputs with differing # of lines", "1 a", "1 b\n2 c", "1 a b\n"},
	}
	for _, test := range tests {
		inputs := [2]Input{
			NewInput(strings.NewReader(test.left), "left"),
			NewInput(strings.NewReader(test.right), "right"),
		}
		var output bytes.Buffer
		err := Join(inputs[:], " ", &output)
		if err != nil {
			t.Errorf("%s: unexpected error: %v", test.intent, err)
		}
		got := output.String()
		if got != test.want {
			t.Errorf("%s: want: %q got: %q", test.intent, test.want, got)
		}
	}
}

func TestJoinThreeInputs(t *testing.T) {
	// Most cases are exercised in TestJoinTwoInputs().
	tests := []struct {
		intent, left, mid, right, want string
	}{
		{"all inputs match", "1 a", "1 b", "1 c", "1 a b c\n"},
		{"only some input match", "1 a", "2 b", "1 c", ""},
	}
	for _, test := range tests {
		inputs := [3]Input{
			NewInput(strings.NewReader(test.left), "left"),
			NewInput(strings.NewReader(test.mid), "mid"),
			NewInput(strings.NewReader(test.right), "right"),
		}
		var output bytes.Buffer
		err := Join(inputs[:], " ", &output)
		if err != nil {
			t.Errorf("%s: unexpected error: %v", test.intent, err)
		}
		got := output.String()
		if got != test.want {
			t.Errorf("%s: want: %q got: %q", test.intent, test.want, got)
		}
	}
}

func TestNonDefaultSeparator(t *testing.T) {
	inputs := [2]Input{
		NewInput(strings.NewReader("foo|bar"), "left"),
		NewInput(strings.NewReader("foo|baz|zap"), "right"),
	}
	var output bytes.Buffer
	err := Join(inputs[:], "|", &output)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	got := output.String()
	want := "foo|bar|baz|zap\n"
	if got != want {
		t.Errorf("want: %q got: %q", want, got)
	}
}

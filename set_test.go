package trig

import (
	"reflect"
	"testing"
)

func TestCleanSlice(t *testing.T) {
	tests := []struct {
		in  []string
		v   string
		out []string
	}{
		{[]string{"a", "b", "c", "d", "e"}, "c", []string{"a", "b", "d", "e"}},
		{[]string{"a", "b", "c", "d", "e"}, "e", []string{"a", "b", "c", "d"}},
		{[]string{"a", "b", "c", "d", "e"}, "a", []string{"b", "c", "d", "e"}},
		{[]string{"a", "b", "c", "d", "e"}, "x", []string{"a", "b", "c", "d", "e"}},
		{[]string{"a", "b", "c", "d", "e"}, "", []string{"a", "b", "c", "d", "e"}},
		{[]string{"a"}, "b", []string{"a"}},
		{[]string{"a"}, "a", []string{}},
		{[]string{}, "a", []string{}},
		{[]string(nil), "a", []string(nil)},
	}
	for i, test := range tests {
		s := cleanSlice(test.in, test.v)
		if !reflect.DeepEqual(s, test.out) {
			t.Errorf("#%d want %#v, got %#v", i, test.out, s)
		}
	}
}

func TestInSlice(t *testing.T) {
	tests := []struct {
		in  []string
		v   string
		out bool
	}{
		{[]string{"a", "b", "c", "d", "e"}, "c", true},
		{[]string{"a", "b", "c", "d", "e"}, "e", true},
		{[]string{"a", "b", "c", "d", "e"}, "a", true},
		{[]string{"a", "b", "c", "d", "e"}, "x", false},
		{[]string{"a", "b", "c", "d", "e"}, "", false},
		{[]string{"a"}, "b", false},
		{[]string{"a"}, "a", true},
		{[]string{}, "a", false},
		{[]string(nil), "a", false},
	}
	for i, test := range tests {
		r := isInSlice(test.in, test.v)
		if test.out != r {
			t.Errorf("#%d want %#v, got %#v", i, test.out, r)
		}
	}

}

func TestAppendIfNecessary(t *testing.T) {
	tests := []struct {
		in  []string
		v   string
		out []string
	}{
		{[]string{"a", "b", "c", "d", "e"}, "c", []string{"a", "b", "c", "d", "e"}},
		{[]string{"a", "b", "c", "d", "e"}, "x", []string{"a", "b", "c", "d", "e", "x"}},
		{[]string{"a", "b", "c", "d", "e"}, "", []string{"a", "b", "c", "d", "e"}},
		{[]string{"a"}, "b", []string{"a", "b"}},
		{[]string{"a"}, "a", []string{"a"}},
		{[]string{}, "a", []string{"a"}},
		{[]string(nil), "a", []string{"a"}},
	}
	for i, test := range tests {
		s := appendIfNecessary(test.in, test.v)
		if !reflect.DeepEqual(s, test.out) {
			t.Errorf("#%d want %#v, got %#v", i, test.out, s)
		}
	}
}

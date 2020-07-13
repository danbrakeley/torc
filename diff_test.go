package main

import (
	"reflect"
	"testing"
)

func Test(t *testing.T) {
	cases := []struct {
		Name          string
		A             []string
		B             []string
		AOnlyExpected []string
		BOnlyExpected []string
	}{
		{"both nil", nil, nil, []string{}, []string{}},
		{"both empty", []string{}, []string{}, []string{}, []string{}},
		{"a nil", nil, []string{"b"}, []string{}, []string{"b"}},
		{"a empty", []string{}, []string{"b"}, []string{}, []string{"b"}},
		{"b nil", []string{"a"}, nil, []string{"a"}, []string{}},
		{"b empty", []string{"a"}, []string{}, []string{"a"}, []string{}},
		{"no diff", []string{"aaaa", "b", "cd"}, []string{"aaaa", "b", "cd"}, []string{}, []string{}},
		{"diff in a only", []string{"aaaa", "b", "cd"}, []string{"b", "cd"}, []string{"aaaa"}, []string{}},
		{"diff in b only", []string{"q", "ran"}, []string{"bob", "q", "ran"}, []string{}, []string{"bob"}},
		{"diff in both", []string{"m", "xyz"}, []string{"pork", "xyz"}, []string{"m"}, []string{"pork"}},
		{
			"multiple diffs",
			[]string{"a", "franklin mint", "m", "xyz"},
			[]string{"a", "barry", "berry", "m", "zzz"},
			[]string{"franklin mint", "xyz"},
			[]string{"barry", "berry", "zzz"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {

			aOnlyActual, bOnlyActual := DiffSortedSlices(tc.A, tc.B)
			if !reflect.DeepEqual(aOnlyActual, tc.AOnlyExpected) {
				t.Errorf("Expected A-only strings to be %v, not %v", tc.AOnlyExpected, aOnlyActual)
			}
			if !reflect.DeepEqual(bOnlyActual, tc.BOnlyExpected) {
				t.Errorf("Expected B-only strings to be %v, not %v", tc.BOnlyExpected, bOnlyActual)
			}
		})
	}
}

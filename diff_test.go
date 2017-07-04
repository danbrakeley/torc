package main

import (
	"reflect"
	"testing"
)

func testDiffSortedSlices(t *testing.T, a, b, aOnlyExpected, bOnlyExpected []string) {
	aOnly, bOnly := DiffSortedSlices(a, b)
	if !reflect.DeepEqual(aOnly, aOnlyExpected) {
		t.Errorf("Expected aOnly to be %v, not %v", aOnlyExpected, aOnly)
	}
	if !reflect.DeepEqual(bOnly, bOnlyExpected) {
		t.Errorf("Expected bOnly to be %v, not %v", bOnlyExpected, bOnly)
	}
}

func TestDiffSortedSlices_BothNil(t *testing.T) {
	testDiffSortedSlices(t, nil, nil, []string{}, []string{})
}

func TestDiffSortedSlices_BothEmpty(t *testing.T) {
	testDiffSortedSlices(t, []string{}, []string{}, []string{}, []string{})
}

func TestDiffSortedSlices_ANil(t *testing.T) {
	b := []string{"a"}
	testDiffSortedSlices(t, nil, b, []string{}, b)
}
func TestDiffSortedSlices_AEmpty(t *testing.T) {
	b := []string{"a"}
	testDiffSortedSlices(t, []string{}, b, []string{}, b)
}

func TestDiffSortedSlices_BNil(t *testing.T) {
	a := []string{"mon key blah", "qqghad"}
	testDiffSortedSlices(t, a, nil, a, []string{})
}

func TestDiffSortedSlices_BEmpty(t *testing.T) {
	a := []string{"mon key blah", "qqghad"}
	testDiffSortedSlices(t, a, []string{}, a, []string{})
}

func TestDiffSortedSlices_NoDiff(t *testing.T) {
	a := []string{"aaaa", "b", "cd"}
	testDiffSortedSlices(t, a, a, []string{}, []string{})
}

func TestDiffSortedSlices_DiffInAOnly(t *testing.T) {
	a := []string{"aaaa", "b", "cd"}
	b := []string{"b", "cd"}
	aOnlyExpected := []string{"aaaa"}
	bOnlyExpected := []string{}
	testDiffSortedSlices(t, a, b, aOnlyExpected, bOnlyExpected)
}

func TestDiffSortedSlices_DiffInBOnly(t *testing.T) {
	a := []string{"q", "ran"}
	b := []string{"bob", "q", "ran"}
	aOnlyExpected := []string{}
	bOnlyExpected := []string{"bob"}
	testDiffSortedSlices(t, a, b, aOnlyExpected, bOnlyExpected)
}

func TestDiffSortedSlices_DiffInBoth(t *testing.T) {
	a := []string{"m", "xyz"}
	b := []string{"pork", "xyz"}
	aOnlyExpected := []string{"m"}
	bOnlyExpected := []string{"pork"}
	testDiffSortedSlices(t, a, b, aOnlyExpected, bOnlyExpected)
}

func TestDiffSortedSlices_MultipleDiffs(t *testing.T) {
	a := []string{"a", "franklin mint", "m", "xyz"}
	b := []string{"a", "barry", "berry", "m", "zzz"}
	aOnlyExpected := []string{"franklin mint", "xyz"}
	bOnlyExpected := []string{"barry", "berry", "zzz"}
	testDiffSortedSlices(t, a, b, aOnlyExpected, bOnlyExpected)
}

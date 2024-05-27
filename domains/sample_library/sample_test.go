package sample_library

import "testing"

func TestSampleFunction(t *testing.T) {
	expected := "Hello, Akehilesh!"
	actual := SampleGoFunction("Akehilesh")

	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

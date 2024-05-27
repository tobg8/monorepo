package sample_library

import "testing"

func TestSampleFunction(t *testing.T) {
	expected := "Hello, Akhilesh!"
	actual := SampleGoFunction("Akhilesh")

	if actual!= expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}
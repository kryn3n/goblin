package arcgis

import (
	"reflect"
	"testing"
)

func TestCreateBatches(t *testing.T) {
	limit := 2
	objectIds := []int{0, 1, 2, 3, 4}
	result := createBatches(objectIds, limit)

	expected := [][]int{{0, 1}, {2, 3}, {4}}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Result was incorrect, got: %v, want: %v", result, expected)
	}
}

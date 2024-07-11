package formula 

import (
	"testing"
	"gotest.tools/v3/assert"
)

// func TestDominanentLoad(t *testing.T){
// 	resourceType := [2]string{
// 		"CPU", "Memory",
// 	};
// 	// app_resource := map[string]int{
// 	// 	"CPU" : 10,
// 	// 	"Memory" : 20,
// 	// };
// 	nodeAvailResource := map[string]float64{
// 		"CPU":    8,
// 		"Memory": 30,
// 	};
// 	nodeCapResource := map[string]float64{
// 		"CPU":    10,
// 		"Memory": 50,
// 	};
// 	dominanentLoad := GetDominantLoad(nodeAvailResource, nodeCapResource, resourceType[:]);
// 	expected := 0.3
// 	if dominanentLoad != expected {
// 		t.Errorf("Got: %v, Expected: %v", dominanentLoad, expected)
// 	}
// }
func TestGetMean (t *testing.T) {
	objects := [...]float64{-1.0, 0.0, 1.0}
	mean := GetMean(objects[:])
	expect := 0.0
	assert.Equal(t, mean, expect)
}

func TestGetStdDev (t *testing.T) {
	objects := [...]float64{-1.0, 0.0, 1.0}
	mean := GetMean(objects[:])
	stdDev := GetStdDev(objects[:], mean)
	expect := 0.816496580927726
	assert.Equal(t, stdDev, expect)
}

func TestGetStandardizationScore (t *testing.T) {
	object := 0.5
	objects := [...]float64{-1.0, 0.0, 1.0}
	mean := GetMean(objects[:])
	stdDev := GetStdDev(objects[:], mean)
	standardizationScore := GetStandardizationScore(object, mean, stdDev)
	expect := 0.6123724356957945
	assert.Equal(t, standardizationScore, expect)
}
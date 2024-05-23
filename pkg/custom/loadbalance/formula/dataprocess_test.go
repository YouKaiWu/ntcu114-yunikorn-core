package formula 

import (
	"testing"
	"gotest.tools/v3/assert"
)


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
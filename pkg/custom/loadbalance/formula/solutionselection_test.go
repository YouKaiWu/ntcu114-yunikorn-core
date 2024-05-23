package formula

import (
	"testing"
	"gotest.tools/v3/assert"
)


func TestGetBestSol(t *testing.T) {
	meansOfObjects := [...]float64{1, 2, 3}
	stdDevsOfObjectsOfObjects := [...]float64{1, 2, 3}
	minMean, minStdDev := GetBestSol(meansOfObjects[:], stdDevsOfObjectsOfObjects[:])
	expectMinMean := 1.0
	expectMinStdDev := 1.0
	assert.Equal(t, minMean, expectMinMean)
	assert.Equal(t, minStdDev, expectMinStdDev)
}

func TestGetWorstSol(t *testing.T) {
	meansOfObjects := [...]float64{1, 2, 3}
	stdDevsOfObjects := [...]float64{1, 2, 3}
	maxMean, maxStdDev := GetWorstSol(meansOfObjects[:], stdDevsOfObjects[:])
	expectMaxMean := 3.0
	expectMaxStdDev := 3.0
	assert.Equal(t, maxMean, expectMaxMean)
	assert.Equal(t, maxStdDev, expectMaxStdDev)
}

func TestGetEuclideanDistance(t *testing.T) {
	x := [...]float64{4, 5, 6, 7}
	y := [...]float64{1, 2, 3, 4}
	distance := GetEuclideanDistance(x[:], y[:])
	expect := 6
	assert.Equal(t, distance, expect)
}

func TestGetRCVal(t *testing.T) {
	best := 1.0
	worst := 4.0
	rc := GetRCVal(best, worst)
	expect := 0.8
	assert.Equal(t, rc, expect)
}
package formula

import (
	"math"
)

func Normalized(q []float64) []float64 {
	result := make([]float64, 0)

	sum := float64(0)
	for _, element := range q {
		sum += math.Pow(element, float64(2))
	}
	base := math.Sqrt(sum)

	for _, element := range q {
		tmp := element / base
		result = append(result, tmp)
	}
	return result
}

func Weight(normalizedValues []float64, objectNames []string) []float64 {
	result := make([]float64, 0)
	w := float64(len(objectNames))
	for _, value := range normalizedValues {
		result = append(result, (value / w))
	}
	return result
}

func APlus(q []float64) float64 {
	min := q[0]
	for _, element := range q {
		if min > element {
			min = element
		}
	}
	return min
}

func AMinus(q []float64) float64 {
	max := q[0]
	for _, element := range q {
		if max < element {
			max = element
		}
	}
	return max
}

func APlusOfUsages(q []float64) float64 {
	return AMinus(q)
}

func AMinusOfUsages(q []float64) float64 {
	return APlus(q)
}

// This fomula only calculate one of MIG or deviation distance
// For example, there are n MIGs and n A(+/-), this gomula would caculate
// If develeoper would add a new objetive, developer
func SM(weighted [][]float64, AObjects []float64) []float64 {
	result := make([]float64, len(weighted[0]))
	for objectiveType, objective := range weighted {
		objectValue := AObjects[objectiveType]
		for index, value := range objective {
			result[index] += math.Pow(value-objectValue, float64(2))
		}
	}
	for index, sum := range result {
		result[index] = math.Sqrt(sum)
	}
	return result
}

func IndexOfMaxRC(SMPlus, SMMinus []float64) (int, []float64) {
	index := 0
	max := float64(0)
	debug := make([]float64, 0)
	for i, value := range SMMinus {
		base := SMPlus[i] + value
		tmp := value / base
		if tmp > max {
			max = tmp
			index = i
		}
		debug = append(debug, tmp)
	}
	return index, debug
}
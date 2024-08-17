package apps

import (
	"container/heap"
	"testing"
	"time"
)

func TestApps(t *testing.T) {
	type input struct {
		id             string
		submissionTime time.Time
	}
	tests := []struct {
		caseName string
		apps     []input
		expected []string
	}{
		{
			"normal",
			[]input{
				{"application-01", time.Date(2020, time.November, 10, 23, 0, 0, 0, time.UTC)},
				{"application-02", time.Date(2009, time.November, 11, 23, 0, 0, 0, time.UTC)},
				{"application-03", time.Date(2010, time.November, 10, 23, 0, 0, 0, time.UTC)},
			},
			[]string{"application-02", "application-03", "application-01"},
		},
		{
			"same time",
			[]input{
				{"application-02", time.Date(2020, time.November, 10, 23, 0, 0, 0, time.UTC)},
				{"application-01", time.Date(2020, time.November, 10, 23, 0, 0, 0, time.UTC)},
				{"application-03", time.Date(2020, time.November, 10, 23, 0, 0, 0, time.UTC)},
			},
			[]string{"application-01", "application-02", "application-03"},
		},
	}
	for _, test := range tests {
		t.Run(test.caseName, func(t *testing.T) {
			h := NewApps()
			for _, element := range test.apps {
				heap.Push(h, NewApp(element.id, element.submissionTime, ""))
			}

			if len(test.expected) != h.Len() {
				t.Errorf("expected len %d,got %d", len(test.expected), h.Len())
			} else {
				for i := 0; h.Len() > 0; i++ {
					tmp := heap.Pop(h).(*App)
					if tmp.Id != test.expected[i] {
						t.Errorf("Expected app: %s, got %s", tmp.Id, test.expected[i])
					}
				}
			}
		})
	}
}
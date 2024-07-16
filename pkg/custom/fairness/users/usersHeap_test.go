package users

import (
	"container/heap"
	"testing"
)

func TestUsersHeap(t *testing.T) {
	type input struct {
		username string
		dominantResourceShare      float64
	}
	tests := []struct {
		caseName string
		users    []input
		expected []string
	}{
		{
			"normal",
			[]input{
				{"user-01", 0.2},
				{"user-02", 0.1},
				{"user-03", 0.3},
			},
			[]string{"user-02", "user-01", "user-03"},
		},
		{
			"same drs",
			[]input{
				{"user-02", 0.1},
				{"user-01", 0.1},
				{"user-03", 0.1},
			},
			[]string{"user-01", "user-02", "user-03"},
		},
	}
	for _, test := range tests {
		t.Run(test.caseName, func(t *testing.T) {
			h := NewUsersHeap()
			for _, element := range test.users {
				heap.Push(h, &UserInfo{element.username, element.dominantResourceShare, "vcore"})
			}

			if len(test.expected) != h.Len() {
				t.Errorf("expected len %d, got %d", len(test.expected), h.Len())
			} else {
				for i := 0; h.Len() > 0; i++ {
					tmp := heap.Pop(h).(*UserInfo)
					if tmp.username != test.expected[i] {
						t.Errorf("Expected user: %s, got %s", test.expected[i], tmp.username)
					}
				}
			}
		})
	}
}

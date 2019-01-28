package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type AllCombinationsData struct {
	Keys    []string
	Results [][]string
}

func TestAllCombinaTions(t *testing.T) {
	allCombinationsData := []AllCombinationsData{
		AllCombinationsData{
			Keys: []string{"v1"},
			Results: [][]string{
				[]string{"v1"},
			},
		},
		AllCombinationsData{
			Keys: []string{"v1", "v2"},
			Results: [][]string{
				[]string{"v1"},
				[]string{"v2"},
				[]string{"v1", "v2"},
			},
		},
		AllCombinationsData{
			Keys: []string{"v1", "v2", "v3"},
			Results: [][]string{
				[]string{"v1"},
				[]string{"v2"},
				[]string{"v1", "v2"},
				[]string{"v3"},
				[]string{"v1", "v3"},
				[]string{"v2", "v3"},
				[]string{"v1", "v2", "v3"},
			},
		},
		AllCombinationsData{
			Keys: []string{"v1", "v2", "v3", "v4"},
			Results: [][]string{
				[]string{"v1"},
				[]string{"v2"},
				[]string{"v1", "v2"},
				[]string{"v3"},
				[]string{"v1", "v3"},
				[]string{"v2", "v3"},
				[]string{"v1", "v2", "v3"},
				[]string{"v4"},
				[]string{"v1", "v4"},
				[]string{"v2", "v4"},
				[]string{"v1", "v2", "v4"},
				[]string{"v3", "v4"},
				[]string{"v1", "v3", "v4"},
				[]string{"v2", "v3", "v4"},
				[]string{"v1", "v2", "v3", "v4"},
			},
		},
	}

	for _, data := range allCombinationsData {
		actual := allCombinations(data.Keys, make([][]string, 0))
		assert.Equal(t, data.Results, actual)

	}

}

package autocomplete

import (
	"slices"
	"testing"
)

func TestUtils(t *testing.T) {
	t.Run("test PATH executable retrieval", func(t *testing.T) {
		ret := GetExecutablesFromPath()

		defaultExecutables := map[string]int{
			"ls":   1,
			"echo": 1,
			"cat":  1,
		}

		for cmd := range defaultExecutables {
			if slices.Index(ret, cmd) == -1 {
				t.Errorf("%s should be a built-in executable but is not found in the list of executables", cmd)
			}
		}
	})
}

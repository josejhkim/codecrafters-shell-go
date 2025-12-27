package autocomplete

import (
	"os"
	"strings"
)

func GetExecutablesFromPath() []string {
	ret := []string{}

	pathString := os.Getenv("PATH")
	pathDirs := strings.Split(pathString, string(os.PathListSeparator))

	for _, pathDir := range pathDirs {
		entries, err := os.ReadDir(pathDir)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if !entry.IsDir() {
				ret = append(ret, entry.Name())
			}
		}
	}
	return ret
}

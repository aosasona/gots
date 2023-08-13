package gots

import (
	"bytes"
	"os"
)

func (g *gots) areSameBytesContent(out string) bool {
	stat, err := os.Stat(g.config.OutputFileOrDefault())
	if err != nil {
		return false
	}

	if stat.IsDir() {
		return false
	}
	file, err := os.ReadFile(g.config.OutputFileOrDefault())
	if err != nil {
		return false
	}

	if bytes.Equal(file, []byte(out)) {
		return true
	}

	return false
}

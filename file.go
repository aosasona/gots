package gots

import (
	"bytes"
	"fmt"
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

func (g *gots) exportToFile(ts string) error {
	out := fmt.Sprintf(`/*
* This file is auto-generated and modified by Gots (https://github.com/aosasona/gots).
* DO NOT MODIFY THE CONTENT OF THIS FILE
*/

%s`, ts)

	// if the file already exists and the content is the same, don't write to it
	if g.areSameBytesContent(out) {
		return nil
	}

	err := os.WriteFile(g.config.OutputFileOrDefault(), []byte(out), 0644)
	if err != nil {
		return err
	}

	return nil
}

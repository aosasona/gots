package gots

import (
	"testing"

	"github.com/aosasona/gots/config"
)

func Test_Fork(t *testing.T) {
	original := Init(config.Config{})

	fork := original.Fork(config.Config{
		Enabled:           Bool(true),
		OutputFile:        String("test.ts"),
		UseTypeForObjects: Bool(true),
	}, false)

	if fork.config.EnabledOrDefault() != true {
		t.Errorf("Expected fork.config.EnabledOrDefault() to be true, got %v", fork.config.EnabledOrDefault())
	}

	if fork.config.OutputFileOrDefault() != "test.ts" {
		t.Errorf("Expected fork.config.OutputFileOrDefault() to be \"test.ts\", got %v", fork.config.OutputFileOrDefault())
	}

	if fork.config.UseTypeForObjectsOrDefault() != true {
		t.Errorf("Expected fork.config.UseTypeForObjectsOrDefault() to be true, got %v", fork.config.UseTypeForObjectsOrDefault())
	}

	if fork.config.ExpandObjectTypesOrDefault() != false {
		t.Errorf("Expected fork.config.ExpandObjectTypesOrDefault() to be false, got %v", fork.config.ExpandObjectTypesOrDefault())
	}
}

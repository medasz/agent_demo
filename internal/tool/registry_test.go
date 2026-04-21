package tool

import "testing"

func TestRegistryDefinitions(t *testing.T) {
	registry := NewRegistry()
	definitions := registry.Definitions()

	if len(definitions) != 5 {
		t.Fatalf("Definitions() length = %d, want %d", len(definitions), 5)
	}

	names := make(map[string]bool)
	for _, definition := range definitions {
		if definition.Function == nil {
			t.Fatal("Definitions() contains tool without function definition")
		}
		names[definition.Function.Name] = true
	}

	for _, name := range []string{"read_file", "list_file", "edit_file", "create_file", "run_command"} {
		if !names[name] {
			t.Fatalf("Definitions() missing %q", name)
		}
	}
}

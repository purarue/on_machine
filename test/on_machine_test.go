package on_machine_test

import (
	"fmt"
	"github.com/purarue/on_machine"
	"testing"
)

func TestPatterns(t *testing.T) {
	res := on_machine.ReplaceFields("no fields here")
	if res != "no fields here" {
		t.Errorf("Expected 'no fields here', got '%s'\n", res)
	}
	res = on_machine.ReplaceFields("")
	if res != "" {
		t.Errorf("Expected empty output, got '%s'\n", res)
	}
	res = on_machine.ReplaceFields("%o/")
	if res != on_machine.GetOS()+"/" {
		t.Errorf("Expected '%s/', got res", on_machine.GetOS())
	}
	res = on_machine.ReplaceFields("s")
	if res != "s" {
		t.Errorf("Expected 's' got '%s'\n", res)
	}
	res = on_machine.ReplaceFields("%o")
	if res != on_machine.GetOS() {
		t.Errorf("Expected '%s', got '%s'\n", on_machine.GetOS(), res)
	}
	res = on_machine.ReplaceFields("ff - %o")
	if res != "ff - "+on_machine.GetOS() {
		t.Errorf("Expected 'ff - %s', got '%s'\n", on_machine.GetOS(), res)
	}
	expected := fmt.Sprintf("on_machine_%s_%s_%s_%s_%s", on_machine.GetOS(), on_machine.GetDistro(), on_machine.GetHostname(), on_machine.GetGolangArch(), on_machine.GetGolangOS())
	res = on_machine.ReplaceFields("on_machine_%o_%d_%h_%a_%O")
	if res != expected {
		t.Errorf("Expected '%s', got '%s'\n", expected, res)
	}
}

package on_machine

import (
	"os/exec"
	"strings"
)

func LsbReleaseSh() (*string, error) {
	result, err, _ := Cache.Memoize("lsb_release", func() (interface{}, error) {
		path, err := exec.LookPath("lsb_release")
		if err != nil {
			return nil, err
		}
		cmd := exec.Command(path, "-si")
		out, err := cmd.Output()
		if err != nil {
			return nil, err
		}
		// remove double quotes from beginning/end
		outStr := strings.ToLower(strings.Trim(strings.TrimSpace(string(out)), "\""))
		return &outStr, nil
	})
	if lsbRes, ok := result.(*string); ok {
		return lsbRes, nil
	} else {
		return nil, err
	}
}

package monitor

import (
	"bytes"
	"encoding/json"
	"os/exec"
)

func ExecPuppeteer(urls []string) (*map[string]string, error) {
	contentMap := make(map[string]string, len(urls))

	// ensure urls are properly quoted
	args := make([]string, len(urls)+1)
	args[0] = "index.js"
	for i, v := range urls {
		args[i+1] = v
	}

	cmd := exec.Command("node", args...)

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return &contentMap, err
	}

	bytes := out.Bytes()
	if err := json.Unmarshal(bytes, &contentMap); err != nil {
		return &contentMap, err
	}

	return &contentMap, nil
}

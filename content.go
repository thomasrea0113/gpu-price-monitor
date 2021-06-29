package monitor

import (
	"bytes"
	"encoding/json"
	"os/exec"
	"strings"
)

func ExecPuppeteer(urls []string) (*map[string]string, error) {
	contentMap := make(map[string]string, len(urls))

	// ensure urls are properly quoted
	for i, v := range urls {
		urls[i] = "'" + v + "'"
	}

	args := strings.Join(urls, " ")
	cmd := exec.Command("node", "index.js", args)

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

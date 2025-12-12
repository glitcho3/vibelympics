package helm

import (
    "bytes"
    "fmt"
    "os/exec"
)

func RenderChart(path string) (string, error) {
    cmd := exec.Command("helm", "template", path)

    var out bytes.Buffer
    var stderr bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &stderr

    err := cmd.Run()
    if err != nil {
        return "", fmt.Errorf("helm error: %s", stderr.String())
    }

    return out.String(), nil
}


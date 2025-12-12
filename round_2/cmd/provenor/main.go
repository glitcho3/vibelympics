package main

import (
    "fmt"
    "os"

    "helm-auditor/internal/provenance"
)

func main() {
    image := os.Getenv("PROV_IMAGE")
    output := os.Getenv("OUTPUT_FOLDER")

    fmt.Printf("[provenor] PROV_IMAGE=%s\n", image)
    fmt.Printf("[provenor] OUTPUT_FOLDER=%s\n", output)

    if image == "" {
        panic("PROV_IMAGE empty")
    }
    if output == "" {
        panic("OUTPUT_FOLDER empty")
    }

    if err := provenance.Run(image, output); err != nil {
        panic(err)
    }

    fmt.Println("[provenor] completed OK")
}


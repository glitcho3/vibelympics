package main

import (
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"
)

type TrivyReport struct {
    Results []struct {
        Target            string `json:"Target"`
        Class             string `json:"Class"`
        Type              string `json:"Type"`
        MisconfSummary    struct {
            Successes int `json:"Successes"`
            Failures  int `json:"Failures"`
        } `json:"MisconfSummary"`
        Misconfigurations []struct {
            ID       string `json:"ID"`
            Type     string `json:"Type"`
            Message  string `json:"Message"`
            Severity string `json:"Severity"`
        } `json:"Misconfigurations"`
    } `json:"Results"`
}

type VulnReport struct {
    Results []struct {
        Vulnerabilities []json.RawMessage `json:"Vulnerabilities"`
    } `json:"Results"`
}

// CycloneDX flexible structure
type Sbom struct {
    Components []struct {
        Name string `json:"name"`
    } `json:"components,omitempty"`

    BOM struct {
        Components []struct {
            Name string `json:"name"`
        } `json:"components,omitempty"`
    } `json:"bom,omitempty"`
}

type AuditSummary struct {
    TotalMisconfigs int `json:"total_misconfigs"`
    TotalFailures   int `json:"total_failures"`
    TotalSuccesses  int `json:"total_successes"`
    Criticals       int `json:"criticals"`
    Highs           int `json:"highs"`
    Components      int `json:"components"`
    Vulns           int `json:"vulns"`
}

func main() {
    promChart := os.Getenv("PROM_CHART")
    reportsPath := os.Getenv("OUTPUT_FOLDER")

    summary := AuditSummary{}
    sbomComponents := 0
    sbomVulns := 0

    // Load Trivy misconfig report
    trivyFile := filepath.Join(promChart + ".report.trivy.json")
    trivyPath := filepath.Join("/reports/", trivyFile)

    trivyData, err := os.ReadFile(trivyPath)
    if err != nil {
        panic(err)
    }

    var trivy TrivyReport
    if err := json.Unmarshal(trivyData, &trivy); err != nil {
        panic(err)
    }

    for _, r := range trivy.Results {
        summary.TotalFailures += r.MisconfSummary.Failures
        summary.TotalSuccesses += r.MisconfSummary.Successes

        for _, m := range r.Misconfigurations {
            summary.TotalMisconfigs++
            switch m.Severity {
            case "CRITICAL":
                summary.Criticals++
            case "HIGH":
                summary.Highs++
            }
        }
    }

    chartFolder := filepath.Join("/reports", promChart)

    // Load SBOM components
    cdxFiles, _ := filepath.Glob(filepath.Join(chartFolder, "*.cdx.json"))
    for _, f := range cdxFiles {
        data, _ := os.ReadFile(f)

        var sbom Sbom
        if err := json.Unmarshal(data, &sbom); err != nil {
            fmt.Println("Skipping invalid SBOM:", f, err)
            continue
        }

        count := 0

        // Prefer .components[]
        if len(sbom.Components) > 0 {
            count = len(sbom.Components)
        }

        // Fallback .bom.components[]
        if count == 0 && len(sbom.BOM.Components) > 0 {
            count = len(sbom.BOM.Components)
        }

        sbomComponents += count
    }

    summary.Components = sbomComponents

    // Load vuln reports
    vulnFiles, _ := filepath.Glob(filepath.Join(chartFolder, "*.vulns.json"))
    for _, f := range vulnFiles {
        data, _ := os.ReadFile(f)

        var vr VulnReport
        if err := json.Unmarshal(data, &vr); err != nil {
            fmt.Println("Skipping invalid vuln report:", f, err)
            continue
        }

        for _, r := range vr.Results {
            sbomVulns += len(r.Vulnerabilities)
        }
    }

    summary.Vulns = sbomVulns

    // Write final JSON
    out, _ := json.MarshalIndent(summary, "", "  ")
    resultPath := filepath.Join(reportsPath, "audit-summary.json")

    if err := os.WriteFile(resultPath, out, 0644); err != nil {
        panic(err)
    }

    fmt.Println("Audit report written to", resultPath)
    fmt.Println(string(out))

    if summary.Criticals > 0 {
        fmt.Println("Critical issues detected. Failing stage.")
        os.Exit(1)
    }

    fmt.Println("No critical issues. Continue.")
    os.Exit(0)
}


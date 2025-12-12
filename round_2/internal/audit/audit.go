package audit

import (
    "strings"

    "gopkg.in/yaml.v3"
    "helm-auditor/internal/types"
)

func AuditYAML(yamlText string, chartPath string) (*types.AuditResult, error) {
    docs := strings.Split(yamlText, "\n---")
    result := &types.AuditResult{
        ChartPath: chartPath,
        Kinds:     map[string]int{},
    }

    for _, d := range docs {
        var obj map[string]any

        // skip empty docs
        trimmed := strings.TrimSpace(d)
        if trimmed == "" {
            continue
        }

        err := yaml.Unmarshal([]byte(trimmed), &obj)
        if err != nil {
            result.Warnings = append(result.Warnings, "Failed to parse a YAML doc")
            continue
        }

        result.AST = append(result.AST, obj)

        // Count resource kinds
        kind, _ := obj["kind"].(string)
        if kind != "" {
            result.Kinds[kind]++
        }
    }

    result.Resources = len(result.AST)

    return result, nil
}


package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type ImageSummary struct {
	Name            string `json:"name"`
	Digest          string `json:"digest"`
	Signed          bool   `json:"signed"`
	Components      int    `json:"components"`
	Vulnerabilities int    `json:"vulnerabilities"`
}

type ExtendedAudit struct {
	Chart struct {
		Name    string `json:"name"`
		URL     string `json:"url"`
		Version string `json:"version"`
	} `json:"chart"`

	ImagesSummary struct {
		Total  int            `json:"total_images"`
		Images []ImageSummary `json:"images"`
	} `json:"images_summary"`

	TotalMisconfigs int `json:"total_misconfigs"`
	TotalFailures   int `json:"total_failures"`
	TotalSuccesses  int `json:"total_successes"`
	Criticals       int `json:"criticals"`
	Highs           int `json:"highs"`
	Components      int `json:"components"`
	Vulns           int `json:"vulns"`
}

func countComponents(sbomData []byte) int {
	compCount := 0
	var raw struct {
		Components []struct {
			Name string `json:"name"`
		} `json:"components,omitempty"`
		BOM struct {
			Components []struct {
				Name string `json:"name"`
			} `json:"components,omitempty"`
		} `json:"bom,omitempty"`
		Results []struct {
			Packages []any `json:"Packages"`
		} `json:"Results"`
	}
	if err := json.Unmarshal(sbomData, &raw); err != nil {
		return 0
	}

	if len(raw.Components) > 0 {
		compCount = len(raw.Components)
	} else if len(raw.BOM.Components) > 0 {
		compCount = len(raw.BOM.Components)
	} else {
		for _, r := range raw.Results {
			compCount += len(r.Packages)
		}
	}
	return compCount
}

func countVulns(vulnData []byte) int {
	vCount := 0
	var raw struct {
		Results []struct {
			Vulnerabilities []any `json:"Vulnerabilities"`
		} `json:"Results"`
	}
	if err := json.Unmarshal(vulnData, &raw); err != nil {
		return 0
	}
	for _, r := range raw.Results {
		vCount += len(r.Vulnerabilities)
	}
	return vCount
}

func main() {
	chartName := os.Getenv("PROM_CHART")
	chartRepo := os.Getenv("PROM_REPO")
	chartVersion := os.Getenv("PROM_VERSION")
	reportsPath := os.Getenv("OUTPUT_FOLDER")
	imagesFile := filepath.Join(reportsPath, "images.txt")

	extended := ExtendedAudit{}
	extended.Chart.Name = chartName
	extended.Chart.URL = fmt.Sprintf("%s%s", chartRepo, chartName)
	extended.Chart.Version = chartVersion

	imagesData, err := os.ReadFile(imagesFile)
	if err != nil {
		fmt.Println("Cannot read images.txt:", err)
		os.Exit(1)
	}

	images := strings.Split(strings.TrimSpace(string(imagesData)), "\n")
	extended.ImagesSummary.Total = len(images)

	totalComponents := 0
	totalVulns := 0

	sbomFiles, _ := filepath.Glob(filepath.Join(reportsPath, "*.cdx.json"))
        vulnFiles, _ := filepath.Glob(filepath.Join(reportsPath, "*.vulns.json"))
        provFiles, _ := filepath.Glob(filepath.Join(reportsPath, "*.prov.json"))
        
        for i, img := range images {
            compCount := 0
            vCount := 0
            signed := false
        
            if i < len(sbomFiles) {
                data, _ := os.ReadFile(sbomFiles[i])
                compCount = countComponents(data)
                totalComponents += compCount
            }
            if i < len(vulnFiles) {
                data, _ := os.ReadFile(vulnFiles[i])
                vCount = countVulns(data)
                totalVulns += vCount
            }
            if i < len(provFiles) {
                signed = true
            }
        
            extended.ImagesSummary.Images = append(extended.ImagesSummary.Images, ImageSummary{
                Name:            img,
                Digest:          "", // opcional: puedes leer digest desde SBOM si quieres
                Signed:          signed,
                Components:      compCount,
                Vulnerabilities: vCount,
            })
        
            fmt.Println("SBOM loaded for", img, "components:", compCount)
            fmt.Println("Vulns loaded:", vCount)
        }

	//for _, img := range images {
	//	if img == "" {
	//		continue
	//	}

	//	h := sha256.Sum256([]byte(img))
	//	hash := fmt.Sprintf("%x", h)

	//	// Variables declaradas aquí
	//	compCount := 0
	//	vCount := 0
	//	signed := false

	//	sbomFiles, _ := filepath.Glob(filepath.Join(reportsPath, hash+".cdx.json"))
	//	vulnFiles, _ := filepath.Glob(filepath.Join(reportsPath, hash+".vulns.json"))
	//	provFiles, _ := filepath.Glob(filepath.Join(reportsPath, hash+".prov.json"))

	//	if len(sbomFiles) > 0 {
	//		data, _ := os.ReadFile(sbomFiles[0])
	//		compCount = countComponents(data)
	//		totalComponents += compCount
	//	}

	//	if len(vulnFiles) > 0 {
	//		data, _ := os.ReadFile(vulnFiles[0])
	//		vCount = countVulns(data)
	//		totalVulns += vCount
	//	}

	//	if len(provFiles) > 0 {
	//		signed = true
	//	}

	//	extended.ImagesSummary.Images = append(extended.ImagesSummary.Images, ImageSummary{
	//		Name:            img,
	//		Digest:          hash,
	//		Signed:          signed,
	//		Components:      compCount,
	//		Vulnerabilities: vCount,
	//	})

	//	fmt.Println("SBOM loaded for", img, "components:", compCount)
	//	fmt.Println("Vulns loaded:", vCount)
	//}

	// Parte auditor Trivy
	trivyFile := filepath.Join(reportsPath, chartName, chartName+".report.trivy.json")
	if data, err := os.ReadFile(trivyFile); err == nil {
		var trivy struct {
			Results []struct {
				MisconfSummary struct {
					Successes int `json:"Successes"`
					Failures  int `json:"Failures"`
				} `json:"MisconfSummary"`
				Misconfigurations []struct {
					Severity string `json:"Severity"`
				} `json:"Misconfigurations"`
			} `json:"Results"`
		}
		if err := json.Unmarshal(data, &trivy); err == nil {
			for _, r := range trivy.Results {
				extended.TotalFailures += r.MisconfSummary.Failures
				extended.TotalSuccesses += r.MisconfSummary.Successes
				for _, m := range r.Misconfigurations {
					extended.TotalMisconfigs++
					switch m.Severity {
					case "CRITICAL":
						extended.Criticals++
					case "HIGH":
						extended.Highs++
					}
				}
			}
		}
	}

	extended.Components = totalComponents
	extended.Vulns = totalVulns

	outFile := filepath.Join(reportsPath, "audit-images.json")
	outData, _ := json.MarshalIndent(extended, "", "  ")
	if err := os.WriteFile(outFile, outData, 0644); err != nil {
		fmt.Println("Error writing audit report:", err)
		os.Exit(1)
	}

	fmt.Println("Extended audit report written to", outFile)
}


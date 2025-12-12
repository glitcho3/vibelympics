package provenance

import (
    "context"
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"

    "github.com/sigstore/cosign/v2/pkg/cosign"
    "github.com/sigstore/cosign/v2/pkg/oci"
    "github.com/google/go-containerregistry/pkg/name"
    "helm-auditor/internal/types"
)

// Run executes verification and writes reports to reportDir
func Run(imageRef, reportDir string) error {
    ctx := context.Background()

    sigs, attes := verifyImageSafe(ctx, imageRef)

    if err := os.MkdirAll(reportDir, 0o755); err != nil {
        return fmt.Errorf("creating report dir: %w", err)
    }

    // Save signatures
    sigFile := filepath.Join(reportDir, "signatures.json")
    sigStrings := []string{}
    if sigs != nil {
        sigStrings = make([]string, len(sigs))
        for i, s := range sigs {
            b, _ := s.Base64Signature()
            sigStrings[i] = string(b)
        }
    }
    sigBytes, _ := json.MarshalIndent(sigStrings, "", "  ")
    if err := os.WriteFile(sigFile, sigBytes, 0o644); err != nil {
        return fmt.Errorf("writing signatures: %w", err)
    }

    // Save attestations
    attResults := []types.AttestationResult{}
    if attes != nil {
        attResults = make([]types.AttestationResult, len(attes))
        payloadsFile := filepath.Join(reportDir, "payloads.txt")
        f, err := os.Create(payloadsFile)
        if err != nil {
            return fmt.Errorf("creating payloads file: %w", err)
        }
        defer f.Close()

        for i, a := range attes {
            payload, err := a.Payload()
            if err != nil {
                fmt.Printf("[provenance] WARNING: reading payload %d failed: %v\n", i, err)
                continue
            }
            f.WriteString(fmt.Sprintf("Attestation %d:\n%s\n\n", i, string(payload)))
            attResults[i] = types.AttestationResult{
                Name:    fmt.Sprintf("Attestation %d", i),
                Payload: string(payload),
            }
        }
    }

    attBytes, _ := json.MarshalIndent(attResults, "", "  ")
    attFile := filepath.Join(reportDir, "attestations.json")
    if err := os.WriteFile(attFile, attBytes, 0o644); err != nil {
        return fmt.Errorf("writing attestations: %w", err)
    }

    fmt.Printf("[provenance] Reports written to: %s\n", reportDir)
    return nil
}

// verifyImageSafe wraps cosign verification and never panics.
// Returns nil slices if verification fails.
func verifyImageSafe(ctx context.Context, imageRef string) ([]oci.Signature, []oci.Signature) {
    ref, err := name.ParseReference(imageRef)
    if err != nil {
        fmt.Printf("[provenance] WARNING: parsing reference failed: %v\n", err)
        return nil, nil
    }

    opts := &cosign.CheckOpts{} // default, safe
    sigs, bundle, err := cosign.VerifyImageSignatures(ctx, ref, opts)
    if err != nil {
        fmt.Printf("[provenance] WARNING: signatures not verified for %s: %v\n", imageRef, err)
        sigs = nil
    } else {
        fmt.Printf("[provenance] Signatures verified: %v, bundle: %v\n", sigs, bundle)
    }

    attes, attBundle, err := cosign.VerifyImageAttestations(ctx, ref, opts)
    if err != nil {
        fmt.Printf("[provenance] WARNING: attestations not verified for %s: %v\n", imageRef, err)
        attes = nil
    } else {
        fmt.Printf("[provenance] Attestations verified: %v, bundle: %v\n", attes, attBundle)
    }

    return sigs, attes
}


package types

//import "github.com/sigstore/cosign/v2/pkg/oci"

// AttestationResult wraps la attestation y su payload
type AttestationResult struct {
    Name    string `json:"name,omitempty"`    // opcional, puedes poner algo del bundle
    Payload string `json:"payload"`          // contenido firmado
    // puedes agregar más campos como firmadoPor, fecha, etc.
}

// AuditResult es el resultado de la auditoría
type AuditResult struct {
    ChartPath    string             `json:"chart_path"`
    Resources    int                `json:"resources"`
    Kinds        map[string]int     `json:"kinds"`
    Warnings     []string           `json:"warnings"`
    AST          []map[string]any   `json:"ast"`
    Signatures   []string           `json:"signatures,omitempty"`    // firmas de la imagen
    Attestations []AttestationResult `json:"attestations,omitempty"` // payloads de attestations
}


package evidence

// EvidenceStore handles local file storage for scan evidence
type EvidenceStore struct {
    basePath string
}

// NewEvidenceStore creates a new evidence store
func NewEvidenceStore(basePath string) *EvidenceStore {
    return &EvidenceStore{basePath: basePath}
}

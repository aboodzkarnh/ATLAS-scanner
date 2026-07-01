package integrations

import (
    "fmt"
    "os"
)

// InteractshConfig holds configuration for self-hosted Interactsh
type InteractshConfig struct {
    URL   string
    Token string
}

// InteractshClient defines the interface for OOB interaction
type InteractshClient interface {
    GenerateURL() (string, error)
    Poll() ([]string, error)
    Close() error
}

// stubClient is a placeholder that does nothing
type stubClient struct{}

func (s *stubClient) GenerateURL() (string, error) {
    return "", fmt.Errorf("interactsh client not configured")
}

func (s *stubClient) Poll() ([]string, error) {
    return nil, fmt.Errorf("interactsh client not configured")
}

func (s *stubClient) Close() error {
    return nil
}

// NewInteractshClient creates a new Interactsh client
// Currently returns a stub - real implementation will be added in Phase 2
func NewInteractshClient(cfg *InteractshConfig) (InteractshClient, error) {
    // ENV placeholders for future use
    url := os.Getenv("INTERACTSH_URL")
    token := os.Getenv("INTERACTSH_TOKEN")

    if url != "" && token != "" {
        return nil, fmt.Errorf("Interactsh integration not yet implemented. URL=%s (token hidden)", url)
    }

    return &stubClient{}, nil
}

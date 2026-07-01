package reports

// ReportGenerator creates scan reports in multiple formats
type ReportGenerator struct {
    format string
}

// NewReportGenerator creates a new report generator
func NewReportGenerator(format string) *ReportGenerator {
    return &ReportGenerator{format: format}
}

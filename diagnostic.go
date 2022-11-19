package main

type Severity int

const (
	Info Severity = iota
	Warning
	Error
)

func (s Severity) String() string {
	switch s {
	case Info:
		return "Info"
	case Warning:
		return "Warning"
	case Error:
		return "Error"
	default:
		return string(s)
	}
}

type Diagnostic struct {
	// The message to display
	Message string
	// The position in the source code where the error occurred
	Position int
	// The length of the error
	Length int
	// The severity of the error
	Severity Severity
}

func (d Diagnostic) IsError() bool {
	return d.Severity == Error
}

func (d Diagnostic) IsWarning() bool {
	return d.Severity == Warning
}

func (d Diagnostic) IsInfo() bool {
	return d.Severity == Info
}

func (d Diagnostic) String() string {
	return d.Severity.String() + ": " + d.Message
}

type DiangosticContainer struct {
	diagnostics []Diagnostic
}

func (dc *DiangosticContainer) Add(diagnostic Diagnostic) {
	dc.diagnostics = append(dc.diagnostics, diagnostic)
}

func (dc *DiangosticContainer) AddRange(diagnostics []Diagnostic) {
	dc.diagnostics = append(dc.diagnostics, diagnostics...)
}

func (dc *DiangosticContainer) Diagnostics() []Diagnostic {
	return dc.diagnostics
}

func (dc *DiangosticContainer) Clear() {
	dc.diagnostics = []Diagnostic{}
}

func (dc *DiangosticContainer) HasErrors() bool {
	for _, diagnostic := range dc.diagnostics {
		if diagnostic.Severity == Error {
			return true
		}
	}
	return false
}

func (dc *DiangosticContainer) HasWarnings() bool {
	for _, diagnostic := range dc.diagnostics {
		if diagnostic.Severity == Warning {
			return true
		}
	}
	return false
}

func (dc *DiangosticContainer) HasInfos() bool {
	for _, diagnostic := range dc.diagnostics {
		if diagnostic.Severity == Info {
			return true
		}
	}
	return false
}

func NewDiagnosticContainer() *DiangosticContainer {
	return &DiangosticContainer{
		diagnostics: []Diagnostic{},
	}
}

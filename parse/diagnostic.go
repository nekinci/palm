package parse

type DiagnosticKind int

const (
	DiagnosticError DiagnosticKind = iota
	DiagnosticWarning
)

type DiagnosticMessage struct {
	Kind     DiagnosticKind
	Position Location
	Message  string
}

func (d *DiagnosticMessage) String() string {
	return d.Message
}

func (d *DiagnosticMessage) IsError() bool {
	return d.Kind == DiagnosticError
}

func (d *DiagnosticMessage) IsWarning() bool {
	return d.Kind == DiagnosticWarning
}

func newDiagnosticError(location Location, message string) *DiagnosticMessage {
	return &DiagnosticMessage{
		Kind:     DiagnosticError,
		Position: location,
		Message:  message,
	}
}

func newDiagnosticWarning(location Location, message string) *DiagnosticMessage {
	return &DiagnosticMessage{
		Kind:     DiagnosticWarning,
		Position: location,
		Message:  message,
	}
}

type Diagnostic struct {
	diagnostics []DiagnosticMessage
}

func (d *Diagnostic) AddError(location Location, message string) {
	d.diagnostics = append(d.diagnostics, *newDiagnosticError(location, message))
}

func (d *Diagnostic) AddWarning(location Location, message string) {
	d.diagnostics = append(d.diagnostics, *newDiagnosticWarning(location, message))
}

func (d *Diagnostic) HasErrors() bool {
	for _, diagnostic := range d.diagnostics {
		if diagnostic.IsError() {
			return true
		}
	}
	return false
}

func (d *Diagnostic) HasWarnings() bool {
	for _, diagnostic := range d.diagnostics {
		if diagnostic.IsWarning() {
			return true
		}
	}
	return false
}

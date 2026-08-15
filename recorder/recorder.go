package recorder

import (
	"fmt"
)

// Recorder buffers intermediate actions and generates formatted scripts.
type Recorder struct {
	actions []ActionIR
}

// NewRecorder instantiates a Recorder buffer.
func NewRecorder() *Recorder {
	return &Recorder{
		actions: make([]ActionIR, 0),
	}
}

// AddAction adds a raw action to the internal buffer, performing merging optimizations if possible.
func (r *Recorder) AddAction(action ActionIR) {
	if len(r.actions) == 0 {
		r.actions = append(r.actions, action)
		return
	}

	// Optimization: Merge consecutive fills on the same element
	lastIdx := len(r.actions) - 1
	lastAction := r.actions[lastIdx]

	if action.Type == ActionFill && lastAction.Type == ActionFill &&
		action.SelectorStrategy == lastAction.SelectorStrategy &&
		action.Selector == lastAction.Selector {
		// Update value of last action rather than adding new ones
		r.actions[lastIdx].Value = action.Value
		return
	}

	r.actions = append(r.actions, action)
}

// Generate compiles the recorded buffer using the target language generator.
func (r *Recorder) Generate(lang string) (string, error) {
	var generator CodeGenerator

	switch lang {
	case "typescript", "ts":
		generator = &TypeScriptGenerator{}
	case "python", "py":
		generator = &PythonGenerator{}
	case "go":
		generator = &GoGenerator{}
	case "java":
		generator = &JavaGenerator{}
	case "csharp", "cs":
		generator = &CSharpGenerator{}
	case "kotlin", "kt":
		generator = &KotlinGenerator{}
	default:
		return "", fmt.Errorf("unsupported language generator type %q", lang)
	}

	return generator.Generate(r.actions)
}

// Clear flushes the event buffer.
func (r *Recorder) Clear() {
	r.actions = make([]ActionIR, 0)
}

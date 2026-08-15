package recorder

// CodeGenerator defines the blueprint that compiles intermediate representation vectors into SDK codes.
type CodeGenerator interface {
	// Generate returns the compiled test code script.
	Generate(actions []ActionIR) (string, error)
}

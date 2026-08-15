package recorder

// ActionType defines the type of event in our Intermediate Representation.
type ActionType string

const (
	ActionClick    ActionType = "click"
	ActionFill     ActionType = "fill"
	ActionSwipe    ActionType = "swipe"
	ActionOpenApp  ActionType = "open_app"
	ActionCloseApp ActionType = "close_app"
)

// ActionIR stores the structured mobile gesture actions for code generation.
type ActionIR struct {
	Type             ActionType
	SelectorStrategy string
	Selector         string
	Value            string
	StartX           int
	StartY           int
	EndX             int
	EndY             int
	DurationMS       int
}

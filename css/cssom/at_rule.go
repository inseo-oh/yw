package cssom

// AtRule represents CSS at-rule (e.g. @media { ... })
type AtRule struct {
	Name    string
	Prelude any
	Value   any
}

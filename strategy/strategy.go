package strategy

// ExecutionStrategy defines what queries are run how
type ExecutionStrategy interface {
	RunCommand()
}

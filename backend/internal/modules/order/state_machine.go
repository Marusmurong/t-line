package order

import "fmt"

// validTransitions defines allowed state transitions for orders.
var validTransitions = map[string][]string{
	StatusPending:   {StatusPaid, StatusCancelled},
	StatusPaid:      {StatusUsed, StatusCancelled, StatusRefunding},
	StatusUsed:      {StatusCompleted},
	StatusRefunding: {StatusRefunded},
	// Terminal states: cancelled, completed, refunded have no transitions
}

// ValidateTransition checks if a state transition is allowed.
func ValidateTransition(from, to string) error {
	allowed, exists := validTransitions[from]
	if !exists {
		return fmt.Errorf("状态 %s 为终态，不可变更", from)
	}

	for _, s := range allowed {
		if s == to {
			return nil
		}
	}

	return fmt.Errorf("不允许从 %s 变更到 %s", from, to)
}

// IsTerminalState returns true if the status is a terminal state.
func IsTerminalState(status string) bool {
	return status == StatusCancelled || status == StatusCompleted || status == StatusRefunded
}

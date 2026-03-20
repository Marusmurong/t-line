package order

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateTransition(t *testing.T) {
	tests := []struct {
		name    string
		from    string
		to      string
		wantErr bool
	}{
		// Legal transitions
		{name: "pending -> paid", from: StatusPending, to: StatusPaid, wantErr: false},
		{name: "pending -> cancelled", from: StatusPending, to: StatusCancelled, wantErr: false},
		{name: "paid -> used", from: StatusPaid, to: StatusUsed, wantErr: false},
		{name: "paid -> refunding", from: StatusPaid, to: StatusRefunding, wantErr: false},
		{name: "paid -> cancelled", from: StatusPaid, to: StatusCancelled, wantErr: false},
		{name: "used -> completed", from: StatusUsed, to: StatusCompleted, wantErr: false},
		{name: "refunding -> refunded", from: StatusRefunding, to: StatusRefunded, wantErr: false},

		// Illegal transitions: terminal states cannot transition
		{name: "cancelled -> paid (terminal)", from: StatusCancelled, to: StatusPaid, wantErr: true},
		{name: "completed -> cancelled (terminal)", from: StatusCompleted, to: StatusCancelled, wantErr: true},
		{name: "refunded -> paid (terminal)", from: StatusRefunded, to: StatusPaid, wantErr: true},

		// Illegal transitions: skip / invalid
		{name: "pending -> completed (skip)", from: StatusPending, to: StatusCompleted, wantErr: true},
		{name: "pending -> used (skip)", from: StatusPending, to: StatusUsed, wantErr: true},
		{name: "pending -> refunding (skip)", from: StatusPending, to: StatusRefunding, wantErr: true},
		{name: "paid -> completed (skip)", from: StatusPaid, to: StatusCompleted, wantErr: true},
		{name: "used -> paid (reverse)", from: StatusUsed, to: StatusPaid, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTransition(tt.from, tt.to)
			if tt.wantErr {
				assert.Error(t, err, "transition %s -> %s should be rejected", tt.from, tt.to)
			} else {
				assert.NoError(t, err, "transition %s -> %s should be allowed", tt.from, tt.to)
			}
		})
	}
}

func TestIsTerminalState(t *testing.T) {
	tests := []struct {
		status   string
		terminal bool
	}{
		{StatusPending, false},
		{StatusPaid, false},
		{StatusUsed, false},
		{StatusRefunding, false},
		{StatusCancelled, true},
		{StatusCompleted, true},
		{StatusRefunded, true},
	}

	for _, tt := range tests {
		t.Run(tt.status, func(t *testing.T) {
			assert.Equal(t, tt.terminal, IsTerminalState(tt.status))
		})
	}
}

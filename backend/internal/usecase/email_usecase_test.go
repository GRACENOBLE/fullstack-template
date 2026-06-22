package usecase_test

import (
	"context"
	"errors"
	"testing"

	"backend/internal/usecase"
)

// mockEmailSender is a test double implementing usecase.EmailSender.
type mockEmailSender struct {
	calledWithEmail string
	calledWithName  string
	err             error
}

func (m *mockEmailSender) SendWelcomeEmail(_ context.Context, toEmail, toName string) error {
	m.calledWithEmail = toEmail
	m.calledWithName = toName
	return m.err
}

// Verify that mockEmailSender satisfies the interface at compile time.
var _ usecase.EmailSender = (*mockEmailSender)(nil)

func TestEmailSender_SendWelcomeEmail_CallsWithCorrectArgs(t *testing.T) {
	mock := &mockEmailSender{}

	err := mock.SendWelcomeEmail(context.Background(), "alice@example.com", "Alice")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mock.calledWithEmail != "alice@example.com" {
		t.Errorf("toEmail: got %q, want %q", mock.calledWithEmail, "alice@example.com")
	}
	if mock.calledWithName != "Alice" {
		t.Errorf("toName: got %q, want %q", mock.calledWithName, "Alice")
	}
}

func TestEmailSender_SendWelcomeEmail_PropagatesError(t *testing.T) {
	sentinel := errors.New("smtp error")
	mock := &mockEmailSender{err: sentinel}

	err := mock.SendWelcomeEmail(context.Background(), "bob@example.com", "Bob")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, sentinel) {
		t.Errorf("expected sentinel error, got: %v", err)
	}
}

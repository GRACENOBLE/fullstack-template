package postgres

import (
	"context"
	"testing"
)

// setupFCMTokensTable creates the fcm_tokens table for tests and drops it afterwards.
// Integration tests do NOT rely on migrations — each test sets up its own schema.
func setupFCMTokensTable(t *testing.T) {
	t.Helper()
	_, err := testDB.Exec(`
		CREATE TABLE IF NOT EXISTS fcm_tokens (
			id          UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id     TEXT         NOT NULL,
			token       TEXT         NOT NULL UNIQUE,
			platform    TEXT         NOT NULL,
			created_at  TIMESTAMPTZ  NOT NULL DEFAULT now()
		)`)
	if err != nil {
		t.Fatalf("setupFCMTokensTable: %v", err)
	}
	t.Cleanup(func() { testDB.Exec("DROP TABLE IF EXISTS fcm_tokens") }) //nolint:errcheck
}

func TestFCMTokenRepository_SaveAndGet(t *testing.T) {
	setupFCMTokensTable(t)
	ctx := context.Background()
	repo := NewFCMTokenRepository(testDB)

	if err := repo.SaveToken(ctx, "user1", "tok-abc", "android"); err != nil {
		t.Fatalf("SaveToken: %v", err)
	}

	tokens, err := repo.GetTokensByUserID(ctx, "user1")
	if err != nil {
		t.Fatalf("GetTokensByUserID: %v", err)
	}
	if len(tokens) != 1 {
		t.Fatalf("expected 1 token, got %d", len(tokens))
	}
	if tokens[0].Token != "tok-abc" || tokens[0].Platform != "android" || tokens[0].UserID != "user1" {
		t.Errorf("unexpected token: %+v", tokens[0])
	}
}

func TestFCMTokenRepository_Upsert(t *testing.T) {
	setupFCMTokensTable(t)
	ctx := context.Background()
	repo := NewFCMTokenRepository(testDB)

	if err := repo.SaveToken(ctx, "userA", "tok-upsert", "web"); err != nil {
		t.Fatalf("first SaveToken: %v", err)
	}
	// Same physical token, different user — upsert should reassign it.
	if err := repo.SaveToken(ctx, "userB", "tok-upsert", "web"); err != nil {
		t.Fatalf("second SaveToken (upsert): %v", err)
	}

	tokensA, err := repo.GetTokensByUserID(ctx, "userA")
	if err != nil {
		t.Fatalf("GetTokensByUserID(userA): %v", err)
	}
	tokensB, err := repo.GetTokensByUserID(ctx, "userB")
	if err != nil {
		t.Fatalf("GetTokensByUserID(userB): %v", err)
	}
	if len(tokensA) != 0 {
		t.Errorf("expected userA to have 0 tokens after upsert, got %d", len(tokensA))
	}
	if len(tokensB) != 1 {
		t.Errorf("expected userB to have 1 token, got %d", len(tokensB))
	}
}

func TestFCMTokenRepository_Delete(t *testing.T) {
	setupFCMTokensTable(t)
	ctx := context.Background()
	repo := NewFCMTokenRepository(testDB)

	if err := repo.SaveToken(ctx, "user2", "tok-del", "ios"); err != nil {
		t.Fatalf("SaveToken: %v", err)
	}
	if err := repo.DeleteToken(ctx, "user2", "tok-del"); err != nil {
		t.Fatalf("DeleteToken: %v", err)
	}

	tokens, err := repo.GetTokensByUserID(ctx, "user2")
	if err != nil {
		t.Fatalf("GetTokensByUserID: %v", err)
	}
	if len(tokens) != 0 {
		t.Errorf("expected 0 tokens after delete, got %d", len(tokens))
	}
}

func TestFCMTokenRepository_DeleteOtherUser(t *testing.T) {
	setupFCMTokensTable(t)
	ctx := context.Background()
	repo := NewFCMTokenRepository(testDB)

	if err := repo.SaveToken(ctx, "user3", "tok-other", "web"); err != nil {
		t.Fatalf("SaveToken: %v", err)
	}

	// Attempt to delete using a different userID — must be a no-op.
	if err := repo.DeleteToken(ctx, "attacker", "tok-other"); err != nil {
		t.Fatalf("DeleteToken with wrong userID returned unexpected error: %v", err)
	}

	tokens, err := repo.GetTokensByUserID(ctx, "user3")
	if err != nil {
		t.Fatalf("GetTokensByUserID: %v", err)
	}
	if len(tokens) != 1 {
		t.Errorf("expected token to remain after cross-user delete attempt, got %d tokens", len(tokens))
	}
}

func TestFCMTokenRepository_GetEmpty(t *testing.T) {
	setupFCMTokensTable(t)
	ctx := context.Background()
	repo := NewFCMTokenRepository(testDB)

	tokens, err := repo.GetTokensByUserID(ctx, "nonexistent")
	if err != nil {
		t.Fatalf("GetTokensByUserID: %v", err)
	}
	if len(tokens) != 0 {
		t.Errorf("expected 0 tokens, got %d", len(tokens))
	}
}

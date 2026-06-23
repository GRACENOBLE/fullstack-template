package postgres

import (
	"context"
	"testing"

	"backend/internal/domain"
)

func setupUsersTable(t *testing.T) {
	t.Helper()
	_, err := testDB.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id           BIGSERIAL    PRIMARY KEY,
			firebase_uid TEXT         UNIQUE,
			name         TEXT         NOT NULL DEFAULT '',
			email        TEXT,
			photo_url    TEXT,
			created_at   TIMESTAMPTZ  NOT NULL DEFAULT now(),
			updated_at   TIMESTAMPTZ  NOT NULL DEFAULT now()
		)`)
	if err != nil {
		t.Fatalf("setupUsersTable: %v", err)
	}
	t.Cleanup(func() {
		_, _ = testDB.Exec(`DROP TABLE IF EXISTS users`)
	})
}

func TestUserRepository_Upsert_Create(t *testing.T) {
	setupUsersTable(t)
	repo := NewUserRepository(testDB)

	u := &domain.User{
		FirebaseUID: "uid-001",
		Name:        "Alice",
		Email:       "alice@example.com",
		PhotoURL:    "https://example.com/alice.png",
	}
	got, err := repo.Upsert(context.Background(), u)
	if err != nil {
		t.Fatalf("Upsert() unexpected error: %v", err)
	}
	if got.ID == 0 {
		t.Error("expected non-zero ID after upsert")
	}
	if got.FirebaseUID != u.FirebaseUID {
		t.Errorf("FirebaseUID: got %q, want %q", got.FirebaseUID, u.FirebaseUID)
	}
	if got.Name != u.Name {
		t.Errorf("Name: got %q, want %q", got.Name, u.Name)
	}
	if got.Email != u.Email {
		t.Errorf("Email: got %q, want %q", got.Email, u.Email)
	}
	if got.PhotoURL != u.PhotoURL {
		t.Errorf("PhotoURL: got %q, want %q", got.PhotoURL, u.PhotoURL)
	}
}

func TestUserRepository_Upsert_Update(t *testing.T) {
	setupUsersTable(t)
	repo := NewUserRepository(testDB)

	first := &domain.User{
		FirebaseUID: "uid-002",
		Name:        "Bob",
		Email:       "bob@example.com",
		PhotoURL:    "",
	}
	created, err := repo.Upsert(context.Background(), first)
	if err != nil {
		t.Fatalf("first Upsert() error: %v", err)
	}

	updated := &domain.User{
		FirebaseUID: "uid-002",
		Name:        "Bob Updated",
		Email:       "bob+new@example.com",
		PhotoURL:    "https://example.com/bob.png",
	}
	got, err := repo.Upsert(context.Background(), updated)
	if err != nil {
		t.Fatalf("second Upsert() error: %v", err)
	}
	if got.ID != created.ID {
		t.Errorf("expected same ID on update: got %d, want %d", got.ID, created.ID)
	}
	if got.Name != "Bob Updated" {
		t.Errorf("Name: got %q, want %q", got.Name, "Bob Updated")
	}
	if got.Email != "bob+new@example.com" {
		t.Errorf("Email: got %q, want %q", got.Email, "bob+new@example.com")
	}
}

func TestUserRepository_DeleteByFirebaseUID(t *testing.T) {
	setupUsersTable(t)
	repo := NewUserRepository(testDB)

	u := &domain.User{
		FirebaseUID: "uid-003",
		Name:        "Carol",
		Email:       "carol@example.com",
	}
	if _, err := repo.Upsert(context.Background(), u); err != nil {
		t.Fatalf("Upsert() error: %v", err)
	}

	if err := repo.DeleteByFirebaseUID(context.Background(), "uid-003"); err != nil {
		t.Fatalf("DeleteByFirebaseUID() error: %v", err)
	}

	// Confirm record is gone.
	var count int
	if err := testDB.QueryRow(`SELECT COUNT(*) FROM users WHERE firebase_uid = $1`, "uid-003").Scan(&count); err != nil {
		t.Fatalf("count query error: %v", err)
	}
	if count != 0 {
		t.Errorf("expected 0 rows after delete, got %d", count)
	}
}

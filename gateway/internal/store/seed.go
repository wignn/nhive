package store

import (
	"crypto/rand"
	"encoding/hex"
	"os"
)

func (s *Store) seed() {
	getPasswordFromEnvOrGenerate := func(key string) string {
		if p := os.Getenv(key); p != "" {
			return p
		}
		b := make([]byte, 8)
		rand.Read(b)
		return "dev-" + hex.EncodeToString(b)
	}
	s.genres = []string{
		"Fantasy", "Action", "Romance", "Adventure", "Sci-Fi", "Mystery",
		"Horror", "Comedy", "Drama", "Slice of Life", "Martial Arts",
		"Isekai", "Wuxia", "Xianxia",
	}

	adminPass := getPasswordFromEnvOrGenerate("ADMIN_PASS")
	admin := &User{
		ID: "admin-001", Username: "admin", Email: "admin@novelhive.com",
		PasswordHash: hashPassword(adminPass), Role: "admin",
	}
	s.users[admin.ID] = admin
	s.userByEmail[admin.Email] = admin.ID

	readerPass := getPasswordFromEnvOrGenerate("READER_PASS")
	reader := &User{
		ID: "reader-001", Username: "bookworm42", Email: "reader@novelhive.com",
		PasswordHash: hashPassword(readerPass), Role: "reader",
	}
	s.users[reader.ID] = reader
	s.userByEmail[reader.Email] = reader.ID
}

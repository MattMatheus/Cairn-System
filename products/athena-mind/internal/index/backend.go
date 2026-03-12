package index

import (
	"os"
	"strings"
)

func selectedStorageBackend() string {
	for _, key := range []string{"ATHENA_INDEX_BACKEND", "ATHENA_STORAGE_BACKEND"} {
		if value := strings.ToLower(strings.TrimSpace(os.Getenv(key))); value != "" {
			switch value {
			case "mongodb":
				return "mongodb"
			case "sqlite":
				return "sqlite"
			}
		}
	}
	return "sqlite"
}

func selectedIndexStore() indexStore {
	switch selectedStorageBackend() {
	case "mongodb":
		return mongoIndexStore{}
	default:
		return sqliteIndexStore{}
	}
}

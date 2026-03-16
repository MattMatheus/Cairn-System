package snapshot

import (
	"testing"

	"memorycli/internal/types"
)

func TestValidateSnapshotManifestRequiresFields(t *testing.T) {
	err := ValidateSnapshotManifest(types.SnapshotManifest{})
	if err == nil {
		t.Fatal("expected validation error")
	}
}

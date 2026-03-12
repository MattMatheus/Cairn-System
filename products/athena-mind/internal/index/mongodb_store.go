package index

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"athenamind/internal/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoIndexStore struct{}

type mongoIndexEntryDocument struct {
	RootPath     string `bson:"root_path"`
	ID           string `bson:"id"`
	Type         string `bson:"type"`
	Domain       string `bson:"domain"`
	Path         string `bson:"path"`
	MetadataPath string `bson:"metadata_path"`
	Status       string `bson:"status"`
	UpdatedAt    string `bson:"updated_at"`
	Title        string `bson:"title"`
}

type mongoEmbeddingDocument struct {
	RootPath    string    `bson:"root_path"`
	EntryID     string    `bson:"entry_id"`
	Vector      []float64 `bson:"vector"`
	ModelID     string    `bson:"model_id,omitempty"`
	Provider    string    `bson:"provider,omitempty"`
	Dim         int       `bson:"dim,omitempty"`
	ContentHash string    `bson:"content_hash,omitempty"`
	CommitSHA   string    `bson:"commit_sha,omitempty"`
	SessionID   string    `bson:"session_id,omitempty"`
	GeneratedAt string    `bson:"generated_at,omitempty"`
	UpdatedAt   string    `bson:"updated_at,omitempty"`
}

func (mongoIndexStore) Load(root string) (types.IndexFile, error) {
	if err := os.MkdirAll(root, 0o755); err != nil {
		return types.IndexFile{}, err
	}

	ctx, cancel := mongoContext()
	defer cancel()
	db, cleanup, err := openMongoDatabase(ctx)
	if err != nil {
		return types.IndexFile{}, err
	}
	defer cleanup()

	coll := db.Collection(mongoEntriesCollection())
	rootKey := mongoRootMarker(root)
	count, err := coll.CountDocuments(ctx, bson.D{{Key: "root_path", Value: rootKey}})
	if err != nil {
		return types.IndexFile{}, fmt.Errorf("mongodb entries count failed: %w", err)
	}
	if count == 0 {
		legacy, legacyErr := loadIndexFromYAML(root)
		if legacyErr != nil {
			return types.IndexFile{}, legacyErr
		}
		if len(legacy.Entries) > 0 {
			if err := saveIndexToMongo(ctx, coll, root, legacy); err != nil {
				return types.IndexFile{}, err
			}
		}
		return legacy, nil
	}

	cursor, err := coll.Find(
		ctx,
		bson.D{{Key: "root_path", Value: rootKey}},
		options.Find().SetSort(bson.D{{Key: "id", Value: 1}}),
	)
	if err != nil {
		return types.IndexFile{}, fmt.Errorf("mongodb entries find failed: %w", err)
	}
	defer cursor.Close(ctx)

	idx := types.IndexFile{
		SchemaVersion: types.DefaultSchema,
		UpdatedAt:     time.Now().UTC().Format(time.RFC3339),
		Entries:       []types.IndexEntry{},
	}
	for cursor.Next(ctx) {
		var doc mongoIndexEntryDocument
		if err := cursor.Decode(&doc); err != nil {
			return types.IndexFile{}, fmt.Errorf("mongodb entry decode failed: %w", err)
		}
		entry := types.IndexEntry{
			ID:           doc.ID,
			Type:         doc.Type,
			Domain:       doc.Domain,
			Path:         doc.Path,
			MetadataPath: doc.MetadataPath,
			Status:       doc.Status,
			UpdatedAt:    doc.UpdatedAt,
			Title:        doc.Title,
		}
		idx.Entries = append(idx.Entries, entry)
		if entry.UpdatedAt > idx.UpdatedAt {
			idx.UpdatedAt = entry.UpdatedAt
		}
	}
	if err := cursor.Err(); err != nil {
		return types.IndexFile{}, fmt.Errorf("mongodb cursor error: %w", err)
	}
	if err := ValidateSchemaVersion(idx.SchemaVersion); err != nil {
		return types.IndexFile{}, err
	}
	if err := ValidateIndex(idx, root); err != nil {
		return types.IndexFile{}, err
	}
	return idx, nil
}

func (mongoIndexStore) Save(root string, idx types.IndexFile) error {
	if err := ValidateSchemaVersion(idx.SchemaVersion); err != nil {
		return err
	}
	ctx, cancel := mongoContext()
	defer cancel()
	db, cleanup, err := openMongoDatabase(ctx)
	if err != nil {
		return err
	}
	defer cleanup()
	return saveIndexToMongo(ctx, db.Collection(mongoEntriesCollection()), root, idx)
}

func saveIndexToMongo(ctx context.Context, coll *mongo.Collection, root string, idx types.IndexFile) error {
	rootKey := mongoRootMarker(root)
	if _, err := coll.DeleteMany(ctx, bson.D{{Key: "root_path", Value: rootKey}}); err != nil {
		return fmt.Errorf("mongodb entries clear failed: %w", err)
	}
	if len(idx.Entries) == 0 {
		return nil
	}
	docs := make([]any, 0, len(idx.Entries))
	for _, e := range idx.Entries {
		docs = append(docs, mongoIndexEntryDocument{
			RootPath:     rootKey,
			ID:           e.ID,
			Type:         e.Type,
			Domain:       e.Domain,
			Path:         e.Path,
			MetadataPath: e.MetadataPath,
			Status:       e.Status,
			UpdatedAt:    e.UpdatedAt,
			Title:        e.Title,
		})
	}
	if _, err := coll.InsertMany(ctx, docs); err != nil {
		return fmt.Errorf("mongodb entries insert failed: %w", err)
	}
	return nil
}

func upsertEmbeddingRecordMongo(root string, record types.EmbeddingRecord) error {
	entryID := strings.TrimSpace(record.EntryID)
	if entryID == "" {
		return errors.New("entry id is required for embedding upsert")
	}
	if len(record.Vector) == 0 {
		return errors.New("embedding vector cannot be empty")
	}
	ctx, cancel := mongoContext()
	defer cancel()
	db, cleanup, err := openMongoDatabase(ctx)
	if err != nil {
		return err
	}
	defer cleanup()

	now := time.Now().UTC().Format(time.RFC3339)
	generatedAt := strings.TrimSpace(record.GeneratedAt)
	if generatedAt == "" {
		generatedAt = now
	}
	doc := mongoEmbeddingDocument{
		RootPath:    mongoRootMarker(root),
		EntryID:     entryID,
		Vector:      record.Vector,
		ModelID:     strings.TrimSpace(record.ModelID),
		Provider:    strings.TrimSpace(record.Provider),
		Dim:         record.Dim,
		ContentHash: strings.TrimSpace(record.ContentHash),
		CommitSHA:   strings.TrimSpace(record.CommitSHA),
		SessionID:   strings.TrimSpace(record.SessionID),
		GeneratedAt: generatedAt,
		UpdatedAt:   now,
	}
	if doc.Dim <= 0 {
		doc.Dim = len(doc.Vector)
	}
	_, err = db.Collection(mongoEmbeddingsCollection()).ReplaceOne(
		ctx,
		bson.M{"root_path": mongoRootMarker(root), "entry_id": entryID},
		doc,
		options.Replace().SetUpsert(true),
	)
	if err != nil {
		return fmt.Errorf("mongodb embedding upsert failed: %w", err)
	}
	return nil
}

func getEmbeddingRecordsMongo(root string, ids []string) (map[string]types.EmbeddingRecord, error) {
	ctx, cancel := mongoContext()
	defer cancel()
	db, cleanup, err := openMongoDatabase(ctx)
	if err != nil {
		return nil, err
	}
	defer cleanup()

	filter := bson.D{{Key: "root_path", Value: mongoRootMarker(root)}}
	if len(ids) > 0 {
		valid := make([]string, 0, len(ids))
		for _, id := range ids {
			if trimmed := strings.TrimSpace(id); trimmed != "" {
				valid = append(valid, trimmed)
			}
		}
		if len(valid) == 0 {
			return map[string]types.EmbeddingRecord{}, nil
		}
		filter = bson.D{
			{Key: "root_path", Value: mongoRootMarker(root)},
			{Key: "entry_id", Value: bson.M{"$in": valid}},
		}
	}

	cursor, err := db.Collection(mongoEmbeddingsCollection()).Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("mongodb embeddings find failed: %w", err)
	}
	defer cursor.Close(ctx)

	out := map[string]types.EmbeddingRecord{}
	for cursor.Next(ctx) {
		var doc mongoEmbeddingDocument
		if err := cursor.Decode(&doc); err != nil {
			return nil, fmt.Errorf("mongodb embedding decode failed: %w", err)
		}
		if len(doc.Vector) == 0 {
			continue
		}
		out[doc.EntryID] = types.EmbeddingRecord{
			EntryID:     doc.EntryID,
			Vector:      doc.Vector,
			ModelID:     doc.ModelID,
			Provider:    doc.Provider,
			Dim:         doc.Dim,
			ContentHash: doc.ContentHash,
			CommitSHA:   doc.CommitSHA,
			SessionID:   doc.SessionID,
			GeneratedAt: doc.GeneratedAt,
			LastUpdated: doc.UpdatedAt,
		}
	}
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("mongodb embeddings cursor error: %w", err)
	}
	return out, nil
}

func mongoContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 5*time.Second)
}

func openMongoDatabase(ctx context.Context) (*mongo.Database, func(), error) {
	uri := strings.TrimSpace(os.Getenv("ATHENA_MONGODB_URI"))
	if uri == "" {
		uri = "mongodb://127.0.0.1:27017"
	}
	database := strings.TrimSpace(os.Getenv("ATHENA_MONGODB_DATABASE"))
	if database == "" {
		database = "athenamind"
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, func() {}, fmt.Errorf("mongodb connect failed: %w", err)
	}
	cleanup := func() {
		disconnectCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_ = client.Disconnect(disconnectCtx)
	}
	if err := client.Ping(ctx, nil); err != nil {
		cleanup()
		return nil, func() {}, fmt.Errorf("mongodb ping failed: %w", err)
	}
	return client.Database(database), cleanup, nil
}

func mongoEntriesCollection() string {
	if name := strings.TrimSpace(os.Getenv("ATHENA_MONGODB_ENTRIES_COLLECTION")); name != "" {
		return name
	}
	return "memory_entries"
}

func mongoEmbeddingsCollection() string {
	if name := strings.TrimSpace(os.Getenv("ATHENA_MONGODB_EMBEDDINGS_COLLECTION")); name != "" {
		return name
	}
	return "memory_embeddings"
}

func mongoRootMarker(root string) string {
	return filepath.Clean(root)
}

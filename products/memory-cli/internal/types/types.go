package types

const (
	DefaultSchema = "1.0"
)

type IndexFile struct {
	SchemaVersion string       `json:"schema_version"`
	UpdatedAt     string       `json:"updated_at"`
	Entries       []IndexEntry `json:"entries"`
}

type IndexEntry struct {
	ID           string `json:"id"`
	Type         string `json:"type"`
	Domain       string `json:"domain"`
	Path         string `json:"path"`
	MetadataPath string `json:"metadata_path"`
	Status       string `json:"status"`
	UpdatedAt    string `json:"updated_at"`
	Title        string `json:"title"`
}

type MetadataFile struct {
	SchemaVersion string     `json:"schema_version"`
	ID            string     `json:"id"`
	Title         string     `json:"title"`
	Status        string     `json:"status"`
	UpdatedAt     string     `json:"updated_at"`
	SourceRef     string     `json:"source_ref,omitempty"`
	SourceKind    string     `json:"source_kind,omitempty"`
	SourceType    string     `json:"source_type,omitempty"`
	Review        ReviewMeta `json:"review"`
}

type ReviewMeta struct {
	ReviewedBy   string `json:"reviewed_by"`
	ReviewedAt   string `json:"reviewed_at"`
	Decision     string `json:"decision"`
	DecisionNote string `json:"decision_notes"`
}

type RetrieveResult struct {
	SelectedID    string              `json:"selected_id"`
	SelectionMode string              `json:"selection_mode"`
	SourcePath    string              `json:"source_path"`
	Confidence    float64             `json:"confidence"`
	Reason        string              `json:"reason"`
	FallbackUsed  bool                `json:"fallback_used,omitempty"`
	SemanticHit   bool                `json:"semantic_hit,omitempty"`
	PrecisionHint float64             `json:"precision_hint,omitempty"`
	Candidates    []RetrieveCandidate `json:"candidates,omitempty"`
}

type RetrieveCandidate struct {
	ID             string  `json:"id"`
	SourcePath     string  `json:"source_path"`
	SelectionMode  string  `json:"selection_mode"`
	Confidence     float64 `json:"confidence"`
	LexicalScore   float64 `json:"lexical_score,omitempty"`
	EmbeddingScore float64 `json:"embedding_score,omitempty"`
	FusedScore     float64 `json:"fused_score,omitempty"`
	HasVector      bool    `json:"has_vector,omitempty"`
	Reason         string  `json:"reason,omitempty"`
}

type EmbeddingRecord struct {
	EntryID     string    `json:"entry_id"`
	Vector      []float64 `json:"vector"`
	ModelID     string    `json:"model_id,omitempty"`
	Provider    string    `json:"provider,omitempty"`
	Dim         int       `json:"dim,omitempty"`
	ContentHash string    `json:"content_hash,omitempty"`
	CommitSHA   string    `json:"commit_sha,omitempty"`
	SessionID   string    `json:"session_id,omitempty"`
	GeneratedAt string    `json:"generated_at,omitempty"`
	LastUpdated string    `json:"updated_at,omitempty"`
}

type MutationAuditRecord struct {
	SchemaVersion string   `json:"schema_version"`
	ID            string   `json:"id"`
	Stage         string   `json:"stage"`
	Decision      string   `json:"decision"`
	ReviewedBy    string   `json:"reviewed_by"`
	ReviewedAt    string   `json:"reviewed_at"`
	DecisionNotes string   `json:"decision_notes"`
	Reason        string   `json:"reason"`
	Risk          string   `json:"risk"`
	ReworkNotes   string   `json:"rework_notes,omitempty"`
	ReReviewedBy  string   `json:"re_reviewed_by,omitempty"`
	ChangedFiles  []string `json:"changed_files"`
	Applied       bool     `json:"applied"`
}

type TelemetryEvent struct {
	EventName       string  `json:"event_name"`
	EventVersion    string  `json:"event_version"`
	TimestampUTC    string  `json:"timestamp_utc"`
	SessionID       string  `json:"session_id"`
	TraceID         string  `json:"trace_id"`
	ScenarioID      string  `json:"scenario_id"`
	Operation       string  `json:"operation"`
	Result          string  `json:"result"`
	PolicyGate      string  `json:"policy_gate"`
	MemoryType      string  `json:"memory_type"`
	LatencyMS       int64   `json:"latency_ms"`
	SelectedID      string  `json:"selected_id,omitempty"`
	SelectionMode   string  `json:"selection_mode,omitempty"`
	SourcePath      string  `json:"source_path,omitempty"`
	OperatorVerdict string  `json:"operator_verdict"`
	ErrorCode       string  `json:"error_code,omitempty"`
	Reason          string  `json:"reason,omitempty"`
	FallbackUsed    bool    `json:"fallback_used,omitempty"`
	SemanticHit     bool    `json:"semantic_hit,omitempty"`
	PrecisionProxy  float64 `json:"precision_proxy,omitempty"`
	SemanticHitRate float64 `json:"semantic_hit_rate,omitempty"`
	FallbackRate    float64 `json:"fallback_rate,omitempty"`
}

type SnapshotChecksum struct {
	Path   string `json:"path"`
	SHA256 string `json:"sha256"`
}

type SnapshotManifest struct {
	SnapshotID    string             `json:"snapshot_id"`
	CreatedAt     string             `json:"created_at"`
	CreatedBy     string             `json:"created_by"`
	SchemaVersion string             `json:"schema_version"`
	IndexVersion  string             `json:"index_version"`
	Scope         string             `json:"scope"`
	Reason        string             `json:"reason"`
	Checksums     []SnapshotChecksum `json:"checksums"`
	PayloadRefs   []string           `json:"payload_refs"`
}

type SnapshotListRow struct {
	SnapshotID    string `json:"snapshot_id"`
	CreatedAt     string `json:"created_at"`
	CreatedBy     string `json:"created_by"`
	SchemaVersion string `json:"schema_version"`
	IndexVersion  string `json:"index_version"`
	Scope         string `json:"scope"`
	Reason        string `json:"reason"`
}

type SnapshotAuditEvent struct {
	EventName   string `json:"event_name"`
	SnapshotID  string `json:"snapshot_id"`
	SessionID   string `json:"session_id"`
	TraceID     string `json:"trace_id"`
	Result      string `json:"result"`
	Timestamp   string `json:"timestamp_utc"`
	ErrorCode   string `json:"error_code,omitempty"`
	Description string `json:"description,omitempty"`
}

type WritePolicyInput struct {
	Stage        string
	Reviewer     string
	ApprovedFlag bool
	Decision     string
	Notes        string
	Reason       string
	Risk         string
	ReworkNotes  string
	ReReviewedBy string
}

type WritePolicyDecision struct {
	Decision     string
	Reviewer     string
	Notes        string
	Reason       string
	Risk         string
	ReworkNotes  string
	ReReviewedBy string
}

type UpsertEntryInput struct {
	ID         string
	Title      string
	Type       string
	Domain     string
	Body       string
	BodyFile   string
	Stage      string
	SourceRef  string
	SourceKind string
	SourceType string
}

type BootstrapMemoryEntry struct {
	ID            string  `json:"id"`
	SelectionMode string  `json:"selection_mode"`
	SourcePath    string  `json:"source_path"`
	Confidence    float64 `json:"confidence"`
	Reason        string  `json:"reason"`
	Type          string  `json:"type"`
	Domain        string  `json:"domain"`
	Title         string  `json:"title"`
}

type BootstrapPayload struct {
	Repo          string                 `json:"repo"`
	SessionID     string                 `json:"session_id"`
	Scenario      string                 `json:"scenario"`
	GeneratedAt   string                 `json:"generated_at"`
	MemoryEntries []BootstrapMemoryEntry `json:"memory_entries"`
}

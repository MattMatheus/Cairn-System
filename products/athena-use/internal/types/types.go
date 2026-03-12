package types

type Registry struct {
	Version int
	Tools   []Tool
}

type Tool struct {
	ID            string
	Name          string
	Description   string
	Tags          []string
	StageAffinity []string
	CredentialRef string
	SupportTier   string
	Call          ToolCall
	Schema        []SchemaField
}

type ToolCall struct {
	Type    string
	Method  string
	URL     string
	Command string
}

type SchemaField struct {
	Name        string
	Type        string
	Required    bool
	Description string
	Enum        []string
}

package types

type Registry struct {
	Version int          `yaml:"version" json:"version"`
	Systems []ToolSystem `yaml:"tools" json:"tools"`
}

type ToolSystem struct {
	ID            string       `yaml:"id" json:"id"`
	Name          string       `yaml:"name" json:"name"`
	Description   string       `yaml:"description" json:"description"`
	Status        string       `yaml:"status" json:"status"`
	Tags          []string     `yaml:"tags" json:"tags"`
	Guidance      string       `yaml:"guidance" json:"guidance"`
	Complements   []string     `yaml:"complements" json:"complements"`
	CredentialRef string       `yaml:"credential" json:"credential_ref"`
	SupportTier   string       `yaml:"support_tier" json:"support_tier"`
	Capabilities  []Capability `yaml:"capabilities" json:"capabilities"`
}

type Capability struct {
	ID            string        `yaml:"id" json:"id"`
	Name          string        `yaml:"name" json:"name"`
	Description   string        `yaml:"description" json:"description"`
	Status        string        `yaml:"status" json:"status"`
	Availability  string        `yaml:"availability" json:"availability"`
	Tags          []string      `yaml:"tags" json:"tags"`
	StageAffinity []string      `yaml:"stage_affinity" json:"stage_affinity"`
	Guidance      string        `yaml:"guidance" json:"guidance"`
	Call          ToolCall      `yaml:"call" json:"call"`
	Schema        []SchemaField `yaml:"schema" json:"schema"`
}

type ToolCall struct {
	Type    string `yaml:"type" json:"type"`
	Method  string `yaml:"method" json:"method"`
	URL     string `yaml:"url" json:"url"`
	Command string `yaml:"command" json:"command"`
}

type SchemaField struct {
	Name        string   `yaml:"name" json:"name"`
	Type        string   `yaml:"type" json:"type"`
	Required    bool     `yaml:"required" json:"required"`
	Description string   `yaml:"description" json:"description"`
	Enum        []string `yaml:"enum" json:"enum"`
}

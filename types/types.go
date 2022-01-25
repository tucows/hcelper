package types

import (
	"github.com/hashicorp/vault/api"
)

type VaultLDAPResponse struct {
	RequestID     string `json:"request_id"`
	LeaseID       string `json:"lease_id"`
	Renewable     bool   `json:"renewable"`
	LeaseDuration int    `json:"lease_duration"`
	Data          struct {
	} `json:"data"`
	WrapInfo interface{} `json:"wrap_info"`
	Warnings interface{} `json:"warnings"`
	Auth     struct {
		ClientToken   string   `json:"client_token"`
		Accessor      string   `json:"accessor"`
		Policies      []string `json:"policies"`
		TokenPolicies []string `json:"token_policies"`
		Metadata      struct {
			Username string `json:"username"`
		} `json:"metadata"`
		LeaseDuration int    `json:"lease_duration"`
		Renewable     bool   `json:"renewable"`
		EntityID      string `json:"entity_id"`
		TokenType     string `json:"token_type"`
		Orphan        bool   `json:"orphan"`
	} `json:"auth"`
}

type Mount struct {
	Accessor string `json:"accessor"`
	Config   struct {
		DefaultLeaseTTL int  `json:"default_lease_ttl"`
		ForceNoCache    bool `json:"force_no_cache"`
		MaxLeaseTTL     int  `json:"max_lease_ttl"`
	} `json:"config"`
	Description           string      `json:"description"`
	ExternalEntropyAccess bool        `json:"external_entropy_access"`
	Local                 bool        `json:"local"`
	Options               interface{} `json:"options"`
	SealWrap              bool        `json:"seal_wrap"`
	Type                  string      `json:"type"`
	UUID                  string      `json:"uuid"`
}

type VaultConfig struct {
	Client *api.Client
}

type ValidMount struct {
	Name string
	Path string
	Type string
}

type MountResponse struct {
	Mounts []ValidMount
}

type ExportResponse struct {
	Exports []CompiledExport
}

type CompiledExport struct {
	URL   string
	Token string
	Type  string
}

type MountConfig struct {
	Address string `json:"address"`
}

type NomadToken struct {
	Accessor_ID string `json:"accessor_id"`
	Secret_ID   string `json:"secret_id"`
}

type ConsulToken struct {
	Token string `json:"token"`
}

type EngineRoles struct {
	Data struct {
		Keys []string `json:"keys"`
	} `json:"data"`
}

package lib

// GovernanceRecord represents a single entry in the governance_record_models.
type GovernanceRecord struct {
	ID        string  `json:"id"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
	DeletedAt *string `json:"deleted_at"` // Use pointer for nullable fields
	Identifier string `json:"identifier"`  // This will be the key in world state
	Name      string  `json:"name"`
	Status    string  `json:"status"`
}

// TrustRecord represents a single entry in the trust_record_models.
type TrustRecord struct {
	ID                     string  `json:"id"` // This will be the key in world state
	CreatedAt              string  `json:"created_at"`
	UpdatedAt              string  `json:"updated_at"`
	DeletedAt              *string `json:"deleted_at"`
	Identifier             string  `json:"identifier"` // Used for validation against GovernanceRecord
	EntityType             string  `json:"entity_type"`
	CredentialType         string  `json:"credential_type"`
	GovernanceFrameworkURI string  `json:"governance_framework_uri"`
	DIDDocument            string  `json:"did_document"`
	ValidFromDt            string  `json:"valid_from_dt"`
	ValidUntilDt           string  `json:"valid_until_dt"`
	Status                 string  `json:"status"`
	StatusDetail           string  `json:"status_detail"`
}

// Structs to handle incoming JSON lists if you submit entire lists in one transaction
// (though typically, chaincode functions handle one asset per call for simplicity)
type GovernanceRecordsWrapper struct {
	Models []GovernanceRecord `json:"governance_record_models"`
}

type TrustRecordsWrapper struct {
	Models []TrustRecord `json:"trust_record_models"`
}

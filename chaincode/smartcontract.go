package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	// "github.com/google/uuid" // We will no longer need uuid.NewRandom() in chaincode
	mylib "com.ps/trust-registry/chaincode/lib"
)

// SmartContract definition
type SmartContract struct {
	contractapi.Contract
}

// InitLedger adds a base governance record to the ledger from a JSON payload
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface, initialRecordJSON string) (string, error)  {
    var record mylib.GovernanceRecord
    err := json.Unmarshal([]byte(initialRecordJSON), &record)
    if err != nil {
       return "", fmt.Errorf("failed to unmarshal initial record JSON: %w", err)
    }

    // Check if a governance record with this identifier already exists
    recordExists, err := s.GovernanceRecordExists(ctx, record.Identifier)
    if err != nil {
       return "", err
    }
    if recordExists {
       return "", fmt.Errorf("initial governance record with identifier '%s' already exists", record.Identifier)
    }

    recordJSON, err := json.Marshal(record)
    if err != nil {
       return "", fmt.Errorf("failed to marshal InitLedger record %s: %w", record.Identifier, err)
    }

    err = ctx.GetStub().PutState(record.Identifier, recordJSON)
    if err != nil {
       return "", fmt.Errorf("failed to put InitLedger record %s: %w", record.Identifier, err)
    }

    return record.ID, nil
}
// CreateGovernanceRecord creates a new governance record on the ledger.
// It expects the ID, CreatedAt, and UpdatedAt to be provided in the input JSON.
//func (s *SmartContract) CreateGovernanceRecord(ctx contractapi.TransactionContextInterface, governanceRecordJSON string) error {
func (s *SmartContract) CreateGovernanceRecord(ctx contractapi.TransactionContextInterface, governanceRecordJSON string) (string, error) { // CORRECT SIGNATURE FOR PAYLOAD

	var record mylib.GovernanceRecord
	err := json.Unmarshal([]byte(governanceRecordJSON), &record)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal governance record JSON: %w", err)
	}

	// Validate essential fields that *must* be provided by the client
	if record.ID == "" {
		return "", fmt.Errorf("governance record ID cannot be empty")
	}
	if record.Identifier == "" {
		return "", fmt.Errorf("governance record identifier cannot be empty")
	}
	if record.CreatedAt == "" {
		return "", fmt.Errorf("governance record CreatedAt cannot be empty")
	}
	if record.UpdatedAt == "" {
		return "", fmt.Errorf("governance record UpdatedAt cannot be empty")
	}

	recordExists, err := s.GovernanceRecordExists(ctx, record.Identifier)
	if err != nil {
		return "", err
	}
	if recordExists {
		return "", fmt.Errorf("governance record with identifier '%s' already exists", record.Identifier)
	}

	recordBytes, err := json.Marshal(record)
	if err != nil {
		return "", fmt.Errorf("failed to marshal governance record to JSON: %w", err)
	}

	//return ctx.GetStub().PutState(record.Identifier, recordBytes)
	err = ctx.GetStub().PutState(record.Identifier, recordBytes)
    if err != nil {
        return "", err // Return empty string on error
    }

    // --- NEW: Return the ID on success ---
    return record.ID, nil
}

// GovernanceRecordExists checks if a governance record exists (using its identifier as key)
func (s *SmartContract) GovernanceRecordExists(ctx contractapi.TransactionContextInterface, identifier string) (bool, error) {
    data, err := ctx.GetStub().GetState(identifier)
    if err != nil {
        return false, fmt.Errorf("failed to read from world state: %v", err)
    }
    return data != nil, nil
}

// ReadGovernanceRecord reads a governance record from the ledger and returns it as a JSON string
func (s *SmartContract) ReadGovernanceRecord(ctx contractapi.TransactionContextInterface, identifier string) (string, error) {
    data, err := ctx.GetStub().GetState(identifier)
    if err != nil {
        return "", fmt.Errorf("failed to read from world state: %v", err)
    }
    if data == nil {
        return "", fmt.Errorf("governance record with identifier %s does not exist", identifier)
    }

    return string(data), nil
}

// GetAllGovernanceRecords returns all governance records found on the world state as a JSON string
 func (s *SmartContract) GetAllGovernanceRecords(ctx contractapi.TransactionContextInterface) (string, error) {
     resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
     if err != nil {
         return "", err
     }
     defer resultsIterator.Close()

     var records []*mylib.GovernanceRecord
     for resultsIterator.HasNext() {
         queryResponse, err := resultsIterator.Next()
         if err != nil {
             return "", err
         }

         // Attempt to unmarshal into a GovernanceRecord
         var governanceRecord mylib.GovernanceRecord
         err = json.Unmarshal(queryResponse.Value, &governanceRecord)
         if err == nil && governanceRecord.Identifier != "" && governanceRecord.Name != "" {
             // It's likely a GovernanceRecord if these key fields exist.
             records = append(records, &governanceRecord)
         }
         // If unmarshalling fails or key fields are missing, it's not a GovernanceRecord, so we ignore it.
     }

     recordsJSON, err := json.Marshal(records)
     if err != nil {
         return "", fmt.Errorf("failed to marshal all governance records to JSON: %v", err)
     }

     return string(recordsJSON), nil
 }

// CreateTrustRecord creates a new trust record on the ledger.
func (s *SmartContract) CreateTrustRecord(ctx contractapi.TransactionContextInterface, trustRecordJSON string) (string, error) { // CORRECT SIGNATURE FOR PAYLOAD

    var record mylib.TrustRecord
    err := json.Unmarshal([]byte(trustRecordJSON), &record)
    if err != nil {
        return "", fmt.Errorf("failed to unmarshal trust record JSON: %w", err)
    }

    // Validate essential fields that *must* be provided by the client
    if record.ID == "" {
		return "", fmt.Errorf("trust record ID cannot be empty")
	}
    if record.CreatedAt == "" {
		return "", fmt.Errorf("trust record CreatedAt cannot be empty")
	}
	if record.UpdatedAt == "" {
		return "", fmt.Errorf("trust record UpdatedAt cannot be empty")
	}

    governanceExists, err := s.GovernanceRecordExists(ctx, record.Identifier)
    if err != nil {
        return "", err
    }
    if !governanceExists {
        return "", fmt.Errorf("trust record validation failed: identifier '%s' not found in any governance record", record.Identifier)
    }

    recordBytes, err := json.Marshal(record)
    if err != nil {
        return "", fmt.Errorf("failed to marshal trust record to JSON: %w", err)
    }

    //return ctx.GetStub().PutState(record.ID, recordBytes)
    // For TrustRecord, you are putting state by record.ID
    	err = ctx.GetStub().PutState(record.ID, recordBytes)
    	if err != nil {
    		return "", err // Return empty string on error
    	}

    	// --- NEW: Return the ID on success ---
    	return record.ID, nil
}

// TrustRecordExists checks if a trust record exists (using its ID as key)
func (s *SmartContract) TrustRecordExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
    data, err := ctx.GetStub().GetState(id)
    if err != nil {
        return false, fmt.Errorf("failed to read from world state: %v", err)
    }
    return data != nil, nil
}

// ReadTrustRecord reads a trust record from the ledger and returns it as a JSON string
func (s *SmartContract) ReadTrustRecord(ctx contractapi.TransactionContextInterface, id string) (string, error) {
    data, err := ctx.GetStub().GetState(id)
    if err != nil {
        return "", fmt.Errorf("failed to read from world state: %v", err)
    }
    if data == nil {
        return "", fmt.Errorf("trust record with ID %s does not exist", id)
    }

    return string(data), nil
}

// GetAllTrustRecords returns all trust records found on the world state as a JSON string
func (s *SmartContract) GetAllTrustRecords(ctx contractapi.TransactionContextInterface) (string, error) {
    resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
    if err != nil {
        return "", err
    }
    defer resultsIterator.Close()

    var records []*mylib.TrustRecord
    for resultsIterator.HasNext() {
        queryResponse, err := resultsIterator.Next()
        if err != nil {
            return "", err
        }

        var record mylib.TrustRecord
        err = json.Unmarshal(queryResponse.Value, &record)
        if err != nil {
            return "", fmt.Errorf("failed to unmarshal record from iterator: %v", err)
        }
        records = append(records, &record)
    }

    recordsJSON, err := json.Marshal(records)
    if err != nil {
        return "", fmt.Errorf("failed to marshal all trust records to JSON: %v", err)
    }

    return string(recordsJSON), nil
}

// GetTrustRecordsByCredentialType returns all trust records found on the world state that match the provided credential type
func (s *SmartContract) GetTrustRecordsByCredentialType(ctx contractapi.TransactionContextInterface, credentialType string) (string, error) {
    resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
    if err != nil {
        return "", err
    }
    defer resultsIterator.Close()

    var records []*mylib.TrustRecord
    for resultsIterator.HasNext() {
        queryResponse, err := resultsIterator.Next()
        if err != nil {
            return "", err
        }

        var record mylib.TrustRecord
        // To handle cases where some records might be GovernanceRecords, we'll unmarshal and check for the CredentialType field.
        err = json.Unmarshal(queryResponse.Value, &record)
        if err != nil {
            // If it's not a TrustRecord, just continue to the next item
            continue
        }

        if record.CredentialType == credentialType {
            records = append(records, &record)
        }
    }

    recordsJSON, err := json.Marshal(records)
    if err != nil {
        return "", fmt.Errorf("failed to marshal filtered trust records to JSON: %v", err)
    }

    return string(recordsJSON), nil
}

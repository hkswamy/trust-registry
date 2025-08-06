I can provide you with a `README.md` file that focuses exclusively on the `trust-registry` chaincode. This document will cover its purpose, data models, chaincode functions, and how to build and deploy it on a Hyperledger Fabric network.

-----

# Trust Registry Chaincode

This project contains the `trust-registry` chaincode, a smart contract written in Go for Hyperledger Fabric. Its purpose is to provide a decentralized, immutable registry for managing governance frameworks and trust records. It is designed to be a single source of truth for trust relationships in a a distributed network.

## Table of Contents

  - [Project Overview](https://www.google.com/search?q=%23project-overview)
  - [Data Models](https://www.google.com/search?q=%23data-models)
      - [GovernanceRecord](https://www.google.com/search?q=%23governancerecord)
      - [TrustRecord](https://www.google.com/search?q=%23trustrecord)
  - [Chaincode Functions (API)](https://www.google.com/search?q=%23chaincode-functions-api)
  - [Setup and Deployment](https://www.google.com/search?q=%23setup-and-deployment)
      - [Prerequisites](https://www.google.com/search?q=%23prerequisites)
      - [Building the Chaincode](https://www.google.com/search?q=%23building-the-chaincode)
      - [Deployment to Hyperledger Fabric](https://www.google.com/search?q=%23deployment-to-hyperledger-fabric)
  - [Troubleshooting](https://www.google.com/search?q=%23troubleshooting)

-----

## Project Overview

The `trust-registry` chaincode provides core business logic for two primary types of records:

1.  **Governance Records:** Define the rules and policies of a trust framework, such as a flight operator's certification or a regulatory body's mandate. These records are identified by a unique `identifier`.
2.  **Trust Records:** Represent specific trust relationships, such as a verifiable credential issued to a particular entity under a specific governance framework. These records are identified by a unique `id`.

The chaincode ensures that these records are stored on the blockchain in an immutable and verifiable manner, providing transparency and auditability for all participants in the network.

-----

## Data Models

The chaincode uses two Go structs to define the data models for the records. Note that the JSON tags use **snake\_case** for compatibility with common API conventions.

### GovernanceRecord

Represents a governance framework.

```go
type GovernanceRecord struct {
    ID              string `json:"id"`
    CreatedAt       string `json:"created_at"`
    UpdatedAt       string `json:"updated_at"`
    DeletedAt       string `json:"deleted_at"`
    Identifier      string `json:"identifier"`
    Name            string `json:"name"`
    Status          string `json:"status"`
}
```

### TrustRecord

Represents a specific trust credential or relationship.

```go
type TrustRecord struct {
    ID                  string `json:"id"`
    CreatedAt           string `json:"created_at"`
    UpdatedAt           string `json:"updated_at"`
    DeletedAt           string `json:"deleted_at"`
    Identifier          string `json:"identifier"`
    EntityType          string `json:"entity_type"`
    CredentialType      string `json:"credential_type"`
    GovernanceFrameworkURI string `json:"governance_framework_uri"`
    DIDDocument         string `json:"did_document"`
    ValidFromDt         string `json:"valid_from_dt"`
    ValidUntilDt        string `json:"valid_until_dt"`
    Status              string `json:"status"`
    StatusDetail        string `json:"status_detail"`
}
```

-----

## Chaincode Functions (API)

The following functions are exposed by the chaincode and can be invoked from a client application.

| Function Name              | Arguments                                        | Description                                                              |
|----------------------------|--------------------------------------------------|--------------------------------------------------------------------------|
| `InitLedger`               | `(context)`                                      | Initializes the ledger with a set of default governance and trust records. |
| `CreateGovernanceRecord`   | `(context, record_json)`                         | Creates a new governance record.                                         |
| `ReadGovernanceRecord`     | `(context, identifier)`                          | Retrieves a governance record using its identifier.                      |
| `GetAllGovernanceRecords`  | `(context)`                                      | Retrieves all governance records from the ledger.                          |
| `CreateTrustRecord`        | `(context, record_json)`                         | Creates a new trust record.                                                |
| `ReadTrustRecord`          | `(context, id)`                                  | Retrieves a trust record using its unique ID.                                |
| `GetAllTrustRecords`       | `(context)`                                      | Retrieves all trust records from the ledger.                               |
| `GetTrustRecordsByCredentialType`| `(context, credential_type)`                     | Retrieves all trust records that match a specific credential type.       |

-----

## Setup and Deployment

### Prerequisites

  - **Go 1.20+**
  - **Hyperledger Fabric Samples** (`v2.x` or `v2.5.x` recommended)
  - **Docker Desktop**

### Building the Chaincode

1.  Navigate to the chaincode directory.
2.  Run `go mod tidy` to ensure all dependencies are resolved.

The chaincode is packaged as a Go module and is ready for deployment. There is no need for a separate build step for deployment to Fabric.

### Deployment to Hyperledger Fabric

Assuming you have a `test-network` running from `fabric-samples`, you can deploy this chaincode using the `network.sh` script.

1.  Make sure your Fabric network is running. If not, start it:
    ```bash
    cd /path/to/fabric-samples/test-network
    ./network.sh up createChannel -c mychannel
    ```
2.  From the `test-network` directory, run the deployment command. Replace `/path/to/chaincode` with the actual path to your chaincode folder.
    ```bash
    ./network.sh deployCC -ccn trustregistry -ccp /path/to/chaincode -ccl go -c trust-registry-channel
    ```
      * `-ccn trustregistry`: The name of the chaincode.
      * `-ccp`: Path to the chaincode's directory.
      * `-ccl go`: The language.
      * `-c trust-registry-channel`: The channel to deploy on.

-----

## Troubleshooting

  - **`Error: failed to deploy chaincode...`**:
      * Check your chaincode path (`-ccp`) in the deployment command. It must point to the directory containing the `go.mod` file and the `smartcontract.go` file.
      * Ensure the chaincode is free of syntax errors by running `go build` in its directory.
  - **`Error: chaincode response 500, governance record CreatedAt cannot be empty`**:
      * This error originates from a Go `fmt.Errorf` call in the chaincode's validation logic.
      * It means the JSON payload sent to the chaincode is missing a required field or the field is empty.
      * Make sure the client application (e.g., a REST API) is sending a valid JSON object with all required fields properly populated. For example, the `created_at` field cannot be an empty string (`""`) or null.
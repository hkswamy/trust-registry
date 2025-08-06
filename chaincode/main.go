package main

import (
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	chaincode, err := contractapi.NewChaincode(&SmartContract{}) // <-- CORRECTED LINE
	if err != nil {
		log.Panicf("Error creating trust-registry chaincode: %v", err)
	}

	// Start the chaincode
	if err := chaincode.Start(); err != nil {
		log.Panicf("Error starting trust-registry chaincode: %v", err)
	}
}
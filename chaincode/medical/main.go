package main

import "github.com/hyperledger/fabric-contract-api-go/v2/contractapi"

func main() {
	chaincode, err := contractapi.NewChaincode(&MedicalRecordContract{})
	if err != nil {
		panic(err.Error())
	}
	if err := chaincode.Start(); err != nil {
		panic(err.Error())
	}
}

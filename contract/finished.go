package main

import (
	"encoding/json"
	"errors"
	"fmt"
	
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type Product struct {

Name   string  `json:"name"`

Qunatity string `json:"qnty"`

Owner string  `json:"owner"`

Batchid string  `json:"batchid"`

}

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	err := stub.PutState("hello_world", []byte(args[0]))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Invoke isur entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "write" {
	fmt.Println("calling wrrite Method")
		return t.write(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
	}
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}

// write - invoke function to write key/value pair
func (t *SimpleChaincode) write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key string
	var err error
	fmt.Println("running write()")

	if len(args) != 4 {
		return nil, errors.New("Incorrect number of arguments. Expecting 4. name of the key and value to set")
	}
	
	key := args[3]
	
	product := Product{
	
	Name:   args[0],
	
	Qunatity: args[1],
	
	Owner: args[2],
	
	Batchid: args[3],
	
	}
	bytes, err := json.Marshal(product)
	
	if err != nil {
	
	fmt.Println("Error marshaling product")
	
	return nil, errors.New("Error marshaling product")
	
	}
fmt.Println("batchId is : " + product.Batchid)
fmt.Println("key is : " + key)
	err = stub.PutState(key, bytes)

	if err != nil {

	return nil, err

	}

	return nil, nil
	
}

// read - query function to read key/value pair
func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}

	key := args[0]
	fmt.Println("key is : " + key)
	valAsbytes, err := stub.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}
	
	
	return valAsbytes, nil
}
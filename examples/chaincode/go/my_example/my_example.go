/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

//WARNING - this chaincode's ID is hard-coded in chaincode_example04 to illustrate one way of
//calling chaincode from a chaincode. If this example is modified, chaincode_example04.go has
//to be modified as well with the new ID of chaincode_example02.
//chaincode_example05 show's how chaincode ID can be passed in as a parameter instead of
//hard-coding.

import (
	"errors"
	"fmt"
	"strings"
	"strconv"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

func (t *SimpleChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	// var A, B string    // Entities
	// var Aval, Bval int // Asset holdings
	// var err error

	// if len(args) < 3 {
	// 	return nil, errors.New("Incorrect number of arguments. Expecting at least 3 arguments")
	// }

	// Initialize the chaincode
	//args[0] contains type of material.
	//Value 0 means Raw material
	//Value 1 means an intermediatory or final product.
	// if args[0] == 0 {
	// 	rawMaterial = args[1]
	// 	//Here we can have a validity test if required.
	// 	dataOrigin, err = args[2]
	// 	if err != nil {
	// 		return nil, errors.New("Expe")
	// 	}
	// }
	
	
	// B = args[2]
	// Bval, err = strconv.Atoi(args[3])
	// if err != nil {
	// 	return nil, errors.New("Expecting integer value for asset holding")
	// }
	// fmt.Printf("Aval = %d, Bval = %d\n", Aval, Bval)

	// Write the state to the ledger
	// err = stub.PutState(A, []byte(strconv.Itoa(Aval)))
	// if err != nil {
	// 	return nil, err
	// }

	// err = stub.PutState(B, []byte(strconv.Itoa(Bval)))
	// if err != nil {
	// 	return nil, err
	// }

	return nil, nil
}

// Transaction makes payment of X units from A to B
func (t *SimpleChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {

	var err error
	var dataOrigin string

	//Checking Validity
	if len(args) < 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting at least 3 arguments")
	}

	// Initialize the chaincode
	//args[0] contains type of material.
	//Value 0 means Raw material
	//Value 1 means an Intermediatory or Final product.
	if args[0] == strconv.Itoa(0) {
		dataKey :=  args[1]
		//Here we can have a validity test for the possible dataOrigins if required.
		dataOrigin = strconv.Itoa(0) + "," + args[2]
		err = stub.PutState(dataKey, []byte(dataOrigin))
		if err != nil {
			return nil, err
		}
	} else if args[0] == strconv.Itoa(1) {
		dataKey := args[1]
		nArgs := len(args)
		dataOrigin = "1"
		var i int
		for i = 2 ; i < nArgs; i++ {
			//checkState checks if the Dependencies of an intermediate are present in the Legder
			switch checkState(stub, args[i]) {
				case -1: fmt.Printf("Warning: %v does not exist in the ledger.", args[i])
				case -2: fmt.Printf("Warning: Provenance of %v not completely available.", args[i])
			}
			dataOrigin = dataOrigin + "," + getOrigin(stub, args[i])		
		}
		//Here we can have a validity test for the possible dataOrigins if required.
		//Value is stored as: 1,NameOfRawMaterial1,NameOfRawMaterial2,...
		err = stub.PutState(dataKey, []byte(dataOrigin))
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}


// Query callback representing the query of a chaincode
func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	if function != "query" {
		return nil, errors.New("Invalid query function name. Expecting \"query\"")
	}
	var item string // Entities
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the item to query")
	}

	item = args[0]

	// Get the state from the ledger
	dataKey, err := stub.GetState(item)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + item + "\"}"
		return nil, errors.New(jsonResp)
	}

	if dataKey == nil {
		jsonResp := "{\"Error\":\"Nil Origin for " + item + "\"}"
		return nil, errors.New(jsonResp)
	}

	dataVal := string(dataKey)
	keyArray := strings.Split(dataVal, ",")
	jsonResp := "{\"Name\":\"" + item + "\",\"Type\":\"Intermediatory\",\"Origin\":\"" + keyArray[1] + "\"}"
	if keyArray[0] == strconv.Itoa(0) {
		jsonResp = "{\"Name\":\"" + item + "\",\"Type\":\"Raw Material\",\"Origin\":\"" + keyArray[1] + "\"}"
	}
	

	
	fmt.Printf("Query Response:%s\n", jsonResp)
	return dataKey, nil
}

func checkState (stub *shim.ChaincodeStub, checkString string) int {

	dataKey, err := stub.GetState(checkString)
	if err != nil {
		return -1
	}
	if dataKey == nil {
		return -1
	}

	dataVal := string(dataKey)
	keyArray := strings.Split(dataVal, ",")

	if keyArray[0] == "1" {
		var i int
		checkValid := true;
		for i = 1 ; i < len(keyArray) ; i++ {
			if checkState(stub, keyArray[i]) < 0 {
				checkValid = (checkValid && false)
			}
		}
		switch checkValid {
			case false: return -2
			case true: return 1
		}
	}

	return 1
}

func getOrigin (stub *shim.ChaincodeStub, checkString string) string {
	dataKey, err := stub.GetState(checkString)
	if err != nil {
		return "{" + checkString + " - Not in Ledger}"
	}
	if dataKey == nil {
		return "{" + checkString + " - Not in Ledger}"
	}

	dataVal := string(dataKey)
	keyArray := strings.Split(dataVal, ",")

	retString := "{" + checkString + " - "

	if keyArray[0] == "0" {
		retString = retString + keyArray[1]
	} else {
		var i int
		for i = 1 ; i < len(keyArray) ; i++ {
			convString := strings.Split(keyArray[i], " ")
			paraString := strings.TrimPrefix(convString[0], "{")
			retString = retString + getOrigin(stub, paraString)
		}
	}

	retString = retString + "}"
	return retString
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

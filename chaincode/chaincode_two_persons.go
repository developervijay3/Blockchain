package main

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// TwoPersonsChaincode example a Chaincode implementation
type TwoPersonsChaincode struct {
}

func (t *TwoPersonsChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Printf("Init called, initializing chaincode")

	var John, Andy string    // Entities
	var JohnVal, AndyVal int // Asset holdings
	var err error

	if len(args) != 4 {
		return nil, errors.New("Incorrect number of arguments. Expecting 4")
	}

	// Initialize the chaincode
	John = args[0]
	JohnVal, err = strconv.Atoi(args[1])
	if err != nil {
		return nil, errors.New("Expecting integer value for asset holding")
	}
	Andy = args[2]
	Andyval, err = strconv.Atoi(args[3])
	if err != nil {
		return nil, errors.New("Expecting integer value for asset holding")
	}
	fmt.Printf("JohnVal = %d, AndyVal = %d\n", JohnVal, AndyVal)

	// Write the state to the ledger
	err = stub.PutState(John, []byte(strconv.Itoa(JohnVal)))
	if err != nil {
		return nil, err
	}

	err = stub.PutState(Andy, []byte(strconv.Itoa(AndyVal)))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Transaction makes payment of X units from John to Andy
func (t *TwoPersonsChaincode) invoke(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Printf("Running invoke")

	var John, Andy string    // Entities
	var JohnVal, AndyVal int // Asset holdings
	var X int          // Transaction value
	var err error

	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3")
	}

	John = args[0]
	Andy = args[1]

	// Get the state from the ledger
	// TODO: will be nice to have a GetAllState call to ledger
	JohnValbytes, err := stub.GetState(John)
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	if JohnValbytes == nil {
		return nil, errors.New("Entity not found")
	}
	JohnVal, _ = strconv.Atoi(string(JohnValbytes))

	AndyValbytes, err := stub.GetState(Andy)
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	if AndyValbytes == nil {
		return nil, errors.New("Entity not found")
	}
	AndyVal, _ = strconv.Atoi(string(AndyValbytes))

	// Perform the execution
	X, err = strconv.Atoi(args[2])
	JohnVal = JohnVal - X
	AndyVal = AndyVal + X
	fmt.Printf("JohnVal = %d, AndyVal = %d\n", JohnVal, AndyVal)

	// Write the state back to the ledger
	err = stub.PutState(John, []byte(strconv.Itoa(JohnVal)))
	if err != nil {
		return nil, err
	}

	err = stub.PutState(Andy, []byte(strconv.Itoa(AndyVal)))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Deletes an entity from state
func (t *TwoPersonsChaincode) delete(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Printf("Running delete")

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3")
	}

	John := args[0]

	// Delete the key from the state in ledger
	err := stub.DelState(John)
	if err != nil {
		return nil, errors.New("Failed to delete state")
	}

	return nil, nil
}

// Invoke callback representing the invocation of a chaincode
// This chaincode will manage two accounts John and Andy and will transfer X units from John to Andy upon invoke
func (t *TwoPersonsChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Printf("Invoke called, determining function")

	// Handle different functions
	if function == "invoke" {
		// Transaction makes payment of X units from John to Andy
		fmt.Printf("Function is invoke")
		return t.invoke(stub, args)
	} else if function == "init" {
		fmt.Printf("Function is init")
		return t.Init(stub, function, args)
	} else if function == "delete" {
		// Deletes an entity from its state
		fmt.Printf("Function is delete")
		return t.delete(stub, args)
	}

	return nil, errors.New("Received unknown function invocation")
}

func (t* TwoPersonsChaincode) Run(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Printf("Run called, passing through to Invoke (same function)")

	// Handle different functions
	if function == "invoke" {
		// Transaction makes payment of X units from John to Andy
		fmt.Printf("Function is invoke")
		return t.invoke(stub, args)
	} else if function == "init" {
		fmt.Printf("Function is init")
		return t.Init(stub, function, args)
	} else if function == "delete" {
		// Deletes an entity from its state
		fmt.Printf("Function is delete")
		return t.delete(stub, args)
	}

	return nil, errors.New("Received unknown function invocation")
}

// Query callback representing the query of a chaincode
func (t *TwoPersonsChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Printf("Query called, determining function")

	if function != "query" {
		fmt.Printf("Function is query")
		return nil, errors.New("Invalid query function name. Expecting \"query\"")
	}
	var John string // Entities
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the person to query")
	}

	John = args[0]

	// Get the state from the ledger
	JohnValbytes, err := stub.GetState(John)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + John + "\"}"
		return nil, errors.New(jsonResp)
	}

	if JohnValbytes == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + John + "\"}"
		return nil, errors.New(jsonResp)
	}

	jsonResp := "{\"Name\":\"" + John + "\",\"Amount\":\"" + string(JohnValbytes) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return JohnValbytes, nil
}

func main() {
	err := shim.Start(new(TwoPersonsChaincode))
	if err != nil {
		fmt.Printf("Error starting TwoPersonsChaincode chaincode: %s", err)
	}
}

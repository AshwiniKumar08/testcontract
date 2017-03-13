package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var logger = shim.NewLogger("CLDChaincode")

const AUTHORITY = "regulator"
const MANUFACTURER = "manufacturer"
const FARMER = "farmer"
const RETAILER = "walmart"
const SLAUGHTERHOUSE = "slaughterhouse"

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

func (t *SimpleChaincode) get_username(stub shim.ChaincodeStubInterface) (string, error) {

	username, err := stub.ReadCertAttribute("username")
	if err != nil {
		return "", errors.New("Couldn't get attribute 'username'. Error: " + err.Error())
	}
	return string(username), nil
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("Initializing cotton id collection")

	var blank []string
	blankBytes, _ := json.Marshal(&blank)

	err := stub.PutState("cottonids", blankBytes)
	err = stub.PutState("fabricids", blankBytes)
	err = stub.PutState("productids", blankBytes)

	if err != nil {
		fmt.Println("Failed to initialize cotton Id collection")
	}

	fmt.Println("Initialization complete")
	return nil, nil
}

// Invoke isur entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "createCotton" {
		return t.createCotton(stub, args)
	} else if function == "createCattleTransfer" {
		return t.createCattleTransfer(stub, args)
	} else if function == "createRM" {
		return t.createRM(stub, args)
	} else if function == "createBatch" {
		return t.createBatch(stub, args)
	} else if function == "createFoodPack" {
		return t.createFoodPack(stub, args)
	} else if function == "updateHdr" {
		return t.updateHdr(stub, args)
	} //updateHdr

	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "getCotton" {
		return t.getCotton(stub, args)
	} else if function == "getAllCotton" {
		return t.getAllCotton(stub, args)
	} else if function == "getCattleTrans" {
		return t.getCattleTrans(stub, args)
	} else if function == "getAllRM" {
		return t.getAllRM(stub, args)
	}

	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}

// Peer one functions

type Cattle struct {
	Species     string  `json:"species"`
	CattleType  string  `json:"cattletype"`
	CattleId    string  `json:"cattleid"`
	CattleTag   string  `json:"cattletag"`
	Birthdate   string  `json:"birthdate"`
	Weight      float64 `json:"weight"`
	FarmerId    string  `json:"farmerid"`
	Status      string  `json:"status"`
	Certificate string  `json:"certificate"`
}

type Farmer struct {
	Cattle []string `json:"cattle"`
}

type Cotton struct {
	LotId	string  `json:"lotId"`
	Quantity	string  `json:"quatity"`
	MaterialType	string  `json:"materialType"`
	WaterConsumption	string  `json:"waterConsumption"`
	ChemicalUsed	string  `json:"chemicalUsed"`
	ContainerType	string `json:"containerType"`
	FarmCode	string  `json:"farmCode"`
	FarmStatus      string  `json:"farmStatus"`
	FarmName	string  `json:"farmName"`
	FarmNumber	string  `json:"farmNumber"`
	FarmType	string  `json:"farmType"`
}


type CottonFarmer struct {
	Cotton []string `json:"cotton"`
}


type CattleHeader struct {
	Blockheader []string `json:"blockheader"`
}

func (t *SimpleChaincode) createCotton(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	var cottonLotId string

	fmt.Println("createCotton Starts")

	//weight, err := strconv.ParseFloat(args[5], 64)

	//if args[6] != FARMER { // Only the farmer can create a cattle
	//	return nil, errors.New(fmt.Sprintf("Permission Denied. Create Cattle. %v === %v", args[6], FARMER))
	//}

	bytes, err := stub.GetState(args[0])

	if bytes != nil {
		err = json.Unmarshal(bytes, &cottonLotId)

		if cottonLotId != "" {
			return nil, errors.New(fmt.Sprintf("Cotton Lot Already Present"))
		}
	}

	cottonCorp := Cotton{
		LotId	:     args[0],
		Quantity:  args[1],
		MaterialType:    args[2],
		WaterConsumption:   args[3],
		ChemicalUsed:   args[4],
		ContainerType:      args[5],
		FarmCode:    args[6],
		FarmStatus:      args[7],
		FarmName: args[8],
		FarmNumber: args[9],
		FarmType: args[10],
	}

	bytes, err = json.Marshal(&cottonCorp)

	if err != nil {
		return nil, err
	}

	err = stub.PutState(cottonCorp.LotId, bytes)

	if err != nil {
		return nil, err
	}

	bytes, err = stub.GetState("cottonids")

	if err != nil {
		return nil, errors.New("Unable to get cattleids")
	}

	// Create Cattle List
	var cottons CottonFarmer

	err = json.Unmarshal(bytes, &cottons)

	if err != nil {
		return nil, errors.New("Corrupt Farmer record")
	}

	cottons.Cotton = append(cottons.Cotton, cottonCorp.LotId)

	fmt.Println("cottons is : " + cottons)

	bytes, err = json.Marshal(cottons)

	err = stub.PutState("cottonids", bytes)

	if err != nil {
		return nil, errors.New("Unable to put the state")
	}
	// Create Empty Blockheader list
	//var blank []string
	//blankBytes, _ := json.Marshal(&blank)
	//var cattletaghdr string

	//cattletaghdr = "cattlehdr-" + args[3]
	// Create Block Header json
	//headerBlock := "\"block\":\"" + args[8] + "\", " // Variables to define the JSON
	//headerType := "\"type\":\"CREATE\", "
	//headerValue := "\"value\":\"" + args[9] + "\", "
	//prevHash := "\"prevHash\":\"" + args[10] + "\""

	//headerjson := "{" + headerBlock + headerType + headerValue + prevHash + "}"

	// save Blockheader
	//var cattleheaders CattleHeader

	//err = json.Unmarshal(blankBytes, &cattleheaders)
	//cattleheaders.Blockheader = append(cattleheaders.Blockheader, headerjson)

	//bytes, err = json.Marshal(cattleheaders)
	//err = stub.PutState(cattletaghdr, bytes)

	fmt.Println("createCotton Ends")
	return nil, nil
}

// read cattle
func (t *SimpleChaincode) getCotton(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error
	
	fmt.Println("getCotton start")
	
	key = args[0]
	valAsbytes, err := stub.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}

	fmt.Println("getCotton Ends")
	return valAsbytes, nil
}

// Get all cotton
func (t *SimpleChaincode) getAllCotton(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var jsonResp string
	var err error

	valAsbytes, err := stub.GetState("cottonids")
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for cottonids \"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}

// Get all Cattle Transaction
func (t *SimpleChaincode) getCattleTrans(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var jsonResp string
	var err error

	valAsbytes, err := stub.GetState("cattlehdr-" + args[0])
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for Cattle Transactions \"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}

type Batch struct {
	Batchid   string `json:"batchid"`
	Taglist   string `json:"taglist"`
	Batchhdr  string `json:"batchhdr"`
	Batchdate string `json:"batchdate"`
	Source    string `json:"source"`
	SourceHdr string `json:"sourcehdr"`
}

type BatchList struct {
	Batch []string
}

// Create Batch
func (t *SimpleChaincode) createBatch(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	// args 0 = "Farmername", 1 = "Batchid" , 2 = "[\"T01\",\"T02\"]" list of tags , 3 = batch date, 3 = sourcehdr

	batchkey := "batchids-" + args[0]

	batchidbytes, _ := stub.GetState(batchkey)

	// Create/update soruce's Batch List
	var batchlist BatchList

	err := json.Unmarshal(batchidbytes, &batchlist)

	batchlist.Batch = append(batchlist.Batch, args[1])

	batchidbytes, err = json.Marshal(batchlist)

	err = stub.PutState(batchkey, batchidbytes)

	// Create Batch

	batch := Batch{
		Batchid:   args[1],
		Taglist:   args[2],
		Batchhdr:  batchkey,
		Batchdate: args[3],
		Source:    args[0],
		SourceHdr: args[4],
	}

	bytes, _ := json.Marshal(&batch)

	err = stub.PutState(batch.Batchid, bytes)

	if err != nil {
		return nil, errors.New("Corrupt Transaction record")
	}

	return nil, nil
}

// Create Cattle Transfer
type TransferDetail struct {
	Transfer []string `json:"transfer"`
}

func (t *SimpleChaincode) createCattleTransfer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	// Creat or update Transaction in From side
	var transferFromdetails TransferDetail

	transferFrombytes, err := stub.GetState(args[1])

	err = json.Unmarshal(transferFrombytes, &transferFromdetails)

	transferFromdetails.Transfer = append(transferFromdetails.Transfer, args[0])
	transferFrombytes, err = json.Marshal(transferFromdetails)
	err = stub.PutState(args[1], transferFrombytes)

	// Creat or update Transaction in To side
	var transferTodetails TransferDetail
	transferTobytes, err := stub.GetState(args[2])

	err = json.Unmarshal(transferTobytes, &transferTodetails)

	transferTodetails.Transfer = append(transferTodetails.Transfer, args[0])
	transferTobytes, err = json.Marshal(transferTodetails)
	err = stub.PutState(args[2], transferTobytes)

	if err != nil {
		return nil, errors.New("Corrupt Transaction record")
	}

	return nil, nil

}

func (t *SimpleChaincode) updateHdr(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	// Create and Update Cattle Header, Rawmeat and food pkg
	hdr := args[1]

	bytes, err := stub.GetState(hdr)

	if err != nil {
		return nil, errors.New("Corrupt Transaction record")
	}

	var headers CattleHeader

	err = json.Unmarshal(bytes, &headers)
	headers.Blockheader = append(headers.Blockheader, args[2])

	bytes, err = json.Marshal(headers)
	err = stub.PutState(hdr, bytes)

	return nil, nil
}

type Rawmeat struct {
	RawmeatId   string  `json:"rawmeatid"`
	Weight      float64 `json:"weight"`
	CreatedDate string  `json:"createddate"`
	SourceTag   string  `json:"sourcetag"`
	ExpireDate  string  `json:"expiredate"`
	Temperature string  `json:"temperature"`
	Company     string  `json:"company"`
	Certificate string  `json:"certificate"`
}

type Slaughter struct {
	Rawmeat []string `json:"rawmeats"`
}

// Get all cattle
func (t *SimpleChaincode) getAllRM(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var jsonResp string
	var err error

	valAsbytes, err := stub.GetState("rmids")
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for rmids \"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}

// Peer Two function
func (t *SimpleChaincode) createRM(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Initializing Raw meat Creation")

	weight, err := strconv.ParseFloat(args[1], 64)

	if args[6] != SLAUGHTERHOUSE { // Only the farmer can create a cattle
		return nil, errors.New(fmt.Sprintf("Permission Denied. Create rawmeat . %v === %v", args[6], SLAUGHTERHOUSE))
	}

	rawmeat := Rawmeat{
		RawmeatId:   args[0],
		Weight:      weight,
		CreatedDate: args[2],
		SourceTag:   args[3],
		ExpireDate:  args[4],
		Temperature: args[5],
		Company:     args[6],
		Certificate: args[7],
	}

	bytes, err := json.Marshal(&rawmeat)

	if err != nil {
		return nil, err
	}

	err = stub.PutState(rawmeat.RawmeatId, bytes)

	if err != nil {
		return nil, err
	}

	bytes, err = stub.GetState("rmids")

	if err != nil {
		return nil, errors.New("Unable to get rmids")
	}

	// Create Cattle List
	var rawmeats Slaughter

	err = json.Unmarshal(bytes, &rawmeats)

	if err != nil {
		return nil, errors.New("Corrupt Farmer record")
	}

	rawmeats.Rawmeat = append(rawmeats.Rawmeat, rawmeat.RawmeatId)

	bytes, err = json.Marshal(rawmeats)

	err = stub.PutState("rmids", bytes)

	if err != nil {
		return nil, errors.New("Unable to put the state")
	}
	// Create Empty Blockheader list
	var blank []string
	blankBytes, _ := json.Marshal(&blank)
	var cattletaghdr string

	cattletaghdr = "rawmeathdr-" + args[0]

	// save Blockheader
	var cattleheaders CattleHeader

	err = json.Unmarshal(blankBytes, &cattleheaders)
	cattleheaders.Blockheader = append(cattleheaders.Blockheader, args[8])

	bytes, err = json.Marshal(cattleheaders)
	err = stub.PutState(cattletaghdr, bytes)

	return nil, nil
}

type FoodPack struct {
	Foodpackid          string  `json:"foodpackid"`
	Weight              float64 `json:"weight"`
	CreatedDate         string  `json:"createddate"`
	SourceTag           string  `json:"sourcetag"`
	ExpireDate          string  `json:"expiredate"`
	Temperature         string  `json:"temperature"`
	Company             string  `json:"company"`
	PerservationProcess string  `json:"perservationprocess"`
	Certificate         string  `json:"certificate"`
	PackageType         string  `json:"packagetype"`
	Productstate        string  `json:"productstate"`
	Primalcut           string  `json:"partname"`
}

type Foodmfg struct {
	Foodpack []string `json:"foodpacks"`
}

func (t *SimpleChaincode) createFoodPack(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Initializing Food pack Creation")

	weight, err := strconv.ParseFloat(args[1], 64)

	foodpack := FoodPack{
		Foodpackid:          args[0],
		Weight:              weight,
		CreatedDate:         args[2],
		SourceTag:           args[3],
		ExpireDate:          args[4],
		Temperature:         args[5],
		Company:             args[6],
		PerservationProcess: args[7],
		Certificate:         args[8],
		PackageType:         args[9],
		Productstate:        args[10],
		Primalcut:           args[11],
	}

	bytes, err := json.Marshal(&foodpack)

	if err != nil {
		return nil, err
	}

	err = stub.PutState(foodpack.Foodpackid, bytes)

	if err != nil {
		return nil, err
	}

	bytes, err = stub.GetState("foodpackids")

	if err != nil {
		return nil, errors.New("Unable to get rmids")
	}

	// Create Cattle List
	var foodpacks Foodmfg

	err = json.Unmarshal(bytes, &foodpacks)

	if err != nil {
		return nil, errors.New("Corrupt Farmer record")
	}

	foodpacks.Foodpack = append(foodpacks.Foodpack, foodpack.Foodpackid)

	bytes, err = json.Marshal(foodpacks)

	err = stub.PutState("foodpackids", bytes)

	if err != nil {
		return nil, errors.New("Unable to put the state")
	}
	// Create Empty Blockheader list
	var blank []string
	blankBytes, _ := json.Marshal(&blank)
	var cattletaghdr string

	cattletaghdr = "foodpkghdr-" + args[0]

	// save Blockheader
	var cattleheaders CattleHeader

	err = json.Unmarshal(blankBytes, &cattleheaders)
	cattleheaders.Blockheader = append(cattleheaders.Blockheader, args[12])

	bytes, err = json.Marshal(cattleheaders)
	err = stub.PutState(cattletaghdr, bytes)

	return nil, nil

}

package main

import (
	"errors"
	"fmt"
	"strings"
	"strconv"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type Marble struct{
	Name string `json:"name"`					
	Color string `json:"color"`
	Size int `json:"size"`
	User string `json:"user"`
}


//Defining Structure for the Tokens
type Token struct{
	Id string `json:"id"`
	Colour string `json:"colour"`
	User string `json:"user"`
	Sell bool `json:"sell"`
	Value int `json:"value"`
	Trade bool `json:"trade"`
	TradeColour string `json:"tradecolour"`
	SaleString string `json:"salestring"`
	TradeString string `json:"tradestring"`
}

//Defining structure for the players
type Player struct {
	Id string `json:"id"`
	Points string `json:"points"`
	Assets string `json:"assets"`
	Name string `json:"name"`
	Photo string `json:"photo"`
}


//"<token1ID>,<token2ID>,<Token3ID>...."
var tokenList = "tokenList"
//"<User1ID>:<TokenWillingID>|<TokenWillingColour>-<ColourRequired>,<User2ID>:<TokenWillingID>|<TokenWillingColour>-<ColourRequired>"
var tradeList = "tradeList"
var userList = "userList"
var saleList = "saleList"

//Init
func (t *SimpleChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {

	//We Initialize all the lists

	//Initializing User List
	err := stub.PutState(userList, []byte(""))
	if err != nil {
		return nil, err
	}

	//Initializing Token List
	err = stub.PutState(tokenList, []byte(""))
	if err != nil {
		return nil, err
	}

	//Initializing Trade List
	err = stub.PutState(tradeList, []byte(""))
	if err != nil {
		return nil, err
	}

	//Initializing Sale List
	err = stub.PutState(saleList, []byte(""))
	if err != nil {
		return nil, err
	}

	return nil, nil
}


func (t *SimpleChaincode) Run(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	
	// Handle different functions
	if function == "init" {											
		return t.Init(stub, "init", args)
	} else if function == "delete" {								
		//return t.delete(stub, args)
	} else if function == "createToken" {							
		return t.createToken(stub, args)
	} else if function == "createUser" {								
		return t.createUser(stub, args)
	} else if function == "claimUser" {										
		//return t.setUser(stub, args)
	} else if function == "claimToken" {										
		//return t.setUser(stub, args)
	} else if function == "setTradeStatus" {									
		//return t.setTradeStatus(stub, args)
	} else if function == "setSellStatus" {										
		//return t.setSellStatus(stub, args)
	} else if function == "trade" {									
		//return t.trade(stub, args)
	} else if function == "buy" {									
		//return t.buy(stub, args)
	} else if function == "redeem" {									
		//return t.buy(stub, args)
	}

	return nil, errors.New("Received unknown function invocation")
}


//args[0] = tokenID, args[1] = Token colour
func (t *SimpleChaincode) createToken(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	//Id, Colour, User, Sell, Value, Trade, TradeColour
	
	//Checking Errors
	if len(args) < 2 {
		return nil, errors.New("Expecting At least 2 arguments")
	}
	if len(args[0]) <= 0 {
		return nil, errors.New("ID must be non null")
	}
	if len(args[1]) <= 0 {
		return nil, errors.New("Colour must be non null")
	}

	//Creating Json
	jsonString := `{"id": "` + args[0] + `", "colour": "` + args[1] + `", "user": "", "sell": "` + "false" + `", "sellvalue": "0", "trade": "false", "tradecolour" : "", "salestring": "", "tradestring": ""}`
	
	//Writing to Ledger
	err := stub.PutState(args[0], []byte(jsonString))
	if err != nil {
		return nil, err
	}
	
	//Updating Token List
	getTokenList, err := stub.GetState(tokenList)
	if err != nil {
		return nil, errors.New("Failed to get Token ID")
	}

	//Updating tokenList
	updatedTokenList := addSubstringtoString(string(getTokenList), args[0])

	//Writing Back Updated List
	err = stub.PutState(tokenList, []byte(updatedTokenList))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

//args[0] = userID
func (t *SimpleChaincode) createUser(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	//Id, Points, Assets, Name, Photo
	
	//Checking Errors
	if len(args) != 1 {
		return nil, errors.New("Expecting Only 1")
	}
	
	//Creating Json
	jsonString := `{"id": "` + args[0] + `", "points": "` + "50" + `", "assets": "", "Name": "` + "" + `", "photo": ""}`
	
	var err error 

	//Writing to Ledger
	err = stub.PutState(args[0], []byte(jsonString))
	if err != nil {
		return nil, err
	}
	
	var updatedUserList string

	//Updating Token List
	getUserList, err := stub.GetState(userList)
	if err != nil {
		return nil, errors.New("Failed to get User List")
	}
	//Writing Back Updated List
	updatedUserList = addSubstringtoString(string(getUserList), args[0])
	err = stub.PutState(userList, []byte(updatedUserList))
	if err != nil {
		return nil, err
	}
	return nil, nil
}

//args[0] = userID, args[1] = name of new User, args[2] = Photo Encoding
func (t *SimpleChaincode) claimUser(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	//Id, Points, Assets, Name, Photo
	userAsBytes, err := stub.GetState(args[0])
	if err != nil {
		return nil, errors.New("Failed to get user details")
	}
	tempUser := Player{}
	json.Unmarshal(userAsBytes, &tempUser)
	tempUser.Name = args[1]
	tempUser.Photo = args[2]

	tempUserWriteBack, _ := json.Marshal(tempUser)
	err = stub.PutState(args[0], tempUserWriteBack)

	if err != nil {
		return nil, err
	}

	return nil, nil
}

//args[0] = tokenID and args[1] is the User Claiming
func (t *SimpleChaincode) claimToken(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	//Id, Colour, User, Sell, Value, Trade, TradeColour
	tokenAsBytes, err := stub.GetState(args[0])
	if err != nil {
		return nil, errors.New("Failed to get user details")
	}
	//Edit Token Entry
	tempToken := Token{}
	json.Unmarshal(tokenAsBytes, &tempToken)

	//Set User
	tempToken.User = args[1]

	tempTokenWriteBack, _ := json.Marshal(tempToken)
	err = stub.PutState(args[0], tempTokenWriteBack)

	if err != nil {
		return nil, err
	}

	//Edit user Entry
	tempUser, err := getPlayerFromID(stub, args[1])
	if err!= nil {
		return nil, err
	}

	//Add Token to his assets
	tempUser.Assets = addSubstringtoString(tempUser.Assets, args[0])

	tempUserWriteBack, _ := json.Marshal(tempUser)
	err = stub.PutState(args[1], tempUserWriteBack)

	if err != nil {
		return nil, err
	}

	return nil, nil
}

//args[0] = UserId, args[1]  = TokenID to trade, args[2] = Colour wanted in return.
func (t *SimpleChaincode) setTradeStatus(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	
	//Get User Data
	tempUser, err := getPlayerFromID(stub, args[0])
	if err != nil {
		return nil, errors.New("Failed to get user details")
	}

	//Checking User Assets
	if searchString(tempUser.Assets, args[1]) == false {
		return nil, errors.New("Asset does not exist")
	}

	//Getting Token Data
	tempToken, err := getTokenFromID(stub, args[1])
	if err != nil {
		return nil, errors.New("Failed to get token details")
	}

	//Updating Trade List
	var putString = args[0] + ":" + args[1] + "|" + tempToken.Colour + "-" + args[2]

	//Updating Token Entry
	tempToken.Trade = true
	tempToken.TradeString = putString
	tempToken.TradeColour = args[2]

	//Writing back token to blockchain
	tempTokenWriteBack, _ := json.Marshal(tempToken)
	err = stub.PutState(args[1], tempTokenWriteBack)
	if err != nil {
		return nil, err
	}

	getTradeList, err := stub.GetState(tradeList)
	if err != nil {
		return nil, errors.New("Failed to get Trade List")
	}
	//Writing Back Updated List
	updatedTradeList := addSubstringtoString(string(getTradeList), putString)

	err = stub.PutState(tradeList, []byte(updatedTradeList))
	if err != nil {
		return nil, err
	}

	return nil,nil

}

//args[0] = UserID, args[1] = TokenID to sell, agds[2] = Value
func (t *SimpleChaincode) setSellStatus(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
		
	//Get User Data
	userAsBytes, err := stub.GetState(args[0])
	if err != nil {
		return nil, errors.New("Failed to get user details")
	}

	tempUser := Player{}
	json.Unmarshal(userAsBytes, &tempUser)

	//Checking User Assets
	if !(searchString(tempUser.Assets, args[1])) {
		return nil, errors.New("Asset does not exist for user")
	}

	//Getting Token Data
	tempToken, err := getTokenFromID(stub, args[1])
	if err != nil {
		return nil, errors.New("Token entry couldn't be fetched")
	}

	//Updating Sale List
	putString := args[0] + ":" + args[1] + "|" + tempToken.Colour + "-" + args[2]

	//Updating Token Entry
	tempToken.Sell = true
	tempToken.Value, _ = strconv.Atoi(args[2])
	tempToken.SaleString = putString
	tempTokenWriteBack, _ := json.Marshal(tempToken)
	err = stub.PutState(args[1], tempTokenWriteBack)
	if err != nil {
		return nil, err
	}

	

	getSaleList, err := stub.GetState(saleList)
	if err != nil {
		return nil, errors.New("Failed to get Saleß List")
	}
	updatedSaleList := addSubstringtoString(string(getSaleList), putString)

	//Writing Back Updated List
	err = stub.PutState(saleList, []byte(updatedSaleList))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

//args[0] = buyerID, args[1] = buyerTokenID, args[2]=buyerTokenColour args[3] = SellerID, args[4] = SellerTokenID, args[5] = SellerTokenColour
func (t *SimpleChaincode) trade(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) != 6 {
		return nil, errors.New("6 Arguements Required")
	}
	
	var tradeString = args[3] + ":" + args[4] + "|" + args[5] + "-" + args[2]
	tradeListString, err := stub.GetState(tradeList)
	if err != nil {
		return nil, err
	}

	//Searching Trade String to check if trade exists
	if (!searchString(string(tradeListString), tradeString)) {
		return nil, errors.New("Trade not found")
	}

	//Checking that both tokens belong to owners and getting Player Data
	checkA, tempBuyer, errA := userCheckTokenOwnership(stub, args[0], args[1])
	checkB, tempSeller, errB := userCheckTokenOwnership(stub, args[3], args[4])
	if errA != nil {
		return nil, errA
	}
	if errB != nil {
		return nil, errB
	}
	if (!(checkA && checkB)) {
		return nil, errors.New("Problem with token ownership")
	}

	//Creating new Trade List after omitting current trade
	var newTradeString = removeSubstringFromString(string(tradeListString), tradeString)

	//Getting Seller Token Data
	tempSellerToken, err := getTokenFromID(stub, args[4])
	if err != nil {
		return nil, err
	}

	//Getting Buyer Token Data
	tempBuyerToken, err := getTokenFromID(stub, args[1])
	if err != nil {
		return nil, err
	}

	//Removing Entries of traded tokens from Sale List
	if (tempSellerToken.Sell || tempBuyerToken.Sell) {
		getSellList, err := stub.GetState(saleList)
		if err != nil {
			return nil, errors.New("Couldn't find saleList")
		}
		var updatedSaleList string
		if (tempSellerToken.Sell) && (tempBuyerToken.Sell) {
			updatedSaleList = removeSubstringFromString(string(getSellList), tempSellerToken.SaleString)
			updatedSaleList = removeSubstringFromString(updatedSaleList, tempBuyerToken.SaleString)
		} else if (tempBuyerToken.Sell) {
			updatedSaleList = removeSubstringFromString(string(getSellList), tempBuyerToken.SaleString)
		} else {
			updatedSaleList = removeSubstringFromString(string(getSellList), tempSellerToken.SaleString)
		}
		err = stub.PutState(saleList, []byte(updatedSaleList))
	}

	//Making the trade
	tempSeller.Assets = removeSubstringFromString(tempSeller.Assets, args[4])
	tempSeller.Assets = addSubstringtoString(tempSeller.Assets, args[1])
	tempBuyer.Assets = removeSubstringFromString(tempBuyer.Assets, args[1])
	tempBuyer.Assets = addSubstringtoString(tempBuyer.Assets, args[4])
	tempSellerToken.User = args[0]
	tempSellerToken.Trade = false
	tempSellerToken.Sell = false
	tempSellerToken.SaleString = ""
	tempSellerToken.TradeString = ""
	tempSellerToken.TradeColour = ""
	tempSellerToken.Value = 0
	tempBuyerToken.User = args[3]
	tempBuyerToken.Trade = false
	tempBuyerToken.Sell = false
	tempBuyerToken.SaleString = ""
	tempBuyerToken.TradeString = ""
	tempBuyerToken.TradeColour = ""
	tempBuyerToken.Value = 0

	//Writing Back all data
	writeBack, _ := json.Marshal(tempSellerToken)
	err = stub.PutState(args[4], writeBack)
	if err != nil {
		return nil, err
	}
	writeBack, _ = json.Marshal(tempBuyerToken)
	err = stub.PutState(args[1], writeBack)
	if err != nil {
		return nil, err
	}
	writeBack, _ = json.Marshal(tempBuyer)
	err = stub.PutState(args[0], writeBack)
	if err != nil {
		return nil, err
	}
	writeBack, _ = json.Marshal(tempSeller)
	err = stub.PutState(args[3], writeBack)
	if err != nil {
		return nil, err
	}
	err = stub.PutState(tradeList, []byte(newTradeString))
	if err != nil {
		return nil, err
	}

	//Trade Complete
	return nil, nil
}

//args[0] = buyerID, args[1] = SellerID, args[2] = SellerTokenID, args[3] = SellerTokenColour, args[4] = Value
func (t *SimpleChaincode) buy(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) != 5 {
		return nil, errors.New("Exactly 5 Arguements Required")
	}

	//Obtaining List from Blockchain
	saleListString, err := stub.GetState(saleList)
	if err != nil {
		return nil, err
	}

	//Searching Sell String to check if Sale exists
	var sellString = args[1] + ":" + args[2] + "|" + args[3] + "-" + args[4]
	if (!searchString(string(saleListString), sellString)) {
		return nil, errors.New("Sale not found")
	}

	//Checking ownership of Token for seller
	check, sellerUser, err := userCheckTokenOwnership(stub, args[1], args[2])
	if err != nil {
		return nil, err
	}	
	if (!check) {
		return nil, errors.New("Token doesn't belong to user")
	}
	
	//Checking if Buyer has sufficient balance to buy
	buyerUser, err := getPlayerFromID(stub, args[0])
	reqBuyerPoints, _ := strconv.Atoi(args[4])
	buyerPoints, _ := strconv.Atoi(buyerUser.Points)
	if buyerPoints < reqBuyerPoints {
		return nil, errors.New("Buyer does not have enough points")
	}
	
	sellerToken, err := getTokenFromID(stub, args[2])
	if err != nil {
		return nil, err
	}

	//Making New Sell String after omitting this sale
	var newSellString = removeSubstringFromString(string(saleListString), sellString)

	//Removing Redundant Entries from Trade List
	if (sellerToken.Trade) {
		getTradeList, err := stub.GetState(tradeList)
		if err != nil {
			return nil, errors.New("Couldn't find tradeList")
		}
		updatedTradeList := removeSubstringFromString(string(getTradeList), sellerToken.TradeString)
		err = stub.PutState(saleList, []byte(updatedTradeList))
	}


	//Performing the sale
	sellerPoints, _ := strconv.Atoi(sellerUser.Points)
	sellerPoints = sellerPoints + reqBuyerPoints
	buyerPoints = buyerPoints - reqBuyerPoints
 	sellerUser.Assets = removeSubstringFromString(sellerUser.Assets, args[2])
 	sellerUser.Points = strconv.Itoa(sellerPoints)
 	buyerUser.Assets = addSubstringtoString (buyerUser.Assets, args[2])
 	buyerUser.Points = strconv.Itoa(buyerPoints)
 	sellerToken.User = args[0]
 	sellerToken.Trade = false
 	sellerToken.TradeString = ""
 	sellerToken.Sell = false
 	sellerToken.SaleString = ""
 	

 	//Writing Back to the blockchain
 	writeBack, _ := json.Marshal(sellerToken)
	err = stub.PutState(args[2], writeBack)
	if err != nil {
		return nil, err
	}
	writeBack, _ = json.Marshal(sellerUser)
	err = stub.PutState(args[1], writeBack)
	if err != nil {
		return nil, err
	}
	writeBack, _ = json.Marshal(buyerUser)
	err = stub.PutState(args[0], writeBack)
	if err != nil {
		return nil, err
	}
	err = stub.PutState(saleList, []byte(newSellString))
	if err != nil {
		return nil, err
	}
	//Completed Sale
	return nil, nil
}

//The function to redeem combinations from the bank
func (t *SimpleChaincode) redeem(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	return nil, nil
}

func (t *SimpleChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {


	// Handle different functions
	if function == "init" {											
		return t.Init(stub, "init", args)
	} else if function == "delete" {								
		//return t.delete(stub, args)
	} else if function == "createToken" {							
		return t.createToken(stub, args)
	} else if function == "createUser" {								
		return t.createUser(stub, args)
	} else if function == "claimUser" {										
		return t.claimUser(stub, args)
	} else if function == "claimToken" {										
		return t.claimToken(stub, args)
	} else if function == "setTradeStatus" {									
		return t.setTradeStatus(stub, args)
	} else if function == "setSellStatus" {										
		return t.setSellStatus(stub, args)
	} else if function == "trade" {									
		return t.trade(stub, args)
	} else if function == "buy" {									
		return t.buy(stub, args)
	} else if function == "redeem" {									
		//return t.buy(stub, args)
	}
	fmt.Println("run did not find func: " + function)						//error

	return nil, errors.New("Received unknown function invocation")
	return nil, nil

}

func userCheckTokenOwnership (stub *shim.ChaincodeStub, userID string, tokenID string) (bool, Player, error) {

	tempUser, err := getPlayerFromID(stub, userID)
	if err != nil {
		errorUser := Player{}
		return false, errorUser, err
	}
	tempUserAssets := strings.Split(tempUser.Assets, ",")
	for i := 0 ; i < len(tempUserAssets) ; i++ {
		if tempUserAssets[i] == tokenID {
			return true, tempUser, err
		}
	}
	return false, tempUser, err
}

func removeSubstringFromString (fullString string, substring string ) string {
	returnString := ""
	listString := strings.Split(fullString, ",")
	for i := 0 ; i < len(listString) ; i++ {
		if listString[i] != substring {
			if i == 0 {
				returnString = returnString + listString[i]
			} else {
				returnString = returnString + "," + listString[i]
			}
		}
	}
	return returnString
}

func getTokenFromID (stub *shim.ChaincodeStub, tokenID string) (Token, error) {
	tokenAsBytes, err := stub.GetState(tokenID)
	if err != nil {
		errorToken := Token{}
		return errorToken, errors.New("Failed to get user details")
	}
	tempToken := Token{}
	json.Unmarshal(tokenAsBytes, &tempToken)
	return tempToken, err
}

func getPlayerFromID (stub *shim.ChaincodeStub, playerID string) (Player, error) {
	userAsBytes, err := stub.GetState(playerID)
	if err != nil {
		errorPlayer := Player{}
		return errorPlayer, err
	}
	tempUser := Player{}
	json.Unmarshal(userAsBytes, &tempUser)
	return tempUser, err
}

func addSubstringtoString (fullString string , substring string ) string {
	if fullString == "" {
		return substring
	} else {
		return fullString + "," + substring
	}
}

func searchString (fullString string , searchVal string ) bool {
	listString := strings.Split(fullString, ",")
	for i:= 0 ; i < len(listString) ; i++ {
		if listString[i] == searchVal {
			return true
		}
	}
	return false
}


func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	//Checking Number of arguements
	if len(args) != 1 {
		return nil, errors.New("Only 1 Argument")
	}

	//Getting value from blockchain
	returnVal, err := stub.GetState(args[0])
	if err != nil {
		return nil, err
	}


	fmt.Printf("Query Response:%s\n", string(returnVal))

	return returnVal, err
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting chaincode: %s", err)
	}
}



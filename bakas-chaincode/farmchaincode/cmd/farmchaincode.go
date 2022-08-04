package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/logging"
)

type SmartContract struct {
	contractapi.Contract
}

type Farm struct {
	FarmId            string  `json:"farmId"`
	OwnerUserId       string  `json:"ownerUserId"`
	Size              float64 `json:"size"`
	PlotArea          float64 `json:"plotArea"`
	OtherTrees        string  `json:"otherTrees"`
	OriginAddress     string  `json:"originAddress"`
	OriginCoordinates string  `json:"originCoordinates"`
}

type QueryResult struct {
	Key    string `json:"Key"`
	Record *Farm
}

//Add Farm
func (s *SmartContract) AddFarm(ctx contractapi.TransactionContextInterface, farmId string, ownerUserId string,
	size float64, plotArea float64, otherTrees string, originAddrees string, originCoordinates string) (string, error) {

	logger := logging.NewLogger("farmchaincode")
	logger.Infoln("Start: Calling AddFarm function.")

	key := farmId
	isExist, err := s.IsExists(ctx, key)
	if err != nil {
		return "error", err
	}

	if isExist {
		return fmt.Sprintf("AddFarm: the farm %v is already existing. ", key), nil
	}

	farm := Farm{
		FarmId:            farmId,
		OwnerUserId:       ownerUserId,
		Size:              size,
		PlotArea:          plotArea,
		OtherTrees:        otherTrees,
		OriginAddress:     originAddrees,
		OriginCoordinates: originCoordinates,
	}

	farmAsBytes, err := json.Marshal(farm)
	if err != nil {
		return "error", fmt.Errorf("AddFarm: unable to Marshal %s ", farmAsBytes)
	}

	err = ctx.GetStub().PutState(key, farmAsBytes)
	if err != nil {
		return "error", fmt.Errorf("AddFarm: unable to invoke function %s ", farmAsBytes)
	}

	return ctx.GetStub().GetTxID(), nil
}

// CHECK IF RECORD IS EXISTING IN WORLD STATE
func (s *SmartContract) IsExists(ctx contractapi.TransactionContextInterface,
	key string) (bool, error) {

	logger := logging.NewLogger("farmchaincode")
	logger.Infoln("Start: Calling IsExists function.")

	queryResult, err := ctx.GetStub().GetState(key)
	if err != nil {
		return false, nil
	}

	logger.Infoln("End: IsExists called with values of: ", key)
	return queryResult != nil, nil
}

//Get by PRODUCT ID
func (s *SmartContract) GetFarm(ctx contractapi.TransactionContextInterface,
	farmId string) (*Farm, error) {

	logger := logging.NewLogger("farmchaincode")
	logger.Infoln("Start: Calling GetFarm function.")

	if farmId == "" {
		return nil, fmt.Errorf("GetFarm: input parameters must not be empty")
	}

	queryResult, err := ctx.GetStub().GetState(farmId)
	if err != nil {
		return nil, fmt.Errorf("GetFarm: failed to read from world state: %v", err)
	}

	if queryResult == nil {
		return nil, nil
	}

	var farm Farm
	err = json.Unmarshal(queryResult, &farm)
	if err != nil {
		return nil, err
	}

	logger.Infoln("End: GetFarm called with key value of: ", farm)
	return &farm, nil
}

// GET ALL RECORDS ON THE LEDGER
func (s *SmartContract) QueryFarms(ctx contractapi.TransactionContextInterface) ([]QueryResult, error) {

	logger := logging.NewLogger("farmchaincode")
	logger.Infoln("Start: Calling QueryFarms function.")
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, fmt.Errorf("QueryFarms: unable to call GetStateByRange function")
	}

	defer resultsIterator.Close()

	var farms []QueryResult
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var farm Farm
		err = json.Unmarshal(queryResponse.Value, &farm)
		if err != nil {
			return nil, err
		}

		queryResult := QueryResult{Key: queryResponse.Key, Record: &farm}
		farms = append(farms, queryResult)
	}

	if farms == nil {
		return nil, nil
	}

	logger.Infoln("End: QueryFarms called.")
	return farms, nil
}

// DELETE ALL FARM
func (s *SmartContract) DeleteAll(ctx contractapi.TransactionContextInterface) error {

	logger := logging.NewLogger("farmchaincode")
	logger.Infoln("Start: Calling DeleteAll function.")
	data, err := s.QueryFarms(ctx)
	if err != nil {
		return err
	}

	for _, product := range data {
		ctx.GetStub().DelState(product.Key)
	}
	logger.Infoln("End: DeleteAll function called.")
	return nil
}

func main() {
	cc, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create chaincode: %s", err.Error())
		return
	}

	if err := cc.Start(); err != nil {
		fmt.Printf("Error starting chaincode: %s", err.Error())
	}
}

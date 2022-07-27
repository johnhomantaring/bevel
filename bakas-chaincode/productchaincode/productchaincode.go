package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/logging"
)

type SmartContract struct {
	contractapi.Contract
}

//INPUT ATTRIBUTES
type Product struct {
	BatchId          string  `json:"batchId"`
	OwnerUserId      string  `json:"ownerUserId"`
	OwnerType        string  `json:"ownerType"`
	FarmId           string  `json:"farmId"`
	HarvestId        string  `json:"harvestId"`
	TraderId         string  `json:"traderId"`
	WholesalerId     string  `json:"wholesalerId"`
	ProductName      string  `json:"productName"`
	TreeAge          int     `json:"treeAge"`
	TreeCount        int     `json:"treeCount"`
	Volume           float64 `json:"volume"`
	HarvestDate      string  `json:"harvestDate"`
	PricePerUnit     float64 `json:"pricePerUnit"`
	ProductUnit      string  `json:"productUnit"`
	ProductGrade     string  `json:"productGrade"`
	ProductCategory  string  `json:"productCategory"`
	TradeCategory    string  `json:"tradeCategory"`
	FermentationDate string  `json:"fermentationDate"`
	FermentationTemp float64 `json:"fermentationTemp"`
	DryingDate       string  `json:"dryingDate"`
	DryingConditions string  `json:"dryingConditions"`
	RoastingDate     string  `json:"roastingDate"`
	RoastingType     string  `json:"roastingType"`
	Status           string  `json:"status"`
	Remarks          string  `json:"remarks"`
	ModifiedDate     string  `json:"modifiedDate"`
	ModifiedBy       string  `json:"modifiedBy"`
}

const (
	Farmer     = "farmer"
	Trader     = "trader"
	Wholesaler = "wholesaler"
	Retailer   = "retailer"
)

type QueryResult struct {
	Key    string `json:"Key"`
	Record *Product
}

func (s *SmartContract) AddProduct(ctx contractapi.TransactionContextInterface,
	product Product) (string, error) {

	logger := logging.NewLogger("productchaincode")
	logger.Infoln("Start: Calling AddProduct function.")

	key := product.BatchId
	isExist, _ := s.IsExists(ctx, key)

	if isExist {
		return fmt.Sprintf("AddProduct: the product %v is already existing. ", key), nil
	}

	product = Product{
		BatchId:          product.BatchId,
		OwnerUserId:      product.OwnerUserId,
		OwnerType:        product.OwnerType,
		FarmId:           product.FarmId,
		HarvestId:        product.HarvestId,
		TraderId:         product.TraderId,
		WholesalerId:     product.WholesalerId,
		ProductName:      product.ProductName,
		TreeAge:          product.TreeAge,
		TreeCount:        product.TreeCount,
		Volume:           product.Volume,
		HarvestDate:      product.HarvestDate,
		PricePerUnit:     product.PricePerUnit,
		ProductUnit:      product.ProductUnit,
		ProductGrade:     product.ProductGrade,
		ProductCategory:  product.ProductCategory,
		TradeCategory:    product.TradeCategory,
		FermentationDate: product.FermentationDate,
		FermentationTemp: product.FermentationTemp,
		DryingDate:       product.DryingDate,
		DryingConditions: product.DryingConditions,
		RoastingDate:     product.RoastingDate,
		RoastingType:     product.RoastingType,
		Status:           product.Status,
		Remarks:          product.Remarks,
		ModifiedDate:     product.ModifiedDate,
		ModifiedBy:       product.ModifiedBy,
	}

	productAsBytes, err := json.Marshal(product)
	if err != nil {
		return "error", fmt.Errorf("AddProduct: unable to Marshal %s ", productAsBytes)
	}

	err = ctx.GetStub().PutState(key, productAsBytes)
	if err != nil {
		return "error", fmt.Errorf("AddProduct: unable to invoke function %s ", productAsBytes)
	}

	return ctx.GetStub().GetTxID(), nil
}

//UPDATE PRODUCT
func (s *SmartContract) UpdateProduct(ctx contractapi.TransactionContextInterface,
	product Product) (string, error) {

	logger := logging.NewLogger("productchaincode")
	logger.Infoln("Start: Calling UpdateProduct function.")

	key := product.BatchId

	queryResult, err := s.GetByBatchId(ctx, key)
	if err != nil {
		return "error", err
	}

	if queryResult == nil {
		return fmt.Sprintf("UpdateProduct: the product %v does not exist. ", key), nil
	}

	queryResult.OwnerUserId = product.OwnerUserId
	queryResult.OwnerType = product.OwnerType
	queryResult.FarmId = product.FarmId
	queryResult.HarvestId = product.TraderId
	queryResult.TraderId = product.HarvestId
	queryResult.WholesalerId = product.WholesalerId
	queryResult.ProductName = product.ProductName
	queryResult.TreeAge = product.TreeAge
	queryResult.TreeCount = product.TreeCount
	queryResult.Volume = product.Volume
	queryResult.HarvestDate = product.HarvestDate
	queryResult.PricePerUnit = product.PricePerUnit
	queryResult.ProductUnit = product.ProductUnit
	queryResult.ProductGrade = product.ProductGrade
	queryResult.ProductCategory = product.ProductCategory
	queryResult.TradeCategory = product.TradeCategory
	queryResult.FermentationDate = product.FermentationDate
	queryResult.FermentationTemp = product.FermentationTemp
	queryResult.DryingDate = product.DryingDate
	queryResult.DryingConditions = product.DryingConditions
	queryResult.RoastingDate = product.RoastingDate
	queryResult.RoastingType = product.RoastingType
	queryResult.Status = product.Status
	queryResult.Remarks = product.Remarks
	queryResult.ModifiedDate = product.ModifiedDate
	queryResult.ModifiedBy = product.ModifiedBy

	productAsBytes, err := json.Marshal(queryResult)
	if err != nil {
		return "error", fmt.Errorf("UpdateProduct: unable to Marshal %s ", productAsBytes)
	}

	err = ctx.GetStub().PutState(key, productAsBytes)
	if err != nil {
		return "error", fmt.Errorf("UpdateProduct: unable to invoke function %s ", productAsBytes)
	}

	return ctx.GetStub().GetTxID(), nil

}

//GET BY BATCH ID
func (s *SmartContract) GetByBatchId(ctx contractapi.TransactionContextInterface,
	batchId string) (*Product, error) {

	logger := logging.NewLogger("productchaincode")
	logger.Infoln("Start: Calling GetByBatchId function.")

	if (batchId) == "" {
		return nil, fmt.Errorf("GetByBatchId: input parameters must not be empty")
	}

	queryResult, err := ctx.GetStub().GetState(batchId)
	if err != nil {
		return nil, fmt.Errorf("GetByBatchId: failed to read from world state: %v", err)
	}

	if queryResult == nil {
		return nil, nil
	}

	var product Product
	err = json.Unmarshal(queryResult, &product)
	if err != nil {
		return nil, err
	}

	logger.Infoln("End: GetByBatchId called with key value of: ", batchId)
	return &product, nil
}

//GET ALL PRODUCTS BY FARM ID
func (s *SmartContract) GetProductsByFarmId(ctx contractapi.TransactionContextInterface,
	farmId string) ([]QueryResult, error) {

	logger := logging.NewLogger("productchaincode")
	logger.Infoln("Start: Calling GetProductsByFarmId function.")

	queryResult := fmt.Sprintf("{\"selector\":{\"farmId\":\"%s\"}}", farmId)

	queryResultsIterator, err := ctx.GetStub().GetQueryResult(queryResult)
	if err != nil {
		return nil, err
	}

	defer queryResultsIterator.Close()

	var products []QueryResult
	for queryResultsIterator.HasNext() {
		responseRange, err := queryResultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var product Product
		err = json.Unmarshal(responseRange.Value, &product)
		if err != nil {
			return nil, err
		}
		result := QueryResult{Key: responseRange.Key, Record: &product}
		products = append(products, result)
	}
	if products == nil {
		return nil, nil
	}

	logger.Infoln("End: GetProductsByFarmId called.")
	return products, nil
}

// GET ALL AVAILABLE DATA IN THE LEDGER
func (s *SmartContract) GetAvailableProducts(ctx contractapi.TransactionContextInterface,
	ownerType string) ([]QueryResult, error) {

	logger := logging.NewLogger("productchaincode")
	logger.Infoln("Start: Calling GetAvailableProducts function.")
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, fmt.Errorf("GetAvailableProducts: unable to call GetStateByRange function")
	}

	defer resultsIterator.Close()

	var products []QueryResult
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var product Product
		err = json.Unmarshal(queryResponse.Value, &product)
		if err != nil {
			return nil, err
		}

		queryResult := QueryResult{Key: queryResponse.Key, Record: &product}
		if strings.ToLower(ownerType) == Trader && strings.ToLower(product.OwnerType) == Farmer {
			products = append(products, queryResult)
		}
		if strings.ToLower(ownerType) == Wholesaler && strings.ToLower(product.OwnerType) == Trader {
			products = append(products, queryResult)
		}
		if strings.ToLower(ownerType) == Retailer {
			if strings.ToLower(product.OwnerType) == Trader || strings.ToLower(product.OwnerType) == Wholesaler {
				products = append(products, queryResult)
			}
		}
	}

	if products == nil {
		return nil, nil
	}

	logger.Infoln("End: GetAvailableProducts called.")
	return products, nil

}

// GET ALL RECORDS ON THE LEDGER
func (s *SmartContract) QueryProducts(ctx contractapi.TransactionContextInterface) ([]QueryResult, error) {

	logger := logging.NewLogger("productchaincode")
	logger.Infoln("Start: Calling QueryProducts function.")
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, fmt.Errorf("QueryProducts: unable to call GetStateByRange function")
	}

	defer resultsIterator.Close()

	var products []QueryResult
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var product Product
		err = json.Unmarshal(queryResponse.Value, &product)
		if err != nil {
			return nil, err
		}

		queryResult := QueryResult{Key: queryResponse.Key, Record: &product}
		products = append(products, queryResult)
	}

	if products == nil {
		return nil, nil
	}

	logger.Infoln("End: QueryProducts called.")
	return products, nil
}

// CHECK IF RECORD IS EXISTING IN WORLD STATE
func (s *SmartContract) IsExists(ctx contractapi.TransactionContextInterface,
	key string) (bool, error) {

	logger := logging.NewLogger("productchaincode")
	logger.Infoln("Start: Calling IsExists function.")

	queryResult, err := ctx.GetStub().GetState(key)
	if err != nil {
		return false, nil
	}

	logger.Infoln("End: IsExists called with values of: ", key)
	return queryResult != nil, nil
}

// DELETE ALL PRODUCTS
func (s *SmartContract) DeleteAll(ctx contractapi.TransactionContextInterface) error {

	logger := logging.NewLogger("productchaincode")
	logger.Infoln("Start: Calling DeleteAll function.")
	data, err := s.QueryProducts(ctx)
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

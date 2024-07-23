package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// struktur SmartContract
type SmartContract struct {
	contractapi.Contract
}

// struktur Asset
type Asset struct {
	ID                string `json:"ID"`
	TnxType           string `json:"tnxType"`
	Date              string `json:"date"`
	Amount            int    `json:"amount"`
	Owner             string `json:"owner"`
	Name              string `json:"name"`
	ProductPrice      int    `json:"productPrice"`
	Location          string `json:"location"`
	TransactionState  string `json:"transactionState"`
	DueDate           string `json:"dueDate"`
}

// struktur AssetHistory untuk menyimpan riwayat aset
type AssetHistory struct {
    TxId      string `json:"txId"`
    Timestamp string `json:"timestamp"`
    IsDelete  bool   `json:"isDelete"`
    Asset     *Asset `json:"asset"`
}

// Fungsi untuk menginisialisasi ledger dengan beberapa data awal
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	assets := []Asset{
		{ID: "asset1", TnxType: "Buy", Date: "2024-05-04 15:04:05", Amount: 200, Owner: "Ventela", Name: "High Kick 25", ProductPrice: 20, Location: "Ventela Plant A", TransactionState: "In Plant", DueDate: "2024-05-05"},
		{ID: "asset2", TnxType: "Sell", Date: "2024-05-04 16:00:05", Amount: 100, Owner: "Bata", Name: "Holo Shoes 30", ProductPrice: 30, Location: "Bata Plant A", TransactionState: "Delivered", DueDate: "2024-05-05"},
		{ID: "asset3", TnxType: "Buy", Date: "2024-05-07 12:00:05", Amount: 400, Owner: "Ventela", Name: "Air Flow 35", ProductPrice: 10, Location: "Ventela Plant B", TransactionState: "In Plant", DueDate: "2024-05-08"},
	}

	for _, asset := range assets {
		assetJSON, err := json.Marshal(asset)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(asset.ID, assetJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state: %v", err)
		}
	}

	return nil
}

// Fungsi untuk membeli aset
func (s *SmartContract) BuyAsset(ctx contractapi.TransactionContextInterface, id string, amountStr string, owner string, name string, productPriceStr string, location string, transactionState string, dueDate string) error {
	if err := checkClientMSP(ctx, []string{"Org1MSP", "Org2MSP"}); err != nil {
		return err
	}

	amount, err := strconv.Atoi(amountStr)
	if err != nil {
		return fmt.Errorf("failed to convert amount to integer: %v", err)
	}

	productPrice, err := strconv.Atoi(productPriceStr)
	if err != nil {
		return fmt.Errorf("failed to convert productPrice to integer: %v", err)
	}

	currentTime := time.Now()
	stringTime := currentTime.Format("2006-01-02")

	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}

	if exists {
		asset, err := s.ReadAsset(ctx, id)
		if err != nil {
			return fmt.Errorf("failed to read existing asset: %v", err)
		}

		// Update aset yang sudah ada
		asset.Amount += amount
		asset.Location = location
		asset.TransactionState = transactionState
		asset.Date = stringTime
		asset.DueDate = dueDate

		updatedAssetJSON, err := json.Marshal(asset)
		if err != nil {
			return fmt.Errorf("error marshalling updated target asset JSON: %v", err)
		}

		err = ctx.GetStub().PutState(asset.ID, updatedAssetJSON)
		if err != nil {
			return fmt.Errorf("error putting updated target asset into state: %v", err)
		}

		return nil
	}

	// Membuat aset baru jika belum ada
	newAsset := &Asset{
		ID:                id,
		TnxType:           "Buy",
		Date:              stringTime,
		Amount:            amount,
		Owner:             owner,
		Name:              name,
		ProductPrice:      productPrice,
		Location:          location,
		TransactionState:  transactionState,
		DueDate:           dueDate,
	}

	newAssetJSON, err := json.Marshal(newAsset)
	if err != nil {
		return fmt.Errorf("error marshalling new asset JSON: %v", err)
	}

	err = ctx.GetStub().PutState(newAsset.ID, newAssetJSON)
	if err != nil {
		return fmt.Errorf("error putting new asset into state: %v", err)
	}

	return nil
}

// Fungsi untuk menjual aset
func (s *SmartContract) SellAsset(ctx contractapi.TransactionContextInterface, id string, amountStr string, owner string, name string, productPriceStr string, location string, transactionState string, dueDate string) error {
	if err := checkClientMSP(ctx, []string{"Org3MSP", "Org2MSP"}); err != nil {
		return err
	}

	amount, err := strconv.Atoi(amountStr)
	if err != nil {
		return fmt.Errorf("failed to convert amount to integer: %v", err)
	}

	productPrice, err := strconv.Atoi(productPriceStr)
	if err != nil {
		return fmt.Errorf("failed to convert productPrice to integer: %v", err)
	}

	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("the asset %s does not exist", id)
	}

	asset, err := s.ReadAsset(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to read existing asset: %v", err)
	}

	if asset.Amount < amount {
		return fmt.Errorf("the asset %s doesn't have enough stock", name)
	}

	// Mengurangi jumlah stok dari aset yang ada
	asset.Amount -= amount

	updatedAssetJSON, err := json.Marshal(asset)
	if err != nil {
		return fmt.Errorf("error marshalling updated target asset JSON: %v", err)
	}

	err = ctx.GetStub().PutState(asset.ID, updatedAssetJSON)
	if err != nil {
		return fmt.Errorf("error putting updated target asset into state: %v", err)
	}

	currentTime := time.Now()
	stringTime := currentTime.Format("2006-01-02")
	
	newId := id+owner

	// Membuat aset baru dengan tnxType "Sell"
	newTransactionAsset := Asset{
		ID:                newId,
		TnxType:           "Sell",
		Date:              stringTime,
		Amount:            amount,
		Owner:             owner,
		Name:              name,
		ProductPrice:      productPrice,
		Location:          location,
		TransactionState:  transactionState,
		DueDate:           dueDate,
	}

	newTransactionAssetJSON, err := json.Marshal(newTransactionAsset)
	if err != nil {
		return fmt.Errorf("error marshalling new transaction asset JSON: %v", err)
	}

	err = ctx.GetStub().PutState(newTransactionAsset.ID, newTransactionAssetJSON)
	if err != nil {
		return fmt.Errorf("error putting new transaction asset into state: %v", err)
	}

	return nil
}

// Fungsi untuk membaca aset berdasarkan ID
func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, id string) (*Asset, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", id)
	}

	var asset Asset
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

// Fungsi untuk memperbarui status aset hanya untuk field location dan transactionState.
func (s *SmartContract) UpdateAssetStatus(ctx contractapi.TransactionContextInterface, id string, location string, transactionState string) error {
	if err := checkClientMSP(ctx, []string{"Org2MSP"}); err != nil {
		return err
	}

	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", id)
	}

	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return err
	}
	if assetJSON == nil {
		return fmt.Errorf("the asset %s does not exist", id)
	}

	var asset Asset
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return err
	}

	asset.Location = location
	asset.TransactionState = transactionState

	updatedAssetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, updatedAssetJSON)
}

// Fungsi untuk menghapus aset
func (s *SmartContract) DeleteAsset(ctx contractapi.TransactionContextInterface, id string) error {
	if err := checkClientMSP(ctx, []string{"Org2MSP"}); err != nil {
		return err
	}

	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", id)
	}

	return ctx.GetStub().DelState(id)
}

// Fungsi untuk memeriksa apakah aset ada atau tidak berdasarkan ID
func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return assetJSON != nil, nil
}

// Fungsi untuk mendapatkan semua aset yang ada di world state ledger
func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*Asset, error) {
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		log.Printf("Error getting client ID: %v", err)
		return nil, err
	}
	log.Printf("Client ID: %s", clientID)

	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		log.Printf("Error getting state by range: %v", err)
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []*Asset
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			log.Printf("Error iterating over query results: %v", err)
			return nil, err
		}

		var asset Asset
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			log.Printf("Error unmarshalling asset JSON: %v", err)
			return nil, err
		}
		assets = append(assets, &asset)
	}

	return assets, nil
}

// Fungsi untuk memeriksa apakah klien memiliki MSP yang diizinkan
func checkClientMSP(ctx contractapi.TransactionContextInterface, allowedMSPs []string) error {
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("failed to get client MSP ID: %v", err)
	}

	for _, msp := range allowedMSPs {
		if clientMSPID == msp || clientMSPID == "Org2MSP" {
			return nil
		}
	}

	return fmt.Errorf("client is not authorized, allowed MSP IDs: %v, client MSP ID: %s", allowedMSPs, clientMSPID)
}

// Fungsi untuk mendapatkan riwayat aset berdasarkan ID
func (s *SmartContract) GetAssetHistory(ctx contractapi.TransactionContextInterface, id string) ([]*AssetHistory, error) {
    resultsIterator, err := ctx.GetStub().GetHistoryForKey(id)
    if err != nil {
        return nil, err
    }
    defer resultsIterator.Close()

    var records []*AssetHistory
    for resultsIterator.HasNext() {
        response, err := resultsIterator.Next()
        if err != nil {
            return nil, err
        }

        var asset Asset
        if len(response.Value) > 0 {
            err = json.Unmarshal(response.Value, &asset)
            if err != nil {
                return nil, err
            }
        } else {
            asset = Asset{
                ID: id,
            }
        }

        record := &AssetHistory{
            TxId:      response.TxId,
            Timestamp: response.Timestamp.String(),
            IsDelete:  response.IsDelete,
            Asset:     &asset,
        }
        records = append(records, record)
    }

    return records, nil
}

func main() {
	assetChaincode, err := contractapi.NewChaincode(&SmartContract{})
	if err != nil {
		log.Panicf("Error creating asset-transfer-basic chaincode: %v", err)
	}

	if err := assetChaincode.Start(); err != nil {
		log.Panicf("Error starting asset-transfer-basic chaincode: %v", err)
	}
}

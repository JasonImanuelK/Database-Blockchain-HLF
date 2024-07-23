package web

import (
	"fmt"
	"encoding/json"
	"net/http"

	"github.com/hyperledger/fabric-gateway/pkg/client"
)

type InvokeRequest struct {
	ChannelID    string   `json:"channelid"`
	ChainCodeID  string   `json:"chaincodeid"`
	Function     string   `json:"function"`
	Args         []string `json:"args"`
}

// Invoke berisi chaincode invoke requests yang mengubah ledger.
func (setup *OrgSetup) Invoke(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received Invoke request")

	/* if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %s", err)
		return
	}
	chainCodeName := r.FormValue("chaincodeid")
	channelID := r.FormValue("channelid")
	function := r.FormValue("function")
	args := r.Form["args"]
	fmt.Printf("channel: %s, chaincode: %s, function: %s, args: %s\n", channelID, chainCodeName, function, args) */

	// Parse JSON request body
	var req InvokeRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing JSON: %s", err), http.StatusBadRequest)
		return
	}

	// Extract fields from JSON request
	chainCodeName := req.ChainCodeID
	channelID := req.ChannelID
	function := req.Function
	args := req.Args

	fmt.Printf("channel: %s, chaincode: %s, function: %s, args: %v\n", channelID, chainCodeName, function, args)

	network := setup.Gateway.GetNetwork(channelID)
	contract := network.GetContract(chainCodeName)
	txn_proposal, err := contract.NewProposal(function, client.WithArguments(args...))
	if err != nil {
		fmt.Fprintf(w, "Error creating txn proposal: %s", err)
		return
	}
	txn_endorsed, err := txn_proposal.Endorse()
	if err != nil {
		fmt.Fprintf(w, "Error endorsing txn: %s", err)
		return
	}
	txn_committed, err := txn_endorsed.Submit()
	if err != nil {
		fmt.Fprintf(w, "Error submitting transaction: %s", err)
		return
	}
	fmt.Fprintf(w, "Transaction ID : %s submitted", txn_committed.TransactionID())
}

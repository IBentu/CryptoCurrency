package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	ec "github.com/IBentu/CryptoCurrency/EClib"
)

// WebServer is resposible for handling wallet (client) requests in http
type WebServer struct {
	server *NodeServer
}

// handlerSendTransaction gets the transaction of the web client, verifies it and add
// it to the transaction pool if it's ok
func (ws *WebServer) handlerSendTransaction(w http.ResponseWriter, r *http.Request) {
	body, err1 := ioutil.ReadAll(r.Body)
	trx, err2 := UnformatTransaction(body)
	if err1 == nil && err2 == nil && ws.server.node.verifyTransaction(trx) {
		ws.server.node.transactionPool.addTransaction(trx)
		w.Write([]byte("Transaction Accepted."))
	} else {
		w.Write([]byte("Transaction Rejected."))
	}
}

// handlerMine gets the mine request from the web client, verifies the signature
// and mine a block
func (ws *WebServer) handlerMine(w http.ResponseWriter, r *http.Request) {
	body, err1 := ioutil.ReadAll(r.Body)
	mineReq := &struct {
		Timestamp int64  `json:"timestamp"`
		Sign      string `json:"sign"`
	}{}
	err2 := json.Unmarshal(body, &mineReq)
	if err1 == nil && err2 == nil {
		hash := ec.ECHashString(fmt.Sprintf("%s%d", ws.server.node.pubKey, mineReq.Timestamp))
		if ec.ECVerify(hash, mineReq.Sign, ws.server.node.pubKey) {
			if ws.server.node.mine() {
				w.Write([]byte("Mined Successfully."))
			} else {
				w.Write([]byte("Could not mine."))
			}
		} else {
			w.Write([]byte("Unautherized request."))
		}
	} else {
		w.Write([]byte("Something went wrong."))
	}
}

// handlerGetBalance gets the public key from the web client checks the balance and sends the
// balance back
func (ws *WebServer) handlerGetBalance(w http.ResponseWriter, r *http.Request) {
	pk := r.URL.Query().Get("pk")
	bal := ws.server.node.checkBalance(pk)
	w.Write([]byte(strconv.Itoa(bal)))
}

// handlerWallet sends the wallet.html file to the web client
func handlerWallet(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "Web Files/wallet.html")
}

// handlerNode sends the node.html file to the web client
func handlerNode(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "Web Files/node.html")
}

// handlerFunctions sends the functions.js file to the web client
func handlerFunctions(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "Web Files/functions.js")
}

// handlerEclib sends the eclib.js file to the web client
func handlerEclib(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "Web Files/eclib.js")
}

// Start initiates the webServer. run with a goroutine
func (ws *WebServer) Start() {
	http.HandleFunc("/static/functions.js", handlerFunctions)
	http.HandleFunc("/static/eclib.js", handlerEclib)
	http.HandleFunc("/wallet", handlerWallet)
	http.HandleFunc("/node", handlerNode)
	http.HandleFunc("/api/sendTransaction", ws.handlerSendTransaction)
	http.HandleFunc("/api/mineRequest", ws.handlerMine)
	http.HandleFunc("/api/getBalance", ws.handlerGetBalance)
	http.ListenAndServe(fmt.Sprintf(":%d", ListenPort+1), nil)
}

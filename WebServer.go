package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

// WebServer is resposible for handling wallet (client) requests in http
type WebServer struct {
	server *NodeServer
}

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

func (ws *WebServer) handlerGetBalance(w http.ResponseWriter, r *http.Request) {
	pk := r.URL.Query().Get("pk")
	bal := ws.server.node.checkBalance(pk)
	w.Write([]byte(strconv.Itoa(bal)))
}

func handlerUI(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/wallet.html")
}

// Start initiates the webServer. run with a goroutine
func (ws *WebServer) Start() {

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/api/getBalance", ws.handlerGetBalance)
	http.HandleFunc("/wallet", handlerUI)
	http.HandleFunc("/api/sendTransaction", ws.handlerSendTransaction)
	http.ListenAndServe(fmt.Sprintf(":%d", ListenPort+1), nil)
}

package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"strings"
)

// WebServer is resposible for handling wallet (client) requests in http
type WebServer struct {
	server *NodeServer
}

func (ws *WebServer) handlerGetBalance(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Access-Control-Allow-Origin", "*")
	pk := r.URL.Query().Get("pk")
	XY := strings.Split(pk, "-")
	res := "Invalid Public Key"
	if len(XY) == 2 {
        X := new(big.Int)
        Y := new(big.Int)
        X, ok1 := X.SetString(XY[0], 10)
        Y, ok2 := Y.SetString(XY[1], 10)
		if ok1 && ok2 {
			key := ecdsa.PublicKey{
				X:     X,
				Y:     Y,
				Curve: elliptic.P256(),
			}
            bal := ws.server.node.checkBalance(key)
			res = strconv.Itoa(bal)
        }
	}
	w.Write([]byte(res))
}

func handlerUI(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`
<html>
    <head title="Wallet">
        <script src="/static/elliptic.min.js></script>
        <script src="/static/functions"></script>
    </head>
    Private Key: <input type="password" id="PrivateKey"></input>
    <br/>
    Public Key: <input type="input" id="PublicKey"></input>
    <br/>
    Balance: <input type="input" id="balance" readonly="true"></input>
    <input type="button" onclick="doCheckBalance()" value="Check Balance"></input>
    <br/>
    <br/>
    Money To Send: <input type="input" id="amount"></input>
    <br/>
    Recipient: <input type="input" id="amount"></input>
    <br/>
    <input type="button" onclick="sendMoney()" value="Send"></input>
</html>
    `))

}

// Start initiates the webServer. run with a goroutine
func (ws *WebServer) Start() {
    fs := http.FileServer(http.Dir("static"))
    http.Handle("/static", fs)
	http.HandleFunc("/api/getBalance", ws.handlerGetBalance)
	http.HandleFunc("/", handlerUI)
	http.ListenAndServe(fmt.Sprintf(":%d", ListenPort+1), nil)
}

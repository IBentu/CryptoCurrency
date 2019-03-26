function doCheckBalance() {
    /* 
    doCheckBalance sends a request for the balance to the node and changes the balance
    text box to the updates balance
    */
    var xhr = new XMLHttpRequest();
    xhr.onreadystatechange = function() {
        if (xhr.readyState == XMLHttpRequest.DONE) {
            var blc = document.getElementById("balance");
            blc.value = xhr.responseText
        }
    }
    var pk = document.getElementById("PublicKey").value;
    xhr.open('GET', '/api/getBalance?pk='+encodeURIComponent(pk), true);
    xhr.send(null);
}

function genKey() {
    /*
    genKey generates a new key pair and changes the text boxes tto the generated pair
    */
    var key = ec.ECGenerateKey();
    var box = document.getElementById("PrivateKey");
    box.value = key[0]
    box = document.getElementById("PublicKey");
    box.value = key[1]
}

function makeTransaction() {
    /*
    makeTransaction sends a transaction to the node based of the values in the texts boxes
    */
    var recp = document.getElementById("Recipient").value;
    var amount = document.getElementById("Amount").value;
    var privKey = document.getElementById("PrivateKey").value;
    var pubKey = document.getElementById("PublicKey").value;
    var now = new Date();
    var millis = now.getTime();
    var hash = ec.ECHashString(pubKey + recp + amount + millis);
    var signature = ec.ECSign(hash, privKey, pubKey);
    var amountInt = Number(amount)
    if(isNaN(amountInt)) {
        alert("Parameters Error!")
        return
    }
    var transaction = {
        "senderKey":pubKey,
        "recipientKey":recp,
        "amount":amountInt,
        "timestamp":millis,
        "hash":hash,
        "sign":signature,
    };
    var xhr = new XMLHttpRequest();
    xhr.onreadystatechange = function() {
        if (xhr.readyState == XMLHttpRequest.DONE) {
            
            alert(xhr.responseText);
        }
    }
    xhr.open('POST', '/api/sendTransaction', true);
    xhr.send(JSON.stringify(transaction));
}

function mine() {
    /*
    mine sends a mine request to the node
    */
    var privKey = document.getElementById("PrivateKey").value;
    var pubKey = document.getElementById("PublicKey").value;
    var now = new Date();
    var millis = now.getTime();
    var hash = ec.ECHashString(pubKey + millis);
    var signature = ec.ECSign(hash, privKey, pubKey);
    var mineRequest = {
        "timestamp":millis,
        "sign":signature,
    };
    var xhr = new XMLHttpRequest();
    xhr.onreadystatechange = function() {
        if (xhr.readyState == XMLHttpRequest.DONE) {
            alert(xhr.responseText)
        }
    }
    var pk = document.getElementById("PublicKey").value;
    xhr.open('POST', '/api/mineRequest', true);
    xhr.send(JSON.stringify(mineRequest));
}



function copyTextBox(boxId) {
    /*
    copyTextBox copies the text in a text box to the clipboard based on the id of the box
    */
    var copyText = document.getElementById(boxId);
    copyText.select();
    document.execCommand("copy");
  }
  
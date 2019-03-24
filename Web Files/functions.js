function doCheckBalance() {
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
    var key = ec.ECGenerateKey();
    var box = document.getElementById("PrivateKey");
    box.value = key[0]
    box = document.getElementById("PublicKey");
    box.value = key[1]
}

function makeTransaction() {
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
    var copyText = document.getElementById(boxId);
    copyText.select();
    document.execCommand("copy");
  }
  
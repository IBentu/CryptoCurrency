function doCheckBalance() {

    var xhr = new XMLHttpRequest();
    xhr.onreadystatechange = function() {
        if (xhr.readyState == XMLHttpRequest.DONE) {
            var blc = document.getElementById("balance");
            blc.value = xhr.responseText
        }
    }
    var pk = document.getElementById("PublicKey").value;
    xhr.open('GET', '/api/getBalance?pk='+pk, true);
    xhr.send(null);
}
// Parse for PDFs and embedded PDFs

// Get binary data

var xhr = new XMLHttpRequest();
xhr.onload = function(e) {
  var data = new Uint8Array(xhr.response);
  console.log(data);
  var hash = sha256(data);
  console.log(hash);

  PDFJS.disableWorker = true;

  PDFJS.getDocument({data: data}).then(function(doc) {
    doc.getMetadata().then(function(metadata) {
      var origin = metadata.metadata.metadata["pdfx:verification_endpoint"];
      xhr = new XMLHttpRequest();
      xhr.onload = function(e) {
        res = JSON.parse(xhr.responseText);
        inclusion_proof = res["inclusion_proof"];
        curHash = hash;
        for (node of inclusion_proof) {
          console.log(curHash);
          components = node.split("_");
          h = components[0];
          dir = components[1];
          if (dir == 'L') {
            console.log("HASHING L: " + h + " AND " + curHash);
            curHash = sha256(h + curHash);
          }else {
            console.log("HASHING R: " + curHash + " AND " + h);
            curHash = sha256(curHash + h);
          }
        }
        console.log("Calculated Hash: " + curHash);
        console.log("Received Hash: " + res["root_hash"]);
      }
      xhr.open('GET', "http://localhost:4000/verify/".concat(origin).concat("/").concat(hash));
      xhr.send();
    });
  });
}

console.log($("embed").attr("src"));
xhr.open('GET', $("embed").attr("src"), true);
xhr.responseType = 'arraybuffer';
xhr.send();

// Make request for proof of inclusion

// Compare proof of inclusion and validate signature

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
        console.log(xhr.responseText);
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

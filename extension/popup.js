// Parse for PDFs and embedded PDFs
PDFJS.workerSrc = 'pdf.worker.js';

// Make XHR to receive binary data
var xhr = new XMLHttpRequest();
xhr.onreadystatechange = function() {
    if (this.readyState == XMLHttpRequest.DONE) {
        // Hash binary data
        var data = new Uint8Array(xhr.response || xhr.mozResponseArrayBuffer);
        console.log(data);
        PDFJS.getDocument(data).then(function(pdf) {
          console.log(pdf);
        });
        console.log(pdf)
        var shaObj = new jsSHA("SHA-256", "BYTES");
        shaObj.update(this.response);
        var hash = shaObj.getHash("HEX");
        console.log(hash);
    }
}
console.log($("embed").attr("src"));
xhr.responseType = 'arraybuffer';
xhr.open('GET', $("embed").attr("src"), true);
xhr.send(null);

// Use PDF.js or PDFium to parse metadata

// Make request for proof of inclusion

// Compare proof of inclusion and validate signature

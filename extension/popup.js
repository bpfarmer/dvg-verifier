// Parse for PDFs and embedded PDFs
//PDFJS.workerSrc = 'pdf.worker.js';

//PDFJS.workerSrc = "pdf.worker.js";
// Make XHR to receive binary data

var xhr = new XMLHttpRequest();
xhr.onload = function(e) {
  var data = new Uint8Array(xhr.response);
  console.log(data);
  var hash = sha256(data);
  console.log(hash);

  PDFJS.disableWorker = true;

  PDFJS.getDocument({data: data}).then(function(doc) {
    console.log(doc.pdfInfo.metadata);
  });
}

console.log($("embed").attr("src"));
xhr.open('GET', $("embed").attr("src"), true);
xhr.responseType = 'arraybuffer';
xhr.send();

// Use PDF.js or PDFium to parse metadata

// Make request for proof of inclusion

// Compare proof of inclusion and validate signature

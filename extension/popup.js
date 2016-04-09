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

  PDFJS.workerSrc = "pdf.worker.js";
  //PDFJS.disableWorker = true;

  PDFJS.getDocument({data: data}).then(function (pdfDoc_) {
    console.log("HERE")
    pdfDoc = pdfDoc_;
    pdfDoc.getMetadata().then(function(stuff) {
        console.log(stuff);
    }).catch(function(err) {
       console.log('Error getting meta data');
       console.log(err);
    });
  }).catch(function(err) {
    console.log('Error getting PDF from ' + url);
    console.log(err);
  });
}

console.log($("embed").attr("src"));
xhr.open('GET', $("embed").attr("src"), true);
xhr.responseType = 'arraybuffer';
xhr.send();

// Use PDF.js or PDFium to parse metadata

// Make request for proof of inclusion

// Compare proof of inclusion and validate signature

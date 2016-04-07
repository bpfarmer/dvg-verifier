// Parse for PDFs and embedded PDFs
//PDFJS.workerSrc = 'pdf.worker.js';

// Make XHR to receive binary data
if (!PDFJS.workerSrc && typeof document !== 'undefined') {
  // workerSrc is not set -- using last script url to define default location
  PDFJS.workerSrc = (function () {
    'use strict';
    var scriptTagContainer = document.body ||
                             document.getElementsByTagName('head')[0];
    var pdfjsSrc = scriptTagContainer.lastChild.src;
    return pdfjsSrc && pdfjsSrc.replace(/\.js$/i, '.worker.js');
  })();
}

var xhr = new XMLHttpRequest();
xhr.onload = function(e) {
  var data = new Uint8Array(xhr.response);
  console.log(data);
  var hash = sha256(data);
  console.log(hash);

  PDFJS.getDocument({data: data}).then(function (pdfDoc_) {
        pdfDoc = pdfDoc_;
        pdfDoc.getMetadata().then(function(stuff) {
            console.log(stuff); // Metadata object here
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

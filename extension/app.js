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

intToByteArray = function(intToConvert) {
  var byteArray = new Array(4)
  for(var i = 0; i < byteArray.length; i++) {
    var byte = intToConvert & 0xff;
    byteArray[i] = byte;
    intToConvert = (intToConvert - byte) / 256 ;
  }
  return byteArray;
};

function hex2bin(str) {
  /*var bytes = new Array(str.length / 2);
  for(var i = 0; i < str.length-1; i+=2) {
    bytes.push(
      intToByteArray(
        parseInt(str.substr(i, 2), 16)
      )
    );
  }

  return bytes;*/
  var buf = new ArrayBuffer(str.length*2); // 2 bytes for each char
  var bufView = new Uint8Array(buf);
  for (var i=0, strLen=str.length; i<strLen; i++) {
    bufView[i] = str.charCodeAt(i);
  }
  return buf;
}

console.log(hex2bin("132bbbb69d7e0ee918481f073d1cb14324e58031eaa78ab2c2423c7cfedf508d"));
console.log(sha256(hex2bin("132bbbb69d7e0ee918481f073d1cb14324e58031eaa78ab2c2423c7cfedf508d")
+ hex2bin("d2d777e00435764c1c711e533c689c0c88d1ebce7dbbe57f1a2eed2b80cc153b")));

// Make request for proof of inclusion

// Compare proof of inclusion and validate signature

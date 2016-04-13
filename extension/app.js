// Parse for PDFs and embedded PDFs

// Get binary data
var xhr = new XMLHttpRequest();
xhr.onload = function(e) {
  var data = new Uint8Array(xhr.response);
  //console.log(data);
  var hash = sha256(data);
  //console.log(hash);

  PDFJS.disableWorker = true;

  PDFJS.getDocument({data: data}).then(function(doc) {
    doc.getMetadata().then(function(metadata) {
      var origin = metadata.metadata.metadata["pdfx:verification_endpoint"];
      xhr = new XMLHttpRequest();
      xhr.onload = function(e) {
        res = JSON.parse(xhr.responseText);
        if(res["error"]) {
          console.log("Invalid Document");
          displayInvalid(origin);
        }else {
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
          if(curHash == res["root_hash"]) {
            console.log("Verifying Signature: " + res["signature"]);
            //console.log(new Uint8Array(res["signature"]))
            //console.log(nacl.sign.detached.verify(hexToUa(res["root_hash"]), hexToUa(res["signature"]), hexToUa(res["public_key"])))
            console.log("Validated Proof");
            $("body").prepend("<div id='invalid-doc-error' style='width:100%;height:50px;background-color:red;'>Document is valid for origin: "+origin+".</div>");
            $('#invalid-doc-error').css("background-color", "#3366CC");
            $('#invalid-doc-error').css("background-image", "-webkit-linear-gradient(135deg, transparent,transparent 25%, hsla(0,0%,0%,.05) 25%,hsla(0,0%,0%,.05) 50%, transparent 50%,transparent 75%, hsla(0,0%,0%,.05) 75%,hsla(0,0%,0%,.05))");
            $('#invalid-doc-error').css("background-image", "-o-linear-gradient(135deg, transparent,transparent 25%, hsla(0,0%,0%,.1) 25%,hsla(0,0%,0%,.1) 50%, transparent 50%,transparent 75%, hsla(0,0%,0%,.1) 75%,hsla(0,0%,0%,.1))");
            $('#invalid-doc-error').css("background-image", "linear-gradient(135deg, transparent,transparent 25%, hsla(0,0%,0%,.1) 25%,hsla(0,0%,0%,.1) 50%, transparent 50%,transparent 75%, hsla(0,0%,0%,.1) 75%,hsla(0,0%,0%,.1))");
            $('#invalid-doc-error').css("box-shadow", "0 5px 0 hsla(0,0%,0%,.1)");
            $('#invalid-doc-error').css("text-align", "center");
            $('#invalid-doc-error').css("text-decoration", "none");
            $('#invalid-doc-error').css("color", "#f6f6f6");
            $('#invalid-doc-error').css("display", "block");
            $('#invalid-doc-error').css("font", "bold 20px/44px sans-serif");
          }else {
            displayInvalid(origin);
          }
        }
      }
      xhr.open('GET', "http://localhost:4000/verify/".concat(origin).concat("/").concat(hash));
      xhr.send();
    });
  });
}

function displayInvalid(origin) {
  $("body").prepend("<div id='invalid-doc-error' style='width:100%;height:50px;background-color:red;'>Warning. Document could not be validated with origin: "+origin+".</div>");
  $('#invalid-doc-error').css("background-color", "#c4453c");
  $('#invalid-doc-error').css("background-image", "-webkit-linear-gradient(135deg, transparent,transparent 25%, hsla(0,0%,0%,.05) 25%,hsla(0,0%,0%,.05) 50%, transparent 50%,transparent 75%, hsla(0,0%,0%,.05) 75%,hsla(0,0%,0%,.05))");
  $('#invalid-doc-error').css("background-image", "-o-linear-gradient(135deg, transparent,transparent 25%, hsla(0,0%,0%,.1) 25%,hsla(0,0%,0%,.1) 50%, transparent 50%,transparent 75%, hsla(0,0%,0%,.1) 75%,hsla(0,0%,0%,.1))");
  $('#invalid-doc-error').css("background-image", "linear-gradient(135deg, transparent,transparent 25%, hsla(0,0%,0%,.1) 25%,hsla(0,0%,0%,.1) 50%, transparent 50%,transparent 75%, hsla(0,0%,0%,.1) 75%,hsla(0,0%,0%,.1))");
  $('#invalid-doc-error').css("box-shadow", "0 5px 0 hsla(0,0%,0%,.1)");
  $('#invalid-doc-error').css("text-align", "center");
  $('#invalid-doc-error').css("text-decoration", "none");
  $('#invalid-doc-error').css("color", "#f6f6f6");
  $('#invalid-doc-error').css("display", "block");
  $('#invalid-doc-error').css("font", "bold 20px/44px sans-serif");
}

//console.log($("embed").attr("src"));
xhr.open('GET', $("embed").attr("src"), true);
xhr.responseType = 'arraybuffer';
xhr.send();

function handleFileSelect(evt) {
    evt.stopPropagation();
    evt.preventDefault();
    $(".box__error:first").hide();
    $(".box__success:first").hide();
    var files = evt.dataTransfer.files; // FileList object.

    // files is a FileList of File objects. List some properties.
    var output = [];
    for (var i = 0, f; f = files[i]; i++) {
        var reader = new FileReader();

        reader.onload = function(e) {
            var data = new Uint8Array(e.target.result);
            verify_file("stanford.edu", data);
        }
        reader.readAsArrayBuffer(f);
    }
    //document.getElementById('list').innerHTML = '<ul>' + output.join('') + '</ul>';
}

function handleDragOver(evt) {
    evt.stopPropagation();
    evt.preventDefault();
    evt.dataTransfer.dropEffect = 'copy'; // Explicitly show this is a copy.
}

function verify_file(origin, data) {
    var hash = sha256(data);
    xhr = new XMLHttpRequest();
    xhr.onload = function(e) {
        res = JSON.parse(xhr.responseText);
        console.log(res);
        if(res["error"]) {
            console.log("VALIDATION FAILED");
            $(".box__error:first").show();
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
                //console.log("Verifying Signature: " + res["signature"]);
                $(".box__success:first").show();
                $('.check_crypto').fadeTo(800, 0.2);
                $('.verify_cred').fadeTo(400, 1);
                PDFJS.disableWorker = true;
                PDFJS.getDocument({data: data}).then(function(doc) {
                    console.log("IN HERE");
                    doc.getMetadata().then(function(metadata) {
                        var origin = JSON.parse(metadata.metadata.metadata["pdfx:verification_endpoint"]);
                        console.log("EVEN IN HERE");
                        console.log(origin);
                        if(origin.length > 1) {
                            $("#multiple-origins").text("Also validated by: " + origin[1]);
                        }
                        //verify(origin, hash);
                    });
                });
                //console.log(new Uint8Array(res["signature"]))
                //console.log(nacl.sign.detached.verify(hexToUa(res["root_hash"]), hexToUa(res["signature"]), hexToUa(res["public_key"])))
            }else {
                $(".box__error:first").show();
            }
        }
    }
    xhr.open('GET', "/verify/".concat(origin).concat("/").concat(hash));
    xhr.send();
}

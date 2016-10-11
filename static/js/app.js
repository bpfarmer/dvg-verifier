var monthNames = ["January", "February", "March", "April", "May", "June",
  "July", "August", "September", "October", "November", "December"
];
var url_components = String(window.location.href).split("?");
var url_params = url_components[1].split("&");
var params = {};
for (param of url_params) {
  param_components = param.split("=");
  params[param_components[0]] = param_components[1];
}
$(document).ready(function() {
  console.log(params);
  console.log($('h3#v-name'));
  console.log(params["n"]);
  var hash_components = [params["c"], params["cred"], params["n"], params["s"], params["e"]];
  console.log(sha256(hash_components.join("__")))
  verify_parameters("stanford.edu", sha256(hash_components.join("__")))
});

function verify_parameters(origin, hash) {
  xhr = new XMLHttpRequest();
  xhr.onload = function(e) {
    if(params["n"] && params["cred"] && params["e"] && params["s"]) {
      $('h3#v-name').html(params["n"].replace("_", " "));
      $('h3#v-cred').html(params["cred"].replace(/_/g, " "));
      var start_date = new Date(params["s"].substr(0,4), params["s"].substr(4,2) - 1);
      var end_date = new Date(params["e"].substr(0,4), params["e"].substr(4,2) - 1);
      $('h3#v-date').html(monthNames[start_date.getMonth()] + " " + start_date.getFullYear() + " - " + monthNames[end_date.getMonth()] + " " + end_date.getFullYear());
    }
    res = JSON.parse(xhr.responseText);
    if(res["error"]) {
      $("h3#v-status").attr("class", "text-danger");
      $("h3#v-status").html("Invalid Credential");
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
        $("h3#v-status").attr("class", "text-success");
        $("h3#v-status").html("Valid Credential");
      }else {
        $("h3#v-status").attr("class", "text-danger");
        $("h3#v-status").html("Invalid Credential");
      }
    }
  }
  xhr.open('GET', "/verify/".concat(origin).concat("/").concat(hash));
  xhr.send();
}
//console.log(sha256())

/*
xhr.onload = function(e) {
  var data = new Uint8Array(xhr.response);
  var hash = sha256(data);
  console.log(hash);
  if((url_components[url_components.length - 1] != 'pdf')) {
      verify("stanford.edu", hash);
  }
}


  xhr.open('GET', "http://localhost:4000/verify/".concat(origin).concat("/").concat(hash));
  xhr.send();
}

function stringToUint(string) {
    var string = btoa(unescape(encodeURIComponent(string))),
        charList = string.split(''),
        uintArray = [];
    for (var i = 0; i < charList.length; i++) {
        uintArray.push(charList[i].charCodeAt(0));
    }
    return new Uint8Array(uintArray);
}

//console.log($("embed").attr("src"));
if(url_components[url_components.length - 1] == 'pdf') {
  xhr.open('GET', $("embed").attr("src"), true);
}else {
  xhr.open('GET', window.location.href, true);
}
xhr.responseType = 'arraybuffer';
xhr.send();
*/

// dec2hex :: Integer -> String
function dec2hex (dec) {
    return ('0' + dec.toString(16)).substr(-2)
}
  
// generateId :: Integer -> String
function generateId (len) {
    var arr = new Uint8Array((len || 40) / 2)
    window.crypto.getRandomValues(arr)
    return Array.from(arr, dec2hex).join('')
}

// const username = process.env.MQUNAME;
const username = "1406568753";

// const password = process.env.MQPW;
const password = "167664";

// const vhost = process.env.MQVHOST;
const vhost = "1406568753";

const host = window.location.href.replace("http://", '');

// const url = process.env.URL + "/compress";
const url = window.location.href;

var id = ""

const form = document.querySelector('form');

form.addEventListener('submit', e => {
    e.preventDefault();

    id = "urlpass"
    console.log(id)
    const url1 = document.getElementsByName('url');
    const url2 = [].map.call(url1, el => el.value).filter(function (el) {return el != ""})
    console.log(url2)
    key = generateId(32)
    const formData = new FormData();
    formData.append('url', url2.join(";"));
    formData.append('id', key);
    const uname = document.getElementsByName('uname')[0].value;
    const pword = document.getElementsByName('pw')[0].value;
    formData.append('username', uname);
    formData.append('password', pword);
    fetch(url, {
        method: 'POST',
        headers: {
            "X-ROUTING-KEY": id,
        },
        body: formData,
    }).then(function(response) {
        return response.json();
    }).then(function(json) {
        console.log(JSON.stringify(json));
        if (json.Code == 200) {
            document.getElementById("status").innerHTML = json.message;
        } else {
            document.getElementById("status").innerHTML = json.message;
        }
    }).then(function() {
        for(i = 0; i< url2.length; i++) {
            WebSocketStats(i);
        };
    }).then(function() {
        WebSocketTest();
    });
});

function WebSocketStats(i) {
    if ("WebSocket" in window) {
		// var ws_stomp_display = new SockJS(process.env.MQURL);
		var ws_stomp_display = new SockJS('http://152.118.148.103:15674/stomp');
		var client_display = Stomp.over(ws_stomp_display);
		// var mq_queue_display = "/exchange/"+ process.env.NPM + "/" + id;
		var mq_queue_display = "/exchange/"+ "1406568753-dl1" + "/" + "dlstatus"+ i;
		var on_connect_display = function() {
			console.log('connected');
			client_display.subscribe(mq_queue_display, on_message_display);
		};
		var on_error_display = function() {
			console.log('error');
		};
		var on_message_display = function(m) {
			console.log('message received');
			document.getElementById("progress" + (i+1)).innerHTML = m.body;
		};
		client_display.connect(username, password, on_connect_display, on_error_display, vhost);
	} else {
		// The browser doesn't support WebSocket
		alert("WebSocket NOT supported by your Browser!");
	}
}

function WebSocketTest() {
	if ("WebSocket" in window) {
		// var ws_stomp_display = new SockJS(process.env.MQURL);
		var ws_stomp_display = new SockJS('http://152.118.148.103:15674/stomp');
		var client_display = Stomp.over(ws_stomp_display);
		// var mq_queue_display = "/exchange/"+ process.env.NPM + "/" + id;
		var mq_queue_display = "/exchange/"+ "1406568753-frontdl" + "/" + "dlpass";
		var on_connect_display = function() {
			console.log('connected');
			client_display.subscribe(mq_queue_display, on_message_display);
		};
		var on_error_display = function() {
			console.log('error');
		};
		var on_message_display = function(m) {
			console.log('message received');
			document.getElementById("progress").innerHTML = m.body;
		};
		client_display.connect(username, password, on_connect_display, on_error_display, vhost);
	} else {
		// The browser doesn't support WebSocket
		alert("WebSocket NOT supported by your Browser!");
	}
}
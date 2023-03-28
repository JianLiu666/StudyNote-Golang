package main

import (
	"fmt"
	"html/template"
	"net/http"
	"websocketbenchmark/internal/config"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	// enable logger modules
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05-07:00",
	})

	viper.SetConfigFile("./conf.d/env.yaml")
	viper.AutomaticEnv()
}

var upgrader = websocket.Upgrader{}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logrus.Error("websocket upgrade:", err)
		return
	}
	defer c.Close()

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			logrus.Error("read:", err)
			break
		}

		logrus.Info("recv:", string(message))

		err = c.WriteMessage(mt, message)
		if err != nil {
			logrus.Error("write:", err)
			break
		}
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/echo")
}

func main() {
	conf := config.NewFromViper()

	http.HandleFunc("/echo", echo)
	http.HandleFunc("/", home)

	addr := fmt.Sprintf("%s:%s", conf.Server.Addr, conf.Server.Port)
	logrus.Fatal(http.ListenAndServe(addr, nil))
}

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>  
window.addEventListener("load", function(evt) {
    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;
    var print = function(message) {
        var d = document.createElement("div");
        d.textContent = message;
        output.appendChild(d);
        output.scroll(0, output.scrollHeight);
    };
    document.getElementById("open").onclick = function(evt) {
        if (ws) {
            return false;
        }
        ws = new WebSocket("{{.}}");
        ws.onopen = function(evt) {
            print("OPEN");
        }
        ws.onclose = function(evt) {
            print("CLOSE");
            ws = null;
        }
        ws.onmessage = function(evt) {
            print("RESPONSE: " + evt.data);
        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        return false;
    };
    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: " + input.value);
        ws.send(input.value);
        return false;
    };
    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };
});
</script>
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<p>Click "Open" to create a connection to the server, 
"Send" to send a message to the server and "Close" to close the connection. 
You can change the message and send multiple times.
<p>
<form>
<button id="open">Open</button>
<button id="close">Close</button>
<p><input id="input" type="text" value="Hello world!">
<button id="send">Send</button>
</form>
</td><td valign="top" width="50%">
<div id="output" style="max-height: 70vh;overflow-y: scroll;"></div>
</td></tr></table>
</body>
</html>
`))

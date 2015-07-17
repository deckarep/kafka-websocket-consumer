package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", wsHandler)

	fmt.Println("Listening on port: 8989")
	if err := http.ListenAndServe(":8989", nil); err != nil {
		log.Fatal(err)
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println("Socket upgraded!")
	go writer(conn)
}

func writer(conn *websocket.Conn) {
	defer func() {
		conn.Close()
		fmt.Println("conn Closed()")
	}()

	count := 0

	for {
		msg := fmt.Sprintf("How are you doing? %d", count)
		if err := conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
			log.Println("An error occured writing to the websocket.")
		}
		count++
		time.Sleep(time.Second * 1)
	}
}
func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", 404)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	fmt.Println("Serving home.")
	fmt.Fprintln(w, homeHTML)
}

const homeHTML = `<!DOCTYPE html>
<html>
   <head>
      <script type="text/javascript">
         function WebSocketTest() {
            if ("WebSocket" in window) {
               console.log("WebSocket is supported by your Browser!");
               
               // Let us open a web socket
               var ws = new WebSocket("ws://localhost:8989/ws");
				
               ws.onopen = function() {
                  // Web Socket is connected, send data using send()
                  ws.send("Message to send");
                  console.log("Message is sent...");
               };
				
               ws.onmessage = function (evt) { 
                  var received_msg = evt.data;
				  console.log(evt);
				  console.log(evt.data);
               };
				
               ws.onclose = function() { 
                  // websocket is closed.
                  console.log("Connection is closed..."); 
               };
            }
            else {
               // The browser doesn't support WebSocket
               console.log("WebSocket NOT supported by your Browser!");
            }
         }
      </script>
</head>
   <body>
      <div id="sse">
        <a href="javascript:WebSocketTest()">Start!</a>
      </div>
   </body>
</html>
`

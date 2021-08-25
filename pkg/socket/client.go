package socket

import (
	"github.com/gorilla/websocket"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	log "github.com/sirupsen/logrus"
	"net/http"
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Client struct {
	conn *websocket.Conn
	send chan []byte
}

func (c *Client) writePump() {

}

// serveWs handles websocket requests from the peer.
func serveWs(job *types.CommandJob, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{conn: conn, send: make(chan []byte, 256)}

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
}

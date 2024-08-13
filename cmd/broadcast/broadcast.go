package broadcast

import (
	"log"
	"net/http"

	"github.com/Cirqach/GoStream/cmd/queueController"
	"github.com/Cirqach/GoStream/internal/database"
)

// BroadcastEngine struct    allow access to control websockets connection and vidoe lifetime update
type BroadcastEngine struct {
	Hub             *Hub
	queueController *queuecontroller.QueueController
}

// HandleWebsocket method    handler for htmx websocket
func (b *BroadcastEngine) HandleWebsocket(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling websocket connection")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Connection established")
	client := &Client{hub: b.Hub, conn: conn}
	client.hub.register <- client
}

// NewBroadcastEngine function    create new BroadcastEngine object
func NewBroadcastEngine() *BroadcastEngine {
	log.Println("Creating new broadcast struct")
	return &BroadcastEngine{
		Hub:             NewHub(),
		queueController: queuecontroller.NewQueueController(),
	}
}

// UpdateVideo method    livetime video update
func (b *BroadcastEngine) UpdateVideo(c chan database.Video) {
	video := <-c
	b.Hub.Stream <- []byte(`
	<div hx-swap-oob="innerHTML:#video-div">
	<video id="video" controls autoplay></video>
	</div>
	<div hx-swap-oob="innerHTML:#videoJS-div">
	<script>
	    if(Hls.isSupported()) {
	    var video = document.getElementById('video');
	    var hls = new Hls();
	    hls.loadSource('http://localhost:8080/` + video.Path + `/index.m3u8');
	    hls.attachMedia(video);
	    hls.on(Hls.Events.MANIFEST_PARSED,function
	      video.play();
	    });
	    }
	</script>
	</div>
		`)
}

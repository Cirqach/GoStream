package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Cirqach/GoStream/cmd/server/api/handler"
	"github.com/Cirqach/GoStream/cmd/server/broadcast"
	videoprocessor "github.com/Cirqach/GoStream/cmd/server/broadcast/videoProcessor"
	"github.com/gorilla/mux"
	// "github.com/Cirqach/GoStream/cmd/server/broadcast/videoProcessor"
)


type server struct {
	router *mux.Router
	port string
}

func NewServer(addr string) *server{
	log.Println("Creating new server")
	return &server{
		router: mux.NewRouter(),
		port: addr,
	}
}


func (s *server)StartServer() {
	vp := videoprocessor.NewVideoProcessor()
	vp.FFmpegEngine.Process("./video/unprocessed/test1.mp4", "test1")
	log.Println("Starting server")
	b := broadcast.NewBroadcast()
	go b.Hub.Run()
	log.Println("Handling websocket")
	s.router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
b.HandleWebsocket(w, r)
	})
	log.Println("Serving static files")
	s.router.PathPrefix("/video/processed/").Handler(http.StripPrefix("/video/processed/", http.FileServer(http.Dir("./video/processed/"))))
	s.router.PathPrefix("/web/static/js/").Handler(http.StripPrefix("/web/static/js/", http.FileServer(http.Dir("./web/static/js/"))))
	s.router.PathPrefix("/web/static/css/").Handler(http.StripPrefix("/web/static/css/", http.FileServer(http.Dir("./web/static/css/"))))
	log.Println("Serving routes files")
	s.router.HandleFunc("/", handler.RootHandler)
	s.router.HandleFunc("/watch", handler.WatchHandler)
	s.router.HandleFunc("/book", handler.BookatimeHandler)
	
go func(b *broadcast.Broadcast){	
	
	time.Sleep(10 * time.Second)
	vp.FFmpegEngine.Process("./video/unprocessed/test2.mp4", "test2")
	b.Hub.Stream <- []byte(`
	<div hx-swap-oob="innerHTML:#video-div">
        <video id="video" controls autoplay></video>
	</div>
	<div hx-swap-oob="innerHTML:#videoJS-div">
<script>
    if(Hls.isSupported()) {
    var video = document.getElementById('video');
    var hls = new Hls();
    hls.loadSource('http://localhost:8080/video/processed/` + "test2" + `/index.m3u8');
    hls.attachMedia(video);
    hls.on(Hls.Events.MANIFEST_PARSED,function() {
      video.play();
    });
    }
    </script>
	</div>
	`)

}(b)

	fmt.Printf("\nListen on %s\n", s.port)
	go log.Fatal(http.ListenAndServe(s.port, s.router))
	
}


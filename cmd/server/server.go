package server

import (
	"fmt"
	"log"
	"time"

	// "net/http"
	//
	// "github.com/Cirqach/GoStream/cmd/server/api/handler"
	"github.com/Cirqach/GoStream/cmd/server/broadcast"
	videoprocessor "github.com/Cirqach/GoStream/cmd/server/broadcast/videoProcessor"
	"github.com/Cirqach/GoStream/internal/database"
	"github.com/gorilla/mux"
	// "github.com/Cirqach/GoStream/cmd/server/broadcast/videoProcessor"
)

type server struct {
	router             *mux.Router
	videoProcessor     *videoprocessor.VideoProcessor
	broadcastEngine    *broadcast.BroadcastEngine
	databaseController *database.DatabaseController
	port               string
}

func NewServer(addr string) *server {
	log.Println("Creating new server")
	return &server{
		router:             mux.NewRouter(),
		videoProcessor:     videoprocessor.NewVideoProcessor(),
		broadcastEngine:    broadcast.NewBroadcastEngine(),
		databaseController: database.NewDatabaseController(),
		port:               addr,
	}
}

func (s *server) StartServer() {

	s.databaseController.MakeConnection()
	data, err := s.databaseController.GetSoonerVideo()
	if err != nil {
		log.Printf("error due receiving data: %s", err)
	}
	fmt.Printf("recieved data: %s", data)

	// go s.videoProcessor.Process("./video/unprocessed/1.mp4", "1")
	// go s.videoProcessor.Process("./video/unprocessed/2.mp4", "2")
	// go s.videoProcessor.Process("./video/unprocessed/3.mp4", "3")
	// go s.videoProcessor.Process("./video/unprocessed/4.mp4", "4")
	// go s.videoProcessor.Process("./video/unprocessed/5.mp4", "5")
	// log.Println("Starting server")
	// go s.broadcastEngine.Hub.Run()
	// log.Println("Handling websocket")
	// s.router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
	// 	s.broadcastEngine.HandleWebsocket(w, r)
	// })
	// log.Println("Serving static files")
	// s.router.PathPrefix("/video/processed/").
	// 	Handler(http.StripPrefix("/video/processed/",
	// 		http.FileServer(http.Dir("./video/processed/"))))
	// s.router.PathPrefix("/web/static/js/").
	// 	Handler(http.StripPrefix("/web/static/js/",
	// 		http.FileServer(http.Dir("./web/static/js/"))))
	// s.router.PathPrefix("/web/static/css/").
	// 	Handler(http.StripPrefix("/web/static/css/",
	// 		http.FileServer(http.Dir("./web/static/css/"))))
	// log.Println("Serving routes files")
	// s.router.HandleFunc("/", handler.RootHandler)
	// s.router.HandleFunc("/watch", handler.WatchHandler)
	// s.router.HandleFunc("/book", handler.BookatimeHandler)
	//
	// go func(b *broadcast.BroadcastEngine) {
	//
	// 		time.Sleep(10 * time.Second)
	// 		s.videoProcessor.Process("./video/unprocessed/test2.mp4", "test2")
	// 		b.Hub.Stream <- []byte(`
	// 	<div hx-swap-oob="innerHTML:#video-div">
	//         <video id="video" controls autoplay></video>
	// 	</div>
	// 	<div hx-swap-oob="innerHTML:#videoJS-div">
	// <script>
	//     if(Hls.isSupported()) {
	//     var video = document.getElementById('video');
	//     var hls = new Hls();
	//     hls.loadSource('http://localhost:8080/video/processed/` + "test2" + `/index.m3u8');
	//     hls.attachMedia(video);
	//     hls.on(Hls.Events.MANIFEST_PARSED,function
	//       video.play();
	//     });
	//     }
	//     </script>
	// 	</div>
	// 	`)
	// 	}(s.broadcastEngine)
	// fmt.Printf("\nListen on %s\n", s.port)
	// go log.Fatal(http.ListenAndServe(s.port, s.router))
}

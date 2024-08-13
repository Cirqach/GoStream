package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Cirqach/GoStream/cmd/api/handler"
	"github.com/Cirqach/GoStream/cmd/broadcast"
	"github.com/Cirqach/GoStream/cmd/logger"
	videoprocessor "github.com/Cirqach/GoStream/cmd/videoProcessor"
	"github.com/Cirqach/GoStream/internal/database"
	"github.com/gorilla/mux"
)

type server struct {
	router             *mux.Router
	videoProcessor     *videoprocessor.VideoProcessor
	broadcastEngine    *broadcast.BroadcastEngine
	databaseController *database.DatabaseController
	protocol           string
	ip                 string
	port               string
}

func NewServer(protocol, ip, port string) *server {
	log.Println("Creating new server")
	return &server{
		router:             mux.NewRouter(),
		videoProcessor:     videoprocessor.NewVideoProcessor(),
		broadcastEngine:    broadcast.NewBroadcastEngine(),
		databaseController: database.NewDatabaseController(),
		protocol:           protocol,
		ip:                 ip,
		port:               port,
	}
}

func (s *server) StartServer() {

	logger.LogMessage("StartServer", "Starting server")
	go s.broadcastEngine.Hub.Run()

	logger.LogMessage("StartServer", "Handling websocket")
	s.router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		s.broadcastEngine.HandleWebsocket(w, r)
	})

	logger.LogMessage("StartServer", "Serving static files")
	s.router.PathPrefix("/video/processed/").
		Handler(http.StripPrefix("/video/processed/",
			http.FileServer(http.Dir("./video/processed/"))))
	s.router.PathPrefix("/web/static/js/").
		Handler(http.StripPrefix("/web/static/js/",
			http.FileServer(http.Dir("./web/static/js/"))))
	s.router.PathPrefix("/web/static/css/").
		Handler(http.StripPrefix("/web/static/css/",
			http.FileServer(http.Dir("./web/static/css/"))))

	logger.LogMessage("StartServer", "Serving routes")
	s.router.HandleFunc("/", handler.RootHandler(s.protocol+s.ip+s.port))
	s.router.HandleFunc("/watch", handler.WatchHandler(s.protocol+s.ip+s.port))
	s.router.HandleFunc("/book", handler.BookatimeHandler(s.protocol+s.ip+s.port))

	// go func(b *broadcast.BroadcastEngine) {
	// 	time.Sleep(10 * time.Second)
	// 	s.videoProcessor.Process("./video/unprocessed/test2.mp4", "test2")
	// 	b.Hub.Stream <- []byte(`
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
	// }(s.broadcastEngine)
	fmt.Printf("\nListen on %s\n", s.port)
	go log.Fatal(http.ListenAndServe(s.port, s.router))
}

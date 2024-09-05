package broadcast

import (
	"github.com/Cirqach/GoStream/cmd/logger"
	"github.com/Cirqach/GoStream/internal/database"
)

// Engine struct    allow access to control websockets connection and vidoe lifetime update
type Engine struct {
	Chan chan database.Video
	Hub  *Hub
}

// NewEngine function    create new BroadcastEngine object
func NewEngine() *Engine {
	logger.LogMessage(logger.GetFuncName(0), "Creating new broadcast struct")
	return &Engine{
		Hub: NewHub(),
	}
}

// TODO: need to create time changing for videos
// UpdateVideo method    livetime video update
func (b *Engine) UpdateVideo(v database.Video) {
	logger.LogMessage(logger.GetFuncName(0), "Start updating video")
	logger.LogMessage(logger.GetFuncName(0), "Updating video: "+v.Name)
	b.Hub.Stream <- []byte(`
	<div hx-swap-oob="innerHTML:#video-div">
	<video id="video" controls autoplay></video>
	</div>
	<div hx-swap-oob="innerHTML:#videoJS-div">
	<script>
	    if(Hls.isSupported()) {
	    var video = document.getElementById('video');
	    var hls = new Hls();
	    hls.loadSource('http://localhost:8080/processed/` + v.Name + `/index.m3u8');
	    hls.attachMedia(video);
	    hls.on(Hls.Events.MANIFEST_PARSED,function
	      video.play();
	    });
	    }
	</script>
	</div>
		`)
}

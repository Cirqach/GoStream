package broadcast

import (
	"github.com/Cirqach/GoStream/cmd/logger"
	"github.com/Cirqach/GoStream/internal/database"
)

// BroadcastEngine struct    allow access to control websockets connection and vidoe lifetime update
type BroadcastEngine struct {
	Hub *Hub
}

// NewBroadcastEngine function    create new BroadcastEngine object
func NewBroadcastEngine() *BroadcastEngine {
	logger.LogMessage(logger.GetFuncName(0), "Creating new broadcast struct")
	return &BroadcastEngine{
		Hub: NewHub(),
	}
}

// TODO: need to create time changing for videos
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

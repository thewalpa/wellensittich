package interfaces

import (
	"fmt"
	"path"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/pion/webrtc/v3/pkg/media/oggwriter"
	"github.com/thewalpa/wellensittich/util"
)

type OGGFileSoundReceiver struct {
	Path      string
	ssrcFiles map[uint32]*oggwriter.OggWriter
	mu        sync.Mutex
}

func NewOGGFileSoundReceiver(path string) (*OGGFileSoundReceiver, error) {
	return &OGGFileSoundReceiver{Path: path, ssrcFiles: make(map[uint32]*oggwriter.OggWriter)}, nil
}

func (osr *OGGFileSoundReceiver) Before() {
}

func (osr *OGGFileSoundReceiver) After() {
	for _, file := range osr.ssrcFiles {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}
	osr.ssrcFiles = make(map[uint32]*oggwriter.OggWriter)
}

func (osr *OGGFileSoundReceiver) HandleOpusPacket(p discordgo.Packet, done chan struct{}) {
	osr.mu.Lock()
	defer osr.mu.Unlock()
	file, ok := osr.ssrcFiles[p.SSRC]
	if !ok {
		newFile, err := oggwriter.New(path.Join(osr.Path, fmt.Sprintf("%d.ogg", p.SSRC)), 48000, 2)
		if err != nil {
			fmt.Println(err)
			return
		}
		osr.ssrcFiles[p.SSRC] = newFile
		file = newFile
	}

	rtp := util.CreatePionRTPPacket(&p)
	err := file.WriteRTP(rtp)
	if err != nil {
		fmt.Println(err)
	}
}

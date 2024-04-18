package discordws

import (
	"fmt"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

type SoundReceiver interface {
	HandleOpusPacket(discordgo.Packet, chan struct{})
	Before()
	After()
}

type GenericSoundReceiver struct {
	GenericHandler func(discordgo.Packet, chan struct{})
	GenericBefore  func()
	GenericAfter   func()
}

func (gsr *GenericSoundReceiver) HandleOpusPacket(p discordgo.Packet, done chan struct{}) {
	if gsr.GenericHandler != nil {
		gsr.GenericHandler(p, done)
	}
}

func (gsr *GenericSoundReceiver) Before() {
	if gsr.GenericBefore != nil {
		gsr.GenericBefore()
	}
}

func (gsr *GenericSoundReceiver) After() {
	if gsr.GenericAfter != nil {
		gsr.GenericAfter()
	}
}

type VoiceReceiver struct {
	mu             sync.Mutex
	done           chan struct{}
	listenPassive  bool
	subscribeEvent chan struct{}
	subscribers    map[string]SoundReceiver
}

func NewVoiceReceiver() *VoiceReceiver {
	return &VoiceReceiver{
		done:           make(chan struct{}),
		subscribers:    make(map[string]SoundReceiver),
		subscribeEvent: make(chan struct{}, 1),
	}
}

func (vr *VoiceReceiver) IsListenPassive() bool {
	vr.mu.Lock()
	defer vr.mu.Unlock()

	return vr.listenPassive
}

func (vr *VoiceReceiver) notifySubscribersChange() {
	select {
	case vr.subscribeEvent <- struct{}{}:
	default:
	}
}

func (vr *VoiceReceiver) Subscribe(featureName string, handler SoundReceiver) {
	vr.mu.Lock()
	defer vr.mu.Unlock()

	go handler.Before()
	vr.subscribers[featureName] = handler
	if len(vr.subscribers) == 1 {
		vr.notifySubscribersChange()
	}
}

func (vr *VoiceReceiver) SubscribeOnce(featureName string, handler SoundReceiver) {
	vr.mu.Lock()
	defer vr.mu.Unlock()

	go handler.Before()
	vr.subscribers[featureName] = &GenericSoundReceiver{
		GenericHandler: func(p discordgo.Packet, done chan struct{}) {
			handler.HandleOpusPacket(p, done)
			vr.Unsubscribe(featureName)
		},
		GenericBefore: handler.Before,
		GenericAfter:  handler.After,
	}
	if len(vr.subscribers) == 1 {
		vr.notifySubscribersChange()
	}
}

func (vr *VoiceReceiver) SubscribeWithTimer(featureName string, handler SoundReceiver, duration time.Duration) {
	vr.mu.Lock()
	defer vr.mu.Unlock()

	go handler.Before()
	vr.subscribers[featureName] = handler
	if len(vr.subscribers) == 1 {
		vr.notifySubscribersChange()
	}
	go func() {
		select {
		case <-time.After(duration):
			vr.Unsubscribe(featureName)
		case <-vr.done:
			return
		}
	}()
}

func (vr *VoiceReceiver) Unsubscribe(featureName string) {
	vr.mu.Lock()
	defer vr.mu.Unlock()

	// could be a potentially long activity, so do it concurrently
	sub, ok := vr.subscribers[featureName]
	if !ok {
		// was already Unsubscribed
		return
	}

	go sub.After()
	delete(vr.subscribers, featureName)
	if len(vr.subscribers) == 0 {
		vr.notifySubscribersChange()
	}
}

func (vr *VoiceReceiver) notifySubscribers(p *discordgo.Packet) {
	vr.mu.Lock()
	defer vr.mu.Unlock()
	for _, handler := range vr.subscribers {
		// send a copy of the opus packet to all subscribers
		// for proper synchronization these have to run after each other
		// HandleOpusPacket should return early therefore
		// maybe add a waitgroup, requires testing what would be more performant
		handler.HandleOpusPacket(*p, vr.done)
	}
}

func (vr *VoiceReceiver) ToggleListen(wsvc *WellensittichVoiceConnection) bool {
	vr.mu.Lock()
	defer vr.mu.Unlock()

	if vr.listenPassive {
		vr.listenPassive = false
		go vr.Stop()
		return false
	} else {
		vr.listenPassive = true
		vr.done = make(chan struct{})
		go vr.Listen(wsvc)
		return true
	}
}

func (vr *VoiceReceiver) Listen(wsvc *WellensittichVoiceConnection) {
	fmt.Println("Starting listening.")
	defer func() {
		vr.mu.Lock()
		vr.listenPassive = false
		vr.mu.Unlock()
	}()
	for {
		select {
		case <-vr.subscribeEvent:
			fmt.Println("First subscriber added.")
		innerloop:
			for {
				select {
				case opusPacket, open := <-wsvc.OpusRecv:
					if !open {
						fmt.Println("OpusRecv channel closed. Stopping listening.")
						return
					}
					vr.notifySubscribers(opusPacket)
				case <-vr.done:
					fmt.Println("Received done signal. Stopping listening.")
					return
				case <-vr.subscribeEvent:
					vr.mu.Lock()
					numSubscribers := len(vr.subscribers)
					vr.mu.Unlock()
					if numSubscribers == 0 {
						fmt.Println("Last subscriber removed.")
						break innerloop
					}
				}
			}
		case <-vr.done:
			fmt.Println("Received done signal. Stopping listening.")
			return
		}
	}
}

func (vr *VoiceReceiver) Stop() {
	vr.mu.Lock()
	defer vr.mu.Unlock()
	if !vr.listenPassive {
		return
	}

	for featureName, subscriber := range vr.subscribers {
		go subscriber.After()
		delete(vr.subscribers, featureName)
	}

	if vr.done != nil {
		close(vr.done)
	}
	vr.listenPassive = false
}

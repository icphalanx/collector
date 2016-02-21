package collector

import (
	pb "github.com/icphalanx/rpc"
	"sync"
)

type Hub struct {
	roomLock sync.RWMutex
	rooms    map[string][]chan pb.LogLine
}

func (h *Hub) Register(ch chan pb.LogLine, rooms []string) {
	h.roomLock.Lock()
	defer h.roomLock.Unlock()

	for _, room := range rooms {
		if _, ok := h.rooms[room]; !ok {
			h.rooms[room] = make([]chan pb.LogLine, 0)
		}

		h.rooms[room] = append(h.rooms[room], ch)
	}
}

func (h *Hub) Send(ll pb.LogLine, rooms []string) {
	h.roomLock.Lock()
	defer h.roomLock.Unlock()

	for _, room := range rooms {
		if roomChans, ok := h.rooms[room]; ok {
			okChans := make([]chan pb.LogLine, 0)
			for _, ch := range roomChans {
				select {
				case ch <- ll:
					okChans = append(okChans, ch)
				default:
				}
			}
			h.rooms[room] = okChans
		}
	}
}

var (
	MainHub *Hub = &Hub{rooms: make(map[string][]chan pb.LogLine)}
)

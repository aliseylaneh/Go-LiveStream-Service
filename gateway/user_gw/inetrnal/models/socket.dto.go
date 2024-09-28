package models

import (
	"context"
	"encoding/json"
	"os"
	"safir/libs/idgen"
	"strings"
	"sync"
	"time"
	"vpeer_usergw/inetrnal/global"
	pb_room "vpeer_usergw/proto/api/room"

	"github.com/gofiber/contrib/websocket"
	"github.com/pion/rtcp"
	"github.com/pion/webrtc/v3"
	"github.com/pion/webrtc/v3/pkg/media/ivfwriter"
	"github.com/pion/webrtc/v3/pkg/media/oggwriter"
)

var (
	Rooms       = make(map[string]*Room)
	ExpiryRooms = make(map[string]*time.Time)
	Recorded    = make(map[string]*OpenFile)
)

func GetRoomCopy() map[string]*Room {
	rooms := Rooms
	return rooms

}

type User struct {
	DisplayName  string  `json:"display_name"`
	UserId       string  `json:"user_id"`
	ConnectionId string  `json:"connection_id"`
	Recording    string  `json:"recorded_file"`
	IvfFileName  *string `json:"video_filename"`
	OggFileName  *string `json:"audio_filename"`
}

type Users struct {
	Users []User `json:"users"`
}

type WebRTCMessage struct {
	Event string `json:"event"`
	Data  string `json:"data"`
	Name  string `json:"name"`
}

type UserDataMessage struct {
	DisplayName  string `json:"DisplayName"`
	ConnectionId string `json:"connection_id"`
	UserId       string `json:"user_id"`
}

func NewUserDataMessage(displayName, connectionId string, UserId string) UserDataMessage {
	return UserDataMessage{
		DisplayName:  displayName,
		ConnectionId: connectionId,
		UserId:       UserId,
	}
}

type PeerConnectionState struct {
	UserId         string  `json:"user_id"`
	User           string  `json:"user"`
	DisplayName    string  `json:"display_name"`
	ConnectionId   string  `json:"connection_id"`
	RecordingState string  `json:"recording_state"`
	IvfFileName    *string `json:"ivf_file_name"`
	OggFileName    *string `json:"ogg_file_name"`
	PeerConnection *webrtc.PeerConnection
	Websocket      *ThreadSafeWriter
}

type MeetRecord struct {
	FileId string
	File   *os.File
}

type Peers struct {
	ListLock    sync.RWMutex
	Connections []PeerConnectionState
	TrackLocals map[string]TrackLocalKey
}

type TrackLocalKey struct {
	// UserID int32
	// Type       string
	TrackLocal *webrtc.TrackLocalStaticRTP
}

type Room struct {
	OwnerId string
	Peers   *Peers
	// RoomExpiry *int64
	UserLength int32
	Filled     bool
	Polls      *Polls
	CreatedAt  time.Time
}

type Polls struct {
	UserPolls map[string]bool
}

type OggFileStruct struct {
	OggFile     *oggwriter.OggWriter
	OggFileName string
}

type IvfFileStruct struct {
	IvfFile     *ivfwriter.IVFWriter
	IvfFileName string
}

type MkvFileStruct struct {
	MkvFile     *os.File
	MkvFileName string
}

// Helper to make Gorilla Websockets threadsafe
type ThreadSafeWriter struct {
	*websocket.Conn
	sync.Mutex
}

type OpenFile struct {
	OggFile     *oggwriter.OggWriter
	IvfFile     *ivfwriter.IVFWriter
	OggFileName string
	IvfFileName string
	UserId      string
	RoomId      string
	MkvFileName string
	WebSocket   *ThreadSafeWriter
}

func CreateNewPeer(res global.TokenInformation, displayName string, threadsafe *ThreadSafeWriter, peerConnection *webrtc.PeerConnection) (*PeerConnectionState, *WebSocketMessage) {
	id, err := idgen.NextAlphanumericString(20)
	if err != nil {
		return nil, &WebSocketMessage{Event: "error", State: "connection", Data: "failed to create a connection id"}
	}
	global.ROOM_SERVER_CLIENT.AddRoomLog(context.Background(), &pb_room.AddRoomLog{
		RoomId:    res.RoomId,
		UserId:    res.UserId,
		UserEvent: "joined",
	})
	return &PeerConnectionState{
		Websocket:      threadsafe,
		PeerConnection: peerConnection,
		UserId:         res.UserId,
		User:           res.UserId,
		DisplayName:    displayName,
		ConnectionId:   id,
		IvfFileName:    nil,
		OggFileName:    nil,
		RecordingState: "not_recorded",
	}, nil

}

func (t *ThreadSafeWriter) WriteJSON(v interface{}) error {
	t.Lock()
	defer t.Unlock()

	return t.Conn.WriteJSON(v)
}

func (t *ThreadSafeWriter) SafeClose() error {
	t.Lock()
	defer t.Unlock()

	return t.Conn.Close()
}

type WebSocketMessage struct {
	Event string      `json:"event"`
	State string      `json:"state"`
	Data  interface{} `json:"message"`
}

func (r *Room) BroadcastNotify(state string, data interface{}) {
	for _, connection := range r.Peers.Connections {
		connection.Websocket.WriteJSON(&WebSocketMessage{
			Event: "notification",
			State: state,
			Data:  data,
		})
	}
}

func (r *Room) SingleNotify(state string, message interface{}, conn *ThreadSafeWriter) {
	conn.WriteJSON(&WebSocketMessage{
		Event: "notification",
		State: state,
		Data:  message,
	})
}

func DispatchKeyFrames() {
	for range time.NewTicker(time.Second * 3).C {
		for _, room := range Rooms {
			room.Peers.DispatchKeyFrame()
		}
	}
}

func (p *Peers) AddTrack(t *webrtc.TrackRemote) TrackLocalKey {
	p.ListLock.Lock()
	defer func() {
		p.ListLock.Unlock()
		p.SignalPeerConnections()
	}()

	// Create a new TrackLocal with the same codec as our incoming
	trackLocal, err := webrtc.NewTrackLocalStaticRTP(t.Codec().RTPCodecCapability, t.ID(), t.StreamID())
	if err != nil {
		panic(err)
	}
	trackLocalState := TrackLocalKey{
		// UserID:     user_id,
		TrackLocal: trackLocal,
		// Type:       track_type,
	}
	p.TrackLocals[t.ID()] = trackLocalState
	return trackLocalState
}

func (p *Peers) RemoveTrack(t *webrtc.TrackLocalStaticRTP) {
	p.ListLock.Lock()
	defer func() {
		p.ListLock.Unlock()
		p.SignalPeerConnections()
	}()

	// delete(p.TrackLocals, t.ID())
	delete(p.TrackLocals, t.ID())
}

func (p *Peers) SignalPeerConnections() {
	p.ListLock.Lock()
	defer func() {
		p.ListLock.Unlock()
		p.DispatchKeyFrame()
	}()

	attemptSync := func() (tryAgain bool) {
		for i := range p.Connections {
			if p.Connections[i].PeerConnection.ConnectionState() == webrtc.PeerConnectionStateClosed {
				p.Connections = append(p.Connections[:i], p.Connections[i+1:]...)
				return true
			}
			existingSenders := map[string]bool{}
			for _, sender := range p.Connections[i].PeerConnection.GetSenders() {
				if sender.Track() == nil {
					continue
				}

				existingSenders[sender.Track().ID()] = true

				if _, ok := p.TrackLocals[sender.Track().ID()]; !ok {
					if err := p.Connections[i].PeerConnection.RemoveTrack(sender); err != nil {
						return true
					}
				}
			}

			for _, receiver := range p.Connections[i].PeerConnection.GetReceivers() {
				if receiver.Track() == nil {
					continue
				}

				existingSenders[receiver.Track().ID()] = true
			}
			for trackID := range p.TrackLocals {
				if _, ok := existingSenders[trackID]; !ok {
					if _, err := p.Connections[i].PeerConnection.AddTrack(p.TrackLocals[trackID].TrackLocal); err != nil {
						return true
					}

				}
				// break
			}
			offer, err := p.Connections[i].PeerConnection.CreateOffer(nil)
			if err != nil {
				return true
			}

			if err = p.Connections[i].PeerConnection.SetLocalDescription(offer); err != nil {
				if strings.Contains(err.Error(), "have-local-offer") {
					if p.Connections[i].PeerConnection.ConnectionState() == webrtc.PeerConnectionStateNew {
						p.Connections[i].PeerConnection.Close()
						p.Connections[i].Websocket.Close()
						return false
					}
				}

				return true
			}
			offerString, err := json.Marshal(offer)
			if err != nil {
				return true
			}

			if err = p.Connections[i].Websocket.WriteJSON(&WebRTCMessage{
				Event: "offer",
				Data:  string(offerString),
			}); err != nil {
				return true
			}
		}
		return
	}
	for syncAttempt := 0; ; syncAttempt++ {
		if syncAttempt == 25 {
			go func() {
				time.Sleep(time.Second * 3)
				p.SignalPeerConnections()
			}()
			return
		}

		if !attemptSync() {
			break
		}
	}
}
func (p *Peers) DispatchKeyFrame() {
	p.ListLock.Lock()
	defer p.ListLock.Unlock()

	for i := range p.Connections {
		for _, receiver := range p.Connections[i].PeerConnection.GetReceivers() {
			if receiver.Track() == nil {
				continue
			}

			_ = p.Connections[i].PeerConnection.WriteRTCP([]rtcp.Packet{
				&rtcp.PictureLossIndication{
					MediaSSRC: uint32(receiver.Track().SSRC()),
				},
			})
		}
	}
}

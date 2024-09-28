package helpers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
	"vpeer_usergw/inetrnal/global"
	"vpeer_usergw/inetrnal/models"
	"vpeer_usergw/inetrnal/services"
	"vpeer_usergw/inetrnal/types"
	pb_room "vpeer_usergw/proto/api/room"

	"github.com/gofiber/contrib/websocket"
	"github.com/pion/interceptor"
	"github.com/pion/interceptor/pkg/intervalpli"
	"github.com/pion/webrtc/v3"
	"github.com/pion/webrtc/v3/pkg/media"
)

func copyString(original string) string {
	// Convert the original string to a byte slice and back to a string.
	copiedString := string([]byte(original))
	return copiedString
}
func VerifyToken(conn *websocket.Conn) (*global.TokenInformation, string, *models.WebSocketMessage) {
	timeoutDuration := 10 * time.Second
	conn.SetReadDeadline(time.Now().Add(timeoutDuration))

	message := &models.WebRTCMessage{}
	_, raw, err := conn.ReadMessage()
	if err != nil {
		return nil, "", &models.WebSocketMessage{
			Event: "error",
			State: "verification",
			Data:  "خطا در درخواست احراز هویت",
		}
	} else if err := json.Unmarshal(raw, &message); err != nil {
		return nil, "", &models.WebSocketMessage{
			Event: "error",
			State: "verification",
			Data:  "درخواست نامعتبر است",
		}
	}

	mres := global.TOKENS[message.Data]
	if mres == nil {
		return nil, "", &models.WebSocketMessage{
			Event: "error",
			State: "verification",
			Data:  "توکن نامعتبر است",
		}
	}
	if mres.ExpireAt.Before(time.Now()) {
		delete(global.TOKENS, message.Data)
		return nil, "", &models.WebSocketMessage{
			Event: "error",
			State: "verification",
			Data:  "توکن منقضی شده است",
		}
	}
	delete(global.TOKENS, message.Data)
	newmemUserId := copyString(mres.UserId)
	mres.UserId = newmemUserId
	return mres, message.Name, nil
}

func GetRoomOrInit(roomId string) (*models.Room, *models.WebSocketMessage) {
	var room models.Room
	room_copy := models.GetRoomCopy()[roomId]

	if room_copy != nil {
		room = *room_copy
	} else {
		p := &models.Peers{}
		polls := &models.Polls{
			UserPolls: make(map[string]bool),
		}
		p.TrackLocals = make(map[string]models.TrackLocalKey)
		res, err := global.ROOM_SERVER_CLIENT.GetRoomByRoomid(context.TODO(), &pb_room.GetRoomByRoomId{
			RoomId: roomId,
		})
		if err != nil {
			return nil, &models.WebSocketMessage{
				Event: "error",
				State: "room",
				Data:  "ساخت جلسه با خطا مواجه شد",
			}
		}
		expendTime := time.Now().Add(time.Minute * 10)
		room = models.Room{
			OwnerId:    res.UserId,
			Peers:      p,
			UserLength: res.UsersLength,
			Filled:     false,
			Polls:      polls,
			CreatedAt:  time.Now(),
		}

		models.Rooms[roomId] = &room
		models.ExpiryRooms[roomId] = &expendTime

		// if res.RoomExpiry != nil {
		// expendTime := time.Now().Add(time.Second * 10)
		// room.RoomExpiry = &expendTime
		// room.IdleExpiry = res.RoomExpiry
		// expireTime := time.Unix(*res.RoomExpiry, 0)
		// models.ExpiryRooms[roomId] = &expendTime
		// }
	}
	return &room, nil
}
func SaveToDisk(i media.Writer, track *webrtc.TrackRemote, trackLocal models.TrackLocalKey) {
	defer func() {
		if err := i.Close(); err != nil {
			fmt.Println("Error closing media writer:", err)
		}
	}()

	for {
		rtpPacket, _, err := track.ReadRTP()
		if err != nil {
			if err == io.EOF {
				// EOF indicates that the track has no more data to read
				fmt.Println("Track has reached EOF. Closing SaveToDisk routine.")
				return
			}
			fmt.Println("Error reading RTP:", err)
			return
		}

		err = trackLocal.TrackLocal.WriteRTP(rtpPacket)
		if err != nil {
			println("Error writing RTP: ")
			return
		}
		// Check if the file is opened before writing the RTP packet
		if err := i.WriteRTP(rtpPacket); err != nil {
			if !strings.Contains(err.Error(), "file not opened") {
				fmt.Println("Error writing RTP to media writer:", err)
			}
			return
		}
	}
}

func OnICECandidate(peerConnection *webrtc.PeerConnection, threadSafeWriter *models.ThreadSafeWriter) {
	peerConnection.OnICECandidate(func(i *webrtc.ICECandidate) {
		if i == nil {
			return
		}

		candidateString, err := json.Marshal(i.ToJSON())
		if err != nil {
			log.Println(err)
			return
		}

		if writeErr := threadSafeWriter.WriteJSON(&models.WebRTCMessage{
			Event: "candidate",
			Data:  string(candidateString),
		}); writeErr != nil {
			log.Println(writeErr)
		}
	})
}
func RandStringRunes(n int) string {
	// var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func CreateMkvFile(minioService services.MinioService, fileName string) (*models.MkvFileStruct, *types.Error) {
	videoFile, cerr := minioService.CreateMkvFile(fileName + ".mkv")
	if cerr != nil {
		return nil, cerr
	}
	mkvFile := &models.MkvFileStruct{
		MkvFile:     videoFile,
		MkvFileName: fileName + ".mkv",
	}

	return mkvFile, nil
}

func CreateIvfFile(minioService services.MinioService, fileName string) (*models.IvfFileStruct, *types.Error) {
	videoFile, cerr := minioService.CreateIvfFile(fileName + ".ivf")
	if cerr != nil {
		return nil, cerr
	}
	ivfFile := &models.IvfFileStruct{
		IvfFile:     videoFile,
		IvfFileName: fileName + ".ivf",
	}

	return ivfFile, nil
}

func CreateOggFile(minioService services.MinioService, fileName string) (*models.OggFileStruct, *types.Error) {
	audioFile, cerr := minioService.CreateOggFile(fileName + ".ogg")
	if cerr != nil {
		return nil, cerr
	}
	oggFile := &models.OggFileStruct{
		OggFile:     audioFile,
		OggFileName: fileName + ".ogg",
	}
	return oggFile, nil
}

func InitializeRecordingFiles(minioService services.MinioService) (*models.IvfFileStruct, *models.OggFileStruct, *types.Error) {
	fileName := RandStringRunes(30)
	videoFile, cerr := minioService.CreateIvfFile(fileName + ".ivf")
	if cerr != nil {
		return nil, nil, cerr
	}
	ivfFile := &models.IvfFileStruct{
		IvfFile:     videoFile,
		IvfFileName: fileName + ".ivf",
	}

	audioFile, cerr := minioService.CreateOggFile(fileName + ".ogg")
	if cerr != nil {
		return nil, nil, cerr
	}
	oggFile := &models.OggFileStruct{
		OggFile:     audioFile,
		OggFileName: fileName + ".ogg",
	}

	return ivfFile, oggFile, nil

}

// CreatePeerConnection initializes a new PeerConnection with specified configurations
func CreatePeerConnection() (*webrtc.PeerConnection, error) {
	mediaEngine := &webrtc.MediaEngine{}
	// RegisterDefaultCodecs(mediaEngine)
	mediaEngine.RegisterDefaultCodecs()
	webRTCInterceptor := &interceptor.Registry{}
	intervalPliFactory, err := intervalpli.NewReceiverInterceptor()
	if err != nil {
		return nil, err
	}
	webRTCInterceptor.Add(intervalPliFactory)

	if err = webrtc.RegisterDefaultInterceptors(mediaEngine, webRTCInterceptor); err != nil {
		return nil, err
	}

	api := webrtc.NewAPI(
		webrtc.WithMediaEngine(mediaEngine),
		webrtc.WithInterceptorRegistry(webRTCInterceptor),
	)
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}
	peerConnection, err := api.NewPeerConnection(config)
	if err != nil {
		return nil, err
	}
	return peerConnection, nil
}

func AddTransceiver(peerConnection *webrtc.PeerConnection) {
	for _, typ := range []webrtc.RTPCodecType{webrtc.RTPCodecTypeVideo, webrtc.RTPCodecTypeAudio} {
		if _, err := peerConnection.AddTransceiverFromKind(typ, webrtc.RTPTransceiverInit{
			Direction: webrtc.RTPTransceiverDirectionRecvonly,
		}); err != nil {
			log.Print(err)
			return
		}
	}
}

// SaveToDiskWithTimeout writes incoming RTP packets to media and local tracks within a specified duration
func SaveToDiskWithTimeout(ctx context.Context, i media.Writer, track *webrtc.TrackRemote, trackLocal *webrtc.TrackLocalStaticRTP, duration time.Duration, introduceChannel *chan struct{}) error {
	ctx, cancel := context.WithTimeout(ctx, duration)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			i.Close()
			*introduceChannel = make(chan struct{})
			return nil
		default:
			if err := processRTPPackets(track, trackLocal, i); err != nil {
				return err
			}
		}
	}
}

func processRTPPackets(track *webrtc.TrackRemote, trackLocal *webrtc.TrackLocalStaticRTP, i media.Writer) error {
	rtpPacket, _, err := track.ReadRTP()
	if err != nil {
		return err
	}

	if err := i.WriteRTP(rtpPacket); err != nil {
		return err
	}

	if err := trackLocal.WriteRTP(rtpPacket); err != nil {
		return err
	}

	return nil
}

// Relay relays RTP packets from a remote track to a local track
func Relay(track *webrtc.TrackRemote, trackLocal *webrtc.TrackLocalStaticRTP) error {
	rtpPacket, _, err := track.ReadRTP()
	if err != nil {
		return err
	}

	if err := trackLocal.WriteRTP(rtpPacket); err != nil {
		return err
	}

	return nil
}

func CheckExpireRooms() {
	for range time.NewTicker(time.Second * 5).C {
		currentTime := time.Now()
		for roomID, expiryTime := range models.ExpiryRooms {
			fmt.Println("RoomId :", roomID, "Current Time :", currentTime, "Expire Time :", expiryTime)
			if expiryTime != nil && currentTime.After(*expiryTime) {
				// Room has expired
				// Perform actions like closing connections, deleting the room, etc.
				SubmitResultPoll(roomID)
				closeRoomConnections(roomID)
				delete(models.Rooms, roomID)
				delete(models.ExpiryRooms, roomID)
				global.ROOM_SERVER_CLIENT.CloseRoom(context.Background(), &pb_room.CloseRoomByRoomIdRequest{
					RoomId: roomID,
				})
				// Add any other necessary handling for the expired room
			}
		}
	}
}

func closeRoomConnections(roomID string) {
	room, ok := models.Rooms[roomID]
	if !ok {
		return
	}

	room.Peers.ListLock.Lock()
	defer room.Peers.ListLock.Unlock()

	for i := range room.Peers.Connections {
		peer := &room.Peers.Connections[i]
		if peer.Websocket != nil {
			// Close WebSocket connection
			room.BroadcastNotify("close", "جلسه به پایان رسید")
			_ = peer.Websocket.Conn.Close()
		}
		if peer.PeerConnection != nil {
			// Close PeerConnection
			_ = peer.PeerConnection.Close()
		}
	}
}

func CloseRoom(roomId string) {
	global.ROOM_SERVER_CLIENT.CloseRoom(context.Background(), &pb_room.CloseRoomByRoomIdRequest{
		RoomId: roomId,
	})
}

type UserPoll struct {
	Status string `json:"status"`
	UserId string `json:"meet_user_id"`
}

func SubmitResultPoll(roomId string) {
	url := fmt.Sprintf("https://estate.sedrehgroup.ir/api/meets/%s/status/", roomId)
	polls := make([]UserPoll, 0)
	approvers := make([]string, 0)
	deniers := make([]string, 0)
	for id, p := range models.Rooms[roomId].Polls.UserPolls {
		var status string
		if p {
			status = "approval"
			approvers = append(approvers, id)

		} else {
			status = "opposition"
			deniers = append(deniers, id)
		}
		polls = append(polls, UserPoll{Status: status, UserId: id})

	}
	_, err := global.ROOM_SERVER_CLIENT.AddRoomResult(context.Background(), &pb_room.AddRoomResult{
		RoomId:    roomId,
		Approvers: approvers,
		Deniers:   deniers,
	})
	if err != nil {
		fmt.Println("Error grpc request:", err)
		// return
	}

	jsonData, err := json.Marshal(polls)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	// Send POST request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "token e8ac642409f2e01d668d7c518eef59ee6da372b2")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Response Status:", resp.Status)
}

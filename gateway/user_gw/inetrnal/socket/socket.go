package socket

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
	"vpeer_usergw/inetrnal/global"
	"vpeer_usergw/inetrnal/helpers"
	"vpeer_usergw/inetrnal/models"
	"vpeer_usergw/inetrnal/services"
	pb_room "vpeer_usergw/proto/api/room"

	"github.com/gofiber/contrib/websocket"
	"github.com/pion/webrtc/v3"
)

type (
	SocketController interface {
		WebsocketHandler(*websocket.Conn)
		UploadRecords()
	}

	socketController struct {
		minioService             services.MinioService
		minioDownloadedFilesPath string
	}
)

func NewSocketController(minioService services.MinioService, minioDownloadedFilesPath string) SocketController {
	return &socketController{
		minioService:             minioService,
		minioDownloadedFilesPath: minioDownloadedFilesPath,
	}
}

func (c *socketController) WebsocketHandler(conn *websocket.Conn) {
	defer func() {
		conn.Close() // Close the WebSocket connection when the function exits
		// Additional cleanup or handling if needed
	}()

	res, displayName, wserr := helpers.VerifyToken(conn)
	if wserr != nil {
		conn.WriteJSON(wserr)
		return
	}
	timeoutDuration := 1 * time.Hour
	conn.SetReadDeadline(time.Now().Add(timeoutDuration))
	conn.SetWriteDeadline(time.Now().Add(timeoutDuration))
	threadSafeWriter := models.ThreadSafeWriter{Conn: conn, Mutex: sync.Mutex{}}
	defer threadSafeWriter.Close()

	peerConnection, err := helpers.CreatePeerConnection()
	if err != nil {
		log.Print(err)
		return
	}

	defer peerConnection.Close()

	helpers.AddTransceiver(peerConnection)

	room, terr := helpers.GetRoomOrInit(res.RoomId)
	if terr != nil {
		conn.WriteJSON(terr)
		return
	}
	newPeer, wserr := models.CreateNewPeer(*res, displayName, &threadSafeWriter, peerConnection)
	if wserr != nil {
		conn.WriteJSON(wserr)
		return
	}

	peers := room.Peers
	peers.ListLock.Lock()
	peers.Connections = append(peers.Connections, *newPeer)
	peers.ListLock.Unlock()

	models.Rooms[res.RoomId].Polls.UserPolls[newPeer.UserId] = false

	helpers.OnICECandidate(peerConnection, &threadSafeWriter)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var ivfFile *models.IvfFileStruct
	var oggFile *models.OggFileStruct

	peerConnection.OnConnectionStateChange(func(cp webrtc.PeerConnectionState) {
		switch cp {
		case webrtc.PeerConnectionStateConnected:
			room.BroadcastNotify("connected", models.NewUserDataMessage(displayName, newPeer.ConnectionId, newPeer.UserId))
			if !models.Rooms[res.RoomId].Filled {
				if room.UserLength == int32(len(models.Rooms[res.RoomId].Peers.Connections)) {
					delete(models.ExpiryRooms, res.RoomId)
					expiryTime := time.Now().Add(time.Second * time.Duration((room.UserLength*15)+30))
					models.ExpiryRooms[res.RoomId] = &expiryTime
					models.Rooms[res.RoomId].Filled = true
				}
			}
		case webrtc.PeerConnectionStateFailed:
			if err := peerConnection.Close(); err != nil {
				log.Print(err)
			}
		case webrtc.PeerConnectionStateClosed:
			peers.SignalPeerConnections()
			if models.Rooms[res.RoomId] != nil && len(models.Rooms[res.RoomId].Peers.Connections) == 0 {
				helpers.SubmitResultPoll(res.RoomId)
				helpers.CloseRoom(res.RoomId)
				delete(models.Rooms, res.RoomId)
				delete(models.ExpiryRooms, res.RoomId)
			} else {
				room.BroadcastNotify("disconnected", models.NewUserDataMessage(displayName, newPeer.ConnectionId, newPeer.UserId))

				// for i := range models.Rooms[res.RoomId].Peers.Connections {
				// models.Rooms[res.RoomId].Peers.Connections = append(models.Rooms[res.RoomId].Peers.Connections[:i], models.Rooms[res.RoomId].Peers.Connections[i+1:]...)
				// break
				// }
			}
			if ivfFile != nil && oggFile != nil {
				models.Recorded[res.UserId] = &models.OpenFile{
					OggFile:     oggFile.OggFile,
					IvfFile:     ivfFile.IvfFile,
					OggFileName: oggFile.OggFileName,
					IvfFileName: ivfFile.IvfFileName,
					UserId:      res.UserId,
					RoomId:      res.RoomId,
					WebSocket:   &threadSafeWriter,
				}
				users := make([]models.User, 0)
				for _, connection := range models.Rooms[res.RoomId].Peers.Connections {
					users = append(users, models.User{DisplayName: connection.DisplayName, UserId: connection.UserId, ConnectionId: connection.ConnectionId, Recording: connection.RecordingState, IvfFileName: connection.IvfFileName, OggFileName: connection.OggFileName})
				}
				room.BroadcastNotify("user_data", models.Users{Users: users})

				ivfFile = nil
				oggFile = nil
			}
			global.ROOM_SERVER_CLIENT.AddRoomLog(context.Background(), &pb_room.AddRoomLog{
				RoomId:    res.RoomId,
				UserId:    res.UserId,
				UserEvent: "left",
			})
		}
	})

	var introduceChannel = make(chan struct{})

	timeDuration := time.Second * 15
	peerConnection.OnTrack(func(t *webrtc.TrackRemote, _ *webrtc.RTPReceiver) {
		trackLocal := peers.AddTrack(t)
		defer peers.RemoveTrack(trackLocal.TrackLocal)
		codec := t.Codec()
		if strings.EqualFold(codec.MimeType, webrtc.MimeTypeOpus) {
			defer func() {
				if oggFile != nil {
					if err := oggFile.OggFile.Close(); err != nil {
						return
					}
				}
			}()
			for {
				select {
				case <-introduceChannel:
					if err := helpers.SaveToDiskWithTimeout(ctx, oggFile.OggFile, t, trackLocal.TrackLocal, timeDuration, &introduceChannel); err != nil {
						return
					}
				default:
					if err := helpers.Relay(t, trackLocal.TrackLocal); err != nil {
						return
					}
				}
			}
		} else if strings.EqualFold(codec.MimeType, webrtc.MimeTypeVP8) {
			defer func() {
				if ivfFile != nil {
					if err := ivfFile.IvfFile.Close(); err != nil {
						return
					}
				}
			}()
			for {
				select {
				case <-introduceChannel:
					if err := helpers.SaveToDiskWithTimeout(ctx, ivfFile.IvfFile, t, trackLocal.TrackLocal, timeDuration, &introduceChannel); err != nil {
						return
					}
				default:
					if err := helpers.Relay(t, trackLocal.TrackLocal); err != nil {
						return
					}
				}
			}
		}
	})

	recording := false
	meetRecording := false
	meetRecordData := models.MkvFileStruct{}

	peers.SignalPeerConnections()
	message := &models.WebRTCMessage{}
	for {
		_, raw, err := conn.ReadMessage()
		if err != nil {
			if meetRecordData.MkvFile != nil {
				meetRecordData.MkvFile.Close()
				models.Recorded[res.UserId] = &models.OpenFile{
					MkvFileName: meetRecordData.MkvFileName,
					UserId:      res.UserId,
					RoomId:      res.RoomId,
					WebSocket:   &threadSafeWriter,
				}
				room.BroadcastNotify("meet_record", meetRecordData.MkvFileName)
				meetRecording = false
				meetRecordData = models.MkvFileStruct{}
			}
			if ivfFile != nil && oggFile != nil {
				models.Recorded[res.UserId] = &models.OpenFile{
					OggFile:     oggFile.OggFile,
					IvfFile:     ivfFile.IvfFile,
					OggFileName: oggFile.OggFileName,
					IvfFileName: ivfFile.IvfFileName,
					UserId:      res.UserId,
					RoomId:      res.RoomId,
					WebSocket:   &threadSafeWriter,
				}
				ivfFile = nil
				oggFile = nil
				users := make([]models.User, 0)
				for _, connection := range models.Rooms[res.RoomId].Peers.Connections {
					users = append(users, models.User{DisplayName: connection.DisplayName, UserId: connection.UserId, ConnectionId: connection.ConnectionId, Recording: connection.RecordingState, IvfFileName: connection.IvfFileName, OggFileName: connection.OggFileName})
				}
				room.BroadcastNotify("user_data", models.Users{Users: users})
			}
			return
		} else if err := json.Unmarshal(raw, &message); err != nil {
			meetRecordData.MkvFile.Write(raw)
		}
		switch message.Event {
		case "candidate":

			candidate := webrtc.ICECandidateInit{}
			if err := json.Unmarshal([]byte(message.Data), &candidate); err != nil {
				log.Println(err)
				return
			}

			if err := peerConnection.AddICECandidate(candidate); err != nil {
				log.Println(err)
				return
			}
		case "answer":
			answer := webrtc.SessionDescription{}
			if err := json.Unmarshal([]byte(message.Data), &answer); err != nil {
				log.Println(err)
				return
			}

			if err := peerConnection.SetRemoteDescription(answer); err != nil {
				log.Println(err)
				return
			}
		case "close":
			if res.UserId == room.OwnerId {
				allRecorded := true
				for _, u := range models.Rooms[res.RoomId].Peers.Connections {
					if u.UserId != room.OwnerId {
						if u.RecordingState != "recorded" {
							threadSafeWriter.WriteJSON(&models.WebSocketMessage{
								Event: "notification",
								State: "error",
								Data:  fmt.Sprintf("کاربر %s هنوز ویدیو خود را ظبط نکرده است", u.DisplayName),
							})
							allRecorded = false
						}
					}
				}
				if allRecorded {
					room.BroadcastNotify("close", "جلسه به پایان رسید")
					closeRoom := models.Rooms[res.RoomId]
					if closeRoom != nil {
						if closeRoom.OwnerId == res.UserId {
							for _, c := range closeRoom.Peers.Connections {
								c.PeerConnection.Close()
							}
							global.ROOM_SERVER_CLIENT.CloseRoom(context.TODO(), &pb_room.CloseRoomByRoomIdRequest{
								RoomId: res.RoomId,
							})
						}
					}
					return
				}
			} else {
				threadSafeWriter.WriteJSON(&models.WebSocketMessage{
					Event: "notification",
					State: "error",
					Data:  "شما ادمین جلسه نیستید",
				})
			}
		case "public_chat":
			room.BroadcastNotify("message", displayName+": "+message.Data)
		case "start_record":
			fmt.Println("#1")
			timeSecond, err := strconv.ParseInt(message.Data, 10, 32)
			if err != nil {
				threadSafeWriter.WriteJSON(&models.WebSocketMessage{
					Event: "error",
					State: "error",
					Data:  "وقت ظبط ویدیو به اشتباه وارد شده است",
				})
			} else {
				fmt.Println("#2")
				if ivfFile != nil {
					c.minioService.RemoveFile(ivfFile.IvfFileName)
				}
				if oggFile != nil {
					c.minioService.RemoveFile(oggFile.OggFileName)
				}
				fmt.Println("#3")
				if !recording {
					fmt.Println("#4")
					fileName := helpers.RandStringRunes(30)

					ivf, rerr := helpers.CreateIvfFile(c.minioService, fileName)
					if rerr != nil {
						return
					}
					ogg, err := helpers.CreateOggFile(c.minioService, fileName)
					if err != nil {
						return
					}

					fmt.Println("#5")

					ivfFile = ivf
					oggFile = ogg
					if timeSecond != 0 {
						timeDuration = time.Second * time.Duration(timeSecond)
					}
					recording = true
					var wg sync.WaitGroup
					wg.Add(1)

					go func() {
						defer wg.Done()
						close(introduceChannel)
					}()

					// Wait for the goroutine to complete
					wg.Wait()
					fmt.Println("#6")
					for i, connection := range models.Rooms[res.RoomId].Peers.Connections {
						if connection.UserId == res.UserId {
							models.Rooms[res.RoomId].Peers.Connections[i].RecordingState = "recording"
							models.Rooms[res.RoomId].Peers.Connections[i].IvfFileName = &ivf.IvfFileName
							models.Rooms[res.RoomId].Peers.Connections[i].OggFileName = &ogg.OggFileName
							break
						}
					}
					fmt.Println("#7")

				}
			}

		case "stop_record":
			fmt.Println("#8")
			if recording {
				fmt.Println("#9")
				recording = false
				var wg sync.WaitGroup
				wg.Add(1)

				go func() {
					defer wg.Done()
					ivfFile.IvfFile.Close()
					oggFile.OggFile.Close()
					introduceChannel = make(chan struct{})
				}()
				// Wait for the goroutine to complete
				wg.Wait()
				fmt.Println("#10")
				for i, connection := range models.Rooms[res.RoomId].Peers.Connections {
					if connection.UserId == res.UserId {
						models.Rooms[res.RoomId].Peers.Connections[i].RecordingState = "recorded"
						models.Rooms[res.RoomId].Peers.Connections[i].IvfFileName = &ivfFile.IvfFileName
						models.Rooms[res.RoomId].Peers.Connections[i].OggFileName = &oggFile.OggFileName
						break
					}
				}
				fmt.Println("#11")
				room.SingleNotify("notification", models.User{DisplayName: displayName, UserId: res.UserId, Recording: "recorded", IvfFileName: &ivfFile.IvfFileName, OggFileName: &oggFile.OggFileName}, &threadSafeWriter)
			}
		case "confirm_file":
			if ivfFile != nil && oggFile != nil {
				for i, connection := range models.Rooms[res.RoomId].Peers.Connections {
					if connection.UserId == res.UserId {
						models.Rooms[res.RoomId].Peers.Connections[i].RecordingState = "confirmed"
						break
					}
				}
				models.Recorded[res.UserId] = &models.OpenFile{
					OggFile:     oggFile.OggFile,
					IvfFile:     ivfFile.IvfFile,
					OggFileName: oggFile.OggFileName,
					IvfFileName: ivfFile.IvfFileName,
					UserId:      res.UserId,
					RoomId:      res.RoomId,
					WebSocket:   &threadSafeWriter,
				}
				ivfFile = nil
				oggFile = nil
			}
			users := make([]models.User, 0)
			for _, connection := range models.Rooms[res.RoomId].Peers.Connections {
				users = append(users, models.User{DisplayName: connection.DisplayName, UserId: connection.UserId, ConnectionId: connection.ConnectionId, Recording: connection.RecordingState, IvfFileName: connection.IvfFileName, OggFileName: connection.OggFileName})
			}
			room.BroadcastNotify("user_data", models.Users{Users: users})
		case "get_users":
			users := make([]models.User, 0)
			for _, connection := range models.Rooms[res.RoomId].Peers.Connections {
				users = append(users, models.User{DisplayName: connection.DisplayName, UserId: connection.UserId, ConnectionId: connection.ConnectionId, Recording: connection.RecordingState, IvfFileName: connection.IvfFileName, OggFileName: connection.OggFileName})
			}
			room.SingleNotify("user_data", models.Users{Users: users}, &threadSafeWriter)
		case "self_mute":
			room.BroadcastNotify("mute", models.NewUserDataMessage(displayName, newPeer.ConnectionId, newPeer.UserId))
		case "self_unmute":
			room.BroadcastNotify("unmute", models.NewUserDataMessage(displayName, newPeer.ConnectionId, newPeer.UserId))
		case "connection":
			room.SingleNotify("self_data", models.NewUserDataMessage(displayName, newPeer.ConnectionId, newPeer.UserId), &threadSafeWriter)
		case "raise_hand":
			room.BroadcastNotify("raise_hand_request", models.NewUserDataMessage(displayName, newPeer.ConnectionId, newPeer.UserId))
		case "expire":
			if room.OwnerId == res.UserId {
				timeSecond, err := strconv.ParseInt(message.Data, 10, 64)
				if err != nil {
					threadSafeWriter.WriteJSON(&models.WebSocketMessage{
						Event: "error",
						State: "error",
						Data:  "وقت اتمام جلسه به اشتباه وارد شده است",
					})
				} else {
					expireTime := time.Unix(timeSecond, 0)
					models.ExpiryRooms[res.RoomId] = &expireTime
					room.SingleNotify("expire_verification", models.NewUserDataMessage(displayName, newPeer.ConnectionId, newPeer.UserId), &threadSafeWriter)
				}
			}
		case "init_meet_record":
			if !meetRecording {
				fileName := helpers.RandStringRunes(30)
				mkv, rerr := helpers.CreateMkvFile(c.minioService, fileName)
				if rerr != nil {
					return
				}
				meetRecordData.MkvFile = mkv.MkvFile
				meetRecordData.MkvFileName = mkv.MkvFileName
				meetRecording = true
			}
		case "stop_meet_record":
			if meetRecording {
				meetRecording = false
				if meetRecordData.MkvFile != nil {
					meetRecordData.MkvFile.Close()
					models.Recorded[res.UserId] = &models.OpenFile{
						MkvFileName: meetRecordData.MkvFileName,
						UserId:      res.UserId,
						RoomId:      res.RoomId,
						WebSocket:   &threadSafeWriter,
					}
					room.BroadcastNotify("meet_record", meetRecordData.MkvFileName)
					meetRecording = false
					meetRecordData = models.MkvFileStruct{}
				}
			}
		case "approve_meeting":
			models.Rooms[res.RoomId].Polls.UserPolls[newPeer.UserId] = true

		}
	}
}

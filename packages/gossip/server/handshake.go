package server

import (
	"bytes"
	"time"

	"github.com/iotaledger/hive.go/autopeering/server"
	"google.golang.org/protobuf/proto"

	pb "github.com/iotaledger/goshimmer/packages/gossip/server/proto"
)

const (
	versionNum          = 0
	handshakeExpiration = 20 * time.Second
)

// isExpired checks whether the given UNIX time stamp is too far in the past.
func isExpired(ts int64) bool {
	return time.Since(time.Unix(ts, 0)) >= handshakeExpiration
}

func newHandshakeRequest(toAddr string) ([]byte, error) {
	println("new request")

	m := &pb.HandshakeRequest{
		Version:   versionNum,
		To:        toAddr,
		Timestamp: time.Now().Unix(),
	}
	bytes, error := proto.Marshal(m)
	// for i := 0; i < len(bytes); i++ {
	// 	println(bytes[i])
	// }
	return bytes, error
	// return proto.Marshal(m)
}

func newHandshakeResponse(reqData []byte) ([]byte, error) {
	println("new response")

	m := &pb.HandshakeResponse{
		ReqHash: server.PacketHash(reqData),
	}
	return proto.Marshal(m)
}

func (t *TCP) validateHandshakeRequest(reqData []byte) bool {
	println("validating request")

	m := new(pb.HandshakeRequest)
	if err := proto.Unmarshal(reqData, m); err != nil {
		t.log.Debugw("invalid handshake",
			"err", err,
		)
		return false
	}
	if m.GetVersion() != versionNum {
		t.log.Debugw("invalid handshake",
			"version", m.GetVersion(),
			"want", versionNum,
		)
		return false
	}
	if isExpired(m.GetTimestamp()) {
		t.log.Debugw("invalid handshake",
			"timestamp", time.Unix(m.GetTimestamp(), 0),
		)
	}

	println("validating request: good")
	return true
}

func (t *TCP) validateHandshakeResponse(resData []byte, reqData []byte) bool {
	println("validating response")

	m := new(pb.HandshakeResponse)
	if err := proto.Unmarshal(resData, m); err != nil {
		t.log.Debugw("invalid handshake",
			"err", err,
		)
		return false
	}
	if !bytes.Equal(m.GetReqHash(), server.PacketHash(reqData)) {
		t.log.Debugw("invalid handshake",
			"hash", m.GetReqHash(),
		)
		return false
	}

	println("validating response: good")
	return true
}

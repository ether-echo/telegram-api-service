package grpc_server

import (
	"context"

	"github.com/ether-echo/telegram-api-service/pkg/logger"

	pb "github.com/ether-echo/protos/messageProcessor"
)

var (
	log = logger.Logger().Named("grpc_server").Sugar()
)

type IMessage interface {
	SendMessage(chatId int64, message, url, command string)
}
type MessageServer struct {
	pb.UnimplementedMessageServiceServer
	IMessage IMessage
}

func (m *MessageServer) SendMessage(ctx context.Context, req *pb.MessageRequest) (*pb.MessageResponse, error) {
	log.Infof("Send message to user %d: %s, %s", req.ChatId, req.Message, req.URL)

	m.IMessage.SendMessage(req.ChatId, req.Message, req.URL, req.Command)

	return &pb.MessageResponse{Success: true}, nil
}

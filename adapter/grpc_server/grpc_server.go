package grpc_server

import (
	"context"
	pb "github.com/ether-echo/protos/messageProcessor"
	"github.com/ether-echo/telegram-api-service/pkg/logger"
)

var (
	log = logger.Logger().Named("grpc_server").Sugar()
)

type IMessage interface {
	SendMessage(chatId int64, message, url string)
}
type MessageServer struct {
	pb.UnimplementedMessageServiceServer
	IMessage IMessage
}

func (m *MessageServer) SendNotification(ctx context.Context, req *pb.MessageRequest) (*pb.MessageResponse, error) {

	log.Infof("Send message to user %d: %s", req.ChatId, req.Message)

	m.IMessage.SendMessage(req.ChatId, req.Message, req.URL)

	return &pb.MessageResponse{Success: true}, nil
}

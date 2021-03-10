package broker_service

import (
	"grpc/stream/prototypes"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
)

type BrokerService struct {
	prototypes.UnimplementedBrokerServiceServer
}

func (x BrokerService) Publish(publisher prototypes.BrokerService_PublishServer) error {
	defer publisher.SendAndClose(&emptypb.Empty{})
	for {
		select {
		case <-publisher.Context().Done():
			return publisher.Context().Err()
		default:
			msg, err := publisher.Recv()
			if err != nil {
				return err
			}
			pushMessageToSubscriber(msg)
		}
	}
}

func (x BrokerService) Subscribe(cfg *prototypes.SubscriberConfig, subscriber prototypes.BrokerService_SubscribeServer) error {
	if cfg.ChannelSize <= 0 {
		cfg.ChannelSize = 128
	}
	ch := make(chan *prototypes.MessageEntry, cfg.ChannelSize)
	defer close(ch)

	topicMapping := make(map[string]struct{})
	for _, topic := range cfg.Topics {
		topicMapping[topic] = struct{}{}
	}

	name := uuid.New().String()
	addSubscriber(name, ch)
	defer delSubscriber(name)

	for {
		select {
		case <-subscriber.Context().Done():
			return subscriber.Context().Err()
		case msg := <-ch:
			_, ok := topicMapping[msg.Topic]
			if !ok {
				continue
			}
			err := subscriber.Send(msg)
			if err != nil {
				return err
			}
		}
	}
}

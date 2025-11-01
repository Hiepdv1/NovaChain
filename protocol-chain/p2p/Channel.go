package p2p

import (
	"context"
	"encoding/json"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/peer"
	log "github.com/sirupsen/logrus"
)

const ChannelBufSize = 1000

type Channel struct {
	ctx   context.Context
	pub   *pubsub.PubSub
	topic *pubsub.Topic
	sub   *pubsub.Subscription

	channelName string
	self        peer.ID
	Content     chan *ChannelContent

	worker *Worker[*pubsub.Message]
}

type ChannelContent struct {
	Message  string
	SendFrom string
	SendTo   string
	Payload  []byte
}

func JoinChannel(ctx context.Context, pub *pubsub.PubSub, selfId peer.ID, channelName string, subscribe bool) (*Channel, error) {
	topic, err := pub.Join(topicName(channelName))
	if err != nil {
		return nil, err
	}

	var sub *pubsub.Subscription
	if subscribe {
		sub, err = topic.Subscribe(pubsub.WithBufferSize(ChannelBufSize))
		if err != nil {
			return nil, err
		}
	} else {
		sub = nil
	}

	Channel := &Channel{
		ctx:         ctx,
		pub:         pub,
		topic:       topic,
		sub:         sub,
		channelName: channelName,
		self:        selfId,
		Content:     make(chan *ChannelContent, ChannelBufSize),
	}

	worker := NewWorker(1000, ctx, Error, Channel.HandleContent)
	worker.Start(1)

	Channel.worker = worker

	go Channel.readLoop()

	return Channel, nil
}

func (channel *Channel) ListPeers() []peer.ID {
	return channel.topic.ListPeers()
}

func (channel *Channel) Publish(message string, payload []byte, SendTo string) error {
	msg := ChannelContent{
		Message:  message,
		SendFrom: ShortID(channel.self),
		SendTo:   SendTo,
		Payload:  payload,
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return channel.topic.Publish(channel.ctx, msgBytes)
}

func (channel *Channel) HandleContent(content *pubsub.Message) {
	if content.ReceivedFrom == channel.self {
		return
	}

	NewContent := new(ChannelContent)
	err := json.Unmarshal(content.Data, NewContent)
	if err != nil {
		log.Errorf("Umashal Messsage with err: %v", err)
		return
	}

	if NewContent.SendTo != "" && NewContent.SendTo != channel.self.String() {
		return
	}

	select {
	case channel.Content <- NewContent:
	default:
		log.Warn("Worker queue full, skip message")
	}
}

func (channel *Channel) readLoop() {
	if channel.sub == nil {
		return
	}

	for {
		content, err := channel.sub.Next(channel.ctx)

		if err != nil {
			log.Errorf("Receiver Message with error: %v", err)
			close(channel.Content)
			return
		}

		channel.worker.Push(content)
	}
}

func topicName(channelName string) string {
	return "channel:" + channelName
}

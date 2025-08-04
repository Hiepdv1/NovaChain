package p2p

import (
	"context"
	"encoding/json"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/peer"
)

const ChannelBufSize = 128

type Channel struct {
	ctx   context.Context
	pub   *pubsub.PubSub
	topic *pubsub.Topic
	sub   *pubsub.Subscription

	channelName string
	self        peer.ID
	Content     chan *ChannelContent
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
		sub, err = topic.Subscribe()
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

	go Channel.readLoop()

	return Channel, nil
}

func (channel *Channel) ListPeers() []peer.ID {
	// return channel.pub.ListPeers(topicName(channel.channelName))
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

func (channel *Channel) readLoop() {
	if channel.sub == nil {
		return
	}

	for {
		content, err := channel.sub.Next(channel.ctx)
		if err != nil {
			close(channel.Content)
			return
		}

		if content.ReceivedFrom == channel.self {
			continue
		}

		NewContent := new(ChannelContent)
		err = json.Unmarshal(content.Data, NewContent)
		if err != nil {
			continue
		}

		if NewContent.SendTo != "" && NewContent.SendTo != channel.self.String() {
			continue
		}

		channel.Content <- NewContent
	}
}

func topicName(channelName string) string {
	return "channel:" + channelName
}

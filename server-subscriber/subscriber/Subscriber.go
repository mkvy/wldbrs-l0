package subscriber

import (
	"fmt"
	"github.com/mkvy/wldbrs-l0/server-subscriber/store"
	"github.com/nats-io/stan.go"
)

type StanSubscriber struct {
	sc stan.Conn
	ss store.StoreService
}

func CreateSub(ss store.StoreService) *StanSubscriber {
	sc := StanSubscriber{
		ss: ss,
	}
	return &sc
}

func (sSub *StanSubscriber) Connect(clusterID string, clientID string, URL string) error {
	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL(URL))
	if err != nil {
		return err
	}
	sSub.sc = sc
	return err
}

func (sSub *StanSubscriber) Close() {
	if sSub != nil {
		sSub.sc.Close()
	}
}

func (sSub *StanSubscriber) SubscribeToChannel(channel string, opts ...stan.SubscriptionOption) (stan.Subscription, error) {
	sub, err := sSub.sc.Subscribe(channel, sSub.handlerMsg, opts...)
	if err != nil {
		fmt.Println("Can't connect")
		sSub.ss.RestoreCache()
	}
	return sub, err
}

func (sSub *StanSubscriber) handlerMsg(msg *stan.Msg) {
	err := sSub.ss.SaveOrderData(msg.Data)
	if err != nil {
		fmt.Println("error while saving")
	}
}

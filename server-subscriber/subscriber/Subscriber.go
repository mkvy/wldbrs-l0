package subscriber

import "github.com/nats-io/stan.go"

type StanSubscriber struct {
	sc stan.Conn
}

func CreateSub() *StanSubscriber {
	sc := StanSubscriber{}
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

func (sSub *StanSubscriber) SubscribeToChannel(channel string, opts ...stan.SubscriptionOption) {

}

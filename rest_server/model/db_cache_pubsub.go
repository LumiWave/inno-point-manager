package model

import (
	"encoding/json"
	"strings"

	"github.com/ONBUFF-IP-TOKEN/basedb"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
)

func (o *DB) PublishEvent(channel string, val interface{}) error {
	msg, _ := json.Marshal(val)
	return o.Cache.GetDB().Publish(MakePubSubKey(channel), string(msg))
}

func (o *DB) ListenSubscribeEvent() error {
	defer func() {
		if recver := recover(); recver != nil {
			log.Error("Recoverd in listenPubSubEvent()", recver)
			go o.ListenSubscribeEvent()
		}
	}()

	log.Info("ListenSubscribeEvent() has been started")

	receiveCh := make(chan basedb.PubSubMessageV8)
	defer close(receiveCh)

	channel := MakePubSubKey(InternalCmd)
	rch, err := o.Cache.GetDB().Subscribe(receiveCh, channel)
	if err != nil {
		log.Error(err)
		return err
	}
	defer o.Cache.GetDB().ClosePubSub()

	for {
		msg, ok := <-rch
		if msg == nil || !ok {
			continue
		}

		if strings.Contains(msg.Channel, MakePubSubKey(InternalCmd)) {
			o.PubSubCmdByInternal(msg)
		}

		log.Debugf("subscribe channel: %v, val: %v", msg.Channel, msg.Payload)
	}

	return nil
}

func (o *DB) PubSubCmdByInternal(msg basedb.PubSubMessageV8) error {

	header := &PSHeader{}
	json.Unmarshal([]byte(msg.Payload), header)

	if strings.EqualFold(header.Type, PubSub_type_maintenance) {
		psPacket := &PSMaintenance{}
		json.Unmarshal([]byte(msg.Payload), psPacket)
		SetMaintenance(psPacket.Value.Enable)
	} else if strings.EqualFold(header.Type, PubSub_type_Swap) {
		psPacket := &PSSwap{}
		json.Unmarshal([]byte(msg.Payload), psPacket)
		SetSwapEnable(psPacket.Value.Enable)
	} else if strings.EqualFold(header.Type, PubSub_type_CoinTransferExternal) {
		psPacket := &PSCoinTransferExternal{}
		json.Unmarshal([]byte(msg.Payload), psPacket)
		SetExternalTransferEnable(psPacket.Value.Enable)
	} else if strings.EqualFold(header.Type, PubSub_type_meta_refresh) {
		// db meta refresh
		o.GetPointList()
		o.GetAppCoins()
		o.GetCoins()
		o.GetApps()
		o.GetAppPoints()
		o.GetBaseCoins()
		log.Infof("pubsub cmd : %v", PubSub_type_meta_refresh)
	} else if strings.EqualFold(header.Type, PubSub_type_point_update) {
		psPacket := &PSPointUpdate{}
		json.Unmarshal([]byte(msg.Payload), psPacket)
		SetPointUpdateEnable(psPacket.Value.Enable)
	}

	return nil
}

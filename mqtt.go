package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/eclipse/paho.mqtt.golang"
)

const mqttUser = `digitraffic`
const mqttPass = `digitrafficPassword`

const mqttTopic = `vessels/#`

type mqttHandler = func(mqtt.Client, mqtt.Message)

type cc struct {
	c     io.WriteCloser
	fmt   outFormatter
	errCh chan error
}

func (c *cc) onMessageReceived(client mqtt.Client, message mqtt.Message) {
	err := c.onMessage(client, message)
	if err != nil {
		c.errCh <- err
	}
}

func (c *cc) writeCombo(bs []byte, err error) error {
	if err != nil {
		return err
	}
	_, err = c.c.Write(bs)
	return err
}

func (c *cc) parseTopic(topic string) (string, error) {
	if topic == `vessels/status` {
		return `status`, nil
	}
	ss := strings.Split(topic, `/`)
	if len(ss) != 3 || ss[0] != `vessels` {
		return "", fmt.Errorf("invalid topic provided")
	}
	switch t := ss[2]; t {
	case `metadata`, `locations`:
		return t, nil
	}
	return "", fmt.Errorf("invalid topic type")
}

func (c *cc) onMessage(client mqtt.Client, message mqtt.Message) error {
	topic := message.Topic()
	bs := message.Payload()

	debugf("Received message on topic: %q\nMessage: %s\n", topic, bs)
	t, err := c.parseTopic(topic)
	if err != nil {
		return fmt.Errorf("fatal MQTT topic: %q %q: %w", topic, bs, err)
	}
	switch t {
	case `metadata`:
		var vmsg vesselMetadata
		err := parseVesselMetadata(bs, &vmsg)
		if err != nil {
			return fmt.Errorf("fatal MQTT metadata: %q %q: %w", topic, bs, err)
		}
		debugf("Decoded into: %#v", vmsg)
		err = c.writeCombo(c.fmt.FormatVesselMetadata(&vmsg))
		if err != nil {
			return err
		}
	case `locations`:
		var vmsg vesselLocation
		err := parseVesselLocation(bs, &vmsg)
		if err != nil {
			return fmt.Errorf("fatal MQTT locations: %q %q: %w", topic, bs, err)
		}
		debugf("Decoded into: %#v", vmsg)
		err = c.writeCombo(c.fmt.FormatVesselLocation(&vmsg))
		if err != nil {
			return err
		}
	case `status`:
		// FIXME parse status messages.
		debugf("Status: %q", bs)
		return nil
	default:
		return fmt.Errorf("invalid message topic: %q %q", topic, bs)

	}
	return nil
}

func dialMqtt(messageCallback mqttHandler) error {
	if *verbose {
		mqtt.DEBUG = log.New(os.Stderr, "", 0)
	}
	mqtt.ERROR = log.New(os.Stderr, "", 0)
	opt := mqtt.NewClientOptions()
	url := `wss://` + *serverName + `:61619/mqtt`
	debugf("Using mqtt url %q", url)
	opt.SetAutoReconnect(true).AddBroker(url).
		SetUsername(mqttUser).SetPassword(mqttPass).SetCleanSession(true)
	opt.OnConnect = func(c mqtt.Client) {}

	c := mqtt.NewClient(opt)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	if token := c.Subscribe(mqttTopic, 0, messageCallback); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

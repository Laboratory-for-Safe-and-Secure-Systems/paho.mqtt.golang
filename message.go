/*
 * Copyright (c) 2021 IBM Corp and others.
 *
 * All rights reserved. This program and the accompanying materials
 * are made available under the terms of the Eclipse Public License v2.0
 * and Eclipse Distribution License v1.0 which accompany this distribution.
 *
 * The Eclipse Public License is available at
 *    https://www.eclipse.org/legal/epl-2.0/
 * and the Eclipse Distribution License is available at
 *   http://www.eclipse.org/org/documents/edl-v10.php.
 *
 * Contributors:
 *    Seth Hoenig
 *    Allan Stockdill-Mander
 *    Mike Robertson
 */

package mqtt

import (
	"net/url"
	"sync"

	"github.com/Laboratory-for-Safe-and-Secure-Systems/paho.mqtt.golang/packets"
)

// Message defines the externals that a message implementation must support
// these are received messages that are passed to the callbacks, not internal
// messages
type Message interface {
	Duplicate() bool
	Qos() byte
	Retained() bool
	Topic() string
	MessageID() uint16
	Payload() []byte
	Ack()
}

type message struct {
	duplicate bool
	qos       byte
	retained  bool
	topic     string
	messageID uint16
	payload   []byte
	once      sync.Once
	ack       func()
}

func (m *message) Duplicate() bool {
	return m.duplicate
}

func (m *message) Qos() byte {
	return m.qos
}

func (m *message) Retained() bool {
	return m.retained
}

func (m *message) Topic() string {
	return m.topic
}

func (m *message) MessageID() uint16 {
	return m.messageID
}

func (m *message) Payload() []byte {
	return m.payload
}

func (m *message) Ack() {
	m.once.Do(m.ack)
}

func messageFromPublish(p *packets.PublishPacket, ack func()) Message {
	return &message{
		duplicate: p.Dup,
		qos:       p.Qos,
		retained:  p.Retain,
		topic:     p.TopicName,
		messageID: p.MessageID,
		payload:   p.Payload,
		ack:       ack,
	}
}

func newConnectMsgFromOptions(options *ClientOptions, broker *url.URL) *packets.ConnectPacket {
	m := packets.NewControlPacket(packets.Connect).(*packets.ConnectPacket)

	m.CleanSession = options.CleanSession
	m.WillFlag = options.WillEnabled
	m.WillRetain = options.WillRetained
	m.ClientIdentifier = options.ClientID

	if options.WillEnabled {
		m.WillQos = options.WillQos
		m.WillTopic = options.WillTopic
		m.WillMessage = options.WillPayload
	}

	username := options.Username
	password := options.Password
	if broker.User != nil {
		username = broker.User.Username()
		if pwd, ok := broker.User.Password(); ok {
			password = pwd
		}
	}
	if options.CredentialsProvider != nil {
		username, password = options.CredentialsProvider()
	}

	if username != "" {
		m.UsernameFlag = true
		m.Username = username
		// mustn't have password without user as well
		if password != "" {
			m.PasswordFlag = true
			m.Password = []byte(password)
		}
	}

	m.Keepalive = uint16(options.KeepAlive)

	return m
}

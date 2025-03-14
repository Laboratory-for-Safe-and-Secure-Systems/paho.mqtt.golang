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
	"bytes"
	"testing"

	"github.com/Laboratory-for-Safe-and-Secure-Systems/paho.mqtt.golang/packets"
)

func Test_NewPingReqMessage(t *testing.T) {
	pr := packets.NewControlPacket(packets.Pingreq).(*packets.PingreqPacket)
	if pr.MessageType != packets.Pingreq {
		t.Errorf("NewPingReqMessage bad msg type: %v", pr.MessageType)
	}
	if pr.RemainingLength != 0 {
		t.Errorf("NewPingReqMessage bad remlen, expected 0, got %d", pr.RemainingLength)
	}

	exp := []byte{
		0xC0,
		0x00,
	}

	var buf bytes.Buffer
	if err := pr.Write(&buf); err != nil {
		t.Fatal(err)
	}
	bs := buf.Bytes()

	if len(bs) != 2 {
		t.Errorf("NewPingReqMessage.Bytes() wrong length: %d", len(bs))
	}

	if exp[0] != bs[0] || exp[1] != bs[1] {
		t.Errorf("NewPingMessage.Bytes() wrong")
	}
}

func Test_DecodeMessage_pingresp(t *testing.T) {
	bs := bytes.NewBuffer([]byte{
		0xD0,
		0x00,
	})
	presp, _ := packets.ReadPacket(bs)
	if presp.(*packets.PingrespPacket).MessageType != packets.Pingresp {
		t.Errorf("DecodeMessage ping response wrong msg type: %v", presp.(*packets.PingrespPacket).MessageType)
	}
	if presp.(*packets.PingrespPacket).RemainingLength != 0 {
		t.Errorf("DecodeMessage ping response wrong rem len: %d", presp.(*packets.PingrespPacket).RemainingLength)
	}
}

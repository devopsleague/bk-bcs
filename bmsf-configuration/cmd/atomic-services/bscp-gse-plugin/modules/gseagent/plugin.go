/*
Tencent is pleased to support the open source community by making Blueking Container Service available.
Copyright (C) 2019 THL A29 Limited, a Tencent company. All rights reserved.
Licensed under the MIT License (the "License"); you may not use this file except
in compliance with the License. You may obtain a copy of the License at
http://opensource.org/licenses/MIT
Unless required by applicable law or agreed to in writing, software distributed under
the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
either express or implied. See the License for the specific language governing permissions and
limitations under the License.
*/

package gseagent

/*
#include <stdlib.h>
#include <stdio.h>

int sendMsg(char * meta, int metalen, char * body, int bodylen);
int regFuncSendMsg(void* pFunc);
*/
import "C"
import "unsafe"

/* DO NOT EDIT THIS FILE, UNLESS YOU KNOW WHAT YOU DOING! */

// import other golang packages here.
import (
	"encoding/json"

	"bk-bscp/pkg/logger"
)

// TransmitType is gse agent message transmit type.
type TransmitType int

const (
	// TTBROADCAST broadcast message to all reachable taskserver.
	TTBROADCAST TransmitType = 0

	// TTP2P send message to target taskserver.
	TTP2P TransmitType = 1

	// TTRAND send message to one taskserver in rand mode.
	TTRAND TransmitType = 2
)

// Meta is agent message metadata.
type Meta struct {
	// SessionID is tunnel session id from taskserver, it's empty
	// if the message from agent to tunnelserver.
	SessionID int64 `json:"session_id"`

	// MessageID is agent message sequence id.
	MessageID int64 `json:"msgseq_id"`

	// TransmitType is gse agent message transmit type.
	TransmitType TransmitType `json:"transmit_type"`
}

// RegFuncSendMsg is gse agent regFuncSendMsg wrapper.
func RegFuncSendMsg(pFunc unsafe.Pointer) int {
	return int(C.regFuncSendMsg(pFunc))
}

// SendMessage is gse agent sendMsg wrapper.
func SendMessage(msg string, sessionID, messageID int64, transmitType TransmitType) int {
	csMsg := C.CString(msg)
	defer C.free(unsafe.Pointer(csMsg))

	// build meta.
	meta := &Meta{
		SessionID:    sessionID,
		MessageID:    messageID,
		TransmitType: transmitType,
	}
	metaInfo, err := json.Marshal(meta)
	if err != nil {
		return -1
	}
	logger.V(4).Infof("GSE Agent| send plugin message, %+v", meta)

	csMeta := C.CString(string(metaInfo))
	defer C.free(unsafe.Pointer(csMeta))

	return int(sendMessage(csMeta, C.int(len(metaInfo)), csMsg, C.int(len(msg))))
}

func sendMessage(meta *C.char, metalen C.int, body *C.char, bodylen C.int) C.int {
	return C.sendMsg(meta, metalen, body, bodylen)
}

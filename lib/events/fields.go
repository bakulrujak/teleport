/*
Copyright 2019 Gravitational, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package events

import (
	"bytes"
	"text/template"
	"time"

	"github.com/gravitational/teleport/lib/utils"

	"github.com/gravitational/trace"
	"github.com/jonboulle/clockwork"
)

// AugmentEventFields updates passed event fields with additional information
// common for all event types such as unique IDs, timestamps, codes, etc.
//
// This method is a "final stop" for various audit log implementations for
// updating event fields before it gets persisted in the backend.
func AugmentEventFields(event string, fields EventFields, clock clockwork.Clock, uid utils.UID) (err error) {
	additionalFields := make(map[string]interface{})
	if fields.GetType() == "" {
		additionalFields[EventType] = event
	}
	if fields.GetID() == "" {
		additionalFields[EventID] = uid.New()
	}
	if fields.GetTimestamp().IsZero() {
		additionalFields[EventTime] = clock.Now().UTC().Round(time.Second)
	}
	additionalFields[EventCode], err = getEventCode(event, fields)
	if err != nil {
		return trace.Wrap(err)
	}
	additionalFields[EventMessage], err = getEventMessage(event, fields)
	if err != nil {
		return trace.Wrap(err)
	}
	for k, v := range additionalFields {
		fields[k] = v
	}
	return nil
}

func getEventCode(event string, fields EventFields) (string, error) {
	switch event {
	case UserLoginEvent:
		switch fields.GetString(LoginMethod) {
		case LoginMethodLocal:
			if fields.GetError() != nil {
				return UserLocalLoginFailedCode, nil
			}
			return UserLocalLoginCode, nil
		default:
			if fields.GetError() != nil {
				return UserSSOLoginFailedCode, nil
			}
			return UserSSOLoginCode, nil
		}
	case SessionStartEvent:
		return SessionStartCode, nil
	case SessionJoinEvent:
		return SessionJoinCode, nil
	case SessionLeaveEvent:
		return SessionLeaveCode, nil
	case SessionEndEvent:
		return SessionEndCode, nil
	case SessionUploadEvent:
		return SessionUploadCode, nil
	case SubsystemEvent:
		return SubsystemCode, nil
	case ExecEvent: // TODO Has failure counterpart?
		return ExecCode, nil
	case PortForwardEvent:
		return PortForwardCode, nil
	case SCPEvent: // TODO Has failure counterpart?
		return SCPCode, nil
	case ResizeEvent:
		return ResizeCode, nil
	case ClientDisconnectEvent:
		return ClientDisconnectCode, nil
	case AuthAttemptEvent:
		return AuthAttemptCode, nil
	}
	return "", trace.BadParameter("unknown event type %q", event)
}

func getEventMessage(eventCode string, fields EventFields) (string, error) {
	messageTemplate := fields.GetMessage()
	if messageTemplate == "" {
		var ok bool
		if messageTemplate, ok = codeToMessage[eventCode]; !ok {
			return "", trace.BadParameter("no message template for event %q", eventCode)
		}
	}
	template, err := template.New("message").Parse(messageTemplate)
	if err != nil {
		return "", trace.Wrap(err)
	}
	var b bytes.Buffer
	if err := template.Execute(&b, fields); err != nil {
		return "", trace.Wrap(err)
	}
	return b.String(), nil
}

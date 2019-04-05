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

const (
	UserLocalLoginCode       = "T1000I"
	UserLocalLoginFailedCode = "T1000W"
	UserSSOLoginCode         = "T1001I"
	UserSSOLoginFailedCode   = "T1001W"
	SessionStartCode         = "T2000I"
	SessionJoinCode          = "T2001I"
	SessionLeaveCode         = "T2002I"
	SessionEndCode           = "T2003I"
	SessionUploadCode        = "T2004I"
	SubsystemCode            = "T3000I"
	ExecCode                 = "T3001I"
	ExecErrorCode            = "T3001E"
	PortForwardCode          = "T3002I"
	SCPCode                  = "T3003I"
	ResizeCode               = "T3004I"
	ClientDisconnectCode     = "T3005I"
	AuthAttemptCode          = "T3006W"
)

var (
	codeToMessage = map[string]string{
		UserLocalLoginCode:       "Local user {{.user}} successfully logged in",
		UserLocalLoginFailedCode: "Local user {{.user}} login failed: {{.error}}",
		UserSSOLoginCode:         "SSO user {{.user}} successfully logged in",
		UserSSOLoginFailedCode:   "SSO user {{.user}} login failed: {{.error}}",
		SessionStartCode:         "User {{.user}} has started a session",
		SessionJoinCode:          "User {{.user}} has joined the session",
		SessionLeaveCode:         "User {{.user}} has left the session",
		SessionEndCode:           "User {{.user}} has ended the session",
		SessionUploadCode:        "Recorded session has been uploaded",
		SubsystemCode:            "User {{.user}} requested subsystem {{.name}}",
		ExecCode:                 `User {{.user}} executed command on node {{index . "addr.remote"}}`,
		ExecErrorCode:            `User {{.user}} command execution on node {{index . "addr.remote"}} failed: {{.exitError}}`,
		PortForwardCode:          "User {{.user}} started port forwarding",
		SCPCode:                  "User {{.user}} {{.action}}ed file {{.path}}",
		ResizeCode:               "User {{.user}} resized the terminal",
		ClientDisconnectCode:     "User {{.user}} has been disconnected: {{.reason}}",
		AuthAttemptCode:          "User {{.user}} failed auth attempt",
	}
)

package dto

type RequestDTO struct {
	Action string      `json:"action"`
	Auth   LoginDTO    `json:"auth,omitempty"`
	Log    LogDTO      `json:"log,omitempty"`
	Mail   SendMailDTO `json:"mail,omitempty"`
}

package main

type answer struct {
	Success bool                   `json:"succes"`
	Error   string                 `json:"error"`
	Res     map[string]interface{} `json:"result"`
}

type message struct {
	From    string `json:"from_name"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

type sendMessageReq struct {
	PeerName string `json:"peer_name"`
	Message  string `json:"message"`
}

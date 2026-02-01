package store

import (
	"encoding/json"
	"net"
)

type PairingAction struct {
	TaskID string `json:"task_id"`
	Action string `json:"action"`
	Source string `json:"source"`
}

func ReportToPairing(action PairingAction) error {
	socketPath := "/tmp/ami-pairing.sock"
	
	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		// Silent fail if daemon not running
		return nil
	}
	defer conn.Close()

	data, err := json.Marshal(action)
	if err != nil {
		return err
	}

	_, err = conn.Write(append(data, '\n'))
	return err
}

func GetSocketPath() string {
	return "/tmp/ami-pairing.sock"
}

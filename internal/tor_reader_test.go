package internal

import "testing"

func TestTorReader(t *testing.T) {
	var tp TorProgress

	tp.Write([]byte(`Read line: 650 STATUS_CLIENT NOTICE BOOTSTRAP PROGRESS=15 TAG=onehop_create SUMMARY="Establishing an encrypted directory connection"` + "\n"))
}

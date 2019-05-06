package internal

import (
	"bytes"
	"regexp"
	"strconv"
	"strings"
)

type TorProgress struct {
	buf        bytes.Buffer
	OnProgress func(info string, percentage int)
}

func (t *TorProgress) Write(p []byte) (int, error) {
	for i := 0; i < len(p); i++ {
		t.buf.WriteByte(p[i])
		if p[i] == '\n' {
			t.gotLine(t.buf.Bytes())
			t.buf.Reset()
		}
	}
	return len(p), nil
}

var re = regexp.MustCompile(`Read line: 650 STATUS_CLIENT NOTICE BOOTSTRAP PROGRESS=(\d+) TAG=[\w_]+ SUMMARY="([\w\s]+)"`)

func (t *TorProgress) gotLine(p []byte) {
	s := strings.TrimSpace(string(p))

	if strings.HasPrefix(s, "Read line: 650 STATUS_CLIENT NOTICE BOOTSTRAP PROGRESS=") {
		matches := re.FindStringSubmatch(s)
		if len(matches) == 3 {
			percent, err := strconv.ParseInt(matches[1], 10, 8)
			if err != nil {
				return
			}

			summary := matches[2]
			t.OnProgress(summary, int(percent))
			return
		}
	}

	if s == "Read line: 650 STATUS_CLIENT NOTICE CIRCUIT_ESTABLISHED" {
		t.OnProgress("Circuit established", -1)
		return
	}

	if s == "Waiting for publication" {
		t.OnProgress("Waiting for publication", -1)
		return
	}
}

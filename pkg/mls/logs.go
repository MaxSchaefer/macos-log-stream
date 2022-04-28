package mls

import (
	"bufio"
	"encoding/json"
	"os/exec"
	"sync"
)

type Log struct {
	TraceID            int64       `json:"traceID"`
	EventMessage       string      `json:"eventMessage"`
	EventType          string      `json:"eventType"`
	Source             interface{} `json:"source"`
	FormatString       string      `json:"formatString"`
	ActivityIdentifier int         `json:"activityIdentifier"`
	Subsystem          string      `json:"subsystem"`
	Category           string      `json:"category"`
	ThreadID           int         `json:"threadID"`
	SenderImageUUID    string      `json:"senderImageUUID"`
	Backtrace          struct {
		Frames []struct {
			ImageOffset int    `json:"imageOffset"`
			ImageUUID   string `json:"imageUUID"`
		} `json:"frames"`
	} `json:"backtrace"`
	BootUUID                 string `json:"bootUUID"`
	ProcessImagePath         string `json:"processImagePath"`
	Timestamp                string `json:"timestamp"`
	SenderImagePath          string `json:"senderImagePath"`
	MachTimestamp            int64  `json:"machTimestamp"`
	MessageType              string `json:"messageType"`
	ProcessImageUUID         string `json:"processImageUUID"`
	ProcessID                int    `json:"processID"`
	SenderProgramCounter     int    `json:"senderProgramCounter"`
	ParentActivityIdentifier int    `json:"parentActivityIdentifier"`
	TimezoneName             string `json:"timezoneName"`
}

type Logs struct {
	m         sync.Mutex
	Predicate string
	Channel   chan Log
	exit      chan bool
}

func NewLogs() *Logs {
	return &Logs{
		Channel: make(chan Log),
		exit:    make(chan bool),
	}
}

func (logs *Logs) StartGathering() error {
	args := []string{
		"stream",
		"--color=none",
		"--style=ndjson",
	}

	if logs.Predicate != "" {
		args = append(args, "--predicate", logs.Predicate)
	}

	cmd := exec.Command("log", args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	go func() {
		logs.m.Lock()
		defer logs.m.Unlock()

		cmd.Start()
		defer cmd.Process.Kill()

		// drop first message
		bufio.NewReader(stdout).ReadLine()

		dec := json.NewDecoder(stdout)

		for {
			select {
			case <-logs.exit:
				return
			default:
				log := Log{}
				if err := dec.Decode(&log); err != nil {
					panic(err)
				} else {
					logs.Channel <- log
				}
			}
		}
	}()

	go func() {
		stderrBuf := bufio.NewReader(stderr)
		for {
			line, _, _ := stderrBuf.ReadLine()
			if len(line) > 0 {
				logs.StopGathering()
				panic(err)
			}
		}
	}()

	return nil
}

func (logs *Logs) StopGathering() {
	logs.exit <- true
}

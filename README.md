# macos-log-stream

Go wrapper around the `> log stream` command of the unified logging system on MacOS.

## pkg/mls

This example shows the main purpose of the package.
If you start gathering logs, the mls package runs `log stream --color none --style ndjson`.
When `logs.Predicate` is set, the mls package appends `... --predicate [logs.Predicate]` to the command.

```go
package main

import (
	"fmt"
	"github.com/MaxSchaefer/go-macos-log-stream/pkg/mls"
)

func main() {
	logs := mls.NewLogs()

	// logs.Predicate = ""

	if err := logs.StartGathering(); err != nil {
		panic(err)
	}

	for log := range logs.Channel {
		fmt.Println(log.EventMessage)
	}
}
```

## cli

To build the cli run `make build`.

```shell
# Show logs when the webcam turns on or off
mls -predicate 'subsystem == "com.apple.UVCFamily" && (eventMessage CONTAINS[c] "start stream" || eventMessage CONTAINS[c] "stop stream")'
```

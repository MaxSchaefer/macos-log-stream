.PHONY: dev build run

BIN := mls

dev:
	go get github.com/githubnemo/CompileDaemon
	CompileDaemon -exclude-dir=.git -color=true -build="go build -o ${BIN}" -command="./${BIN}"

build:
	go build -o ${BIN}
	du -h ${BIN}

run: build
	./${BIN} -predicate 'subsystem == "com.apple.UVCFamily" && (eventMessage CONTAINS[c] "start stream" || eventMessage CONTAINS[c] "stop stream")'

mls: build

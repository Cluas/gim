# Go parameters
GOCMD=GO111MODULE=on go
GOBUILD=$(GOCMD) build

build:
	rm -rf target/
	mkdir target/
	cp cmd/comet/comet-example.toml target/comet.toml
	cp cmd/logic/logic-example.toml target/logic.toml
	cp cmd/job/job-example.toml target/job.toml
	$(GOBUILD) -o target/comet cmd/comet/main.go
	$(GOBUILD) -o target/logic cmd/logic/main.go
	$(GOBUILD) -o target/job cmd/job/main.go

clean:
	rm -rf target/

run:
	nohup target/logic -p=target &
	nohup target/comet -p=target &
	nohup target/job -p=target &

stop:
	pkill -f target/logic
	pkill -f target/job
	pkill -f target/comet

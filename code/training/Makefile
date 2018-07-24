all: compile docker push clean

compile:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o trainModel

docker:
	docker build --force-rm=true -t pskreter/train-req-model:multi-1 .

push:
	docker push pskreter/train-req-model:multi-1

clean:
	rm trainModel

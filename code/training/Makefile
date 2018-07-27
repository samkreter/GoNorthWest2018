all: compile docker push clean

compile:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o trainModel

docker:
	docker build --force-rm=true -t pskreter/train-req-model:test .

push:
	docker push pskreter/train-req-model:test

clean:
	rm trainModel

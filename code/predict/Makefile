all: compile docker push clean

compile:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o predict

docker:
	docker build --force-rm=true -t pskreter/predict-req:1.0.1 .

push:
	docker push pskreter/predict-req:1.0.1

clean:
	rm predict
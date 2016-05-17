GOPATH=$(shell pwd)/../../../../
export GOPATH

.PHONY : clean

clean :
	rm -Rf $(GOPATH)bin/*
	rm -Rf bin

build :
	go install

docker-clean:
	docker rm blackflowhub
	docker rmi alivinco/blackflowhub

dist-docker :
	mkdir -p bin
	env GOOS=linux GOARCH=amd64 go build -o bin/blackflowhub
	docker build -t alivinco/blackflowhub .

docker-publish : dist-docker
	docker push alivinco/blackflowhub

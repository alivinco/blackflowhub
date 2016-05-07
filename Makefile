GOPATH=$(shell pwd)/../../../..
export GOPATH

.PHONY : clean

clean : ;
	rm $(GOPATH)/bin/* -Rf

build :
	go install

dist-docker : build
	docker rmi alivinco/blackflowhub
	docker build -t alivinco/blackflowhub .

docker-publish : dist-docker
	docker push alivinco/blackfly

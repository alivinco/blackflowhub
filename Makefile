GOPATH=./../../../

.PHONY: clean
clean:
    rm $(GOPATH)/bin/* -Rf

build:
    go install

dist-docker:
    docker rmi alivinco/blackflowhub

docker-publish:
    docker push alivinco/blackfly

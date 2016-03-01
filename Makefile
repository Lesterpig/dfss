REVISION := $(shell git rev-parse HEAD || echo )
VERSION := $(shell git tag --points-at HEAD | grep -m1 v[0-9] | sed -e 's/^v//g' )
ifeq ($(VERSION),)
	VERSION := master
endif

.PHONY:

release: clean build_all package

clean:
	rm -rf release

build_all:
	go get github.com/mitchellh/gox
	gox -os "linux darwin windows" -parallel 1 -output "release/dfss_${VERSION}_{{.OS}}_{{.Arch}}/{{.Dir}}" dfss/dfssc dfss/dfssd dfss/dfssp dfss/dfsst

package:
	echo "$(VERSION) $(REVISION)" > build/embed/VERSION
	cd release && ls -1 . | xargs -n1 -I {} cp ../build/embed/* {}/
	cd release && ls -1 . | xargs -n1 -I {} tar zcvf {}.tar.gz {}

deploy:
	mkdir -p /deploy/$(VERSION)
	cp release/*.tar.gz /deploy/$(VERSION)/

protobuf:
	cd .. && \
	protoc --go_out=plugins=grpc:. dfss/dfssc/api/client.proto && \
	protoc --go_out=plugins=grpc:. dfss/dfssd/api/demonstrator.proto && \
	protoc --go_out=plugins=grpc:. dfss/dfssp/api/platform.proto
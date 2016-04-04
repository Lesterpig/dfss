REVISION := $(shell git rev-parse HEAD || echo )
VERSION := $(shell git tag --points-at HEAD | grep -m1 v[0-9] | sed -e 's/^v//g' )
ifeq ($(VERSION),)
	VERSION := master
endif

.PHONY:

install: nocache
	go install .
	go install ./dfssc
	go install ./dfssd
	go install ./dfssp
	go install ./dfsst

release: clean build_all package

clean:
	rm -rf release

# GUI Build (Docker required)

# prepare_gui builds a new container from the goqt image, adding DFSS dependencies for faster builds.
# call it once or after dependency addition.
prepare_gui: nocache
	docker run --name dfss-builder -v ${PWD}:/go/src/dfss -w /go/src/dfss lesterpig/goqt /bin/bash -c \
		"cp -r ../github.com/visualfc/goqt/bin . ; ./build/deps.sh"
	docker commit dfss-builder dfss:builder
	docker rm dfss-builder

# gui builds the gui component of the dfss project into a docker container, outputing the result in bin/ directory.
gui: nocache
	docker run --rm -v ${PWD}:/go/src/dfss -w /go/src/dfss/gui dfss:builder \
		../bin/goqt_rcc -go main -o application.qrc.go application.qrc
	docker run --rm -v ${PWD}:/go/src/dfss -w /go/src/dfss/gui dfss:builder \
		go build -ldflags "-r ." -o ../bin/gui

# Release internals

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
	protoc --go_out=plugins=grpc:. dfss/dfssp/api/platform.proto && \
	protoc --go_out=plugins=grpc:. dfss/dfsst/api/resolution.proto

nocache:

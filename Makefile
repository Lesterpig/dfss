REVISION := $(shell git rev-parse HEAD || echo )
VERSION := $(shell git tag --points-at HEAD | grep -m1 v[0-9] | sed -e 's/^v//g' )
ifeq ($(VERSION),)
	VERSION := master
endif

.PHONY:

install: nocache
	go install ./dfssc
	go install ./dfssp
	go install ./dfsst

# install_all installs everything, including libraries. It's mandatory for linter, but should be improved in the future.
install_all: install
	git stash
	rm -rf gui
	rm -rf dfssd/cmd
	rm -rf dfssd/gui
	rm dfssd/main.go
	go install ./...
	git reset --hard


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

# dfssd builds the demonstrator into a docker container, outputing the result in bin/ directory
dfssd: nocache
	docker run --rm -v ${PWD}:/go/src/dfss -w /go/src/dfss/dfssd/gui dfss:builder \
		../../bin/goqt_rcc -go gui -o application.qrc.go application.qrc
	docker run --rm -v ${PWD}:/go/src/dfss -w /go/src/dfss/dfssd dfss:builder \
		go build -ldflags "-r ." -o ../bin/dfssd

protobuf:
	cd .. && \
	protoc --go_out=plugins=grpc:. dfss/dfssc/api/client.proto && \
	protoc --go_out=plugins=grpc:. dfss/dfssd/api/demonstrator.proto && \
	protoc --go_out=plugins=grpc:. dfss/dfssp/api/platform.proto && \
	protoc --go_out=plugins=grpc:. dfss/dfsst/api/resolution.proto && \
	protoc --go_out=plugins=grpc:. dfss/net/fixtures/test.proto

# Release internals
# Do not run these commands on your personal computer
release: clean build_x build_g package

build_x:
	go get github.com/mitchellh/gox
	gox -ldflags "-s -w" -osarch "linux/amd64 linux/386 linux/arm windows/386 darwin/amd64" -parallel 1 -output "release/dfss_${VERSION}_{{.OS}}_{{.Arch}}/{{.Dir}}" dfss/dfssc dfss/dfssp dfss/dfsst

build_g:
	cp $(GOPATH)/src/github.com/visualfc/goqt/bin/goqt_rcc /bin/
	cp $(GOPATH)/src/github.com/visualfc/goqt/bin/lib* /lib/
	cd gui && goqt_rcc -go main -o a.qrc.go application.qrc
	cd dfssd/gui && goqt_rcc -go gui -o a.qrc.go application.qrc
	cd gui && go build -ldflags "-s -w -r ." -o ../release/dfss_${VERSION}_linux_amd64/dfssc_gui
	cd dfssd && go build -ldflags "-s -w -r ." -o ../release/dfss_${VERSION}_linux_amd64/dfssd
	cp /lib/libqtdrv.ui.so.1.0.0 release/dfss_${VERSION}_linux_amd64/libqtdrv.ui.so.1

package:
	echo "$(VERSION) $(REVISION)" > build/embed/VERSION
	cd release && ls -1 . | xargs -n1 -I {} cp ../build/embed/* {}/
	cd release && ls -1 . | xargs -n1 -I {} tar zcvf {}.tar.gz {}

deploy:
	mkdir -p /deploy/$(VERSION)
	cp release/*.tar.gz /deploy/$(VERSION)/

clean:
	rm -rf release

nocache:

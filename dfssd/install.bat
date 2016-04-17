set GOARCH=386
set CGO_ENABLED=1
cd gui
..\..\..\github.com\visualfc\goqt\bin\goqt_rcc.exe -go gui -o application.qrc.go application.qrc
cd ..
go build -v -ldflags "-H windowsgui" -o ..\bin\dfssd.exe
set GOARCH=386
set CGO_ENABLED=1
..\..\github.com\visualfc\goqt\bin\goqt_rcc.exe -go main -o application.qrc.go application.qrc
go build -v -ldflags "-H windowsgui" -o ..\bin\gui.exe
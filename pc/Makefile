.PHONY:dev prod sys
dev:sys
	GOOS=windows GOARCH=386 go build
prod:sys
	GOOS=windows GOARCH=386 go build -ldflags="-H windowsgui"
sys:
	rsrc -manifest main.manifest -ico caller.ico -o rsrc.syso

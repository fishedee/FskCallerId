.PHONY:dev prod sys
dev:sys
	GOOS=windows GOARCH=386 go build -ldflags="-H windowsgui"
prod:sys
	GOOS=windows GOARCH=386 go build -ldflags="-H windowsgui -w"
sys:
	rsrc -manifest=main.manifest -ico=caller.ico,head.ico -o rsrc.syso

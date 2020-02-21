foo:
	@cat Makefile.help
	@exit 1

pi3:
	GOOS=linux GOARCH=arm GOARM=7 go build smon.go

here:
	go build smon.go
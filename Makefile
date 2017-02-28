build:
	go build -o bin/alexaops

updateassets:
	go get -u github.com/bouk/staticfiles
	staticfiles -o files/files.go static/

crosscompile:
	GOOS=linux GOARCH=amd64 go build -o bin/alexaops_linux_amd64
	GOOS=linux GOARCH=arm GOARM=5 go build -o bin/alexaops_linux_arm_5
	GOOS=linux GOARCH=arm GOARM=6 go build -o bin/alexaops_linux_arm_6
	GOOS=linux GOARCH=arm GOARM=7 go build -o bin/alexaops_linux_arm_7
	GOOS=darwin GOARCH=amd64 go build -o bin/alexaops_darwin_amd64
	GOOS=windows GOARCH=amd64 go build -o bin/alexaops_windows_amd64

docker:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/alexaops
	docker build -t andreaskoch/alexa-ops:latest .

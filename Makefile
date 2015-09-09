main:
	PATH="${PATH}:${HOME}/Documents/src/go/bin" go generate bitbucket.org/yargevad/go-sparkpost/config
	go build

default: douceur

.PHONY: douceur
douceur:
	CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' .

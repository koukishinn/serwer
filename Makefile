CC=go
EXEC=serwer

all: build

build:
	$(CC) build -o $(EXEC) cmd/main.go

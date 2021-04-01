TARGET = docbase
REV_PARSE = git rev-parse --short HEAD

$(TARGET) : cmd/$(TARGET)/main.go
	@ go build -ldflags "-X main.githash=`$(REV_PARSE)`" -o ./bin/$(TARGET) ./cmd/$(TARGET)

clean :
	@ rm -f ./bin/$(TARGET)

test : 
	@ go test ./...

install :
	@ go install ./cmd/$(TARGET)

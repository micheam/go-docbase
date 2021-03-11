TARGET = docbase

$(TARGET) : cmd/$(TARGET)/main.go
	go build -o ./bin/$(TARGET) ./cmd/$(TARGET)

clean :
	rm -f ./bin/$(TARGET)

test : 
	go test ./...

install :
	go install ./cmd/$(TARGET)

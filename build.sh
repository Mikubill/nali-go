# Windows
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-w -s -X main.githash=$(git rev-parse HEAD)" . && mv main.exe nali.exe && zip nali-go-windows-x64.zip nali.exe && rm nali.exe
CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -ldflags "-w -s -X main.githash=$(git rev-parse HEAD)" . && mv main.exe nali.exe && zip nali-go-windows-x86.zip nali.exe && rm nali.exe
CGO_ENABLED=0 GOOS=windows GOARCH=arm go build -ldflags "-w -s -X main.githash=$(git rev-parse HEAD)" . && mv main.exe nali.exe && zip nali-go-windows-arm.zip nali.exe && rm nali.exe

# Linux
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w -s -X main.githash=$(git rev-parse HEAD)" . && mv main nali && zip nali-go-linux-amd64.zip nali && rm nali
CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -ldflags "-w -s -X main.githash=$(git rev-parse HEAD)" . && mv main nali && zip nali-go-linux-386.zip nali && rm nali
CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -ldflags "-w -s -X main.githash=$(git rev-parse HEAD)" . && mv main nali && zip nali-go-linux-arm.zip nali && rm nali
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags "-w -s -X main.githash=$(git rev-parse HEAD)" . && mv main nali && zip nali-go-linux-arm64.zip nali && rm nali

# Freebsd
CGO_ENABLED=0 GOOS=freebsd GOARCH=amd64 go build -ldflags "-w -s -X main.githash=$(git rev-parse HEAD)" . && mv main nali && zip nali-go-freebsd-amd64.zip nali && rm nali
CGO_ENABLED=0 GOOS=freebsd GOARCH=386 go build -ldflags "-w -s -X main.githash=$(git rev-parse HEAD)" . && mv main nali && zip nali-go-freebsd-386.zip nali && rm nali

# Darwin
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "-w -s -X main.githash=$(git rev-parse HEAD)" . && mv main nali && zip nali-go-darwin-amd64.zip nali && rm nali
CGO_ENABLED=0 GOOS=darwin GOARCH=386 go build -ldflags "-w -s -X main.githash=$(git rev-parse HEAD)" . && mv main nali && zip nali-go-darwin-386.zip nali && rm nali

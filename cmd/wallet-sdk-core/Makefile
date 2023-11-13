.PHONY: ios android clean

BIN_DIR=../../bin/

clean:
	env GO111MODULE=on go clean -cache
	gomobile clean
	rm -fr build

ios:
	go get golang.org/x/mobile
	go mod download golang.org/x/exp
	GOARCH=arm64 gomobile bind -v -trimpath -ldflags "-s -w" -o build/WalletSdkCore.xcframework -target=ios ./cold-wallet/ ./hot-wallet/

android:
	go get golang.org/x/mobile
	go mod download golang.org/x/exp
	gomobile bind -v -trimpath -ldflags "-s -w" -o build/WalletSdkCore.aar -target=android ./cold-wallet/ ./hot-wallet/
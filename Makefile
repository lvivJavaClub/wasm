build-wasm:
	GOOS=js GOARCH=wasm go build -o main.wasm wasm.go

build-proxy-wasm:
	tinygo build -o plugin.wasm -scheduler=none -target=wasi ./proxyfilter.go

envoy: build-proxy-wasm
	envoy -c envoy-wasm.yaml

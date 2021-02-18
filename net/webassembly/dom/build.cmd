set GOOS=js
set GOARCH=wasm
go build -o hello.wasm hello.go
copy %GOROOT%\misc\wasm\wasm_exec.js .
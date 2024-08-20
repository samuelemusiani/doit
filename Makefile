all: go 

go: vue
	cp -r dist cmd/front
	go build -o doit ./cmd

vue:
	npm run build

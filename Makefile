all: go 

go: vue
	cp -r dist cmd/_front
	go build -o doit ./cmd

vue:
	npm install
	npm run build

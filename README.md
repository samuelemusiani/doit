# DOIT

DOIT is a simple todo app written primarily for practicing Go and Vue. Despite 
this, I regularly use it as my todo app.

## Compile

To compile the project simply call make:
```bash
make
```
This will produce a single binary `./doit` that you can execute

## Development

To develop the project you can start the backend and frontend separately:

To start the backend simply run go:
```bash
go run ./cmd
```

To start the frontend first you need to install all the npm dependencies, and 
then you can start the vite server:
```bash
npm install
npm run dev
```

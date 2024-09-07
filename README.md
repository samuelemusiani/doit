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

## Deploy

To deploy DOIT simply run the binary. If you need to adjust the configuration 
you can put a `config.toml` file in the same directory as the binary. There is
an example of the config in the root project directory.

### First user and password

DOIT need a **first user**. If no config is provided his username will be 
`admin`, if a config is provided the username is specified in the config. The 
password is random and printed on the console only the first time DOIT start. 
You should change the default password as soon as possible by logging in.
Note that the password is printed as INFO, if you specify WARN or ERROR as 
`log_level` the password will not be printed.

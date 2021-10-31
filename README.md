# SRV database MS

## How to run the server

1. Install `go` and `docker` in you local device
2. Once you have installed `go` check the `GOPATH` of your machine by using the `go env` commmand in the terminal
3. Navigate to the `GOPATH`
4. In `GOPATH` create a new directories `src`, `github.com`, and `SchulichRacingElectrical`. You should have this folder structure

```
GOPATH
└───src
    └───github.com
        └───SchulichRacingElectrical
```

5. Inside `SchulichRacingElectrical` directory, clone this repository
6. Navigate to `srv-database-ms`
7. Add the `firebase_config.json` folder in the `config` directory
8. Run the command `docker compose up` to start the server inside a docker container
   - The server will be using port `8080` so make sure that no other processes is using that port
   - Once the server starts, it should display all available endpoints
9. `Crt-d` to kill the container

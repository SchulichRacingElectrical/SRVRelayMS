# specify the base image to use for the application
FROM golang:1.17

#download the watcher code on startup
RUN go get -v github.com/canthefason/go-watcher/cmd/watcher

# specify the working directory with the name of this project
WORKDIR /go/src/github.com/SchulichRacingElectrical/srv-database-ms

# copy all the files to the container
COPY . .

RUN export GIN_MODE=release 

# run watcher
ENTRYPOINT /go/bin/watcher
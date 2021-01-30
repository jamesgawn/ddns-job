FROM golang:alpine as build
WORKDIR $GOPATH/src/github.com/jamesgawn/ddns-job

# Copy the rest of the project and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -o /ddns-job .

# Reset to scratch to drop all of the above layers and only copy over the final binary
FROM scratch
ENV HOME=/home
COPY --from=build /ddns-job /bin/ddns-job

CMD ["/bin/ddns-job"]

EXPOSE 8080
FROM golang:1.17 as build
WORKDIR /go/src/github.com/shawntoffel/convector
COPY Makefile go.mod go.sum ./
RUN make deps
ADD . .
RUN make build-linux
RUN echo "convector:x:100:101:/" > passwd

FROM scratch
COPY --from=build /go/src/github.com/shawntoffel/convector/passwd /etc/passwd
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build --chown=100:101 /go/src/github.com/shawntoffel/convector/bin/convector .
USER convector
ENTRYPOINT ["./convector"]

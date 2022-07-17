FROM golang:1.18 as builder

ARG LD_FLAGS="-s -w"
ARG TARGETPLATFORM

WORKDIR /app
COPY . .

RUN go get -d -v \
    && go install -v

RUN export GOOS=$(echo ${TARGETPLATFORM} | cut -d / -f1) && \
    export GOARCH=$(echo ${TARGETPLATFORM} | cut -d / -f2)

RUN go env

RUN CGO_ENABLED=0 go build -ldflags="${LD_FLAGS}" -o /app/build/tracee-polr-adapter -v

FROM scratch
LABEL MAINTAINER "Frank Jogeleit <frank.jogeleit@gweb.de>"

WORKDIR /app

USER 1234

COPY --from=builder /app/LICENSE.md .
COPY --from=builder /app/build/tracee-polr-adapter /app/tracee-polr-adapter

EXPOSE 2112

ENTRYPOINT ["/app/tracee-polr-adapter", "run"]

FROM golang:alpine
WORKDIR /go/src/app
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x
ENV USER=go \
    UID=1000 \
    GID=1000 \
    CGO_ENABLED=0 \
    GOCACHE=/root/.cache/go-build

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=cache,target="/root/.cache/go-build" \
    go mod tidy && \
    go build -ldflags="-s -w" \
    -o go-rest-template && \
    addgroup --gid "$GID" "$USER" && \
    adduser \
    --disabled-password \
    --gecos "" \
    --home "$(pwd)" \
    --ingroup "$USER" \
    --no-create-home \
    --uid "$UID" \
    "$USER" && \
    chown "$UID":"$GID" /go/src/app/go-rest-template && \
    mkdir /cache && \
    chown -R "$UID":"$GID" /cache

FROM scratch
COPY --from=0 /etc/passwd /etc/passwd
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=0 --chown=1000:1000 /cache /cache
COPY --from=0 /go/src/app/go-rest-template /
VOLUME /cache
# Persist certs from certmagic
ENV HOME=/cache
USER 1000
EXPOSE 8080/tcp
EXPOSE 80/tcp
EXPOSE 443/tcp
ENTRYPOINT ["/go-rest-template"]
CMD ["server"]

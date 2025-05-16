FROM golang:1.22.9
ENTRYPOINT ["emboxd"]
CMD ["-c", "/config/config.yaml"]

WORKDIR /code
COPY go.mod go.sum ./
RUN go mod download

RUN go install github.com/playwright-community/playwright-go/cmd/playwright
RUN playwright install --with-deps

COPY . .
RUN go install .

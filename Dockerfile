# The base image is a ubuntu image with Homebrew pre-installed.
FROM homebrew/brew:latest as base
RUN brew update-reset
# HOMEBREW_DIR is the directory where Homebrew is installed.
ENV HOMEBREW_DIR=/home/linuxbrew/.linuxbrew

# Install all system-level dependencies using brew
FROM base as base-brew-install
WORKDIR /src

COPY Brewfile Brewfile.lock.json ./
RUN brew bundle

# Create a workspace for the project and assign permissions
FROM base-brew-install as base-workspace
USER root
RUN mkdir -p /src
RUN chmod +rwx /src
RUN git config --global --add safe.directory '/src'
WORKDIR /src

COPY Taskfile.yaml ./

ENV GOOS=linux
ENV GO111MODULE=on
ENV PATH="/root/go/bin:$PATH"

# Install all npm dependencies
FROM base-workspace as node-install
WORKDIR /src
COPY package.json package-lock.json ./

RUN npm ci

# Install all Go dependencies
FROM node-install as go-install
WORKDIR /src
COPY go.mod go.sum ./

RUN go mod download

RUN go install github.com/tylermmorton/tmpl/cmd/tmpl@latest

# Build the Go binary
FROM go-install as go-build
WORKDIR /src
COPY ./main.go ./
COPY ./app/model ./app/model
COPY ./app/routes/ ./app/routes/
COPY ./app/services/ ./app/services/
COPY ./app/styles/ ./app/styles/
COPY ./app/templates/ ./app/templates/

COPY tailwind.config.js ./
RUN task build:css

RUN go env

RUN go generate ./...
RUN go build -v -o ./.build/bin/testmail ./main.go
RUN chmod +x ./.build/bin/testmail

RUN ls -la ./.build
RUN ls -la ./.build/bin

FROM ubuntu:latest as prod
COPY --from=go-build /src/.build/bin/testmail /bin/app

EXPOSE 8080 1025
ENTRYPOINT [ "/bin/app" ]
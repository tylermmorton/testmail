[![Build status](https://github.com/tylermmorton/testmail/actions/workflows/image.yaml/badge.svg?branch=main&event=push)](https://github.com/tylermmorton/testmail/actions/workflows/image.yaml)

# testmail

`testmail` is a drop-in replacement for your production SMTP server that you can run in your local development environment. It catches all emails sent by your application and allows you to inspect them in a web interface. This is useful for testing end-to-end flows that involve sending and opening emails. To that end, testmail also provides crawler interfaces for [go-rod]() and Cypress to help you implement automated email tests quickly.

Why? Sending test emails from your domain to a real inbox may affect your deliverability and cause your emails to be marked as spam. This is especially true if you are sending a large volume of emails during development and testing. Using a local SMTP server allows you to test your email flows without affecting your production email reputation.

## Features
- [x] SMTP server that catches all emails and stores them in an embedded mongodb collection
- [x] Web interface for viewing and managing emails received by testmail
- [x] Docker image for easy deployment and integration into your testing stack
- [ ] go-rod & cypress integrations for driving automated testing

## Secondary Objectives

`testmail` is not only a useful development tool, but also an experiment in building a hypermedia-based web application in Go using the [torque](https://lbft.dev) framework. This project is meant to dogfood the framework and provide examples of how to build interactive pages when using torque with [htmx](https://htmx.org/).

> 👋🏻 The project is open to contributions and feedback, so feel free to open an issue or PR if you have any ideas or suggestions!

## Getting Started
The easiest way of getting started with `testmail` is to use the Docker image. You can run the image with the following command:

```bash
docker run -p 8080:8080 -p 1025:1025 ghcr.io/tylermmorton/testmail:main
```

This will start the SMTP server on port `1025` and the web interface on port `8080`. You can then configure your application to send emails to `localhost:1025` and view them in the web interface at `localhost:8080`.

You can also use `docker-compose` to  integrate `testmail`:

```yaml
  testmail:
    container_name: testmail
    image: ghcr.io/tylermmorton/testmail:main
    ports:
      - "8080:8080"
      - "2500:2500"
```

## The Stack
This project was generated from the `create-torque-app` template and is using the following technologies:
- [torque](https://lbft.dev) - Webserver framework
- [htmx](https://htmx.org/) - Frontend framework
- [tmpl](https://github.com/tylermmorton/tmpl) - Go `html/template` compiler and renderer
- [TailwindCSS](https://tailwindcss.com/) - CSS framework
- [Docker](https://www.docker.com/) - Container runtime
- [eslint](https://eslint.org/) - JavaScript & HTML linter
- [prettier](https://prettier.io/) - JavaScript & HTML formatter
- [Taskfile](https://taskfile.dev/) - Task runner & mini build system

# testmail

`testmail` is a drop-in replacement for your production SMTP server that you can run in your local development environment. It catches all emails sent by your application and allows you to inspect them in a web interface.

## Features
- [x] SMTP server that catches all emails and stores them in a mongodb collection
- [ ] Web interface for viewing emails
- [ ] Docker image for easy deployment
- [ ] go-rod & cypress integrations for driving automated testing

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

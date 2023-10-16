# Lenslocked

This is a web application written in Go that I built while taking the [Web Development with Go](https://www.usegolang.com/) course.

## Getting Started

### Prerequisites

You will need to have Go installed on your machine. You can find instructions on how to do that [here](https://golang.org/doc/install).

### Installing

To install the application run the following command in your terminal:

```
go install github.com/terrorsquad/lenslocked@latest
```

### Running the application

To run the application run the following command in your terminal:

```
lenslocked
```


### Development

This project uses [modd](https://github.com/cortesi/modd) to run the application and automatically reload it when changes
are made. To run the application in development mode run the following command in your terminal:

```bash
modd
```

Note: TailwindCSS is used for styling and the [standalone CLI](https://tailwindcss.com/blog/standalone-cli#get-started) build is used. There is a helper script `download_tailwind.sh`
that will download the latest version of TailwindCSS and place it in the root of the project. This script is run automatically
when the application is built.

### Built With
- [Go](https://golang.org/) - The programming language used
- [TailwindCSS](https://tailwindcss.com/) - The CSS framework used
- [modd](https://github.com/cortesi/modd) - The tool used to run the application in development mode
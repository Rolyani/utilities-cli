# utilities-cli

A terminal-based utility for quickly transforming text files. Inspired by Programmer's File Editor (pfe) 
https://www.lancaster.ac.uk/~steveb/cpaap/pfe/ 


`utilities-cli` is designed to make repetitive text cleanup and formatting tasks easier through a guided interactive interface. It is built in Go using Bubble Tea and the Charm ecosystem, with a focus on being approachable for non-technical users while still useful for developers and power users.

## Current features

* Add a comma to the end of every line
* Add text to the beginning of every line
* Add text to the end of every line
* Split text after spaces, commas, or both
* Preview input and output before saving
* Save as a new file or overwrite the original
* Automatic backup when overwriting
* Choose files from a default `files/` directory
* Enter a custom file path manually

## Project goals

The project aims to be:

* Simple to use
* Safe by default
* Easy to extend
* Portable as a single binary
* Friendly to contributors

## Built with

* [Go](https://go.dev/)
* [Bubble Tea](https://github.com/charmbracelet/bubbletea)
* [Bubbles](https://github.com/charmbracelet/bubbles)
* [Lip Gloss](https://github.com/charmbracelet/lipgloss)

## Getting started

### Requirements

* Go 1.25+ recommended

### Run locally

```bash
go run .
```

### Build locally

```bash
go build -o utilities-cli
```

On Windows:

```powershell
go build -o utilities-cli.exe
```

## How to use

When the program starts, you can choose one of two ways to load a file:

* Select a file from the default `files/` directory
* Enter a custom file path

Once a file is loaded, choose a transformation, preview the result, then save the output.

### Default files directory

The application supports a default working directory named:

```text
files/
```

If it does not exist, the program will create it.

This is intended to give non-technical users a simple "drop files here" workflow.

## Example workflow

1. Start the app
2. Choose a file source
3. Load a file
4. Select a text transformation
5. Preview the result
6. Save as a new file or overwrite the original

## Project structure

```text
.
├── main.go
├── model.go
├── model_init.go
├── update.go
├── view.go
├── io.go
├── main_test.go
└── transform
    ├── transform.go
```

### File overview

* `main.go` — program entrypoint
* `model.go` — shared model and state definitions
* `model_init.go` — initial model construction and item definitions
* `update.go` — Bubble Tea update logic
* `view.go` — Bubble Tea rendering
* `io.go` — file and directory helpers
* `transform/transrm.go` — pure text transformation logic

## Testing

Run all tests with:

```bash
go test ./...
```


## Contributing

Contributions are welcome.

If you want to contribute:

1. Fork the repository
2. Create a feature branch
3. Make focused changes
4. Add or update tests where appropriate
5. Open a pull request with a clear description

### Contribution guidelines

Please try to keep changes:

* small and well-scoped
* consistent with the existing code style
* easy to review
* backed by tests where possible

For larger changes, open an issue or discussion first.


## Status

This project is currently in active development.

v0.9.0 is for beta tests on platforms. 


## License

MIT

If you are contributing and something is unclear, open an issue or start a discussion.


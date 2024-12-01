# burrow

Quick reloading Go scripts made easy.

## Running

See [releases](https://github.com/KittenzExe/burrow/releases) for downloads

### Windows

On windows, place the executable in the same directory as your project, and run:

``` bash
./burrow.exe app.go
```

Or, if you want to reload when any Go files are modified, run:

``` bash
./burrow.exe --all app.go
```

### MacOS and Linux

On macOS or Linux, place the executable in the same directory as your project, and run:

``` bash
./burrow app.go
```

Or, if you want to reload when any Go files are modified, run:

``` bash
./burrow --all app.go
```

## Building from source

### Windows

To build for Windows, run:

``` bash
go build -o burrow.exe main.go
```

### MacOS and Linux

To build for MacOS and Linux, run:

``` bash
go build -o burrow main.go
```

## License

Under the [GNU Affero General Public License v3.0](https://github.com/KittenzExe/burrow?tab=AGPL-3.0-1-ov-file)

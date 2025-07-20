# goxorro

A simple Go command line application for compressing files.

## Features

- Copy files from source to destination
- Check if source file exists before copying
- Simple command line interface with `-c` flag
- Self-contained binary build

## Building

To build the application, run the build script:

```bash
./build.sh
```

This will create a self-contained binary named `goxorro`.

## Usage

```bash
./goxorro [-c] <source> <destination>
```

### Options

- `-c`: Compress source file to destination file (default behavior)

### Examples

Copy a file:
```bash
./goxorro source.txt destination.txt
```

Copy a file with explicit `-c` flag:
```bash
./goxorro -c source.txt destination.txt
```

## Error Handling

The application will:
- Check if the source file exists before attempting to compress
- Display appropriate error messages if the source file doesn't exist
- Exit with error code 1 on failure
- Exit with error code 0 on success

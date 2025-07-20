# goxorro

A high-performance Go command line application for XOR-based compression using prime number patterns.

## Features

- **XOR Pattern Compression**: Uses 1029 precomputed prime numbers (2-8209) to find optimal bit patterns
- **Intelligent Prime Selection**: Tests primes to achieve maximum zeros (95% threshold for early termination)
- **Gzip Final Compression**: Applies gzip to the optimized bit streams for maximum efficiency
- **Compact 11-bit Metadata**: Efficient prime index encoding using bit packing
- **Lossless Round-trip**: Perfect data integrity with reversible compression/decompression
- **Verbose Logging**: Detailed algorithm analysis with `-v` flag
- **1KB Chunk Processing**: Processes files in 1024-byte chunks for optimal memory usage

## Algorithm Overview

1. **Bit Stream Analysis**: Converts data to bit streams and analyzes distribution
2. **Conditional Negation**: Negates bits if more ones than zeros initially
3. **Prime XOR Testing**: Tests up to 1029 primes to find patterns that maximize zeros
4. **Sparse Detection**: Stops when bit stream becomes <5% ones (95% zeros)
5. **Gzip Compression**: Final compression of optimized bit streams
6. **Metadata Encoding**: 11-bit packed prime indices for minimal overhead

## Building

To build the application, run the build script:

```bash
./build.sh
```

This will create a self-contained binary named `goxorro`.

## Usage

```bash
./goxorro [-c|-d] [-v] <source> <destination>
```

### Options

- `-c`: Compress source file to destination file (default behavior)
- `-d`: Decompress source file to destination file
- `-v`: Enable verbose logging to debug.log

### Examples

Compress a file:
```bash
./goxorro source.txt compressed.xor
```

Decompress a file:
```bash
./goxorro -d compressed.xor restored.txt
```

Compress with verbose logging:
```bash
./goxorro -v -c source.txt compressed.xor
```

## Performance

- **Startup**: Instant - uses precomputed primes (no generation overhead)
- **Compression Ratio**: Achieves 14:1+ ratios on sparse data patterns
- **Metadata Overhead**: ~31% reduction with 11-bit prime index packing
- **Processing Speed**: ~25ms per 1KB chunk with 1029 prime testing

## Error Handling

The application will:
- Check if the source file exists before attempting compression/decompression
- Verify round-trip data integrity during compression
- Display appropriate error messages for invalid operations
- Exit with error code 1 on failure, 0 on success

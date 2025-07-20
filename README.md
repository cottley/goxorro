# goxorro

A high-performance Go command line application for XOR-based compression using prime number patterns.

## Features

- **XOR Pattern Compression**: Uses 1029 precomputed prime numbers (2-8209) to find optimal bit patterns
- **Configurable Prime Testing**: Control compression speed vs thoroughness with `-primes` flag
- **Intelligent Prime Selection**: Tests primes to achieve maximum zeros (95% threshold for early termination)
- **Gzip Final Compression**: Applies gzip to the optimized bit streams for maximum efficiency
- **Delimited File Format**: Separated metadata and data sections with clear delimiters
- **Compact 11-bit Metadata**: Efficient prime index encoding using bit packing
- **Lossless Round-trip**: Perfect data integrity with reversible compression/decompression
- **Verbose Logging**: Detailed algorithm analysis with `-v` flag
- **1KB Chunk Processing**: Processes files in 1024-byte chunks for optimal memory usage

## Algorithm Overview

1. **Bit Stream Analysis**: Converts data to bit streams and analyzes distribution
2. **Conditional Negation**: Negates bits if more ones than zeros initially
3. **Prime XOR Testing**: Tests configurable number of primes to find patterns that maximize zeros
4. **Iterative Optimization**: Continues XOR operations until no improvement possible
5. **Gzip Compression**: Final compression of entire optimized file
6. **Delimited Storage**: Separates data chunks and metadata with clear section delimiters
7. **Metadata Encoding**: 11-bit packed prime indices with chunk delimiters for minimal overhead

## Building

To build the application, run the build script:

```bash
./build.sh
```

This will create a self-contained binary named `goxorro`.

## Usage

```bash
./goxorro [-c|-d] [-v] [-primes N] <source> <destination>
```

### Options

- `-c`: Compress source file to destination file (default behavior)
- `-d`: Decompress source file to destination file
- `-v`: Enable verbose logging to debug.log
- `-primes N`: Number of primes to test per chunk (default: 50, 0 = test all 1029 primes)

### Examples

Compress a file (default 50 primes):
```bash
./goxorro source.txt compressed.xor
```

Compress with fast processing (10 primes):
```bash
./goxorro -primes 10 -c source.txt compressed.xor
```

Compress with maximum thoroughness (all 1029 primes):
```bash
./goxorro -primes 0 -c source.txt compressed.xor
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
- **Processing Speed**: Configurable from ~0.4ms (10 primes) to ~25ms (all 1029 primes) per 1KB chunk
- **Speed vs Quality**: `-primes 25` offers good balance, `-primes 0` for maximum compression

### Prime Testing Performance

| Prime Limit | Speed (per 1KB) | Use Case |
|-------------|-----------------|----------|
| 10          | ~0.4ms         | Fast processing, acceptable quality |
| 25          | ~1ms           | Balanced speed/quality |
| 50 (default)| ~2ms           | Good compression with reasonable speed |
| 200         | ~8ms           | High quality compression |
| 0 (all 1029)| ~25ms          | Maximum compression quality |

## Error Handling

The application will:
- Check if the source file exists before attempting compression/decompression
- Verify round-trip data integrity during compression
- Display appropriate error messages for invalid operations
- Exit with error code 1 on failure, 0 on success

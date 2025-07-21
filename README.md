# goxorro

A Go command line application for XOR-based compression using prime number patterns with experimental sparse data compression techniques.

## Features

- **XOR Pattern Compression**: Uses 1029 precomputed prime numbers (2-8209) to find optimal bit patterns
- **Configurable Prime Testing**: Control compression speed vs thoroughness with `-primes` flag
- **Intelligent Prime Selection**: Tests primes to achieve maximum zeros (95% threshold for early termination)
- **Adaptive Sparse Compression**: Experimental techniques for post-XOR sparse data compression
- **Multiple File Formats**: Format v1 (basic) and v2 (streaming with experimental compression)
- **Compact 11-bit Metadata**: Efficient prime index encoding using bit packing
- **Lossless Round-trip**: Perfect data integrity with reversible compression/decompression
- **Verbose Logging**: Detailed algorithm analysis with `-v` flag
- **1KB Chunk Processing**: Processes files in 1024-byte chunks for optimal memory usage

## Algorithm Overview

1. **Bit Stream Analysis**: Converts data to bit streams and analyzes distribution
2. **Conditional Negation**: Negates bits if more ones than zeros initially
3. **Prime XOR Testing**: Tests configurable number of primes to find patterns that maximize zeros
4. **Iterative Optimization**: Continues XOR operations until no improvement possible
5. **Sparse Data Compression**: Applies experimental compression to post-XOR sparse data
6. **Streaming Format**: Format v2 stores metadata with each chunk for streaming decompression
7. **Adaptive Compression**: Automatically selects best compression method based on data characteristics

## Experimental Compression Research

This project has explored several approaches to compress sparse binary data remaining after XOR optimization:

### Tested Approaches
- **Combinadic Numbers**: Mathematical representation of bit combinations - effective for very sparse data but computationally expensive
- **Bit Position Encoding**: Stores positions of 1-bits with delta compression - good for sparse data (92% of original size)
- **Run-Length Encoding**: Encodes consecutive runs of 0s and 1s - ineffective for random data
- **Adaptive Compression**: Automatically selects bit positions (<30% density) or RLE (≥30% density)

### Key Findings
- **Random dense data** (~50% ones): No compression method is effective (entropy is already maximal)
- **Sparse data** (<30% ones): Bit position encoding provides ~8% space savings
- **Structured data**: RLE can be effective for patterns with long runs
- **Raw storage** is often most efficient for post-XOR data due to its already-compressed nature

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

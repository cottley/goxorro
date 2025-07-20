package main

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

type CompressionStep struct {
	Operation string
	Prime     int
	Applied   bool
}

var primes []int
var debugLogger *log.Logger
var verboseMode bool

func init() {
	primes = generatePrimes(1029)
}

func generatePrimes(count int) []int {
	var primes []int
	num := 2
	
	for len(primes) < count {
		if isPrime(num) {
			primes = append(primes, num)
		}
		num++
	}
	
	return primes
}

func isPrime(n int) bool {
	if n < 2 {
		return false
	}
	if n == 2 {
		return true
	}
	if n%2 == 0 {
		return false
	}
	
	for i := 3; i*i <= n; i += 2 {
		if n%i == 0 {
			return false
		}
	}
	
	return true
}


func getPrimeIndex(prime int) int {
	for i, p := range primes {
		if p == prime {
			return i
		}
	}
	return -1
}

func main() {
	var compressFlag bool
	var decompressFlag bool
	flag.BoolVar(&compressFlag, "c", false, "Compress source file to destination file")
	flag.BoolVar(&decompressFlag, "d", false, "Decompress source file to destination file")
	flag.BoolVar(&verboseMode, "v", false, "Enable verbose logging to debug.log")
	flag.Parse()

	if verboseMode {
		initDebugLogging()
	}

	args := flag.Args()
	if len(args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s [-c|-d] <source> <destination>\n", os.Args[0])
		os.Exit(1)
	}

	if compressFlag && decompressFlag {
		fmt.Fprintf(os.Stderr, "Error: Cannot use both -c and -d flags\n")
		os.Exit(1)
	}

	sourceFile := args[0]
	destFile := args[1]

	if _, err := os.Stat(sourceFile); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: Source file '%s' does not exist\n", sourceFile)
		os.Exit(1)
	}

	if decompressFlag {
		logDebug("Starting decompression: %s -> %s", sourceFile, destFile)
		if err := decompressFile(sourceFile, destFile); err != nil {
			logDebug("Decompression failed: %v", err)
			fmt.Fprintf(os.Stderr, "Error decompressing file: %v\n", err)
			os.Exit(1)
		}
		logDebug("Decompression completed successfully")
		fmt.Printf("Successfully decompressed '%s' to '%s'\n", sourceFile, destFile)
	} else {
		logDebug("Starting compression: %s -> %s", sourceFile, destFile)
		if err := compressFile(sourceFile, destFile); err != nil {
			logDebug("Compression failed: %v", err)
			fmt.Fprintf(os.Stderr, "Error compressing file: %v\n", err)
			os.Exit(1)
		}
		logDebug("Compression completed successfully")
		fmt.Printf("Successfully compressed '%s' to '%s'\n", sourceFile, destFile)
	}
}

func compressFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	stat, err := sourceFile.Stat()
	if err == nil {
		logDebug("Source file size: %d bytes", stat.Size())
	}

	buffer := make([]byte, 1024)
	var allCompressionSteps [][]CompressionStep
	chunkCount := 0

	for {
		n, err := sourceFile.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		chunk := buffer[:n]
		logDebug("Processing chunk %d: %d bytes", chunkCount, len(chunk))
		
		start := time.Now()
		compressedData, steps := compressChunk(chunk)
		duration := time.Since(start)
		
		logDebug("Chunk %d compressed: %d -> %d bytes (%.3f ratio) in %v", 
			chunkCount, len(chunk), len(compressedData), 
			float64(len(compressedData))/float64(len(chunk)), duration)
		logDebug("Chunk %d compression steps: %d operations", chunkCount, len(steps))
		
		allCompressionSteps = append(allCompressionSteps, steps)

		if err := binary.Write(destFile, binary.LittleEndian, int32(len(compressedData))); err != nil {
			return err
		}
		if _, err := destFile.Write(compressedData); err != nil {
			return err
		}
		
		chunkCount++
	}

	metadataBytes := encodeMetadata(allCompressionSteps)
	logDebug("Metadata encoded: %d bytes for %d chunks", len(metadataBytes), len(allCompressionSteps))

	if _, err := destFile.Write(metadataBytes); err != nil {
		return err
	}
	
	if err := binary.Write(destFile, binary.LittleEndian, int32(len(metadataBytes))); err != nil {
		return err
	}

	finalStat, err := destFile.Stat()
	if err == nil {
		logDebug("Final compressed file size: %d bytes", finalStat.Size())
	}

	return nil
}

func decompressFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	stat, err := sourceFile.Stat()
	if err != nil {
		return err
	}
	logDebug("Compressed file size: %d bytes", stat.Size())

	var metadataSize int32
	if _, err := sourceFile.Seek(-4, io.SeekEnd); err != nil {
		return err
	}
	if err := binary.Read(sourceFile, binary.LittleEndian, &metadataSize); err != nil {
		return err
	}
	logDebug("Metadata size: %d bytes", metadataSize)

	metadataStart := stat.Size() - int64(metadataSize) - 4
	if _, err := sourceFile.Seek(metadataStart, io.SeekStart); err != nil {
		return err
	}

	metadataBytes := make([]byte, metadataSize)
	if _, err := io.ReadFull(sourceFile, metadataBytes); err != nil {
		return err
	}

	allCompressionSteps, err := decodeMetadata(metadataBytes)
	if err != nil {
		return err
	}
	logDebug("Decoded metadata for %d chunks", len(allCompressionSteps))

	if _, err := sourceFile.Seek(0, io.SeekStart); err != nil {
		return err
	}

	chunkIndex := 0
	totalDecompressed := 0
	for {
		pos, _ := sourceFile.Seek(0, io.SeekCurrent)
		if pos >= metadataStart {
			break
		}

		var chunkSize int32
		if err := binary.Read(sourceFile, binary.LittleEndian, &chunkSize); err != nil {
			break
		}

		compressedData := make([]byte, chunkSize)
		if _, err := sourceFile.Read(compressedData); err != nil {
			return err
		}

		if chunkIndex >= len(allCompressionSteps) {
			return fmt.Errorf("chunk index out of range")
		}

		logDebug("Decompressing chunk %d: %d bytes compressed", chunkIndex, chunkSize)
		start := time.Now()
		originalData := decompressChunk(compressedData, allCompressionSteps[chunkIndex])
		duration := time.Since(start)
		
		logDebug("Chunk %d decompressed: %d -> %d bytes in %v", 
			chunkIndex, chunkSize, len(originalData), duration)
		
		if _, err := destFile.Write(originalData); err != nil {
			return err
		}

		totalDecompressed += len(originalData)
		chunkIndex++
	}

	logDebug("Decompression complete: %d chunks, %d total bytes", chunkIndex, totalDecompressed)
	return nil
}

func compressChunk(data []byte) ([]byte, []CompressionStep) {
	var steps []CompressionStep
	bitStream := bytesToBits(data)
	
	ones, zeros := countBits(bitStream)
	logDebug("Initial bit distribution: %d ones, %d zeros (%.1f%% ones)", 
		ones, zeros, float64(ones)/float64(len(bitStream))*100)
	
	if ones > zeros {
		bitStream = negateBits(bitStream)
		steps = append(steps, CompressionStep{Operation: "negate", Applied: true})
		logDebug("Applied negation: %d zeros, %d ones", ones, zeros)
	}

	
	for iteration := 0; iteration < 10; iteration++ {
		logDebug("XOR iteration %d: testing primes for bit stream with %d bits", 
			iteration, len(bitStream))
		
		bestPrime := -1
		bestZeros := 0
		var bestResult []int
		primesChecked := 0
		
		for i, prime := range primes {
			if prime > 1000 && iteration == 0 {
				continue
			}
			
			xorPattern := createXORPattern(len(bitStream), prime)
			testResult := xorBits(bitStream, xorPattern)
			
			_, zeros := countBits(testResult)
			primesChecked++
			
			if zeros > bestZeros {
				bestZeros = zeros
				bestPrime = prime
				bestResult = testResult
				logDebug("New best prime %d: %d zeros (%.1f%%)", 
					prime, zeros, float64(zeros)/float64(len(bitStream))*100)
				
				if zeros > len(bitStream)*85/100 {
					logDebug("Excellent pattern found (>85%% zeros), stopping early")
					break
				}
			}
			
			if i > 100 && bestZeros > len(bitStream)*60/100 {
				logDebug("Good pattern found after 100 primes (>60%% zeros), stopping")
				break
			}
		}
		
		logDebug("Iteration %d complete: checked %d primes, best was %d with %d zeros", 
			iteration, primesChecked, bestPrime, bestZeros)
		
		if bestPrime == -1 || bestZeros <= len(bitStream)/2 {
			logDebug("No improvement found, stopping iterations")
			break
		}
		
		bitStream = bestResult
		steps = append(steps, CompressionStep{Operation: "xor", Prime: bestPrime, Applied: true})
		
		if isSparseBitStream(bitStream) {
			logDebug("Bit stream is now sparse (<%0.1f%% ones), stopping", 
				float64(countOnes(bitStream))/float64(len(bitStream))*100)
			break
		}
	}

	compressed := gzipCompress(bitStream)
	return compressed, steps
}

func bytesToBits(data []byte) []int {
	var bits []int
	for _, b := range data {
		for i := 7; i >= 0; i-- {
			bits = append(bits, int((b>>i)&1))
		}
	}
	return bits
}

func countBits(bits []int) (ones, zeros int) {
	for _, bit := range bits {
		if bit == 1 {
			ones++
		} else {
			zeros++
		}
	}
	return
}

func countOnes(bits []int) int {
	ones := 0
	for _, bit := range bits {
		if bit == 1 {
			ones++
		}
	}
	return ones
}

func negateBits(bits []int) []int {
	result := make([]int, len(bits))
	for i, bit := range bits {
		result[i] = 1 - bit
	}
	return result
}

func createXORPattern(length, prime int) []int {
	pattern := make([]int, length)
	for i := prime - 1; i < length; i += prime {
		pattern[i] = 1
	}
	return pattern
}

func xorBits(a, b []int) []int {
	result := make([]int, len(a))
	for i := range a {
		result[i] = a[i] ^ b[i]
	}
	return result
}

func isSparseBitStream(bits []int) bool {
	ones, _ := countBits(bits)
	return float64(ones)/float64(len(bits)) < 0.1
}


func decompressChunk(compressedData []byte, steps []CompressionStep) []byte {
	bitStream := gzipDecompress(compressedData)
	
	for i := len(steps) - 1; i >= 0; i-- {
		step := steps[i]
		if step.Applied {
			switch step.Operation {
			case "xor":
				xorPattern := createXORPattern(len(bitStream), step.Prime)
				bitStream = xorBits(bitStream, xorPattern)
			case "negate":
				bitStream = negateBits(bitStream)
			}
		}
	}
	
	return bitsToBytes(bitStream)
}


func bitsToBytes(bits []int) []byte {
	var result []byte
	
	for i := 0; i < len(bits); i += 8 {
		var b byte
		for j := 0; j < 8 && i+j < len(bits); j++ {
			if bits[i+j] == 1 {
				b |= 1 << (7 - j)
			}
		}
		result = append(result, b)
	}
	
	return result
}

func encodeMetadata(allSteps [][]CompressionStep) []byte {
	var result []byte
	
	result = append(result, byte(len(allSteps)))
	
	for _, steps := range allSteps {
		var negateApplied byte
		var xorSteps []CompressionStep
		
		for _, step := range steps {
			if step.Operation == "negate" && step.Applied {
				negateApplied = 1
			} else if step.Operation == "xor" && step.Applied {
				xorSteps = append(xorSteps, step)
			}
		}
		
		result = append(result, negateApplied, byte(len(xorSteps)))
		
		for _, step := range xorSteps {
			primeIndex := getPrimeIndex(step.Prime)
			if primeIndex >= 0 && primeIndex < len(primes) {
				indexBytes := make([]byte, 2)
				binary.LittleEndian.PutUint16(indexBytes, uint16(primeIndex))
				result = append(result, indexBytes...)
			}
		}
	}
	
	return result
}

func decodeMetadata(data []byte) ([][]CompressionStep, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty metadata")
	}
	
	numChunks := int(data[0])
	var allSteps [][]CompressionStep
	
	offset := 1
	for i := 0; i < numChunks; i++ {
		if offset+1 >= len(data) {
			return nil, fmt.Errorf("invalid metadata: insufficient data")
		}
		
		negateApplied := data[offset]
		numXorSteps := int(data[offset+1])
		offset += 2
		
		if offset+numXorSteps*2 > len(data) {
			return nil, fmt.Errorf("invalid metadata: insufficient XOR step data")
		}
		
		var steps []CompressionStep
		
		if negateApplied == 1 {
			steps = append(steps, CompressionStep{Operation: "negate", Applied: true})
		}
		
		for j := 0; j < numXorSteps; j++ {
			primeIndex := int(binary.LittleEndian.Uint16(data[offset+j*2 : offset+j*2+2]))
			if primeIndex < len(primes) {
				steps = append(steps, CompressionStep{
					Operation: "xor",
					Prime:     primes[primeIndex],
					Applied:   true,
				})
			}
		}
		
		offset += numXorSteps * 2
		allSteps = append(allSteps, steps)
	}
	
	return allSteps, nil
}


func gzipCompress(bits []int) []byte {
	bitBytes := bitsToBytes(bits)
	
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	
	lengthBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(lengthBytes, uint32(len(bits)))
	gz.Write(lengthBytes)
	gz.Write(bitBytes)
	gz.Close()
	
	return buf.Bytes()
}

func gzipDecompress(data []byte) []int {
	buf := bytes.NewReader(data)
	gz, err := gzip.NewReader(buf)
	if err != nil {
		return []int{}
	}
	defer gz.Close()
	
	lengthBytes := make([]byte, 4)
	if _, err := io.ReadFull(gz, lengthBytes); err != nil {
		return []int{}
	}
	originalBitLength := int(binary.LittleEndian.Uint32(lengthBytes))
	
	compressedBytes, err := io.ReadAll(gz)
	if err != nil {
		return []int{}
	}
	
	allBits := bytesToBits(compressedBytes)
	if len(allBits) >= originalBitLength {
		return allBits[:originalBitLength]
	}
	
	return allBits
}

func initDebugLogging() {
	logFile, err := os.OpenFile("debug.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Could not create debug.log: %v\n", err)
		return
	}
	
	debugLogger = log.New(logFile, "", log.LstdFlags|log.Lmicroseconds)
	debugLogger.Printf("=== goxorro debug session started ===")
	debugLogger.Printf("Command: %v", os.Args)
	debugLogger.Printf("Generated %d prime numbers (2 to %d)", len(primes), primes[len(primes)-1])
}

func logDebug(format string, args ...interface{}) {
	if verboseMode && debugLogger != nil {
		debugLogger.Printf(format, args...)
	}
}
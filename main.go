package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"os"
)

type CompressionStep struct {
	Operation string
	Prime     int
	Applied   bool
}

func main() {
	var compressFlag bool
	flag.BoolVar(&compressFlag, "c", false, "Compress source file to destination file")
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s [-c] <source> <destination>\n", os.Args[0])
		os.Exit(1)
	}

	sourceFile := args[0]
	destFile := args[1]

	if _, err := os.Stat(sourceFile); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: Source file '%s' does not exist\n", sourceFile)
		os.Exit(1)
	}

	if compressFlag {
		if err := compressFile(sourceFile, destFile); err != nil {
			fmt.Fprintf(os.Stderr, "Error compressing file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Successfully compressed '%s' to '%s'\n", sourceFile, destFile)
	} else {
		if err := compressFile(sourceFile, destFile); err != nil {
			fmt.Fprintf(os.Stderr, "Error compressing file: %v\n", err)
			os.Exit(1)
		}
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

	buffer := make([]byte, 1024)
	var compressionSteps []CompressionStep

	for {
		n, err := sourceFile.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		chunk := buffer[:n]
		compressedData, steps := compressChunk(chunk)
		compressionSteps = append(compressionSteps, steps...)

		if err := binary.Write(destFile, binary.LittleEndian, int32(len(compressedData))); err != nil {
			return err
		}
		if _, err := destFile.Write(compressedData); err != nil {
			return err
		}
	}

	return nil
}

func compressChunk(data []byte) ([]byte, []CompressionStep) {
	var steps []CompressionStep
	bitStream := bytesToBits(data)
	
	ones, zeros := countBits(bitStream)
	if ones > zeros {
		bitStream = negateBits(bitStream)
		steps = append(steps, CompressionStep{Operation: "negate", Applied: true})
	}

	primes := []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31}
	
	for _, prime := range primes {
		xorPattern := createXORPattern(len(bitStream), prime)
		testResult := xorBits(bitStream, xorPattern)
		
		testOnes, testZeros := countBits(testResult)
		if testZeros > testOnes {
			bitStream = testResult
			steps = append(steps, CompressionStep{Operation: "xor", Prime: prime, Applied: true})
		}
		
		if isSparseBitStream(bitStream) {
			break
		}
	}

	compressed := runLengthEncode(bitStream)
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

func runLengthEncode(bits []int) []byte {
	if len(bits) == 0 {
		return []byte{}
	}
	
	var result []byte
	currentBit := bits[0]
	count := 1
	
	for i := 1; i < len(bits); i++ {
		if bits[i] == currentBit && count < 255 {
			count++
		} else {
			result = append(result, byte(currentBit), byte(count))
			currentBit = bits[i]
			count = 1
		}
	}
	
	result = append(result, byte(currentBit), byte(count))
	return result
}
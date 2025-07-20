package main

import (
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

var primes = []int{
	2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 
	31, 37, 41, 43, 47, 53, 59, 61, 67, 71, 
	73, 79, 83, 89, 97, 101, 103, 107, 109, 113, 
	127, 131, 137, 139, 149, 151, 157, 163, 167, 173, 
	179, 181, 191, 193, 197, 199, 211, 223, 227, 229, 
	233, 239, 241, 251, 257, 263, 269, 271, 277, 281, 
	283, 293, 307, 311, 313, 317, 331, 337, 347, 349, 
	353, 359, 367, 373, 379, 383, 389, 397, 401, 409, 
	419, 421, 431, 433, 439, 443, 449, 457, 461, 463, 
	467, 479, 487, 491, 499, 503, 509, 521, 523, 541, 
	547, 557, 563, 569, 571, 577, 587, 593, 599, 601, 
	607, 613, 617, 619, 631, 641, 643, 647, 653, 659, 
	661, 673, 677, 683, 691, 701, 709, 719, 727, 733, 
	739, 743, 751, 757, 761, 769, 773, 787, 797, 809, 
	811, 821, 823, 827, 829, 839, 853, 857, 859, 863, 
	877, 881, 883, 887, 907, 911, 919, 929, 937, 941, 
	947, 953, 967, 971, 977, 983, 991, 997, 1009, 1013, 
	1019, 1021, 1031, 1033, 1039, 1049, 1051, 1061, 1063, 1069, 
	1087, 1091, 1093, 1097, 1103, 1109, 1117, 1123, 1129, 1151, 
	1153, 1163, 1171, 1181, 1187, 1193, 1201, 1213, 1217, 1223, 
	1229, 1231, 1237, 1249, 1259, 1277, 1279, 1283, 1289, 1291, 
	1297, 1301, 1303, 1307, 1319, 1321, 1327, 1361, 1367, 1373, 
	1381, 1399, 1409, 1423, 1427, 1429, 1433, 1439, 1447, 1451, 
	1453, 1459, 1471, 1481, 1483, 1487, 1489, 1493, 1499, 1511, 
	1523, 1531, 1543, 1549, 1553, 1559, 1567, 1571, 1579, 1583, 
	1597, 1601, 1607, 1609, 1613, 1619, 1621, 1627, 1637, 1657, 
	1663, 1667, 1669, 1693, 1697, 1699, 1709, 1721, 1723, 1733, 
	1741, 1747, 1753, 1759, 1777, 1783, 1787, 1789, 1801, 1811, 
	1823, 1831, 1847, 1861, 1867, 1871, 1873, 1877, 1879, 1889, 
	1901, 1907, 1913, 1931, 1933, 1949, 1951, 1973, 1979, 1987, 
	1993, 1997, 1999, 2003, 2011, 2017, 2027, 2029, 2039, 2053, 
	2063, 2069, 2081, 2083, 2087, 2089, 2099, 2111, 2113, 2129, 
	2131, 2137, 2141, 2143, 2153, 2161, 2179, 2203, 2207, 2213, 
	2221, 2237, 2239, 2243, 2251, 2267, 2269, 2273, 2281, 2287, 
	2293, 2297, 2309, 2311, 2333, 2339, 2341, 2347, 2351, 2357, 
	2371, 2377, 2381, 2383, 2389, 2393, 2399, 2411, 2417, 2423, 
	2437, 2441, 2447, 2459, 2467, 2473, 2477, 2503, 2521, 2531, 
	2539, 2543, 2549, 2551, 2557, 2579, 2591, 2593, 2609, 2617, 
	2621, 2633, 2647, 2657, 2659, 2663, 2671, 2677, 2683, 2687, 
	2689, 2693, 2699, 2707, 2711, 2713, 2719, 2729, 2731, 2741, 
	2749, 2753, 2767, 2777, 2789, 2791, 2797, 2801, 2803, 2819, 
	2833, 2837, 2843, 2851, 2857, 2861, 2879, 2887, 2897, 2903, 
	2909, 2917, 2927, 2939, 2953, 2957, 2963, 2969, 2971, 2999, 
	3001, 3011, 3019, 3023, 3037, 3041, 3049, 3061, 3067, 3079, 
	3083, 3089, 3109, 3119, 3121, 3137, 3163, 3167, 3169, 3181, 
	3187, 3191, 3203, 3209, 3217, 3221, 3229, 3251, 3253, 3257, 
	3259, 3271, 3299, 3301, 3307, 3313, 3319, 3323, 3329, 3331, 
	3343, 3347, 3359, 3361, 3371, 3373, 3389, 3391, 3407, 3413, 
	3433, 3449, 3457, 3461, 3463, 3467, 3469, 3491, 3499, 3511, 
	3517, 3527, 3529, 3533, 3539, 3541, 3547, 3557, 3559, 3571, 
	3581, 3583, 3593, 3607, 3613, 3617, 3623, 3631, 3637, 3643, 
	3659, 3671, 3673, 3677, 3691, 3697, 3701, 3709, 3719, 3727, 
	3733, 3739, 3761, 3767, 3769, 3779, 3793, 3797, 3803, 3821, 
	3823, 3833, 3847, 3851, 3853, 3863, 3877, 3881, 3889, 3907, 
	3911, 3917, 3919, 3923, 3929, 3931, 3943, 3947, 3967, 3989, 
	4001, 4003, 4007, 4013, 4019, 4021, 4027, 4049, 4051, 4057, 
	4073, 4079, 4091, 4093, 4099, 4111, 4127, 4129, 4133, 4139, 
	4153, 4157, 4159, 4177, 4201, 4211, 4217, 4219, 4229, 4231, 
	4241, 4243, 4253, 4259, 4261, 4271, 4273, 4283, 4289, 4297, 
	4327, 4337, 4339, 4349, 4357, 4363, 4373, 4391, 4397, 4409, 
	4421, 4423, 4441, 4447, 4451, 4457, 4463, 4481, 4483, 4493, 
	4507, 4513, 4517, 4519, 4523, 4547, 4549, 4561, 4567, 4583, 
	4591, 4597, 4603, 4621, 4637, 4639, 4643, 4649, 4651, 4657, 
	4663, 4673, 4679, 4691, 4703, 4721, 4723, 4729, 4733, 4751, 
	4759, 4783, 4787, 4789, 4793, 4799, 4801, 4813, 4817, 4831, 
	4861, 4871, 4877, 4889, 4903, 4909, 4919, 4931, 4933, 4937, 
	4943, 4951, 4957, 4967, 4969, 4973, 4987, 4993, 4999, 5003, 
	5009, 5011, 5021, 5023, 5039, 5051, 5059, 5077, 5081, 5087, 
	5099, 5101, 5107, 5113, 5119, 5147, 5153, 5167, 5171, 5179, 
	5189, 5197, 5209, 5227, 5231, 5233, 5237, 5261, 5273, 5279, 
	5281, 5297, 5303, 5309, 5323, 5333, 5347, 5351, 5381, 5387, 
	5393, 5399, 5407, 5413, 5417, 5419, 5431, 5437, 5441, 5443, 
	5449, 5471, 5477, 5479, 5483, 5501, 5503, 5507, 5519, 5521, 
	5527, 5531, 5557, 5563, 5569, 5573, 5581, 5591, 5623, 5639, 
	5641, 5647, 5651, 5653, 5657, 5659, 5669, 5683, 5689, 5693, 
	5701, 5711, 5717, 5737, 5741, 5743, 5749, 5779, 5783, 5791, 
	5801, 5807, 5813, 5821, 5827, 5839, 5843, 5849, 5851, 5857, 
	5861, 5867, 5869, 5879, 5881, 5897, 5903, 5923, 5927, 5939, 
	5953, 5981, 5987, 6007, 6011, 6029, 6037, 6043, 6047, 6053, 
	6067, 6073, 6079, 6089, 6091, 6101, 6113, 6121, 6131, 6133, 
	6143, 6151, 6163, 6173, 6197, 6199, 6203, 6211, 6217, 6221, 
	6229, 6247, 6257, 6263, 6269, 6271, 6277, 6287, 6299, 6301, 
	6311, 6317, 6323, 6329, 6337, 6343, 6353, 6359, 6361, 6367, 
	6373, 6379, 6389, 6397, 6421, 6427, 6449, 6451, 6469, 6473, 
	6481, 6491, 6521, 6529, 6547, 6551, 6553, 6563, 6569, 6571, 
	6577, 6581, 6599, 6607, 6619, 6637, 6653, 6659, 6661, 6673, 
	6679, 6689, 6691, 6701, 6703, 6709, 6719, 6733, 6737, 6761, 
	6763, 6779, 6781, 6791, 6793, 6803, 6823, 6827, 6829, 6833, 
	6841, 6857, 6863, 6869, 6871, 6883, 6899, 6907, 6911, 6917, 
	6947, 6949, 6959, 6961, 6967, 6971, 6977, 6983, 6991, 6997, 
	7001, 7013, 7019, 7027, 7039, 7043, 7057, 7069, 7079, 7103, 
	7109, 7121, 7127, 7129, 7151, 7159, 7177, 7187, 7193, 7207, 
	7211, 7213, 7219, 7229, 7237, 7243, 7247, 7253, 7283, 7297, 
	7307, 7309, 7321, 7331, 7333, 7349, 7351, 7369, 7393, 7411, 
	7417, 7433, 7451, 7457, 7459, 7477, 7481, 7487, 7489, 7499, 
	7507, 7517, 7523, 7529, 7537, 7541, 7547, 7549, 7559, 7561, 
	7573, 7577, 7583, 7589, 7591, 7603, 7607, 7621, 7639, 7643, 
	7649, 7669, 7673, 7681, 7687, 7691, 7699, 7703, 7717, 7723, 
	7727, 7741, 7753, 7757, 7759, 7789, 7793, 7817, 7823, 7829, 
	7841, 7853, 7867, 7873, 7877, 7879, 7883, 7901, 7907, 7919, 
	7927, 7933, 7937, 7949, 7951, 7963, 7993, 8009, 8011, 8017, 
	8039, 8053, 8059, 8069, 8081, 8087, 8089, 8093, 8101, 8111, 
	8117, 8123, 8147, 8161, 8167, 8171, 8179, 8191, 8209,
}
var debugLogger *log.Logger
var verboseMode bool


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
	} else {
		logDebug("Starting compression: %s -> %s", sourceFile, destFile)
		if err := compressFile(sourceFile, destFile); err != nil {
			logDebug("Compression failed: %v", err)
			fmt.Fprintf(os.Stderr, "Error compressing file: %v\n", err)
			os.Exit(1)
		}
		logDebug("Compression completed successfully")
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
	if err != nil {
		return err
	}
	fileSize := stat.Size()
	totalChunks := (fileSize + 1023) / 1024
	logDebug("Source file size: %d bytes (%d chunks)", fileSize, totalChunks)

	fmt.Printf("Compressing %d bytes in %d chunks...\n", fileSize, totalChunks)

	buffer := make([]byte, 1024)
	var allCompressionSteps [][]CompressionStep
	chunkCount := 0
	bytesProcessed := int64(0)
	var avgChunkTime time.Duration

	for {
		n, err := sourceFile.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		chunk := buffer[:n]
		bytesProcessed += int64(n)
		progress := float64(bytesProcessed) / float64(fileSize) * 100
		
		start := time.Now()
		compressedData, steps := compressChunk(chunk)
		duration := time.Since(start)
		
		// Calculate time estimates
		if chunkCount == 0 {
			avgChunkTime = duration
		} else {
			// Moving average of chunk processing time
			avgChunkTime = (avgChunkTime*time.Duration(chunkCount) + duration) / time.Duration(chunkCount+1)
		}
		
		remainingChunks := int64(totalChunks) - int64(chunkCount + 1)
		estimatedRemaining := avgChunkTime * time.Duration(remainingChunks)
		
		// Format time estimates
		timeStr := ""
		if chunkCount > 2 { // Show estimates after a few chunks for accuracy
			if estimatedRemaining < time.Minute {
				timeStr = fmt.Sprintf(" ETA: %ds", int(estimatedRemaining.Seconds()))
			} else if estimatedRemaining < time.Hour {
				timeStr = fmt.Sprintf(" ETA: %dm%ds", int(estimatedRemaining.Minutes()), int(estimatedRemaining.Seconds())%60)
			} else {
				timeStr = fmt.Sprintf(" ETA: %dh%dm", int(estimatedRemaining.Hours()), int(estimatedRemaining.Minutes())%60)
			}
		}
		
		fmt.Printf("\rProgress: %.1f%% (%d/%d chunks)%s", progress, chunkCount+1, totalChunks, timeStr)
		logDebug("Processing chunk %d: %d bytes", chunkCount, len(chunk))
		
		logDebug("Chunk %d compressed: %d -> %d bytes (%.3f ratio) in %v", 
			chunkCount, len(chunk), len(compressedData), 
			float64(len(compressedData))/float64(len(chunk)), duration)
		logDebug("Chunk %d compression steps: %d operations", chunkCount, len(steps))
		
		allCompressionSteps = append(allCompressionSteps, steps)

		// Write compressed data with length delimiter
		if err := binary.Write(destFile, binary.LittleEndian, uint32(len(compressedData))); err != nil {
			return err
		}
		if _, err := destFile.Write(compressedData); err != nil {
			return err
		}
		
		chunkCount++
	}

	fmt.Printf("\rProgress: 100.0%% (%d/%d chunks) - Writing metadata...\n", totalChunks, totalChunks)

	// Write delimiter marker to separate data section from metadata section
	delimiter := []byte("GOXORRO_META_START")
	if _, err := destFile.Write(delimiter); err != nil {
		return err
	}

	// Write metadata section
	metadataBytes := encodeMetadataWithDelimiters(allCompressionSteps)
	logDebug("Metadata encoded: %d bytes for %d chunks", len(metadataBytes), len(allCompressionSteps))

	if _, err := destFile.Write(metadataBytes); err != nil {
		return err
	}

	// Close the raw XOR file
	destFile.Close()
	
	// Get raw file size
	rawStat, err := os.Stat(dst)
	if err != nil {
		return err
	}
	rawSize := rawStat.Size()
	logDebug("Raw XOR file size: %d bytes", rawSize)
	fmt.Printf("XOR processing complete: %d -> %d bytes\n", fileSize, rawSize)
	
	// Now gzip the entire file
	fmt.Printf("Applying gzip compression to entire file...\n")
	gzipDst := dst + ".gz"
	if err := gzipEntireFile(dst, gzipDst); err != nil {
		return err
	}
	
	// Get final gzipped size
	gzipStat, err := os.Stat(gzipDst)
	if err != nil {
		return err
	}
	finalSize := gzipStat.Size()
	logDebug("Final gzipped file size: %d bytes", finalSize)
	
	fmt.Printf("Final compression: %d -> %d bytes (%.1f%% reduction)\n", 
		fileSize, finalSize, (1.0-float64(finalSize)/float64(fileSize))*100)
	fmt.Printf("Gzip improvement: %d -> %d bytes (%.1f%% reduction)\n", 
		rawSize, finalSize, (1.0-float64(finalSize)/float64(rawSize))*100)
	
	// Compare to original (for JPEG testing)
	if fileSize > finalSize {
		fmt.Printf("SUCCESS: Compressed file is smaller than original!\n")
	} else {
		fmt.Printf("Original file was already more compact.\n")
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

	// Find the metadata delimiter
	delimiter := []byte("GOXORRO_META_START")
	file_content, err := io.ReadAll(sourceFile)
	if err != nil {
		return err
	}
	
	// Find delimiter position
	delimiterPos := -1
	for i := 0; i <= len(file_content)-len(delimiter); i++ {
		if string(file_content[i:i+len(delimiter)]) == string(delimiter) {
			delimiterPos = i
			break
		}
	}
	
	if delimiterPos == -1 {
		return fmt.Errorf("metadata delimiter not found")
	}
	
	// Extract metadata section
	metadataStart := delimiterPos + len(delimiter)
	metadataBytes := file_content[metadataStart:]
	logDebug("Metadata starts at position %d, size: %d bytes", metadataStart, len(metadataBytes))

	allCompressionSteps, err := decodeMetadataWithDelimiters(metadataBytes)
	if err != nil {
		return err
	}
	totalChunks := len(allCompressionSteps)
	logDebug("Decoded metadata for %d chunks", totalChunks)

	fmt.Printf("Decompressing %d chunks...\n", totalChunks)

	// Process data section (before delimiter)
	dataSection := file_content[:delimiterPos]
	logDebug("Data section size: %d bytes", len(dataSection))
	
	chunkIndex := 0
	totalDecompressed := 0
	var avgDecompTime time.Duration
	dataOffset := 0
	
	for chunkIndex < totalChunks {
		if dataOffset >= len(dataSection) {
			break
		}

		// Read chunk size (4 bytes, uint32 little endian)
		if dataOffset+4 > len(dataSection) {
			break
		}
		chunkSize := binary.LittleEndian.Uint32(dataSection[dataOffset:dataOffset+4])
		dataOffset += 4

		// Read compressed data
		if dataOffset+int(chunkSize) > len(dataSection) {
			return fmt.Errorf("chunk data extends beyond data section")
		}
		compressedData := dataSection[dataOffset:dataOffset+int(chunkSize)]
		dataOffset += int(chunkSize)

		if chunkIndex >= len(allCompressionSteps) {
			return fmt.Errorf("chunk index out of range")
		}

		logDebug("Decompressing chunk %d: %d bytes compressed", chunkIndex, chunkSize)
		start := time.Now()
		originalData := decompressChunk(compressedData, allCompressionSteps[chunkIndex])
		duration := time.Since(start)
		
		// Calculate decompression time estimates
		if chunkIndex == 0 {
			avgDecompTime = duration
		} else {
			avgDecompTime = (avgDecompTime*time.Duration(chunkIndex) + duration) / time.Duration(chunkIndex+1)
		}
		
		remainingChunks := totalChunks - (chunkIndex + 1)
		estimatedRemaining := avgDecompTime * time.Duration(remainingChunks)
		progress := float64(chunkIndex+1) / float64(totalChunks) * 100
		
		// Format time estimates for decompression
		timeStr := ""
		if chunkIndex > 1 && remainingChunks > 0 {
			if estimatedRemaining < time.Minute {
				timeStr = fmt.Sprintf(" ETA: %ds", int(estimatedRemaining.Seconds()))
			} else {
				timeStr = fmt.Sprintf(" ETA: %dm%ds", int(estimatedRemaining.Minutes()), int(estimatedRemaining.Seconds())%60)
			}
		}
		
		fmt.Printf("\rProgress: %.1f%% (%d/%d chunks)%s", progress, chunkIndex+1, totalChunks, timeStr)
		
		logDebug("Chunk %d decompressed: %d -> %d bytes in %v", 
			chunkIndex, chunkSize, len(originalData), duration)
		
		if _, err := destFile.Write(originalData); err != nil {
			return err
		}

		totalDecompressed += len(originalData)
		chunkIndex++
	}

	fmt.Printf("\rProgress: 100.0%% (%d/%d chunks) - Complete!\n", totalChunks, totalChunks)

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

	
	// Continue XOR operations until no improvement is possible (limit iterations for performance)
	maxIterations := 10
	if len(bitStream) >= 8192 { // 1KB chunks - optimize for JPEG testing
		maxIterations = 5
	}
	
	for iteration := 0; iteration < maxIterations; iteration++ {
		logDebug("XOR iteration %d: testing primes for bit stream with %d bits", 
			iteration, len(bitStream))
		
		bestPrime := -1
		bestZeros := 0
		var bestResult []int
		primesChecked := 0
		currentZeros := countZeros(bitStream)
		
		// For performance, test much fewer primes for large files
		maxPrimesToTest := 50
		if len(bitStream) < 2048 { // Small chunks get more prime testing
			maxPrimesToTest = 200
		}
		
		for i, prime := range primes {
			if i >= maxPrimesToTest {
				break
			}
			
			xorPattern := createXORPattern(len(bitStream), prime)
			testResult := xorBits(bitStream, xorPattern)
			
			zeros := countZeros(testResult)
			primesChecked++
			
			// Only accept if we get MORE zeros than current (strict improvement)
			if zeros > currentZeros && zeros > bestZeros {
				bestZeros = zeros
				bestPrime = prime
				bestResult = testResult
				logDebug("New best prime %d: %d zeros (%.1f%% improvement)", 
					prime, zeros, float64(zeros-currentZeros)/float64(len(bitStream))*100)
			}
		}
		
		logDebug("Iteration %d complete: checked %d primes, best improvement: %d -> %d zeros", 
			iteration, primesChecked, currentZeros, bestZeros)
		
		// Stop if no improvement found (no prime increases zeros)
		if bestPrime == -1 {
			logDebug("No further improvement possible, stopping XOR iterations")
			break
		}
		
		bitStream = bestResult
		steps = append(steps, CompressionStep{Operation: "xor", Prime: bestPrime, Applied: true})
	}

	// Return raw bit stream as bytes (no gzip compression)
	compressed := bitsToBytes(bitStream)
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

func countZeros(bits []int) int {
	zeros := 0
	for _, bit := range bits {
		if bit == 0 {
			zeros++
		}
	}
	return zeros
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
	return float64(ones)/float64(len(bits)) < 0.05
}


func decompressChunk(compressedData []byte, steps []CompressionStep) []byte {
	// Convert raw bytes back to bit stream (no gzip decompression)
	bitStream := bytesToBits(compressedData)
	
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

func encodeMetadataWithDelimiters(allSteps [][]CompressionStep) []byte {
	var result []byte
	
	// Write number of chunks
	result = append(result, byte(len(allSteps)))
	
	// Write metadata for each chunk with delimiters
	for chunkIndex, steps := range allSteps {
		if chunkIndex > 0 {
			// Add delimiter between chunks
			delimiter := []byte("CHUNK_SEP")
			result = append(result, delimiter...)
		}
		
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
		
		var primeIndices []int
		for _, step := range xorSteps {
			primeIndex := getPrimeIndex(step.Prime)
			if primeIndex >= 0 && primeIndex < len(primes) {
				primeIndices = append(primeIndices, primeIndex)
			}
		}
		
		packedBits := pack11BitIndices(primeIndices)
		result = append(result, packedBits...)
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
		
		var primeIndices []int
		for _, step := range xorSteps {
			primeIndex := getPrimeIndex(step.Prime)
			if primeIndex >= 0 && primeIndex < len(primes) {
				primeIndices = append(primeIndices, primeIndex)
			}
		}
		
		packedBits := pack11BitIndices(primeIndices)
		result = append(result, packedBits...)
	}
	
	return result
}

func decodeMetadataWithDelimiters(data []byte) ([][]CompressionStep, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty metadata")
	}
	
	numChunks := int(data[0])
	var allSteps [][]CompressionStep
	
	offset := 1
	delimiter := []byte("CHUNK_SEP")
	
	for i := 0; i < numChunks; i++ {
		// Skip delimiter if not first chunk
		if i > 0 {
			if offset+len(delimiter) <= len(data) && 
			   string(data[offset:offset+len(delimiter)]) == string(delimiter) {
				offset += len(delimiter)
			}
		}
		
		if offset+1 >= len(data) {
			return nil, fmt.Errorf("invalid metadata: insufficient data")
		}
		
		negateApplied := data[offset]
		numXorSteps := int(data[offset+1])
		offset += 2
		
		var steps []CompressionStep
		
		if negateApplied == 1 {
			steps = append(steps, CompressionStep{Operation: "negate", Applied: true})
		}
		
		if numXorSteps > 0 {
			packedBitsLength := (numXorSteps*11 + 7) / 8
			if offset+packedBitsLength > len(data) {
				return nil, fmt.Errorf("invalid metadata: insufficient XOR step data")
			}
			
			packedBits := data[offset : offset+packedBitsLength]
			primeIndices := unpack11BitIndices(packedBits, numXorSteps)
			
			for _, primeIndex := range primeIndices {
				if primeIndex < len(primes) {
					steps = append(steps, CompressionStep{
						Operation: "xor",
						Prime:     primes[primeIndex],
						Applied:   true,
					})
				}
			}
			
			offset += packedBitsLength
		}
		
		allSteps = append(allSteps, steps)
	}
	
	return allSteps, nil
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
		
		var steps []CompressionStep
		
		if negateApplied == 1 {
			steps = append(steps, CompressionStep{Operation: "negate", Applied: true})
		}
		
		if numXorSteps > 0 {
			packedBitsLength := (numXorSteps*11 + 7) / 8
			if offset+packedBitsLength > len(data) {
				return nil, fmt.Errorf("invalid metadata: insufficient XOR step data")
			}
			
			packedBits := data[offset : offset+packedBitsLength]
			primeIndices := unpack11BitIndices(packedBits, numXorSteps)
			
			for _, primeIndex := range primeIndices {
				if primeIndex < len(primes) {
					steps = append(steps, CompressionStep{
						Operation: "xor",
						Prime:     primes[primeIndex],
						Applied:   true,
					})
				}
			}
			
			offset += packedBitsLength
		}
		
		allSteps = append(allSteps, steps)
	}
	
	return allSteps, nil
}

func gzipEntireFile(srcPath, dstPath string) error {
	// Open source file
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	
	// Create destination file
	dstFile, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()
	
	// Create gzip writer
	gzWriter := gzip.NewWriter(dstFile)
	defer gzWriter.Close()
	
	// Copy and compress
	_, err = io.Copy(gzWriter, srcFile)
	return err
}

func pack11BitIndices(indices []int) []byte {
	if len(indices) == 0 {
		return []byte{}
	}
	
	totalBits := len(indices) * 11
	totalBytes := (totalBits + 7) / 8
	result := make([]byte, totalBytes)
	
	bitOffset := 0
	for _, index := range indices {
		for i := 0; i < 11; i++ {
			bit := (index >> (10 - i)) & 1
			byteIndex := bitOffset / 8
			bitPosition := bitOffset % 8
			
			if bit == 1 {
				result[byteIndex] |= 1 << (7 - bitPosition)
			}
			bitOffset++
		}
	}
	
	return result
}

func unpack11BitIndices(packedBits []byte, count int) []int {
	if count == 0 {
		return []int{}
	}
	
	result := make([]int, count)
	bitOffset := 0
	
	for i := 0; i < count; i++ {
		index := 0
		for j := 0; j < 11; j++ {
			byteIndex := bitOffset / 8
			bitPosition := bitOffset % 8
			
			if byteIndex < len(packedBits) {
				bit := (packedBits[byteIndex] >> (7 - bitPosition)) & 1
				index = (index << 1) | int(bit)
			}
			bitOffset++
		}
		result[i] = index
	}
	
	return result
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
	debugLogger.Printf("Using %d precomputed prime numbers (2 to %d)", len(primes), primes[len(primes)-1])
}

func logDebug(format string, args ...interface{}) {
	if verboseMode && debugLogger != nil {
		debugLogger.Printf(format, args...)
	}
}
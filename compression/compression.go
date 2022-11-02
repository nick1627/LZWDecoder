package compression

import (
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	// "bytes"
)

type Dictionary struct {
	length  uint16
	entries [4096][]byte
}

func (dict *Dictionary) GetLength() uint16 {
	return dict.length
}

func (dict *Dictionary) Clear() {
	// Just have to reset the length
	dict.length = 256
}

func (dict *Dictionary) Initialise() {
	var i uint16
	for i = 0; i < 256; i++ {
		dict.entries[i] = []byte{uint8(i)}
	}
	dict.length = 256
}

func (dict *Dictionary) AddEntry(newElement []byte) {
	if dict.length == 4096 {
		dict.Clear()
	}
	dict.entries[dict.length] = newElement
	dict.length += 1
}

func (dict *Dictionary) GetEntry(index uint16) ([]byte, error) {
	if index >= dict.length {

		return []byte{}, errors.New("index out of range of dictionary")
	} else {
		return dict.entries[index], nil
	}
}

func getCodes(byteList []byte) [2]uint16 {
	/*
		This function takes three consecutive bytes.  It
		splits the second byte in half and concatenates
		left and right such that there are now two lots
		of 12 bits, treated as numbers.  These numbers
		are returned as two integers in an array.

		TODO: Is there a better way using splitting
		and concatenation of the bits?
	*/

	var total uint32
	total += uint32(byteList[2])
	total += (uint32(byteList[1]) * uint32(math.Pow(2, 8)))
	total += (uint32(byteList[0]) * uint32(math.Pow(2, 16)))

	total2 := total % uint32(math.Pow(2, 12))
	total -= total2
	total /= uint32(math.Pow(2, 12))

	// firstCode, secondCode := getCodes(buffer)
	codes := [2]uint16{uint16(total), uint16(total2)}

	return codes
}

func Decompress(encodedFile string) {
	// encodedFile:		The path of the file to be decoded

	newFilePath := encodedFile[:len(encodedFile)-2]

	// Open the encoded file
	file, errMsg := os.Open(encodedFile)
	if errMsg != nil {
		fmt.Println(errMsg)
	}
	defer file.Close()

	newFile, errMsg := os.Create(newFilePath)
	if errMsg != nil {
		fmt.Println(errMsg)
	}
	defer newFile.Close()

	// Will now read out 3 bytes at a time into the bufferArray.
	buffer := make([]byte, 3)

	// Create and fill the dictionary using previously defined struct
	var dictionary Dictionary
	dictionary.Initialise()

	// Need to keep track of what was emitted previously
	var lastEmitted []byte
	var endReached bool = false

	// Loop until end of file
	for {
		// Fill buffer with new bytes from file
		_, errMsg := file.Read(buffer)
		if errMsg != nil {
			if errMsg != io.EOF {
				fmt.Println(errMsg)
				break
			} else {
				endReached = true
			}
			// Error or end of file, so break loop
			break
		}

		// Deal with EOF edge cases
		var startPoint int = 0
		if endReached {
			// Either all bytes in the buffer are zero,
			// or only the last one is zero, because
			// the file must have an even number of
			// bytes.
			if buffer[0] == 0 {
				// have reached end of file
				break
			} else {
				// Third byte is zero.  Special treatment
				// required to allow the code to be
				// extracted from the first two bytes
				// properly.
				if buffer[1] == 0 && buffer[2] == 0 {
					fmt.Println("Something's broken")
				}
				buffer[2] = buffer[1]
				buffer[1] = buffer[0]
				buffer[0] = 0

				// First code contains no data, so set startPoint
				// to 1 so that it is skipped over
				startPoint = 1
			}
		}

		// Extract the two codes from these three bytes
		codes := getCodes(buffer)
		for i := startPoint; i < 2; i++ {
			// For every code, we apply the rules of the LZW decoding algorithm.
			// See https://en.wikipedia.org/wiki/Lempel–Ziv–Welch
			if codes[i] >= dictionary.GetLength() {
				// The code is not in the dictionary
				v := append(lastEmitted, lastEmitted[0])
				dictionary.AddEntry(v)
				fmt.Println(string(v))

				_, errMsg = newFile.Write(v)
				if errMsg != nil {
					fmt.Println(errMsg)
				}

				lastEmitted = v
			} else {
				// The code is in the dictionary
				w, _ := dictionary.GetEntry(codes[i])

				_, errMsg = newFile.Write(w)
				if errMsg != nil {
					fmt.Println(errMsg)
				}

				if len(lastEmitted) != 0 {
					newEntry := append(lastEmitted, w[0])
					dictionary.AddEntry(newEntry)
					fmt.Println(string(newEntry))
				}

				lastEmitted = w
			}
		}
	}
}

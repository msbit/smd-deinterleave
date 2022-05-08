package main

import (
	"fmt"
	"os"
)

const SIZE = 16384

type smd_header_t struct {
	size_of_file   byte // Byte 00h : Size of file.
	file_data_type byte // Byte 01h : File data type.
	status_flags_1 byte // Byte 02h : Status flags.
	status_flags_2 byte //
	status_flags_3 byte //
	status_flags_4 byte //
	status_flags_5 byte //
	status_flags_6 byte //
	identifier_1   byte // Byte 08h : Identifier 1.
	identifier_2   byte // Byte 09h : Identifier 2.
	file_type      byte // Byte 0Ah : File type.
}

func NewHeader(data [11]byte) smd_header_t {
	var result smd_header_t

	result.size_of_file = data[0]
	result.file_data_type = data[1]
	result.status_flags_1 = data[2]
	result.status_flags_2 = data[3]
	result.status_flags_3 = data[4]
	result.status_flags_4 = data[5]
	result.status_flags_5 = data[6]
	result.status_flags_6 = data[7]
	result.identifier_1 = data[8]
	result.identifier_2 = data[9]
	result.file_type = data[10]

	return result
}

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "usage: %s <input-file> <output-file>\n", os.Args[0])
		os.Exit(-1)
	}

	input, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Open input: %s\n", err)
		os.Exit(-1)
	}
	defer input.Close()

	header_buffer := make([]byte, 11)
	_, err = input.Read(header_buffer)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Read header: %s\n", err)
		os.Exit(-1)
	}

	header := NewHeader(*(*[11]byte)(header_buffer))
	if header.identifier_1 != 0xaa || header.identifier_2 != 0xbb || header.file_type != 0x06 {
		fmt.Fprintf(
			os.Stderr,
			"invalid file: identifier_1 = 0x%02x identifier_2 = 0x%02x file_type = 0x%02x\n",
			header.identifier_1,
			header.identifier_2,
			header.file_type)
		os.Exit(-1)
	}

	_, err = input.Seek(512, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Seek input: %s\n", err)
		os.Exit(-1)
	}

	output, err := os.Create(os.Args[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Create output: %s\n", err)
		os.Exit(-1)
	}
	defer output.Close()

	var inputBuffer [SIZE]byte
	var outputBuffer [SIZE]byte

	for i := byte(0); i < header.size_of_file; i++ {
		_, err := input.Read(inputBuffer[:])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Read chunk: %s\n", err)
			os.Exit(-1)
		}

		for j, k := 0, 0; j < SIZE/2; j, k = j+1, k+2 {
			outputBuffer[k] = inputBuffer[j+(SIZE/2)]
			outputBuffer[k+1] = inputBuffer[j]
		}

		_, err = output.Write(outputBuffer[:])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Write chunk: %s\n", err)
			os.Exit(-1)
		}
	}
}

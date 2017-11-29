package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

type LLRSHeader struct {
	Magic      uint32
	Version    uint32
	FilesCount uint32
	Unknown1   uint32
	Unknown2   uint32
}

type FileRecord struct {
	Name   [12]byte
	Unk    uint16
	Unk1   uint32
	Offset uint32
	Length uint32
}

func main() {

	var (
		i   uint32
		err error
	)

	if len(os.Args) != 2 {
		fmt.Printf("LLRS [Leo the Lion / Lew Leon] resource files extractor\n"+
			"-------------------------------------------------------\n"+
			"Extracts all available files to the current directory.\n\nUsage: %s <filename>\n"+
			"Example: %s SOUND\n", os.Args[0], os.Args[0])
		os.Exit(0)
	}

	file, err := os.Open(os.Args[1])

	if err != nil {
		panic(err)
	}

	defer file.Close()

	header := LLRSHeader{}
	binary.Read(file, binary.LittleEndian, &header)

	if header.Magic != 0x73726c6c {
		fmt.Printf("Unknown magick: 0x%X, must be 'llrs'\n", header.Magic)
		os.Exit(-1)
	}

	fmt.Printf("LLRS File version %d, %d files\n", header.Version, header.FilesCount)

	for i = 0; i < header.FilesCount; i++ {
		fr := FileRecord{}
		err = binary.Read(file, binary.LittleEndian, &fr)

		if err != nil {
			fmt.Println(err)
		} else {

			// Skip special records
			if fr.Unk == 0 {
				fmt.Printf("%d: File: %s  Offset: 0x%x Length: %d\n", i, fr.Name, fr.Offset, fr.Length)

				buf := make([]byte, fr.Length)
				count, err := file.ReadAt(buf, int64(fr.Offset))
				if err != nil {
					fmt.Printf("ReadError: %s, readed: %d bytes instead of %d, skipping ...\n", err, count, fr.Length)
				} else {

					// Truncate trailing zeros in the filename
					s := fmt.Sprintf("%s", bytes.Trim(fr.Name[:], "\x00"))

					part, err := os.Create(s)

					if err != nil {
						fmt.Printf("Error: %s\n", err)

					} else {
						part.Write(buf)
					}

					part.Close()
				}

			}
		}
	}
	fmt.Println()
}

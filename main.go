package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"reflect"

	"github.com/alecthomas/kong"
)

const (
	NeroV1 NeroFormat = iota
	NeroV2
)

var (
	NeroV1FormatCode = []byte("NERO")
	NeroV2FormatCode = []byte("NER5")
)

var cli struct {
	SourceNRG string `arg`
	TargetISO string `arg`
}

func main() {
	kong.Parse(&cli,
		kong.Description("This command-line tool converts NERO image files (.nrg) to ISO image format (.iso)."))

	err := convert()
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}

	fmt.Printf("Done!\n")
}

func convert() error {
	inFilePath := cli.SourceNRG
	outFilePath := cli.TargetISO
	inFile, err := os.Open(inFilePath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer inFile.Close()

	format, err := getFormat(inFile)
	if err != nil {
		return fmt.Errorf("failed parsing source file: %w", err)
	}

	firstChunkOffset, err := getFirstChunkOffset(inFile, format)
	if err != nil {
		return fmt.Errorf("failed parsing source file: %w", err)
	}

	_, err = inFile.Seek(firstChunkOffset, io.SeekStart)
	if err != nil {
		return fmt.Errorf("failed reading source file: %w", err)
	}

	for {
		chunk, err := readChunk(inFile)
		if err != nil {
			return fmt.Errorf("failed reading source file: %w", err)
		}

		if daox, ok := chunk.(*DAOX); ok {
			outFile, err := os.Create(outFilePath)
			if err != nil {
				return fmt.Errorf("create target file error: %w", err)
			}
			defer outFile.Close()

			_, err = inFile.Seek(daox.Index1, io.SeekStart)
			if err != nil {
				return fmt.Errorf("seek error: %w", err)
			}

			_, err = io.CopyN(outFile, inFile, daox.EndOfTrack-daox.Index1)
			if err != nil {
				return fmt.Errorf("copy error: %w", err)
			}

			break
		}

		if _, ok := chunk.(*END); ok {
			return fmt.Errorf("image format not supported: %w", err)
		}
	}

	return nil
}

func readChunk(file *os.File) (any, error) {
	chunkHeader := ChunkHeader{}
	err := binary.Read(file, binary.BigEndian, &chunkHeader)
	if err != nil {
		return nil, fmt.Errorf("read chunk header error: %w", err)
	}

	chunkID := string(chunkHeader.ID[:])

	var chunk any
	switch chunkID {
	case "CUEX":
		chunk = &CUEX{}
	case "DAOX":
		chunk = &DAOX{}
	case "SINF":
		chunk = &SINF{}
	case "MTYP":
		chunk = &MTYP{}
	case "END!":
		chunk = &END{}
	default:
		return nil, fmt.Errorf("uknown chunk ID: %s", chunkID)
	}

	err = binary.Read(file, binary.BigEndian, chunk)
	if err != nil {
		return nil, fmt.Errorf("read chunk error: %w", err)
	}

	chunkBodySize := reflect.Indirect(reflect.ValueOf(chunk)).Type().Size()
	offset := int64(chunkHeader.Size) - int64(chunkBodySize)

	if offset > 0 {
		_, err = file.Seek(offset, io.SeekCurrent)
		if err != nil {
			return nil, fmt.Errorf("seek error: %w", err)
		}
	}

	return chunk, nil
}

func getFirstChunkOffset(file *os.File, format NeroFormat) (int64, error) {
	size := 4
	if format == NeroV2 {
		size = 8
	}

	_, err := file.Seek(-int64(size), io.SeekEnd)
	if err != nil {
		return 0, err
	}

	var offset int64
	err = binary.Read(file, binary.BigEndian, &offset)
	if err != nil {
		return 0, err
	}

	return offset, nil
}

func getFormat(file *os.File) (NeroFormat, error) {
	formats := []struct {
		offset int64
		format NeroFormat
		code   []byte
	}{
		{
			offset: -8,
			code:   NeroV1FormatCode,
			format: NeroV1,
		},
		{
			offset: -12,
			code:   NeroV2FormatCode,
			format: NeroV2,
		},
	}

	for _, f := range formats {
		_, err := file.Seek(f.offset, io.SeekEnd)
		if err != nil {
			return 0, fmt.Errorf("seek error: %w", err)
		}

		buffer := make([]byte, 4)
		_, err = file.Read(buffer)
		if err != nil {
			return 0, err
		}

		if bytes.Equal(buffer, f.code) {
			return f.format, nil
		}
	}

	return 0, fmt.Errorf("unable to determine format")
}

type NeroFormat int

func (f NeroFormat) String() string {
	switch f {
	case NeroV1:
		return string(NeroV1FormatCode)
	case NeroV2:
		return string(NeroV2FormatCode)
	default:
		return fmt.Sprintf("%d", int(f))
	}
}

type ChunkHeader struct {
	ID   [4]byte
	Size uint32
}

type CUEX struct {
	Mode        byte
	TrackNumber byte
	IndexNumber byte
	Padding     byte
	LbaPosition int32
}

type DAOX struct {
	Size       uint32
	Upc        [13]byte
	Padding    byte
	TocType    uint16
	FirstTrack byte
	LastTrack  byte
	Isrc       [12]byte
	SectorSize uint16
	Mode       uint16
	Uknown     uint16
	Index0     int64
	Index1     int64
	EndOfTrack int64
}

type SINF struct {
	TracksInSession int32
}

type MTYP struct {
	Uknown int32
}

type END struct{}

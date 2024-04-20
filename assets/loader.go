// Code from bwmarrin on GitHub, thanks for your repository <3
// https://github.com/bwmarrin/discordgo/blob/master/examples/airhorn/main.go

package assets

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

func LoadSoundFile(path string) ([][]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening dca file :", err)
		return nil, err
	}

	buffer := make([][]byte, 0)

	var opuslen int16

	for {
		// Read opus frame length from dca file.
		err = binary.Read(file, binary.LittleEndian, &opuslen)

		// If this is the end of the file, just return.
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			err := file.Close()
			if err != nil {
				return nil, err
			}
			return buffer, nil
		}

		if err != nil {
			fmt.Println("Error reading from dca file :", err)
			return nil, err
		}

		// Read encoded pcm from dca file.
		InBuf := make([]byte, opuslen)
		err = binary.Read(file, binary.LittleEndian, &InBuf)

		// Should not be any end of file errors
		if err != nil {
			fmt.Println("Error reading from dca file :", err)
			return nil, err
		}

		// Append encoded pcm data to the buffer.
		buffer = append(buffer, InBuf)
	}
}

package bin2paper

import (
	"errors"
	"fmt"
	"hash"
	"io"
	"os"
)

type BinDecoder struct {
	Filename  string
	Hash      string
	Input     io.ReadSeeker
	Output    io.StringWriter
	checkHash hash.Hash
}

func BinFromTxt(input_path string, output_path string) error {
	input, err := os.Open(input_path)
	if err != nil {
		return errors.New("failed to open " + input_path + ": " + err.Error())
	}
	stat, err := input.Stat()
	if err != nil {
		return errors.New("failed to open " + input_path + ": " + err.Error())
	}
	decoder := BinDecoder{
		Input: input,
	}


	decoder.Output, err = os.OpenFile(output_path, os.O_RDWR|os.O_CREATE, stat.Mode())
	if err != nil {
		return errors.New("failed to open " + input_path + ": " + err.Error())
	}

	decoder.GetBasics()
	return nil
}

func (this *BinDecoder) GetBasics() error {
	_, err := this.Input.Seek(0, io.SeekStart)
	if err != nil {
		return errors.New("failed to return to the begining of the input: " + err.Error())
	}
	fmt.Println(ReadLine(this.Input))

	return nil
}

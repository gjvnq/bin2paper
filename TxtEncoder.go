package bin2paper

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"hash"
	"hash/adler32"
	"io"
	"math"
	"os"
	"path/filepath"
	"strconv"
)

type TxtEncoder struct {
	Filename   string
	Width      int
	Height     int
	Input      io.ReadSeeker
	Output     io.StringWriter
	size       int64
	pos        int64
	hash       string
	checkHash  hash.Hash
	pageNum    int
	totalPages int
	nDataChars int
	nBytes     int
}

func TXTFromFile(input_path string, output_path string) error {
	input, err := os.Open(input_path)
	if err != nil {
		return errors.New("failed to open " + input_path + ": " + err.Error())
	}
	stat, err := input.Stat()
	if err != nil {
		return errors.New("failed to open " + input_path + ": " + err.Error())
	}

	encoder := TxtEncoder{
		Filename: filepath.Base(input_path),
		Input:    input,
	}
	if output_path == "" {
		encoder.Output = os.Stdout
	} else {
		encoder.Output, err = os.OpenFile(output_path, os.O_RDWR|os.O_CREATE, stat.Mode())
		if err != nil {
			return errors.New("failed to open " + input_path + ": " + err.Error())
		}
	}

	return encoder.Encode()
}

func (this *TxtEncoder) prepare() error {
	if this.Width <= 0 {
		this.Width = 72
	}
	if this.Height <= 0 {
		this.Height = 60
	}

	// Compute input hash
	var err error
	h := sha1.New()
	if this.size, err = io.Copy(h, this.Input); err != nil {
		return errors.New("failed to hash input file: " + err.Error())
	}
	this.hash = ToHex(h.Sum(nil))
	_, err = this.Input.Seek(0, io.SeekStart)
	if err != nil {
		return errors.New("failed to return to the begining of the input: " + err.Error())
	}

	// Find how much data each line and page contain
	this.checkHash = adler32.New()
	example_checksum := this.calcCheckHash(nil)
	this.nDataChars = this.Width - len("00: ") - len(" "+example_checksum)
	this.nBytes = 6 * this.nDataChars / 8

	bytes_per_page := (this.Height - 5) * this.nBytes
	this.totalPages = int(math.Ceil(float64(this.size) / float64(bytes_per_page)))
	this.pageNum = 1

	return nil
}

func (this *TxtEncoder) writeHeaders() error {
	s := ""
	s += fmt.Sprintf("File: %s\n", this.Filename)
	s += fmt.Sprintf("Start/Size: %7d/%d bytes | Checksum: ADLER32\n", this.pos, this.size)
	s += fmt.Sprintf("Hash: SHA1:%s\n", this.hash)
	s += fmt.Sprintf("Page: %d/%d\n", this.pageNum, this.totalPages)
	s += "\n"
	_, err := this.Output.WriteString(s)
	if err != nil {
		return err
	}
	return nil
}

func (this *TxtEncoder) calcCheckHash(buf []byte) string {
	this.checkHash.Reset()
	this.checkHash.Write(buf)
	return base64.RawStdEncoding.EncodeToString(this.checkHash.Sum(nil))
}

func (this *TxtEncoder) writeDataPage() error {
	for line := 5; line < this.Height; line++ {
		// Read data
		buf := make([]byte, this.nBytes)
		n, err := this.Input.Read(buf)
		if n != 0 && err != nil {
			return err
		}
		if n == 0 {
			return nil
		}
		buf = buf[:n] // Trim it too!
		this.pos += int64(n)

		s := fmt.Sprintf("%02d: %-"+strconv.Itoa(this.nDataChars)+"s %s\n", line-5, ToB64str(buf), this.calcCheckHash(buf))

		_, err = this.Output.WriteString(s)
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *TxtEncoder) writePage() error {
	if err := this.writeHeaders(); err != nil {
		return err
	}
	if err := this.writeDataPage(); err != nil {
		return err
	}
	this.pageNum++
	return nil
}

func (this *TxtEncoder) Encode() error {
	this.prepare()
	for this.pos < this.size {
		if err := this.writePage(); err != nil {
			return err
		}
	}
	return nil
}

package bin2paper

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"path/filepath"
	"hash/adler32"
	"crypto/sha1"
	"strconv"
	"errors"
	"os"
	"math"
	"hash"
)

const BUF_SIZE_FACTOR = 64

// 72 columns by 60 rows
type Config struct {
	Filename string
	Width int
	Height int
	Input io.ReadSeeker
	Output io.StringWriter
	size int64
	pos int64
	hash string
	checkHash hash.Hash
	pageNum int
	totalPages int
	nDataChars int
	nBytes int
}

func ToB64(src []byte) (dst []byte) {
	base64.StdEncoding.Encode(dst, src)
	return
}

func ToB64str(src []byte) string {
	return base64.StdEncoding.EncodeToString(src)
}

func ToHex(src []byte) string {
	return hex.EncodeToString(src)
}

func TXTFromFile(input_path string) error {
	input, err := os.Open(input_path)
	if err != nil {
		return errors.New("failed to open "+input_path+": "+err.Error())
	}

	config := Config{
		Filename: filepath.Base(input_path),
		Input: input,
		Output: os.Stdout,
	}

	return config.Do()
}

func (this *Config) prepare() error {
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
		return errors.New("failed to hash input file: "+err.Error())
	}
	this.hash = ToHex(h.Sum(nil))
	_, err = this.Input.Seek(0, io.SeekStart)
	if err != nil {
		return errors.New("failed to return to the begining of the input: "+err.Error())
	}

	// Find how much data each line and page contain
	this.checkHash = adler32.New()
	example_checksum := this.calcCheckHash(nil)
	this.nDataChars = this.Width - len("00: ") - len(" "+example_checksum)
	this.nBytes = 6*this.nDataChars/8

	bytes_per_page := (this.Height - 5)*this.nBytes
	this.totalPages = int(math.Ceil(float64(this.size)/float64(bytes_per_page)))
	this.pageNum = 1

	return nil
}

func (this *Config) writeHeaders() error {
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

func (this *Config) calcCheckHash(buf []byte) string {
	this.checkHash.Reset()
	this.checkHash.Write(buf)
	return base64.RawStdEncoding.EncodeToString(this.checkHash.Sum(nil))
}

func (this *Config) writeDataPage() error {
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

func (this *Config) writePage() error {
	if err := this.writeHeaders(); err != nil {
		return err
	}
	if err := this.writeDataPage(); err != nil {
		return err
	}
	this.pageNum++
	return nil
}

func (this *Config) Do() error {
	this.prepare()
	for this.pos < this.size {
		if err := this.writePage(); err != nil {
			return err
		}
	}
	return nil
}
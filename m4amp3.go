package webpbin

import (
	"errors"
	"io"

	"github.com/nickalie/go-binwrapper"
)

// Mp4mp3 compresses an image using the WebP format. Input format can be either PNG, JPEG, TIFF, WebP or raw Y'CbCr samples.
// https://developers.google.com/speed/webp/docs/cwebp
type Mp4mp3 struct {
	*binwrapper.BinWrapper
	input      io.Reader
	inputFile  string
	output     io.Writer
}

// NewCWebP creates new Mp4mp3 instance.
func NewMp4mp3(optionFuncs ...OptionFunc) *Mp4mp3 {
	bin := &Mp4mp3{
		BinWrapper: createBinWrapper(optionFuncs...),
	}
	bin.ExecPath("ffmpeg")

	return bin
}

// Version returns Mp4mp3 version.
func (c *Mp4mp3) Version() (string, error) {
	return version(c.BinWrapper)
}

// Input sets reader to convert.
// InputFile or InputImage called before will be ignored.
func (c *Mp4mp3) Input(reader io.Reader) *Mp4mp3 {
	c.input = reader
	return c
}

// InputFile sets image file to convert.
// Input or InputImage called before will be ignored.
func (c *Mp4mp3) InputFile(file string) *Mp4mp3 {
	c.input = nil
	c.inputFile = file
	return c
}

// Output specify writer to write webp file content.
// OutputFile called before will be ignored.
func (c *Mp4mp3) Output(writer io.Writer) *Mp4mp3 {
	c.output = writer
	return c
}

func (c *Mp4mp3) Run() error {
	defer c.BinWrapper.Reset()

	err := c.setInput()

	if err != nil {
		return err
	}

	c.Arg("-codec:v", "copy").Arg("-codec:a", "libmp3lame").Arg("-qscale:a", "4")
	c.Arg("-f", "mp3")

	c.Arg("-")

	if c.output != nil {
		c.SetStdOut(c.output)
	}

	err = c.BinWrapper.Run()

	if err != nil {
		return errors.New(err.Error() + ". " + string(c.StdErr()))
	}

	return nil
}

func (c *Mp4mp3) setInput() error {
	if c.input != nil {
		c.Arg("-i").Arg("-")
		c.StdIn(c.input)
	} else if c.inputFile != "" {
		c.Arg(c.inputFile)
	} else {
		return errors.New("Undefined input")
	}

	return nil
}

func (c *Mp4mp3) getOutput() (string, error) {
	if c.output != nil {
		return "-", nil
	} else {
		return "", errors.New("Undefined output")
	}
}

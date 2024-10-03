package webpbin

import (
	"errors"
	"io"

	"github.com/Duckduckgot/go-binwrapper"
)

// IWebP compresses an image using the WebP format. Input format can be either PNG, JPEG, TIFF, WebP or raw Y'CbCr samples.
// https://developers.google.com/speed/webp/docs/cwebp
type IWebP struct {
	*binwrapper.BinWrapper
	input  io.Reader
	inputFile  string
	inputFile2  string
	output io.Writer
}

// NewIWebP creates new IWebP instance.
func NewIWebP(optionFuncs ...OptionFunc) *IWebP {
	bin := &IWebP{
		BinWrapper: createBinWrapper(optionFuncs...),
	}
	bin.ExecPath("img2webp")

	return bin
}

// Version returns img2WebP version.
func (c *IWebP) Version() (string, error) {
	return version(c.BinWrapper)
}

// Input sets reader to convert.
// InputFile or InputImage called before will be ignored.
func (c *IWebP) Input(reader, reader2 io.Reader) *IWebP {
	c.input = reader
	return c
}

// InputFile sets image file to convert.
// Input or InputImage called before will be ignored.
func (c *IWebP) InputFile(file, file2 string) *IWebP {
	c.input = nil
	c.inputFile = file
	c.inputFile2 = file2
	return c
}

// Output specify writer to write webp file content.
// OutputFile called before will be ignored.
func (c *IWebP) Output(writer io.Writer) *IWebP {
	c.output = writer
	return c
}

func (c *IWebP) Run() error {
	defer c.BinWrapper.Reset()

	output, err := c.getOutput()

	if err != nil {
		return err
	}

	c.Arg("-loop", "0").Arg("-lossy", "")

	err = c.setInput()

	if err != nil {
		return err
	}

	c.Arg("-o", output)

	if c.output != nil {
		c.SetStdOut(c.output)
	}

	err = c.BinWrapper.Run()

	if err != nil {
		return errors.New(err.Error() + ". " + string(c.StdErr()))
	}

	return nil
}

func (c *IWebP) setInput() error {
	if c.input != nil {
		c.Arg("-d").Arg("5000")
		c.Arg("--").Arg("-")
		c.StdIn(c.input)
	} else if c.inputFile != "" {
		c.Arg("-d").Arg("5000")
		c.Arg(c.inputFile)
		if c.inputFile2 != "" {
			c.Arg("-d").Arg("3000")
			c.Arg(c.inputFile2)
		}
	} else {
		return errors.New("Undefined input")
	}

	return nil
}

func (c *IWebP) getOutput() (string, error) {
	if c.output != nil {
		return "-", nil
	} else {
		return "", errors.New("Undefined output")
	}
}

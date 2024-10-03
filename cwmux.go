package webpbin

import (
	"errors"
	"io"

	"github.com/Duckduckgot/go-binwrapper"
)

// CWMux compresses an image using the WebP format. Input format can be either PNG, JPEG, TIFF, WebP or raw Y'CbCr samples.
// https://developers.google.com/speed/webp/docs/cwebp
type CWMux struct {
	*binwrapper.BinWrapper
	input      io.Reader
	output     io.Writer
}

// NewCWebP creates new CWMux instance.
func NewCWMux(optionFuncs ...OptionFunc) *CWMux {
	bin := &CWMux{
		BinWrapper: createBinWrapper(optionFuncs...),
	}
	bin.ExecPath("webpmux")

	return bin
}

// Version returns CWMux version.
func (c *CWMux) Version() (string, error) {
	return version(c.BinWrapper)
}

// Input sets reader to convert.
// InputFile or InputImage called before will be ignored.
func (c *CWMux) Input(reader io.Reader) *CWMux {
	c.input = reader
	return c
}

// Output specify writer to write webp file content.
// OutputFile called before will be ignored.
func (c *CWMux) Output(writer io.Writer) *CWMux {
	c.output = writer
	return c
}

func (c *CWMux) Run() error {
	defer c.BinWrapper.Reset()

	output, err := c.getOutput()

	if err != nil {
		return err
	}

	c.Arg("-o", output)

	err = c.setInput()

	if err != nil {
		return err
	}

	if c.output != nil {
		c.SetStdOut(c.output)
	}

	err = c.BinWrapper.Run()

	if err != nil {
		return errors.New(err.Error() + ". " + string(c.StdErr()))
	}

	return nil
}

func (c *CWMux) setInput() error {
	if c.input != nil {
		c.Arg("--").Arg("-")
		c.StdIn(c.input)
	} else {
		return errors.New("Undefined input")
	}

	return nil
}

func (c *CWMux) getOutput() (string, error) {
	if c.output != nil {
		return "-", nil
	} else {
		return "", errors.New("Undefined output")
	}
}

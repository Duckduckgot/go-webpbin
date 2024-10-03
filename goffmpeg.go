package webpbin

import (
	"errors"
	"io"

	"github.com/Duckduckgot/go-binwrapper"
)

// Ffmpeg compresses an image using the WebP format. Input format can be either PNG, JPEG, TIFF, WebP or raw Y'CbCr samples.
// https://developers.google.com/speed/webp/docs/cwebp
type Ffmpeg struct {
	*binwrapper.BinWrapper
	input      io.Reader
	inputFile  string
	output     io.Writer
}

// NewCWebP creates new Ffmpeg instance.
func NewFfmpeg(optionFuncs ...OptionFunc) *Ffmpeg {
	bin := &Ffmpeg{
		BinWrapper: createBinWrapper(optionFuncs...),
	}
	bin.ExecPath("ffmpeg")

	return bin
}

// Version returns Ffmpeg version.
func (c *Ffmpeg) Version() (string, error) {
	return version(c.BinWrapper)
}

// Input sets reader to convert.
// InputFile or InputImage called before will be ignored.
func (c *Ffmpeg) Input(reader io.Reader) *Ffmpeg {
	c.input = reader
	return c
}

// InputFile sets image file to convert.
// Input or InputImage called before will be ignored.
func (c *Ffmpeg) InputFile(file string) *Ffmpeg {
	c.input = nil
	c.inputFile = file
	return c
}

// Output specify writer to write webp file content.
// OutputFile called before will be ignored.
func (c *Ffmpeg) Output(writer io.Writer) *Ffmpeg {
	c.output = writer
	return c
}

func (c *Ffmpeg) Run() error {
	defer c.BinWrapper.Reset()

	err := c.setInput()

	if err != nil {
		return err
	}

	c.Arg("-pix_fmt", "yuva420p").Arg("-t", "10").Arg("-r", "8")
	c.Arg("-filter:v", "fps=15,scale=512:512:force_original_aspect_ratio=decrease:flags=lanczos,format=rgba,pad=512:512:-1:-1:color=#00000000")
	c.Arg("-lossless", "0").Arg("-loop", "0").Arg("-preset", "photo").Arg("-quality", "20")
	c.Arg("-f", "gif")

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

func (c *Ffmpeg) setInput() error {
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

func (c *Ffmpeg) getOutput() (string, error) {
	if c.output != nil {
		return "-", nil
	} else {
		return "", errors.New("Undefined output")
	}
}

// ffmpeg -i uer.mp4 -r 10 -vf "fps=fps=20,scale=512:512:force_original_aspect_ratio=decrease,format=rgba,pad=512:512:-1:-1:color=#00000000" out.gif

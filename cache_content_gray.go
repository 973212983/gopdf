package gopdf

import (
	"fmt"
	"io"
)

const grayTypeFill = "g"
const grayTypeStroke = "G"

type cacheContentGray struct {
	grayType string
	scale    float64
}

func (c *cacheContentGray) Clone(f func() *GoPdf) ICacheContent {
	cl := new(cacheContentGray)
	cl.grayType = c.grayType
	cl.scale = c.scale
	return cl
}

func (c *cacheContentGray) write(w io.Writer, protection *PDFProtection) error {
	fmt.Fprintf(w, "%.2f %s\n", c.scale, c.grayType)
	return nil
}

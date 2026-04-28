package gopdf

import (
	"fmt"
	"io"
)

type cacheContentLineWidth struct {
	width float64
}

func (c *cacheContentLineWidth) Clone(f func() *GoPdf) ICacheContent {
	cl := new(cacheContentLineWidth)
	cl.width = c.width
	return cl
}

func (c *cacheContentLineWidth) write(w io.Writer, protection *PDFProtection) error {
	fmt.Fprintf(w, "%.2f w\n", c.width)
	return nil
}

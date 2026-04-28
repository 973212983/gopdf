package gopdf

import (
	"fmt"
	"io"
)

type cacheContentTextColorCMYK struct {
	c, m, y, k uint8
}

func (c cacheContentTextColorCMYK) write(w io.Writer, protection *PDFProtection) error {
	fmt.Fprintf(w, "%.2f %.2f %.2f %.2f %s\n", float64(c.c)/100, float64(c.m)/100, float64(c.y)/100, float64(c.k)/100, colorTypeFillCMYK)
	return nil
}

func (c cacheContentTextColorCMYK) equal(obj ICacheColorText) bool {
	cmyk, ok := obj.(cacheContentTextColorCMYK)
	if !ok {
		return false
	}

	return c.c == cmyk.c && c.m == cmyk.m && c.y == cmyk.y && c.k == cmyk.k
}

func (c cacheContentTextColorCMYK) Clone(f func() *GoPdf) ICacheContent {
	cl := new(cacheContentTextColorCMYK)
	cl.c = c.c
	cl.m = c.m
	cl.y = c.y
	cl.k = c.k
	return *cl
}

func (c cacheContentTextColorCMYK) CloneText() ICacheColorText {
	cl := new(cacheContentTextColorCMYK)
	cl.c = c.c
	cl.m = c.m
	cl.y = c.y
	cl.k = c.k
	return *cl
}

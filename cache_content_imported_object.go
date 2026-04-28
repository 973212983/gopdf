package gopdf

import (
	"fmt"
	"io"
)

type cacheContentImportedTemplate struct {
	pageHeight float64
	tplName    string
	scaleX     float64
	scaleY     float64
	tX         float64
	tY         float64
}

func (c *cacheContentImportedTemplate) Clone(f func() *GoPdf) ICacheContent {
	cl := new(cacheContentImportedTemplate)
	cl.pageHeight = c.pageHeight
	cl.tplName = c.tplName
	cl.scaleX = c.scaleX
	cl.scaleY = c.scaleY
	cl.tX = c.tX
	cl.tY = c.tY
	return cl
}

func (c *cacheContentImportedTemplate) write(w io.Writer, protection *PDFProtection) error {
	c.tY += c.pageHeight
	fmt.Fprintf(w, "q 0 J 1 w 0 j 0 G 0 g q %.4F 0 0 %.4F %.4F %.4F cm %s Do Q Q\n", c.scaleX, c.scaleY, c.tX, c.tY, c.tplName)
	return nil
}

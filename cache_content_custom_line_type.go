package gopdf

import (
	"fmt"
	"io"
)

type cacheContentCustomLineType struct {
	dashArray []float64
	dashPhase float64
}

func (c *cacheContentCustomLineType) Clone(f func() *GoPdf) ICacheContent {
	cl := new(cacheContentCustomLineType)
	cl.dashArray = make([]float64, len(c.dashArray))
	copy(cl.dashArray, c.dashArray)
	cl.dashPhase = c.dashPhase
	return cl
}

func (c *cacheContentCustomLineType) write(w io.Writer, protection *PDFProtection) error {
	fmt.Fprintf(w, "%0.2f %0.2f d\n", c.dashArray, c.dashPhase)
	return nil
}

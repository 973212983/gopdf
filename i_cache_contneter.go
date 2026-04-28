package gopdf

import (
	"io"
)

type ICacheContent interface {
	write(w io.Writer, protection *PDFProtection) error
	Clone(f func() *GoPdf) ICacheContent
}

package gopdf

type ICacheColorText interface {
	ICacheContent
	equal(obj ICacheColorText) bool
	CloneText() ICacheColorText
}

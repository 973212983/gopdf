package gopdf

// Current current state
type Current struct {
	setXCount int //many times we go func SetX()
	X         float64
	Y         float64

	//font
	IndexOfFontObj int
	CountOfFont    int
	CountOfL       int

	FontSize      float64
	FontStyle     int // Regular|Bold|Italic|Underline
	FontFontCount int
	FontType      int // CURRENT_FONT_TYPE_IFONT or  CURRENT_FONT_TYPE_SUBSET

	IndexOfColorSpaceObj int
	CountOfColorSpace    int

	CharSpacing float64

	FontISubset *SubsetFontObj // FontType == CURRENT_FONT_TYPE_SUBSET

	//page
	IndexOfPageObj int

	//img
	CountOfImg int
	//cache of image in pdf file
	ImgCaches map[int]ImageCache

	//text color mode
	txtColorMode string //color, gray

	//text color
	txtColor ICacheColorText

	//text grayscale
	grayFill float64
	//draw grayscale
	grayStroke float64

	lineWidth float64

	//current page size
	pageSize *Rect

	//current trim box
	trimBox *Box

	sMasksMap       SMaskMap
	extGStatesMap   ExtGStatesMap
	transparency    *Transparency
	transparencyMap TransparencyMap
}

func (c *Current) Clone(f func() *GoPdf) *Current {
	cl := new(Current)
	cl.setXCount = c.setXCount
	cl.X = c.X
	cl.Y = c.Y
	cl.IndexOfFontObj = c.IndexOfFontObj
	cl.CountOfFont = c.CountOfFont
	cl.CountOfL = c.CountOfL
	cl.FontSize = c.FontSize
	cl.FontStyle = c.FontStyle
	cl.FontFontCount = c.FontFontCount
	cl.FontType = c.FontType
	cl.IndexOfColorSpaceObj = c.IndexOfColorSpaceObj
	cl.CountOfColorSpace = c.CountOfColorSpace
	cl.CharSpacing = c.CharSpacing
	if c.FontISubset != nil {
		cl.FontISubset = c.FontISubset.clone(f).(*SubsetFontObj)
	}

	cl.IndexOfPageObj = c.IndexOfPageObj
	cl.CountOfImg = c.CountOfImg
	cl.ImgCaches = make(map[int]ImageCache)
	for k, chache := range c.ImgCaches {
		cl.ImgCaches[k] = chache
	}
	cl.txtColorMode = c.txtColorMode
	if c.txtColor != nil {
		cl.txtColor = c.txtColor.CloneText()
	}
	cl.grayFill = c.grayFill
	cl.grayStroke = c.grayStroke
	cl.lineWidth = c.lineWidth
	if c.pageSize != nil {
		cl.pageSize = c.pageSize.Clone()
	}
	if c.trimBox != nil {
		cl.trimBox = c.trimBox.Clone()
	}
	cl.sMasksMap = *c.sMasksMap.Clone(f)
	cl.extGStatesMap = *c.extGStatesMap.Clone(f)
	if c.transparency != nil {
		tr := c.transparency.Clone()
		cl.transparency = &tr
	}
	cl.transparencyMap = *c.transparencyMap.Clone(f)
	return cl
}

func (c *Current) setTextColor(color ICacheColorText) {
	c.txtColor = color
}

func (c *Current) textColor() ICacheColorText {
	return c.txtColor
}

// ImageCache is metadata for caching images.
type ImageCache struct {
	Path  string //ID or Path
	Index int
	Rect  *Rect
}

func (c *ImageCache) Clone() *ImageCache {
	cl := new(ImageCache)
	cl.Path = c.Path
	cl.Index = c.Index
	if c.Rect != nil {
		cl.Rect = c.Rect.Clone()
	}
	return cl
}

package gopdf

type imgInfo struct {
	w, h int
	//src              string
	formatName       string
	colspace         string
	bitsPerComponent string
	filter           string
	decodeParms      string
	trns             []byte
	smask            []byte
	smarkObjID       int
	pal              []byte
	deviceRGBObjID   int
	data             []byte
}

func (info imgInfo) Clone() imgInfo {
	cl := imgInfo{
		w:                info.w,
		h:                info.h,
		formatName:       info.formatName,
		colspace:         info.colspace,
		bitsPerComponent: info.bitsPerComponent,
		filter:           info.filter,
		decodeParms:      info.decodeParms,
		trns:             info.trns,
		smask:            info.smask,
		smarkObjID:       info.smarkObjID,
		pal:              info.pal,
		deviceRGBObjID:   info.deviceRGBObjID,
		data:             info.data,
	}
	return cl
}

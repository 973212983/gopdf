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
		trns:             make([]byte, len(info.trns)),
		smask:            make([]byte, len(info.smask)),
		smarkObjID:       info.smarkObjID,
		pal:              make([]byte, len(info.pal)),
		deviceRGBObjID:   info.deviceRGBObjID,
		data:             make([]byte, len(info.data)),
	}
	copy(cl.trns, info.trns)
	copy(cl.smask, info.smask)
	copy(cl.pal, info.pal)
	copy(cl.data, info.data)
	return cl
}

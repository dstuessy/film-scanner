package tiff

import (
	"log"
	"unsafe"

	"gocv.io/x/gocv"
)

const (
	EndianII = 0x4949
	EndianMM = 0x4D4D
)

func encodeHex(num uint16) [2]byte {
	return [2]byte{byte(num), byte(num >> 8)}
}

func encodeHex32(num uint32) [4]byte {
	return [4]byte{byte(num), byte(num >> 8), byte(num >> 16), byte(num >> 24)}
}

type TiffHeader struct {
	Endian         uint16
	TiffIdentifier uint16
	IfdOffset      uint32
}

func (th *TiffHeader) Encode() []byte {
	buf := make([]byte, 8)
	en := encodeHex(th.Endian)
	copy(buf[0:2], en[:])
	id := encodeHex(th.TiffIdentifier)
	copy(buf[2:4], id[:])
	o := encodeHex32(th.IfdOffset)
	copy(buf[4:8], o[:])
	return buf
}

func (th *TiffHeader) Len() int {
	return int(unsafe.Sizeof(th.Endian) + unsafe.Sizeof(th.TiffIdentifier) + unsafe.Sizeof(th.IfdOffset))
}

type TiffIfd struct {
	NumFields     uint16       // [2]byte
	Fields        []*TiffField // [TiffIfd.NumFields*12]byte
	NextIfdOffset uint32       // [4]byte
}

func (ifd *TiffIfd) Encode() []byte {
	fieldLen := 12
	buf := make([]byte, 2+uint16(fieldLen)*ifd.NumFields+4)
	nf := encodeHex(ifd.NumFields)
	copy(buf[0:2], nf[:])

	for i, field := range ifd.Fields {
		start := 2 + 12*i
		end := start + field.Len()
		copy(buf[start:end], field.Encode()) // 2 is the offset from which the NumFields ended
	}

	ifdOffsetStart := 2 + fieldLen*int(ifd.NumFields) // 2 is the offset from which the NumFields ended
	ifdOffsetEnd := ifdOffsetStart + 4
	nextIfd := encodeHex32(ifd.NextIfdOffset)
	copy(buf[ifdOffsetStart:ifdOffsetEnd], nextIfd[:])

	return buf
}

func (ifd *TiffIfd) Len() int {
	return int(unsafe.Sizeof(ifd.NumFields)) + len(ifd.Fields)*12 + int(unsafe.Sizeof(ifd.NextIfdOffset))
}

type TiffField struct {
	Tag         uint16 // [2]byte
	Type        uint16 // [2]byte
	Count       uint32 // [4]byte
	Value       uint32 // [4]byte represents value if spec allows, otherwise it represents the offset at which TiffField.OffsetValue os located
	OffsetValue []byte // bytes for the value in a location represented by TiffField.Value. These are not encoded in the tiff field itself, and also omitted from the TiffField.Len() calculation
}

func (tf *TiffField) Encode() []byte {
	field := make([]byte, 12)
	tag := encodeHex(tf.Tag)
	copy(field[0:2], tag[:])
	tp := encodeHex(tf.Type)
	copy(field[2:4], tp[:])
	c := encodeHex32(tf.Count)
	copy(field[4:8], c[:])
	v := encodeHex32(tf.Value)
	copy(field[8:12], v[:])
	return field
}

func (tf *TiffField) Len() int {
	return int(unsafe.Sizeof(tf.Tag) + unsafe.Sizeof(tf.Type) + unsafe.Sizeof(tf.Count) + unsafe.Sizeof(tf.Value))
}

func EncodeTiff(img gocv.Mat) ([]byte, error) {
	imgData := img.ToBytes()

	// HEADER
	h := TiffHeader{
		Endian:         EndianII,
		TiffIdentifier: 0x2a,
		IfdOffset:      0x8,
	}

	fields := make([]*TiffField, 0)

	// ImageWidth
	fields = append(fields, &TiffField{
		Tag:   0x100,
		Type:  0x4,
		Count: 0x1,
		Value: uint32(img.Cols()), // horizontal Length
	})

	// ImageLength
	fields = append(fields, &TiffField{
		Tag:   0x101,
		Type:  0x4,
		Count: 0x1,
		Value: uint32(img.Rows()), // vertical Length
	})

	// BitsPerSample
	fields = append(fields, &TiffField{
		Tag:   0x102,
		Type:  0x3,
		Count: 0x1,
		Value: 0x8, // 8 bits per channel
	})

	// Compression
	fields = append(fields, &TiffField{
		Tag:   0x103,
		Type:  0x3,
		Count: 0x1,
		Value: 0x1, // no compression
	})

	// PhotometricInterpretation
	fields = append(fields, &TiffField{
		Tag:   0x106,
		Type:  0x3,
		Count: 0x1,
		Value: 0x2, // ZeroIsBlack for full rgb
	})

	// StripOffsets
	stripOffsets := &TiffField{
		Tag:         0x111,
		Type:        0x4,
		Count:       0x1,
		OffsetValue: imgData,
	}
	fields = append(fields, stripOffsets)

	// SamplesPerPixel
	fields = append(fields, &TiffField{
		Tag:   0x115,
		Type:  0x3,
		Count: 0x1,
		Value: 0x3, // 3 channels per pixel
	})

	// RowsPerStrip
	fields = append(fields, &TiffField{
		Tag:   0x116,
		Type:  0x4,
		Count: 0x1,
		Value: uint32(img.Rows()), // number of rows per strip
	})

	// StripByteCounts
	stripByteCounts := &TiffField{
		Tag:   0x117,
		Type:  0x4,
		Count: 0x1,
		Value: uint32(len(imgData)),
	}
	fields = append(fields, stripByteCounts)

	// PlanarConfiguration
	fields = append(fields, &TiffField{
		Tag:   0x11C,
		Type:  0x3,
		Count: 0x1,
		Value: 0x1, // how color channels are stored per pixel
	})

	// IFD 0
	ifd := TiffIfd{
		NumFields:     uint16(len(fields)),
		Fields:        fields,
		NextIfdOffset: 0x0, // this is the last ifd
	}

	stripOffsets.Value = uint32(h.Len() + ifd.Len())

	buf := make([]byte, 0)
	buf = append(buf, h.Encode()...)
	buf = append(buf, ifd.Encode()...)
	buf = append(buf, stripOffsets.OffsetValue...)

	return buf, nil
}

package records

import (
	//"image"
	//"bytes"
	"io"
	//"github.com/leotaku/mobi/jfif"
)

/*
type ImageRecord struct {
	img image.Image
}

func NewImageRecord(img image.Image) ImageRecord {
	return ImageRecord{
		img: img,
	}
}

func (r ImageRecord) Write(w io.Writer) error {
	return jfif.Encode(w, r.img, nil)
}
*/

type BlobRecord struct {
    data []byte
}

func NewBlobRecord(data []byte) BlobRecord {
    return BlobRecord{
        data: data,
    }
}

func (r BlobRecord) Write(w io.Writer) error {
    _, err := w.Write(r.data)
    return err
}



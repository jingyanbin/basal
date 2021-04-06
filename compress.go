package base

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"io/ioutil"
)

func ZLibCompress(data []byte) []byte {
	var in bytes.Buffer
	w := zlib.NewWriter(&in)
	w.Write(data)
	w.Close()
	return in.Bytes()
}

func ZLibUnCompress(data []byte) ([]byte, error) {
	b := bytes.NewReader(data)
	r, err := zlib.NewReader(b)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	unData, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return unData, err
}

func GZipCompress(data []byte) []byte {
	var in bytes.Buffer
	w := gzip.NewWriter(&in)
	w.Write(data)
	w.Close()
	return in.Bytes()
}

func GZipUnCompress(data []byte) ([]byte, error) {
	b := bytes.NewReader(data)
	r, err := gzip.NewReader(b)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	unData, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return unData, err
}

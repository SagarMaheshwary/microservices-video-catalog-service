package helper

import (
	"crypto/rand"
	"encoding/base32"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"google.golang.org/grpc/metadata"
)

func GetRootDir() string {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))

	return filepath.Dir(d)
}

func GetGRPCMetadataValue(md metadata.MD, k string) (string, bool) {
	v := md.Get(k)

	if len(v) == 0 {
		return "", false
	}

	return v[0], true
}

func UniqueString(length int) string {
	b := make([]byte, 32)

	rand.Read(b)

	return base32.StdEncoding.EncodeToString(b)[:length]
}

func Slug(str string) string {
	return strings.ReplaceAll(strings.ToLower(str), " ", "-")
}

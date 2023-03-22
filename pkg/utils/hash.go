package utils

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
)

func ReadAndCalcHash(path string) ([]byte, string) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	buf, err := io.ReadAll(file)
	if err != nil {
		return nil, ""
	}
	md5hash := md5.New()
	//sha1hash := sha1.New()
	//sha256hash := sha256.New()

	if _, err := io.Copy(md5hash, file); err != nil {
		panic(err)
	}
	if _, err := file.Seek(0, 0); err != nil {
		panic(err)
	}
	//if _, err := io.Copy(sha1hash, file); err != nil {
	//	panic(err)
	//}
	//if _, err := file.Seek(0, 0); err != nil {
	//	panic(err)
	//}
	//if _, err := io.Copy(sha256hash, file); err != nil {
	//	panic(err)
	//}

	md5sum := fmt.Sprintf("%x", md5hash.Sum(nil))
	//sha1sum := fmt.Sprintf("%x", sha1hash.Sum(nil))
	//sha256sum := fmt.Sprintf("%x", sha256hash.Sum(nil))

	// fmt.Printf("MD5: %s\nSHA1: %s\nSHA256: %s\n", md5sum, sha1sum, sha256sum)
	return buf, md5sum
}

package backup

import (
	"crypto/md5"
	"fmt"
	"io"
	//	"io/fs"
	"os"
	"path/filepath"
)

func DirHash(path string) (string, error) {
	hash := md5.New()
	// err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
	// 	fmt.Println("1")
	// 	if err != nil {
	// 		return err
	// 	}
	// 	io.WriteString(hash, path)
	// 	info, err := d.Info()
	// 	fmt.Println("2")
	// 	if err != nil {
	// 		return err
	// 	}
	// 	fmt.Fprintf(hash, "%v", info.IsDir())
	// 	fmt.Fprintf(hash, "%v", info.ModTime())
	// 	fmt.Fprintf(hash, "%v", info.Mode())
	// 	fmt.Fprintf(hash, "%v", info.Name())
	// 	fmt.Fprintf(hash, "%v", info.Size())
	// 	return nil
	// })

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		fmt.Println("1")

		// fmt.Println("isDir=", info.IsDir())
		// fmt.Println("ModTime=", info.ModTime())
		// fmt.Println("Mode=", info.Mode())
		// fmt.Println("Name=", info.Name())
		// fmt.Println("Size=", info.Size())

		if err != nil {
			return err
		}
		fmt.Println("2")
		io.WriteString(hash, path)

		fmt.Fprintf(hash, "%v", info.IsDir())
		fmt.Fprintf(hash, "%v", info.ModTime())
		fmt.Fprintf(hash, "%v", info.Mode())
		fmt.Fprintf(hash, "%v", info.Name())
		fmt.Fprintf(hash, "%v", info.Size())
		return nil
	})

	fmt.Println("3")
	if err != nil {
		return "", err
	}
	fmt.Println("4")
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

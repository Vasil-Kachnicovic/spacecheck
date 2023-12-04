package main

import (
	"os"
	"log"
	"io/fs"
	"fmt"
	"strings"
)

type DirSize struct {
	Name string
	Size int64
}

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}

	cwdFs := os.DirFS(cwd)

	sizes, err := IterDirs(cwdFs)
	if err != nil {
		log.Fatalln(err)
	}
	
	PrintSizes(sizes)
}

func IterDirs(cwdFs fs.FS) ([]*DirSize ,error) {
	dirs, err := cwdFs.(fs.ReadDirFS).ReadDir(".")
	if err != nil {
		return nil, err
	}

	sizes := make([]*DirSize, 0, len(dirs))

	for _, dir := range dirs {
		fi, err := dir.Info()
		if err != nil {
			fmt.Println(err)
			continue
		}
		if !(fi.Mode().IsRegular() || fi.IsDir()) {
			continue
		}
		size, err := CheckSize(cwdFs, dir.Name())
		if err != nil {
			return sizes, err
		}

		sizes = append(sizes, &DirSize{dir.Name(), size})
	}

	return sizes, nil
}

func CheckSize(cwdFs fs.FS, dir string) (int64, error) {
	var size int64 = 0
	var err error

	err = fs.WalkDir(cwdFs, dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			if d == nil {
				return err
			}
			fmt.Println(err)
		}

		// Skip directories
		if d.IsDir() {
			return nil
		}

		fi, err := d.Info()
		if err != nil {
			fmt.Println(err)
		}

		// Skip non-regular files
		if !fi.Mode().IsRegular() {
			return nil
		}

		size += fi.Size()

		return nil
	})
	
	return size, err
}

func PrintSizes(sizes []*DirSize) {
	var total int64 = 0

	biggestLen := 0

	for _, size := range sizes {
		if len(size.Name) > biggestLen {
			biggestLen = len(size.Name)
		}
	}

	for _, size := range sizes {
		fmt.Println(size.Format(biggestLen))
		total += size.Size
	}

	ts := &DirSize{"Total is", total}
	fmt.Printf("\n%v\n", ts.Format(biggestLen))
}

func (ds *DirSize) Format(biggestLen int) string {
	out := ds.Name
	nSpaces := biggestLen - len(out)
	out += strings.Repeat(" ", nSpaces)

	units := []string{"B", "K", "M", "G", "T", "P"}

	for i, unit := range units {
		if i+1 == len(units) {
			out += fmt.Sprintf(" %v%v", ds.Size >> (10*i), unit)
			break
		}
		formated_value := ds.Size >> (10*(i+1))
		if formated_value < 1 {
			out += fmt.Sprintf(" %6v%v", ds.Size >> (10*i), unit)
			break
		}
		if formated_value < 999 {
			out += fmt.Sprintf(" %6v%v", formated_value, units[i+1])
			break
		}
	}
	return out
}
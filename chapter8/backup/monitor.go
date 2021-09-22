package backup

import (
	"log"
	"path/filepath"
	"time"
)

type Monitor struct {
	Paths       map[string]string
	Archiver    Archiver
	Destination string
}

func (m *Monitor) Now() (int, error) {
	var counter int
	for path, lastHash := range m.Paths {
		log.Println(path)
		newHash, err := DirHash(path)
		if err != nil {
			return 0, err
		}
		if newHash != lastHash {
			err := m.act(path)
			if err != nil {
				return counter, err
			}
			m.Paths[path] = newHash // update hash
			counter++
		}
	}
	return counter, nil
}

// TODO: change extension .zip to arbitrary extension with argument
func (m *Monitor) act(path string) error {
	dirname := filepath.Base(path)
	filename := m.Archiver.DestFmt()(time.Now().UnixNano())
	log.Println(dirname)
	log.Println(filename)
	log.Println(filepath.Join(dirname, filename))
	return m.Archiver.Archive(path, filepath.Join(m.Destination, dirname, filename))
}

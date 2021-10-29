package daas

import (
	"crypto/md5"
	"fmt"
	"os"
	"strings"
	"time"
)

type CSV struct {
	column    []string
	name      string
	delimiter string
	raw       []byte
}

func (file *CSV) AddHeader(data ...string) error {
	file.delimiter = "|"
	file.column = data
	head := []byte(strings.Join(data, file.delimiter))
	head = append(head, byte('\n'))
	file.raw = append(head, file.raw...)
	return nil
}

func (file *CSV) Appendln(data ...string) error {
	file.raw = append(file.raw, []byte(strings.Join(data, file.delimiter))...)
	file.raw = append(file.raw, byte('\n'))
	return nil
}

func (file *CSV) Save() error {
	h := md5.New()
	file.name = fmt.Sprintf("%s/%x", os.TempDir(), h.Sum([]byte(fmt.Sprint(time.Now().Unix()))))
	rt, err := os.Create(file.name)
	if err != nil {
		return err
	}
	defer rt.Close()

	_, err = rt.Write(file.raw)
	if err != nil {
		return err
	}
	return nil
}

func (file *CSV) Remove() error {
	if file.name == "" {
		return nil
	}
	return os.Remove(file.name)
}

func (file *CSV) GetFilePath() string {
	return file.name
}

func (file *CSV) GetColumns() string {
	return strings.Join(file.column, ",")
}

func (file *CSV) GetDelimiter() string {
	return file.delimiter
}

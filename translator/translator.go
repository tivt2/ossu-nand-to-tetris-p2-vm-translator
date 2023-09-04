package translator

import (
	"github.com/tivt2/vm-translator/fio"
	"github.com/tivt2/vm-translator/parser"
)

func Translate(filePath string) {
	reader := fio.NewReader(filePath)
	go reader.Read()

	parser := parser.New(reader.ReadChan)
	go parser.Parse()

	writer := fio.NewWriter(filePath)
	writer.Wg.Add(1)
	go writer.Write(parser.WriteChan)
	writer.Wg.Wait()
}

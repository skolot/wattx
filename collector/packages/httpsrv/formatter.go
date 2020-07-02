package httpsrv

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"wattx/collector/packages/api"
)

type FormaterFunc func(data []api.Data) (*bytes.Buffer, error)

const (
	delim string = "\t"
	crlf         = "\n"

	plainFormat string = "plain"
	jsonFormat  string = "json"
	csvFormat   string = "csv"

	defaultFormat = plainFormat

	contentTypeHeader = "Content-Type"
)

var (
	header []string = []string{"RANK", "NAME", "FULLNAME", "PRICE", "CURRENCY"}

	formatters map[string]FormaterFunc = map[string]FormaterFunc{
		plainFormat: formatPlain,
		jsonFormat:  formatJSON,
		csvFormat:   formatCSV,
	}

	contentTypes map[string]string = map[string]string{
		plainFormat: "text/plain",
		jsonFormat:  "application/json",
		csvFormat:   "text/csv",
	}
)

func formatPlain(data []api.Data) (*bytes.Buffer, error) {
	strs := []string{
		strings.Join(header, delim),
	}

	format := "%d" + delim + "%s" + delim + "%s" + delim + "%f" + delim + "%s"

	for _, d := range data {
		line := fmt.Sprintf(format, d.Rank, d.Name, d.FullName, d.Price, d.Currency)
		strs = append(strs, line)
	}

	return bytes.NewBufferString(strings.Join(strs, crlf)), nil
}

func formatJSON(data []api.Data) (*bytes.Buffer, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(b), nil
}

func formatCSV(data []api.Data) (*bytes.Buffer, error) {
	buf := &bytes.Buffer{}
	wbuf := bufio.NewWriter(buf)
	wcsv := csv.NewWriter(wbuf)

	if err := wcsv.Write(header); err != nil {
		return nil, err
	}

	for _, d := range data {
		record := []string{
			strconv.Itoa(d.Rank),
			d.Name,
			d.FullName,
			strconv.FormatFloat(float64(d.Price), 'f', -1, 32),
			d.Currency,
		}

		if err := wcsv.Write(record); err != nil {
			return nil, err
		}
	}

	wcsv.Flush()
	if err := wcsv.Error(); err != nil {
		return nil, err
	}

	return buf, nil
}

func writeFormattedData(w http.ResponseWriter, format string, data []api.Data) error {
	formatFunc, ok := formatters[format]
	if !ok {
		return errors.New("unknow format: " + format)
	}

	buf, err := formatFunc(data)
	if err != nil {
		return err
	}

	w.Header().Set(contentTypeHeader, contentTypes[format])
	buf.WriteTo(w)

	return nil
}

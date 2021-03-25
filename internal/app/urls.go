package app

import (
	"strings"

	"github.com/vimoppa/turl.to/internal/storage"
)

// RecordsItem is a single url source records exported.
type RecordsItem struct {
	Hash    string `json:"hash,omitempty"`
	LongURL string `json:"long_url,omitempty"`
}

// GetAllRecords reads and prepares records in a consumable format.
func GetAllRecords(s storage.Accessor) ([]RecordsItem, error) {
	records := make([]RecordsItem, 0)

	results, err := s.ReadAll()
	if err != nil {
		return nil, err
	}

	for _, v := range results {
		r := strings.Split(v, " ")
		item := RecordsItem{
			Hash:    r[0],
			LongURL: r[1],
		}
		records = append(records, item)
	}

	return records, nil
}

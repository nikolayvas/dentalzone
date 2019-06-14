package mssql

import (
	m "dental_hub/models"
	"io"
)

/*
import (
	ex "dental_hub/exceptions"
	m "dental_hub/models"
	u "dental_hub/utils"
	"encoding/json"
	"sync"
	"time"

	"github.com/gocql/gocql"
	uuid "github.com/satori/go.uuid"
)
*/

// InsertImage ...
func (r *Repository) InsertImage(string /*patient*/, io.Reader /*imageFileLocation*/, []string /*tags*/, string /*file name*/, int64 /*file size*/) error {
	return nil
}

// GetImagesByTags ...
func (r *Repository) GetImagesByTags(string /*patient*/, []string /*tags*/) (*[]m.FileDetails, error) {
	return nil, nil
}

// GetImage ...
func (r *Repository) GetImage(string /*patient*/, string /*S3/minio image id*/) (io.Reader /*file*/, error) {
	return nil, nil
}

// GetTagsByPatient ...
func (r *Repository) GetTagsByPatient(string /*patient*/) ([]string, error) {
	return nil, nil
}

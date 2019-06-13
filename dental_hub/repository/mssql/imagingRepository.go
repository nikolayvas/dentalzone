package mssql

import "io"

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
func (r *Repository) InsertImage(string /*patient*/, io.Reader /*imageFileLocation*/, []string /*tags*/, int64 /*file size*/) error {
	return nil
}

// GetImageIdsByTags ...
func (r *Repository) GetImageIdsByTags(string /*patient*/, []string /*tags*/) ([]string /*S3/minio file id's*/, error) {
	var empty []string
	return empty, nil
}

// GetImage ...
func (r *Repository) GetImage(string /*patient*/, string /*S3/minio image id*/) (io.Reader /*file*/, error) {
	return nil, nil
}

package mssql

import (
	m "dental_hub/models"
	"dental_hub/utils"
	"fmt"
	"io"
	"math/rand"
	"strings"
	"sync"
	"time"

	minio "github.com/minio/minio-go/v6"
	uuid "github.com/satori/go.uuid"
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
func (r *Repository) InsertImage(patientID string, file io.Reader, tags []string, fileName string, fileSize int64) error {

	bucket := strings.ToLower(patientID)

	found, err := r.MinioClient.BucketExists(bucket)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if !found {
		err = r.MinioClient.MakeBucket(bucket, "")
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Println("Successfully created '", bucket, "' busket")
	}

	var uid = utils.UniqueID()

	n, err := r.MinioClient.PutObject(bucket, uid, file, fileSize, minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("Successfully uploaded bytes: ", n)

	for _, tag := range tags {

		_, err = r.Connection.Exec("exec [InsertFileAndTags_SP] ?, ?, ?, ?, ?, ?",
			patientID,
			tag,
			tag,
			uid,
			fileName,
			0)

		if err != nil {
			return err
		}
	}

	return nil
}

// GetImagesByTags ...
func (r *Repository) GetImagesByTags(patientID string, tags []string) (*[]m.FileDetails, error) {
	if len(tags) == 0 {
		return nil, nil
	}

	var intersect []m.FileDetails

	for i, tag := range tags {
		var id string
		var name string
		var size int64

		filesByTag := make([]m.FileDetails, 0)

		rows, err := r.Connection.Query(`SELECT imageid, filename, filesize FROM tag t inner join tags t2 on t2.id = t.[TagKey] WHERE t2.[PatientId] = ? and t2.[TagKey] = ?`, patientID, tag)
		if err != nil {
			return nil, err
		}

		defer rows.Close()
		for rows.Next() {

			err := rows.Scan(
				&id,
				&name,
				&size)

			if err != nil {
				return nil, err
			}

			filesByTag = append(filesByTag, m.FileDetails{
				ID:   id,
				Name: name,
				Size: size,
			})
		}

		if i > 0 {
			intersect = Intersect(intersect, filesByTag)
			if len(intersect) == 0 {
				break
			}
		} else {
			intersect = filesByTag
		}
	}

	return &intersect, nil
}

// GetImage ...
func (r *Repository) GetImage(patientID string, imageID string) (io.Reader /*file*/, error) {
	object, err := r.MinioClient.GetObject(patientID, imageID, minio.GetObjectOptions{})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return object, nil
}

// GetTagsByPatient ...
func (r *Repository) GetTagsByPatient(patientID string) ([]string, error) {
	tags := make([]string, 0)

	rows, err := r.Connection.Query(`select [TagKey] from tags where [PatientId] = ?`, patientID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {

		var tag string
		err := rows.Scan(
			&tag)

		if err != nil {
			return nil, err
		}

		tags = append(tags, tag)
	}

	return tags, nil
}

// Intersect ...
func Intersect(a []m.FileDetails, b []m.FileDetails) []m.FileDetails {
	set := make([]m.FileDetails, 0)
	hash := make(map[string]bool)

	for _, file := range a {
		hash[file.ID] = true
	}

	for _, file := range b {
		if _, found := hash[file.ID]; found {
			set = append(set, file)
		}
	}

	return set
}

// Insert1000TestImagesWith100Tags ...
func (r *Repository) Insert1000TestImagesWith100Tags(patientID string, tags []string) error {

	fmt.Println("Insert1000TestImagesWith100Tags started: ", time.Now())

	var repeat = 100
	var randomTagsPerImage = 5
	var wg sync.WaitGroup
	wg.Add(repeat)

	for i := 0; i < repeat; i++ {
		go func() {

			defer wg.Done()

			for j := 0; j < 1000; j++ {

				var id uuid.UUID
				id, err := uuid.NewV4()
				if err != nil {

				}

				fileName := id.String()

				for j := 0; j < randomTagsPerImage; j++ {

					var tag = tags[rand.Intn(len(tags)-1)]

					_, err = r.Connection.Exec("exec [InsertFileAndTags_SP] ?, ?, ?, ?, ?, ?",
						patientID,
						tag,
						tag,
						fileName,
						fileName,
						0)

					if err != nil {
						fmt.Println(err, time.Now())
					}
				}
			}
		}()
	}

	wg.Wait()

	fmt.Println("Insert1000TestImagesWith100Tags ended: ", time.Now())

	return nil
}

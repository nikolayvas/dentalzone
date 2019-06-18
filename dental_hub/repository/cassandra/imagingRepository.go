package cassandra

import (
	"fmt"
	"io"
	"sync"
	"time"

	m "dental_hub/models"
	"dental_hub/utils"
	"math/rand"

	"github.com/minio/minio-go/v6"
	uuid "github.com/satori/go.uuid"
)

// InsertImage inserts image to minio/S3
func (r *Repository) InsertImage(patientID string, file io.Reader, tags []string, fileName string, fileSize int64) error {

	found, err := r.MinioClient.BucketExists(patientID)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if !found {
		err = r.MinioClient.MakeBucket(patientID, "")
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Println("Successfully created '", patientID, "' busket")
	}

	var uid = utils.UniqueID()

	n, err := r.MinioClient.PutObject(patientID, uid, file, fileSize, minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("Successfully uploaded bytes: ", n)

	for _, tag := range tags {
		if err = r.Session.Query(`INSERT INTO tags(patientid, tagkey) Values(?, ?)`,
			patientID,
			tag,
		).Exec(); err != nil {
			return err
		}
	}

	for _, tag := range tags {
		if err = r.Session.Query(`INSERT INTO tag(patientid, tagkey, tagvalue, imageid, filename, filesize) Values(?, ?, ?, ?, ?, ?)`,
			patientID,
			tag,
			tag,
			uid,
			fileName,
			fileSize,
		).Exec(); err != nil {
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

		iter := r.Session.Query(`SELECT imageid, filename, filesize FROM tag WHERE patientid = ? and tagkey = ?`, patientID, tag).Iter()

		for iter.Scan(&id, &name, &size) {
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

	/*
		allFilesByTags := make([][]m.FileDetails, len(tags))

		var wg sync.WaitGroup
		wg.Add(len(tags))

		for i, tag := range tags {
			go func(i int, tag string) {

				var id string
				var name string
				var size int64

				defer wg.Done()

				filesByTag := make([]m.FileDetails, 0)

				iter := r.Session.Query(`SELECT imageid, filename, filesize FROM tag WHERE patientid = ? and tagkey = ?`, patientId, tag).Iter()

				for iter.Scan(&id, &name, &size) {
					filesByTag = append(filesByTag, m.FileDetails{
						ID:   id,
						Name: name,
						Size: size,
					})
				}

				allFilesByTags[i] = filesByTag
			}(i, tag)
		}

		wg.Wait()

		intersect := allFilesByTags[0]

		for _, files := range allFilesByTags {
			intersect = Intersect(intersect, files)
		}

	*/

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

	var tag string
	tags := make([]string, 0)

	iter := r.Session.Query(`SELECT tagkey from tags WHERE patientid=?`, patientID).Iter()

	for iter.Scan(&tag) {
		tags = append(tags, tag)
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

	return tags, nil
}

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
					if err := r.Session.Query(`INSERT INTO tags(patientid, tagkey) Values(?, ?)`,
						patientID,
						tag,
					).Exec(); err != nil {
						fmt.Println(err, time.Now())
					}

					if err := r.Session.Query(`INSERT INTO tag(patientid, tagkey, tagvalue, imageid, filename, filesize) Values(?, ?, ?, ?, ?, ?)`,
						patientID,
						tag,
						tag,
						fileName,
						fileName,
						0,
					).Exec(); err != nil {
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

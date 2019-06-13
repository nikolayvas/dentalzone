package cassandra

import (
	"fmt"
	"io"

	"dental_hub/utils"

	"github.com/minio/minio-go/v6"
)

// InsertImage inserts image to minio/S3
func (r *Repository) InsertImage(patientID string, file io.Reader, tags []string, size int64) error {

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

	n, err := r.MinioClient.PutObject(patientID, uid, file, size, minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("Successfully uploaded bytes: ", n)

	for _, tag := range tags {
		if err = r.Session.Query(`INSERT INTO tag(patientid, tagkey, tagvalue, imageid) Values(?, ?, ?, ?)`,
			patientID,
			tag,
			tag,
			uid,
		).Exec(); err != nil {
			return err
		}
	}

	return nil
}

// GetImageIdsByTags ...
func (r *Repository) GetImageIdsByTags(patientID string /*patient*/, tag []string /*tags*/) ([]string /*S3/minio file id's*/, error) {

	var imageID string
	tags := make([]string, 0)

	iter := r.Session.Query(`SELECT imageid FROM tag WHERE patientid=? and tagkey=? and tagvalue = ?`, patientID, "one", "one").Iter()

	for iter.Scan(&imageID) {
		tags = append(tags, imageID)
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

	return tags, nil
}

// GetImage ...
func (r *Repository) GetImage(patientID string, imageID string) (io.Reader /*file*/, error) {

	object, err := r.MinioClient.GetObject( /*patientID*/ "70ba9c89-c1ff-4873-9006-ed5214828be2" /*imageID*/, "4b55f125-a3fc-477b-beff-ac0b3f4f754f", minio.GetObjectOptions{})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return object, nil
}

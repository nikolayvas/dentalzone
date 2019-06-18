package mongodb

import (
	m "dental_hub/models"
	"dental_hub/utils"
	"fmt"
	"io"
	"math/rand"

	"sync"
	"time"

	minio "github.com/minio/minio-go/v6"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo/options"
	uuid "github.com/satori/go.uuid"
)

// Tag read model
type Tag struct {
	PatientID primitive.ObjectID `bson:"patientId"`
	TagKey    string             `bson:"tagKey"`
}

// FileMetadata read model
type FileMetadata struct {
	PatientID primitive.ObjectID `bson:"patientId"`
	TagKey    string             `bson:"tagKey"`
	TagValue  string             `bson:"tagValue"`
	ImageID   string             `bson:"imageId"`
	FileName  string             `bson:"fileName"`
	FileSize  int64              `bson:"fileSize"`
}

//INSERT INTO tag(patientid, tagkey, tagvalue, imageid, filename, filesize

// InsertImage ...
func (r *Repository) InsertImage(patientID string, file io.Reader, tags []string, fileName string, fileSize int64) error {

	ctx, cancel := defaultContextWithTimeout()
	defer cancel()

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

	hex, err := primitive.ObjectIDFromHex(patientID)

	if err != nil {
		return err
	}

	t := true
	opt := options.UpdateOptions{Upsert: &t}

	tagsCollection := r.Client.Database(MongoDbSchema.DatabaseName).Collection(MongoDbSchema.TagsCollection)

	for _, tag := range tags {

		filter := bson.M{"patientId": hex, "tagKey": tag}
		_, err := tagsCollection.UpdateOne(ctx, filter, bson.M{"$set": bson.M{"tagValue": tag}}, &opt)
		if err != nil {
			return err
		}
	}

	tagCollection := r.Client.Database(MongoDbSchema.DatabaseName).Collection(MongoDbSchema.TagCollection)

	for _, tag := range tags {

		filter := bson.M{"patientId": hex, "tagKey": tag, "tagValue": tag, "imageId": fileName}
		_, err := tagCollection.UpdateOne(ctx, filter, bson.M{"$set": bson.M{"fileName": fileName, "fileSize": 0}}, &opt)
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

	hex, err := primitive.ObjectIDFromHex(patientID)

	if err != nil {
		return nil, err
	}

	var intersect []m.FileDetails

	ctx, cancel := contextWithTimeout(500)
	defer cancel()

	tagCollection := r.Client.Database(MongoDbSchema.DatabaseName).Collection(MongoDbSchema.TagCollection)

	for i, tag := range tags {
		filesByTag := make([]m.FileDetails, 0)

		filter := bson.M{"patientId": hex, "tagKey": tag}
		cursor, err := tagCollection.Find(ctx, filter)

		defer cursor.Close(ctx)

		if err != nil {
			return nil, err
		}

		for cursor.Next(ctx) {
			var metadata FileMetadata

			err = cursor.Decode(&metadata)

			if err != nil {
				return nil, err
			}

			filesByTag = append(filesByTag, m.FileDetails{
				ID:   metadata.ImageID,
				Name: metadata.FileName,
				Size: metadata.FileSize,
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

	hex, err := primitive.ObjectIDFromHex(patientID)

	if err != nil {
		return nil, err
	}

	tags := make([]string, 0)

	ctx, cancel := contextWithTimeout(500)
	defer cancel()

	tagsCollection := r.Client.Database(MongoDbSchema.DatabaseName).Collection(MongoDbSchema.TagsCollection)

	filter := bson.M{"patientId": hex}
	cursor, err := tagsCollection.Find(ctx, filter)

	defer cursor.Close(ctx)

	if err != nil {
		return nil, err
	}

	for cursor.Next(ctx) {
		var tag Tag

		err = cursor.Decode(&tag)

		if err != nil {
			return nil, err
		}

		tags = append(tags, tag.TagKey)
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

	ctx, cancel := contextWithTimeout(500)
	defer cancel()

	tagsCollection := r.Client.Database(MongoDbSchema.DatabaseName).Collection(MongoDbSchema.TagsCollection)
	tagCollection := r.Client.Database(MongoDbSchema.DatabaseName).Collection(MongoDbSchema.TagCollection)

	var repeat = 100
	var randomTagsPerImage = 5
	var wg sync.WaitGroup
	wg.Add(repeat)

	hex, err := primitive.ObjectIDFromHex(patientID)

	if err != nil {
		return err
	}

	for i := 0; i < repeat; i++ {
		go func() {

			defer wg.Done()

			for j := 0; j < 1000; j++ {

				var id uuid.UUID
				id, err := uuid.NewV4()
				if err != nil {
					fmt.Println(err, time.Now())
				}

				fileName := id.String()

				for j := 0; j < randomTagsPerImage; j++ {

					var tag = tags[rand.Intn(len(tags)-1)]

					filter := bson.M{"patientId": hex, "tagKey": tag}

					t := true
					opt := options.UpdateOptions{Upsert: &t}

					_, err := tagsCollection.UpdateOne(ctx, filter, bson.M{"$set": bson.M{"tagValue": tag}}, &opt)
					if err != nil {
						fmt.Println(err, time.Now())
					}

					filter = bson.M{"patientId": hex, "tagKey": tag, "tagValue": tag, "imageId": fileName}
					_, err = tagCollection.UpdateOne(ctx, filter, bson.M{"$set": bson.M{"fileName": fileName, "fileSize": 0}}, &opt)
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

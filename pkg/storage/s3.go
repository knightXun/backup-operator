package storage

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	//"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	//"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"io/ioutil"
	"k8s.io/klog"
	"strings"
)

var (
	maxRetries = 6
)
type S3ReadWriter struct {
	S3api             *s3.S3
	BackupDir      string
	Bucket         string
}

func NewS3ReadWriter( endpoint, region, bucket, accessKey, SecretAccessKey, backupdir string) (*S3ReadWriter, error){
	awsConfig := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(accessKey, SecretAccessKey,""),
		Endpoint:         aws.String(endpoint),
		Region:           aws.String(region),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	}

	sess := session.New(awsConfig)
	s3api := s3.New(sess)

	return &S3ReadWriter{
		S3api: s3api,
		BackupDir: backupdir,
		Bucket: bucket,
	}, nil
}

func (readwriter *S3ReadWriter) WriteFile(name string, data string)  error  {
	input := &s3.PutObjectInput{
		Body:   aws.ReadSeekCloser(bytes.NewReader([]byte(data) )),
		Bucket: aws.String(readwriter.Bucket),
		Key:    aws.String(name),
	}

	_, err := readwriter.S3api.PutObject(input)

	return err
}

func (readwriter *S3ReadWriter) ReadFile(name string) ([]byte, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(readwriter.Bucket),
		Key:    aws.String(readwriter.BackupDir + "/" +  name),
	}

	output, err := readwriter.S3api.GetObject(input)
	if err != nil {
		klog.Errorf("Fetch S3 Object Failed: %v", err)
	}

	data, err := ioutil.ReadAll(output.Body)
	if err != nil {
		klog.Errorf("Read S3 Object Contents: %v", err)
		return nil, err
	}

	return data, err
}

func (readwriter *S3ReadWriter) LoadFiles(dir string) *Files {
	klog.Info("List Files with prefix: ", dir)

	input := &s3.ListObjectsInput{
		Bucket: aws.String(readwriter.Bucket),
		Prefix: aws.String(dir),
	}

	output, err := readwriter.S3api.ListObjects(input)
	if err != nil {
		klog.Fatalf("List Files Failed: %v", err)
		return nil
	}

	files := &Files{}
	for _, f := range output.Contents {
		key := strings.TrimPrefix(*f.Key, dir + "/")
		klog.V(3).Info("file name is: ", key)
		switch {
		case strings.HasSuffix(key, dbSuffix):
			files.Databases = append(files.Databases, key)
		case strings.HasSuffix(key, schemaSuffix):
			files.Schemas = append(files.Schemas, key)
		default:
			if strings.HasSuffix(key, tableSuffix) {
				files.Tables = append(files.Tables, key)
			}
		}
	}
	return &Files{}
}
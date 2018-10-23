package framework

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
)

func StorageUpload(file string,bucket,filename string) (string,error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return "",err
	}

	f,err := os.Open(file)
	if err != nil {
		return "",err
	}
	defer f.Close()
	wc := client.Bucket(bucket).Object(filename).NewWriter(ctx)
	defer wc.Close()
	_,err = io.Copy(wc,f);
	if  err != nil {
		return "",err
	}
	url := fmt.Sprintf("https://storage.googleapis.com/%s/%s",bucket,filename)
	return url,nil
}

func StorageUploadFile(r *http.Request,img,bucket,filename string) (string,error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return "",err
	}
	f, fh, err := r.FormFile(img)
	wc := client.Bucket(bucket).Object(filename).NewWriter(ctx)
	wc.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}
	wc.ContentType = fh.Header.Get("Content-Type")
	wc.CacheControl = "public, max-age=86400"
	defer wc.Close()
	_,err = io.Copy(wc,f);
	if  err != nil {
		return "",err
	}
	url := fmt.Sprintf("https://storage.googleapis.com/%s/%s",bucket,filename)
	return url,nil
}

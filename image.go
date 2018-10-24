package framework

import (
	"cloud.google.com/go/storage"
	"context"
	"github.com/disintegration/imaging"
	"github.com/satori/go.uuid"
	"image"
	"mime/multipart"
	"path"
	"strconv"
	"time"
)

type Image struct {
	img image.Image
	ext string
	file *multipart.FileHeader
	Encrypt bool
	Width int
	Height int
}

type Filename struct {
	Filename string
	Dir string
	Fullpath string
}


func(img *Image) Set(file *multipart.FileHeader) error {
	img.Encrypt = false
	f,err := file.Open()
	if err!= nil {
		return err
	}
	defer f.Close()
	src, err := imaging.Decode(f)
	if err!= nil {
		return err
	}
	img.img = src
	img.file = file
	return nil
}

func(img *Image) ResizeSave(location string)  (Filename,error){
	var f Filename
	img.img = imaging.Resize(img.img,img.Width,img.Height, imaging.Lanczos)
	filename := img.file.Filename

	if img.Encrypt {
		filename = unix() + path.Ext(filename)
	}
	f.Dir = location
	f.Filename = filename
	f.Fullpath = location+filename
	err := imaging.Save(img.img, f.Fullpath)
	if err != nil {
		return Filename{},err
	}
	return f,nil
}

func(img *Image) ResizeUpload(bucket string)  (Filename,error){
	var f Filename
	img.img = imaging.Resize(img.img,img.Width,img.Height, imaging.Lanczos)
	//filename := img.file.Filename
	filename := unix() + ".png"
	f.Dir = "https://storage.googleapis.com/"+bucket+"/"
	f.Filename = filename
	f.Fullpath = f.Dir+f.Filename


	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return Filename{},err
	}

	wc := client.Bucket(bucket).Object(f.Filename).NewWriter(ctx)
	wc.CacheControl = "public, max-age=86400"

	err = imaging.Encode(wc, img.img,imaging.PNG )
	if err != nil {
		return Filename{},err
	}

	if err = wc.Close(); err != nil {
		return Filename{},err
	}

	return f,nil
}

func(img *Image) ResizeMultiUpload(bucket string,size map[string]uint)  (Filename,error){
	var f Filename
	filename := unix() + ".png"
	f.Dir = "https://storage.googleapis.com/"+bucket+"/"
	f.Filename = filename
	f.Fullpath = f.Dir+f.Filename
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return Filename{},err
	}
	for i,v := range size {
		img.img = imaging.Resize(img.img,int(v),0, imaging.Lanczos)
		wc := client.Bucket(bucket).Object(i+f.Filename).NewWriter(ctx)
		wc.CacheControl = "public, max-age=86400"
		err = imaging.Encode(wc, img.img,imaging.PNG )
		if err != nil {
			return Filename{},err
		}
		if err = wc.Close(); err != nil {
			return Filename{},err
		}
	}
	return f,nil
}

func unix() string {
	t := strconv.Itoa(int(time.Now().UnixNano()))
	return t+uuid.Must(uuid.NewV4()).String()
}

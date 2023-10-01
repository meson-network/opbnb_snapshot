package uploader_greenfield

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/bnb-chain/greenfield-go-sdk/client"
	"github.com/bnb-chain/greenfield-go-sdk/types"
	storageTypes "github.com/bnb-chain/greenfield/x/storage/types"
	"github.com/meson-network/opbnb_snapshot/src/utils/custom_reader"
	"github.com/vbauerster/mpb/v8"
)

type UploaderWorker struct {
	client          client.IClient
	bucketName      string
	additional_path string
	fileName        string
	bar             *mpb.Bar
}

func newUploadWorker(cli client.IClient, bucketName string, additional_path string,
	fileName string, bar *mpb.Bar) *UploaderWorker {

	return &UploaderWorker{
		client: cli, bucketName: bucketName, additional_path: additional_path,
		fileName: fileName, bar: bar,
	}
}

func (u *UploaderWorker) uploadFile(localFilePath string) error {

	keyInRemote := u.fileName
	if u.additional_path != "" {
		keyInRemote = u.additional_path + "/" + u.fileName
	}

	fileInfo, err := os.Stat(localFilePath)
	if err != nil {
		return err
	}

	// check remote
	exists, sameSize := validateRemoteChunk(u.client, u.bucketName, keyInRemote, fileInfo.Size())
	if exists && sameSize {
		return nil
	} else if exists && !sameSize {
		err := deleteRemoteFile(u.client, u.bucketName, keyInRemote)
		if err != nil {
			return err
		}
	}

	// upload new one
	uploadFile, err := os.OpenFile(localFilePath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer uploadFile.Close()

	// use custom reader to show upload progress
	reader := &custom_reader.CustomReader{
		Reader:      uploadFile,
		Size:        fileInfo.Size(),
		DownloadBar: nil,
		UploadBar:   u.bar,
	}

	// upload file
	// create and put object
	txnHash, err := u.client.CreateObject(context.Background(), u.bucketName, keyInRemote, uploadFile, types.CreateObjectOptions{})
	if err != nil {
		return err
	}
	uploadFile.Seek(0, 0)

	// Put your object
	err = u.client.PutObject(context.Background(), u.bucketName, keyInRemote, fileInfo.Size(),
		reader, types.PutObjectOptions{TxnHash: txnHash})
	if err != nil {
		// cancel create
		u.client.CancelCreateObject(context.Background(),u.bucketName,keyInRemote,types.CancelCreateOption{})

		// return err
		return err
	}

	_, err = u.client.UpdateObjectVisibility(context.Background(), u.bucketName, keyInRemote, storageTypes.VISIBILITY_TYPE_PUBLIC_READ, types.UpdateObjectOption{})
	if err != nil {

	}

	//wait for SP to seal your object
	err = waitObjectSeal(u.client, u.bucketName, keyInRemote)
	if err != nil {

	}

	return nil
}

func waitObjectSeal(cli client.IClient, bucketName, objectName string) error {
	ctx := context.Background()
	// wait for the object to be sealed
	timeout := time.After(15 * time.Second)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			err := errors.New("object not sealed after 15 seconds")
			return err
		case <-ticker.C:
			objectDetail, err := cli.HeadObject(ctx, bucketName, objectName)
			if err != nil {
				return err
			}
			if objectDetail.ObjectInfo.GetObjectStatus().String() == "OBJECT_STATUS_SEALED" {
				fmt.Printf("put object %s successfully \n", objectName)
				return nil
			}
		}
	}
}

func validateRemoteChunk(cli client.IClient, bucketName string, keyInRemote string, localFileSize int64) (sameSize bool, exist bool) {
	// get fileInfo from bucket
	reader, info, err := cli.GetObject(context.Background(), bucketName, keyInRemote, types.GetObjectOptions{})
	if err != nil {
		return false, false
	}
	defer reader.Close()
	return info.Size == localFileSize, true
}

func deleteRemoteFile(cli client.IClient, bucketName string, keyInRemote string) error {
	// delete object
	delTx, err := cli.DeleteObject(context.Background(), bucketName, keyInRemote, types.DeleteObjectOption{})
	if err != nil {
		return err
	}
	_, err = cli.WaitForTx(context.Background(), delTx)
	if err != nil {
		return err
	}
	return nil
}

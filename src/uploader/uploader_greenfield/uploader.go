package uploader_greenfield

import (
	"errors"
	"fmt"
	"math"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"

	"github.com/bnb-chain/greenfield-go-sdk/client"
	"github.com/meson-network/opbnb_snapshot/src/config"
	"github.com/meson-network/opbnb_snapshot/src/model"
)

const (
	DEFAULT_RETRY_TIMES = 5
	DEFAULT_THREAD      = 5
)

func Upload_greenfield(originDir string, thread int, bucketName string, additional_path string,
	privateKey string,rpcAddr string, chainId string,retryTimes int) error {

	if thread <= 0 {
		thread = DEFAULT_THREAD
	}

	if retryTimes <= 0 {
		retryTimes = DEFAULT_RETRY_TIMES
	}

	// read json from originDir
	configFilePath := filepath.Join(originDir, model.DEFAULT_CONFIG_NAME)
	fileConfig, err := config.LoadFile4Upload(configFilePath)
	if err != nil {
		fmt.Println("[ERROR] ", err.Error())
		return err
	}

	client, err := genGreenfieldClient(privateKey,rpcAddr,chainId)
	if err != nil {
		fmt.Println("[ERROR] gen r2 client err:", err)
		return err
	}

	fmt.Println("[INFO] start upload...")

	if err := upload_file(originDir, thread, retryTimes, fileConfig,
		client, bucketName, additional_path); err != nil {
		return err
	}

	upload_config(originDir,
		client, bucketName, additional_path)

	return nil
}

func upload_file(originDir string, thread int, retryTimes int, fileConfig *model.FileConfig,
	client client.IClient, bucketName string, additional_path string) error {

	fileList := fileConfig.ChunkedFileList
	errorFiles := []*model.ChunkedFileInfo{}
	var errorFilesLock sync.Mutex

	var wg sync.WaitGroup
	progressBar := mpb.New(mpb.WithAutoRefresh())
	counter := int64(0)

	threadChan := make(chan struct{}, thread)
	for _, v := range fileList {
		fileInfo := v

		threadChan <- struct{}{}
		wg.Add(1)
		go func() {
			defer func() {
				<-threadChan
				wg.Done()
			}()

			c := atomic.AddInt64(&counter, 1)
			bar := progressBar.AddBar(
				int64(100),
				mpb.BarRemoveOnComplete(),
				mpb.BarFillerClearOnComplete(),
				mpb.PrependDecorators(
					// simple name decorator
					decor.Name(fmt.Sprintf(" %d / %d %s ", c, len(fileList), fileInfo.FileName)),
					// decor.DSyncWidth bit enables column width synchronization
					decor.Percentage(decor.WCSyncSpace),
				),
				mpb.AppendDecorators(
					decor.OnComplete(
						decor.Name(""), "SUCCESS ",
					),
					decor.OnAbort(
						decor.Elapsed(decor.ET_STYLE_GO), "FAILED ",
					),
				),
			)

			bar.SetPriority(math.MaxInt - len(fileList) + int(c))

			uploadWorker := newUploadWorker(client, bucketName, additional_path, fileInfo.FileName,  bar)

			// try some times if upload failed
			for try := 0; try < retryTimes; try++ {
				bar.SetCurrent(0)

				localFilePath := filepath.Join(originDir, fileInfo.FileName)

				err := uploadWorker.uploadFile(localFilePath)
				if err != nil {
					if try < retryTimes-1 {
						time.Sleep(3 * time.Second)
						continue
					}

					errorFilesLock.Lock()
					errorFiles = append(errorFiles, &fileInfo)
					errorFilesLock.Unlock()

					bar.Abort(false)
					bar.SetPriority(math.MaxInt - int(c) - len(fileList))
				} else {
					if !bar.Completed() {
						bar.SetCurrent(100)
					}
					bar.SetPriority(int(c))
					return
				}
			}
		}()
	}
	// must wait wg first
	wg.Wait()
	progressBar.Wait()

	if len(errorFiles) > 0 {
		fmt.Println("[ERROR] the following files upload failed, please try again:")
		for _, v := range errorFiles {
			fmt.Println(v.FileName)
		}
		return errors.New("upload error")
	}

	return nil
}

func upload_config(originDir string, client client.IClient, bucketName string, additional_path string) {

	fileDir, fileName := originDir, model.DEFAULT_CONFIG_NAME

	progressBar := mpb.New(mpb.WithAutoRefresh())
	bar := progressBar.AddBar(
		int64(100),
		mpb.PrependDecorators(
			// simple name decorator
			decor.Name(fmt.Sprintf(" %s ", fileName)),
			// decor.DSyncWidth bit enables column width synchronization
			decor.Percentage(decor.WCSyncSpace),
		),
		mpb.AppendDecorators(
			decor.OnComplete(
				decor.Name(""), "SUCCESS ",
			),
			decor.OnAbort(
				decor.Elapsed(decor.ET_STYLE_GO), "FAILED ",
			),
		),
	)

	uploadWorker := newUploadWorker(client, bucketName, additional_path, fileName, bar)

	localFilePath := filepath.Join(fileDir, fileName)
	err := uploadWorker.uploadFile(localFilePath)

	progressBar.Wait()
	if err != nil {
		bar.Abort(false)
		fmt.Println("[ERROR] upload json file error")
	} else {
		if !bar.Completed() {
			bar.SetCurrent(100)
		}
		fmt.Println("[INFO] upload job finish")
	}
}
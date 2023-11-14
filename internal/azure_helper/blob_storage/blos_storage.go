package blob_storage

import (
	"context"
	"fmt"
	"github.com/dyammarcano/eventReceiverQuel/internal/azure_helper/event_hub"
	"github.com/dyammarcano/eventReceiverQuel/internal/logger"

	"go.uber.org/zap"
	"os"
	"strings"
	"time"
)

//https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/storage/azblob
//https://www.thorsten-hans.com/azure-blob-storage-using-azure-sdk-for-golang
//https://github.com/ThorstenHans/azb
//https://github.com/Azure-Samples/storage-blobs-go-quickstart

type (
	BlobItem struct {
		ETag          string             `xml:"Etag"`
		Name          string             `xml:"Name"`
		Metadata      map[string]*string `xml:"Metadata"`
		LastModified  time.Time          `xml:"Last-Modified"`
		AccessTier    string             `xml:"AccessTier"`
		ContentLength int64              `xml:"Content-Length"`
		ContentMD5    []byte             `xml:"Content-MD5"`
		ContentType   string             `xml:"Content-MimeType"`
		CreationTime  time.Time          `xml:"Creation-Time"`
		LeaseState    string             `xml:"LeaseState"`
		LeaseStatus   string             `xml:"LeaseStatus"`
		RootDir       bool               `xml:"RootDir"`
	}

	StorageClient struct {
		*azblob.Client
		context.Context
	}
)

func NewStorageClient(ctx context.Context) *StorageClient {
	config, ok := ctx.Value("azure").(event_hub.AzureCredentials)
	if !ok {
		logger.Error("error getting azureConfig from context")
		os.Exit(1)
	}

	if config.Storage.AccountName == "" || config.Storage.AccountKey == "" {
		panic("azure.storage.account_name or azure.storage.account_key not found")
	}

	cred, err := azblob.NewSharedKeyCredential(config.Storage.AccountName, config.Storage.AccountKey)
	if err != nil {
		return nil
	}

	client, err := azblob.NewClientWithSharedKeyCredential(fmt.Sprintf("https://%s.blob.core.windows.net/", config.Storage.AccountName), cred, nil)
	if err != nil {
		return nil
	}

	return &StorageClient{
		Client:  client,
		Context: ctx,
	}
}

func (s *StorageClient) UploadBlob(filePath string, blobName string, containerName string) (response BlobItem, err error) {
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0)

	if err != nil {
		return BlobItem{}, err
	}

	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			logger.Error("Error", zap.Any("msg", err))
		}
	}(file)

	data, err := s.UploadFile(context.Background(), containerName, blobName, file, nil)

	if err != nil {
		return BlobItem{}, err
	}

	return BlobItem{
		ETag:         string(*data.ETag),
		Name:         blobName,
		LastModified: *data.LastModified,
		ContentMD5:   data.ContentMD5,
	}, nil
}

func (s *StorageClient) DownloadBlob(filePath string, blobName string, containerName string) (response int64, err error) {
	file, err := os.Create(filePath)

	if err != nil {
		return 0, err
	}

	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			logger.Error("Error", zap.Any("msg", err))
		}
	}(file)

	response, err = s.DownloadFile(context.Background(), containerName, blobName, file, nil)

	if err != nil {
		return 0, err
	}

	return 0, nil
}

func (s *StorageClient) BlobDelete(blobName string, containerName string) (response azblob.DeleteBlobResponse, err error) {
	response, err = s.DeleteBlob(context.Background(), containerName, blobName, nil)
	return
}

func (s *StorageClient) ContainerDelete(containerName string) (response azblob.DeleteContainerResponse, err error) {
	response, err = s.DeleteContainer(context.Background(), containerName, nil)
	return
}

func (s *StorageClient) ContainerCreate(containerName string) (response azblob.CreateContainerResponse, err error) {
	response, err = s.CreateContainer(context.Background(), containerName, nil)
	return
}

func (s *StorageClient) ListBlobs(containerName string) (blobItems []*BlobItem, err error) {
	pager := s.NewListBlobsFlatPager(containerName, &azblob.ListBlobsFlatOptions{
		Include: azblob.ListBlobsInclude{Snapshots: true, Versions: true},
	})

	for pager.More() {
		page, err := pager.NextPage(context.Background())

		if err != nil {
			return nil, err
		}

		for _, blob := range page.Segment.BlobItems {
			if blob.Metadata == nil {
				blob.Metadata = make(map[string]*string)
			}

			rootDir := !strings.Contains(*blob.Name, "/")

			blobItems = append(blobItems, &BlobItem{
				Name:          *blob.Name,
				LastModified:  *blob.Properties.LastModified,
				ContentLength: *blob.Properties.ContentLength,
				ContentType:   *blob.Properties.ContentType,
				CreationTime:  *blob.Properties.CreationTime,
				ETag:          string(*blob.Properties.ETag),
				AccessTier:    string(*blob.Properties.AccessTier),
				LeaseState:    string(*blob.Properties.LeaseState),
				LeaseStatus:   string(*blob.Properties.LeaseStatus),
				RootDir:       rootDir,
				Metadata:      blob.Metadata,
				ContentMD5:    blob.Properties.ContentMD5,
			})
			// for debugging
			//s, _ := json.Marshal(blobItems)
			//fmt.Println(string(s))
		}
	}
	return
}

//func OpenBlob(containerURL azblob.ContainerURL, blobName string) (*azblob.Blob, error) {
//	ctx := context.Background()
//	blobURL := containerURL.NewBlobURL(blobName)
//	properties, err := blobURL.GetProperties(ctx, azblob.BlobAccessConditions{}, azblob.ClientProvidedKeyOptions{})
//	if err != nil {
//		return nil, fmt.Errorf("failed to get blob properties: %v", err)
//	}
//	blobSize := properties.ContentLength()
//	blobReader, err := blobURL.Download(ctx, 0, azblob.CountToEnd, azblob.BlobAccessConditions{}, false, azblob.ClientProvidedKeyOptions{})
//	if err != nil {
//		return nil, fmt.Errorf("failed to download blob: %v", err)
//	}
//	defer blobReader.Close()
//	blob := azblob.NewBlob(blobURL, blobSize)
//	return &blob, nil
//}

//func ListContainers(client *azblob.Client) (response azblob.ListContainersResponse, err error) {
//	response, err = client.ListContainers(context.Background(), nil)
//	return
//}
//
//func BlobExists(client *azblob.Client, blobName string, containerName string) (response azblob.BlobExistsResponse, err error) {
//	response, err = client.BlobExists(context.Background(), containerName, blobName, nil)
//	return
//}

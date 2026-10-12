package googledrive

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

const ParentFolderID = "1zDqP4XB1LOJT5GJblhNSAgTk35QOTgjR"

type Uploader struct {
	service *drive.Service
}

func NewUploader(ctx context.Context, client *http.Client) (*Uploader, error) {
	service, err := drive.NewService(
		ctx,
		option.WithHTTPClient(client),
	)
	if err != nil {
		return nil, err
	}

	return &Uploader{service: service}, nil
}

func (u *Uploader) findFolder(ctx context.Context, parentID, name string) (string, error) {
	query := fmt.Sprintf(
		"name = '%s' and mimeType = 'application/vnd.google-apps.folder' and '%s' in parents and trashed = false",
		escapeQuery(name), escapeQuery(parentID),
	)

	result, err := u.service.Files.List().
		Q(query).
		Spaces("drive").
		PageSize(100).
		Fields("files(id,name)").
		IncludeItemsFromAllDrives(true).
		SupportsAllDrives(true).
		Do()
	if err != nil {
		return "", err
	}

	if len(result.Files) > 0 {
		return result.Files[0].Id, nil
	}

	folder, err := u.service.Files.Create(&drive.File{
		Name:     name,
		MimeType: "application/vnd.google-apps.folder",
		Parents:  []string{parentID},
	}).
		SupportsAllDrives(true).
		Fields("id").
		Do()
	if err != nil {
		return "", fmt.Errorf("create folder %q: %w", name, err)
	}

	return folder.Id, nil
}

func escapeQuery(s string) string {
	// Escape backslashes and apostrophes for Drive query strings.
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		if s[i] == '\\' || s[i] == '\'' {
			out = append(out, '\\')
		}
		out = append(out, s[i])
	}
	return string(out)
}

func (u *Uploader) Upload(ctx context.Context, path string) error {
	filename := filepath.Base(path)

	destination, err := ClassifyLog(filename)
	if err != nil {
		return err
	}

	yearID, err := u.findFolder(ctx, ParentFolderID, destination.Year)
	if err != nil {
		return err
	}

	destinationID, err := u.findFolder(ctx, yearID, destination.Event)
	if err != nil {
		return err
	}

	// Skip duplicates rather than silently creating another copy.
	query := fmt.Sprintf(
		"name = '%s' and '%s' in parents and trashed = false",
		escapeQuery(filename), escapeQuery(destinationID),
	)

	existing, err := u.service.Files.List().
		Q(query).
		Spaces("drive").
		PageSize(1).
		Fields("files(id)").
		IncludeItemsFromAllDrives(true).
		SupportsAllDrives(true).
		Do()
	if err != nil {
		return err
	}

	if len(existing.Files) > 0 {
		fmt.Printf("Skipping existing file: %s\n", filename)
		return nil
	}

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	upload := &drive.File{
		Name:    filename,
		Parents: []string{destinationID},
	}

	created, err := u.service.Files.Create(upload).
		Media(file).
		SupportsAllDrives(true).
		Fields("id,name,size").
		Do()
	if err != nil {
		return fmt.Errorf("upload %s: %w", filename, err)
	}

	fmt.Printf("Uploaded %s (%d bytes; Drive ID %s)\n",
		created.Name, info.Size(), created.Id)

	return nil
}

func (u *Uploader) UploadMany(ctx context.Context, paths []string) error {
	var failures []error

	for _, path := range paths {
		if err := u.Upload(ctx, path); err != nil {
			fmt.Fprintf(os.Stderr, "Failed %s: %v\n", path, err)
			failures = append(failures, fmt.Errorf("%s: %w", path, err))
		}
	}

	if len(failures) > 0 {
		return fmt.Errorf("%d upload(s) failed", len(failures))
	}

	return nil
}

// Keep io imported if you later add a progress-reporting reader.
var _ io.Reader

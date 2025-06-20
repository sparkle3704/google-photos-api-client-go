package uploader_test

import (
	"context"
	"github.com/sparkle3704/google-photos-api-client-go/mocks"
	"github.com/sparkle3704/google-photos-api-client-go/uploader"
	"net/http"
	"testing"
)

func TestNewSimpleUploader(t *testing.T) {
	got, err := uploader.NewSimpleUploader(http.DefaultClient)
	if err != nil {
		t.Fatalf("error was not expected at this point: %s", err)
	}
	want := "https://photoslibrary.googleapis.com/v1/uploads"

	if want != got.BaseURL {
		t.Errorf("want: %s, got: %s", want, got)
	}
}

func TestSimpleUploader_UploadFile(t *testing.T) {
	testCases := []struct {
		name        string
		path        string
		uploadToken string
		errExpected bool
	}{
		{name: "Upload should be successful", path: "testdata/upload-success", uploadToken: mocks.UploadToken, errExpected: false},
		{name: "Upload existing file with errors should be a failure", path: "testdata/upload-should/fail", uploadToken: "", errExpected: true},
		{name: "Upload a non-existing file should be a failure", path: "non-existent", uploadToken: "", errExpected: true},
	}
	srv := mocks.NewMockedGooglePhotosService()
	defer srv.Close()

	u, err := uploader.NewSimpleUploader(http.DefaultClient)
	if err != nil {
		t.Fatalf("error was not expected at this point: %s", err)
	}
	u.BaseURL = srv.URL() + "/v1/uploads"

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := u.UploadFile(context.Background(), tc.path)
			if tc.errExpected && err == nil {
				t.Fatalf("error was expected, but not produced")
			}
			if !tc.errExpected && err != nil {
				t.Fatalf("error was not expected, err: %s", err)
			}
			if err == nil && tc.uploadToken != got {
				t.Errorf("want: %s, got: %s", tc.uploadToken, got)
			}
		})
	}
}

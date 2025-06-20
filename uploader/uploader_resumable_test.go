package uploader_test

import (
	"context"
	"github.com/sparkle3704/google-photos-api-client-go/mocks"
	"github.com/sparkle3704/google-photos-api-client-go/uploader"
	"net/http"
	"testing"
)

func TestNewResumableUploader(t *testing.T) {
	u, err := uploader.NewResumableUploader(http.DefaultClient)
	if err != nil {
		t.Fatalf("error was not expected at this point: %s", err)
	}
	want := "https://photoslibrary.googleapis.com/v1/uploads"

	if want != u.BaseURL {
		t.Errorf("want: %s, got: %s", want, u.BaseURL)
	}
}

func TestResumableUploader_UploadFile(t *testing.T) {
	testCases := []struct {
		name           string
		path           string
		alreadyStarted bool
		errExpected    bool
	}{
		{"Should be successful when file is uploaded", "testdata/upload-success", false, false},
		{"Should be successful when file is resuming upload ", "testdata/upload-resume-success", true, false},
		{"Should fail when file is not uploaded", "testdata/upload-should-fail", false, true},
		{"Should fail if file doesn't exist", "non-existent", false, true},
	}
	srv := mocks.NewMockedGooglePhotosService()
	defer srv.Close()

	store := NewMockStore()

	u, err := uploader.NewResumableUploader(http.DefaultClient)
	u.BaseURL = srv.URL() + "/v1/uploads"
	u.Store = store

	if err != nil {
		t.Fatalf("error was not expected at this point, err: %s", err)
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := u.UploadFile(context.Background(), tc.path)
			if tc.errExpected && err == nil {
				t.Fatalf("error was expected, but not produced")
			}
			if !tc.errExpected && err != nil {
				t.Fatalf("error was not expected, err: %s", err)
			}
			want := mocks.UploadToken
			if err == nil && want != got {
				t.Errorf("want: %s, got: %s", want, got)
			}
		})
	}
}

func TestResumableUploader_IsResumeEnabled(t *testing.T) {
	t.Run("Should return False by default", func(t *testing.T) {
		u, err := uploader.NewResumableUploader(http.DefaultClient)
		if err != nil {
			t.Fatalf("error not expected at this point: %s", err)
		}

		if u.IsResumeEnabled() {
			t.Errorf("want: false, got: true")
		}
	})

	t.Run("Should return True when an Store is configured", func(t *testing.T) {
		store := NewMockStore()

		u, err := uploader.NewResumableUploader(http.DefaultClient)
		if err != nil {
			t.Fatalf("error not expected at this point: %s", err)
		}
		u.Store = store

		if !u.IsResumeEnabled() {
			t.Errorf("want: true, got: false")
		}
	})
}

type MockStore struct {
	m map[string]string
}

func NewMockStore() uploader.Store {
	return &MockStore{
		make(map[string]string),
	}
}

func (s *MockStore) Get(fingerprint string) (string, bool) {
	url, ok := s.m[fingerprint]
	return url, ok
}

func (s *MockStore) Set(fingerprint, url string) {
	s.m[fingerprint] = url
}

func (s *MockStore) Delete(fingerprint string) {
	delete(s.m, fingerprint)
}

func (s *MockStore) Close() {
	for k := range s.m {
		delete(s.m, k)
	}
}

package env

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"os/user"
	"testing"

	"github.com/spf13/afero"
)

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	// fmt.Println(os.Args[3:])
	os.Exit(0)
}

func mockExecCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

func TestDescribe(t *testing.T) {
	cases := []struct {
		Body                io.Reader
		UnreachableEndpoint bool
		ExpectedErr         error
	}{
		{bytes.NewBuffer([]byte(`{"origin":"127.0.0.1"}`)), false, nil},
		{bytes.NewBuffer([]byte(`{"test:"127.0.0.1"}`)), false, ErrJSONPayloadInvalidFormat},
		{nil, false, ErrJSONPayloadInvalidFormat},
		{nil, true, ErrHTTPError},
	}
	for _, tc := range cases {
		httpGet = http.Get

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, tc.Body)
		}))
		defer ts.Close()

		if tc.UnreachableEndpoint {
			httpGet = func(url string) (resp *http.Response, err error) {
				return nil, ErrHTTPError
			}
		}

		IPCheckerURL = ts.URL
		defer func() { IPCheckerURL = "http://httpbin.org/ip" }()

		if err := Describe(); err != tc.ExpectedErr {
			t.Errorf("Describe func returned wrong err: got %#v want %#v",
				err, tc.ExpectedErr)
		}
	}
}

func TestEnsureDotfilesDir(t *testing.T) {
	AppFs = afero.NewMemMapFs()
	defer func() { AppFs = afero.NewOsFs() }()

	cases := []struct {
		DotfilesDirPath string
		PermissionOk    bool
		ExpectedErr     error
	}{
		{"/Users/thylong/.dotfiles", true, nil},
		{"/Users/thylong/.dotfiles", true, nil},
		{"\\Users\\thylong\\.dotfiles", true, nil},
		{"/Users/thylong/.dotfiles", false, ErrOperationNotPermitted},
	}
	for _, tc := range cases {
		if !tc.PermissionOk {
			AppFs = afero.NewReadOnlyFs(AppFs)
		}
		if err := EnsureDotfilesDir(tc.DotfilesDirPath); err != tc.ExpectedErr {
			t.Errorf("Describe func returned wrong err: got %#v want %#v",
				err, tc.ExpectedErr)
		}
		AppFs = afero.NewMemMapFs()
	}
}

func TestImportIntoDotfilesDir(t *testing.T) {
	AppFs = afero.NewMemMapFs()
	defer func() { AppFs = afero.NewOsFs() }()

	cases := []struct {
		DotfilesToSave  []string
		DotfilesDirPath string
		FileExists      bool
		PermissionOk    bool
		ExpectedErr     error
	}{
		{[]string{"testong.yml"}, "/Users/thylong/.dotfiles", false, true, ErrCannotMoveDotfile},
		{[]string{"testong.yml"}, "/Users/thylong/.dotfiles", true, false, ErrCannotMoveDotfile},
		{[]string{"testong.yml"}, "/Users/thylong/.dotfiles", true, true, ErrCannotSymlink},
	}
	for _, tc := range cases {
		if tc.FileExists {
			usr, _ := user.Current()
			AppFs.Mkdir(tc.DotfilesDirPath, 0766)
			afero.WriteFile(AppFs, fmt.Sprintf("%s/%s", usr.HomeDir, tc.DotfilesToSave[0]), []byte("test"), 0644)
		}
		if !tc.PermissionOk {
			AppFs = afero.NewReadOnlyFs(AppFs)
		}
		if err := ImportIntoDotfilesDir(tc.DotfilesToSave, tc.DotfilesDirPath); err != tc.ExpectedErr {
			t.Errorf("ImportIntoDotfilesDir func returned wrong err: got %#v want %#v",
				err, tc.ExpectedErr)
		}
		AppFs = afero.NewMemMapFs()
	}
}

// func TestEnsureDotfilesRepository(t *testing.T) {
// 	AppFs = afero.NewMemMapFs()
// 	execCommand = mockExecCommand
// 	defer func() {
// 		AppFs = afero.NewOsFs()
// 		execCommand = exec.Command
// 	}()
//
// 	cases := []struct {
// 		DotfilesRepository string
// 		DotfilesDirPath    string
// 		ExpectedErr        error
// 	}{
// 		{"thylong/dotfiles", "/Users/thylong/.dotfiles", nil},
// 	}
// 	for _, tc := range cases {
// 		if err := EnsureDotfilesRepository(tc.DotfilesRepository, tc.DotfilesDirPath); err != tc.ExpectedErr {
// 			t.Errorf("EnsureDotfilesRepository func returned wrong err: got %#v want %#v",
// 				err, tc.ExpectedErr)
// 		}
// 		AppFs = afero.NewMemMapFs()
// 	}
// }
//
// func TestPersistDotfiles(t *testing.T) {
// 	AppFs = afero.NewMemMapFs()
// 	execCommand = mockExecCommand
// 	defer func() {
// 		AppFs = afero.NewOsFs()
// 		execCommand = exec.Command
// 	}()
//
// 	cases := []struct {
// 		Message         string
// 		DotfilesDirPath string
// 		FileExists      bool
// 		PermissionOk    bool
// 		ExpectedErr     error
// 	}{
// 		{"coucou", "/Users/thylong/.dotfiles", true, true, nil},
// 	}
// 	for _, tc := range cases {
// 		if tc.FileExists {
// 			AppFs.Mkdir(tc.DotfilesDirPath, 0766)
// 		}
// 		if !tc.PermissionOk {
// 			AppFs = afero.NewReadOnlyFs(AppFs)
// 		}
// 		if err := PersistDotfiles(tc.Message, tc.DotfilesDirPath); err != tc.ExpectedErr {
// 			fmt.Println(AppFs.Name())
// 			t.Errorf("PersistDotfiles func returned wrong err: got %#v want %#v",
// 				err, tc.ExpectedErr)
// 		}
// 		AppFs = afero.NewMemMapFs()
// 	}
// }

package env

import (
	"testing"

	"github.com/spf13/afero"
)

func TestCopyFile(t *testing.T) {
	AppFs = afero.NewMemMapFs()
	defer func() { AppFs = afero.NewOsFs() }()

	cases := []struct {
		SrcPath      string
		SrcExists    bool
		IsSrcSymlink bool
		DstPath      string
		DstExists    bool
		ExpectedErr  error
	}{
		{"test.yml", true, false, "test.yml", false, nil},
		{"test.yml", true, false, "test.yml", true, nil},
		{"test.yml", true, false, "test2.yml", true, nil},
		{"test.yml", false, false, "test2.yml", false, ErrCannotStatFile},
	}
	for _, tc := range cases {
		if tc.SrcExists {
			afero.WriteFile(AppFs, tc.SrcPath, []byte("test"), 0644)
		}
		if tc.DstExists {
			afero.WriteFile(AppFs, tc.DstPath, []byte("test"), 0644)
		}

		err := CopyFile(tc.SrcPath, tc.DstPath)
		if err != tc.ExpectedErr {
			t.Errorf("CopyFile func returned wrong error: got %#v want %#v",
				err, tc.ExpectedErr)
		}

		AppFs.Remove(tc.SrcPath)
		AppFs.Remove(tc.DstPath)
	}
}

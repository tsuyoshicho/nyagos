package dos

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
	"unicode"
	"unsafe"
)

var msvcrt = syscall.NewLazyDLL("msvcrt")
var _chdrive = msvcrt.NewProc("_chdrive")
var _wchdir = msvcrt.NewProc("_wchdir")

func chdrive_(n rune) uintptr {
	rc, _, _ := _chdrive.Call(uintptr(n & 0x1F))
	return rc
}

func getFirst(s string) (rune, error) {
	reader := strings.NewReader(s)
	drive, _, err := reader.ReadRune()
	if err != nil {
		return 0, err
	}
	return unicode.ToUpper(drive), nil
}

// Change drive without changing the working directory there.
func Chdrive(drive string) error {
	driveLetter, driveErr := getFirst(drive)
	if driveErr != nil {
		return driveErr
	}
	chdrive_(driveLetter)
	return nil
}

var rxPath = regexp.MustCompile("^([a-zA-Z]):(.*)$")

// Change the current working directory
// without changeing the working directory
// in the last drive.
func Chdir(folder_ string) error {
	folder := folder_
	if m := rxPath.FindStringSubmatch(folder_); m != nil {
		status := chdrive_(rune(m[1][0]))
		if status != 0 {
			return fmt.Errorf("%s: no such directory", folder_)
		}
		folder = m[2]
		if len(folder) <= 0 {
			return nil
		}
	}
	utf16, err := syscall.UTF16PtrFromString(folder)
	if err == nil {
		status, _, _ := _wchdir.Call(uintptr(unsafe.Pointer(utf16)))
		if status != 0 {
			err = fmt.Errorf("%s: no such directory", folder_)
		}
	}
	return err
}

var rxRoot = regexp.MustCompile("^([a-zA-Z]:)?[/\\\\]")
var rxDrive = regexp.MustCompile("^[a-zA-Z]:")

func joinPath2(a, b string) string {
	if len(a) <= 0 || rxRoot.MatchString(b) || rxDrive.MatchString(b) {
		return b
	}
	switch a[len(a)-1] {
	case '\\', '/', ':':
		return a + b
	default:
		return a + "\\" + b
	}
}

// Equals filepath.Join but this works right when path has drive-letter.
func Join(paths ...string) string {
	result := paths[len(paths)-1]
	for i := len(paths) - 2; i >= 0; i-- {
		result = joinPath2(paths[i], result)
	}
	return result
}

// Expand filenames matching with wildcard-pattern.
func Glob(pattern string) ([]string, error) {
	name := filepath.Base(pattern)
	if strings.IndexAny(name, "*?") < 0 {
		return nil, nil
	}
	findf, err := FindFirst(pattern)
	if err != nil {
		return nil, err
	}
	dirname := filepath.Dir(pattern)
	match := make([]string, 0, 100)
	for {
		match = append(match, filepath.Join(dirname, findf.Name()))

		err := findf.FindNext()
		if err != nil {
			return match, nil
		}
	}
}

package basal

import (
	"os"
	"path"
	"path/filepath"
	"runtime"
)

//打开文件 自动创建目录
func OpenFile(folderPath string, fileName string, flag int, perm os.FileMode) (file *os.File, err error) {
	path := filepath.Join(folderPath, fileName)
	file, err = os.OpenFile(path, flag, perm)
	if err != nil {
		var has bool
		has, err = IsExistFolder(folderPath)
		if err != nil {
			return
		}
		if has == false {
			err = os.MkdirAll(folderPath, os.ModeDir)
			if err != nil {
				return
			}
			file, err = os.OpenFile(path, flag, perm)
		}
	}
	return
}

//打开文件 自动创建目录
func OpenFileB(filePath string, flag int, perm os.FileMode) (file *os.File, err error) {
	folderPath, fileName := filepath.Split(filePath)
	return OpenFile(folderPath, fileName, flag, perm)
}

//文件或文件夹是否存在
func IsExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// 文件夹是否存在
func IsExistFolder(path string) (bool, error) {
	handle, err := os.Stat(path)
	if err == nil {
		if handle.IsDir() {
			return true, nil
		} else {
			return false, NewError("not is folder: %s", path)
		}
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

type HandleFile struct {
	folderPath string
	fileName   string
	flag       int
	perm       os.FileMode
	handle     *os.File
}

func (my *HandleFile) PathName() string {
	return filepath.Join(my.folderPath, my.fileName)
}

func (my *HandleFile) SetPathName(folderPath, fileName string) bool {
	if my.folderPath == folderPath && my.fileName == fileName {
		return false
	}
	my.Close()
	my.folderPath = folderPath
	my.fileName = fileName
	return true
}

func (my *HandleFile) WriteString(s string) (n int, err error) {
	n, err = my.Write([]byte(s))
	return
}

func (my *HandleFile) Write(b []byte) (n int, err error) {
	if my.handle == nil {
		my.handle, err = OpenFile(my.folderPath, my.fileName, my.flag, my.perm)
		if err != nil {
			return
		}
	} else {
		var hasLog bool
		hasLog, err = IsExist(filepath.Join(my.folderPath, my.fileName))
		if !hasLog {
			my.handle.Close()
			my.handle, err = OpenFile(my.folderPath, my.fileName, my.flag, my.perm)
			if err != nil {
				return
			}
		}
	}
	n, err = my.handle.Write(b)
	return
}

func (my *HandleFile) Close() {
	if my.handle == nil {
		return
	}
	my.handle.Close()
	my.handle = nil
}

const HANDLE_FILE_FLAG_WRITER = os.O_WRONLY | os.O_APPEND | os.O_CREATE
const HANDLE_FILE_PERM_ALL = 0777

func NewHandleFile(flag int, perm os.FileMode) *HandleFile {
	return &HandleFile{flag: flag, perm: perm}
}

func OpenHandleFile(folderPath string, fileName string, flag int, perm os.FileMode) (*HandleFile, error) {
	var err error
	if folderPath == "" {
		err = NewError("OpenHandleFile Error: folderPath is nil")
		return nil, err
	}
	if fileName == "" {
		err = NewError("OpenHandleFile Error: fileName is nil")
		return nil, err
	}

	hf := &HandleFile{folderPath: folderPath, fileName: fileName, flag: flag, perm: perm}

	hf.handle, err = OpenFile(folderPath, fileName, flag, perm)
	if err != nil {
		return nil, err
	}
	return hf, nil
}

var programDir string
var execDir string

//程序所在路径
func ProgramDir() string {
	if programDir == "" {
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			panic(err)
		}
		programDir = dir
		return dir
	} else {
		return programDir
	}
}

//程序执行路径
func ExecDir() string {
	if execDir == "" {
		dir, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		execDir = dir
		return dir
	} else {
		return execDir
	}
}

func PathBase(p string) string {
	if runtime.GOOS == "windows" {
		return path.Base(filepath.ToSlash(p))
	} else {
		return path.Base(p)
	}
}

func init() {
	ProgramDir()
	ExecDir()
}

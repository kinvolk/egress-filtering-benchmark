// Code generated for package bpf by go-bindata DO NOT EDIT. (@generated)
// sources:
// datapath/bpf.o
package bpf

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

// Name return file name
func (fi bindataFileInfo) Name() string {
	return fi.name
}

// Size return file size
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}

// Mode return file mode
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}

// Mode return file modify time
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}

// IsDir return file whether a directory
func (fi bindataFileInfo) IsDir() bool {
	return fi.mode&os.ModeDir != 0
}

// Sys return file is sys mode
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _datapathBpfO = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x94\x92\xb1\x4a\x33\x41\x10\xc7\xff\xb7\x97\xef\x4b\x24\x46\x62\x02\x12\x05\x45\x8c\x45\x10\x3c\x82\x92\x52\x08\x81\x68\x73\x62\x0a\xc5\xf2\x88\xc7\x19\x02\x17\x89\x77\x57\xc4\xca\x4a\xdf\xc0\x5e\xac\x7c\x83\x94\x69\x7d\x04\x4b\x4b\x1f\x21\x85\x70\x72\xbb\xb3\xba\xec\x25\x04\x07\xb2\xb3\xf3\xcb\x0c\x33\xf3\xbf\xbd\x6f\xdb\xc7\xcc\x30\x20\xcd\xa0\x9f\x6e\x75\xf6\x7b\x6f\xd2\x99\x83\x81\x31\xb1\x6e\xa9\x23\xfc\xaa\xcd\xfd\xa4\x2c\x78\xd6\x04\x56\x00\xec\x57\xd7\x52\x7c\x87\x73\x01\x6e\x4b\x05\xee\x7b\x0c\xc8\x25\x71\x79\x99\xc7\x67\x07\x22\xff\x92\x01\x71\x0c\x5c\xb0\x02\xff\x7f\xcc\x80\x6d\x00\xee\xde\x34\x16\x7d\xb7\x78\x9e\xbb\xf1\xc5\xe3\xc9\x0b\xf5\x61\xc0\x34\x8e\xe3\x8a\xb6\xd4\x03\xed\x39\x21\x2e\xf7\xd8\x24\x2d\x92\x38\x03\xe0\x89\x78\x1e\x62\xae\x84\x3d\x3f\x1a\x3f\x1a\x29\xb2\xe0\xa4\x63\xcf\x50\x4e\xd8\x29\x3f\x4d\x7c\x68\xfc\x10\x40\x11\xff\x53\xf9\x35\xce\xff\xa5\xf8\x3a\xe7\x66\x8a\xbf\x92\x37\x78\x17\x00\x56\xe4\x8d\x22\x58\x81\xe7\x5f\xf7\xfd\xc8\x0b\x1c\xaf\x17\x78\x61\x88\xc8\x95\xb7\x41\x77\x18\xc2\x1f\x0e\x1c\x91\x00\xc7\xf1\xfb\xae\x77\x13\x7a\xb0\xc2\x28\x88\xba\x57\xb0\xc2\xbb\x41\xe2\xed\x56\xab\xee\x34\xe6\xee\xf6\x17\x3b\x92\xf3\x69\xf6\x4e\x82\x9e\x6b\x5c\x7f\x8b\xf2\x7d\xea\x8a\x35\xe7\xf4\xcb\x68\x71\x7e\x41\xfd\xa7\xc6\x73\x5a\x9c\x05\xb0\x34\xa3\xcf\x1b\x0d\x5a\x54\xf2\x4c\xa5\x5e\xf2\xaa\xfa\x8d\x14\xab\x51\xfd\xee\x82\xf9\x1b\x73\xea\x6d\x63\x76\xbe\xae\x5f\x5b\x7b\xb3\xd2\x3a\x94\x38\x52\xea\x98\x32\x7f\x85\xfc\x77\x00\x00\x00\xff\xff\x77\xce\x3b\x80\x30\x04\x00\x00")

func datapathBpfOBytes() ([]byte, error) {
	return bindataRead(
		_datapathBpfO,
		"datapath/bpf.o",
	)
}

func datapathBpfO() (*asset, error) {
	bytes, err := datapathBpfOBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "datapath/bpf.o", size: 1072, mode: os.FileMode(436), modTime: time.Unix(1, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"datapath/bpf.o": datapathBpfO,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"datapath": &bintree{nil, map[string]*bintree{
		"bpf.o": &bintree{datapathBpfO, map[string]*bintree{}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}

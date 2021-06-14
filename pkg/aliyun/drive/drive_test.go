package drive

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var fs Fs

func setup(t *testing.T) context.Context {
	token := ""
	cb, err := ioutil.ReadFile("../../../.config")
	if err == nil {
		token = string(cb)
	}
	config := &Config{
		RefreshToken: token,
	}

	ctx := context.Background()
	fs, err = NewFs(ctx, config)
	require.NoError(t, err)
	return ctx
}

func TestList(t *testing.T) {
	ctx := setup(t)
	names, err := fs.List(ctx, "/media")
	require.NoError(t, err)
	println(fmt.Sprintf("size: %v, %v", len(names), names))
}

func TestCreateFolder(t *testing.T) {
	ctx := setup(t)
	_, err := fs.CreateFolder(ctx, "/")
	require.NoError(t, err)
	_, err = fs.CreateFolder(ctx, "/test3/test4")
	require.NoError(t, err)
}

func TestRename(t *testing.T) {
	ctx := setup(t)
	node, err := fs.Get(ctx, "/test2222", FolderKind)
	require.NoError(t, err)
	err = fs.Rename(ctx, node, "test4")
	require.NoError(t, err)
}

func TestMove(t *testing.T) {
	ctx := setup(t)
	node, err := fs.Get(ctx, "/test3/test4", FolderKind)
	require.NoError(t, err)
	newNode, err := fs.Get(ctx, "/", FolderKind)
	require.NoError(t, err)
	err = fs.Move(ctx, node, newNode)
	require.NoError(t, err)
}

func TestRemove(t *testing.T) {
	ctx := setup(t)
	node, err := fs.Get(ctx, "/test4", FolderKind)
	require.NoError(t, err)
	err = fs.Remove(ctx, node)
	require.NoError(t, err)
}

func TestOpen(t *testing.T) {
	ctx := setup(t)
	node, err := fs.Get(ctx, "/media/1.flac", FileKind)
	require.NoError(t, err)
	fd, err := fs.Open(ctx, node, map[string]string{})
	require.NoError(t, err)
	data, err := ioutil.ReadAll(fd)
	require.NoError(t, err)
	fo, err := os.Create("1.flac")
	require.NoError(t, err)
	fo.Write(data)
	require.NoError(t, fd.Close())
	require.NoError(t, fo.Close())
}

func TestCreateFile(t *testing.T) {
	ctx := setup(t)
	fd, err := os.Open("1.mp3")
	require.NoError(t, err)
	info, err := fd.Stat()
	require.NoError(t, err)
	_, err = fs.CreateFile(ctx, "/media/1.mp3", info.Size(), fd, true)
	require.NoError(t, err)
	defer fd.Close()
}

func TestCopy(t *testing.T) {
	ctx := setup(t)
	node, err := fs.Get(ctx, "/media/1.mp3", FileKind)
	require.NoError(t, err)
	parent, err := fs.Get(ctx, "/", FolderKind)
	require.NoError(t, err)
	err = fs.Copy(ctx, node, parent)
	require.NoError(t, err)
}

func TestSha1(t *testing.T) {
	fd, err := os.Open("1.mp3")
	require.NoError(t, err)
	rd, s, err := CalcSha1(fd)
	assert.Equal(t, "462FD5A7D4B12EE8A88CF0881D811BD224DB79FE", s)
	buf := make([]byte, 4)
	_, _ = rd.Read(buf)
	assert.Equal(t, []byte{0x49, 0x44, 0x33, 0x03}, buf)
}

func TestCalcProof(t *testing.T) {
	fd, err := os.Open("1.mp3")
	accessToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9..."
	fileSize := int64(4087117)
	require.NoError(t, err)
	rd, proofCode, err := calcProof(accessToken, fileSize, fd)
	assert.Equal(t, "dj66UE3TEFM=", proofCode)
	buf2 := make([]byte, 4)
	_, _ = rd.Read(buf2)
	assert.Equal(t, []byte{0x49, 0x44, 0x33, 0x03}, buf2)
}

func TestRapidUpload(t *testing.T) {
	ctx := setup(t)
	file, err := os.Open("../../../rapid_upload.txt")
	require.NoError(t, err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		pair := scanner.Text()
		if strings.HasPrefix(pair, "# ") {
			continue
		}

		idx := strings.Index(pair, ";")
		if idx < 0 {
			continue
		}

		realPath := strings.TrimSpace(pair[:idx])
		drivePath := strings.TrimSpace(pair[idx+1:])
		if realPath == "" || drivePath == "" {
			continue
		}

		fmt.Printf("Uploading %s to %s ...\n", realPath, drivePath)
		fd, err := os.Open(realPath)
		require.NoError(t, err)
		info, err := fd.Stat()
		require.NoError(t, err)
		_, err = fs.CreateFile(ctx, drivePath, info.Size(), fd, true)
		require.NoError(t, err)
		fd.Close()
	}
}

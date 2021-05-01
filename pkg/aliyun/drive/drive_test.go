package drive

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

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
	names, err := fs.List(ctx, "/test4")
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

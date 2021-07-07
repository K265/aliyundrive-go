// This is a go lang package written for https://www.aliyundrive.com/
package drive

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
)

type Fs interface {
	Get(ctx context.Context, path string, kind string) (*Node, error)
	List(ctx context.Context, path string) ([]Node, error)
	CreateFolder(ctx context.Context, path string) (*Node, error)
	Rename(ctx context.Context, node *Node, newName string) error
	Move(ctx context.Context, node *Node, parent *Node) error
	Remove(ctx context.Context, node *Node) error
	Open(ctx context.Context, node *Node, headers map[string]string) (io.ReadCloser, error)
	CreateFile(ctx context.Context, path string, size int64, in io.Reader, overwrite bool) (*Node, error)
	CalcProof(fileSize int64, in *os.File) (*os.File, string, error)
	CreateFileWithProof(ctx context.Context, path string, size int64, in io.Reader, sha1Code string, proofCode string, overwrite bool) (*Node, error)
	Copy(ctx context.Context, node *Node, parent *Node) error
}

type Config struct {
	RefreshToken string
}

func (config Config) String() string {
	return fmt.Sprintf("Config{RefreshToken: %s}", config.RefreshToken)
}

type Drive struct {
	token
	config     Config
	driveId    string
	rootId     string
	rootNode   Node
	httpClient *http.Client
	mutex      sync.Mutex
}

type token struct {
	accessToken string
	expireAt    int64
}

func (drive *Drive) String() string {
	return fmt.Sprintf("AliyunDrive{driveId: %s}", drive.driveId)
}

func (drive *Drive) request(ctx context.Context, method, url string, headers map[string]string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	req.Header.Set("Referer", "https://www.aliyundrive.com/")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	res, err2 := drive.httpClient.Do(req)
	if err2 != nil {
		return nil, errors.WithStack(err2)
	}

	return res, nil
}

func (drive *Drive) refreshToken(ctx context.Context) error {
	headers := map[string]string{
		"content-type": "application/json;charset=UTF-8",
	}
	data := map[string]string{
		"refresh_token": drive.config.RefreshToken,
		"grant_type":    "refresh_token",
	}

	var body io.Reader
	b, err := json.Marshal(&data)
	if err != nil {
		return errors.New("error marshalling data")
	}

	body = bytes.NewBuffer(b)
	res, err := drive.request(ctx, "POST", "https://auth.aliyundrive.com/v2/account/token", headers, body)
	if err != nil {
		return errors.WithStack(err)
	}
	defer res.Body.Close()

	var token Token
	b, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.Wrap(err, `error reading res.Body`)
	}
	err = json.Unmarshal(b, &token)
	if err != nil {
		return errors.Wrapf(err, "error parsing responseModel, response: %s", string(b))
	}

	drive.accessToken = token.AccessToken
	drive.expireAt = token.ExpiresIn + time.Now().Unix()
	return nil
}

func (drive *Drive) jsonRequest(ctx context.Context, method, url string, requestModel interface{}, responseModel interface{}) error {

	//Token expired, refresh access
	if drive.expireAt < time.Now().Unix() {
		err := drive.refreshToken(ctx)
		if err != nil {
			return errors.WithStack(err)
		}
	}
	headers := map[string]string{
		"content-type":  "application/json;charset=UTF-8",
		"authorization": "Bearer " + drive.accessToken,
	}

	var bodyBytes []byte
	if requestModel != nil {
		b, err := json.Marshal(requestModel)
		if err != nil {
			return errors.New("error marshalling requestModel")
		}
		bodyBytes = b
	}

	res, err := drive.request(ctx, method, url, headers, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return errors.WithStack(err)
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return errors.Errorf(`error request "%s", getting "%d"`, url, res.StatusCode)
	}

	if responseModel != nil {
		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return errors.Wrap(err, `error reading res.Body`)
		}
		err = json.Unmarshal(b, &responseModel)
		if err != nil {
			return errors.Wrapf(err, "error parsing responseModel, response: %s", string(b))
		}
	}

	return nil
}

func NewFs(ctx context.Context, config *Config) (Fs, error) {
	client := &http.Client{}
	drive := &Drive{
		config:     *config,
		httpClient: client,
	}

	// get driveId
	{
		var user User
		data := map[string]string{}
		err := drive.jsonRequest(ctx, "POST", "https://api.aliyundrive.com/v2/user/get", &data, &user)
		if err != nil {
			return nil, errors.Wrap(err, "error getting driveId")
		}

		drive.driveId = user.DriveId
		drive.rootId = "root"
		drive.rootNode = Node{
			NodeId: "root",
			Type:   "folder",
			Name:   "root",
		}
	}

	return drive, nil
}

// path must start with "/" and must not end with "/"
func normalizePath(s string) string {
	separator := "/"
	if !strings.HasPrefix(s, separator) {
		s = separator + s
	}

	if len(s) > 1 && strings.HasSuffix(s, separator) {
		s = s[:len(s)-1]
	}
	return s
}

func (drive *Drive) listNodes(ctx context.Context, node *Node) ([]Node, error) {
	url := "https://api.aliyundrive.com/v2/file/list"
	data := map[string]interface{}{
		"drive_id":       drive.driveId,
		"parent_file_id": node.NodeId,
		"limit":          200,
		"marker":         "",
	}
	var nodes []Node
	var lNodes *ListNodes
	for {
		if lNodes != nil && lNodes.NextMarker == "" {
			break
		}

		err := drive.jsonRequest(ctx, "POST", url, &data, &lNodes)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		nodes = append(nodes, lNodes.Items...)
		data["marker"] = lNodes.NextMarker
	}

	return nodes, nil
}

const FolderKind = "folder"
const FileKind = "file"
const AnyKind = "any"

func (drive *Drive) findNameNode(ctx context.Context, node *Node, name string, kind string) (*Node, error) {
	nodes, err := drive.listNodes(ctx, node)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for _, d := range nodes {
		if d.Name == name && (kind == AnyKind || d.Type == kind) {
			return &d, nil
		}
	}

	return nil, errors.Errorf(`can't find "%s", kind: "%s" under "%s"`, name, kind, node)
}

// https://help.aliyun.com/document_detail/175927.html#pdsgetfilebypathrequest
func (drive *Drive) Get(ctx context.Context, path string, kind string) (*Node, error) {
	path = normalizePath(path)

	if path == "/" || path == "" {
		return &drive.rootNode, nil
	}

	url := "https://api.aliyundrive.com/v2/file/get_by_path"
	data := map[string]interface{}{
		"drive_id":  drive.driveId,
		"file_path": path,
	}

	var node *Node
	err := drive.jsonRequest(ctx, "POST", url, &data, &node)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if node.Type == kind {
		return node, nil
	}

	i := strings.LastIndex(path, "/")
	parent := path[:i]
	name := path[i+1:]

	parentNode, err := drive.Get(ctx, parent, FolderKind)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return drive.findNameNode(ctx, parentNode, name, kind)
}

func findNodeError(err error, path string) error {
	return errors.Wrapf(err, `error finding node of "%s"`, path)
}

func (drive *Drive) List(ctx context.Context, path string) ([]Node, error) {
	path = normalizePath(path)
	node, err := drive.Get(ctx, path, FolderKind)
	if err != nil {
		return nil, findNodeError(err, path)
	}

	nodes, err2 := drive.listNodes(ctx, node)
	if err2 != nil {
		return nil, errors.Wrapf(err2, `error listing nodes of "%s"`, node)
	}

	return nodes, nil
}

func (drive *Drive) createFolderInternal(ctx context.Context, parent string, name string) (*Node, error) {
	drive.mutex.Lock()
	defer drive.mutex.Unlock()

	node, err := drive.Get(ctx, parent+"/"+name, FolderKind)
	if err == nil {
		return node, nil
	}

	node, err = drive.Get(ctx, parent, FolderKind)
	if err != nil {
		return nil, findNodeError(err, parent)
	}
	body := map[string]string{
		"drive_id":        drive.driveId,
		"check_name_mode": "refuse",
		"name":            name,
		"parent_file_id":  node.NodeId,
		"type":            "folder",
	}
	var createdNode Node
	err = drive.jsonRequest(ctx, "POST", "https://api.aliyundrive.com/v2/file/create_with_proof", &body, &createdNode)
	if err != nil {
		return nil, errors.Wrap(err, "error posting create folder request")
	}
	createdNode.Name = name
	return &createdNode, nil
}

func (drive *Drive) CreateFolder(ctx context.Context, path string) (*Node, error) {
	path = normalizePath(path)
	pathLen := len(path)
	i := 0
	var createdNode *Node
	for i < pathLen {
		parent := path[:i]
		remain := path[i+1:]
		j := strings.Index(remain, "/")
		if j < 0 {
			// already at last position
			return drive.createFolderInternal(ctx, parent, remain)
		} else {
			node, err := drive.createFolderInternal(ctx, parent, remain[:j])
			if err != nil {
				return nil, err
			}
			i = i + j + 1
			createdNode = node
		}
	}

	return createdNode, nil
}

func (drive *Drive) checkRoot(node *Node) error {
	if node == nil {
		return errors.New("empty node")
	}
	if node.NodeId == drive.rootId {
		return errors.New("can't operate on root ")
	}
	return nil
}

func (drive *Drive) Rename(ctx context.Context, node *Node, newName string) error {
	if err := drive.checkRoot(node); err != nil {
		return err
	}

	data := map[string]interface{}{
		"check_name_mode": "refuse",
		"drive_id":        drive.driveId,
		"file_id":         node.NodeId,
		"name":            newName,
	}
	err := drive.jsonRequest(ctx, "POST", "https://api.aliyundrive.com/v2/file/update", &data, nil)
	if err != nil {
		return errors.Wrap(err, `error posting rename request`)
	}
	return nil
}

func (drive *Drive) Move(ctx context.Context, node *Node, parent *Node) error {
	if err := drive.checkRoot(node); err != nil {
		return err
	}

	if parent == nil {
		return errors.New("parent node is empty")
	}
	body := map[string]string{
		"drive_id":          drive.driveId,
		"file_id":           node.NodeId,
		"to_parent_file_id": parent.NodeId,
	}
	err := drive.jsonRequest(ctx, "POST", "https://api.aliyundrive.com/v2/file/move", &body, nil)
	if err != nil {
		return errors.Wrap(err, `error posting move request`)
	}
	return nil
}

func (drive *Drive) Remove(ctx context.Context, node *Node) error {
	if err := drive.checkRoot(node); err != nil {
		return err
	}

	body := map[string]string{
		"drive_id": drive.driveId,
		"file_id":  node.NodeId,
	}

	err := drive.jsonRequest(ctx, "POST", "https://api.aliyundrive.com/v2/recyclebin/trash", &body, nil)
	if err != nil {
		return errors.Wrap(err, `error posting remove request`)
	}
	return nil
}

func (drive *Drive) getDownloadUrl(ctx context.Context, node *Node) (*DownloadUrl, error) {
	var detail DownloadUrl
	data := map[string]string{
		"drive_id": drive.driveId,
		"file_id":  node.NodeId,
	}
	err := drive.jsonRequest(ctx, "POST", "https://api.aliyundrive.com/v2/file/get_download_url", &data, &detail)
	if err != nil {
		return nil, errors.Wrapf(err, "error getting node detail, node: %s", node)
	}
	return &detail, nil
}

func (drive *Drive) Open(ctx context.Context, node *Node, headers map[string]string) (io.ReadCloser, error) {
	if err := drive.checkRoot(node); err != nil {
		return nil, err
	}

	downloadUrl, err := drive.getDownloadUrl(ctx, node)
	if err != nil {
		return nil, err
	}

	url := downloadUrl.Url
	if url == "" {
		return nil, errors.Errorf(`error getting url of "%s"`, node)
	}

	res, err := drive.request(ctx, "GET", url, headers, nil)
	if err != nil {
		return nil, errors.Wrapf(err, `error downloading "%s"`, url)
	}

	return res.Body, nil
}

func CalcSha1(in *os.File) (*os.File, string, error) {
	h := sha1.New()
	_, err := io.Copy(h, in)
	if err != nil {
		return nil, "", errors.Wrap(err, "error calculating sha1")
	}

	_, _ = in.Seek(0, 0)
	return in, fmt.Sprintf("%X", h.Sum(nil)), nil
}

func calcProof(accessToken string, fileSize int64, in *os.File) (*os.File, string, error) {
	start := CalcProofOffset(accessToken, fileSize)
	sret, _ := in.Seek(start, 0)
	if sret != start {
		_, _ = in.Seek(0, 0)
		return in, "", errors.Errorf("error seeking file to %d", start)
	}

	buf := make([]byte, 8)
	_, _ = in.Read(buf)
	proofCode := base64.StdEncoding.EncodeToString(buf)
	_, _ = in.Seek(0, 0)
	return in, proofCode, nil
}

func (drive *Drive) CalcProof(fileSize int64, in *os.File) (*os.File, string, error) {
	return calcProof(drive.accessToken, fileSize, in)
}

func (drive *Drive) CreateFile(ctx context.Context, path string, size int64, in io.Reader, overwrite bool) (*Node, error) {
	sha1Code := ""
	proofCode := ""

	fin, ok := in.(*os.File)
	if ok {
		in, sha1Code, _ = CalcSha1(fin)
		in, proofCode, _ = drive.CalcProof(size, fin)
	}

	return drive.CreateFileWithProof(ctx, path, size, in, sha1Code, proofCode, overwrite)
}

func (drive *Drive) CreateFileWithProof(ctx context.Context, path string, size int64, in io.Reader, sha1Code string, proofCode string, overwrite bool) (*Node, error) {
	path = normalizePath(path)
	i := strings.LastIndex(path, "/")
	parent := path[:i]
	name := path[i+1:]
	_, err := drive.CreateFolder(ctx, parent)
	if err != nil {
		return nil, errors.New("error creating folder")
	}

	node, err := drive.Get(ctx, parent, FolderKind)
	if err != nil {
		return nil, findNodeError(err, parent)
	}

	var uploadResult UploadResult

	preUpload := func() error {
		body := map[string]interface{}{
			"check_name_mode": "auto_rename",
			"drive_id":        drive.driveId,
			"name":            name,
			"parent_file_id":  node.NodeId,
			"part_info_list": []map[string]interface{}{
				{
					"part_number": 1,
				},
			},
			"content_hash":      sha1Code,
			"content_hash_name": "sha1",
			"proof_code":        proofCode,
			"proof_version":     "v1",
			"size":              size,
			"type":              "file",
		}
		err := drive.jsonRequest(ctx, "POST", "https://api.aliyundrive.com/v2/file/create_with_proof", &body, &uploadResult)
		if err != nil {
			return errors.Wrap(err, `error posting create file request`)
		}

		if !uploadResult.RapidUpload && len(uploadResult.PartInfoList) < 1 {
			return errors.New(`error extracting uploadUrl'`)
		}

		return nil
	}

	err = preUpload()
	if err != nil {
		return nil, err
	}

	if uploadResult.RapidUpload {
		// rapid upload
		return drive.Get(ctx, path, FileKind)
	}

	if name != uploadResult.FileName && overwrite {
		node, err := drive.Get(ctx, parent+"/"+name, FileKind)
		if err == nil {
			err = drive.Remove(ctx, node)
			if err == nil {
				err = preUpload()
				if err != nil {
					return nil, err
				}
			}
		}
	}

	uploadUrl := uploadResult.PartInfoList[0].UploadUrl
	{
		req, err := http.NewRequestWithContext(ctx, "PUT", uploadUrl, in)
		if err != nil {
			return nil, errors.Wrap(err, "error creating upload request")
		}
		req.Header.Set("Content-Length", strconv.FormatInt(size, 10))
		req.Header.Set("Content-Type", "")
		ursp, err := drive.httpClient.Do(req)
		if err != nil {
			return nil, errors.Wrap(err, "error uploading file")
		}
		defer ursp.Body.Close()
	}

	var createdNode Node
	{
		body := map[string]interface{}{
			"drive_id":  drive.driveId,
			"file_id":   uploadResult.FileId,
			"upload_id": uploadResult.UploadId,
		}

		err := drive.jsonRequest(ctx, "POST", "https://api.aliyundrive.com/v2/file/complete", &body, &createdNode)
		if err != nil {
			return nil, errors.Wrap(err, `error posting upload complete request`)
		}
	}
	return &createdNode, nil
}

// https://help.aliyun.com/document_detail/175927.html#pdscopyfilerequest
func (drive *Drive) Copy(ctx context.Context, node *Node, parent *Node) error {
	if parent == nil {
		return errors.New("parent node is empty")
	}
	body := map[string]string{
		"drive_id":          drive.driveId,
		"file_id":           node.NodeId,
		"to_parent_file_id": parent.NodeId,
	}
	err := drive.jsonRequest(ctx, "POST", "https://api.aliyundrive.com/v2/file/copy", &body, nil)
	if err != nil {
		return errors.Wrap(err, `error posting move request`)
	}

	return nil
}

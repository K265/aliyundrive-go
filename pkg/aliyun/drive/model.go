package drive

import (
	"fmt"
	"time"
)

type Node struct {
	DownloadUrl string `json:"download_url,omitempty"`
	Type        string `json:"type"`                   // folder | file
	Hash        string `json:"content_hash,omitempty"` // sha1
	Name        string `json:"name"`
	NodeId      string `json:"file_id"`
	ParentId    string `json:"parent_file_id,omitempty"`
	Size        int64  `json:"size,omitempty"`
	Updated     string `json:"updated_at"`
}

func (n Node) String() string {
	return fmt.Sprintf("Node{Name: %s, NodeId: %s}", n.Name, n.NodeId)
}

func (n *Node) GetName() string {
	return n.Name
}

func (n *Node) IsDirectory() bool {
	return n.Type == "folder"
}

func (n *Node) GetTime() (time.Time, error) {
	layout := "2006-01-02T15:04:05.000Z"
	t, err := time.Parse(layout, n.Updated)

	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

type Nodes struct {
	Items []Node `json:"items"`
}

type User struct {
	DriveId string `json:"default_drive_id"`
}

type Token struct {
	AccessToken string `json:"access_token"`
}

type DownloadUrl struct {
	Size int64  `json:"size"`
	Url  string `json:"url"`
}

type PartInfo struct {
	UploadUrl string `json:"upload_url"`
}

type UploadResult struct {
	PartInfoList []PartInfo `json:"part_info_list"`
	FileId       string     `json:"file_id"`
	UploadId     string     `json:"upload_id"`
	FileName     string     `json:"file_name"`
}

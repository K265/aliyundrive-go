package drive

import (
	"fmt"
	"time"
)

type NodeId struct {
	NodeId string `json:"file_id"`
}

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

type ListNodes struct {
	Items      []Node `json:"items"`
	NextMarker string `json:"next_marker"`
}

type User struct {
	DriveId string `json:"default_drive_id"`
}

type AlbumInfo struct {
	Data struct {
		DriveId string `json:"driveId"`
	} `json:"data"`
}

type Token struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

type DownloadUrl struct {
	Size       int64             `json:"size"`
	StreamsUrl map[string]string `json:"streams_url,omitempty"`
	Url        string            `json:"url"`
}

type FileProof struct {
	DriveID         string      `json:"drive_id"`
	PartInfoList    []*PartInfo `json:"part_info_list"`
	ParentFileID    string      `json:"parent_file_id"`
	Name            string      `json:"name"`
	Type            string      `json:"type"`
	CheckNameMode   string      `json:"check_name_mode"`
	Size            int64       `json:"size"`
	ContentHash     string      `json:"content_hash"`
	ContentHashName string      `json:"content_hash_name"`
	ProofCode       string      `json:"proof_code"`
	ProofVersion    string      `json:"proof_version"`
}

type PartInfo struct {
	PartNumber int    `json:"part_number"`
	UploadUrl  string `json:"upload_url"`
}

type ProofResult struct {
	PartInfoList []PartInfo `json:"part_info_list,omitempty"`
	Exist        bool       `json:"exist,omitempty"`
	FileId       string     `json:"file_id"`
	RapidUpload  bool       `json:"rapid_upload"`
	UploadId     string     `json:"upload_id"`
	FileName     string     `json:"file_name"`
}

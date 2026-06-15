package model

type HostFileInfo struct {
	Name    string `json:"name"`
	IsDirty bool   `json:"isDirty"`
}

type SaveHostFileRequest struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

package domain

type FileMetadata struct {
	Metadata  []Metadata `json:"metadata"`
	GUID      string     `json:"guid"`
	KEYB4GUID string     `json:"keyb_4_guid"`
}

type Metadata struct {
	Filename string   `json:"filename" csv:"Filename"`
	EnvGUID  string   `json:"env_guid" csv:"EnvGUID"`
	GUID     string   `json:"guid" csv:"GUID"`
	Ext      string   `json:"ext" csv:"Ext"`
	SHA1     string   `json:"sha1" csv:"SHA1"`
	DocType  string   `json:"doc_type" csv:"DocType"`
	Type     string   `json:"type" csv:"Type"`
	Content  []string `json:"content" csv:"-"` // exclude from CSV
}

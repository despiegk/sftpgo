package httpd

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/drakkan/sftpgo/v2/pkg/dataprovider"
	"github.com/go-chi/render"
)

var supportedOnlyOfficeExtensions = []string{
	"doc", "docx", "odt", "ppt", "pptx", "xls", "xlsx", "ods",
}

// only office environment variables
const (
	// ServerAddressEnvKey Key for ServerAddress env variable
	ServerAddressEnvKey = "SFTP_SERVER_ADDR"
	// OnlyOfficeServerAddressEnvKey Key for OnlyOfficeServerAddress env variable
	OnlyOfficeServerAddressEnvKey = "ONLYOFFICE_SERVER_ADDR"
)

type onlyOfficeCallbackData struct {
	Status int    `json:"status"`
	URL    string `json:"url"`
}

type userInfo struct {
	Name string
	ID   string
}

type editOnlyOfficeFilePage struct {
	BaseURL       string
	OnlyOfficeURL string
	FilePath      string
	FileName      string
	FileKey       string
	Ext           string
	Token         string
	User          userInfo
	ShareID       string
	DocumentURL   string
}

type onlyOfficeCallbackResponse struct {
	Error int `json:"error"`
}

func getServerAddress() string {
	return os.Getenv(ServerAddressEnvKey)
}

func getOnlyOfficeServerAddress() string {
	return os.Getenv(OnlyOfficeServerAddressEnvKey)
}

func generateOnlyOfficeFileKey(fileName string, modTime time.Time) string {
	h := sha256.New()
	value := fmt.Sprintf("%s.%d", fileName, modTime.Unix())
	h.Write([]byte(value))
	bs := h.Sum(nil)
	key := hex.EncodeToString(bs[:20])
	return key
}

func checkOnlyOfficeExt(fileName string) bool {
	ext := path.Ext(path.Base(fileName))[1:]
	for _, supportedExt := range supportedOnlyOfficeExtensions {
		if ext == supportedExt {
			return true
		}
	}
	return false
}

func (s *httpdServer) onlyOfficeWriteCallback(w http.ResponseWriter, r *http.Request) {
	var connection *Connection
	var err error

	shareID := r.URL.Query().Get("id")
	if shareID != "" {
		validScopes := []dataprovider.ShareScope{dataprovider.ShareScopeRead, dataprovider.ShareScopeReadWrite}
		_, connection, err = s.checkPublicShare(w, r, validScopes)
		if err != nil {
			return
		}
	} else {
		connection, err = getUserConnection(w, r)
		if err != nil {
			return
		}
	}

	fileName := connection.User.GetCleanedPath(r.URL.Query().Get("path"))

	callbackData := onlyOfficeCallbackData{}

	err = render.DecodeJSON(r.Body, &callbackData)
	if err != nil {
		sendAPIResponse(w, r, err, "", http.StatusBadRequest)
		return
	}

	if callbackData.Status == 2 {
		fs, fsPath, err := connection.GetFsAndResolvedPath(fileName)
		if err != nil {
			sendAPIResponse(w, r, err, fmt.Sprintf("Unable to save file from only office %#v", fileName), getMappedStatusCode(err))
			return
		}

		file, _, _, err := fs.Create(fsPath, os.O_WRONLY|os.O_CREATE, connection.GetCreateChecks(fileName, true))
		if err != nil {
			sendAPIResponse(w, r, err, fmt.Sprintf("Unable to save file from only office %#v", fileName), getMappedStatusCode(err))
			return
		}

		resp, err := http.Get(callbackData.URL)
		if err != nil {
			sendAPIResponse(w, r, err, fmt.Sprintf("Unable to save file from only office %#v", fileName), getMappedStatusCode(err))
			return
		}
		defer resp.Body.Close()

		_, err = io.Copy(file, resp.Body)
		if err != nil {
			sendAPIResponse(w, r, err, fmt.Sprintf("Unable to save file from only office %#v", fileName), getMappedStatusCode(err))
			return
		}
	}

	render.JSON(w, r, onlyOfficeCallbackResponse{Error: 0})
}

package ihaveahugewangUploader

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gabriel-vasile/mimetype"
)

type Response struct {
	FileType string `json:"filetype"`
	Slug     string `json:"slug"`
}

func PrepareUploadBody(fileName string) (string, io.Reader) {
	mime, _ := mimetype.DetectFile(fileName)
	fileReader, _ := os.Open(fileName)
	boundary := fmt.Sprintf("%x", md5.Sum([]byte(strconv.FormatInt(time.Now().UnixNano(), 10))))
	fileHeader := fmt.Sprintf("Content-type: %s", mime.String())
	fileFormat := "--%s\r\nContent-Disposition: form-data; name=\"file\"; filename=\"%s\"\r\n%s\r\n\r\n"
	filePart := fmt.Sprintf(fileFormat, boundary, fileName, fileHeader)
	bodyBottom := fmt.Sprintf("\r\n--%s--\r\n", boundary)
	body := io.MultiReader(strings.NewReader(filePart), fileReader, strings.NewReader(bodyBottom))
	contentType := fmt.Sprintf("multipart/form-data; boundary=%s", boundary)
	return contentType, body
}

func UploadToSite(contentType string, body io.Reader, uploadURL string) (jsonData *Response, err error) {
	response, err := http.Post(uploadURL, contentType, body)
	if err != nil {
		return
	}
	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &jsonData)

	return
}

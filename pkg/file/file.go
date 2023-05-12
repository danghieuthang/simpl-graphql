package file

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/aspose-words-cloud/aspose-words-cloud-go/dev/api"
	"github.com/aspose-words-cloud/aspose-words-cloud-go/dev/api/models"
)

const (
	TemplatePath = "assets/sample.docx"
	ResultPath   = "result.docx"
)

func DownloadFile(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open(TemplatePath) //Open the file to be downloaded later
	defer file.Close()                 // Close after function return

	if err != nil {
		http.Error(w, "File not found.", 404) //return 404 if file is not found
		return
	}
	http.ServeFile(w, r, TemplatePath)

	// tempBuffer := make([]byte, 512)                       //Create a byte array to read the file later
	// file.Read(tempBuffer)                                 //Read the file into  byte
	// FileContentType := http.DetectContentType(tempBuffer) //Get file header

	// FileStat, _ := file.Stat()                         //Get info from file
	// FileSize := strconv.FormatInt(FileStat.Size(), 10) //Get file size as a string

	// fileName := "demo_download"
	// fmt.Print(FileContentType)
	// //Set the headers
	// w.Header().Set("Content-Type", FileContentType+";"+fileName)
	// w.Header().Set("Content-Length", FileSize)

	// file.Seek(0, 0)  //We read 512 bytes from the file already so we reset the offset back to 0
	// io.Copy(w, file) //'Copy' the file to the client
}
func DownloadFormatFile(w http.ResponseWriter, r *http.Request) {
	FormatTemplate(TemplatePath)

	tempBuffer := make([]byte, 512)  //Create a byte array to read the file later
	file, err := os.Open(ResultPath) //Open the file to be downloaded later
	defer file.Close()               // Close after function return
	if err != nil {
		http.Error(w, "Generate fail.", 400) //return 404 if file is not found
		return
	}
	file.Read(tempBuffer) //Read the file into  byte

	FileStat, _ := file.Stat()                         //Get info from file
	FileSize := strconv.FormatInt(FileStat.Size(), 10) //Get file size as a string

	//Set the headers
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", FileSize)
	w.Header().Set("Content-Disposition", "attachment; filename=result.docx")

	file.Seek(0, 0)  //We read 512 bytes from the file already so we reset the offset back to 0
	io.Copy(w, file) //'Copy' the file to the client
}

func FormatTemplate(templatePath string) (*http.Response, error) {
	config, _ := models.NewConfiguration("config.json")
	wordsApi, ctx, _ := api.CreateWordsApi(config)
	requestTemplate, _ := os.Open(templatePath)
	defer requestTemplate.Close()
	DataSourceType := "Json"
	DataSourceName := "persons"
	requestReportEngineSettings := models.ReportEngineSettings{
		DataSourceType: &DataSourceType,
		DataSourceName: &DataSourceName,
	}
	requestDataPayload := `{
		"Title": "Tesststs",
		"Items": [
			{
				"Text": "adsadhsaihi 1"
			},
			{
				"Text": "adsadhsaihi 2"
			},
			{
				"Text": "adsadhsaihi 3"
			}
		]
	}`
	buildReportRequestOptions := map[string]interface{}{}
	buildReportRequest := &models.BuildReportOnlineRequest{
		Template:             requestTemplate,
		Data:                 &requestDataPayload,
		ReportEngineSettings: &requestReportEngineSettings,
		Optionals:            buildReportRequestOptions,
	}
	response, err := wordsApi.BuildReportOnline(ctx, buildReportRequest)
	if err != nil {
		return nil, err

	}
	defer response.Body.Close()
	out, err := os.Create(ResultPath)
	if err != nil {
		return nil, err
	}
	defer out.Close()
	if err != nil {
		return nil, err
	}
	content, err := ioutil.ReadAll(response.Body)
	_ = ioutil.WriteFile(ResultPath, content, 0644)
	return response, nil
}

package steam

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

var steamBaseUrl = "https://api.steampowered.com"

var  ErrNonOkStatus = errors.New("non 200 status code returned")
func GetPublishedFileDetails(fileIds []string) (PublishedFileDetailsResponse, error) {
	fileUrl := fmt.Sprintf("%v/ISteamRemoteStorage/GetPublishedFileDetails/v1/", steamBaseUrl)
	values := url.Values{"itemcount": {strconv.Itoa(len(fileIds))}}
	for i, id := range fileIds {
		values.Add(fmt.Sprintf("publishedfileids[%v]", i), id)
	}

	var response struct{
		Response PublishedFileDetailsResponse `json:"response"`
	}

	resp, err  := http.PostForm(fileUrl, values)
	if err != nil {
		return PublishedFileDetailsResponse{}, fmt.Errorf("Error getting published  file details: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return PublishedFileDetailsResponse{}, ErrNonOkStatus
	}

	data, _ := ioutil.ReadAll(resp.Body)


	err = json.Unmarshal(data, &response)
	return response.Response, err
}
type PublishedFileDetailsResponse struct {
	Result               int `json:"result"`
	ResultCount          int `json:"resultcount"`
	PublishedFileDetails []struct {
		PublishedFileId       string `json:"publishedfileid"`
		Result                int    `json:"result"`
		Creator               string `json:"creator"`
		CreatorAppID          int    `json:"creator_app_id"`
		ConsumerAppID         int    `json:"consumer_app_id"`
		Filename              string `json:"filename"`
		FileSize              int    `json:"file_size"`
		FileURL               string `json:"file_url"`
		HContentFile          string `json:"hcontent_file"`
		PreviewURL            string `json:"preview_url"`
		HContentPreview       string `json:"hcontent_preview"`
		Title                 string `json:"title"`
		Description           string `json:"description"`
		TimeCreated           int    `json:"time_created"`
		TimeUpdated           int    `json:"time_updated"`
		Visibility            int    `json:"visibility"`
		Banned                int    `json:"banned"`
		BanReason             string `json:"ban_reason"`
		Subscriptions         int    `json:"subscriptions"`
		Favorited             int    `json:"favorited"`
		LifetimeSubscriptions int    `json:"lifetime_subscriptions"`
		LifetimeFavorited     int    `json:"lifetime_favorited"`
		Views                 int    `json:"views"`
		Tags                  []struct {
			Tag string `json:"tag"`
		} `json:"tags"`
	} `json:"publishedfiledetails"`
}
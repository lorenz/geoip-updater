package main

import (
	"net/http"
	"net/url"

	"compress/gzip"
	"strconv"

	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"

	"strings"
)

var baseURL, _ = url.Parse("https://updates.maxmind.com/app/update_secure")
var basePath = "."

// UpdateReader gets a reader for a new update if there is one available
func UpdateReader(editionID string, userID uint64, oldMD5 []byte) (io.Reader, error) {
	params := baseURL.Query()
	params.Set("user_id", strconv.FormatUint(userID, 10))
	params.Set("edition_id", editionID)
	params.Set("db_md5", hex.EncodeToString(oldMD5))

	baseURL.RawQuery = params.Encode()

	req := http.Request{
		Method: "GET",
		URL:    baseURL,
	}
	resp, err := http.DefaultClient.Do(&req)
	if err != nil {
		return nil, fmt.Errorf("Failed to request database update: %v", err)
	}

	contentType := resp.Header.Get("content-type")

	if contentType == "application/gzip" {
		decomp, err := gzip.NewReader(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("Failed to start decompression of database: %v", err)
		}
		return decomp, nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("Failed to download response: %v", err)
	}
	if strings.Contains(string(body), "No new updates available") {
		return nil, nil
	}
	return nil, fmt.Errorf("Failed to update database: %v", string(body))
}

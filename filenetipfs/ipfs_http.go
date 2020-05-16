package filenetipfs

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fnv3/test/merkle"
	"github.com/mr-tron/base58"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"time"
)

type IpfsFileInfo struct {
	Name string
	Hash string
	Size string
}

type IpfsLink struct {
	Name   string `json:"Name"`
	Hash   string `json:"Hash"`
	Size   uint64 `json:"Size"`
	Type   int    `json:"Type"`
	Target string `json:"Target"`
}

type IpfsPathList struct {
	Hash  string
	Links []IpfsLink
}

type IpfsBlock struct {
	Objects []IpfsPathList
}

var (
	HttpClient = &http.Client{
		Timeout: 1000 * time.Second,
	}
)

func SaveFileToIpfs(fileName string, file io.Reader) (*merkle.LeafNodes, *IpfsFileInfo, error) {
	body, err := uploadFile(IpfsAddPath, nil, "path", fileName, file)
	if err != nil {
		return nil, nil, err
	}
	var result = new(IpfsFileInfo)
	var list = new(merkle.LeafNodes)
	err = json.Unmarshal(body, result)
	if err != nil {
		return nil, nil, err
	}
	err = hashList(result.Hash, list)
	if err != nil {
		return nil, nil, err
	}
	return list, result,nil
}

func hashList(hash string, list *merkle.LeafNodes) error {
	var Object = new(IpfsBlock)
	url := IpfsLsPath + "?arg=" + hash
	_, err := IpfsHttpPostJson(url, Object, true)
	if err != nil {
		return err
	}
	if len(Object.Objects[0].Links) > 0 {
		for _, hashInfo := range Object.Objects[0].Links {
			err := hashList(hashInfo.Hash, list)
			if err != nil {
				return err
			}
		}
		return nil
	}
	hb, err := base58.Decode(hash)
	if err != nil {
		return err
	}
	h := merkle.LeafNode(hex.EncodeToString(hb[2:]))
	*list = append(*list, &h)
	return nil
}

func uploadFile(url string, params map[string]string, nameField, fileName string, file io.Reader) ([]byte, error) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	formFile, err := writer.CreateFormFile(nameField, fileName)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(formFile, file)
	if err != nil {
		return nil, err
	}
	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())
	resp, err := HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return content, nil
}

func IpfsHttpPostJson(url string, result interface{}, flag bool) ([]byte, error) {
	req, err := http.NewRequest("get", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if flag {
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}
	return body, nil
}

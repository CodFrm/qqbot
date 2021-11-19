package aria2

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/CodFrm/qqbot/config"
	"github.com/CodFrm/qqbot/utils"
	"github.com/golang/glog"
	"github.com/gosuri/uitable"
	"github.com/zyxar/argo/rpc"
)

var defaultRpc rpc.Client

func DefaultRpc() rpc.Client {
	if defaultRpc == nil {
		defaultRpc, _ = rpc.New(context.Background(),
			fmt.Sprintf("http://%s/jsonrpc", config.AppConfig.Aira2.Addr), config.AppConfig.Aira2.Secret,
			time.Second*10, nil)
	}
	return defaultRpc
}

//{"jsonrpc":"2.0","method":"aria2.addUri","id":"QXJpYU5nXzE2MzczMDg0NzBfMC44MzkzNDU2MDYzNDQwNTc=",
//"params":["token:e6KHePx6fdBs",
//["https://ccp-bj29-video-preview.oss-cn-beijing.aliyuncs.com/lt/65DE5B8FE99480FDECD0C8FA8E656AB287455C3B_4550788332__sha1_bj29/FHD
///media.m3u8?di=bj29&dr=1144982&f=60f7c6aab074ac1b836d4d2f9d6f1ec891de8d9d&u=7e626d1b16544779a93a9b119c54ef4e&x-oss-access-key-id
//=LTAIsE5mAn2F493Q&x-oss-additional-headers=referer&x-oss-expires=1637309305&x-oss-process=hls%2Fsign&x-oss-signature=mkZp7H9VQM%2
//Blo7G3rDKW9KQB1K3yA%2FcV84duok0Ty2I%3D&x-oss-signature-version=OSS2"],
//{"header":["sec-ch-ua: \"Microsoft Edge\";v=\"95\", \"Chromium\";v=\"95\", \";Not A Brand\";v=\"99\"","sec-ch-ua-mobile: ?0","sec-ch-ua-platf
//orm: \"Windows\"","Upgrade-Insecure-Requests: 1","User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko)
//Chrome/95.0.4638.69 Safari/537.36 Edg/95.0.1020.53","Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng
//,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9","Sec-Fetch-Site: cross-site","Sec-Fetch-Mode: navigate","Sec-Fetch-User: ?1","Sec-Fetc
//h-Dest: document","Referer: https://www.aliyundrive.com/","Accept-Encoding: gzip, deflate, br","Accept-Language: zh-CN,zh;q=0.9"]}]}

func Download(uri, platform, dir string) (string, error) {
	param := map[string]interface{}{}
	if platform == "阿里云盘" {
		param["header"] = []string{
			"sec-ch-ua: \"Microsoft Edge\";v=\"95\", \"Chromium\";v=\"95\", \";Not A Brand\";v=\"99\"",
			"sec-ch-ua-mobile: ?0",
			"sec-ch-ua-platform: \"Windows\"",
			"Upgrade-Insecure-Requests: 1",
			"User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.69 Safari/537.36 Edg/95.0.1020.53",
			"Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
			"Sec-Fetch-Site: cross-site",
			"Sec-Fetch-Mode: navigate",
			"Sec-Fetch-User: ?1",
			"Sec-Fetch-Dest: document",
			"Referer: https://www.aliyundrive.com/",
			"Accept-Encoding: gzip, deflate, br",
			"Accept-Language: zh-CN,zh;q=0.9",
		}
	}
	dir, _ = filepath.Abs(dir)
	param["dir"] = dir
	glog.Infof("download file %v %+v", uri, param)
	return defaultRpc.AddURI([]string{uri}, param)
}

type DownloadListItems []*DownloadListItem

type DownloadListItem struct {
	Gid      string
	Name     string
	Progress float32
}

func (d DownloadListItems) Table() string {
	table := uitable.New()
	table.AddRow("ID", "NAME", "PROGRESS")
	table.MaxColWidth = 20
	for _, v := range d {
		table.AddRow(v.Gid, v.Name, v.Progress)
	}
	return table.String()
}

func DownloadList() (DownloadListItems, error) {
	list, err := defaultRpc.TellActive("gid", "files", "totalLength", "completedLength")
	if err != nil {
		return nil, err
	}
	ret := make([]*DownloadListItem, 0)
	for _, v := range list {
		ret = append(ret, &DownloadListItem{
			Gid:      v.Gid,
			Name:     filepath.Base(v.Files[0].Path),
			Progress: float32(utils.StringToInt64(v.CompletedLength)) / float32(utils.StringToInt64(v.TotalLength)) * 100,
		})
	}
	return ret, nil
}

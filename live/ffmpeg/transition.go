package transition

import (
	"bufio"
	"io"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/golang/glog"
)

func ToFlv(infile, outfile string, progress chan float32) error {
	//开始执行命令
	cmd := exec.Command("ffmpeg", "-i", infile, "-vcodec", "libx264", "-f", "flv", "-y", outfile)
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	defer stderrPipe.Close()
	if err := cmd.Start(); err != nil {
		return err
	}
	reader := bufio.NewReader(stderrPipe)
	CurrentTime, Duration := 0, 0
	for {
		line, err := reader.ReadBytes('\r')
		if err != nil || err == io.EOF {
			break
		}
		//匹配视频时长
		reg1 := regexp.MustCompile(`Duration:(.*?),`)
		snatch1 := reg1.FindStringSubmatch(string(line))
		if len(snatch1) > 1 {
			Duration = timeEncode(snatch1[1])
		}
		//匹配视频转码进度时间
		reg2 := regexp.MustCompile(`frame=(.*?)fps=(.*?)q=(.*?)size=(.*?)time=(.*?)bitrate=`)
		snatch2 := reg2.FindStringSubmatch(string(line))
		if len(snatch2) > 5 {
			CurrentTime = timeEncode(snatch2[5])
			glog.Infof("%s to %s progress: %.2f", infile, outfile, float32(CurrentTime)/float32(Duration)*100)
			progress <- float32(CurrentTime) / float32(Duration) * 100
		}
	}
	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}

/**
时间解析
*/
func timeEncode(t string) int {
	time := strings.Trim(t, " ")
	hour, _ := strconv.Atoi(time[:2])
	minute, _ := strconv.Atoi(time[3:5])
	second, _ := strconv.Atoi(time[6:8])
	return second + minute*60 + hour*3600
}

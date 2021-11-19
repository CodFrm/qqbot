package live

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/CodFrm/qqbot/cqhttp"
	"github.com/golang/glog"
	"github.com/zhangpeihao/goflv"
	rtmp "github.com/zhangpeihao/gortmp"
)

type live struct {
	ctx context.Context

	url    string
	secret string

	cancel                       context.CancelFunc
	guild                        int64
	channel                      int64
	user                         int64
	createStreamChan             chan rtmp.OutboundStream
	stream                       rtmp.OutboundStream
	conn                         rtmp.OutboundConn
	status                       uint
	videoDataSize, audioDataSize int64
	playqueue                    []string
}

func newLive(guild, channel, user int64, url, secret string) *live {
	live := &live{
		url:              url,
		secret:           secret,
		guild:            guild,
		channel:          channel,
		user:             user,
		createStreamChan: make(chan rtmp.OutboundStream, 100),
		playqueue:        make([]string, 0),
	}

	return live
}

func (h *live) connect() error {
	h.createStreamChan = make(chan rtmp.OutboundStream, 100)
	client, err := rtmp.Dial(h.url, h, 100)
	if err != nil {
		return err
	}
	if err := client.Connect(); err != nil {
		return err
	}
	return nil
}

func (h *live) Close() {
	if h.cancel != nil {
		h.cancel()
	}
	if h.stream != nil {
		h.stream.Close()
	}
	if h.conn != nil {
		h.conn.Close()
	}
}

func (h *live) OnReceived(conn rtmp.Conn, message *rtmp.Message) {

}

func (h *live) OnReceivedRtmpCommand(conn rtmp.Conn, command *rtmp.Command) {

}

func (h *live) OnClosed(conn rtmp.Conn) {
	glog.Infof("OnClose: %v", h.user)
	cqhttp.SendGuildChannelMsg(h.guild, h.channel, "[CQ,at="+strconv.FormatInt(h.user, 10)+"] 推流结束")
	h.conn = nil
	h.stream = nil
	if h.cancel != nil {
		h.cancel()
		h.cancel = nil
	}
}

func (h *live) OnStatus(obConn rtmp.OutboundConn) {
	s, _ := obConn.Status()
	h.status = s
	glog.Infof("OnStatus: %v %v", h.user, s)
}

func (h *live) OnStreamCreated(obConn rtmp.OutboundConn, stream rtmp.OutboundStream) {
	h.conn = obConn
	stream.Attach(h)
	if err := stream.Publish(h.secret, "live"); err != nil {
		cqhttp.SendGuildChannelMsg(h.guild, h.channel, "[CQ,at="+strconv.FormatInt(h.user, 10)+"] 推流失败:"+err.Error())
	}
}

func (h *live) OnPlayStart(stream rtmp.OutboundStream) {

}

func (h *live) OnPublishStart(stream rtmp.OutboundStream) {
	glog.Infof("OnPublishStart: %v", h.user)
	h.stream = stream
	if err := h.play(h.ctx, h.stream); err != nil {
		cqhttp.SendGuildChannelMsg(h.guild, h.channel, "[CQ,at="+strconv.FormatInt(h.user, 10)+"] 播放失败:"+err.Error())
	}
}

func (h *live) Play(filename string) error {
	if h.cancel != nil {
		h.cancel()
		h.cancel = nil
	}
	if h.conn != nil {
		h.conn.Close()
		time.Sleep(time.Second * 2)
	}
	h.ctx, h.cancel = context.WithCancel(context.Background())
	if err := h.connect(); err != nil {
		return err
	}
	h.playqueue = []string{filename}
	return nil
}

func (h *live) popup() string {
	if len(h.playqueue) > 0 {
		ret := h.playqueue[0]
		h.playqueue = h.playqueue[1:]
		return ret
	}
	return ""
}

func (h *live) play(ctx context.Context, stream rtmp.OutboundStream) error {
	h.audioDataSize, h.videoDataSize = 0, 0
	var err error
	filename := h.popup()
	if filename == "" {
		return errors.New("队列无文件")
	}
	flvFile, err := flv.OpenFile("./data/live/flv/" + filename)
	if err != nil {
		glog.Errorf("Open FLV dump file error:", err)
		return err
	}
	go func() {
		defer flvFile.Close()
		go func() {
			<-ctx.Done()
			flvFile.Close()
		}()
		startTs := uint32(0)
		startAt := time.Now().UnixNano()
		preTs := uint32(0)
		for h.status == rtmp.OUTBOUND_CONN_STATUS_CREATE_STREAM_OK {
			if flvFile.IsFinished() {
				glog.Infof("播放完成: %v", filename)
				cqhttp.SendGuildChannelMsg(h.guild, h.channel, "[CQ:at,qq="+strconv.FormatInt(h.user, 10)+"] "+filename+" 播放完成")
				filename := h.popup()
				if filename == "" {
					cqhttp.SendGuildChannelMsg(h.guild, h.channel, "[CQ:at,qq="+strconv.FormatInt(h.user, 10)+"] "+filename+" 播放队列完成")
					return
				}
				tmpFlvFile, err := flv.OpenFile("./data/live/flv/" + filename)
				if err != nil {
					cqhttp.SendGuildChannelMsg(h.guild, h.channel, "[CQ:at,qq="+strconv.FormatInt(h.user, 10)+"] "+filename+" 播放失败,进入下一篇")
					continue
				}
				flvFile.Close()
				flvFile = tmpFlvFile
			}
			header, data, err := flvFile.ReadTag()
			if err != nil {
				glog.Errorf("flvFile.ReadTag(%v) error:", filename, err)
				break
			}
			switch header.TagType {
			case flv.VIDEO_TAG:
				h.videoDataSize += int64(len(data))
			case flv.AUDIO_TAG:
				h.audioDataSize += int64(len(data))
			}

			if startTs == uint32(0) {
				startTs = header.Timestamp
			}
			diff1 := uint32(0)
			if header.Timestamp > startTs {
				diff1 = header.Timestamp - startTs
			}
			if diff1 > preTs {
				preTs = diff1
			}
			if err = stream.PublishData(header.TagType, data, diff1); err != nil {
				glog.Errorf("PublishData(%v) error:", filename, err)
				break
			}
			diff2 := uint32((time.Now().UnixNano() - startAt) / 1000000)
			if diff1 > diff2+100 {
				time.Sleep(time.Millisecond * time.Duration(diff1-diff2))
			}
		}
	}()
	return nil
}

func (h *live) PlayQueue(filename string) error {
	if len(h.playqueue) == 0 {
		return h.Play(filename)
	}
	h.playqueue = append(h.playqueue, filename)
	return nil
}

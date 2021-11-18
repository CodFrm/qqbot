package live

import (
	"context"
	"strconv"
	"time"

	"github.com/CodFrm/qqbot/cqhttp"
	"github.com/golang/glog"
	"github.com/zhangpeihao/goflv"
	rtmp "github.com/zhangpeihao/gortmp"
)

type live struct {
	ctx                          context.Context
	cancel                       context.CancelFunc
	guild                        int64
	channel                      int64
	user                         int64
	createStreamChan             chan rtmp.OutboundStream
	stream                       rtmp.OutboundStream
	conn                         rtmp.OutboundConn
	status                       uint
	videoDataSize, audioDataSize int64
}

func newLive(guild, channel, user int64) *live {
	return &live{
		guild:            guild,
		channel:          channel,
		user:             user,
		createStreamChan: make(chan rtmp.OutboundStream),
	}
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
	cqhttp.SendGuildChannelMsg(h.guild, h.channel, "[CQ,at="+strconv.FormatInt(h.user, 10)+"]推流失败")
}

func (h *live) OnStatus(obConn rtmp.OutboundConn) {
	s, _ := obConn.Status()
	h.status = s
	glog.Infof("OnStatus: %v %v", h.user, s)
}

func (h *live) OnStreamCreated(obConn rtmp.OutboundConn, stream rtmp.OutboundStream) {
	h.createStreamChan <- stream
	h.conn = obConn
}

func (h *live) OnPlayStart(stream rtmp.OutboundStream) {

}

func (h *live) OnPublishStart(stream rtmp.OutboundStream) {
	glog.Infof("OnPublishStart: %v", h.user)
	h.stream = stream
}

func (h *live) Play(filename string) error {
	if h.cancel != nil {
		h.cancel()
		time.Sleep(time.Second * 2)
	}
	h.ctx, h.cancel = context.WithCancel(context.Background())
	return h.play(h.ctx, filename, h.stream)
}

func (h *live) play(ctx context.Context, filename string, stream rtmp.OutboundStream) error {
	h.audioDataSize, h.videoDataSize = 0, 0
	var err error
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
				break
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

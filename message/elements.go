package message

// from https://github.com/Mrs4s/MiraiGo/blob/master/message/elements.go

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"

	"github.com/tidwall/gjson"

	"github.com/LagrangeDev/LagrangeGo/client/packets/pb/message"
	"github.com/LagrangeDev/LagrangeGo/client/packets/pb/service/oidb"
	"github.com/LagrangeDev/LagrangeGo/utils"
	"github.com/LagrangeDev/LagrangeGo/utils/audio"
	"github.com/LagrangeDev/LagrangeGo/utils/crypto"
)

var DefaultThumb, _ = base64.StdEncoding.DecodeString("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")

type (
	TextElement struct {
		Content string
	}

	AtElement struct {
		TargetUin uint32
		TargetUid string
		Display   string
		SubType   AtType
	}

	FaceElement struct {
		FaceID      uint16
		ResultID    uint16 // 猜拳和骰子的值
		isLargeFace bool
	}

	ReplyElement struct {
		ReplySeq  uint32
		SenderUin uint32
		SenderUid string
		GroupUin  uint32 // 私聊回复群聊时
		Time      uint32
		Elements  []IMessageElement
	}

	VoiceElement struct {
		Name string
		Uuid string
		Size uint32
		Url  string
		Md5  []byte
		Sha1 []byte
		Node *oidb.IndexNode

		// --- sending ---
		MsgInfo  *oidb.MsgInfo
		Compat   []byte
		Duration uint32
		Stream   io.ReadSeeker
		Summary  string
	}

	ImageElement struct {
		ImageId  string
		FileUUID string // only in new protocol photo
		Size     uint32
		Width    uint32
		Height   uint32
		Url      string
		SubType  int32

		// EffectID show pic effect id.
		EffectID int32 // deprecated
		Flash    bool

		// send & receive
		Summary string
		Md5     []byte // only in old protocol photo
		IsGroup bool

		Sha1        []byte
		MsgInfo     *oidb.MsgInfo
		Stream      io.ReadSeeker
		CompatFace  *message.CustomFace     // GroupImage
		CompatImage *message.NotOnlineImage // FriendImage
	}

	FileElement struct {
		FileSize uint64
		FileName string
		FileMd5  []byte
		FileUrl  string
		FileId   string // group
		FileUUID string // private
		FileHash string

		// send
		FileStream io.ReadSeeker
		FileSha1   []byte
	}

	ShortVideoElement struct {
		Name     string
		Uuid     []byte
		Size     uint32
		Url      string
		Duration uint32

		// send
		Thumb   *VideoThumb
		Summary string
		Md5     []byte
		Sha1    []byte
		Stream  io.ReadSeeker
		MsgInfo *oidb.MsgInfo
		Compat  *message.VideoFile
	}

	VideoThumb struct {
		Stream io.ReadSeeker
		Size   uint32
		Md5    []byte
		Sha1   []byte
		Width  uint32
		Height uint32
	}

	LightAppElement struct {
		AppName string
		Content string
	}

	ForwardMessage struct {
		IsGroup bool
		SelfId  uint32
		ResID   string
		Nodes   []*ForwardNode
	}

	AtType int
)

const (
	AtTypeGroupMember = 0 // At群成员
)

func NewText(s string) *TextElement {
	return &TextElement{Content: s}
}

func NewAt(target uint32, display ...string) *AtElement {
	dis := "@" + strconv.FormatInt(int64(target), 10)
	if target == 0 {
		dis = "@全体成员"
	}
	if len(display) != 0 {
		dis = display[0]
	}
	return &AtElement{
		TargetUin: target,
		Display:   dis,
	}
}

func NewGroupReply(m *GroupMessage) *ReplyElement {
	return &ReplyElement{
		ReplySeq:  uint32(m.Id),
		SenderUin: m.Sender.Uin,
		Time:      uint32(m.Time),
		Elements:  m.Elements,
	}
}

func NewPrivateReply(m *PrivateMessage) *ReplyElement {
	return &ReplyElement{
		ReplySeq:  uint32(m.Id),
		SenderUin: m.Sender.Uin,
		Time:      uint32(m.Time),
		Elements:  m.Elements,
	}
}

func NewRecord(data []byte, Summary ...string) *VoiceElement {
	return NewStreamRecord(bytes.NewReader(data), Summary...)
}

func NewStreamRecord(r io.ReadSeeker, Summary ...string) *VoiceElement {
	var summary string
	if len(Summary) != 0 {
		summary = Summary[0]
	}
	md5, sha1, length := crypto.ComputeMd5AndSha1AndLength(r)
	info, err := audio.Decode(r)
	if err != nil {
		return &VoiceElement{
			Size:     uint32(length),
			Summary:  summary,
			Stream:   r,
			Md5:      md5,
			Sha1:     sha1,
			Duration: uint32(length),
		}
	}
	return &VoiceElement{
		Size:     uint32(length),
		Summary:  summary,
		Stream:   r,
		Md5:      md5,
		Sha1:     sha1,
		Duration: uint32(info.Time),
	}
}

func NewFileRecord(path string, Summary ...string) (*VoiceElement, error) {
	voice, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return NewStreamRecord(voice, Summary...), nil
}

func NewImage(data []byte, Summary ...string) *ImageElement {
	return NewStreamImage(bytes.NewReader(data), Summary...)
}

func NewStreamImage(r io.ReadSeeker, Summary ...string) *ImageElement {
	var summary string
	if len(Summary) != 0 {
		summary = Summary[0]
	}
	md5, sha1, length := crypto.ComputeMd5AndSha1AndLength(r)
	return &ImageElement{
		Size:    uint32(length),
		Summary: summary,
		Stream:  r,
		Md5:     md5,
		Sha1:    sha1,
	}
}

func NewFileImage(path string, Summary ...string) (*ImageElement, error) {
	img, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return NewStreamImage(img, Summary...), nil
}

func NewVideo(data, thumb []byte, Summary ...string) *ShortVideoElement {
	return NewStreamVideo(bytes.NewReader(data), bytes.NewReader(thumb), Summary...)
}

func NewStreamVideo(r io.ReadSeeker, thumb io.ReadSeeker, Summary ...string) *ShortVideoElement {
	var summary string
	if len(Summary) != 0 {
		summary = Summary[0]
	}
	md5, sha1, length := crypto.ComputeMd5AndSha1AndLength(r)
	return &ShortVideoElement{
		Size:    uint32(length),
		Thumb:   NewVideoThumb(thumb),
		Summary: summary,
		Md5:     md5,
		Sha1:    sha1,
		Stream:  r,
		Compat:  &message.VideoFile{},
	}
}

func NewFileVideo(path string, thumb []byte, Summary ...string) (*ShortVideoElement, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return NewStreamVideo(file, bytes.NewReader(thumb), Summary...), nil
}

func NewVideoThumb(r io.ReadSeeker) *VideoThumb {
	width := uint32(1920)
	height := uint32(1080)
	md5, sha1, size := crypto.ComputeMd5AndSha1AndLength(r)
	_, imgSize, err := utils.ImageResolve(r)
	if err == nil {
		width = uint32(imgSize.Width)
		height = uint32(imgSize.Height)
	}
	return &VideoThumb{
		Stream: r,
		Size:   uint32(size),
		Md5:    md5,
		Sha1:   sha1,
		Width:  width,
		Height: height,
	}
}

func NewFile(data []byte, fileName string) *FileElement {
	return NewStreamFile(bytes.NewReader(data), fileName)
}

func NewStreamFile(r io.ReadSeeker, fileName string) *FileElement {
	md5, sha1, length := crypto.ComputeMd5AndSha1AndLength(r)
	return &FileElement{
		FileName:   fileName,
		FileSize:   length,
		FileStream: r,
		FileMd5:    md5,
		FileSha1:   sha1,
	}
}

func NewLocalFile(path string, name ...string) (*FileElement, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return NewStreamFile(file, utils.LazyTernary(len(name) == 0, func() string {
		return filepath.Base(file.Name())
	}, func() string {
		return name[0]
	})), nil
}

func NewLightApp(content string) *LightAppElement {
	return &LightAppElement{
		AppName: gjson.Get(content, "app").Str,
		Content: content,
	}
}

func NewForward(resid string, nodes []*ForwardNode) *ForwardMessage {
	return &ForwardMessage{
		ResID: resid,
		Nodes: nodes,
	}
}

func NewForwardWithResID(resid string) *ForwardMessage {
	return &ForwardMessage{
		ResID: resid,
	}
}

func NewForwardWithNodes(nodes []*ForwardNode) *ForwardMessage {
	return &ForwardMessage{
		Nodes: nodes,
	}
}

func NewFace(id uint16) *FaceElement {
	return &FaceElement{FaceID: id}
}

func NewDice(value uint16) *FaceElement {
	if value > 6 {
		value = uint16(crypto.RandU32()%3) + 1
	}
	return &FaceElement{
		FaceID:      358,
		ResultID:    value,
		isLargeFace: true,
	}
}

type FingerGuessingType uint16

const (
	FingerGuessingRock     FingerGuessingType = 3 // 石头
	FingerGuessingScissors FingerGuessingType = 2 // 剪刀
	FingerGuessingPaper    FingerGuessingType = 1 // 布
)

func (m FingerGuessingType) String() string {
	switch m {
	case FingerGuessingRock:
		return "石头"
	case FingerGuessingScissors:
		return "剪刀"
	case FingerGuessingPaper:
		return "布"
	}
	return fmt.Sprint(int(m))
}

func NewFingerGuessing(value FingerGuessingType) *FaceElement {
	return &FaceElement{
		FaceID:      359,
		ResultID:    uint16(value),
		isLargeFace: true,
	}
}

func (e *TextElement) Type() ElementType {
	return Text
}

func (e *AtElement) Type() ElementType {
	return At
}

func (e *FaceElement) Type() ElementType {
	return Face
}

func (e *ReplyElement) Type() ElementType {
	return Reply
}

func (e *VoiceElement) Type() ElementType {
	return Voice
}

func (e *ImageElement) Type() ElementType {
	return Image
}

func (e *FileElement) Type() ElementType { return File }

func (e *ShortVideoElement) Type() ElementType {
	return Video
}

func (e *LightAppElement) Type() ElementType {
	return LightApp
}

func (e *ForwardMessage) Type() ElementType {
	return Forward
}

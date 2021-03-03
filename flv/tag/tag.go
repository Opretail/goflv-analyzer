//
// Copyright (c) 2018- yutopp (yutopp@gmail.com)
//
// Distributed under the Boost Software License, Version 1.0. (See accompanying
// file LICENSE_1_0.txt or copy at  https://www.boost.org/LICENSE_1_0.txt)
//

package tag

import (
	"goflv-analyzer/flv/amf0"
	"io"
	"io/ioutil"
)

// ========================================
// FLV tags

// Type ...
type Type uint8

const (
	// TypeAudio ...
	TypeAudio Type = 8
	// TypeVideo ...
	TypeVideo Type = 9
	// TypeScriptData ...
	TypeScriptData Type = 18
)

// FlvTag ...
type FlvTag struct {
	Type
	Timestamp uint32
	StreamID  uint32      // 24bit
	Data      interface{} // *AudioData | *VideoData | *ScriptData
}

// Close ..
func (t *FlvTag) Close() {
	// TODO: wrap an error?
	switch data := t.Data.(type) {
	case *AudioData:
		data.Close()
	case *VideoData:
		data.Close()
	}
}

// ========================================
// Audio tags

// SoundFormat ...
type SoundFormat uint8

const (
	// SoundFormatLinearPCMPlatformEndian ...
	SoundFormatLinearPCMPlatformEndian SoundFormat = 0
	// SoundFormatADPCM ...
	SoundFormatADPCM = 1
	// SoundFormatMP3 ...
	SoundFormatMP3 = 2
	// SoundFormatLinearPCMLittleEndian ...
	SoundFormatLinearPCMLittleEndian = 3
	// SoundFormatNellymoser16kHzMono ...
	SoundFormatNellymoser16kHzMono = 4
	// SoundFormatNellymoser8kHzMono ...
	SoundFormatNellymoser8kHzMono = 5
	// SoundFormatNellymoser ...
	SoundFormatNellymoser = 6
	// SoundFormatG711ALawLogarithmicPCM ...
	SoundFormatG711ALawLogarithmicPCM = 7
	// SoundFormatG711muLawLogarithmicPCM ...
	SoundFormatG711muLawLogarithmicPCM = 8
	// SoundFormatReserved ...
	SoundFormatReserved = 9
	// SoundFormatAAC ...
	SoundFormatAAC = 10
	// SoundFormatSpeex ...
	SoundFormatSpeex = 11
	// SoundFormatMP3_8kHz ...
	SoundFormatMP3_8kHz = 14
	// SoundFormatDeviceSpecificSound ...
	SoundFormatDeviceSpecificSound = 15
)

func (f SoundFormat) String() string {
	switch f {
	case SoundFormatLinearPCMPlatformEndian:
		return "LinearPCMPlatformEndian"
	case SoundFormatADPCM:
		return "ADPCM"
	case SoundFormatMP3:
		return "MP3"
	case SoundFormatLinearPCMLittleEndian:
		return "LinearPCMLittleEndian"
	case SoundFormatNellymoser16kHzMono:
		return "Nellymoser16kHzMono"
	case SoundFormatNellymoser8kHzMono:
		return "Nellymoser8kHzMono"

	case SoundFormatNellymoser:
		return "Nellymoser"
	case SoundFormatG711ALawLogarithmicPCM:
		return "G711ALawLogarithmicPCM"
	case SoundFormatG711muLawLogarithmicPCM:
		return "G711muLawLogarithmicPCM"

	case SoundFormatReserved:
		return "Reserved"

	case SoundFormatAAC:
		return "AAC"
	case SoundFormatSpeex:
		return "Speex"

	case SoundFormatMP3_8kHz:
		return "MP3_8kHz"
	case SoundFormatDeviceSpecificSound:
		return "DeviceSpecificSound"
	}
	return ""
}

// SoundRate ...
type SoundRate uint8

const (
	// SoundRate5_5kHz ...
	SoundRate5_5kHz SoundRate = 0
	// SoundRate11kHz ...
	SoundRate11kHz = 1
	// SoundRate22kHz ...
	SoundRate22kHz = 2
	// SoundRate44kHz ...
	SoundRate44kHz = 3
)

// SoundSize ...
type SoundSize uint8

const (
	// SoundSize8Bit ...
	SoundSize8Bit SoundSize = 0
	// SoundSize16Bit ...
	SoundSize16Bit = 1
)

// SoundType ...
type SoundType uint8

const (
	// SoundTypeMono ...
	SoundTypeMono SoundType = 0
	// SoundTypeStereo ...
	SoundTypeStereo = 1
)

// AudioData ...
type AudioData struct {
	SoundFormat   SoundFormat
	SoundRate     SoundRate
	SoundSize     SoundSize
	SoundType     SoundType
	AACPacketType AACPacketType
	Data          io.Reader
}

func (d *AudioData) Read(buf []byte) (int, error) {
	return d.Read(buf)
}

// Close ...
func (d *AudioData) Close() {
	_, _ = io.Copy(ioutil.Discard, d.Data) //  // TODO: wrap an error?
}

// AACPacketType ...
type AACPacketType uint8

const (
	// AACPacketTypeSequenceHeader ...
	AACPacketTypeSequenceHeader AACPacketType = 0
	// AACPacketTypeRaw ...
	AACPacketTypeRaw = 1
)

// AACAudioData ...
type AACAudioData struct {
	AACPacketType AACPacketType
	Data          io.Reader
}

// ========================================
// Video Tags

// FrameType ...
type FrameType uint8

const (
	// FrameTypeKeyFrame ...
	FrameTypeKeyFrame FrameType = 1
	// FrameTypeInterFrame ...
	FrameTypeInterFrame = 2
	// FrameTypeDisposableInterFrame ...
	FrameTypeDisposableInterFrame = 3
	// FrameTypeGeneratedKeyFrame ...
	FrameTypeGeneratedKeyFrame = 4
	// FrameTypeVideoInfoCommandFrame ...
	FrameTypeVideoInfoCommandFrame = 5
)

func (t FrameType) String() string {
	switch t {
	case FrameTypeKeyFrame:
		return "KeyFrame"
	case FrameTypeInterFrame:
		return "InterFrame"
	case FrameTypeDisposableInterFrame:
		return "DisposableInterFrame"
	case FrameTypeGeneratedKeyFrame:
		return "GeneratedKeyFrame"
	case FrameTypeVideoInfoCommandFrame:
		return "VideoInfoCommandFrame"
	}
	return ""
}

// CodecID ...
type CodecID uint8

const (
	// CodecIDJPEG ...
	CodecIDJPEG CodecID = 1
	// CodecIDSorensonH263 ...
	CodecIDSorensonH263 = 2
	// CodecIDScreenVideo ...
	CodecIDScreenVideo = 3
	// CodecIDOn2VP6 ...
	CodecIDOn2VP6 = 4
	// CodecIDOn2VP6WithAlphaChannel ...
	CodecIDOn2VP6WithAlphaChannel = 5
	// CodecIDScreenVideoVersion2 ...
	CodecIDScreenVideoVersion2 = 6
	// CodecIDAVC ...
	CodecIDAVC = 7
)

func (c CodecID) String() string {
	switch c {
	case CodecIDJPEG:
		return "JPEG"
	case CodecIDSorensonH263:
		return "SorensonH263"
	case CodecIDScreenVideo:
		return "ScreenVideo"
	case CodecIDOn2VP6:
		return "On2VP6"
	case CodecIDOn2VP6WithAlphaChannel:
		return "On2VP6WithAlphaChannel"
	case CodecIDScreenVideoVersion2:
		return "ScreenVideoVersion2"
	case CodecIDAVC:
		return "AVC"
	}
	return ""
}

// VideoData ...
type VideoData struct {
	FrameType       FrameType
	CodecID         CodecID
	AVCPacketType   AVCPacketType
	CompositionTime int32
	Data            io.Reader
}

func (d *VideoData) Read(buf []byte) (int, error) {
	return d.Read(buf)
}

// Close ...
func (d *VideoData) Close() {
	_, _ = io.Copy(ioutil.Discard, d.Data) //  // TODO: wrap an error?
}

// AVCPacketType ...
type AVCPacketType uint8

const (
	// AVCPacketTypeSequenceHeader ...
	AVCPacketTypeSequenceHeader AVCPacketType = 0
	// AVCPacketTypeNALU ...
	AVCPacketTypeNALU = 1
	// AVCPacketTypeEOS ...
	AVCPacketTypeEOS = 2
)

// AVCVideoPacket ...
type AVCVideoPacket struct {
	AVCPacketType   AVCPacketType
	CompositionTime int32
	Data            io.Reader
}

// ========================================
// Data tags

// ScriptData ...
type ScriptData struct {
	// all values are represented as subset of AMF0
	Objects map[string]amf0.ECMAArray
}

package sonic

//#include "sonic.h"
import "C"

import (
	"errors"
	"io"
	"runtime"
	"unsafe"
)

const (
	DEFAULT_SPEED  = 1.0
	DEFAULT_PITCH  = 1.0
	DEFAULT_RATE   = 1.0
	DEFAULT_VOLUME = 1.0
)

type Stream struct {
	sampleRate                 int
	numChannels                int
	sampleSize                 int
	useChordPitch              bool
	speed, pitch, rate, volume float64
	stream                     C.sonicStream
}

func NewStream(sampleRate, numChannels int) *Stream {
	s := &Stream{
		sampleRate:  sampleRate,
		numChannels: numChannels,
		sampleSize:  numChannels * C.sizeof_short,
		speed:       DEFAULT_SPEED,
		pitch:       DEFAULT_PITCH,
		rate:        DEFAULT_RATE,
		volume:      DEFAULT_VOLUME,
		stream:      C.sonicCreateStream(C.int(sampleRate), C.int(numChannels)),
	}

	if s.stream == nil {
		panic("sonicCreateStream returned NULL")
	}

	runtime.SetFinalizer(s, func(s *Stream) { C.sonicDestroyStream(s.stream) })
	return s
}

func (s *Stream) Write(data []byte) (int, error) {
	nSamples := len(data) / s.sampleSize
	if nSamples == 0 {
		return 0, nil
	}
	ok := C.sonicWriteShortToStream(s.stream, (*C.short)(unsafe.Pointer(&data[0])), C.int(nSamples))
	if ok == 0 {
		return 0, errors.New("memory realloc failed")
	}
	return nSamples * s.sampleSize, nil
}

func (s *Stream) Read(data []byte) (int, error) {
	nSamples := len(data) / s.sampleSize
	if nSamples == 0 {
		return 0, io.ErrShortBuffer
	}
	readSamples := C.sonicReadShortFromStream(s.stream, (*C.short)(unsafe.Pointer(&data[0])), C.int(nSamples))
	if readSamples == 0 {
		return 0, io.EOF
	}
	return int(readSamples) * s.sampleSize, nil
}

func (s *Stream) SamplesAvailable() int {
	nSamples := C.sonicSamplesAvailable(s.stream)
	return int(nSamples)
}

func (s *Stream) Speed() float64 {
	return s.speed
}

func (s *Stream) SetSpeed(speed float64) {
	s.speed = speed
	C.sonicSetSpeed(s.stream, C.float(s.speed))
}

func (s *Stream) Pitch() float64 {
	return s.pitch
}

func (s *Stream) SetPitch(pitch float64) {
	s.pitch = pitch
	C.sonicSetPitch(s.stream, C.float(s.pitch))
}

func (s *Stream) Rate() float64 {
	return s.rate
}

func (s *Stream) SetRate(rate float64) {
	s.rate = rate
	C.sonicSetRate(s.stream, C.float(s.rate))
}

func (s *Stream) Volume() float64 {
	return s.volume
}

func (s *Stream) SetVolume(volume float64) {
	s.volume = volume
	C.sonicSetVolume(s.stream, C.float(s.volume))
}

func (s *Stream) NumChannels() int {
	return s.numChannels
}

func (s *Stream) SetNumChannels(numChannels int) {
	s.numChannels = numChannels
	C.sonicSetNumChannels(s.stream, C.int(s.numChannels))
}

func (s *Stream) ChordPitch() bool {
	return s.useChordPitch
}

func (s *Stream) SetChordPitch(useChordPitch bool) {
	s.useChordPitch = useChordPitch
	if s.useChordPitch {
		C.sonicSetChordPitch(s.stream, C.int(1))
	} else {
		C.sonicSetChordPitch(s.stream, C.int(0))
	}
}

func (s *Stream) Flush() int {
	return int(C.sonicFlushStream(s.stream))
}

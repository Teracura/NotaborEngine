package main

import (
	"NotaborEngine/internal/notasound"
)

// AudioFormat specifies the encoding of a sound file.
type AudioFormat = notasound.AudioFormat

const (
	AudioMP3 = notasound.MP3
	AudioWAV = notasound.WAV
	AudioOGG = notasound.OGG
)

// SoundManager handles loading, caching, and playing sound effects.
type SoundManager struct {
	handle *notasound.SoundManager
}

// SetSoundsFolder configures the base folder used when resolving sound file paths.
func (sm *SoundManager) SetSoundsFolder(path string) {
	sm.handle.SetSoundsFolder(path)
}

// Play starts playback of a sound file. If a sound with the same name is already
// playing, it will be stopped and restarted. Loop causes the sound to repeat
// until Stop is called.
func (sm *SoundManager) Play(sound string, format AudioFormat, volume float32, loop bool) error {
	return sm.handle.Play(sound, notasound.AudioFormat(format), volume, loop)
}

// Stop halts a currently playing sound and releases its resources.
func (sm *SoundManager) Stop(sound string) {
	sm.handle.Stop(sound)
}

// Mute globally silences all audio output when set to true.
func (sm *SoundManager) Mute(muted bool) {
	sm.handle.Mute = muted
	sm.handle.UpdateLiveVolume()
}

// Volume sets the master volume level (0.0 = silent, 1.0 = full).
func (sm *SoundManager) Volume(level float32) {
	sm.handle.MasterVolume = level
	sm.handle.UpdateLiveVolume()
}

package event

import "github.com/draganm/emission"

type EventEmitter interface {
	Emit(event interface{}, arguments ...interface{}) EventEmitter
	RemoveListener(event, listener interface{}) EventEmitter
	AddListener(event, listener interface{}) EventEmitter
}

func NewEmitterAdapter() *EmitterAdapter {
	return &EmitterAdapter{emission.NewEmitter()}
}

type EmitterAdapter struct {
	*emission.Emitter
}

func (e *EmitterAdapter) Emit(event interface{}, arguments ...interface{}) EventEmitter {
	e.Emitter.Emit(event, arguments...)
	return e
}

func (e *EmitterAdapter) RemoveListener(event, listener interface{}) EventEmitter {
	e.Emitter.RemoveListener(event, listener)
	return e
}

func (e *EmitterAdapter) AddListener(event, listener interface{}) EventEmitter {
	e.Emitter.AddListener(event, listener)
	return e
}

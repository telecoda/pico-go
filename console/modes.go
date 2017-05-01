package console

func (c *console) initModes() map[ModeType]Mode {
	return map[ModeType]Mode{
		CLI:         newCLIMode(c),
		CODE_EDITOR: newCodeEditorMode(c),
		RUNTIME:     newRuntimeMode(c),
	}
}

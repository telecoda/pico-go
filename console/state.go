package console

import (
	"encoding/json"

	"io/ioutil"
	"os"
)

type ConsoleState struct {
	X int
	Y int
	W int
	H int
}

type Persister interface {
	SaveState(console Console) error
	LoadState() (*ConsoleState, error)
}

type stateManager struct {
	dir string
}

func NewStateManager() (Persister, error) {

	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	mgr := &stateManager{
		dir: dir,
	}
	return mgr, nil
}

func (s *stateManager) SaveState(c Console) error {

	bounds := c.GetBounds()
	// window := c.GetWindow()
	// if window == nil {
	// 	return fmt.Errorf("Console window is nil, cannot be saved")
	// }

	// x, y := window.GetPosition()
	// w, h := window.GetSize()

	x := bounds.Min.X
	y := bounds.Min.Y
	w := bounds.Max.X
	h := bounds.Max.Y

	state := ConsoleState{
		X: x,
		Y: y,
		W: w,
		H: h,
	}

	// save state to local dir
	bytes, err := json.Marshal(&state)
	if err != nil {
		return err
	}

	filename := s.getFilename()
	return ioutil.WriteFile(filename, bytes, 0600)
}

func (s *stateManager) LoadState() (*ConsoleState, error) {
	filename := s.getFilename()
	f, _ := os.Open(filename)
	defer f.Close()

	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	state := &ConsoleState{}
	err = json.Unmarshal(bytes, state)
	if err != nil {
		return nil, err
	}

	return state, nil
}

func (s *stateManager) getFilename() string {
	return s.dir + string(os.PathSeparator) + ".pico-go-state"
}

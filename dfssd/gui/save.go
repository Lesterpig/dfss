package gui

import (
	"encoding/json"
	"io/ioutil"
	"strconv"
)

func (w *Window) Save(filename string) {
	data, err := json.Marshal(w.scene)
	if err != nil {
		w.StatusBar().ShowMessage(err.Error())
		return
	}

	err = ioutil.WriteFile(filename, data, 0600)
	if err != nil {
		w.StatusBar().ShowMessage(err.Error())
		return
	}
	w.StatusBar().ShowMessage("Saved file as " + filename)
}

func (w *Window) Open(filename string) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		w.StatusBar().ShowMessage(err.Error())
		return
	}

	newScene := &Scene{}
	err = json.Unmarshal(data, newScene)
	if err != nil {
		w.StatusBar().ShowMessage(err.Error())
		return
	}

	w.scene = newScene
	w.StatusBar().ShowMessage("Imported file from " + filename + " (" + strconv.Itoa(len(w.scene.Events)) + " events)")
	w.initScene()
}

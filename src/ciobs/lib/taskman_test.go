package ciobs

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestLoadTask(t *testing.T) {
	f, e := ioutil.TempDir("", "TtestLoadTask")
	if e != nil {
		panic(e)
	}

	defer os.RemoveAll(f)

	// create task dir
	td := path.Join(f, taskPath)

	if e := os.Mkdir(td, 655); e != nil {
		panic(e)
	}
	
	// list task
}

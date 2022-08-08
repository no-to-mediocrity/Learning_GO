package main

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/spf13/afero"
	//"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestFilecopy(t *testing.T) {
	//assert := assert.New(t)
	//fmt.Printf(assert)
	var g string
	for i := 0; i < 1000; i++ {
		g += gofakeit.LoremIpsumParagraph(1, 10, 20, " ")
	}
	var appFS = afero.NewMemMapFs()
	appFS.MkdirAll("/src/a", 0755)
	afero.WriteFile(appFS, "src/a/b", []byte(g), 0644)
	os.Args = append(os.Args, "-from \"/src/a/b\"")
	os.Args = append(os.Args, "-offset 100")
	os.Args = append(os.Args, "-limit 1000")
	main()
	//Filecopy("src/a/b", "src/a/c", 1000, 0)
	_, err := appFS.Stat("src/a/c")
	if os.IsNotExist(err) {
		t.Errorf("file \"%s\" does not exist.\n", "src/a/c")
	}
}

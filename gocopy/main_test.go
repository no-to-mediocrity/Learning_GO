package main

func Test() TestCopyFunc(t *testing.T) {
	dir, err := os.MkdirTemp("", "gotest")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir) // clean up

	file := filepath.Join(dir, "test.txt")

	
	if err := os.WriteFile(file, []byte("Hello world! Hello world!"), 0666); err != nil {
		log.Fatal(err)
	}
 	//Arrange

	//Act

	//Assert
}

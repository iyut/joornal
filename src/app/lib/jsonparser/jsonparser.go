package jsonparser

import (
  "io"
  "io/ioutil"
  "log"
  "os"
)

// Parser must implement ParseJSON
type Parser interface {
  ParseJSON([]byte) error
}

// Load the JSON config file
func Load(filePath string, p Parser){
  var err error
  input := io.ReadCloser(os.Stdin)

  input, err = os.Open(filePath)

  if err != nil {
    log.Fatalln(err)
  }
  defer input.Close()
  // Read the config file
  jsonBytes, err := ioutil.ReadAll(input)
  if err != nil {
    log.Fatalln(err)
  }

  // Parse the config file
  if err := p.ParseJSON(jsonBytes); err != nil {
    log.Fatalln("Could not parse %q : %v", filePath, err)
  }
}

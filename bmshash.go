package main

import (
  "fmt"
  "os"
  "flag"
  "io/ioutil"
  "path/filepath"
  "crypto/sha256"
  "crypto/md5"
)

type myError struct {
  msg string
}
func (e myError) Error() string {
  return e.msg
}

const BUFSIZE = 1024

func isBmsPath(path string) bool {
  ext := filepath.Ext(path)
  bmsExts := []string{".bms", ".bme", ".bml", ".pms", ".bmson"}
  for _, be := range bmsExts {
    if ext == be {
      return true
    }
  }
  return false
}

func loadBms(path string) (string, error) {
  file, err := os.Open(path)
  if err != nil {
    return "", myError{"BMS open error : " + err.Error()}
  }
  defer file.Close()

  var str string
  buf := make([]byte, BUFSIZE)
  for {
    n, err := file.Read(buf)
    if n == 0 {
      break
    }
    if err != nil {
      return "", myError{"BMS read error : " + err.Error()}
    }

    str += string(buf[:n])
  }
  return str, nil
}

func printBmsHash(path string) error {
  bmsStr, err := loadBms(path)
  if err != nil {
    return err
  }
  fmt.Println("# bms path:", path)
  fmt.Printf("md5: %x\n", md5.Sum([]byte(bmsStr)))
  fmt.Printf("sha256: %x\n", sha256.Sum256([]byte(bmsStr)))
  fmt.Printf("\n")

  return nil
}

func findBmsInDirectory(path string, recursive bool) error {
  files, _ := ioutil.ReadDir(path)
  for _, f := range files {
    if f.IsDir() && recursive {
      err := findBmsInDirectory(filepath.Join(path, f.Name()), true)
      if err != nil {
        return err
      }
    } else if isBmsPath(f.Name()) {
      err := printBmsHash(filepath.Join(path, f.Name()))
      if err != nil {
        return err
      }
    }
  }
  return nil
}

func main() {
  var (
    recuresive = flag.Bool("r", false, "recursively find bms files in directory")
  )
  flag.Parse()

  if len(flag.Args()) != 1 {
    fmt.Println("Usage: bmshash [option] <bmspath/dirpath>")
    os.Exit(1)
  }

  path := flag.Arg(0)
  fInfo, err := os.Stat(path)
  if err != nil {
    fmt.Println("Path is wrong: ", err.Error())
    os.Exit(1)
  }
  if fInfo.IsDir() {
    err = findBmsInDirectory(path, *recuresive)
    if err != nil {
      fmt.Println(err.Error())
      os.Exit(1)
    }
  } else if isBmsPath(path) {
    err = printBmsHash(path)
    if err != nil {
      fmt.Println(err.Error())
      os.Exit(1)
    }
  }
}

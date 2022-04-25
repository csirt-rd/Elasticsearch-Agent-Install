package main

import (
  "os"
  "io"
  "log"
  "fmt"
  "regexp"
  "runtime"
  "os/exec"
  "strings"
  "net/http"
  "io/ioutil"
  "archive/zip"
  "path/filepath"
)

func main() {
  fmt.Println("Installing...")
  cert := `-----BEGIN CERTIFICATE-----
...
-----END CERTIFICATE-----`
  path, err := os.Getwd()
  if err != nil {
    log.Fatalln(err)
  }
  // Get Elastic agent version number
  resp, err := http.Get("https://www.elastic.co/es/downloads/elastic-agent")
  if err != nil {
    log.Fatalln(err)
  }
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    log.Fatalln(err)
  }
  // Regexp to get version number
  r, _ := regexp.Compile("Version:.*[0-9]+\\.[0-9]+\\.[0-9]+")
  // Parse version number
  version := regexp.MustCompile("<.*?>").ReplaceAllString(r.FindString(string(body)), "")[9:]
  // Get OS
  arch := runtime.GOOS
  switch arch {
  case "windows":
    // if windows
    agentzip := "elastic-agent-" + version + "-windows-x86_64.zip"
    agentfolder := filepath.Join(path, agentzip[:len(agentzip)-4])
    // create file
    out, err := os.Create(agentzip)
    if err != nil {
      log.Fatalln(err)
    }
    // get file data
    resp, err := http.Get("https://artifacts.elastic.co/downloads/beats/elastic-agent/" + agentzip)
    if err != nil {
      log.Fatalln(err)
    }
    defer resp.Body.Close()
    // copy http content to local file
    _, err = io.Copy(out, resp.Body)
    if err != nil {
      log.Fatalln(err)
    }
    if err = out.Close(); err != nil {
      log.Fatal(err)
    }
    // unpack file
    _, err = Unzip(agentzip, path)
    if err != nil {
      log.Fatal(err)
    }
    // create file
    out, err = os.Create("C:\\Windows\\System32\\drivers\\etc\\ca.crt")
    if err != nil {
      log.Fatalln(err)
    } else {
      out.WriteString(cert)
    }
    if err = out.Close(); err != nil {
      log.Fatal(err)
    }
    // install agent
    cmd := exec.Command(filepath.Join(agentfolder, "elastic-agent.exe"), "install", "--url=", "--enrollment-token=", "--certificate-authorities=C:\\Windows\\System32\\drivers\\etc\\ca.crt", "-f")
    output, err := cmd.CombinedOutput()
    if err != nil {
      fmt.Println(fmt.Sprint(err) + ": " + string(output))
      return
    }
    fmt.Println(string(output))
    // remove folder
    err = os.RemoveAll(agentfolder)
    if err != nil {
      log.Fatal(err)
    }
    // remove files
    err = os.Remove(agentzip)
    if err != nil {
      log.Fatal(err)
    }
  case "darwin":
    // if MAC
    agentzip := "elastic-agent-" + version + "-darwin-x86_64.tar.gz"
    agentfolder := filepath.Join(path, agentzip[:len(agentzip)-7])
    // create file
    out, err := os.Create(agentzip)
    if err != nil {
      log.Fatalln(err)
    }
    defer out.Close()
    // get file data
    resp, err := http.Get("https://artifacts.elastic.co/downloads/beats/elastic-agent/" + agentzip)
    if err != nil {
      log.Fatalln(err)
    }
    defer resp.Body.Close()
    // copy http content to local file
    _, err = io.Copy(out, resp.Body)
    if err != nil {
      log.Fatalln(err)
    }
    // unpack file
    cmd := exec.Command("tar", "xzvf", agentzip)
    err = cmd.Run()
    if err != nil {
      log.Fatalln(err)
    }
    // create file
    out, err = os.Create("/etc/elastic-ca.crt")
    if err != nil {
      log.Fatalln(err)
    } else {
      out.WriteString(cert)
    }
    if err = out.Close(); err != nil {
      log.Fatal(err)
    }
    // install agent
    cmd = exec.Command(filepath.Join(agentfolder, "elastic-agent"), "install", "--url=", "--enrollment-token=", "--certificate-authorities=/etc/elastic-ca.crt", "-f")
    output, err := cmd.CombinedOutput()
    if err != nil {
      fmt.Println(fmt.Sprint(err) + ": " + string(output))
      return
    }
    fmt.Println(string(output))
    // remove files
    err = os.Remove(agentzip)
    if err != nil {
      log.Fatal(err)
    }
    err = os.Remove("ca.crt")
    if err != nil {
      log.Fatal(err)
    }
    // remove folder
    err = os.RemoveAll(agentfolder)
    if err != nil {
      log.Fatal(err)
    }
  case "linux":
    // if Unix
    agentzip := "elastic-agent-" + version + "-linux-x86_64.tar.gz"
    agentfolder := filepath.Join(path, agentzip[:len(agentzip)-7])
    // create file
    out, err := os.Create(agentzip)
    if err != nil {
      log.Fatalln(err)
    }
    defer out.Close()
    // get file data
    resp, err := http.Get("https://artifacts.elastic.co/downloads/beats/elastic-agent/" + agentzip)
    if err != nil {
      log.Fatalln(err)
    }
    defer resp.Body.Close()
    // copy http content to local file
    _, err = io.Copy(out, resp.Body)
    if err != nil {
      log.Fatalln(err)
    }
    // unpack file
    cmd := exec.Command("tar", "xzvf", agentzip)
    err = cmd.Run()
    if err != nil {
      log.Fatalln(err)
    }
    // create file
    out, err = os.Create("/etc/elastic-ca.crt")
    if err != nil {
      log.Fatalln(err)
    } else {
      out.WriteString(cert)
    }
    if err = out.Close(); err != nil {
      log.Fatal(err)
    }
    // install agent
    cmd = exec.Command(filepath.Join(agentfolder, "elastic-agent"), "install", "--url=", "--enrollment-token=", "--certificate-authorities=/etc/elastic-ca.crt", "-f")
    output, err := cmd.CombinedOutput()
    if err != nil {
      fmt.Println(fmt.Sprint(err) + ": " + string(output))
      return
    }
    fmt.Println(string(output))
    // remove files
    err = os.Remove(agentzip)
    if err != nil {
      log.Fatal(err)
    }
    // remove folder
    err = os.RemoveAll(agentfolder)
    if err != nil {
      log.Fatal(err)
    }
  }
}

func Unzip(src string, destination string) ([]string, error) {
  var filenames []string
  r, err := zip.OpenReader(src)
  if err != nil {
    return filenames, err
  }
  defer r.Close()
  for _, f := range r.File {
    fpath := filepath.Join(destination, f.Name)
    if !strings.HasPrefix(fpath, filepath.Clean(destination)+string(os.PathSeparator)){
      return filenames, fmt.Errorf("%s is an illegal filepath", fpath)
    }
    filenames = append(filenames, fpath)
    if f.FileInfo().IsDir() {
      os.MkdirAll(fpath, os.ModePerm)
      continue
    }
    if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
      return filenames, err
    }
    outFile, err := os.OpenFile(fpath, 
      os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
      f.Mode())
    if err != nil {
      return filenames, err
    }
    rc, err := f.Open()
    if err != nil {
      return filenames, err
    }
    _, err = io.Copy(outFile, rc)
    outFile.Close()
    rc.Close()
    if err != nil {
      return filenames, err
    }
  }
  return filenames, nil
}

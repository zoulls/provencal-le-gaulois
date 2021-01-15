package extRequest

import (
  "fmt"
  "io/ioutil"
  "net/http"
  "strconv"
  "time"
)

const (
  lcid string = "2"
  sgid string = "203"
  suid string = "10"

  // sgid string = "242"
  // suid string = "18"
)

var userAgentList = []string{
  "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:84.0) Gecko/20100101 Firefox/84.0",
  "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36",
  "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.198 Safari/537.36 OPR/72.0.3815.465",
}

var LastSync = time.Now()

func CallExtRequest() (bool, error) {
    userAgent := userAgentList[0]

	cookies, err := entete(userAgent)
    if err != nil {
      return false, err
    }

    idPage, err := rdv(userAgent, cookies)
    if err != nil {
     return false, err
    }

    return hours(userAgent, cookies, idPage)
}

// 1 step
func entete(userAgent string) ([]*http.Cookie, error) {
  url := fmt.Sprintf("https://pastel.diplomatie.gouv.fr/rdvinternet/flux/protected/frameset/entete.xml?sgid=%s&suid=%s&lcid=%s", sgid, suid, lcid)
  method := "GET"

  client := &http.Client{}
  req, err := http.NewRequest(method, url, nil)

  if err != nil {
    return nil, err
  }
  req.Header.Add("sec-ch-ua", "\"Google Chrome\";v=\"87\", \" Not;A Brand\";v=\"99\", \"Chromium\";v=\"87\"")
  req.Header.Add("sec-ch-ua-mobile", "?0")
  req.Header.Add("User-Agent", userAgent)
  req.Header.Add("Accept", "application/xml")

  res, err := client.Do(req)
  if err != nil {
    return nil, err
  }
  defer res.Body.Close()

  return res.Cookies(), nil
}

// 2 step
func rdv(userAgent string, cookies []*http.Cookie) (int64, error) {
  var idPage int64
  url := fmt.Sprintf("https://pastel.diplomatie.gouv.fr/rdvinternet/flux/protected/RDV/prise/prendreRDVCg.xml?lcid=%s&sgid=%s&suid=%s", lcid, sgid, suid)
  method := "GET"

  client := &http.Client{}
  req, err := http.NewRequest(method, url, nil)

  if err != nil {
    return idPage, err
  }
  req.Header.Add("sec-ch-ua", "\"Google Chrome\";v=\"87\", \" Not;A Brand\";v=\"99\", \"Chromium\";v=\"87\"")
  req.Header.Add("sec-ch-ua-mobile", "?0")
  req.Header.Add("User-Agent", userAgent)
  req.Header.Add("Accept", "application/xml")

  for _, c := range cookies {
    req.AddCookie(c)
  }

  res, err := client.Do(req)
  if err != nil {
    return idPage, err
  }
  defer res.Body.Close()

  _, err = ioutil.ReadAll(res.Body)
  if err != nil {
    return idPage, err
  }

  //var data rdvStc
  //reader := bytes.NewReader(body)
  //decoder := xml.NewDecoder(reader)
  //decoder.CharsetReader = charset.NewReaderLabel
  //err = decoder.Decode(&data)

  //err = xml.Unmarshal(body, &data)
  //if err != nil {
  //  fmt.Println(err)
  //  return idPage, err
  //}

  // idPage = data.page.metadata.idPage

  now := time.Now()
  idPage = now.Unix() * 1000

  return idPage, err
}

// 3 step
func hours(userAgent string, cookies []*http.Cookie, idPage int64) (bool, error) {
  idPageStr := strconv.FormatInt(idPage, 10)
  url := fmt.Sprintf("https://pastel.diplomatie.gouv.fr/rdvinternet/flux/protected/RDV/prise/horaires.xml?idPage=%s", idPageStr)
  method := "GET"

  client := &http.Client{}
  req, err := http.NewRequest(method, url, nil)

  if err != nil {
    return false, err
  }
  req.Header.Add("sec-ch-ua", "\"Google Chrome\";v=\"87\", \" Not;A Brand\";v=\"99\", \"Chromium\";v=\"87\"")
  req.Header.Add("sec-ch-ua-mobile", "?0")
  req.Header.Add("User-Agent", userAgent)
  req.Header.Add("Accept", "application/xml")

  for _, c := range cookies {
    req.AddCookie(c)
  }

  res, err := client.Do(req)
  if err != nil {
    return false, err
  }
  defer res.Body.Close()

  body, err := ioutil.ReadAll(res.Body)
  if err != nil {
    return false, err
  }
  fmt.Println(string(body))

  return true, nil
}

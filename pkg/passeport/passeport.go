package passeport

import (
  "fmt"
  "io/ioutil"
  "net/http"
  "time"
)

type rdvStc struct {
  page pageStc `xml:"PAGE"`
}

type pageStc struct {
  lock statut `xml:"LOCK"`
  session statut `xml:"SESSION"`
  metadata metadata `xml:"METADATA"`
  cg value `xml:"CG"`
  ccg value `xml:"CCG"`
  poste poste `xml:"POSTE"`
}

type statut struct {
  statut string `xml:"STATUT"`
}

type metadata struct {
  lang string `xml:"LANG"`
  idPage int64 `xml:"ID_PAGE"`
}

type value struct {
  value string `xml:"VALUE"`
}

type poste struct {
  ferm value `xml:"FERM"`
}

func CheckPasseport() (bool, error) {
	cookies, err := entete()
    if err != nil {
      return false, err
    }

    idPage, err := rdv(cookies)
    check, err := hours(cookies, idPage)

	return check, err
}

// 1 step
func entete() ([]*http.Cookie, error) {
  url := "https://pastel.diplomatie.gouv.fr/rdvinternet/flux/protected/frameset/entete.xml?sgid=203&suid=10&lcid=2"
  method := "GET"

  client := &http.Client {
  }
  req, err := http.NewRequest(method, url, nil)

  if err != nil {
    fmt.Println(err)
    return nil, err
  }
  req.Header.Add("sec-ch-ua", "\"Google Chrome\";v=\"87\", \" Not;A Brand\";v=\"99\", \"Chromium\";v=\"87\"")
  req.Header.Add("sec-ch-ua-mobile", "?0")
  req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36")
  req.Header.Add("Accept", "application/xml")

  res, err := client.Do(req)
  if err != nil {
    fmt.Println(err)
    return nil, err
  }
  defer res.Body.Close()

  return res.Cookies(), nil
}

// 2 step
func rdv(cookies []*http.Cookie) (int64, error) {
  var idPage int64
  url := "https://pastel.diplomatie.gouv.fr/rdvinternet/flux/protected/RDV/prise/prendreRDVCg.xml?lcid=2&sgid=203&suid=10"
  method := "GET"

  client := &http.Client {
  }
  req, err := http.NewRequest(method, url, nil)

  if err != nil {
    fmt.Println(err)
    return idPage, err
  }
  req.Header.Add("sec-ch-ua", "\"Google Chrome\";v=\"87\", \" Not;A Brand\";v=\"99\", \"Chromium\";v=\"87\"")
  req.Header.Add("sec-ch-ua-mobile", "?0")
  req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36")
  req.Header.Add("Accept", "application/xml")

  for _, c := range cookies {
    req.AddCookie(c)
  }

  res, err := client.Do(req)
  if err != nil {
    fmt.Println(err)
    return idPage, err
  }
  defer res.Body.Close()

  body, err := ioutil.ReadAll(res.Body)
  if err != nil {
    fmt.Println(err)
    return idPage, err
  }

  //var data rdvStc
  //
  //reader := bytes.NewReader(body)
  //decoder := xml.NewDecoder(reader)
  //decoder.CharsetReader = charset.NewReaderLabel
  //err = decoder.Decode(&data)
  //
  ////err = xml.Unmarshal(body, &data)
  //if err != nil {
  //  fmt.Println(err)
  //  return idPage, err
  //}
  fmt.Println(string(body))

  // idPage = data.page.metadata.idPage
  now := time.Now()
  idPage = now.Unix()

  return idPage, err
}

// 3 step
func hours(cookies []*http.Cookie, idPage int64) (bool, error) {
  url := fmt.Sprintf("https://pastel.diplomatie.gouv.fr/rdvinternet/flux/protected/RDV/prise/horaires.xml?idPage=%s", idPage)
  method := "GET"

  client := &http.Client {
  }
  req, err := http.NewRequest(method, url, nil)

  if err != nil {
    fmt.Println(err)
    return false, err
  }
  req.Header.Add("sec-ch-ua", "\"Google Chrome\";v=\"87\", \" Not;A Brand\";v=\"99\", \"Chromium\";v=\"87\"")
  req.Header.Add("sec-ch-ua-mobile", "?0")
  req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36")
  req.Header.Add("Accept", "application/xml")

  for _, c := range cookies {
    req.AddCookie(c)
  }

  res, err := client.Do(req)
  if err != nil {
    fmt.Println(err)
    return false, err
  }
  defer res.Body.Close()

  body, err := ioutil.ReadAll(res.Body)
  if err != nil {
    fmt.Println(err)
    return false, err
  }
  fmt.Println(string(body))

  return true, nil
}
package apple

import (
   "bytes"
   "encoding/json"
   "github.com/89z/format"
   "github.com/89z/format/http"
   "github.com/89z/format/xml"
   "io"
   "net/url"
   "os"
   "strconv"
)

type Server_Parameters struct {
   Adam_ID string `json:"adamId"`
   Svc_ID string `json:"svcId"`
}

type Auth []*http.Cookie

func (a Auth) Create(name string) error {
   file, err := format.Create(name)
   if err != nil {
      return err
   }
   defer file.Close()
   enc := json.NewEncoder(file)
   enc.SetIndent("", " ")
   return enc.Encode(a)
}

func Open_Auth(name string) (Auth, error) {
   file, err := os.Open(name)
   if err != nil {
      return nil, err
   }
   defer file.Close()
   var auth Auth
   if err := json.NewDecoder(file).Decode(&auth); err != nil {
      return nil, err
   }
   return auth, nil
}

func (s Signin) Auth() (Auth, error) {
   req, err := http.NewRequest(
      "POST", "https://buy.tv.apple.com/account/web/auth", nil,
   )
   if err != nil {
      return nil, err
   }
   req.AddCookie(s.my_ac_info())
   req.Header.Set("Origin", "https://tv.apple.com")
   res, err := Client.Do(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   return res.Cookies(), nil
}

func (a Auth) media_user_token() *http.Cookie {
   for _, cookie := range a {
      if cookie.Name == "media-user-token" {
         return cookie
      }
   }
   return nil
}

const (
   sf_max = 143499
   sf_min = 143441
   v_max = 58
   v_min = 50
)

var Client = http.Default_Client

type Episode struct {
   Data struct {
      Playables map[string]struct {
         Assets Asset
      }
   }
}

func New_Episode(content_ID string) (*Episode, error) {
   req, err := http.NewRequest(
      "GET", "https://tv.apple.com/api/uts/v3/episodes/" + content_ID, nil,
   )
   if err != nil {
      return nil, err
   }
   req.URL.RawQuery = url.Values{
      "caller": {"web"},
      "locale": {"en-US"},
      "pfm": {"web"},
      "sf": {strconv.Itoa(sf_max)},
      "v": {strconv.Itoa(v_max)},
   }.Encode()
   res, err := Client.Do(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   epi := new(Episode)
   if err := json.NewDecoder(res.Body).Decode(epi); err != nil {
      return nil, err
   }
   return epi, nil
}

func (e Episode) Asset() *Asset {
   for _, play := range e.Data.Playables {
      return &play.Assets
   }
   return nil
}

type Config struct {
   WebBag struct {
      AppIdKey string
   }
}

func New_Config() (*Config, error) {
   res, err := Client.Get("https://amp-account.tv.apple.com/account/web/config")
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   con := new(Config)
   if err := json.NewDecoder(res.Body).Decode(con); err != nil {
      return nil, err
   }
   return con, nil
}

type Asset struct {
   FpsKeyServerQueryParameters Server_Parameters
   FpsKeyServerUrl string
   HlsUrl string
}

type Environment struct {
   Media_API struct {
      Token string // authorization: Bearer eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXV...
   }
}

func New_Environment() (*Environment, error) {
   res, err := Client.Get("https://tv.apple.com")
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   var scan xml.Scanner
   scan.Data, err = io.ReadAll(res.Body)
   if err != nil {
      return nil, err
   }
   scan.Sep = []byte(`"web-tv-app/config/environment"`)
   scan.Scan()
   scan.Sep = []byte("<meta")
   var meta struct {
      Content string `xml:"content,attr"`
   }
   if err := scan.Decode(&meta); err != nil {
      return nil, err
   }
   content, err := url.PathUnescape(meta.Content)
   if err != nil {
      return nil, err
   }
   env := new(Environment)
   if err := json.Unmarshal([]byte(content), env); err != nil {
      return nil, err
   }
   return env, nil
}

func (c Config) Signin(email, password string) (Signin, error) {
   buf := new(bytes.Buffer)
   err := json.NewEncoder(buf).Encode(map[string]string{
     "accountName": email,
     "password": password,
   })
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", "https://idmsa.apple.com/appleauth/auth/signin", buf,
   )
   if err != nil {
      return nil, err
   }
   req.Header = http.Header{
      "Content-Type": {"application/json"},
      "X-Apple-Widget-Key": {c.WebBag.AppIdKey},
   }
   res, err := Client.Do(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   return res.Cookies(), nil
}

func (s Signin) my_ac_info() *http.Cookie {
   for _, cookie := range s {
      if cookie.Name == "myacinfo" {
         return cookie
      }
   }
   return nil
}

type Signin []*http.Cookie

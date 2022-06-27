package amc

import (
   "bytes"
   "github.com/89z/format/http"
   "github.com/89z/format/json"
   "strconv"
   "strings"
)

func (a Auth) Create(name string) error {
   return json.Encode(name, a)
}

func Open_Auth(name string) (*Auth, error) {
   auth := new(Auth)
   err := json.Decode(name, auth)
   if err != nil {
      return nil, err
   }
   return auth, nil
}

var Client = http.Default_Client

func Get_NID(input string) (int64, error) {
   _, nID, found := strings.Cut(input, "--")
   if found {
      input = nID
   }
   return strconv.ParseInt(input, 10, 64)
}

type playback_request struct {
   Ad_Tags struct {
      Lat int `json:"lat"`
      Mode string `json:"mode"`
      PPID int `json:"ppid"`
      Player_Height int `json:"playerHeight"`
      Player_Width int `json:"playerWidth"`
      URL string `json:"url"`
   } `json:"adtags"`
}

func Unauth() (*Auth, error) {
   req, err := http.NewRequest(
      "POST", "https://gw.cds.amcn.com/auth-orchestration-id/api/v1/unauth", nil,
   )
   if err != nil {
      return nil, err
   }
   req.Header = http.Header{
      "X-Amcn-Device-Id": {"!"},
      "X-Amcn-Language": {"en"},
      "X-Amcn-Network": {"amcplus"},
      "X-Amcn-Platform": {"web"},
      "X-Amcn-Tenant": {"amcn"},
   }
   res, err := Client.Do(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   auth := new(Auth)
   if err := json.NewDecoder(res.Body).Decode(auth); err != nil {
      return nil, err
   }
   return auth, nil
}

func (a *Auth) Login(email, password string) error {
   buf := new(bytes.Buffer)
   err := json.NewEncoder(buf).Encode(map[string]string{
      "email": email,
      "password": password,
   })
   if err != nil {
      return err
   }
   req, err := http.NewRequest(
      "POST", "https://gw.cds.amcn.com/auth-orchestration-id/api/v1/login", buf,
   )
   if err != nil {
      return err
   }
   req.Header = http.Header{
      "Authorization": {"Bearer " + a.Data.Access_Token},
      "Content-Type": {"application/json"},
      "X-Amcn-Device-Ad-Id": {"!"},
      "X-Amcn-Language": {"en"},
      "X-Amcn-Network": {"amcplus"},
      "X-Amcn-Platform": {"web"},
      "X-Amcn-Service-Group-Id": {"10"},
      "X-Amcn-Service-Id": {"amcplus"},
      "X-Amcn-Tenant": {"amcn"},
      "X-Ccpa-Do-Not-Sell": {"doNotPassData"},
   }
   res, err := Client.Do(req)
   if err != nil {
      return err
   }
   defer res.Body.Close()
   return json.NewDecoder(res.Body).Decode(a)
}

func (a *Auth) Refresh() error {
   req, err := http.NewRequest(
      "POST",
      "https://gw.cds.amcn.com/auth-orchestration-id/api/v1/refresh",
      nil,
   )
   if err != nil {
      return err
   }
   req.Header.Set("Authorization", "Bearer " + a.Data.Refresh_Token)
   res, err := Client.Do(req)
   if err != nil {
      return err
   }
   defer res.Body.Close()
   return json.NewDecoder(res.Body).Decode(a)
}

type Auth struct {
   Data struct {
      Access_Token string
      Refresh_Token string
   }
}

type Source struct {
   Key_Systems *struct {
      Widevine struct {
         License_URL string
      } `json:"com.widevine.alpha"`
   }
   Src string
   Type string
}

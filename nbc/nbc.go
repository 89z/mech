package nbc

import (
   "bytes"
   "crypto/hmac"
   "crypto/sha256"
   "encoding/hex"
   "encoding/json"
   "github.com/89z/format"
   "github.com/89z/format/m3u"
   "io"
   "net/http"
   "strconv"
   "strings"
   "time"
)

const (
   mpxAccountID = 2410887629
   queryVideo = "31b459298c0bf48c3b6300ee4922eaf2c9bea4be1cb15a7ab1fe210cd210f779"
)

var secretKey = []byte("2b84a073ede61c766e4c0b3f1e656f7f")

// nbc.com/saturday-night-live/video/october-2-owen-wilson/9000199358
func Parse(id string) (uint64, error) {
   return strconv.ParseUint(id, 10, 64)
}

func generateHash(text string, key []byte) string {
   mac := hmac.New(sha256.New, key)
   io.WriteString(mac, text)
   sum := mac.Sum(nil)
   return hex.EncodeToString(sum)
}

type AccessVOD struct {
   ManifestPath string // this is only valid for one minute
}

func NewAccessVOD(guid uint64) (*AccessVOD, error) {
   var body vodRequest
   body.Device = "android"
   body.DeviceID = "android"
   body.ExternalAdvertiserID = "NBC"
   body.Mpx.AccountID = mpxAccountID
   buf := new(bytes.Buffer)
   err := json.NewEncoder(buf).Encode(body)
   if err != nil {
      return nil, err
   }
   addr := []byte("http://access-cloudpath.media.nbcuni.com")
   addr = append(addr, "/access/vod/nbcuniversal/"...)
   addr = strconv.AppendUint(addr, guid, 10)
   req, err := http.NewRequest("POST", string(addr), buf)
   if err != nil {
      return nil, err
   }
   unix := strconv.FormatInt(time.Now().UnixMilli(), 10)
   var auth strings.Builder
   auth.WriteString("NBC-Security key=android_nbcuniversal,version=2.4")
   auth.WriteString(",time=")
   auth.WriteString(unix)
   auth.WriteString(",hash=")
   auth.WriteString(generateHash(unix, secretKey))
   req.Header = http.Header{
      "Authorization": {auth.String()},
      "Content-Type": {"application/json"},
   }
   format.Log.Dump(req)
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   vod := new(AccessVOD)
   if err := json.NewDecoder(res.Body).Decode(vod); err != nil {
      return nil, err
   }
   return vod, nil
}

type Video struct {
   Data struct {
      BonanzaPage struct {
         Analytics struct {
            ConvivaAssetName string
         }
      }
   }
}

func NewVideo(guid uint64) (*Video, error) {
   var body videoRequest
   body.Extensions.PersistedQuery.Sha256Hash = queryVideo
   body.Variables.App = "nbc"
   body.Variables.Name = strconv.FormatUint(guid, 10)
   body.Variables.Platform = "android"
   body.Variables.Type = "VIDEO"
   buf := new(bytes.Buffer)
   err := json.NewEncoder(buf).Encode(body)
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", "https://friendship.nbc.co/v2/graphql", buf,
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("Content-Type", "application/json")
   format.Log.Dump(req)
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   vid := new(Video)
   if err := json.NewDecoder(res.Body).Decode(vid); err != nil {
      return nil, err
   }
   return vid, nil
}

func (v Video) Name() string {
   return v.Data.BonanzaPage.Analytics.ConvivaAssetName
}

type videoRequest struct {
   Extensions struct {
      PersistedQuery struct {
         Sha256Hash string `json:"sha256Hash"`
      } `json:"persistedQuery"`
   } `json:"extensions"`
   Variables struct {
      App string `json:"app"`
      Name string `json:"name"`
      Platform string `json:"platform"`
      Type string `json:"type"`
      UserID string `json:"userId"` // can be empty
   } `json:"variables"`
}

type vodRequest struct {
   Device string `json:"device"`
   DeviceID string `json:"deviceId"`
   ExternalAdvertiserID string `json:"externalAdvertiserId"`
   Mpx struct {
      AccountID int `json:"accountId"`
   } `json:"mpx"`
}

func (s Stream) String() string {
   buf := []byte("ID:")
   buf = strconv.AppendInt(buf, s.ID, 10)
   buf = append(buf, " Resolution:"...)
   buf = append(buf, s.Resolution...)
   buf = append(buf, " Bandwidth:"...)
   buf = strconv.AppendInt(buf, s.Bandwidth, 10)
   buf = append(buf, " Codecs:"...)
   buf = append(buf, s.Codecs...)
   return string(buf)
}

////////////////////////////////////////////////////////////////////////////////

// #EXT-X-STREAM-INF
type Stream struct {
   ID int64
   Resolution string
   Bandwidth int64
   Codecs string
   URI string
}

// #EXTINF
type Information struct {
   FrameRate float64
   URI string
}

func (s Stream) Information() ([]Information, error) {
   req, err := http.NewRequest("GET", s.URI, nil)
   if err != nil {
      return nil, err
   }
   format.Log.Dump(req)
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   forms, err := m3u.Decode(res.Body, "")
   if err != nil {
      return nil, err
   }
   var infos []Information
   for _, form := range forms {
      var info Information
      info.FrameRate, err = strconv.ParseFloat(form["frame-rate"], 64)
      if err != nil {
         return nil, err
      }
      info.URI = form["URI"]
      infos = append(infos, info)
   }
   return infos, nil
}

func (a AccessVOD) Streams() ([]Stream, error) {
   req, err := http.NewRequest("GET", a.ManifestPath, nil)
   if err != nil {
      return nil, err
   }
   format.Log.Dump(req)
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   forms, err := m3u.Decode(res.Body, "")
   if err != nil {
      return nil, err
   }
   var streams []Stream
   for i, form := range forms {
      var stream Stream
      stream.Bandwidth, err = strconv.ParseInt(form["BANDWIDTH"], 10, 64)
      if err != nil {
         return nil, err
      }
      stream.Codecs = form["CODECS"]
      stream.ID = int64(i)
      stream.Resolution = form["RESOLUTION"]
      stream.URI = form["URI"]
      streams = append(streams, stream)
   }
   return streams, nil
}


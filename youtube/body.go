package youtube

import (
   "bytes"
   "encoding/json"
   "github.com/89z/mech"
   "net/http"
)

const goog_API = "AIzaSyAO_FJ2SlqU8Q4STEHLGCilw_Y9_11qcW8"

type Request struct {
   Body Body
   Header *Header
}

type Body struct {
   Content_Check_OK bool `json:"contentCheckOk,omitempty"`
   Context Context `json:"context"`
   Query string `json:"query,omitempty"`
   Racy_Check_OK bool `json:"racyCheckOk,omitempty"`
   Video_ID string `json:"videoId,omitempty"`
   Params []byte `json:"params,omitempty"`
}

func (b Body) Search(query string) (*Search, error) {
   b.Query = query
   filter := New_Filter()
   filter.Type(Type["Video"])
   param := New_Params()
   param.Filter(filter)
   b.Params = param.Marshal()
   buf := new(bytes.Buffer)
   if err := json.NewEncoder(buf).Encode(b); err != nil {
      return nil, err
   }
   req, err := http.NewRequest("POST", origin + "/youtubei/v1/search", buf)
   if err != nil {
      return nil, err
   }
   req.Header.Set("X-Goog-Api-Key", goog_API)
   res, err := HTTP_Client.Do(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   search := new(Search)
   if err := json.NewDecoder(res.Body).Decode(search); err != nil {
      return nil, err
   }
   return search, nil
}

type Client struct {
   Client_Name string `json:"clientName"`
   Client_Version string `json:"clientVersion"`
}

type Context struct {
   Client Client `json:"client"`
}

const android_version = "17.23.35"

var Android = Body{
   Context: Context{
      Client: Client{"ANDROID", android_version},
   },
}

var Android_Embed = Body{
   Context: Context{
      Client: Client{"ANDROID_EMBEDDED_PLAYER", android_version},
   },
}

var Android_Racy = Body{
   Context: Context{
      Client: Client{"ANDROID", android_version},
   },
   Racy_Check_OK: true,
}

var Android_Content = Body{
   Context: Context{
      Client: Client{"ANDROID", android_version},
   },
   Racy_Check_OK: true,
   Content_Check_OK: true,
}

var Mweb = Body{
   Context: Context{
      Client: Client{"MWEB", "2.20220322.05.00"},
   },
}

func (r Request) Player(id string) (*Player, error) {
   r.Body.Video_ID = id
   buf, err := mech.Encode(r.Body)
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest("POST", origin + "/youtubei/v1/player", buf)
   if err != nil {
      return nil, err
   }
   if r.Header != nil {
      req.Header.Set("Authorization", "Bearer " + r.Header.Access_Token)
   } else {
      req.Header.Set("X-Goog-Api-Key", goog_API)
   }
   res, err := HTTP_Client.Do(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   play := new(Player)
   if err := json.NewDecoder(res.Body).Decode(play); err != nil {
      return nil, err
   }
   return play, nil
}

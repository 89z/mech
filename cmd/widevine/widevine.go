package main

import (
   "bytes"
   "github.com/89z/format/http"
   "github.com/89z/mech/widevine"
   "io"
   "os"
   "strings"
)

var http_client = http.Default_Client

type flags struct {
   address string
   client_id string
   header string
   key_id string
   private_key string
}

func (f flags) contents() (widevine.Contents, error) {
   var (
      client widevine.Client
      err error
   )
   client.ID, err = os.ReadFile(f.client_id)
   if err != nil {
      return nil, err
   }
   client.Private_Key, err = os.ReadFile(f.private_key)
   if err != nil {
      return nil, err
   }
   client.Raw = f.key_id
   module, err := client.Key_ID()
   if err != nil {
      return nil, err
   }
   buf, err := module.Marshal()
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", f.address, bytes.NewReader(buf),
   )
   if err != nil {
      return nil, err
   }
   key, val, ok := strings.Cut(f.header, ":")
   if ok {
      req.Header.Set(key, val)
   }
   res, err := http_client.Do(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   buf, err = io.ReadAll(res.Body)
   if err != nil {
      return nil, err
   }
   return module.Unmarshal(buf)
}
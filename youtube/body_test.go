package youtube

import (
   "fmt"
   "os"
   "testing"
   "time"
)

func Test_Search(t *testing.T) {
   search, err := Mweb.Search("oneohtrix point never along")
   if err != nil {
      t.Fatal(err)
   }
   for _, item := range search.Items() {
      fmt.Println(item.CompactVideoRenderer)
   }
}

const android = "zv9NimPx3Es"

func Test_Android(t *testing.T) {
   play, err := Request{Body: Android}.Player(android)
   if err != nil {
      t.Fatal(err)
   }
   if play.PlayabilityStatus.Status != "OK" {
      t.Fatal(play)
   }
}

var android_embeds = []string{
   "HtVdAasjOgU",
   "WaOKSUlf4TM",
}

func Test_Android_Embed(t *testing.T) {
   req := Request{Body: Android_Embed}
   for _, embed := range android_embeds {
      play, err := req.Player(embed)
      if err != nil {
         t.Fatal(err)
      }
      if play.PlayabilityStatus.Status != "OK" {
         t.Fatal(play)
      }
      time.Sleep(time.Second)
   }
}

var android_racys = []string{
   "Cr381pDsSsA",
   "HsUATh_Nc2U", // signatureCipher
   "SZJvDhaSDnc", // url
   "Tq92D6wQ1mg", // url
   "dqRZDebPIGs", // signatureCipher
}

func Test_Android_Racy(t *testing.T) {
   home, err := os.UserHomeDir()
   if err != nil {
      t.Fatal(err)
   }
   req := Request{Body: Android_Racy}
   req.Header, err = Open_Header(home + "/mech/youtube.json")
   if err != nil {
      t.Fatal(err)
   }
   for _, racy := range android_racys {
      play, err := req.Player(racy)
      if err != nil {
         t.Fatal(err)
      }
      if play.PlayabilityStatus.Status != "OK" {
         t.Fatal(play)
      }
      time.Sleep(time.Second)
   }
}

const android_content = "nGC3D_FkCmg"

func Test_Android_Content(t *testing.T) {
   home, err := os.UserHomeDir()
   if err != nil {
      t.Fatal(err)
   }
   req := Request{Body: Android_Content}
   req.Header, err = Open_Header(home + "/mech/youtube.json")
   if err != nil {
      t.Fatal(err)
   }
   play, err := req.Player(android_content)
   if err != nil {
      t.Fatal(err)
   }
   if play.PlayabilityStatus.Status != "OK" {
      t.Fatal(play)
   }
}

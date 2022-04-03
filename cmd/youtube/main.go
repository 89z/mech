package main

import (
   "flag"
   "github.com/89z/mech/youtube"
)

func main() {
   var vid video
   // a
   flag.StringVar(&vid.address, "a", "", "address")
   // b
   flag.StringVar(&vid.id, "b", "", "video ID")
   // e
   flag.BoolVar(&vid.embed, "e", false, "use embed client")
   // f
   flag.IntVar(&vid.height, "f", 720, "target video height")
   // g
   flag.StringVar(&vid.audio, "g", "AUDIO_QUALITY_MEDIUM", "target audio")
   // i
   flag.BoolVar(&vid.info, "i", false, "information only")
   // r
   var refresh bool
   flag.BoolVar(&refresh, "r", false, "create OAuth refresh token")
   // s
   var access bool
   flag.BoolVar(&access, "s", false, "create OAuth access token")
   // t
   flag.BoolVar(&vid.token, "t", false, "use OAuth access token")
   // v
   var verbose bool
   flag.BoolVar(&verbose, "v", false, "verbose")
   flag.Parse()
   if verbose {
      youtube.LogLevel = 1
   }
   if refresh {
      err := doRefresh()
      if err != nil {
         panic(err)
      }
   } else if access {
      err := doAccess()
      if err != nil {
         panic(err)
      }
   } else if vid.id != "" || vid.address != "" {
      err := vid.do()
      if err != nil {
         panic(err)
      }
   } else {
      flag.Usage()
   }
}

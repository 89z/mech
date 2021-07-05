package main

import (
   "flag"
   "fmt"
   "github.com/89z/mech/youtube"
   "net/url"
   "os"
   "strings"
   "time"
)

func numberFormat(d int64, a ...string) string {
   var (
      e = float64(d)
      f int
   )
   for e >= 1000 {
      e /= 1000
      f++
   }
   return fmt.Sprintf("%.1f", e) + a[f]
}

func clean(r rune) rune {
   if strings.ContainsRune(`"*/:<>?\|`, r) {
      return -1
   }
   return r
}

func main() {
   var (
      atag, vtag int
      info bool
   )
   flag.BoolVar(&info, "i", false, "info only")
   flag.IntVar(&atag, "a", 0, "audio (-1 to skip)")
   flag.IntVar(&vtag, "v", 0, "video (-1 to skip)")
   flag.Parse()
   if len(os.Args) == 1 {
      fmt.Println("youtube [flags] [URL]")
      flag.PrintDefaults()
      return
   }
   // check URL
   if flag.NArg() != 1 {
      panic("missing URL")
   }
   arg := flag.Arg(0)
   watch, err := url.Parse(arg)
   if err != nil {
      panic(err)
   }
   id := watch.Query().Get("v")
   play, err := youtube.PlayerAndroid(id)
   if err != nil {
      panic(err)
   }
   // info
   if info {
      getInfo(play)
      return
   }
   // sort
   play.AdaptiveFormats.Sort()
   type formatFn map[bool]func(youtube.Format)bool
   // filter audio
   if atag >= 0 {
      f := formatFn{
         true: func(f youtube.Format) bool {
            return f.Itag == atag
         },
         false: func(f youtube.Format) bool {
            return f.Height == 0
         },
      }[atag > 0]
      form := play.AdaptiveFormats.Filter(f)[0]
      err := download(play, form)
      if err != nil {
         panic(err)
      }
   }
   // filter video
   if vtag >= 0 {
      f := formatFn{
         true: func(f youtube.Format) bool {
            return f.Itag == vtag
         },
         false: func(f youtube.Format) bool {
            return f.Height <= 720
         },
      }[vtag > 0]
      form := play.AdaptiveFormats.Filter(f)[0]
      err := download(play, form)
      if err != nil {
         panic(err)
      }
   }
}

func getInfo(play *youtube.Player) {
   fmt.Println("author:", play.Author)
   fmt.Println("title:", play.Title)
   fmt.Println()
   for _, f := range play.AdaptiveFormats {
      fmt.Printf(
         "itag %v, height %v, %v, %v, %v\n",
         f.Itag,
         f.Height,
         numberFormat(f.Bitrate, "", " kb/s", " mb/s", " gb/s"),
         numberFormat(f.ContentLength, "", " kB", " MB", " GB"),
         f.MimeType,
      )
   }
}

func download(a *youtube.Player, f youtube.Format) error {
   create := strings.Map(clean, a.Author + "-" + a.Title + f.Ext())
   file, err := os.Create(create)
   if err != nil {
      return err
   }
   defer file.Close()
   begin := time.Now()
   if err := f.Write(file); err != nil {
      return err
   }
   fmt.Println(time.Since(begin))
   return nil
}

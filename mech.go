package mech
// github.com/89z

import (
   "bytes"
   "encoding/json"
   "mime"
   "strconv"
   "strings"
)

func Clean(in string) string {
   var out strings.Builder
   for _, r := range in {
      switch r {
      case
      '"',
      '*',
      '/',
      ':',
      '<',
      '>',
      '?',
      '\\',
      '|',
      '’': // github.com/PowerShell/PowerShell/issues/16084
      default:
         out.WriteRune(r)
      }
   }
   return out.String()
}

func ExtensionByType(typ string) (string, error) {
   media, _, err := mime.ParseMediaType(typ)
   if err != nil {
      return "", err
   }
   switch media {
   case "audio/mpeg":
      return ".mp3", nil
   case "audio/mp4":
      return ".m4a", nil
   case "audio/webm":
      return ".weba", nil
   case "video/mp4":
      return ".m4v", nil
   case "video/webm":
      return ".webm", nil
   }
   return "", notFound{typ}
}

type notFound struct {
   value string
}

func (n notFound) Error() string {
   var buf []byte
   buf = strconv.AppendQuote(buf, n.value)
   buf = append(buf, " is not found"...)
   return string(buf)
}

func Encode[T any](value T) (*bytes.Buffer, error) {
   buf := new(bytes.Buffer)
   enc := json.NewEncoder(buf)
   enc.SetIndent("", " ")
   err := enc.Encode(value)
   if err != nil {
      return nil, err
   }
   return buf, nil
}

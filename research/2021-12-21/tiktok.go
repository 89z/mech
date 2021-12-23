package tiktok

import (
   "crypto/md5"
   "encoding/binary"
   "encoding/hex"
   "io"
   "math/bits"
   "net/url"
   "time"
)

const byteTable1 =
   "D6283B717076BE1BA4FE19575E6CBC21B214377D8CA2FA67556A95E3FA6778ED" +
   "8E553389A8CE36B35CD6B26F96C434B96AEC3495C4FA72FFB8428DFBEC70F085" +
   "46D8B2A1E0CEAE4B7DAEA487CEE3AC5155C436ADFCC4EA97706A85376AC868FA" +
   "FEB033B9677ECEE3CC86D69F767489E9DA9C78C595AAB034B3F27DB2A2EDE0B5" +
   "B68895D151D69E7DD1C8F9B770CC9CB692C5FADD9F28DAC7E0CA95B2DA3497CE" +
   "74FA37E97DC4A237FBFAF1CFAA897D55AE87BCF5E96AC468C7FA768514D0D0E5" +
   "CEFF19D6E5D6CCF1F46CE9E789B2B7AE2889BE5EDC876CF751F26778AEB34BA2" +
   "B3213B55F8B376B2CFB3B3FFB35E717DFAFCFFA87DFED89C1BC46AF988B5E5"

func xGorgon(inputBytes []byte) ([]byte, error) {
   data2 := append(inputBytes[:4], 0, 0, 0, 0)
   data2 = append(data2, inputBytes[32:36]...)
   data2 = append(data2, 0, 0, 0, 0)
   unix := time.Now().Unix()
   var tempByte [4]byte
   binary.BigEndian.PutUint32(tempByte[:], uint32(unix))
   data2 = append(data2, tempByte[:]...)
   data3, err := initialize(data2)
   if err != nil {
      return nil, err
   }
   for i := range data3 {
      byte1 := data3[i]
      byte1 = bits.RotateLeft8(byte1, 4)
      if i == len(data3)-1 {
         byte1 ^= data3[0]
      } else {
         byte1 ^= data3[i+1]
      }
      byte2 := ((byte1 & 0x55) * 2) | ((byte1 & 0xAA) / 2)
      byte2 = ((byte2 & 0x33) * 4) | ((byte2 & 0xCC) / 4)
      byte3 := bits.RotateLeft8(byte2, 4)
      byte3 ^= 0xFF
      data3[i] = byte3 ^ 0x14
   }
   return append([]byte{0x3, 0x61, 0x41, 0x10, 0x80, 0x0}, data3...), nil
}

// THIS IS CORRECT
func initialize(data []byte) ([]byte, error) {
   byteTable2, err := hex.DecodeString(byteTable1)
   if err != nil {
      return nil, err
   }
   var myhex byte
   for i := range data {
      var hex1 byte
      if i == 0 {
         hex1 = byteTable2[byteTable2[0] - 1]
         byteTable2[i] = hex1
      } else if i == 1 {
         var temp byte = 0xD6 + 0x28
         hex1 = byteTable2[temp - 1]
         myhex = temp
         byteTable2[i] = hex1
      } else {
         temp := myhex + byteTable2[i]
         hex1 = byteTable2[temp - 1]
         myhex = temp
         byteTable2[i] = hex1
      }
      hex2 := byteTable2[hex1*2 - 1]
      data[i] = hex2 ^ data[i]
   }
   return data, nil
}

func genXGorgon(query string) ([]byte, error) {
   null_md5 := make([]byte, 16)
   obj := md5.New()
   io.WriteString(obj, query)
   sb := obj.Sum(nil)
   sb = append(sb, null_md5...)
   sb = append(sb, null_md5...)
   sb = append(sb, null_md5...)
   return xGorgon(sb)
}

var deviceParams = url.Values{
   "aid": {"1180"},
   "app_name": {"trill"},
   "channel": {"channel"},
   "device_id": {"7031670777339250182"},
   "device_platform": {"android"},
   "device_type": {"ONEPLUS A3010"},
   "iid": {"7032045377013942018"},
   "os_version": {"7.1"},
   "version_code": {"170804"},
}
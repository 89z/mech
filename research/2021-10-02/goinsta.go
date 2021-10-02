package goinsta

import (
   "encoding/json"
   "fmt"
   "io"
   "net/http"
   "strconv"
   "time"
)

func defaultHandler(args ...interface{}) {
   fmt.Println(args...)
}

// Instagram represent the main API handler
//
// We recommend to use Export and Import functions after first Login.
//
// Also you can use SetProxy and UnsetProxy to set and unset proxy.
// Golang also provides the option to set a proxy using HTTP_PROXY env var.
type Instagram struct {
   user string
   pass string
   // id: android-1923fjnma8123
   dID string
   // family id, v4 uuid: 8b13e7b3-28f7-4e05-9474-358c6602e3f8
   fID string
   // uuid: 8493-1233-4312312-5123
   uuid string
   // rankToken
   rankToken string
   // token -- I think this is depricated, as I don't see any csrf tokens being used anymore, but not 100% sure
   token string
   // phone id v4 uuid: fbf767a4-260a-490d-bcbb-ee7c9ed7c576
   pid string
   // ads id: 5b23a92b-3228-4cff-b6ab-3199f531f05b
   adid string
   // pigeonSessionId
   psID string
   // contains header options set by Instagram
   headerOptions map[string]string
   // expiry of X-Mid cookie
   xmidExpiry int64
   // User-Agent
   userAgent string
   // Account stores all personal data of the user and his/her options.
   Account *Account
   c *http.Client
   // Set to true to debug reponses
   Debug bool
   // Non-error message handlers. By default they will be printed out,
   // alternatively you can e.g. pass them to a logger
   infoHandler  func(...interface{})
   warnHandler  func(...interface{})
   debugHandler func(...interface{})
}

// Default
var GalaxyS10 = Device{
   AndroidRelease:   11,
   AndroidVersion:   30,
   Chipset:          "exynos9820",
   CodeName:         "beyond2",
   Manufacturer:     "samsung",
   Model:            "SM-G975F",
   ScreenDpi:        "560dpi",
   ScreenResolution: "1440x2898",
}

// New creates Instagram structure
func New(username, password string) *Instagram {
   insta := &Instagram{
      c: &http.Client{
         Transport: &http.Transport{Proxy: http.ProxyFromEnvironment},
      },
      dID: generateDeviceID(
         generateMD5Hash(username + password),
      ),
      debugHandler: defaultHandler,
      fID:           generateUUID(),
      headerOptions: map[string]string{},
      infoHandler:  defaultHandler,
      pass: password,
      pid:           generateUUID(),
      psID:          "UFS-" + generateUUID() + "-0",
      user: username,
      userAgent:     createUserAgent(GalaxyS10),
      uuid:          generateUUID(),
      warnHandler:  defaultHandler,
      xmidExpiry:    -1,
   }
   for k, v := range defaultHeaderOptions {
      insta.headerOptions[k] = v
   }
   return insta
}

// Export exports selected *Instagram object options to an io.Writer
func (insta *Instagram) ExportIO(writer io.Writer) error {
   config := ConfigFile{
      Account:       insta.Account,
      FamilyID:      insta.fID,
      HeaderOptions: map[string]string{},
      ID:            insta.Account.ID,
      PhoneID:       insta.pid,
      RankToken:     insta.rankToken,
      Token:         insta.token,
      UUID:          insta.uuid,
      User:          insta.user,
      XmidExpiry:    insta.xmidExpiry,
   }
   for key, value := range insta.headerOptions {
      config.HeaderOptions[key] = value
   }
   bytes, err := json.Marshal(config)
   if err != nil {
      return err
   }
   if _, err := writer.Write(bytes); err != nil {
      return err
   }
   return nil
}

// Login performs instagram login sequence in close resemblance to the android
// apk. Password will be deleted after login.
func (insta *Instagram) Login() (err error) {
   err = insta.sync()
   if err != nil {
      return
   }
   return insta.login()
}

func (insta *Instagram) login() error {
   timestamp := strconv.Itoa(int(time.Now().Unix()))
   encrypted := fmt.Sprintf("#PWD_INSTAGRAM:0:%s:%s", timestamp, insta.pass)
   result, err := json.Marshal(
      map[string]interface{}{
         "adid":                insta.adid,
         "country_code":        "[{\"country_code\":\"44\",\"source\":[\"default\"]}]",
         "device_id":           insta.dID,
         "enc_password":        encrypted,
         "google_tokens":       "[]",
         "guid":                insta.uuid,
         "login_attempt_count": 0,
         "phone_id":            insta.fID,
         "username":            insta.user,
      },
   )
   if err != nil {
      return err
   }
   body, _, err := insta.sendRequest(
      &reqOptions{
         Endpoint: urlLogin,
         IsPost:   true,
         Query:    map[string]string{"signed_body": "SIGNATURE." + string(result)},
      },
   )
   if err != nil {
      return err
   }
   return insta.verifyLogin(body)
}

func (insta *Instagram) sync(args ...map[string]string) error {
   query := map[string]string{
      "id":                      insta.uuid,
      "server_config_retrieval": "1",
   }
   data, err := json.Marshal(query)
   if err != nil {
      return err
   }
   _, _, err = insta.sendRequest(
      &reqOptions{
         Endpoint: urlSync,
         Query:    generateSignature(data),
         IsPost:   true,
         IgnoreHeaders: []string{"Authorization"},
      },
   )
   return err
}

func (insta *Instagram) verifyLogin(body []byte) error {
   res := accountResp{}
   err := json.Unmarshal(body, &res)
   if err != nil {
      return fmt.Errorf("failed to parse json from login response %q", err)
   }
   if res.Status != "ok" {
      switch res.ErrorType {
      case "bad_password":
         return ErrBadPassword
      }
      return fmt.Errorf("Failed to login: %v, %v", res.ErrorType, res.Message)
   }
   insta.Account = &res.Account
   insta.Account.insta = insta
   insta.rankToken = strconv.FormatInt(insta.Account.ID, 10) + "_" + insta.uuid
   return nil
}
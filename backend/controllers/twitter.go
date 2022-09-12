package controllers

import (
	"net/http"

	"example.com/satsoverflow-backend/models"
	"github.com/dghubble/gologin/twitter"
)

// func generateNonce() string {
// 	// Does not need to be secure
// 	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
// 	sb := strings.Builder{}
// 	sb.Grow(32)
// 	for i := 0; i < 32; i++ {
// 		sb.WriteByte(charset[rand.Intn(len(charset))])
// 	}
// 	return sb.String()
// }

// func getOauthSig(sigBase string) string {
// 	sb := strings.Builder{}
// 	sb.WriteString(url.QueryEscape(os.Getenv("TWITTER_API_SECRET")))
// 	sb.WriteString("&")
// 	// sb.WriteString(url.QueryEscape(os.Getenv("TWITTER_ACCESS_SECRET")))
// 	key := sb.String()

// 	mac := hmac.New(sha1.New, []byte(key))
// 	mac.Write([]byte(sigBase))
// 	res := mac.Sum(nil)
// 	dst := make([]byte, base64.StdEncoding.EncodedLen(len(res)))
// 	base64.StdEncoding.Encode(dst, res)
// 	return string(dst)
// }

// func getSigningKeyDemo() string {
// 	sb := strings.Builder{}
// 	sb.WriteString(url.QueryEscape("kAcSOqF21Fu85e7zjz7ZN2U4ZRhfV3WpwPAoE3Z7kBw"))
// 	sb.WriteString("&")
// 	sb.WriteString(url.QueryEscape("LswwdoUaIvS8ltyTt5jkRh4J50vUPVVHtR2YPi5kE"))
// 	key := sb.String()

// 	sigBase := "POST&https%3A%2F%2Fapi.twitter.com%2F1.1%2Fstatuses%2Fupdate.json&include_entities%3Dtrue%26oauth_consumer_key%3Dxvz1evFS4wEEPTGEFPHBog%26oauth_nonce%3DkYjzVBB8Y0ZFabxSWbWovY3uYSQ2pTgmZeNu2VS4cg%26oauth_signature_method%3DHMAC-SHA1%26oauth_timestamp%3D1318622958%26oauth_token%3D370773112-GmHxMAgYyLbNEtIKZeRNFsMKPR9EyMZeS9weJAEb%26oauth_version%3D1.0%26status%3DHello%2520Ladies%2520%252B%2520Gentlemen%252C%2520a%2520signed%2520OAuth%2520request%2521"
// 	mac := hmac.New(sha1.New, []byte(key))
// 	mac.Write([]byte(sigBase))
// 	res := mac.Sum(nil)
// 	dst := make([]byte, base64.StdEncoding.EncodedLen(len(res)))
// 	base64.StdEncoding.Encode(dst, res)
// 	return string(dst)
// }

// func getAuthorizationHeader(method string, base_url string) string {
// 	// Gets Oauth1 twitter authorization header
// 	// type oauthHeader struct {
// 	// 	oauth_consumer_key     string
// 	// 	oauth_nonce            string
// 	// 	oauth_signature_method string
// 	// 	oauth_timestamp        string
// 	// 	oauth_version          string
// 	// }
// 	rand.Seed(time.Now().UnixNano())
// 	var oauthFields [][]string
// 	oauthFields = append(oauthFields, []string{"oauth_callback", "localhost%3A8080%2Ftwitter"})
// 	oauthFields = append(oauthFields, []string{"oauth_consumer_key", os.Getenv("TWITTER_API_KEY")})
// 	oauthFields = append(oauthFields, []string{"oauth_nonce", generateNonce()})
// 	oauthFields = append(oauthFields, []string{"oauth_signature_method", "HMAC-SHA1"})
// 	oauthFields = append(oauthFields, []string{"oauth_timestamp", strconv.FormatInt(time.Now().Unix(), 10)})
// 	oauthFields = append(oauthFields, []string{"oauth_version", "1.0"})

// 	sb := strings.Builder{}
// 	sb.WriteString(strings.ToUpper(method))
// 	sb.WriteString("&")
// 	sb.WriteString(url.QueryEscape(base_url))
// 	sb.WriteString("&")
// 	for i, field := range oauthFields {
// 		// urlencode the key
// 		sb.WriteString(url.QueryEscape(field[0]))
// 		sb.WriteString("=")
// 		sb.WriteString(url.QueryEscape(field[1]))
// 		if i+1 < len(oauthFields) {
// 			sb.WriteString("&")
// 		}
// 	}

// 	sigBase := sb.String()
// 	sig := getOauthSig(sigBase)
// 	oauthFields = nil
// 	oauthFields = append(oauthFields, []string{"oauth_consumer_key", os.Getenv("TWITTER_API_KEY")})
// 	oauthFields = append(oauthFields, []string{"oauth_nonce", generateNonce()})
// 	oauthFields = append(oauthFields, []string{"oauth_signature_method", "HMAC-SHA1"})
// 	oauthFields = append(oauthFields, []string{"oauth_timestamp", strconv.FormatInt(time.Now().Unix(), 10)})
// 	oauthFields = append(oauthFields, []string{"oauth_version", "1.0"})
// 	oauthFields = append(oauthFields, []string{"oauth_signature", sig})

// 	sb = strings.Builder{}
// 	sb.WriteString("OAuth ")
// 	for i, field := range oauthFields {
// 		sb.WriteString(url.QueryEscape(field[0]))
// 		sb.WriteString("=\"")
// 		sb.WriteString(url.QueryEscape(field[1]))
// 		sb.WriteString("\"")
// 		if i+1 < len(oauthFields) {
// 			sb.WriteString(", ")
// 		}
// 	}
// 	return sb.String()
// 	// auth := oauthHeader{
// 	// 	oauth_consumer_key:     os.Getenv("TWITTER_API_KEY"),
// 	// 	oauth_nonce:            generateNonce(),
// 	// 	oauth_signature_method: "HMAC-SHA1",
// 	// 	oauth_timestamp:        strconv.FormatInt(time.Now().Unix(), 10),
// 	// 	oauth_version:          "1.0",
// 	// }
// 	// sb := strings.Builder{}
// 	// sb.WriteString("OAuth ")
// 	// sb.WriteString("oauth_consumer_key")
// }
// issueSession issues a cookie session after successful Twitter login
func (server *Server) issueSession() http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		twitterUser, err := twitter.UserFromContext(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// 2. Implement a success handler to issue some form of session
		session, _ := server.Store.Get(req, "sessionID")
		username := session.Values["username"]
		server.DB.Model(&models.User{}).Where("username = ?", username).Update("twitter_handle", twitterUser.ScreenName)
		// session.Values["twitterID"] = twitterUser.ID
		// session.Values["twitterName"] = twitterUser.ScreenName
		// session.Save(req, w)
		http.Redirect(w, req, "/profile", http.StatusFound)
	}
	return http.HandlerFunc(fn)
}

// func (server *Server) TwitterCallback(c *gin.Context) {
// 	config := &oauth1.Config{
// 		ConsumerKey:    os.Getenv("TWITTER_API_KEY"),
// 		ConsumerSecret: os.Getenv("TWITTER_API_SECRET"),
// 		CallbackURL:    "localhost:8080/twitter/callback",
// 		Endpoint:       twitterOAuth1.AuthorizeEndpoint,
// 	}
// 	twitter.CallbackHandler(config, issueSession(), nil)
// }

// func (server *Server) Twitter(c *gin.Context) {
// 	fmt.Println("Hit this function!")
// 	config := &oauth1.Config{
// 		ConsumerKey:    os.Getenv("TWITTER_API_KEY"),
// 		ConsumerSecret: os.Getenv("TWITTER_API_SECRET"),
// 		CallbackURL:    "localhost:8080/twitter/callback",
// 		Endpoint:       twitterOAuth1.AuthorizeEndpoint,
// 	}

// 	twitter.LoginHandler(config, nil)
// 	// requestToken, requestSecret, _ := config.RequestToken()
// 	// authorizationURL, _ := config.AuthorizationURL(requestToken)
// 	// http.Redirect(c.Writer, c.Request, authorizationURL.String(), http.StatusFound)

// 	// requestToken, verifier, _ := oauth1.ParseAuthorizationCallback(c.Request)
// 	// accessToken, accessSecret, _ := config.AccessToken(requestToken, requestSecret, verifier)
// 	// // handle error
// 	// oauth1.NewToken(accessToken, accessSecret)
// 	// params := url.Values{}
// 	// params.Add("oauth_callback", "localhost%3A8080%2Ftwitter")
// 	// resp, err := http.Post("https://api.twitter.com/oauth/request_token", "text/html", nil)
// 	// req_method := "POST"
// 	// req, _ := http.NewRequest(req_method, "https://api.twitter.com/oauth/request_token", nil)
// 	// query := url.Values{}
// 	// query.Add("oauth_callback", "localhost%3A8080%2Ftwitter")
// 	// req.URL.RawQuery = query.Encode()
// 	// // req.Header.Add("Authorization", )

// 	// // apiCreds := oauth.Credentials{
// 	// // 	Token:  os.Getenv("TWITTER_API_KEY"),
// 	// // 	Secret: os.Getenv("TWITTER_API_SECRET"),
// 	// // }
// 	// // accessCreds := oauth.Credentials{
// 	// // 	Token:  os.Getenv("TWITTER_ACCESS_KEY"),
// 	// // 	Secret: os.Getenv("TWITTER_ACCESS_SECRET"),
// 	// // }
// 	// // oauthClient := oauth.Client{Credentials: apiCreds}
// 	// // // fmt.Printf("Creds = %v\n", creds)
// 	// // err := oauthClient.SetAuthorizationHeader(req.Header, &accessCreds, req_method, req.URL, nil)
// 	// // if err != nil {
// 	// // 	log.Printf("Error setting authorization header: %v\n", err)
// 	// // }
// 	// authHeader := getAuthorizationHeader("POST", "https://api.twitter.com/oauth/request_token")
// 	// req.Header.Add("Authorization", authHeader)
// 	// httpClient := &http.Client{}
// 	// resp, err := httpClient.Do(req)
// 	// if err != nil {
// 	// 	log.Printf("Error calling twitter oauth/request_token: %v", err)
// 	// }
// 	// defer resp.Body.Close()
// 	// // var p []byte
// 	// // resp.Body.Read(p)
// 	// fmt.Printf("Response: %s", resp.Status)
// 	// var p []byte
// 	// resp.Body.Read(p)
// 	// fmt.Printf("Response extra: %v", string(p))
// 	// c.JSON(http.StatusOK, gin.H{
// 	// 	"message": req.Header.Get("Authorization"),
// 	// 	"url":     req.URL.String(),
// 	// 	"sig":     authHeader,
// 	// 	"sig2":    getSigningKeyDemo(),
// 	// })
// }

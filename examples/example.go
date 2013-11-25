package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/zachlatta/go-localbitcoins/localbitcoins"
	"github.com/zachlatta/goauth2-localbitcoins/oauth"
)

var (
	clientId     = flag.String("id", "", "Client ID")
	clientSecret = flag.String("secret", "", "Client Secret")
	scope        = flag.String("scope", "read+write", "OAuth scope")
	redirectURL  = flag.String("redirect_url", "oob", "Redirect URL")
	authURL      = flag.String("auth_url", "https://localbitcoins.com/oauth2/authorize", "Authentication URL")
	tokenURL     = flag.String("token_url", "https://localbitcoins.com/oauth2/access_token", "Token URL")
	requestURL   = flag.String("request_url", "https://localbitcoins.com/api/myself", "API request")
	code         = flag.String("code", "", "Authorization Code")
	cachefile    = flag.String("cache", "cache.json", "Token cache file")
)

const usageMsg = `
To obtain a request token you must specify both -id and -secret.

To obtain a client and secret ID, head over to
https://localbitcoins.com/accounts/api/ and create an application.

Once you have completed the OAuth flow, the credentials should be stored inside
the file specified by -cache and you may run without the -id and -secret flags.
`

func main() {
	flag.Parse()

	// Set up a configuration.
	config := &oauth.Config{
		ClientId:     *clientId,
		ClientSecret: *clientSecret,
		RedirectURL:  *redirectURL,
		Scope:        *scope,
		AuthURL:      *authURL,
		TokenURL:     *tokenURL,
		TokenCache:   oauth.CacheFile(*cachefile),
	}

	// Set up a Transport using the config.
	transport := &oauth.Transport{Config: config}

	// Try to pull the token from the cache; if this fails, we need to get one.
	token, err := config.TokenCache.Token()
	if err != nil {
		if *clientId == "" || *clientSecret == "" {
			flag.Usage()
			fmt.Fprint(os.Stderr, usageMsg)
			os.Exit(2)
		}
		if *code == "" {
			// Get an authorization code from the data provider.
			// ("Please ask the user if I can access this resource.")
			url := config.AuthCodeURL("")
			fmt.Println("Visit this URL to get a code, then run again with -code=YOUR_CODE\n")
			fmt.Println(url)
			return
		}
		// Exchange the authorization code for an access token.
		// ("Here's the code you gave the user, now give me a token!")
		token, err = transport.Exchange(*code)
		if err != nil {
			log.Fatal("Exchange:", err)
		}
		// (The Exchange method will automatically cache the token.)
		fmt.Printf("Token is cached in %v\n", config.TokenCache)
	}

	// Make the actual request using the cached token to authenticate.
	// ("Here's the token, let me in!")
	transport.Token = token

	// Time to create our LocalBitcoins API client.
	client := localbitcoins.NewClient(transport.Client())

	// Fetch the current user account and print out its details.
	acc, _, err := client.Accounts.Get("")
	if err != nil {
		panic(err)
	}
	fmt.Println(acc)

	// Get the account details of the "zrl" account.
	zrl, _, err := client.Accounts.Get("zrl")
	if err != nil {
		panic(err)
	}
	fmt.Println(zrl)
}

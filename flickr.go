package main

import (
	"fmt"
	"os"

	"gopkg.in/masci/flickr.v3"
	"gopkg.in/masci/flickr.v3/test"
)

func flickrAPIClientInit() (string, error) {
	var err error

	client := flickr.NewFlickrClient(config.Flickr.Key, config.Flickr.Secret)

	// Первым делом достанем request token
	tok, err := flickr.GetRequestToken(client)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	// Соорудим авторизационный URL
	url, err := flickr.GetAuthorizeUrl(client, tok)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(3)
	}

	// А теперь ручные действия пользователя.
	// Выдаём ему url, который надо ввести в адресную строку браузера, в браузере сгенерится verifier code, который нам
	// нужен, чтобы вынуть из flicker api oauth token, secret, access token
	var oauthVerifier string
	fmt.Println("Open your browser at this url:", url)
	fmt.Print("Then, insert the code:")
	fmt.Scanln(&oauthVerifier)

	// finally, get the access token
	accessTok, err := flickr.GetAccessToken(client, tok, oauthVerifier)
	fmt.Println("Successfully retrieved OAuth token", accessTok.OAuthToken, accessTok.OAuthTokenSecret)

	// check everything works
	resp, err := test.Login(client)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(resp.Status, resp.User)
	}

	// И где-то тут надо сохранить токен куда-то на диск

	return "", nil
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */

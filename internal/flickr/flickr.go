package flickr

import (
	"aleesa-webapp-go/internal/config"
	"aleesa-webapp-go/internal/pcachedb"
	"errors"
	"strconv"

	"context"
	"fmt"
	"math/rand/v2"
	"net/url"
	"strings"

	"github.com/carlmjohnson/requests"
	"gopkg.in/masci/flickr.v3"
	"gopkg.in/masci/flickr.v3/test"
)

// SearchResult структура, описывающая ответ api на поисковой запрос.
type SearchResult struct {
	Photos struct {
		Page    int64 `json:"page,omitempty"`
		Pages   int64 `json:"pages,omitempty"`
		Perpage int64 `json:"perpage,omitempty"`
		Total   int64 `json:"total,omitempty"`
		Photo   []struct {
			Farm     int64  `json:"farm,omitempty"`
			ID       string `json:"id,omitempty"`
			Isfamily int64  `json:"isfamily,omitempty"`
			Isfriend int64  `json:"isfriend,omitempty"`
			Ispublic int64  `json:"ispublic,omitempty"`
			Owner    string `json:"owner,omitempty"`
			Secret   string `json:"secret,omitempty"`
			Server   string `json:"server,omitempty"`
			Title    string `json:"title,omitempty"`
		} `json:"photo,omitempty"`
	} `json:"photos,omitempty"`
	Stat string `json:"stat,omitempty"`
}

// APIClientInit генерит и сохраняет token и token secret с помощью даденных пользователем key и secret.
// В процессе выдаём юзеру в консоль урл и запрос открыть его в браузере, чтобы получить verifier key, который нам
// понадобится в дальшейшем для генерации рабчочих token и token secret.
func APIClientInit(cfg *config.MyConfig) error {
	var err error

	client := flickr.NewFlickrClient(cfg.Flickr.Key, cfg.Flickr.Secret)

	// Первым делом достанем request token
	tok, err := flickr.GetRequestToken(client)

	if err != nil {
		return fmt.Errorf("unable to get Request Token: %w", err)
	}

	// Соорудим авторизационный URL
	url, err := flickr.GetAuthorizeUrl(client, tok)

	if err != nil {
		return fmt.Errorf("unable to get auth URL: %w", err)
	}

	// А теперь ручные действия пользователя.
	// Выдаём ему url, который надо ввести в адресную строку браузера, в браузере сгенерится verifier code, который нам
	// нужен, чтобы вынуть из flicker api oauth token, secret, access token
	var oauthVerifier string

	fmt.Println("Open your browser at this url:", url)
	fmt.Print("Then, insert the code: ")
	_, _ = fmt.Scanln(&oauthVerifier)

	// Наконец, достаём Access Token и Aceess Token Secret и сохраняем их в базку на диске.
	accessTok, err := flickr.GetAccessToken(client, tok, oauthVerifier)

	if err != nil {
		return fmt.Errorf("unable to make Aceess Token and Access Token Secret, API answers: %w", err)
	}

	fmt.Println("Successfully retrieved OAuth token and secret")

	// Проверим, что всё работает.
	resp, err := test.Login(client)

	if err != nil {
		return fmt.Errorf("unable to perform test login: %w", err)
	}

	fmt.Printf("Test login successful: %s, %s\n", resp.Status, resp.User)

	// И где-то тут надо сохранить токен куда-то на диск.
	if err = pcachedb.SaveKeyWithValue(cfg, "flickr", "token", accessTok.OAuthToken); err != nil {
		return fmt.Errorf("unable to save Aceess Token: %w", err)
	}

	if err = pcachedb.SaveKeyWithValue(cfg, "flickr", "secret", accessTok.OAuthTokenSecret); err != nil {
		return fmt.Errorf("unable to save Aceess Token Secret: %w", err)
	}

	return nil
}

// Populate заполняет кэш ключами.
func Populate(cfg *config.MyConfig) error {
	if err := pcachedb.SaveKeyWithValue(cfg, "flickr", "token", cfg.Flickr.OAuthToken); err != nil {
		return fmt.Errorf("unable to save Aceess Token: %w", err)
	}

	if err := pcachedb.SaveKeyWithValue(cfg, "flickr", "secret", cfg.Flickr.OAuthTokenSecret); err != nil {
		return fmt.Errorf("unable to save Aceess Token Secret: %w", err)
	}

	return nil
}

// SearchByTags делает поисковый запрос по даденным тэгам.
func SearchByTags(cfg *config.MyConfig, tags []string) (string, error) {
	var (
		apiKey      = cfg.Flickr.Key
		apiSecret   = cfg.Flickr.Secret
		oauthToken  = pcachedb.GetValue(cfg, "flickr", "token")
		oauthSecret = pcachedb.GetValue(cfg, "flickr", "secret")
		c           = flickr.NewFlickrClient(apiKey, apiSecret)
		userAgent   = cfg.UserAgents[rand.IntN(len(cfg.UserAgents))]
		ctx         = context.Background()
		tagsString  = strings.Join(tags, ",")
	)

	// Заполняем поля для оаутха.
	c.Init() // c.EndpointUrl
	c.OAuthToken = oauthToken
	c.OAuthTokenSecret = oauthSecret

	// Добавляем параметры. В данном случае это всё, что описано в доке вот тут:
	// https://www.flickr.com/services/api/flickr.photos.search.html
	args := url.Values{
		"format":         {"json"},
		"method":         {"flickr.photos.search"},
		"nojsoncallback": {"1"},
		"content_types":  {"0"},
		"media":          {"photos"},
		"per_page":       {"100"},
		"tags":           {tagsString},
	}

	c.Args = args

	// Подписываем запрос oauth-говном.
	c.OAuthSign()

	// А мы не будем пользоваться flickr.DoGet() по 2-м причинам, во-первых он хочет xml, а это говно, а не язык
	// разметки, а во-вторых, этот самый flickr.DoGet() валится с segfault-ом, просто потому что. Схема не такая, на
	// которую он рассчитан. Но нам повезло, что можно просто вытащить подписанный url и заправить его в свой
	// http-клиент.
	url := c.GetUrl()

	var response SearchResult

	err := requests.
		URL(url).
		UserAgent(userAgent).
		ToJSON(&response).
		Fetch(ctx)

	if response.Stat == "fail" {
		return "", errors.New("search request returns no results and fails")
	}

	// Ну, вот мы и получили ответ на вопрос - а сколько же у нас результатов поиска всего, по заданному критерию.
	// Теперь мы повторим запрос, но на cей раз попросим flickr выдать нам рандомный результат из этой кучи.

	if err != nil {
		return "", fmt.Errorf("unable to get results for tags %s from flickr api: %w", tagsString, err)
	}

	randomPage := rand.Int64N(response.Photos.Pages)
	c.Args.Set("page", strconv.FormatInt(randomPage, 10))

	// Подписываем запрос oauth-говном.
	c.OAuthSign()

	url = c.GetUrl()

	err = requests.
		URL(url).
		UserAgent(userAgent).
		ToJSON(&response).
		Fetch(ctx)

	if err != nil {
		return "", fmt.Errorf("unable to get results for tags %s from flickr api: %w", tagsString, err)
	}

	if response.Stat == "fail" {
		return "", errors.New("search request returns no results and fails")
	}

	// Частенько выборка не очень уникальна, много дублей. Чтобы почистить выборку, неплохо бы выбирать только те фотки,
	// у которых response.Photos.Photo[exactResult].Title уникальный.

	var u SearchResult

	for i := 0; i < len(response.Photos.Photo); i++ {
		skip := false

		for j := 0; j < len(u.Photos.Photo); j++ {
			if response.Photos.Photo[i].Title == u.Photos.Photo[j].Title {
				skip = true

				break
			}
		}

		if skip {
			continue
		}

		u.Photos.Photo = append(u.Photos.Photo, response.Photos.Photo[i])
	}

	exactResult := rand.IntN(len(u.Photos.Photo))

	result := fmt.Sprintf(
		"https://live.staticflickr.com/%s/%s_%s_z.jpg",
		u.Photos.Photo[exactResult].Server,
		u.Photos.Photo[exactResult].ID,
		u.Photos.Photo[exactResult].Secret,
	)

	return result, err
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */

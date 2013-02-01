package infrastructure

import (
	"encoding/base32"
	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"net/http"
	"net/url"
	"strings"
)

var redisUrl *url.URL

type RedisStoreHandler struct {
	keyPairs [][]byte
	pool     *redis.Pool
}

func (h *RedisStoreHandler) GetSessionStore() sessions.Store {
	return &RedisStore{
		Codecs: securecookie.CodecsFromPairs(h.keyPairs...),
		Options: &sessions.Options{
			Path:   "/",
			MaxAge: 86400 * 30,
		},
		storeHandler: h,
	}
}

func (h *RedisStoreHandler) GetRedisConnection() redis.Conn {
	return h.pool.Get()
}

type RedisStore struct {
	storeHandler *RedisStoreHandler
	Codecs       []securecookie.Codec
	Options      *sessions.Options
}

func (s *RedisStore) Get(r *http.Request, name string) (*sessions.Session, error) {
	return sessions.GetRegistry(r).Get(s, name)
}

func (s *RedisStore) New(r *http.Request, name string) (*sessions.Session, error) {
	session := sessions.NewSession(s, name)
	session.Options = &(*s.Options)
	session.IsNew = true
	var err error
	if c, errCookie := r.Cookie(name); errCookie == nil {
		// TEMP: Check if previous CookieStore
		err = securecookie.DecodeMulti(name, c.Value, &session.Values, s.Codecs...)
		if err == nil {
			err = s.save(session)
			if err == nil {
				session.IsNew = false
			}
		} else {
			err = securecookie.DecodeMulti(name, c.Value, &session.ID, s.Codecs...)
			if err == nil {
				err = s.load(session)
				if err == nil {
					session.IsNew = false
				}
			}
		}
	}
	return session, err
}

func (s *RedisStore) Save(r *http.Request, w http.ResponseWriter, session *sessions.Session) error {
	if session.ID == "" {
		// Because the ID is used in the filename, encode it to
		// use alphanumeric characters only.
		session.ID = strings.TrimRight(base32.StdEncoding.EncodeToString(securecookie.GenerateRandomKey(32)), "=")
	}
	if err := s.save(session); err != nil {
		return err
	}
	encoded, err := securecookie.EncodeMulti(session.Name(), session.ID, s.Codecs...)
	if err != nil {
		return err
	}
	http.SetCookie(w, sessions.NewCookie(session.Name(), encoded, session.Options))
	return nil
}

func (s *RedisStore) save(session *sessions.Session) error {
	encoded, err := securecookie.EncodeMulti(session.Name(), session.Values, s.Codecs...)
	if err != nil {
		return err
	}
	c := s.storeHandler.GetRedisConnection()
	defer c.Close()
	if session.ID == "" {
		// Because the ID is used in the filename, encode it to
		// use alphanumeric characters only.
		session.ID = strings.TrimRight(base32.StdEncoding.EncodeToString(securecookie.GenerateRandomKey(32)), "=")
	}
	c.Send("SET", "morioka_sess_"+session.ID, encoded)
	if err = c.Flush(); err != nil {
		return err
	}
	if _, err = c.Receive(); err != nil {
		return err
	}
	c.Send("EXPIRE", "morioka_sess_"+session.ID, 86400)
	if err = c.Flush(); err != nil {
		return err
	}
	if _, err = c.Receive(); err != nil {
		return err
	}
	return nil
}

func (s *RedisStore) load(session *sessions.Session) error {
	c := s.storeHandler.GetRedisConnection()
	defer c.Close()
	c.Send("GET", "morioka_sess_"+session.ID)
	c.Flush()
	data, err := redis.String(c.Receive())
	if err != nil {
		return err
	}
	if err = securecookie.DecodeMulti(session.Name(), data, &session.Values, s.Codecs...); err != nil {
		return err
	}
	return nil
}

func NewRedisStoreHandler(inputUrl string, keyPairs ...[]byte) *RedisStoreHandler {
	redisUrl, _ = url.Parse(inputUrl)
	return &RedisStoreHandler{
		pool: redis.NewPool(func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", redisUrl.Host)
			if err != nil {
				return nil, err
			}
			password, has := redisUrl.User.Password()
			if !has {
				return c, nil
			}
			if _, err := c.Do("AUTH", password); err != nil {
				c.Close()
				return nil, err
			}
			return c, err
		}, 3),
		keyPairs: keyPairs,
	}
}

package secrets

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	CacheTimeOut = "CACHE_EXTENSION_TTL"
)

var (
	ExtensionName = filepath.Base(os.Args[0])
	PrintPrefix   = fmt.Sprintf("[%s] ", ExtensionName)
	secretCache   = make(map[string]Secret)
)

type Secret struct {
	CacheData CacheData
}

type CacheData struct {
	Data        string
	CacheExpiry time.Time
}

var mutex = &sync.Mutex{}

func IsExpired(cacheExpiry time.Time) bool {
	return cacheExpiry.Before(time.Now())
}

func GetCacheExpiry() time.Time {
	timeOut := os.Getenv(CacheTimeOut)
	if timeOut == "" {
		timeOut = "60m"
	}

	timeOutInMinutes, err := time.ParseDuration(timeOut)
	if err != nil {
		panic("Error while converting CACHE_EXTENSION_TTL env variable " + timeOut)
	}

	return time.Now().Add(timeOutInMinutes)
}

func GetSecretCache(name string, refresh string) string {
	secret := secretCache[name]

	if IsExpired(secret.CacheData.CacheExpiry) || refresh == "1" {
		mutex.Lock()
		secretCache[name] = Secret{
			CacheData: CacheData{
				Data:        GetSecret(name),
				CacheExpiry: GetCacheExpiry(),
			},
		}
		mutex.Unlock()
	}

	return secretCache[name].CacheData.Data
}

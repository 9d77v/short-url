package app

import (
	"context"
	"log"
	"time"

	"github.com/9d77v/go-pkg/cache/redis"
)

//ShortURL 短地址
type ShortURL struct {
	ID        uint   `gorm:"primarykey"`
	URL       string `gorm:"size:500;NOT NULL;comment:长地址"`
	ShortCode string `gorm:"size:10;NOT NULL;comment:短码"`
	Deadline  time.Time
	CreatedAt time.Time
}

const (
	offset           = 9999
	URLPrefix        = "SURL:"
	URLLockPrefix    = "SURL_LOCK:"
	ShortId          = "SURL_SHORT_ID"
	CODEPrefix       = "SURL_CODE:"
	CODEVistedPrefix = "SURL_CODE_VISITED:"
)

var (
	lenMap = map[int]int{}
	arr    = []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'}
)

//GetCode 获取地址短码
func GetURL(code string) (string, error) {
	return redis.GetClient().Get(context.Background(), CODEPrefix+code).Result()
}

//ConvertURL 转换地址
func ConvertURL(url string, h int) string {
	ctx := context.Background()
	shortCode, _ := redis.GetClient().Get(context.Background(), URLPrefix+url).Result()
	if shortCode != "" {
		return shortCode
	}
	c := redis.GetClient()
	c.DLock(ctx, URLLockPrefix+url, 5*time.Second, func() {
		id, _ := c.Incr(ctx, ShortId).Result()
		shortCode = ConvertIntToStr(int(id + offset))
		now := time.Now()
		if h <= 0 {
			err := c.Set(ctx, URLPrefix+url, shortCode, -1).Err()
			if err != nil {
				log.Println("redis set error:", err)
				return
			}
			err = c.Set(ctx, CODEPrefix+shortCode, url, -1).Err()
			if err != nil {
				log.Println("redis set error:", err)
				return
			}
		} else {
			err := c.SetEX(ctx, URLPrefix+url, shortCode, time.Duration(h*int(time.Hour))).Err()
			if err != nil {
				log.Println("redis set error:", err)
				return
			}
			err = c.SetEX(ctx, CODEPrefix+shortCode, url, time.Duration(h*int(time.Hour))).Err()
			if err != nil {
				log.Println("redis set error:", err)
				return
			}
		}
		err := GetDB().Create(&ShortURL{
			ID:        uint(id),
			URL:       url,
			ShortCode: shortCode,
			CreatedAt: now,
			Deadline:  now.Add(time.Duration(h * int(time.Hour))),
		}).Error
		if err != nil {
			log.Println("save shor url failed:", err)
		}
	})
	return shortCode
}

func ConvertIntToStr(num int) string {
	ans, n := []byte{}, 62
	for num != 0 {
		t := num % n
		num /= n
		ans = append(ans, arr[t])
	}
	m := len(ans)
	for i := 0; i < m>>1; i++ {
		ans[i], ans[m-1-i] = ans[m-1-i], ans[i]
	}
	return string(ans)
}

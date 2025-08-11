# 短網址系統完整規格書

## 1. 專案概述

此系統為一個具備角色管理（RBAC）、短網址生成與管理、點擊分析（國家/系統/瀏覽器/時間統計）的 Web 應用，採用 **乾淨架構 (Clean Architecture)** 開發，後端使用 Go + Gin，前端使用 Astro + TailwindCSS + daisyUI，資料庫採 SQLite + sqlc，並支援 Telegram Bot 操作。

---

## 2. 技術棧

### 後端（Go）

- `github.com/gin-gonic/gin`
- `github.com/golang-jwt/jwt/v5`
- `golang.org/x/crypto/bcrypt`
- `github.com/mattn/go-sqlite3`
- `github.com/kyleconroy/sqlc`
- `github.com/joho/godotenv`（可選）
- `github.com/golang-migrate/migrate/v4`（可選）
- `github.com/go-playground/validator/v10`
- `github.com/gin-contrib/cors`
- `github.com/didip/tollbooth/v6`
- `github.com/go-telegram-bot-api/telegram-bot-api/v5`
- **GeoIP**：`github.com/oschwald/geoip2-golang`
- **User-Agent 解析**（二擇一）：
    - `github.com/mssola/user_agent`
    - `github.com/ua-parser/uap-go/uaparser`
- **佇列/工作池（可選）**：`github.com/alitto/pond`

### 前端（Astro）

- `astro`
- `tailwindcss`
- `daisyui`
- `chart.js`

---

## 3. 資料夾與檔案規劃

```
cmd/
  server/main.go
config/
  config.go
internal/
  domain/
    shorturl.go
    user.go
    permission.go
    click.go
  application/
    url_usecase.go
    user_usecase.go
    analytics_usecase.go
infrastructure/
  repository/
    shorturl_repository.go
    user_repository.go
    click_repository.go
  geo/geoip_service.go
  ua/ua_parser.go
  worker/click_worker.go
  telegram/bot_handler.go
  web/http_handler.go
  web/middleware.go
presentation/
  gin/router.go
  gin/handlers.go
  gin/admin_handlers.go
  telegram/commands.go
sqlc/
  queries/
    users.sql
    short_urls.sql
    url_clicks.sql
    telegram_auth_tokens.sql
  generated.go
web/   # Astro 專案
  src/
  dist/
```

---

## 4. 主要檔案內容規劃

### config/config.go

```go
type Config struct {
    DBPath       string
    BotToken     string
    JWTSecret    string
    ServerPort   string
    Environment  string
}
func LoadConfig() (*Config, error)
```

### internal/domain/shorturl.go

```go
type ShortURL struct {
    ID          int64
    ShortPath   string
    OriginalURL string
    UserID      int64
    CreatedAt   time.Time
}
type ShortURLRepository interface {
    Create(ctx context.Context, shortURL ShortURL) (int64, error)
    GetByPath(ctx context.Context, path string) (*ShortURL, error)
    Delete(ctx context.Context, id int64) error
    ListByUser(ctx context.Context, userID int64) ([]ShortURL, error)
}
```

### internal/domain/user.go

```go
type User struct {
    ID             int64
    Username       string
    PasswordHash   string
    Permissions    int
    TelegramID     int64
    CreatedAt      time.Time
}
type UserRepository interface {
    Create(ctx context.Context, user User) (int64, error)
    GetByUsername(ctx context.Context, username string) (*User, error)
    GetByID(ctx context.Context, id int64) (*User, error)
    UpdateTelegramID(ctx context.Context, tgid int64) error
    UpdatePermissions(ctx context.Context, id int64, permissions int) error
}
```

### internal/domain/permission.go

```go
package domain

// Permission is a bitmask type for user permissions.
// Creating random URLs is a baseline action available to guests and all authenticated users,
// and thus is not governed by a specific permission bit.
type Permission int

const (
	PermCreatePrefix Permission = 1 << iota // 1
	PermCreateAny    Permission             // 2
	PermDeleteOwn    Permission             // 4
	PermDeleteAny    Permission             // 8
	PermViewOwnStats Permission             // 16
	PermViewAnyStats Permission             // 32
	PermUserManage   Permission             // 64
)

// Permission sets (pre-configured permission bundles)
const (
	RoleGuest      Permission = 0
	RoleRegular    Permission = PermCreatePrefix | PermDeleteOwn | PermViewOwnStats
	RolePrivileged Permission = RoleRegular | PermCreateAny
	RoleEditor     Permission = RolePrivileged | PermDeleteAny | PermViewAnyStats
	RoleAdmin      Permission = RoleEditor | PermUserManage
)

// Has checks if the user's permissions (p) include the required permission (required).
func (p Permission) Has(required Permission) bool {
	// If required is 0, it's an invalid permission to check.
	if required == 0 {
		return false
	}
	return (p & required) == required
}

// Add grants a new permission.
func (p Permission) Add(perm Permission) Permission {
	return p | perm
}

// Remove revokes a permission.
func (p Permission) Remove(perm Permission) Permission {
	return p &^ perm
}
```

### internal/domain/click.go

```go
type URLClick struct {
    ID          int64
    ShortURLID  int64
    ClickedAt   time.Time
    CountryCode *string
    OSName      *string
    BrowserName *string
    DeviceType  *string
    UserAgent   *string
}
type ClickRepository interface {
    Insert(ctx context.Context, c URLClick) (int64, error)
    CountByShortURL(ctx context.Context, shortURLID int64) (int64, error)
    AggregateByTimeRange(ctx context.Context, shortURLID int64, from, to time.Time, bucket string) ([]TimeBucketCount, error)
    AggregateByCountry(ctx context.Context, shortURLID int64, from, to time.Time) ([]KeyCount, error)
    AggregateByOS(ctx context.Context, shortURLID int64, from, to time.Time) ([]KeyCount, error)
    AggregateByBrowser(ctx context.Context, shortURLID int64, from, to time.Time) ([]KeyCount, error)
}
```

### internal/application/url_usecase.go

```go
type URLUseCase struct {
    repo      ShortURLRepository
    clickRepo ClickRepository
    geo       GeoService
    ua        UAService
    clickSink ClickSink
}
func (uc *URLUseCase) CreateShortURL(ctx context.Context, user User, originalURL, customPath string) (*ShortURL, error)
func (uc *URLUseCase) DeleteShortURL(ctx context.Context, user User, shortPath string) error
func (uc *URLUseCase) ListByUser(ctx context.Context, user User) ([]ShortURL, error)
func (uc *URLUseCase) GetStats(ctx context.Context, shortPath string) (*URLStats, error)
func (uc *URLUseCase) RecordClickAsync(ctx context.Context, short *ShortURL, meta ClickMeta) error
```

### internal/application/user_usecase.go

```go
type UserUseCase struct {
    repo UserRepository
}
func (uc *UserUseCase) Register(ctx context.Context, username, password string) (*User, error)
func (uc *UserUseCase) Login(ctx context.Context, username, password string) (string, error)
func (uc *UserUseCase) GetMe(ctx context.Context, userID int64) (*User, error)
func (uc *UserUseCase) LinkTelegram(ctx context.Context, userID int64, chatID int64) error
func (uc *UserUseCase) UpdateUserRole(ctx context.Context, operator User, targetUserID int64, newRole string) error
```

### internal/application/analytics_usecase.go

```go
type AnalyticsUseCase struct {
    clickRepo ClickRepository
    urlRepo   ShortURLRepository
}
func (a *AnalyticsUseCase) GetOverview(ctx context.Context, shortPath string, from, to time.Time) (*URLStats, error)
func (a *AnalyticsUseCase) GetTimeSeries(ctx context.Context, shortPath string, from, to time.Time, bucket string) ([]TimeBucketCount, error)
func (a *AnalyticsUseCase) GetByCountry(ctx context.Context, shortPath string, from, to time.Time) ([]KeyCount, error)
func (a *AnalyticsUseCase) GetByOS(ctx context.Context, shortPath string, from, to time.Time) ([]KeyCount, error)
func (a *AnalyticsUseCase) GetByBrowser(ctx context.Context, shortPath string, from, to time.Time) ([]KeyCount, error)
```

---

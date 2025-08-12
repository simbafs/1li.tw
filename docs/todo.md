# 開發實作順序 (Implementation Order)

本文件根據系統規格書，將開發任務切分為數個主要模組，並以 checklist 形式排列建議的實作順序。

## 階段一：核心架構與資料庫基礎 (Phase 1: Core Architecture & DB Foundation)

此階段專注於建立專案骨架、定義核心業務邏輯，並完成資料庫的初始化設定。

- [x] **專案結構初始化**
    - [x] 建立 Go Modules (`go mod init`)
    - [x] 根據規格書建立 `cmd`, `config`, `domain`, `application`, `infrastructure`, `presentation` 等資料夾結構
- [x] **Domain (實體層) 定義**
    - [x] `domain/permission.go`: 定義權限位元與角色
    - [x] `domain/user.go`: 定義 `User` 實體與 `UserRepository` 介面
    - [x] `domain/shorturl.go`: 定義 `ShortURL` 實體與 `ShortURLRepository` 介面
    - [x] `domain/click.go`: 定義 `URLClick` 實體與 `ClickRepository` 介面
- [x] **資料庫 Schema 與 sqlc**
    - [x] `sql/schema.sql`: 根據規格書撰寫所有資料表的 SQL `CREATE TABLE` 語法
    - [x] `sql/queries/*.sql`: 撰寫 `users`, `short_urls`, `url_clicks`, `telegram_auth_tokens` 的 CRUD SQL 查詢語句
    - [x] `sqlc.yaml`: 設定 sqlc，並執行 `sqlc generate` 產生類型安全的資料庫操作程式碼

## 階段二：基礎設施與應用邏輯 (Phase 2: Infrastructure & Application Logic)

此階段專注於實作與外部服務（如資料庫）溝通的具體邏輯，並建立核心應用程式的 Use Cases。

- [ ] 根據 @docs/spec.md 撰寫單元測試：**Infrastructure (基礎設施層) 實作**
- [ ] **Infrastructure (基礎設施層) 實作**
    - [ ] 根據 @docs/spec.md 撰寫單元測試：`infrastructure/repository`: 實作 `UserRepository`, `ShortURLRepository`, `ClickRepository` 等介面，串接 sqlc 產生的程式碼
    - [ ] `infrastructure/repository`: 實作 `UserRepository`, `ShortURLRepository`, `ClickRepository` 等介面，串接 sqlc 產生的程式碼
    - [ ] 根據 @docs/spec.md 撰寫單元測試：`infrastructure/worker`: 建立非同步處理點擊事件的背景 Worker (`ClickSink`)
    - [ ] `infrastructure/worker`: 建立非同步處理點擊事件的背景 Worker (`ClickSink`)
    - [ ] 根據 @docs/spec.md 撰寫單元測試：`infrastructure/geo`: 實作 GeoIP 服務，用於從 IP 解析國家代碼
    - [ ] `infrastructure/geo`: 實作 GeoIP 服務，用於從 IP 解析國家代碼
    - [ ] 根據 @docs/spec.md 撰寫單元測試：`infrastructure/ua`: 實作 User-Agent 解析服務
    - [ ] `infrastructure/ua`: 實作 User-Agent 解析服務
- [ ] 根據 @docs/spec.md 撰寫單元測試：**Application (應用層) Use Cases**
- [ ] **Application (應用層) Use Cases**
    - [ ] 根據 @docs/spec.md 撰寫單元測試：`application/user_usecase.go`: 實作使用者註冊、登入、綁定 Telegram 等邏輯
    - [ ] `application/user_usecase.go`: 實作使用者註冊、登入、綁定 Telegram 等邏輯
    - [ ] 根據 @docs/spec.md 撰寫單元測試：`application/url_usecase.go`: 實作短網址的建立、刪除、列表查詢等邏輯
    - [ ] `application/url_usecase.go`: 實作短網址的建立、刪除、列表查詢等邏輯
    - [ ] 根據 @docs/spec.md 撰寫單元測試：`application/analytics_usecase.go`: 實作點擊數據的統計與分析邏輯
    - [ ] `application/analytics_usecase.go`: 實作點擊數據的統計與分析邏輯

## 階段三：API 端點與核心功能 (Phase 3: API Endpoints & Core Features)

此階段將透過 Gin 框架將後端邏輯暴露為 HTTP API，並完成系統最核心的重定向功能。

- [ ] 根據 @docs/spec.md 撰寫單元測試：**Presentation (表現層) - Gin Web Server**
- [ ] **Presentation (表現層) - Gin Web Server**
    - [ ] 根據 @docs/spec.md 撰寫單元測試：`config/config.go`: 實作設定檔載入
    - [ ] `config/config.go`: 實作設定檔載入
    - [ ] 根據 @docs/spec.md 撰寫單元測試：`presentation/gin/router.go`: 設定 Gin 路由
    - [ ] `presentation/gin/router.go`: 設定 Gin 路由
    - [ ] 根據 @docs/spec.md 撰寫單元測試：**使用者認證 API**
    - [ ] **使用者認證 API**
        - [ ] 根據 @docs/spec.md 撰寫單元測試：`POST /api/auth/register`
        - [ ] `POST /api/auth/register`
        - [ ] 根據 @docs/spec.md 撰寫單元測試：`POST /api/auth/login` (包含 JWT 發行與 HttpOnly Cookie 設定)
        - [ ] `POST /api/auth/login` (包含 JWT 發行與 HttpOnly Cookie 設定)
        - [ ] 根據 @docs/spec.md 撰寫單元測試：實作 JWT 驗證的 Middleware
        - [ ] 實作 JWT 驗證的 Middleware
    - [ ] 根據 @docs/spec.md 撰寫單元測試：**核心重定向功能**
    - [ ] **核心重定向功能**
        - [ ] 根據 @docs/spec.md 撰寫單元測試：`GET /:short_path`
        - [ ] `GET /:short_path`
        - [ ] 根據 @docs/spec.md 撰寫單元測試：`GET /@:username/:custom_path`
        - [ ] `GET /@:username/:custom_path`
        - [ ] 根據 @docs/spec.md 撰寫單元測試：整合非同步點擊紀錄 Worker
        - [ ] 整合非同步點擊紀錄 Worker
    - [ ] 根據 @docs/spec.md 撰寫單元測試：**短網址管理 API**
    - [ ] **短網址管理 API**
        - [ ] 根據 @docs/spec.md 撰寫單元測試：`POST /api/urls`
        - [ ] `POST /api/urls`
        - [ ] 根據 @docs/spec.md 撰寫單元測試：`GET /api/urls`
        - [ ] `GET /api/urls`
        - [ ] 根據 @docs/spec.md 撰寫單元測試：`DELETE /api/urls/:id`
        - [ ] `DELETE /api/urls/:id`
    - [ ] 根據 @docs/spec.md 撰寫單元測試：**統計數據 API**
    - [ ] **統計數據 API**
        - [ ] 根據 @docs/spec.md 撰寫單元測試：`GET /api/urls/:id/stats`
        - [ ] `GET /api/urls/:id/stats`
    - [ ] 根據 @docs/spec.md 撰寫單元測試：**管理員 API**
    - [ ] **管理員 API**
        - [ ] 根據 @docs/spec.md 撰寫單元測試：`GET /api/admin/urls`
        - [ ] `GET /api/admin/urls`
        - [ ] 根據 @docs/spec.md 撰寫單元測試：`DELETE /api/admin/urls/:id`
        - [ ] `DELETE /api/admin/urls/:id`

## 階段四：介面整合 (Phase 4: Interface Integration)

此階段專注於完成與使用者互動的介面，包括 Telegram Bot 和 Web 前端。

- [ ] 根據 @docs/spec.md 撰寫單元測試：**Presentation (表現層) - Telegram Bot**
- [ ] **Presentation (表現層) - Telegram Bot**
    - [ ] 根據 @docs/spec.md 撰寫單元測試：`infrastructure/telegram/bot_handler.go`: 初始化 Bot 並設定 Webhook 或長輪詢
    - [ ] `infrastructure/telegram/bot_handler.go`: 初始化 Bot 並設定 Webhook 或長輪詢
    - [ ] 根據 @docs/spec.md 撰寫單元測試：**Web 授權流程**
    - [ ] **Web 授權流程**
        - [ ] 根據 @docs/spec.md 撰寫單元測試：`GET /auth/telegram`: 實作授權頁面
        - [ ] `GET /auth/telegram`: 實作授權頁面
        - [ ] 根據 @docs/spec.md 撰寫單元測試：`POST /api/auth/telegram/link`: 實作完成綁定的 API
        - [ ] `POST /api/auth/telegram/link`: 實作完成綁定的 API
    - [ ] 根據 @docs/spec.md 撰寫單元測試：`presentation/telegram/commands.go`: 串接 Use Cases，實作所有 Bot 指令
    - [ ] `presentation/telegram/commands.go`: 串接 Use Cases，實作所有 Bot 指令
        - [ ] 根據 @docs/spec.md 撰寫單元測試：`/auth`
        - [ ] `/auth`
        - [ ] 根據 @docs/spec.md 撰寫單元測試：`/new`, `/delete`, `/list`, `/stats`
        - [ ] `/new`, `/delete`, `/list`, `/stats`
        - [ ] 根據 @docs/spec.md 撰寫單元測試：`/whoami`, `/start`
        - [ ] `/whoami`, `/start`
- [ ] 根據 @docs/spec.md 撰寫單元測試：**Presentation (表現層) - Astro Web 前端**
- [ ] **Presentation (表現層) - Astro Web 前端**
    - [ ] 根據 @docs/spec.md 撰寫單元測試：初始化 Astro 專案 (`web/`)
    - [ ] 初始化 Astro 專案 (`web/`)
    - [ ] 根據 @docs/spec.md 撰寫單元測試：**訪客頁面**
    - [ ] **訪客頁面**
        - [ ] 根據 @docs/spec.md 撰寫單元測試：首頁 (快速建立短網址)
        - [ ] 首頁 (快速建立短網址)
        - [ ] 根據 @docs/spec.md 撰寫單元測試：登入/註冊頁面
        - [ ] 登入/註冊頁面
    - [ ] 根據 @docs/spec.md 撰寫單元測試：**使用者 Dashboard**
    - [ ] **使用者 Dashboard**
        - [ ] 根據 @docs/spec.md 撰寫單元測試：短網址列表 (含總點擊)
        - [ ] 短網址列表 (含總點擊)
        - [ ] 根據 @docs/spec.md 撰寫單元測試：新增/刪除短網址功能
        - [ ] 新增/刪除短網址功能
    - [ ] 根據 @docs/spec.md 撰寫單元測試：**統計頁面**
    - [ ] **統計頁面**
        - [ ] 根據 @docs/spec.md 撰寫單元測試：使用 `chart.js` 將統計 API 數據視覺化
        - [ ] 使用 `chart.js` 將統計 API 數據視覺化
    - [ ] 根據 @docs/spec.md 撰寫單元測試：**管理後台**
    - [ ] **管理後台**
        - [ ] 根據 @docs/spec.md 撰寫單元測試：全系統短網址列表
        - [ ] 全系統短網址列表
        - [ ] 根據 @docs/spec.md 撰寫單元測試：搜尋與刪除功能
        - [ ] 搜尋與刪除功能

## 階段五：安全、測試與部署 (Phase 5: Security, Testing & Deployment)

此階段為最後的收尾工作，確保系統的穩定性、安全性並準備上線。

- [ ] 根據 @docs/spec.md 撰寫單元測試：**安全性強化**
- [ ] **安全性強化**
    - [ ] 根據 @docs/spec.md 撰寫單元測試：實作 API 速率限制 (Rate Limiting)
    - [ ] 實作 API 速率限制 (Rate Limiting)
    - [ ] 根據 @docs/spec.md 撰寫單元測試：加上必要的 HTTP 安全標頭 (CSP, HSTS 等)
    - [ ] 加上必要的 HTTP 安全標頭 (CSP, HSTS 等)
    - [ ] 根據 @docs/spec.md 撰寫單元測試：確保資料庫檔案位於非 Web 可存取目錄
    - [ ] 確保資料庫檔案位於非 Web 可存取目錄
- [ ] 根據 @docs/spec.md 撰寫單元測試：**測試**
- [ ] **測試**
    - [ ] 根據 @docs/spec.md 撰寫單元測試：為各層級的關鍵邏輯撰寫單元測試 (Unit Tests)
    - [ ] 為各層級的關鍵邏輯撰寫單元測試 (Unit Tests)
    - [ ] 根據 @docs/spec.md 撰寫單元測試：撰寫 API 整合測試 (Integration Tests)
    - [ ] 撰寫 API 整合測試 (Integration Tests)
- [ ] 根據 @docs/spec.md 撰寫單元測試：**部署準備**
- [ ] **部署準備**
    - [ ] 根據 @docs/spec.md 撰寫單元測試：`cmd/server/main.go`: 整合 `embed` 套件以嵌入前端靜態檔案
    - [ ] `cmd/server/main.go`: 整合 `embed` 套件以嵌入前端靜態檔案
    - [ ] 根據 @docs/spec.md 撰寫單元測試：撰寫 Dockerfile (可選)
    - [ ] 撰寫 Dockerfile (可選)
    - [ ] 根據 @docs/spec.md 撰寫單元測試：撰寫部署文件與腳本
    - [ ] 撰寫部署文件與腳本

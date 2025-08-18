# Telegram Bot 整合計畫

## 1. 總覽

本文件旨在規劃將 Telegram Bot 整合到 1li.tw 服務中。使用者將能夠透過 Telegram Bot 進行帳號綁定，並執行一些基本操作，例如新增和查詢自己的短網址。

此整合將遵循現有的 Clean Architecture 架構，將 Bot 視為一個新的 `Presentation` 層，並重複使用現有的 `Application` 層 Use Cases，以確保程式碼的模組化和可維護性。

## 2. 使用者認證流程

為了確保安全性並將 Telegram 使用者與 1li.tw 的帳號對應，我們將採用基於驗證連結的流程。

1.  **使用者發起綁定**
    *   使用者在 Telegram 中找到我們的 Bot，並發送 `/start` 指令。
    *   Bot 回應歡迎訊息，並生成一個一次性的、有時效性的認證 Token。
    *   Bot 建立一個驗證連結 (例如 `https://1li.tw/telegram/bind?token=[UNIQUE_AUTH_TOKEN]`) 並回傳給使用者。
    *   回覆訊息範例：「歡迎使用 1li.tw！請點擊以下連結以綁定您的帳號。此連結將在 10 分鐘後失效：[驗證連結]」

2.  **使用者點擊連結並完成綁定**
    *   使用者點擊 Telegram 中的驗證連結，在其瀏覽器中開啟網頁。
    *   網頁會要求使用者登入 1li.tw 帳號（如果尚未登入）。
    *   登入後，網頁會自動將連結中的 Token 送到後端進行驗證。

3.  **後端處理綁定**
    *   後端 API (例如 `POST /api/telegram/bind`) 接收到 Token。
    *   後端驗證 Token 的有效性，並從中找到對應的 Telegram Chat ID。
    *   驗證成功後，後端會將該 Telegram Chat ID 與當前登入的 1li.tw 使用者 ID 進行綁定，並更新 `users` 資料表。
    *   綁定成功後，Bot 會主動發送一條訊息給該使用者，通知他們已成功綁定。
    *   訊息範例：「您的 1li.tw 帳號已成功綁定！現在您可以使用 `/help` 來查看所有可用指令。」

## 3. 功能指令

初期我們將支援以下核心功能：

*   **/start**
    *   **說明**: 開始與 Bot 互動，並取得帳號綁定連結。
    *   **用法**: `/start`

*   **/help**
    *   **說明**: 顯示所有可用的指令及其用法。
    *   **用法**: `/help`

*   **/new `<url>` `[custom_path]`**
    *   **說明**: 建立一個新的短網址。
    *   **參數**:
        *   `<url>`: (必要) 您想要縮短的原始網址。
        *   `[custom_path]`: (可選) 自訂的短網址路徑。如果留空，將會隨機生成。
    *   **成功回覆**: 「✅ 短網址已建立：`https://1li.tw/[path]`」
    *   **失敗回覆**: 「❌ 建立失敗：[錯誤原因]」
    *   **用法範例**:
        *   `/new https://github.com/simba-fs/1li.tw`
        *   `/new https://google.com my-google`

*   **/list `[page]`**
    *   **說明**: 列出使用者建立的所有短網址。考慮到訊息長度限制，此功能將支援分頁。
    *   **參數**:
        *   `[page]`: (可選) 頁碼，預設為 1。
    *   **回覆**:
        ```
        您的短網址 (第 1/5 頁):
        1. /path1 -> https://example.com/...
        2. /path2 -> https://google.com/...
        ...
        10. /path10 -> https://yahoo.com/...
        ```
    *   **用法範例**: `/list` 或 `/list 2`

*   **/revoke**
    *   **說明**: 解除 Telegram 帳號與 1li.tw 帳號的綁定。Bot 將會刪除對應的 Chat ID 並通知使用者。
    *   **用法**: `/revoke`

## 4. 技術架構規劃

1.  **設定 (Configuration)**
    *   計畫在 `.env` 檔案中新增 `TELEGRAM_BOT_TOKEN` 來存放 Bot 的 API Token。
    *   計畫更新 `config/config.go` 來讀取此設定。

2.  **資料庫 (Database)**
    *   計畫修改 `sql/schema.sql` 中的 `users` 資料表，新增一個 `telegram_chat_id` 欄位 (INTEGER, NULLABLE, UNIQUE)，用來儲存綁定後的 Telegram Chat ID。
    *   計畫新增一個 `telegram_auth_tokens` 資料表，用來儲存 `/start` 產生的暫時性 Token，包含 `token`, `chat_id`, `expires_at` 等欄位。

3.  **後端邏輯 (Backend Logic)**
    *   **Telegram Bot 函式庫**: 我們將使用 `github.com/PaulSonOfLars/gotgbot/v2` 作為與 Telegram Bot API 互動的主要函式庫。
    *   **主要進入點**: 計畫在 `main.go` 中，初始化並啟動一個新的 Goroutine 來處理 Telegram Bot 的 long polling。
    *   **Bot 處理器**: 計畫在 `infrastructure/telegram/bot_handler.go` 中建立一個新的 Bot 處理器。
    *   **指令解析**: 處理器會接收來自 Telegram 的更新，解析訊息中的指令 (例如 `/new`)。
    *   **重複使用 Use Cases**: 處理器將會呼叫 `application` 層中已經存在的 Use Cases (例如 `UrlUsecase.CreateShortURL`, `UrlUsecase.GetUserURLs`) 來執行核心業務邏輯。這確保了業務邏輯的一致性，避免了重複開發。
    *   **API 端點**: 計畫在 `presentation/gin/router.go` 中新增一個 API 端點 (例如 `POST /api/telegram/bind`) 來處理來自網頁的 Token 綁定請求。

## 5. 未來擴充性

此架構預留了良好的擴充空間：

*   **新增指令**: 可以輕易地在 `bot_handler.go` 中新增指令的 case，並對應到新的或現有的 Use Case。
*   **進階功能**: 未來可以擴充更多功能，例如：
    *   `/stats <path>`: 查詢特定短網址的點擊統計。
    *   `/delete <path>`: 刪除一個短網址。
    *   `/set <path> [new_target]` : 修改短網址的目標。
*   **互動式對話**: 對於較複雜的操作，可以利用 Telegram Bot 的 conversation 功能，進行多步驟的互動式問答。
*   **管理員指令**: 可以為管理員角色新增特有的指令，例如 `/admin_stats` 或 `/admin_ban_user`。

## 6. 開發階段

1.  **階段一：基礎建設**
    *   修改資料庫 Schema。
    *   完成設定的讀取。
    *   建立基本的 Bot Goroutine 和指令路由器。

2.  **階段二：認證流程**
    *   實作 `/start` 指令與 Token/連結 生成邏輯。
    *   在前端和後端新增處理綁定連結的功能。
    *   完成完整的帳號綁定與解除綁定流程。

3.  **階段三：核心功能**
    *   實作 `/new` 和 `/list` 指令，並確保其能正確地與 `UrlUsecase` 互動。
    *   實作 `/help` 指令。

4.  **階段四：測試與部署**
    *   撰寫單元測試與整合測試。
    *   部署並進行真人測試。

## 7. 具體實現步驟

### 階段一：基礎建設 (Infrastructure)

1.  **資料庫遷移 (Database Migration)**
    *   **檔案**: `sql/schema.sql`
    *   **操作**:
        *   在 `users` 資料表中新增 `telegram_chat_id BIGINT UNIQUE` 欄位。
        *   新增 `telegram_auth_tokens` 資料表，包含 `token TEXT PRIMARY KEY`, `telegram_chat_id BIGINT NOT NULL`, `expires_at TIMESTAMP NOT NULL` 欄位。
    *   **驗證**: 確認 Schema 修改無誤。

2.  **SQL 查詢生成 (SQL Query Generation)**
    *   **檔案**: `sql/queries/telegram_auth_tokens.sql` (新建)
    *   **操作**:
        *   新增 `CreateTelegramAuthToken`, `GetTelegramAuthToken`, `DeleteTelegramAuthToken` 的 SQLC 查詢。
    *   **檔案**: `sql/queries/users.sql`
    *   **操作**:
        *   新增 `GetUserByTelegramID` 查詢。
        *   新增 `UpdateUserTelegramID` 更新操作。
    *   **指令**: 執行 `sqlc generate` 來生成對應的 Go 程式碼。
    *   **驗證**: 檢查 `sqlc/` 目錄下是否已生成新的 Go 檔案和函式。

3.  **環境變數設定 (Configuration)**
    *   **檔案**: `config/config.go`
    *   **操作**:
        *   在 `Config` struct 中新增 `TelegramBotToken string`。
        *   更新 `LoadConfig` 函式，使其能從 `.env` 檔案讀取 `TELEGRAM_BOT_TOKEN`。
    *   **檔案**: `.env.example`
    *   **操作**: 新增 `TELEGRAM_BOT_TOKEN=""`。
    *   **驗證**: 確保程式啟動時能正確讀取到 Token。

### 階段二：認證流程 (Authentication Flow)

1.  **Telegram Use Case 建立**
    *   **檔案**: `application/telegram_usecase.go` (新建)
    *   **操作**:
        *   定義 `TelegramUsecase` interface，包含 `GenerateAuthToken`, `BindUserToTelegram` 等方法。
        *   實作 `telegramUsecase` struct，它將依賴 `UserRepository` 和 `TelegramTokenRepository`。

2.  **Bot 進入點與初始化**
    *   **檔案**: `main.go`
    *   **操作**:
        *   在 `main` 函式中，初始化 `TelegramUsecase`。
        *   初始化 `gotgbot` 的 Bot instance。
        *   啟動一個新的 Goroutine，執行 `infrastructure/telegram/bot_handler.go` 中的 Bot 主迴圈。

3.  **Bot 指令處理器 (`/start`)**
    *   **檔案**: `infrastructure/telegram/bot_handler.go` (新建)
    *   **操作**:
        *   建立 `NewBotHandler` 函式，接收 `TelegramUsecase` 作為依賴。
        *   實作 `Start` 方法，作為 Bot 的主迴圈，監聽來自 Telegram 的更新。
        *   實作 `/start` 指令的處理邏輯：
            *   呼叫 `telegramUsecase.GenerateAuthToken` 來生成一個新的 Token。
            *   建立驗證連結 `https://1li.tw/telegram/bind?token=...`。
            *   使用 `gotgbot` client 將此連結發送給使用者。

4.  **網頁綁定頁面 (Frontend)**
    *   **檔案**: `web/src/pages/telegram/bind.astro` (新建)
    *   **操作**:
        *   建立一個新的 Astro 頁面。
        *   此頁面在客戶端（瀏覽器）執行：
            *   從 URL 讀取 `token` 參數。
            *   檢查使用者是否已登入 (透過 `localStorage`)。若未登入，提示使用者先登入。
            *   若已登入，發送一個 `POST` 請求到後端的 `/api/telegram/bind`，並在 body 中附上 `token`。
            *   根據後端的回應，顯示成功或失敗的訊息。

5.  **後端綁定 API 端點**
    *   **檔案**: `presentation/gin/handler/http_handler.go`
    *   **操作**:
        *   新增 `BindTelegram` 處理函式。
        *   此函式從請求中讀取 `token`。
        *   從 JWT 中取得當前登入的 `userID`。
        *   呼叫 `telegramUsecase.BindUserToTelegram(userID, token)`。
    *   **檔案**: `presentation/gin/router.go`
    *   **操作**:
        *   在 `/api` 路由群組下，新增 `POST /telegram/bind` 路由，並將其指向 `BindTelegram` 處理函式。此路由需要 JWT 中介軟體保護。

### 階段三：核心功能 (Core Features)

1.  **指令路由與參數解析**
    *   **檔案**: `infrastructure/telegram/bot_handler.go`
    *   **操作**:
        *   擴充 Bot 主迴圈的邏輯，使其能根據指令 (e.g., `/new`, `/list`) 路由到不同的處理函式。
        *   為 `/new` 和 `/list` 指令實作參數解析邏輯。

2.  **`/new` 指令實作**
    *   **操作**:
        *   從訊息中解析出 `target_url` 和可選的 `custom_path`。
        *   透過 `telegram_chat_id` 查詢 `users` 資料表，取得對應的 `user_id`。
        *   呼叫 `UrlUsecase.CreateShortURL(userID, target_url, custom_path)`。
        *   將結果（成功或失敗訊息）回傳給 Telegram 使用者。

3.  **`/list` 指令實作**
    *   **操作**:
        *   解析出可選的 `page` 參數。
        *   透過 `telegram_chat_id` 查詢 `user_id`。
        *   呼叫 `UrlUsecase.GetUserURLs(userID, page, pageSize)`。
        *   將短網址列表格式化後，回傳給 Telegram 使用者。

4.  **`/help` 和 `/revoke` 指令實作**
    *   **操作**:
        *   為 `/help` 實作一個靜態的回覆，列出所有可用指令。
        *   為 `/revoke` 實作解除綁定的邏輯，將 `users` 表中的 `telegram_chat_id` 設為 `NULL`。

### 階段四：測試與部署 (Testing & Deployment)

1.  **單元測試 (Unit Testing)**
    *   **目標**: 針對 `application/telegram_usecase.go` 撰寫單元測試，使用 mock 的 repository。
    *   **驗證**: 確保 Token 生成和使用者綁定的核心邏輯正確。

2.  **整合測試 (Integration Testing)**
    *   **目標**: 建立一個測試用的 Bot Token，並在本地環境中完整地測試整個認證和指令流程。
    *   **驗證**:
        *   `/start` 是否能成功收到連結。
        *   點擊連結後，是否能成功綁定帳號。
        *   綁定後，`/new` 和 `/list` 是否能正常運作。

3.  **部署 (Deployment)**
    *   **操作**:
        *   將正式的 `TELEGRAM_BOT_TOKEN` 加入到生產環境的環境變數中。
        *   重新建置並部署應用程式。

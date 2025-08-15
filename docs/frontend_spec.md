# 前端網站相關規格

#### 3.5.1. Web 介面 (Astro)

- **訪客:**
    - 首頁，提供一個快速建立匿名、隨機短網址的表單。
    - 使用者登入與註冊頁面。
- **已登入使用者:**
    - **Dashboard:** 顯示使用者建立的短網址列表，包含原始網址、短網址以及**總點擊次數**。
    - **統計頁面:** 每個短網址都有一個獨立的詳細統計頁面，以視覺化圖表展示時間分佈、地理分佈等分析數據。
    - **新增短網址表單:** 提供選項以建立隨機或自訂路徑的短網址。
    - **刪除按鈕:** 在 Dashboard 列表中的每個短網址旁提供刪除功能。
- **管理者:**
    - **管理後台:** 一個獨立的頁面，顯示系統中所有的短網址，並提供搜尋和刪除功能。管理者也可以查看任何短網址的詳細統計頁面。
    - **使用者管理界面** 如果現在登入的使用者可以管理其他使用者，那他可以看到一個使用者管理界面，功能包括但不限於修改其他使用者的權限、重設其他使用者的密碼（會隨機建立一個足夠長的字串）

---

#### 前端（Astro）

- `astro`
- `react`
- `tailwindcss`
- `daisyui`
- `chart.js`

---

web/ # Astro 專案
src/
dist/

---

- **前端建置:** Astro 前端專案位於 `web/` 目錄。執行 `npm run build` 後，會將靜態資源輸出到 `web/dist/`。

---

- **資源嵌入:** Go 應用程式將使用 `embed` 套件來嵌入 `web/dist/` 目錄下的所有前端靜態檔案。

---

| 方法   | 路徑                                           | 簡介                                                               | 是否須驗證             | 輸入（請求） / 輸出（回應）                                                                                                                                                               | 對應的 Telegram 操作                                 |
| ------ | ---------------------------------------------- | ------------------------------------------------------------------ | ---------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------- |
| GET    | `/:short_path` <br> `/@:username/:custom_path` | 依短碼重定向至原始網址，並**非同步**紀錄點擊資訊（時間、國家、UA） | 否                     | **輸入**：Path 參數。<br>**輸出**：`301/302` Redirect 至 `original_url`。                                                                                                                 | 無（使用者直接點連結）                               |
| POST   | `/api/auth/register`                           | 使用者註冊                                                         | 否                     | **輸入**：JSON `{ "username": string, "password": string }`。<br>**輸出**：`201`，JSON `{ "id": number, "username": string }`                                                             | 無                                                   |
| POST   | `/api/auth/login`                              | 使用者登入並簽發 JWT                                               | 否                     | **輸入**：JSON `{ "username": string, "password": string }`。<br>**輸出**：`200`，Set-Cookie: JWT（HttpOnly）；JSON `{ "username": string, "permissions": int }`                          | 無                                                   |
| GET    | `/auth/telegram`                               | Telegram 授權頁（透過 `token` 完成帳號綁定流程的視覺化頁面）       | 否（但頁面動作需登入） | **輸入**：Query `token`（一次性、限時）。<br>**輸出**：HTML 頁面（導向登入／確認綁定）。                                                                                                  | `/auth` 會提供此頁的 URL                             |
| POST   | `/api/auth/telegram/link`                      | 確認並完成 Telegram 帳號綁定                                       | 是（Web 已登入）       | **輸入**：JSON `{ "token": string }`（並需 CSRF 保護）。<br>**輸出**：`200`，JSON `{ "ok": true }`                                                                                        | `/auth` 綁定流程的最終步驟                           |
| GET    | `/api/urls`                                    | 取得自己建立的短網址列表                                           | 是                     | **輸入**：無。<br>**輸出**：`200`，JSON `[{ "id": number, "short_path": string, "original_url": string, "total_clicks": number, "created_at": ISO8601 }]`                                 | `/list`（透過 Use Case）                             |
| POST   | `/api/urls`                                    | 新增短網址（隨機或自訂，依角色規則）                               | 是                     | **輸入**：JSON `{ "original_url": string, "custom_path": string(optional) }`。<br>**輸出**：`201`，JSON `{ "id": number, "short_path": string, "original_url": string }`                  | `/new <original_url> [custom_path]`（透過 Use Case） |
| DELETE | `/api/urls/:id`                                | 刪除自己建立的短網址                                               | 是                     | **輸入**：Path 參數 `id`。<br>**輸出**：`204` 無內容                                                                                                                                      | `/delete`（透過 Use Case）                           |
| GET    | `/api/urls/:id/stats`                          | 取得指定短網址的統計（時間、國家、系統、瀏覽器）                   | 是                     | **輸入**：Path 參數 `id`；可選 Query：`from`、`to`、`bucket`。<br>**輸出**：`200`，JSON `{ "total": number, "by_time": [...], "by_country": [...], "by_os": [...], "by_browser": [...] }` | `/stats`（透過 Use Case）                            |
| GET    | `/api/admin/urls`                              | 取得全系統短網址列表（管理功能）                                   | 管理者                 | **輸入**：可選 Query：分頁／搜尋。<br>**輸出**：`200`，JSON `[{ "id": number, "owner": string, "short_path": string, "original_url": string, "total_clicks": number }]`                   | 無                                                   |
| DELETE | `/api/admin/urls/:id`                          | 刪除任一短網址（管理功能）                                         | 管理者                 | **輸入**：Path 參數 `id`。<br>**輸出**：`204` 無內容                                                                                                                                      | 無                                                   |

---

為了方便測試，動態判斷 API URL，

```js
const API_URL = import.meta.env.PROD
	? `${window.location.origin}${window.location.pathname}/api`
	: 'http://localhost:8080/api'
```

# 功能規劃：管理員查看所有網址

## 1. 總覽

此功能旨在為具有特定權限的管理員提供一個集中的管理介面，讓他們可以查看、搜尋、並管理系統中所有使用者建立的短網址。這個頁面將類似於現有的使用者儀表板（Dashboard），但會顯示全站的資料。

**目標使用者：** 擁有 `CanViewAllStats` 或 `CanDeleteAny` 權限的管理員。

**主要功能：**

- 以列表形式展示系統中所有的短網址紀錄。
- 提供刪除任何一筆短網址的功能。
- 提供查看任一短網址統計資料的入口。
- （可選）提供搜尋和篩選功能。

---

## 2. 所需權限

此頁面的存取權限將嚴格受後端 API 控制。使用者必須至少擁有以下權限之一：

- `PermissionViewAllStats` (能夠查看所有統計資料)
- `PermissionDeleteAny` (能夠刪除任何短網址)

前端將根據使用者的權限來決定是否顯示此頁面的入口連結。

---

## 3. 後端實作規劃 (Go)

後端需要建立一個新的 API 端點來提供全站的短網址資料。

### 3.1. 資料庫層 (SQLC)

我們需要在 `sql/queries/short_urls.sql` 中新增一個查詢，以獲取所有的短網址以及其擁有者的資訊。現有的 `ListAllShortURLs` 僅返回 short_urls 表的欄位，我們需要擴充它來 JOIN users 表。

```sql
-- name: ListAllURLsWithUser :many
SELECT
    su.*,
    u.username
FROM
    short_urls su
JOIN
    users u ON su.user_id = u.id
WHERE su.deleted_at IS NULL
ORDER BY
    su.created_at DESC;
```

執行 `sqlc generate` 後，這將在 `sqlc/` 目錄下生成對應的 Go 函數 `ListAllURLsWithUser`。

### 3.2. 儲存庫層 (Repository)

在 `infrastructure/repository/shorturl_repository.go` 中，我們需要新增一個方法來呼叫由 sqlc 生成的 `ListAllURLsWithUser` 函數。

```go
// In ShortURLRepository interface
ListAllURLsWithUser(ctx context.Context) ([]sqlc.ListAllURLsWithUserRow, error)

// In postgresShortURLRepository struct
func (r *postgresShortURLRepository) ListAllURLsWithUser(ctx context.Context) ([]sqlc.ListAllURLsWithUserRow, error) {
    return r.queries.ListAllURLsWithUser(ctx)
}
```

### 3.3. 應用層 (Usecase)

在 `application/url_usecase.go` 中，新增一個處理業務邏輯的函數。這個函數會呼叫儲存庫的方法，然後將 `user_id` 和 `username` 組合成最終的 DTO。

```go
// In URLUsecase struct
func (uc *urlUsecase) GetAllURLs(ctx context.Context) ([]URLDetailsDTO, error) {
    // 1. 呼叫 repository 的 ListAllURLsWithUser 方法
    // 2. 迭代結果，建立一個包含 short_url 所有欄位以及 user_id 和 username 的 DTO 列表
    // 3. 返回 DTO 列表
}
```

### 3.4. 表現層 (API Handler & Router)

我們需要建立一個新的 API 端點。

1.  **Router (`presentation/gin/router.go`):**
    在 `private` 路由群組中新增一個管理員專用的子群組，並定義新的路由。

    ```go
    adminRoutes := private.Group("/admin")
    adminRoutes.Use(PermissionMiddleware(domain.PermissionViewAllStats | domain.PermissionDeleteAny)) // 使用權限中介軟體
    {
        adminRoutes.GET("/urls", httpHandler.GetAllURLs)
    }
    ```

2.  **Handler (`presentation/gin/handler/http_handler.go`):**
    新增對應的處理函數。

    ```go
    func (h *HTTPHandler) GetAllURLs(c *gin.Context) {
        // Call the urlUsecase.GetAllURLs
        // Return the list of URLs as JSON
    }
    ```

    回應的 JSON 格式將包含 `user_id` 和 `username`：

    ```json
    [
        {
            "id": 1,
            "short_path": "abcdef",
            "original_url": "https://example.com",
            "created_at": "2023-10-27T10:00:00Z",
            "user_id": 101,
            "owner_username": "admin"
        },
        ...
    ]
    ```

---

## 4. 前端實作規劃 (Astro + React)

前端需要建立一個新頁面和對應的元件來展示和操作資料。

### 4.1. API 客戶端

在 `web/src/lib/api.ts` 中新增一個函數，用來呼叫後端新的 `GET /api/admin/urls` 端點。

```typescript
// In api.ts
export const getAllUrls = async (): Promise<AdminUrlData[]> => {
	const res = await apiClient.get('/admin/urls')
	return res.data
}

// Define the AdminUrlData type
export interface AdminUrlData {
	id: number
	short_path: string
	original_url: string
	created_at: string
	user_id: number
	owner_username: string
}
```

### 4.2. 新頁面與元件

1.  **建立新頁面 (`web/src/pages/admin/all-urls.astro`):**
    這個 Astro 頁面將作為容器，並引入 React 元件。它會處理頁面的基本佈局和標題。

2.  **建立新元件 (`web/src/components/AllUrlsDashboard.tsx`):**
    這個 React 元件是功能的核心。
    - 使用 `useEffect` 和 `useState` 來呼叫 `getAllUrls` API 並儲存資料。
    - 將資料渲染成一個表格，類似於現有的 `Dashboard.tsx`。
    - 表格欄位應包含：原始網址、短網址、擁有者 (`owner_username`)、建立時間、操作。
    - **操作欄** 包含兩個按鈕：
        - **統計 (Stats):** 點擊後導航至 `dashboard/stats?p=<short_path>`。
        - **刪除 (Delete):** 點擊後會跳出確認對話框，確認後呼叫現有的 `deleteUrl` API，並在成功後重新整理列表。

### 4.3. 導航與入口

在 `web/src/layouts/Layout.astro` 或其他共享的導航元件中，新增一個指向 `/admin/all-urls` 的連結。這個連結需要使用 `useUser` hook 來檢查使用者的權限，只有當 `user.permissions` 包含 `CanViewAllStats` 或 `CanDeleteAny` 時才顯示。

```typescript
// Example logic in a React component inside the layout
import { useUser } from '../hooks/useUser';
import { PERMISSIONS } from '../lib/permissions';

const NavLinks = () => {
  const { user } = useUser();
  const canViewAdminPage = user && (
    (user.permissions & PERMISSIONS.CanViewAllStats) !== 0 ||
    (user.permissions & PERMISSIONS.CanDeleteAny) !== 0
  );

  return (
    <nav>
      {/* other links */}
      {canViewAdminPage && <a href="/admin/all-urls">Admin Panel</a>}
    </nav>
  );
};
```

---

## 5. 總結：需變更或新增的檔案

- **Backend:**
    - `sql/queries/short_urls.sql` (新增或修改 SQL 查詢)
    - `infrastructure/repository/shorturl_repository.go` (新增 Repository 方法)
    - `application/url_usecase.go` (新增 Usecase 方法)
    - `presentation/gin/router.go` (新增路由)
    - `presentation/gin/handler/http_handler.go` (新增 Handler)
- **Frontend:**
    - `web/src/lib/api.ts` (新增 API 呼叫函數)
    - `web/src/pages/admin/all-urls.astro` (新頁面)
    - `web/src/components/AllUrlsDashboard.tsx` (新元件)
    - `web/src/layouts/Layout.astro` (或相關導航元件，新增入口連結)

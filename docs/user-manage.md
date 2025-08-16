# 使用者管理介面建立計畫

本文檔旨在規劃和定義建立一個使用者管理介面所需的功能與實作步驟。

## 1. 目標

建立一個後台管理介面，允許管理員 (Admin) 執行以下操作：

- 查看所有使用者列表。
- 修改指定使用者的權限。
- 刪除使用者（軟刪除）。
- 保留未來擴充更多管理操作的彈性（例如：重設密碼、查看使用者詳細資訊等）。

## 2. 功能分析與現況

### 現有程式碼分析

- **`sql/queries/users.sql`**:

    - `UpdateUserPermissions`: 已定義，可用於修改權限。
    - `DeleteUser`: 已定義，使用軟刪除。
    - `ListUsers`: **缺少**，需要一個查詢來獲取所有使用者列表。

- **`infrastructure/repository/user_repository.go`**:

    - `Update`: 已存在，但邏輯較為混雜（同時處理 Telegram ID 和權限）。為了清晰起見，可以考慮拆分或新增一個專門更新權限的方法。
    - `Delete`: **缺少** `Delete` 方法的實作。
    - `List`: **缺少** `List` 方法的實作。

- **`presentation/gin/handler/http_handler.go`**:
    - 目前只有 `AuthHandler` 和 `URLHandler`。
    - **缺少** 用於管理使用者的 `UserHandler`，以及對應的 API 端點（List, Update, Delete）。

## 3. 實作計畫

### 第一階段：後端 API 建置

#### 步驟 1: 資料庫查詢層 (SQL)

- **檔案**: `sql/queries/users.sql`
- **操作**: 新增一個 `ListUsers` 查詢。

```sql
-- name: ListUsers :many
SELECT *
FROM users
WHERE deleted_at IS NULL
ORDER BY id;
```

#### 步驟 2: Domain 層

- **檔案**: `domain/user.go`
- **操作**: 更新 `UserRepository` 介面，加入必要的方法。

```go
type UserRepository interface {
    // ... (現有方法)
    List(ctx context.Context) ([]*User, error)
    UpdatePermissions(ctx context.Context, userID int64, permissions Permission) error
    Delete(ctx context.Context, userID int64) error
}
```

#### 步驟 3: Repository 層

- **檔案**: `infrastructure/repository/user_repository.go`
- **操作**: 實作 `UserRepository` 介面中新增的方法。

1.  **`List`**: 實作 `ListUsers` 查詢。
2.  **`Delete`**: 實作 `DeleteUser` 查詢。
3.  **`UpdatePermissions`**: 建立一個新方法，專門呼叫 `UpdateUserPermissions` 查詢，讓職責更清晰。

#### 步驟 4: Usecase 層

- **檔案**: `application/user_usecase.go`
- **操作**: 在 `UserUseCase` 中新增業務邏輯方法。

1.  **`ListUsers(ctx context.Context, requester *domain.User) ([]*domain.User, error)`**:

    - 檢查 `requester` 是否有 `PermissionAdmin` 權限。
    - 若無權限，返回錯誤。
    - 若有權限，呼叫 `userRepository.List()`。

2.  **`UpdateUserPermissions(ctx context.Context, requester *domain.User, targetUserID int64, permissions domain.Permission) error`**:

    - 檢查 `requester` 是否有 `PermissionAdmin` 權限。
    - 呼叫 `userRepository.UpdatePermissions()`。

3.  **`DeleteUser(ctx context.Context, requester *domain.User, targetUserID int64) error`**:
    - 檢查 `requester` 是否有 `PermissionAdmin` 權限。
    - 檢查是否正在刪除自己，如果是，則不允許。
    - 呼叫 `userRepository.Delete()`。

#### 步驟 5: API Handler 層

- **檔案**: `presentation/gin/handler/http_handler.go`
- **操作**: 建立一個新的 `UserHandler`。

```go
type UserHandler struct {
    userUseCase *application.UserUseCase
}

func NewUserHandler(userUC *application.UserUseCase) *UserHandler {
    // ...
}

func (h *UserHandler) ListUsers(c *gin.Context) {
    // ...
}

func (h *UserHandler) UpdateUserPermissions(c *gin.Context) {
    // ...
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
    // ...
}
```

#### 步驟 6: 路由層 (Router)

- **檔案**: `presentation/gin/router.go`
- **操作**: 新增使用者管理相關的 API 路由，並使用中介軟體進行權限驗證。

- 建立一個新的路由群組 `/api/admin`，並套用 `AuthMiddleware` 和一個新的 `AdminRequiredMiddleware`。
- `AdminRequiredMiddleware` 將檢查 `c.Get("user")` 的權限是否包含 `domain.PermissionAdmin`。

```go
adminRoutes := router.Group("/api/admin")
adminRoutes.Use(AuthMiddleware(tokenSvc))
adminRoutes.Use(AdminRequiredMiddleware()) // 新增的中介軟體
{
    userHandler := handler.NewUserHandler(userUseCase)
    adminRoutes.GET("/users", userHandler.ListUsers)
    adminRoutes.PUT("/users/:id/permissions", userHandler.UpdateUserPermissions)
    adminRoutes.DELETE("/users/:id", userHandler.DeleteUser)
}
```

### 第二階段：前端介面開發

#### 步驟 1: API 客戶端

- **檔案**: `web/src/lib/api.ts`
- **操作**: 新增呼叫後端 Admin API 的函式。
    - `listUsers()`
    - `updateUserPermissions(userId, permissions)`
    - `deleteUser(userId)`

#### 步驟 2: UI 元件

- **檔案**: `web/src/components/UserManagement.tsx`
- **操作**: 建立一個新的 React 元件。
    - 使用 `useEffect` 呼叫 `listUsers` 獲取使用者列表。
    - 將使用者列表渲染成表格或列表。
    - 每一行包含：`ID`, `Username`, `Permissions`, `Created At` 以及操作按鈕。
    - **權限修改**: 提供一個下拉選單或一組 checkboxes 來修改權限，並在變更時呼叫 `updateUserPermissions`。
    - **刪除按鈕**: 點擊後跳出確認對話框，確認後呼叫 `deleteUser`。

#### 步驟 3: 頁面

- **檔案**: `web/src/pages/dashboard/admin/users.astro`
- **操作**: 建立一個新的頁面來渲染 `UserManagement` 元件。
    - 此頁面應檢查使用者登入狀態與權限，只有管理員可見。可以在 `Layout.astro` 或頁面本身進行權限檢查。

#### 步驟 4: 導覽列

- **檔案**: `web/src/layouts/Layout.astro` 或相關導覽列元件
- **操作**: 在導覽列中新增一個「使用者管理」的連結。
    - 此連結應使用 `useUser` hook 檢查使用者權限，僅在使用者為管理員時顯示。

## 4. 預期 API 設計

- `GET /api/admin/users`
    - **Success Response (200)**: `[{ "id": 1, "username": "admin", "permissions": 3, ... }]`
- `PUT /api/admin/users/:id/permissions`
    - **Request Body**: `{ "permissions": 1 }`
    - **Success Response (204)**: No Content
- `DELETE /api/admin/users/:id`
    - **Success Response (204)**: No Content

# Todo List MCP Tool

這是一個基於 MCP (Mission Control Protocol) 的待辦事項管理工具，提供簡單且高效的待辦事項管理功能。該工具使用 Go 語言開發，並完全在記憶體中運行，適合用於臨時性的任務管理。

## 功能特點

- 純記憶體儲存，輕量快速
- 支援多個同時操作的執行緒安全設計
- 簡潔的命令列介面
- 完整的待辦事項生命週期管理

## 可用指令

該工具提供以下四個主要功能：

### 1. 新增待辦事項 (add_todo)
新增一個待辦事項到清單中。
- 參數：`title` - 待辦事項的標題
- 範例：新增一個待辦事項 "買牛奶"

### 2. 刪除待辦事項 (delete_todo)
從清單中刪除指定的待辦事項。
- 參數：`id` - 待辦事項的唯一識別碼
- 範例：刪除 ID 為 "123456789" 的待辦事項

### 3. 查看所有待辦事項 (get_todos)
顯示所有待辦事項的清單。
- 無需參數
- 顯示格式：`[x]` 表示已完成，`[ ]` 表示未完成

### 4. 更新待辦事項狀態 (set_todo_completion)
更改待辦事項的完成狀態。
- 參數：
  - `id` - 待辦事項的唯一識別碼
  - `completed` - 完成狀態（true/false）

## 技術特點

- 採用 Unix 時間戳作為唯一識別碼
- 支援優雅的錯誤處理
- 內建日誌記錄和錯誤恢復機制

## 資料結構

每個待辦事項包含以下資訊：
- ID：唯一識別碼
- 標題：待辦事項的內容
- 完成狀態：是否已完成
- 創建時間：待辦事項的建立時間

## 注意事項

- 此工具為記憶體儲存，程式重啟後資料會被清空
- 建議用於臨時性的任務管理
- 所有操作都是即時的，無需手動保存

## 系統需求

- Go 1.16 或更高版本
- 支援標準輸入輸出的終端機環境

## 伺服器配置

使用以下配置來設定 MCP Server：

```json
{
  "todo-list-mcp-server": {
    "command": "/Users/lidongying/Documents/Projects/todo-list-mcp/main"
  }
}
```

## 授權

本專案採用 MIT 授權條款。
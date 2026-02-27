# GiftFlow - Internal Gift Application System

## Quick Start

### 1. Start the Environment
Run the following command in the project root:
```bash
docker compose up --build
```
Access the application at: **http://<SERVER_IP>:3000**

### 2. Stop the Environment
```bash
docker compose down
```

---

## Demo Script

### Scenario A: Normal Gift Approval (Dept Admin -> Dept Head)

1.  **Login as Applicant**
    *   Open `http://<SERVER_IP>:3000`.
    *   Click **"Dept Admin (Applicant)"**.
    *   You will see the "Browse Gifts" page.

2.  **Submit Application**
    *   Choice **"Any gift image icon"** (labeled `NORMAL`).
    *   Click **Apply**.
    *   Confirm the dialog.
    *   Click **"My Applications"** in the navbar to verify status is `PENDING_HEAD`.
    *   Click **Logout**.

3.  **Dept Head Approval**
    *   Login as **"Dept Head (Approver)"**.
    *   You will see the "Pending Approvals" dashboard.
    *   Find the "Any gift image icon" application.
    *   Click **Approve**.
    *   Enter a comment (e.g., "Approved for dept").

### Scenario B: VIP Gift Approval (Dept Admin -> Dept Head -> CPRO Admin)

1.  **Submit VIP Application**
    *   Login as **"Dept Admin (Applicant)"**.
    *   Choice **"Any gift image icon"** (labeled `VIP`).
    *   Click **Apply**.
    *   Verify status is `PENDING_HEAD` in "My Applications".
    *   Logout.

2.  **Dept Head Approval (Level 1)**
    *   Login as **"Dept Head (Approver)"**.
    *   Find the "Any gift image icon" application.
    *   Click **Approve**.
    *   Enter comment: "Endorsed for VIP guest".
    *   **Note:** Status changes to `PENDING_CPRO` (Orange), not Completed.

3.  **CPRO Admin Approval (Level 2)**
    *   Login as **"CPRO Admin (VIP Approver)"**.
    *   You will see the VIP application pending your approval.
    *   Click **Approve**.
    *   Enter comment: "Stock allocated".

---

## Technical Notes

### Architecture
*   **Backend:** Go (Gin) + MongoDB.
*   **Frontend:** VueJS + TailwindCSS.
*   **Infrastructure:** Docker Compose (Backend, Frontend, MongoDB).
*   **Deployment Host:** Frontend defaults API calls to `http://<YOUR_HOST>:8080` based on the browser hostname. You can also set `VITE_API_BASE_URL` at build time if backend is on a different host.
*   **Public File URL:** Backend upload URL is generated from request host, or `PUBLIC_BASE_URL` if provided.

### OIDC Integration (Future)
The authentication logic is centralized in `backend/middleware/auth.go`.
To integrate Microsoft OIDC:
1.  Replace `AuthMiddleware` to validate JWT tokens from Azure AD.
2.  Update `MockUsers` lookup to query a real User database or use claims from the token.

## Go Packages Reference
Key Go libraries used in this project:

| Package | Import Path | Description | Functionality |
| :--- | :--- | :--- | :--- |
| **Gin** | `github.com/gin-gonic/gin` | [Web Framework](https://github.com/gin-gonic/gin) | Provides the HTTP server, routing, middleware (CORS, Auth), and JSON response handling. High-performance web framework. |
| **Mongo Driver** | `go.mongodb.org/mongo-driver` | [MongoDB Driver for Go](https://github.com/mongodb/mongo-go-driver) | Handles all interactions with the MongoDB database, including connection pooling, BSON document mapping, and CRUD operations. |
| **Validator** | `github.com/go-playground/validator/v10` | [Struct and Field Validation](https://github.com/go-playground/validator) | Validates incoming request data structures (struct tags) to ensure required fields and data formats are correct before processing. |
| **Sonic** | `github.com/bytedance/sonic` | [High-performance JSON Library](https://github.com/bytedance/sonic) | Used internally by Gin (where available) for faster JSON serialization and deserialization, optimizing API response times. |


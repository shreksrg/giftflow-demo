# GiftFlow - Internal Gift Application System

## Quick Start

### 1. Start the Environment
Run the following command in the project root:
```bash
docker compose up --build
```
Access the application at: **http://localhost:3000**

### 2. Stop the Environment
```bash
docker compose down
```

---

## Demo Script

### Scenario A: Normal Gift Approval (Dept Admin -> Dept Head)

1.  **Login as Applicant**
    *   Open `http://localhost:3000`.
    *   Click **"Dept Admin (Applicant)"**.
    *   You will see the "Browse Gifts" page.

2.  **Submit Application**
    *   Find **"Notebook Set"** (labeled `NORMAL`).
    *   Click **Apply**.
    *   Confirm the dialog.
    *   Click **"My Applications"** in the navbar to verify status is `PENDING_HEAD`.
    *   Click **Logout**.

3.  **Dept Head Approval**
    *   Login as **"Dept Head (Approver)"**.
    *   You will see the "Pending Approvals" dashboard.
    *   Find the "Notebook Set" application.
    *   Click **Approve**.
    *   Enter a comment (e.g., "Approved for team").
    *   Status changes to `COMPLETED` (Green).

### Scenario B: VIP Gift Approval (Dept Admin -> Dept Head -> CPRO Admin)

1.  **Submit VIP Application**
    *   Login as **"Dept Admin (Applicant)"**.
    *   Find **"VIP Hamper"** (labeled `VIP`).
    *   Click **Apply**.
    *   Verify status is `PENDING_HEAD` in "My Applications".
    *   Logout.

2.  **Dept Head Approval (Level 1)**
    *   Login as **"Dept Head (Approver)"**.
    *   Find the "VIP Hamper" application.
    *   Click **Approve**.
    *   Enter comment: "Endorsed for VIP guest".
    *   **Note:** Status changes to `PENDING_CPRO` (Orange), not Completed.

3.  **CPRO Admin Approval (Level 2)**
    *   Login as **"CPRO Admin (VIP Approver)"**.
    *   You will see the VIP application pending your approval.
    *   Click **Approve**.
    *   Enter comment: "Stock allocated".
    *   Status changes to `COMPLETED` (Green).

---

## Technical Notes

### Architecture
*   **Backend:** Go (Gin) + MongoDB.
*   **Frontend:** VueJS + TailwindCSS.
*   **Infrastructure:** Docker Compose (Backend, Frontend, MongoDB).

### OIDC Integration (Future)
The authentication logic is centralized in `backend/middleware/auth.go`.
To integrate Microsoft OIDC:
1.  Replace `AuthMiddleware` to validate JWT tokens from Azure AD.
2.  Update `MockUsers` lookup to query a real User database or use claims from the token.

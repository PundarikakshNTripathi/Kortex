---
description: Standard Git Feature Branch Workflow for Kortex
---

This workflow ensures a clean history and avoids "ahead/behind" confusion by keeping `main` as the source of truth.

1. **Start Fresh**: Switch to `main` and get the latest changes.
   ```bash
   git checkout main
   git pull origin main
   ```

2. **Create Feature Branch**: Create a new branch for your task.
   ```bash
   # Replace 'feature/my-task-name' with a descriptive name
   git checkout -b feature/my-task-name
   ```

3. **Work & Commit**: Make your changes and commit them.
   ```bash
   git add .
   git commit -m "feat: description of changes"
   ```

4. **Sync with Main** (Optional but recommended for long tasks):
   If `main` has updated while you were working, rebase your feature branch on top of it.
   ```bash
   git checkout main
   git pull origin main
   git checkout feature/my-task-name
   git rebase main
   ```

5. **Merge**: Switch back to `main` and merge your feature.
   ```bash
   git checkout main
   git pull origin main  # Just to be sure
   git merge feature/my-task-name
   ```

6. **Push**: Push the updated `main` to the server.
   ```bash
   git push origin main
   ```

7. **Cleanup**: Delete the feature branch.
   ```bash
   git branch -d feature/my-task-name
   ```

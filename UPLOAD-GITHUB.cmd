@echo off
setlocal
chcp 65001 >nul
cd /d "%~dp0"
echo ===============================================
echo  LUAN GROUP - Upload project to GitHub
 echo  Repository: luanlengoc1991-art/APP-TACH-NEN-LUAN
 echo ===============================================
echo.
where git >nul 2>&1
if errorlevel 1 (
  echo [ERROR] Git chua duoc cai dat. Cai Git for Windows tai: https://git-scm.com/downloads/win
  pause
  exit /b 1
)
if not exist ".git\" (
  git init -b main
  if errorlevel 1 goto failure
)
git branch -M main
set "REMOTE=https://github.com/luanlengoc1991-art/APP-TACH-NEN-LUAN.git"
git remote get-url origin >nul 2>&1
if errorlevel 1 (
  git remote add origin "%REMOTE%"
) else (
  git remote set-url origin "%REMOTE%"
)
git add README.md go.mod main.go web .gitignore .github/workflows/windows-build.yml UPLOAD-GITHUB.cmd
if errorlevel 1 goto failure
git diff --cached --quiet
if not errorlevel 1 goto push
git commit -m "Add LUAN GROUP AI Layer Studio desktop source and Windows auto-build"
if errorlevel 1 goto failure
:push
echo.
echo Uploading branch main - Git may ask you to sign in with GitHub...
git push -u origin main
if errorlevel 1 goto failure
echo.
echo [SUCCESS] Da cap nhat code len GitHub.
echo https://github.com/luanlengoc1991-art/APP-TACH-NEN-LUAN
start "" "https://github.com/luanlengoc1991-art/APP-TACH-NEN-LUAN"
pause
exit /b 0
:failure
echo.
echo [ERROR] Upload chua thanh cong. Kiem tra dang nhap GitHub, Git name/email va quyen Push.
pause
exit /b 1

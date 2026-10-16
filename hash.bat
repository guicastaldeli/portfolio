@echo off
setlocal enabledelayedexpansion

cd /d "%~dp0app"

if not exist "hash.go" (
    echo [hash] hash.go not found
    exit /b 1
)

if not exist ".env\.env.pw" (
    echo [hash] .env\.env.pw not found in %CD%
    exit /b 1
)

set "PW="
for /f "usebackq tokens=1,* delims==" %%A in (".env\.env.pw") do (
    set "KEY=%%A"
    set "VAL=%%B"
    set "VAL=!VAL:"=!"
    if /i "!KEY!"=="ADMIN_PASSWORD_TEXT" (
        set "PW=!VAL!"
    )
)

if "!PW!"=="" (
    echo [hash] ADMIN_PASSWORD_TEXT not found in .env\.env.pw
    exit /b 1
)

for /f "delims=" %%H in ('go run hash.go "!PW!"') do set "HASH=%%H"

echo.
echo Hash:
echo !HASH!
echo.
echo !HASH!| clip
echo Copied to clipboard. Paste into ADMIN_PASSWORD in .env.server, then restart.
endlocal
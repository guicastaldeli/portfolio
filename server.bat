@echo off
setlocal

if "%1"=="" (
    set APP_ENV=dev
) else (
    set APP_ENV=%1
)

cd /d "%~dp0app"

echo Running in %APP_ENV% environment...
set APP_ENV=%APP_ENV%
go run .
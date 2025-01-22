@echo off
setlocal enabledelayedexpansion

:: Check if cleanup parameter is passed
if "%1"=="cleanup" (
    echo Cleaning up development processes...
    taskkill /F /IM templ.exe 2>nul
    taskkill /F /IM tailwindcss-windows-x64.exe 2>nul
    taskkill /F /IM air.exe 2>nul
    echo Cleanup complete!
    exit /b 0
)

:: Run the non-watch commands first
echo Running formatters...
call templ fmt .
call go fmt ./...

:: Start each watch process in a new tab of Windows Terminal with custom names and colors
:: Colors: 
:: - Templ: Dark blue background (#000080)
:: - Tailwind: Purple background (#012456)
:: - Air: Dark green background (#004400)
wt -w 0 new-tab --title "Templ Watcher" --tabColor "#000080" -d "%CD%" cmd /k "title Templ Watcher && color 1F && echo Starting Templ watcher... && templ generate --watch" ^
; new-tab --title "Tailwind Watcher" --tabColor "#012456" -d "%CD%" cmd /k "title Tailwind Watcher && color 5F && echo Starting Tailwind watcher... && tailwindcss\tailwindcss-windows-x64.exe -i ./static/css/input.css -o ./static/css/output.css --watch" ^
; new-tab --title "Air Live Reload" --tabColor "#004400" -d "%CD%" cmd /k "title Air Live Reload && color 2F && echo Starting Air... && air"

:: Create a cleanup script if it doesn't exist
echo @echo off > end.bat
echo echo Cleaning up development processes... >> end.bat
echo taskkill /F /IM templ.exe 2^>nul >> end.bat
echo taskkill /F /IM tailwindcss-windows-x64.exe 2^>nul >> end.bat
echo taskkill /F /IM air.exe 2^>nul >> end.bat
echo echo Cleanup complete! >> end.bat

echo Development environment started!
echo.
echo Navigation:
echo - Ctrl + Tab: Cycle through tabs
echo - Ctrl + Number: Switch to specific tab
echo - Ctrl + Shift + W: Close current tab
echo.
echo To stop all processes later, run: end.bat
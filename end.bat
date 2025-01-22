@echo off 
echo Cleaning up development processes... 
taskkill /F /IM templ.exe 2>nul 
taskkill /F /IM tailwindcss-windows-x64.exe 2>nul 
taskkill /F /IM air.exe 2>nul 
echo Cleanup complete 

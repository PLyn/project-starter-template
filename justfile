set shell := ["cmd.exe", "/c"]

# Run all development tools
dev:
    templ fmt .
    templ generate .
    go fmt ./...
    tailwindcss\tailwindcss-windows-x64.exe -i ./static/css/input.css -o ./static/css/output.css
    go run .

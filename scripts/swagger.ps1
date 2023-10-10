# Generates Swagger documention

swag fmt

swag init -g ./cmd/server/main.go

# Replace 'x-nullable' with 'nullable' to proper display it in the Swagger

(Get-Content -Path .\docs\docs.go -Raw) -replace 'x-nullable','nullable' | Set-Content -Path .\docs\docs.go
(Get-Content -Path .\docs\swagger.json -Raw) -replace 'x-nullable','nullable' | Set-Content -Path .\docs\swagger.json
(Get-Content -Path .\docs\swagger.yaml -Raw) -replace 'x-nullable','nullable' | Set-Content -Path .\docs\swagger.yaml

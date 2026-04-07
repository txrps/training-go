Get-Content .env | ForEach-Object {
    if ($_ -match "^([^#][^=]+)=(.*)$") {
        Set-Item -Path "Env:$($matches[1])" -Value $matches[2]
    }
}

$command = $args[0]
$name = $args[1]

# (optional) fix path migrate ให้ชัวร์
$migrate = "migrate" 
# หรือใช้แบบนี้ถ้า PATH มั่ว:
# $migrate = "$env:USERPROFILE\go\bin\migrate.exe"

switch ($command) {
    "create" {
        & $migrate create -ext sql -dir migrations -seq $name
    }
    "up" {
        & $migrate -path migrations -database postgresql://postgres:P@ssw0rd@localhost:5434/testing_db?sslmode=disable up
    }
    "down" {
        & $migrate -path migrations -database postgresql://postgres:P@ssw0rd@localhost:5434/testing_db?sslmode=disable down
    }
    "force" {
        & $migrate -path migrations -database postgresql://postgres:P@ssw0rd@localhost:5434/testing_db?sslmode=disable force $name
    }
    default {
        Write-Host "Usage: .\migrate.ps1 <create|up|down|force> <migration_name>"
    }
}
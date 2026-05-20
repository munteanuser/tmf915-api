# Run all TMF915 migrations against DESKTOP-V2M1TS5 using Windows Integrated Auth
param(
    [string]$Server = "DESKTOP-V2M1TS5"
)

$ErrorActionPreference = "Stop"

Write-Host "Running TMF915 migrations on $Server..." -ForegroundColor Cyan

$scripts = @(
    "migrations\001_create_database.sql",
    "migrations\002_create_aimodelspecification.sql",
    "migrations\003_create_aimodel.sql",
    "migrations\004_create_aicontractspecification.sql",
    "migrations\005_create_aicontract.sql",
    "migrations\006_create_aicontractviolation.sql",
    "migrations\007_create_alarm.sql",
    "migrations\008_create_rule.sql",
    "migrations\009_create_monitor.sql",
    "migrations\010_create_topic.sql",
    "migrations\011_create_event.sql",
    "migrations\012_create_hub_listener.sql"
)

foreach ($script in $scripts) {
    Write-Host "  Applying $script..." -ForegroundColor Yellow
    # 001 runs against master to create the DB; the rest run against TMF915
    if ($script -match "001") {
        sqlcmd -S $Server -E -C -i $script
    } else {
        sqlcmd -S $Server -d TMF915 -E -C -i $script
    }
    if ($LASTEXITCODE -ne 0) {
        Write-Error "Migration $script failed with exit code $LASTEXITCODE"
        exit 1
    }
}

Write-Host "All migrations applied successfully!" -ForegroundColor Green

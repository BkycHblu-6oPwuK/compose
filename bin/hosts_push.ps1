param (
    [string]$HOST_FILE_WINDOWS,
    [string]$HOST_ENTRY
)

if (-not (Select-String -Path $HOST_FILE_WINDOWS -Pattern ([regex]::Escape($HOST_ENTRY)))) {
    Start-Process powershell -Verb RunAs -Wait -PassThru -ArgumentList '-NoProfile', '-ExecutionPolicy', 'Bypass', '-Command', "Start-Sleep -Seconds 1; Add-Content -Path '$HOST_FILE_WINDOWS' -Value '`n$HOST_ENTRY'"
    Write-Output 'The entry has been added to the hosts file.'
} else {
    Write-Output 'The entry already exists in the hosts file.'
}

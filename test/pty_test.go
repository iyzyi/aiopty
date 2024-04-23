package test

// TODO: Supplementary test cases :-)

// Command for testing pty window size
// Most Unix-like systems:
// yes = | head -n$(($(tput cols) * $(tput lines))) | tr -d '\n'
// eval printf '=%.0s' {1..$[$COLUMNS*$LINES]}
// Powershell:
// Write-Host ("=" * $(Get-Host).UI.RawUI.WindowSize.Height * $(Get-Host).UI.RawUI.WindowSize.Width)

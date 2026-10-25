# Privacy/portability notice: personal paths and credentials are intentionally omitted.
# Bundled file references are relative to this file; machine tools are resolved at runtime.
# Provide local-only values through environment variables; never commit secrets.

<#
.DESCRIPTION
sudo command but for Windowsユーザーお嬢様
.EXAMPLE
# Usage example: .\ojosudo.ps1 "notepad .\notes.txt"
#>
param
(
    [Parameter(Mandatory=$true, ValueFromPipeline=$true, ValueFromPipelineByPropertyName=$false)]
    [string] $Command
)

<#
Main block
#>
[ScriptBlock] $mainBlock = {
    Write-OjosamaMessage
    $evacuatedTitle = $Host.UI.RawUI.WindowTitle

    # FIXME
    #Start-AttachedConsole -ParentPID $PID
    $base64Cmd = Convert-CommandToBase64String -Command $Command
    $proc = New-ElevatedProcess -WorkingDirectory . -EncodedCommand $base64Cmd
    $exitCode = Invoke-Process -Process $proc
    $Host.UI.RawUI.WindowTitle = $evacuatedTitle

    exit $exitCode
}

<#
function Start-AttachedConsole
{
    param
    (
        [int] $ParentPID
    )
    $kernel = Add-Type -PassThru -TypeDefinition '
    using System.Runtime.InteropServices;
    public class Kernel
    {
        [DllImport("kernel32.dll", SetLastError=true)]
        public static extern bool AttachConsole(uint dwProcessId);
        [DllImport("kernel32.dll", SetLastError=true, ExactSpelling=true)]
        public static extern bool FreeConsole();
    }
    '
    $kernel::FreeConsole()
    $kernel::AttachConsole($ParentPID)
}
#>

function New-ElevatedProcess
{
    param
    (
        [string] $WorkingDirectory
        ,
        [string] $EncodedCommand
    )
#-NoProfile -ExecutionPolicy unrestricted -Command "Start-Process wt -Verb runas -ArgumentLis
    $proc = [System.Diagnostics.Process]::new()
    $startInfo = $proc.StartInfo
    $startInfo.FileName = "pwsh.exe"
    $startInfo.Verb = "runas"
    #$startInfo.WindowStyle = "Hidden"
    $startInfo.UseShellExecute = $true
    $startInfo.Arguments = @(
        "-NoProfile",
        #"-EncodedCommand", $EncodedCommand,
        "-NoExit"
        "-Command Start-Process wt"
    ) -join " "

    Write-Output $proc
}

function Invoke-Process
{
    param
    (
        [System.Diagnostics.Process] $Process
    )

    try
    {
        $ok = $Process.Start(); if (-not $ok) {
            Write-Error "Failed to start a process: $($Process | Format-List)"
        }
        $Process.WaitForExit()
    } catch {
        Write-Error "Failed to start a process: $($Process | Format-List)"
        throw
    } finally
    {
        Write-Output $Process.ExitCode
    }
}

function Convert-CommandToBase64String
{
    param
    (
        [string] $Command
    )

    $bytes = [System.Text.Encoding]::Unicode.GetBytes($Command)
    [System.Convert]::ToBase64String($bytes) | Write-Output
}

function Test-DoIHaveAdministratorPrivledge
{
    $me = [System.Security.Principal.WindowsIdentity]::GetCurrent() -as [System.Security.Principal.WindowsPrincipal]
    $admin = [System.Security.Principal.WindowsBuiltInRole]::Administrator

    $me.IsInRole($admin) | Write-Output
}

function Write-OjosamaMessage
{
    $white = "$([char]0x1b)[37;1m"
    $cyan = "$([char]0x1b)[36;1m"
    $neutral = "$([char]0x1b)[m"
    @(
        "${white}貴方は領主様から通常の講習を受けたはずですわ。${neutral}"
        "${white}それらは通常に以下3点に要約されますの。${neutral}"
        "${cyan}#1${neutral}) 市民の皆様のプライバシーを尊重すること。"
        "${cyan}#2${neutral}) おタイプする前に考えること。"
        "${cyan}#3${neutral}) そしてノブレス・オブリージュを肝に銘じておくことですわ！"
    ) -join "`n" | Write-Output
}

<#
Invoke main block
#>
& $mainBlock
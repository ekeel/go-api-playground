Import-Module Pester -Force -ErrorAction SilentlyContinue

Describe 'Basic Functional Tests' {
  1..1 | % {
    It 'Test Go plugin' {
      # $ret = curl -s http://localhost:8000/test | ConvertFrom-Json
      $ret = curl http://localhost:8000/test | ConvertFrom-Json

      (($ret -contains '1111') -and ($ret -contains '2222') -and ($ret -contains '3333') -and ($ret -contains '4444')) | Should -Be $true
    }
    
    It 'Test Python plugin' {
      # $ret = curl -s http://localhost:8000/testpy/val2 | ConvertFrom-Json
      $ret = curl http://localhost:8000/testpy/val2 | ConvertFrom-Json

      $ret.test | Should -Be 'val2'
    }

    It 'Test Powershell plugin' {
      # $ret = curl -s http://localhost:8000/testps/val1/val2 | ConvertFrom-Json
      $ret = curl http://localhost:8000/testps/val1/val2 | ConvertFrom-Json

      (($ret.arg1 -eq 'val1') -and ($ret.arg2 -eq 'val2')) | Should -Be $true
    }
  }
}

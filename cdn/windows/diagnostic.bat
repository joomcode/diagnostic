@echo on
setlocal EnableExtensions

:uniqLoop
set "out=%tmp%\joom~%RANDOM%"
if exist "%out%" goto :uniqLoop

set rep=%out%\report
set cat=143277b4305cfcb23573b35ba9d26448e71d8eb4_100_100.jpeg
mkdir "%out%"
mkdir "%rep%"
del /q "%out%\*"
del "%rep%\report.cab"
ver > "%out%\version.txt"
ipconfig /all > "%out%\ipconfig.txt"
route print > "%out%\route.txt"
nslookup myip.opendns.com. resolver1.opendns.com > "%out%\dns.myip.txt"
nslookup img.joomcdn.net > "%out%\dns.geo.txt"
nslookup img-alt.joomcdn.net > "%out%\dns.alt.txt"
nslookup img-amz.joomcdn.net > "%out%\dns.amz.txt"

tracert img-alt.joomcdn.net > "%out%\trace.alt.txt"
tracert img-amz.joomcdn.net > "%out%\trace.amz.txt"
certutil.exe -urlcache -f "http://iam.gcdn.co/info" "%out%\gcore.txt"
certutil.exe -urlcache -f "https://ifconfig.co/json" "%out%\external_ip.txt"
certutil.exe -urlcache -f "http://img-alt.joomcdn.net/%cat%" "%out%\cat_alt_http.jpg" > "%out%\cat_alt_http.out" 2> "%out%\cat_alt_http.err"
certutil.exe -urlcache -f "http://img-amz.joomcdn.net/%cat%" "%out%\cat_amz_http.jpg" > "%out%\cat_amz_http.out" 2> "%out%\cat_amz_http.err"
certutil.exe -urlcache -f "https://img-alt.joomcdn.net/%cat%" "%out%\cat_alt_https.jpg" > "%out%\cat_alt_https.out" 2> "%out%\cat_alt_https.err"
certutil.exe -urlcache -f "https://img-amz.joomcdn.net/%cat%" "%out%\cat_amz_https.jpg" > "%out%\cat_amz_https.out" 2> "%out%\cat_amz_https.err"

rem G-Core specific diagnostics
nslookup iam.gcdn.co > "%out%\dns.gcdn.txt"
nslookup -q=TXT o-o.myaddr.l.google.com > "%out%\dns.txt-myadd-isp.txt"
nslookup -q=TXT o-o.myaddr.l.google.com 8.8.8.8 > "%out%\dns.txt-myadd-gdns.txt"
nslookup -q=TXT dns-debug.d.gcdn.co > "%out%\dns.txt-gcdn.txt"
nslookup d.gcdn.co > "%out%\dns.gcdn-isp.txt"
nslookup d.gcdn.co 8.8.8.8 > "%out%\dns.gcdn-gdns.txt"
nslookup d.gcdn.co 92.223.100.100 > "%out%\dns.gcdn-gcore.txt"
tracert 92.223.100.200 > "%out%\trace.gcore.txt"
ping d.gcdn.co > "%out%\ping.gcdn.txt"
tracert d.gcdn.co > "%out%\trace.gcdn.txt"

cd "%out%"
dir /b /a-d > "%out%\files.lst"
makecab /d "DiskDirectoryTemplate=%rep%" /d "CabinetNameTemplate=report.cab" /f "%out%\files.lst"
explorer "%rep%"

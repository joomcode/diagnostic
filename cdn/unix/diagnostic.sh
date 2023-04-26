#!/bin/bash -x
set -x

CAT=143277b4305cfcb23573b35ba9d26448e71d8eb4_100_100.jpeg
CDN[0]=alt
CDN[1]=amz
CDN[2]=edg

# Begin diagnostic
OUT=`mktemp -d`
# OS version
uname -v | tee "$OUT/ver.txt"
# Network configuration
ifconfig | tee "$OUT/ifconfig.txt"

route > "$OUT/route.txt"
dig @resolver1.opendns.com myip.opendns.com. > "$OUT/dns.myip.txt"

# Geo DNS resolution
dig img.joomcdn.net | tee "$OUT/dns.geo.txt"

# G-Core diagnostic
curl --connect-timeout 30 http://iam.gcdn.co/info -v -o "$OUT/gcore.txt" 2> "$OUT/gcore.err"
curl --connect-timeout 30 http://iam.edgecenter.ru/info -v -o "$OUT/edgecenter.txt" 2> "$OUT/edgecenter.err"
curl --connect-timeout 30 https://ifconfig.co/json -v -o "$OUT/external_ip.txt" 2> "$OUT/external_ip.err"

# Check CDN providers
for i in ${CDN[*]}; do
  # DNS resolution
  dig img-$i.joomcdn.net | tee "$OUT/dns.$i.txt"
  # Try download image
  curl --connect-timeout 30 -v -o "$OUT/cat.${i}_https.jpg" "https://img-${i}.joomcdn.net/$CAT" > "$OUT/cat.${i}_https.txt" 2>&1
  curl --connect-timeout 30 -v -o "$OUT/cat.${i}_http.jpg"  "http://img-${i}.joomcdn.net/$CAT"  > "$OUT/cat.${i}_http.txt"  2>&1
  # Trace routing
  traceroute -w 1 img-$i.joomcdn.net | tee "$OUT/trace.$i.txt"
done

# G-Core specific diagnostics
dig iam.gcdn.co | tee "$OUT/dns.gcdn.txt"
dig o-o.myaddr.l.google.com TXT | tee "$OUT/dns.txt-myadd-isp.txt"
dig o-o.myaddr.l.google.com TXT @8.8.8.8 | tee "$OUT/dns.txt-myadd-gdns.txt"
dig dns-debug.d.gcdn.co TXT | tee "$OUT/dns.txt-gcdn.txt"
dig d.gcdn.co | tee "$OUT/dns.gcdn-isp.txt"
dig d.gcdn.co @8.8.8.8 | tee "$OUT/dns.gcdn-gdns.txt"
dig d.gcdn.co @92.223.100.100 | tee "$OUT/dns.gcdn-gcore.txt"
traceroute -w 1 92.223.100.200 | tee "$OUT/trace.gcore.txt"
ping -c5 d.gcdn.co | tee "$OUT/ping.gcdn.txt"
traceroute -w 1 d.gcdn.co | tee "$OUT/trace.gcdn.txt"

# Create archive with generated report
REP=`mktemp -d`
pushd "$OUT"
tar -czf "$REP/report.tgz" *
popd
rm -fR "$OUT"

# Open directory with generated report
if which xdg-open ; then
  xdg-open "$REP"
elif which open ; then
  open "$REP"
fi
echo "Report file $REP/report.tgz is successfully generated"

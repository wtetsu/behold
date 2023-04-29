#! /bin/sh

gaze="${1:-gaze}"

dir=$(cd $(dirname $0); pwd)
filedir=$dir/files

cd $dir
rm -f test.*.log

timeout -sKILL 5 ${gaze} -v files/*.* | tee test.log &

sleep 1.0
echo >> $filedir/hello.rb
sleep 0.2
echo >> $filedir/hello.go
sleep 0.2
echo >> $filedir/hello.py
sleep 0.2
echo >> $filedir/hello.rs
sleep 0.2
echo >> $filedir/hello.rb
sleep 0.2
echo >> $filedir/hello.go
sleep 0.2
echo >> $filedir/hello.py
sleep 0.2
echo >> $filedir/hello.rs
sleep 0.2
echo >> $filedir/hello.rb
sleep 0.2
echo >> $filedir/hello.py
sleep 0.2
echo >> $filedir/hello.rs
sleep 0.2

wait

num=`cat test.log | grep "hello, world!" | wc -l`

if [ $num -ne 11 ]; then
  echo "Failed:${num}"
  exit 1
fi

git checkout -- files

echo "OK"
exit 0

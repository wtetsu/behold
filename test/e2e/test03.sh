#! /bin/sh

gaze="${1:-gaze}"

dir=$(cd $(dirname $0); pwd)
filedir=$dir/files

cd $dir
rm -f test.*.log

cp $filedir/hello.py "$filedir/he'llo.py"
cp $filedir/hello.py "$filedir/he&llo.py"
cp $filedir/hello.py "$filedir/he llo.py"
cp $filedir/hello.py "$filedir/he(llo.py"

timeout -sKILL 5 ${gaze} -v files/*.* | tee test.log &

sleep 1.0
echo >> "$filedir/he'llo.py"
sleep 0.2
echo >> "$filedir/he&llo.py"
sleep 0.2
echo >> "$filedir/he llo.py"
sleep 0.2
echo >> "$filedir/he(llo.py"
sleep 0.2
echo >> "$filedir/he'llo.py"
sleep 0.2
echo >> "$filedir/he&llo.py"
sleep 0.2
echo >> "$filedir/he llo.py"
sleep 0.2
echo >> "$filedir/he'llo.py"
sleep 0.2
echo >> "$filedir/he(llo.py"
sleep 0.2
echo >> "$filedir/he&llo.py"
sleep 0.2
echo >> "$filedir/he llo.py"
sleep 0.2
echo >> "$filedir/he(llo.py"

wait

rm "$filedir/he'llo.py"
rm "$filedir/he&llo.py"
rm "$filedir/he llo.py"
rm "$filedir/he(llo.py"

num=`cat test.log | grep "hello, world!" | wc -l`

if [ $num -ne 9 ]; then
  echo "Failed:${num}"
  exit 1
fi

echo "OK"
exit 0

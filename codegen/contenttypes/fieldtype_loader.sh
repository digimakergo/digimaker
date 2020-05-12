#! /bin/bash
package=$1
if test -z "$package"
then
  echo "Usage: fieldtype.sh <project package name>"
  exit
fi
dir=`dirname $0`
if [ ! -d "${dir}/temp" ]
then
echo `mkdir ${dir}/temp`
fi
echo "Generate loader for project field types."
echo  "package temp
import _ \"${package}/fieldtype\"" > ./$dir/temp/project.go

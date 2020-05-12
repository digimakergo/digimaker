package=$1
if test -z "$package"
then
  echo "Usage: go to <GOPATH> and run <dir>/gen.sh <project package name>"
  exit
fi

dir=`dirname $0`
path=$1
sh $dir/fieldtype_loader.sh $path
go run $dir/gen.go src/$path
if [ -f "$dir/temp/project.go" ]; then
  rm $dir/temp/project.go
fi

#!/bin/bash

program="api-manager"
program_version=`cat VERSION`

function gobuild() {
	build_time=`date`
	author=`whoami`
	compiler_version=`go version`

	go build -v -ldflags \
		"-X 'main.ProgramVersion=${program_version}' \
		 -X 'main.CompilerVersion=${compiler_version}' \
		 -X 'main.BuildTime=${build_time}' \
		 -X 'main.Author=${author}'" \
		-o ${program}
}

function webbuild() {
	echo Todo
}

function checkbuild() {
	if [ ! -e "./${program}" ]; then
		build
		if [ ! -e "./${program}" ]; then
			echo "build ${program} failed!!!"
			exit 1
		fi
	fi
}

function pkg() {
	checkbuild
	if [[ $? -ne 0 ]]; then
		exit 1
	fi

	pkg_name="${program}_${program_version}.tar.gz"

	tar -zcf $pkg_name ${program} --exclude=logs/*.log \
		logs \
		etc \
		VERSION \
		project.sh \
		public
}

function publish() {
	checkbuild
	if [[ $? -ne 0 ]]; then
		exit 1
	fi
	sudo mkdir -p /data/${program}/
	# supervisorctl stop ${program}
	sudo rm -f /data/${program}/${program}
	sudo cp ${program} /data/${program}/${program}_${program_version}
	sudo ln -s /data/${program}/${program}_${program_version} /data/${program}/${program}
	# supervisorctl start ${program}
}

if [[ $1 = "pkg" ]]; then
	pkg
elif [[ $1 = "pub" ]]; then
	publish
elif [[ $1 = "web" ]]; then
	webbuild
else
	gobuild
fi

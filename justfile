alias bc := bash-create
alias clear := clearall
alias k := kill

pkg := "github.com/opencontainers/runc"

clearall:
   @sudo runc list | awk '$1 !~ /ID/ {printf "%s\n", $1}' | xargs  -I{} sh -c 'sudo runc kill {}; sudo runc delete {}'

install:
    @make
    @sudo -E env "PATH=$PATH" make install

test:
    go test -v -count=1 {{pkg}}/comm
    

bash-create:
    @sudo runc --debug create  -b ${HOME}/oci-images/bash-bundle bash

kill:
    kill -9 $(pgrep -f systemr)

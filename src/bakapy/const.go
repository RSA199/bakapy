package bakapy

import (
	"text/template"
)

// Waiting for client authentication
const STORAGE_AUTH_TIMEOUT = 30 // seconds
const STORAGE_TASK_ID_LEN = 36
const STORAGE_READ_BUFSIZE = 4096

// Length of filename length header
const STORAGE_FILENAME_LEN_LEN = 4

// Storage connection states
const (
	STATE_WAIT_TASK_ID = iota
	STATE_WAIT_FILENAME
	STATE_WAIT_DATA
	STATE_RECEIVING
	STATE_END
)

const JOB_FINISH = "_@!_JOB_FINISH_!@_"

var JOB_TEMPLATE = template.Must(template.New("job").Parse(`
##
# Common header
##
set -e

FILENAME_LEN_LEN='{{.FILENAME_LEN_LEN}}'
TO_HOST='{{.ToHost}}'
TO_PORT='{{.ToPort}}'
TASK_ID='{{.Meta.TaskId}}'
TASK_NAME='{{.Meta.JobName}}'

_send_file(){
    name="$1"

    exec 3<>/dev/tcp/${TO_HOST}/${TO_PORT}
    echo -n ${TASK_ID}$(printf "%0${FILENAME_LEN_LEN}d" ${#name})${name} >&3
    cat - >&3
    exec 3>&-
}

_finish(){
    echo > /dev/null
}

_fail(){
    test ! -z "$1" && echo "command failed at line $1" >&2
    exit 1
}

function error_exit {
    if [ "$?" != "0" ]; then
        echo "FAILED"
        _fail
        exit 1
    fi
}

trap '_fail ${LINENO}' ERR

##
# Command
##
`))

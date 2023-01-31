export POSH_THEME=::CONFIG::
export POWERLINE_COMMAND="oh-my-posh"
export POSH_PID=$$
export CONDA_PROMPT_MODIFIER=false
omp_start_time=""

# start timer on command start
PS0='${omp_start_time:0:$((omp_start_time="$(_omp_start_timer)",0))}'
# set secondary prompt
PS2="$(::OMP:: print secondary --config="$POSH_THEME" --shell=bash --shell-version="$BASH_VERSION")"

function _set_posh_cursor_position() {
    # not supported in Midnight Commander
    # see https://github.com/JanDeDobbeleer/oh-my-posh/issues/3415
    if [[ -v MC_SID ]];then
        return
    fi
    local pos
    echo -ne "\033[6n"            # ask the terminal for the position
    read -s -d\[ garbage          # discard the first part of the response
    read -s -d R pos              # store the position in bash variable 'pos'
    export POSH_CURSOR_LINE=${pos%;*}
    export POSH_CURSOR_COLUMN=${pos#*;}
}

function _omp_start_timer() {
    ::OMP:: get millis
}

# template function for context loading
function set_poshcontext() {
    return
}

function _omp_hook() {
    local ret=$?
    local omp_stack_count=$((${#DIRSTACK[@]} - 1))
    local omp_elapsed=-1
    if [[ -n "$omp_start_time" ]]; then
        local omp_now=$(::OMP:: get millis --shell=bash)
        omp_elapsed=$((omp_now-omp_start_time))
        omp_start_time=""
    fi
    set_poshcontext
    _set_posh_cursor_position
    PS1="$(::OMP:: print primary --config="$POSH_THEME" --shell=bash --shell-version="$BASH_VERSION" --error="$ret" --execution-time="$omp_elapsed" --stack-count="$omp_stack_count" | tr -d '\0')"
    return $ret
}

if [ "$TERM" != "linux" ] && [ -x "$(command -v ::OMP::)" ] && ! [[ "$PROMPT_COMMAND" =~ "_omp_hook" ]]; then
    PROMPT_COMMAND="_omp_hook; $PROMPT_COMMAND"
fi

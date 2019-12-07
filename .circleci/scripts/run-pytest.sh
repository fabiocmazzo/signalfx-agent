#!/bin/bash

set -eo pipefail

[ -f ~/.skip ] && echo "Found ~/.skip, skipping tests!" && exit 0

[ -n "$TESTS_DIR" ] || (echo "TESTS_DIR not defined!" && exit 1)
[ -d "$TESTS_DIR" ] || (echo "Directory '$TESTS_DIR' not found!" && exit 1)

[ -f $BASH_ENV ] && source $BASH_ENV
export BUNDLE_DIR=${BUNDLE_DIR:-$(pwd)/bundle}
export AGENT_BIN=${AGENT_BIN:-${BUNDLE_DIR}/bin/signalfx-agent}
export TEST_SERVICES_DIR=${TEST_SERVICES_DIR:-$(pwd)/test-services}

mkdir -p /tmp/scratch
mkdir -p ~/testresults

if [[ $CIRCLE_NODE_TOTAL -gt 1 && -n "$MARKERS" && $SPLIT -eq 1 ]]; then
    # Collect tests based on MARKERS and split them for parallelism
    TESTS=$(pytest --collect-only -m "$MARKERS" $TESTS_DIR | grep '<Module.*>' | cut -d "'" -f2 | \
        circleci tests split --split-by=timings --total=$CIRCLE_NODE_TOTAL --index=$CIRCLE_NODE_INDEX)
    [ -n "$TESTS" ] || (echo "No test files found after splitting based on '$MARKERS' marker(s)!" && exit 1)
else
    TESTS=$TESTS_DIR
fi

PYTEST_PATH="pytest"
if [ $WITH_SUDO -eq 1 ]; then
    PYTEST_PATH="sudo -E $PYENV_ROOT/shims/pytest"
fi
FIRSTRUN_OPTIONS="--verbose --junitxml=~/testresults/results.xml --html=~/testresults/results.html --self-contained-html"
RERUN_OPTIONS="--last-failed --verbose --junitxml=~/testresults/rerun-results.xml --html=~/testresults/rerun-results.html --self-contained-html"
WORKERS=${WORKERS-}
if [ -n "$WORKERS" ]; then
    WORKERS="-n $WORKERS"
fi

set -x
if [ -n "$MARKERS" ]; then
    $PYTEST_PATH -m "$MARKERS" $WORKERS $PYTEST_OPTIONS $FIRSTRUN_OPTIONS $TESTS || \
        $PYTEST_PATH -m "$MARKERS" $PYTEST_OPTIONS $RERUN_OPTIONS $TESTS
else
    $PYTEST_PATH $WORKERS $PYTEST_OPTIONS $FIRSTRUN_OPTIONS $TESTS || \
        $PYTEST_PATH $PYTEST_OPTIONS $RERUN_OPTIONS $TESTS
fi

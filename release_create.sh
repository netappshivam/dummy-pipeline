#!/bin/bash
set -e

SPRINT='2510'
echo "Creating ${SPRINT} Branch"

new_release_branch_name=$1

LAST_TAG=$new_release_branch_name

if [ -n "$LAST_TAG" ]; then
    echo "Latest and gretest tag is : ${LAST_TAG}"
    REMOTE_BRANCH=$(git ls-remote --heads origin ${SPRINT})
    if [ -n "${REMOTE_BRANCH}" ]; then
            echo "Remote branch ${SPRINT} is already exists hence skipping branch creation!!!"
    else
        echo "Creating the branch - ${SPRINT}"
        git checkout -b ${SPRINT} tags/${LAST_TAG}
        git push origin ${SPRINT}
    fi
else
    echo "Since there is no tag created for the sprint - ${SPRINT}, cutting the release branch based out of master branch"
    echo "Creating the branch - ${SPRINT}"
    git checkout -b ${SPRINT}
    git push origin ${SPRINT}
    #REMOTE_BRANCH=$(git branch -r --list | grep origin/${SPRINT})
    #if [ -n "${REMOTE_BRANCH}" ]; then
            #echo "Remote branch ${SPRINT} is already exists hence skipping branch creation!!!"
    #else
        #echo "Creating the branch - ${SPRINT}"
        #git checkout -b ${SPRINT}
        #git push origin ${SPRINT}
    #fi
fi
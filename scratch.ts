// problem: we have a set of data and need to get a subset and return a boolean based on the presence/lack of presence of
//
// FAIL == warn cannot override linter
// PASS == no override found

// checkDiff() looping through just the diffs
// if we find '#nosec', early exit FAIL
// if we don't find 'eslint-disable', PASS (exit after all diff chunks pass)

// assume we can split into lines?
// eslint-disable-next-line adf-asdgasdg, asgerhijerh-iedfjer

// checkLine() looping through lines for each diff
// look for disable marker. when found:
//     - if no rule specified, early exit FAIL
// - if any rule not in acceptable list, early exit FAIL
// - otherwise (if all rules acceptable), PASS (this chunk; cannot exit unless final chunk)

import { danger, warn, fail } from 'danger';

const bypassingLinterChecks = async () => {
    const allFiles = danger.git.modified_files.concat(danger.git.created_files);
    const diffsByFile = await Promise.all(allFiles.map((f) => danger.git.diffForFile(f)));
    const showDanger = checkPRHasUnpermittedLinterOverride(diffsByFile);
    if (showDanger) {
        // throw dangerjs warning
    }
}

function checkPRHasUnpermittedLinterOverride(dangerJSDiffCollection) {
    let hasBadOverride = false
    for (let d in dangerJSDiffCollection) {
        const diffFile = dangerJSDiffCollection[d]
        const diff = diffFile.diff
        if (diffContainsNosec(diff)) {
            hasBadOverride = true;
            break;
        }
        if (!diffContainsEslint(diff)) {
            continue;
        }

        // magic split into lines happen here
        const lines = magicSplit(diff)
        for (let l in lines) {
            const line = lines[l]
            if (diffContainsEslint(line)) {
                // check for comment marker (// or /*)
                // then parse line after comment chars
                // eg line: 'const whatever = something() // eslint-disable-line'
                let lineParts = line.split('//')
                if (lineParts.length === 1) { // this is where marker isn't found
                    lineParts = line.split('/*')
                    if (lineParts.length === 1) {
                        throw new Error('uhhhh, how did we find eslint disable but no // or /*')
                    }
                }

                // eg lineParts: ['const whatever = something()', 'eslint-disable-line']
                if (doesLineHaveUnpermittedOverride(lineParts[1])) {
                    // fail because user shouldn't add new overrides without security / moose approval
                    hasBadOverride = true;
                } // else continue
            }
        }
        return hasBadOverride;
    }
}
function diffContainsNosec(diffForFile) {
    return !!diffForFile.includes('#nosec');
}

function diffContainsEslint(diffForFile) {
    return !!diffForFile.includes('eslint-disable');
}


// comment characters location (where // or /* is in line string)
function doesLineHaveUnpermittedOverride(line) {
    const okBypassRules = [
        'no-underscore-dangle',
        'prefer-object-spread',
        'object-shorthand',
        'camelcase',
        'react/jsx-props-no-spreading',
        'react/destructuring-assignment',
        'react/forbid-prop-types',
        'react/prefer-stateless-function',
        'react/sort-comp',
        'import/no-extraneous-dependencies',
        'import/order',
        'import/prefer-default-export',
        'import/no-named-as-default',
    ];
    // eslint-disable-next-line no-jsx, no-default
    // use regex to grab any word-like groups
    const magicWordGroupsRegex = /??????/;
    const matches = line.matches(magicWordGroupsRegex);

    // matches: ['eslint-disable-next-line', 'no-jsx', 'no-default']
    if (matches[0] === 'eslint-disable') {
        // fail because don't disable whole file please!
    }

    if (matches.length === 1) {
        // fail because please specify rule
    }

    // rules: ['no-jsx', 'no-default']
    let rules = matches.slice(1);
    let hasUnpermittedOverride = false;
    for (let r in rules) {
        const rule = rules[r]
        if (!okBypassRules.includes(rule)) {
            hasUnpermittedOverride = true
            break
        }
    }
    return hasUnpermittedOverride
}

console.log(bypassingLinterChecks());



















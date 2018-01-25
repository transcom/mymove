import { SHOW_ISSUES } from 'actions/index.js';

function showIssues(state = { issues: null, hasError: false }, action) {
  switch (action.type) {
    case SHOW_ISSUES:
      return { issues: null, hasError: false };
    default:
      return state;
  }
}

export default showIssues;

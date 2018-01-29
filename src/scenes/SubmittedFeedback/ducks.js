// Action

export const SHOW_ISSUES = 'SHOW_ISSUES';

export function showIssues() {
  return { type: SHOW_ISSUES };
}

// Action Creator

// Reducer
function showIssues(state = { issues: null, hasError: false }, action) {
  switch (action.type) {
    case SHOW_ISSUES:
      return loadIssues();
    default:
      return state;
  }
}

export default showIssues;

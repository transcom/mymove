import React from 'react';
export const defaultMyMoveContext = {
  siteName: 'my.move.mil',
  showLoginWarning: true,
  flags: {},
};
export const defaultOfficeContext = {
  siteName: 'office.move.mil',
  showLoginWarning: false,
  flags: {},
};
export const defaultTspContext = {
  siteName: 'tsp.move.mil',
  showLoginWarning: false,
};
export const defaultAdminContext = {
  siteName: 'admin.move.mil',
  showLoginWarning: false,
};
export const AppContext = React.createContext(defaultMyMoveContext);

export function withContext(Component) {
  // ...and returns another component...
  return function ContextualComponent(props) {
    // ... and renders the wrapped component with the context theme!
    // Notice that we pass through any additional props as well
    return <AppContext.Consumer>{context => <Component {...props} context={context} />}</AppContext.Consumer>;
  };
}

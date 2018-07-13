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
export const AppContext = React.createContext(defaultMyMoveContext);

import React from 'react';
export const myMoveContext = {
  siteName: 'my.move.mil',
  showLoginWarning: true,
};
export const officeContext = {
  siteName: 'office.move.mil',
  showLoginWarning: false,
};
export const tspContext = {
  siteName: 'tsp.move.mil',
  showLoginWarning: false,
};
export const AppContext = React.createContext(myMoveContext);

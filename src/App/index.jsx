import React from 'react';
import { Provider } from 'react-redux';

import store from 'shared/store';
import './index.css';

import Loadable from 'react-loadable';

const Office = Loadable({
  loader: () => import('scenes/Office'),
  loading: () => <div>Loading...</div>,
});

const MyMove = Loadable({
  loader: () => import('scenes/MyMove'),
  loading: () => <div>Loading...</div>,
});

const officeSiteNames = [
  'bolocal',
  'office.staging.dp3.us',
  'office.prod.dp3.us',
];
const hostname = window && window.location && window.location.hostname;
console.log(hostname);
const isOfficeSite = officeSiteNames.includes(hostname);
const App = () => (
  <Provider store={store}>{!isOfficeSite ? <MyMove /> : <Office />}</Provider>
);

export default App;

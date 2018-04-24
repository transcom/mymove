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

const hostname = window && window.location && window.location.hostname;
const isOfficeSite = hostname === 'bolocal' || hostname.startsWith('office');
const App = () => (
  <Provider store={store}>{!isOfficeSite ? <MyMove /> : <Office />}</Provider>
);

export default App;

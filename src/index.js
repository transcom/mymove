// This must be the first line in src/index.js
import 'react-app-polyfill/ie11';
import 'core-js/stable';
import 'regenerator-runtime/runtime';
import React from 'react';
import ReactDOM from 'react-dom';
import { isSystemAdminSite } from 'shared/constants.js';
import 'uswds';

import App from './App';
import registerServiceWorker from './registerServiceWorker';
import './index.css';

if (!isSystemAdminSite) {
  require('../node_modules/uswds/dist/css/uswds.css');
}
ReactDOM.render(<App />, document.getElementById('root'));
registerServiceWorker();

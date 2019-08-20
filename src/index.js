// This must be the first line in src/index.js
import 'react-app-polyfill/ie11';
import 'core-js/stable';
import 'regenerator-runtime/runtime';
import React from 'react';
import ReactDOM from 'react-dom';
import 'uswds';
import '../node_modules/uswds/dist/css/uswds.css';

import App from './App';
import registerServiceWorker from './registerServiceWorker';
import './index.css';
ReactDOM.render(<App />, document.getElementById('root'));
registerServiceWorker();

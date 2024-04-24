// This must be the first line in src/index.js
import 'react-app-polyfill/ie11';
import 'core-js/stable';
import 'regenerator-runtime/runtime';
import React from 'react';
import ReactDOM from 'react-dom';

import { configureGlobalLogger } from './utils/milmoveLog';
import App from './App';
// import registerServiceWorker from './registerServiceWorker';
import './index.scss';

// configure the global logger once
configureGlobalLogger();

ReactDOM.render(<App />, document.getElementById('root'));

// MilMove does not have or use a service worker as of 2023-09-07
// so disable this to prevent logging errors
// registerServiceWorker();

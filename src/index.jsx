// This must be the first line in src/index.js
import 'react-app-polyfill/ie11';
import 'core-js/stable';
import 'regenerator-runtime/runtime';
import React from 'react';
import ReactDOM from 'react-dom/client';

import { configureGlobalLogger } from './utils/milmoveLog';
import App from './App';
// import registerServiceWorker from './registerServiceWorker';
import './index.scss';

// configure the global logger once
configureGlobalLogger();

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(<App />);

// so disable this to prevent logging errors
// registerServiceWorker();

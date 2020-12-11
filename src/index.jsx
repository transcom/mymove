// This must be the first line in src/index.js
import 'react-app-polyfill/ie11';
import 'core-js/stable';
import 'regenerator-runtime/runtime';
import React from 'react';
import ReactDOM from 'react-dom';

import App from './App';
import registerServiceWorker from './registerServiceWorker';
import './index.scss';

// eslint-disable-next-line no-console
console.log('RENDER APP');

ReactDOM.render(<App />, document.getElementById('root'));
registerServiceWorker();

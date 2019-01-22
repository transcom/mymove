// This must be the first line in src/index.js
import 'react-app-polyfill/ie11';
import 'babel-polyfill'; //todo: is this still needed?
import React from 'react';
import ReactDOM from 'react-dom';
import 'uswds';

import App from './App';
import registerServiceWorker from './registerServiceWorker';
import './index.css';
import '../node_modules/uswds/dist/css/uswds.css';

ReactDOM.render(<App />, document.getElementById('root'));
registerServiceWorker();

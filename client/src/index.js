import React from 'react';
import ReactDOM from 'react-dom';
import 'uswds';
import App from './App';
import registerServiceWorker from './registerServiceWorker';
import 'index.css';
import '../node_modules/uswds/dist/css/uswds.css';

ReactDOM.render(<App />, document.getElementById('root'));
registerServiceWorker();

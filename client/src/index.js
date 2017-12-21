import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import 'uswds';
import '../node_modules/uswds/dist/css/uswds.css';
import App from './App';
import registerServiceWorker from './registerServiceWorker';

ReactDOM.render(<App />, document.getElementById('root'));
registerServiceWorker();

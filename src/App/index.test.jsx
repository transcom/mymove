import React from 'react';
import ReactDOM from 'react-dom';

import App from './index';

it('renders without crashing', () => {
  const div = document.createElement('div');
  ReactDOM.render(<App />, div);
  // Until we come up with a better plan, this prevents our tests from crashing.
  // 1. The Feedback component is mounted at / and so mounted when we mount <App>
  // 2. The Feedback component uses JSONSchemaForm which means it attempts to
  //    load swagger.yaml when it is mounted
  // 3. This attempt makes our test asynchonous, which without the proper handling causes
  //    the test runner to crash. Immediately unmounting the component prevents the crash
  //    and still does the bare minimum of confirming that the whole app mounts without error.
  ReactDOM.unmountComponentAtNode(div);
});

// todo: add tests for routing

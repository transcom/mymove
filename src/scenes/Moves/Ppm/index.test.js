import React from 'react';
import ReactDOM from 'react-dom';
import { Provider } from 'react-redux';
import PpmSize from '.';
import store from 'shared/store';

const dummyFunc = () => {};
const hasSubmitError = false;
const hasSubmitSuccess = false;
const currentPpm = null;

it('renders without crashing', () => {
  const div = document.createElement('div');
  ReactDOM.render(
    <Provider store={store}>
      <PpmSize
        createPpm={dummyFunc}
        currentPpm={currentPpm}
        hasSubmitSuccess={hasSubmitSuccess}
        hasSubmitError={hasSubmitError}
        match={{ match: 'match' }}
      />
    </Provider>,
    div,
  );
});

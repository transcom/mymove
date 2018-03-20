import React from 'react';
import ReactDOM from 'react-dom';
import { Provider } from 'react-redux';
import { Uploader } from '.';
import store from 'shared/store';

const dummyFunc = () => {};
const hasSubmitError = false;
const hasSubmitSuccess = false;

it('renders without crashing', () => {
  const div = document.createElement('div');
  ReactDOM.render(
    <Provider store={store}>
      <Uploader
        hasSubmitSuccess={hasSubmitSuccess}
        hasSubmitError={hasSubmitError}
        confirmationText=""
        createDocument={dummyFunc}
      />
    </Provider>,
    div,
  );
});

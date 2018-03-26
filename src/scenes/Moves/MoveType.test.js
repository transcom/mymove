import React from 'react';
import ReactDOM from 'react-dom';
import { Provider } from 'react-redux';
import MoveType from './MoveType';
import store from 'shared/store';

const hasSubmitError = false;
const hasSubmitSuccess = false;
const currentMove = null;

it('renders without crashing', () => {
  const div = document.createElement('div');
  ReactDOM.render(
    <Provider store={store}>
      <MoveType
        currentMove={currentMove}
        hasSubmitSuccess={hasSubmitSuccess}
        hasSubmitError={hasSubmitError}
        match={{ match: 'match' }}
      />
    </Provider>,
    div,
  );
});

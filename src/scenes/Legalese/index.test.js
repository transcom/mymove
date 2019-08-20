import React from 'react';
import ReactDOM from 'react-dom';
import { Provider } from 'react-redux';
import SignedCertification from '.';
import store from 'shared/store';

const schema = {};
const uiSchema = {};
const match = { params: { moveId: 'someID' } };
const hasSchemaError = false;
const hasSubmitError = false;
const hasSubmitSuccess = false;

it('renders without crashing', () => {
  const div = document.createElement('div');
  ReactDOM.render(
    <Provider store={store}>
      <SignedCertification
        hasSchemaError={hasSchemaError}
        hasSubmitSuccess={hasSubmitSuccess}
        hasSubmitError={hasSubmitError}
        schema={schema}
        uiSchema={uiSchema}
        confirmationText=""
        match={match}
      />
    </Provider>,
    div,
  );
});

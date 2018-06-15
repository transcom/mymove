import React from 'react';
import ReactDOM from 'react-dom';
import { Provider } from 'react-redux';
import { Feedback } from '.';
import store from 'shared/store';

const dummyFunc = () => {};
const schema = {};
const uiSchema = {};
const hasSchemaError = false;
const hasSubmitError = false;
const hasSubmitSuccess = false;

it('renders without crashing', () => {
  const div = document.createElement('div');
  ReactDOM.render(
    <Provider store={store}>
      <Feedback
        hasSchemaError={hasSchemaError}
        hasSucceeded={hasSubmitSuccess}
        hasErrored={hasSubmitError}
        schema={schema}
        uiSchema={uiSchema}
        createIssue={dummyFunc}
      />
    </Provider>,
    div,
  );
});

import React from 'react';
import ReactDOM from 'react-dom';

import { Provider } from 'react-redux';
import { reducer as formReducer, reduxForm } from 'redux-form';
import { createStore, combineReducers } from 'redux';

import JsonSchemaField from './JsonSchemaField';

import { shallow, mount, render } from 'enzyme';

describe('SchemaField tests', () => {
  const formHolster = field => {
    return props => {
      return <form className="default">{field}</form>;
    };
  };

  const reduxHolster = form => {
    return reduxForm({ form: 'holster' })(formHolster(form));
  };

  const mountField = (store, field) => {
    let Holster = reduxHolster(field);
    let subject = mount(
      <Provider store={store}>
        <Holster />
      </Provider>,
    );
    return subject;
  };

  const testField = (field, tests) => {
    const store = createStore(combineReducers({ form: formReducer }));
    const testField = JsonSchemaField.createSchemaField(
      'test_field',
      field,
      '',
    );
    let subject = mountField(store, testField);

    tests.forEach(testCase => {
      const [testValue, expectedValue, expectedError] = testCase;

      it(`${testValue} results in ${
        !expectedError ? 'no error' : `error: ${expectedError}`
      }`, () => {
        const input = subject.find('input').first();
        input.simulate('change', { target: { value: testValue } });

        let storeData = store.getState().form.holster;
        console.log(storeData);
        let values = storeData.values;
        let errors = storeData.syncErrors;

        if (expectedValue != null) {
          expect(values).not.toBeUndefined();
          expect(values.test_field).toEqual(expectedValue);
        } else {
          expect(values).toBeUndefined();
        }

        if (expectedError) {
          expect(errors).not.toBeUndefined();
          expect(errors.test_field).toEqual(expectedError);
        } else {
          expect(errors).toBeUndefined();
        }
      });
    });
  };

  // First test. A number field?
  describe('number field', () => {
    describe('integer with limits', () => {
      const numberField = {
        type: 'integer',
        example: 10,
        maximum: 11,
        minimum: 0,
        'x-nullable': true,
        title: 'Annoyance Level',
      };

      const numberTests = [
        ['11', 11, null],
        ['1', 1, null],
        ['0', 0, null],
        ['2.', 2, null],
        ['2a', 2, null],
        ['', null, null],
        ['a2', 'a2', 'Must be an integer'],
        ['22.2', 22.2, 'Must be 11 or less'],
        ['100', 100, 'Must be 11 or less'],
        ['1.3', 1.3, 'Must be an integer'],
        ['-1', -1, 'Must be 0 or more'],
      ];

      testField(numberField, numberTests);
    });

    describe('integer without limits', () => {
      const numberField = {
        type: 'integer',
        example: 10,
        'x-nullable': true,
        title: 'Annoyance Level',
      };

      const numberTests = [
        ['11', 11, null],
        ['1', 1, null],
        ['0', 0, null],
        ['2.', 2, null],
        ['2a', 2, null],
        ['-1', -1, null],
        ['100', 100, null],
        ['', null, null],
        ['a2', 'a2', 'Must be an integer'],
        ['22.2', 22.2, 'Must be an integer'],
        ['1.3', 1.3, 'Must be an integer'],
      ];

      testField(numberField, numberTests);
    });

    describe('number with limits', () => {
      const numberField = {
        type: 'number',
        example: 10,
        'x-nullable': true,
        title: 'Annoyance Level',
      };

      const numberTests = [
        ['11', 11, null],
        ['1', 1, null],
        ['0', 0, null],
        ['2.', 2, null],
        ['2a', 2, null],
        ['-1', -1, null],
        ['100', 100, null],
        ['', null, null],
        ['a2', 'a2', 'Must be a number'],
        ['22.2', 22.2, null],
        ['1.3', 1.3, null],
      ];

      testField(numberField, numberTests);
    });
  });
});

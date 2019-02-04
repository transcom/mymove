import React from 'react';

import { Provider } from 'react-redux';
import { reducer as formReducer, reduxForm } from 'redux-form';
import { createStore, combineReducers } from 'redux';

import JsonSchemaField from './JsonSchemaField';
import { recursivelyValidateRequiredFields } from './index';

import { mount } from 'enzyme';

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
    const testField = JsonSchemaField.createSchemaField('test_field', field, '');
    let subject = mountField(store, testField);

    tests.forEach(testCase => {
      const [testValue, expectedValue, expectedError] = testCase;

      it(`${testValue} results in ${!expectedError ? 'no error' : `error: ${expectedError}`}`, () => {
        let input = subject.find('input');
        if (input.length === 0) {
          input = subject.find('textarea');
        }
        input.simulate('change', { target: { value: testValue } });

        let storeData = store.getState().form.holster;
        let values = storeData.values;
        let errors = storeData.syncErrors;

        if (expectedValue === undefined) {
          expect(values).toBeUndefined();
        } else if (expectedValue === null) {
          expect(values).not.toBeUndefined();
          expect(values.test_field).toBeNull();
        } else {
          expect(values).not.toBeUndefined();
          expect(values.test_field).toEqual(expectedValue);
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

  ['string', 'textarea'].forEach(fieldType => {
    describe(fieldType + ' text field', () => {
      describe('with limits', () => {
        const textFieldWithLimits = {
          type: 'string',
          minLength: 2,
          maxLength: 20,
          example: 'I am the very model of a modern Major General',
          'x-nullable': true,
          title: 'Introduction',
        };

        if (fieldType === 'textarea') {
          textFieldWithLimits['format'] = 'textarea';
        }

        const stringTests = [
          ['Hello', 'Hello', null],
          ['1', '1', 'Must be at least 2 characters long.'],
          ['🌝🤩🍕❤️', '🌝🤩🍕❤️', null],
          [
            'This is the song that never ends, it just goes on and on my friends',
            'This is the song that never ends, it just goes on and on my friends',
            'Cannot exceed 20 characters.',
          ],
        ];

        testField(textFieldWithLimits, stringTests);
      });

      describe('without limits', () => {
        const textField = {
          type: 'string',
          example: 'I am the very model of a modern Major General',
          'x-nullable': true,
          title: 'Introduction',
        };

        if (fieldType === 'textarea') {
          textField['format'] = 'textarea';
        }

        const stringTests = [
          ['Hello', 'Hello', null],
          ['1', '1', null],
          ['🌝🤩🍕❤️', '🌝🤩🍕❤️', null],
          [
            'This is the song that never ends, it just goes on and on my friends',
            'This is the song that never ends, it just goes on and on my friends',
            null,
          ],
        ];

        testField(textField, stringTests);
      });
    });
  });

  describe('telephone field', () => {
    const telephoneField = {
      type: 'string',
      format: 'telephone',
      pattern: /^[2-9]\d{2}-\d{3}-\d{4}$/,
      example: '615-222-3323',
      'x-nullable': true,
      title: 'Telephone No.',
    };

    const stringTests = [
      ['615-222-3323', '615-222-3323', null],
      ['6152223323', '615-222-3323', null],
      ['615-222-332sdfsdfsd3', '615-222-3323', null],
      ['615-222-332', '615-222-332', 'Number must have 10 digits and a valid area code.'],
      ['615-222-33233', '615-222-3323', null],
      ['115-222-33233', '115-222-3323', 'Number must have 10 digits and a valid area code.'],
    ];

    testField(telephoneField, stringTests);
  });

  describe('email field', () => {
    const emailField = {
      type: 'string',
      format: 'x-email',
      pattern: '^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$', // eslint-disable-line
      example: 'john_bob@example.com',
      'x-nullable': true,
      title: 'Personal Email Address',
    };

    const stringTests = [
      ['john_bob@example.com', 'john_bob@example.com', null],
      ['macrae.linton@gmail.com', 'macrae.linton@gmail.com', null],
      ['john_BOB@examPLE.co.uk', 'john_BOB@examPLE.co.uk', null],
      ['john_bot', 'john_bot', 'Must be a valid email address'],
      ['john_bot@foo', 'john_bot@foo', 'Must be a valid email address'],
      ['john_bot.com', 'john_bot.com', 'Must be a valid email address'],
    ];

    testField(emailField, stringTests);
  });

  describe('zip field', () => {
    const zipField = {
      type: 'string',
      format: 'zip',
      pattern: /^(\d{5}([-]\d{4})?)$/, // eslint-disable-line
      example: '61522-3323',
      'x-nullable': true,
      title: 'ZIP Code',
    };

    const stringTests = [
      ['61522-3323', '61522-3323', null],
      ['61522', '61522', null],
      ['615223323', '61522-3323', null],
      ['615-22-332sdfsdfsd3', '61522-3323', null],
      ['615-22-332', '61522-332', 'Zip code must have 5 or 9 digits.'],
      ['615-22-33233', '61522-3323', null],
    ];

    testField(zipField, stringTests);
  });

  describe('base quantity field', () => {
    const baseQuantityField = {
      type: 'integer',
      format: 'basequantity',
      example: 120000,
    };

    const baseQuantityTests = [
      ['1.0000', '1.0000', null],
      ['121.9548', '121.9548', null],
      ['12.12345', '12.1234', null],
      ['12.12', '12.12', null],
      ['12.abcd', '12.', null],
      ['1.1..1', '1.11', null],
    ];

    testField(baseQuantityField, baseQuantityTests);
  });

  describe('dimensions field', () => {
    const dimensionsField = {
      type: 'integer',
      format: 'dimension',
      example: 1200,
    };

    const dimensionTests = [
      ['1.00', '1.00', null],
      ['121.95', '121.95', null],
      ['12.12345', '12.12', null],
      ['12.1', '12.1', null],
      ['12.abcd', '12.', null],
      ['1..1.1', '1.11', null],
    ];

    testField(dimensionsField, dimensionTests);
  });
});

describe('fields required tests', () => {
  const testSchema = {
    title: 'A registration form',
    description: 'A simple form example.',
    type: 'object',
    required: ['firstName', 'lastName'],
    properties: {
      firstName: {
        type: 'string',
        title: 'First name',
      },
      lastName: {
        type: 'string',
        title: 'Last name',
      },
      birthday: {
        type: 'string',
        format: 'date',
        title: 'Birthday',
      },
      sex: {
        type: 'string',
        title: 'sex',
        enum: ['Male', 'Female', 'Non-binary', 'Other'],
        'x-display-value': {
          Male: 'male',
          Female: 'female',
          'Non-binary': 'non-binary',
          Other: 'other',
        },
      },
      address: {
        $$ref: '#/definitions/Address',
        type: 'object',
        required: ['address1', 'city'],
        properties: {
          address1: {
            type: 'string',
            title: 'Address 1',
          },
          city: {
            type: 'string',
            title: 'City',
          },
        },
      },
    },
  };

  const testData1 = {
    firstName: 'james',
    address: {
      address1: '1333 Minna',
    },
  };

  const expectedError1 = {
    lastName: 'Required.',
    address: { city: 'Required.' },
  };

  const testData2 = {
    firstName: 'james',
  };

  const expectedError2 = {
    lastName: 'Required.',
  };

  const testData3 = {
    firstName: 'james',
    address: {
      address1: '1333 Minna',
      city: 'SF',
    },
  };

  const expectedError3 = {
    lastName: 'Required.',
  };

  const testData4 = {
    firstName: 'james',
    lastName: 'franco',
    address: {
      address1: '1444 Minna',
      city: 'SF',
    },
  };

  const expectedError4 = {};

  const testData5 = {
    firstName: 'james',
    lastName: 'franco',
  };

  const expectedError5 = {};

  const tests = [
    [testData1, expectedError1, 'partial address errors'],
    [testData2, expectedError2, 'omitted address is fine'],
    [testData3, expectedError3, 'complete address is fine'],
    [testData4, expectedError4, 'complete field is fine'],
    [testData5, expectedError5, 'missing address is complete'],
  ];

  tests.forEach(testCase => {
    const [testData, expectedError, name] = testCase;
    it(name, () => {
      const errors = recursivelyValidateRequiredFields(testData, testSchema);

      expect(errors).toEqual(expectedError);
    });
  });
});

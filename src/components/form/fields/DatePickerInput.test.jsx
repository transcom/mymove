import React from 'react';
import { shallow } from 'enzyme';

import { DatePickerInput } from './DatePickerInput';

import SingleDatePicker from 'shared/JsonSchemaForm/SingleDatePicker';
import { ErrorMessage } from '..';

const mockSetValue = jest.fn();
const formik = require('formik');

const getShallowWrapper = (withError = false) => {
  const meta = withError ? { touched: true, error: 'sample error' } : { touched: false, error: '' };
  formik.useField = jest.fn(() => [{}, meta, { setValue: mockSetValue }]);
  return shallow(<DatePickerInput name="name" label="title" />);
};

describe('DatePickerInput', () => {
  describe('with all required props', () => {
    it('renders no ErrorMessage', () => {
      const errorMessage = getShallowWrapper().find(ErrorMessage);
      expect(errorMessage.length).toBe(1);
      expect(errorMessage.prop('display')).toBe(false);
      expect(errorMessage.prop('children')).toBe('');
    });

    const wrapper = getShallowWrapper(true);

    it('renders an ErrorMessage', () => {
      const errorMessage = wrapper.find(ErrorMessage);
      expect(errorMessage.length).toBe(1);
      expect(errorMessage.prop('display')).toBe(true);
      expect(errorMessage.prop('children')).toBe('sample error');
    });

    it('renders a SingleDatePicker input', () => {
      const input = wrapper.find(SingleDatePicker);
      expect(input.length).toBe(1);
    });

    it('triggers onChange properly', () => {
      const input = wrapper.find(SingleDatePicker);
      input.simulate('change', '16 Jun 2020');
      expect(mockSetValue).toHaveBeenCalledWith('16 Jun 2020');
    });
  });

  afterEach(jest.resetAllMocks);
});

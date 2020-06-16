import React from 'react';
import { shallow } from 'enzyme';
import SingleDatePicker from 'shared/JsonSchemaForm/SingleDatePicker';
import { ErrorMessage } from '..';
import { DatePickerInput } from './DatePickerInput';

const mockSetValue = jest.fn();
// mock out formik hook as we are not testing formik
// needs to be before first describe
jest.mock('formik', () => {
  return {
    ...jest.requireActual('formik'),
    useField: () => [{}, { touched: true, error: 'sample error' }, { setValue: mockSetValue }],
  };
});

describe('DatePickerInput', () => {
  describe('with all required props', () => {
    const wrapper = shallow(<DatePickerInput name="name" label="title" />);

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

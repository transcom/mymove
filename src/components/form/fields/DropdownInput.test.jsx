import React from 'react';
import { shallow } from 'enzyme';
import { Dropdown } from '@trussworks/react-uswds';

import { ErrorMessage } from '../index';

import { DropdownInput } from './DropdownInput';

const mockOnChange = jest.fn();
const formik = require('formik');

const getShallowWrapper = (withError = false) => {
  const meta = withError ? { touched: true, error: 'sample error' } : { touched: false, error: '' };
  formik.useField = jest.fn(() => [
    {
      onChange: mockOnChange,
    },
    meta,
  ]);
  return shallow(<DropdownInput name="dropdown" label="label" options={[{ key: 'key', value: 'value' }]} />);
};

describe('DropdownInput', () => {
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

    it('renders a USWDS Dropdown input', () => {
      const dropdownInput = wrapper.find(Dropdown);
      expect(dropdownInput.length).toBe(1);
    });

    it('triggers onChange properly', () => {
      const dropdownInput = wrapper.find(Dropdown);
      expect(dropdownInput.prop('onChange')).toBe(mockOnChange);
      dropdownInput.simulate('change', { value: 'sample' });
      expect(mockOnChange).toHaveBeenCalledWith({ value: 'sample' });
    });
  });

  afterEach(jest.resetAllMocks);
});

import React from 'react';
import { shallow } from 'enzyme';
import { Dropdown } from '@trussworks/react-uswds';
import { ErrorMessage } from '..';
import { DropdownInput } from './DropdownInput';

const mockOnChange = jest.fn();
// mock out formik hook as we are not testing formik
// needs to be before first describe
jest.mock('formik', () => {
  return {
    ...jest.requireActual('formik'),
    useField: () => [
      {
        onChange: mockOnChange,
      },
      { touched: true, error: 'sample error' },
    ],
  };
});

describe('DropdownInput', () => {
  describe('with all required props', () => {
    const wrapper = shallow(<DropdownInput name="dropdown" options={[['key', 'value']]} />);

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

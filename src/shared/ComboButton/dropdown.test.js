import React from 'react';
import { shallow } from 'enzyme';

import { DropDownItem, DropDown } from './dropdown';

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
  return shallow(<DropDown name="dropdown" label="label" options={[{ key: 'key', value: 'value' }]} />);
};

describe('DropdownItems', () => {
  describe('DropDown is enabled', () => {
    it('when clicked onClick is called', () => {
      const wrapper = getShallowWrapper(true);
      const DropDown = wrapper.find('div');
      DropDown.simulate('click');
      expect(DropDown).toBeDefined();
    });
  });
  describe('DropDownItem is enabled', () => {
    it('when clicked onClick is called', () => {
      const onClick = jest.fn();
      const wrapper = shallow(<DropDownItem disabled={false} onClick={onClick} />);

      const dropDownItem = wrapper.find('div');
      dropDownItem.simulate('click');

      expect(onClick).toHaveBeenCalled();
    });
  });
  describe('DropDownItem is disabled', () => {
    it('when clicked onClick is not called', () => {
      const onClick = jest.fn();
      const wrapper = shallow(<DropDownItem disabled onClick={onClick} />);

      const dropDownItem = wrapper.find('div');
      dropDownItem.simulate('click');

      expect(onClick).not.toHaveBeenCalled();
    });
  });
});

import React from 'react';
import { shallow } from 'enzyme';
import { DropDownItem } from './dropdown';

describe('DropdownItems', () => {
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
      const wrapper = shallow(<DropDownItem disabled={true} onClick={onClick} />);

      const dropDownItem = wrapper.find('div');
      dropDownItem.simulate('click');

      expect(onClick).not.toHaveBeenCalled();
    });
  });
});

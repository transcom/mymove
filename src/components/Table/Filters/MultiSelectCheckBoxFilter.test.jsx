import React from 'react';
import Select from 'react-select';
import { mount } from 'enzyme';

import MultiSelectCheckBoxFilter from './MultiSelectCheckBoxFilter';

const column = {
  filterValue: '',
  setFilter: jest.fn(),
};
const optionsConstants = [
  { value: 'ARMY', label: 'Army' },
  { value: 'NAVY', label: 'Navy' },
];
const optionsStrings = [
  { value: 'approval requested', label: 'Approval Requested' },
  { value: 'paid', label: 'Paid' },
];

describe('MultiSelectCheckBoxFilter', () => {
  it('renders without crashing', () => {
    const wrapper = mount(<MultiSelectCheckBoxFilter options={[{ label: 'test', value: 'test' }]} column={{}} />);
    expect(wrapper.find('[data-testid="MultiSelectCheckBoxFilter"]').length).toBe(1);
    expect(wrapper.find('.MultiSelectCheckBoxFilter__placeholder').at(1).text('Select...')).toBeTruthy();
  });

  describe('It renders the expected placeholder text', () => {
    it('from an array of constant valued objects when a value is chosen', () => {
      const wrapper = mount(<MultiSelectCheckBoxFilter options={optionsConstants} column={column} />);
      const input = wrapper.find(Select).find('input');
      input.simulate('keyDown', { key: 'ArrowDown', keyCode: 40 });
      input.simulate('keyDown', { key: 'Enter', keyCode: 13 });

      expect(wrapper.find('[data-testid="multi-value-container"]').text()).toEqual('Army');
    });

    it('from an array of string valued objects when a value is chosen', () => {
      const wrapper = mount(<MultiSelectCheckBoxFilter options={optionsStrings} column={column} />);
      const input = wrapper.find(Select).find('input');
      input.simulate('keyDown', { key: 'ArrowDown', keyCode: 40 });
      input.simulate('keyDown', { key: 'Enter', keyCode: 13 });

      expect(wrapper.find('[data-testid="multi-value-container"]').text()).toEqual('Approval Requested');
    });

    it('from an array of objects when multiple values are chosen', () => {
      const wrapper = mount(<MultiSelectCheckBoxFilter options={optionsStrings} column={column} />);
      const input = wrapper.find(Select).find('input');
      input.simulate('keyDown', { key: 'ArrowDown', keyCode: 40 });
      input.simulate('keyDown', { key: 'Enter', keyCode: 13 });
      input.simulate('keyDown', { key: 'ArrowDown', keyCode: 40 });
      input.simulate('keyDown', { key: 'ArrowDown', keyCode: 40 });
      input.simulate('keyDown', { key: 'Enter', keyCode: 13 });

      expect(wrapper.find('[data-testid="multi-value-container"]').at(0).text()).toEqual('Approval Requested');
      expect(wrapper.find('[data-testid="multi-value-container"]').at(1).text()).toEqual('Paid');
    });
  });
});

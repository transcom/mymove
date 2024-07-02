import React from 'react';
import Select from 'react-select';
import { mount } from 'enzyme';

import MutliSelectTypeAheadCheckboxFilter from './MutliSelectTypeAheadCheckboxFilter';

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

describe('MutliSelectTypeAheadCheckboxFilter', () => {
  it('renders without crashing', () => {
    const wrapper = mount(
      <MutliSelectTypeAheadCheckboxFilter
        options={[{ label: 'test', value: 'test' }]}
        placeholder="Select..."
        column={{}}
      />,
    );
    expect(wrapper.find('[data-testid="MultiSelectTypeAheadCheckBoxFilter"]').length).toBe(1);
  });

  describe('It renders the expected placeholder text', () => {
    it('from an array of constant valued objects when a value is chosen', () => {
      const wrapper = mount(<MutliSelectTypeAheadCheckboxFilter options={optionsConstants} column={column} />);
      const input = wrapper.find(Select).find('input');
      input.simulate('keyDown', { key: 'ArrowDown', keyCode: 40 });
      input.simulate('keyDown', { key: 'Enter', keyCode: 13 });

      expect(wrapper.find('[data-testid="multi-value-container"]').text()).toEqual('Army');
    });

    it('from an array of string valued objects when a value is chosen', () => {
      const wrapper = mount(<MutliSelectTypeAheadCheckboxFilter options={optionsStrings} column={column} />);
      const input = wrapper.find(Select).find('input');
      input.simulate('keyDown', { key: 'ArrowDown', keyCode: 40 });
      input.simulate('keyDown', { key: 'Enter', keyCode: 13 });

      expect(wrapper.find('[data-testid="multi-value-container"]').text()).toEqual('Approval Requested');
    });

    it('from an array of objects when multiple values are chosen', () => {
      const wrapper = mount(<MutliSelectTypeAheadCheckboxFilter options={optionsStrings} column={column} />);
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

import React from 'react';
import { mount } from 'enzyme';

import Table from './Table';

import { createHeader } from 'components/Table/utils';

describe('React table', () => {
  it('renders without crashing', () => {
    const wrapper = mount(<Table />);
    expect(wrapper.find('[data-testid="react-table"]').length).toBe(1);
  });

  it('renders with data', () => {
    const data = [{ col1: 'Column 1 value' }];
    const cols = [createHeader('Column 1 header', 'col1')];
    const wrapper = mount(<Table data={data} columns={cols} />);

    // checking to see if we get expected lengths
    expect(wrapper.find('th').length).toBe(1);
    expect(wrapper.find('td').length).toBe(1);

    // checking data
    expect(wrapper.find('th[data-testid="col1"]').text()).toBe('Column 1 header');
    // data-testid has a format of ${columnKey}-${rowIndex}
    expect(wrapper.find('td[data-testid="col1-0"]').text()).toBe('Column 1 value');
  });
});

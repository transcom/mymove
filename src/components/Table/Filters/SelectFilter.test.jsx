import React from 'react';
import { mount } from 'enzyme';

import SelectFilter from './SelectFilter';

describe('React table', () => {
  it('renders without crashing', () => {
    const wrapper = mount(
      <SelectFilter
        column={{
          filterValue: '',
          preFilteredRows: [],
          setFilter: jest.fn(),
        }}
        options={[
          { value: 'ARMY', label: 'Army' },
          { value: 'NAVY', label: 'Navy' },
        ]}
      />,
    );
    expect(wrapper.find('select[data-testid="SelectFilter"]').length).toBe(1);
  });
});

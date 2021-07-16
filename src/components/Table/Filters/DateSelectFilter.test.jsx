import React from 'react';
import { mount } from 'enzyme';

import DateSelectFilter from './DateSelectFilter';

describe('React table', () => {
  it('renders without crashing as a Date', () => {
    const wrapper = mount(
      <DateSelectFilter
        column={{
          filterValue: '',
          setFilter: jest.fn(),
        }}
      />,
    );
    expect(wrapper.find('[data-testid="DateSelectFilter"]').length).toBe(1);
  });
  it('renders without crashing as a DateTime', () => {
    const wrapper = mount(
      <DateSelectFilter
        dateTime
        column={{
          filterValue: '',
          setFilter: jest.fn(),
        }}
      />,
    );
    expect(wrapper.find('[data-testid="DateSelectFilter"]').length).toBe(1);
  });
});

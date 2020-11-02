import React from 'react';
import { mount } from 'enzyme';

import DateSelectFilter from './DateSelectFilter';

describe('React table', () => {
  it('renders without crashing', () => {
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
});

import React from 'react';
import { mount } from 'enzyme';

import TextBoxFilter from './TextBoxFilter';

describe('React table', () => {
  it('renders without crashing', () => {
    const wrapper = mount(
      <TextBoxFilter
        column={{
          filterValue: '',
          setFilter: jest.fn(),
        }}
      />,
    );
    expect(wrapper.find(TextBoxFilter).length).toBe(1);
  });
});

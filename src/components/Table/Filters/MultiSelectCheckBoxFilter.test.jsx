import React from 'react';
import { mount } from 'enzyme';

import MultiSelectCheckBoxFilter from './MultiSelectCheckBoxFilter';

describe('MultiSelectCheckBoxFilter', () => {
  it('renders without crashing', () => {
    const wrapper = mount(<MultiSelectCheckBoxFilter options={[{ label: 'test', value: 'test' }]} column={{}} />);
    expect(wrapper.find('[data-testid="MultiSelectCheckBoxFilter"]').length).toBe(1);
  });
});

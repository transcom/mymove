/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { mount } from 'enzyme';

import SelectableCard from './SelectableCard';

const defaultProps = {
  id: '123',
  label: 'My Favorite Card',
  name: 'card',
  value: 'card1',
  cardText: 'This is the best card in the world because it is just the best',
  onChange: jest.fn(),
};
function mountSelectableCard(props = defaultProps) {
  return mount(<SelectableCard {...props} />);
}
describe('SelectableCard component', () => {
  it('renders without crashing', () => {
    const wrapper = mountSelectableCard();
    expect(wrapper.find('SelectableCard').length).toBe(1);
  });
});

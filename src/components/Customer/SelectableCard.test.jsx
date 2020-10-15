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

const mountSelectableCard = (props) => mount(<SelectableCard {...defaultProps} {...props} />);

describe('SelectableCard component', () => {
  describe('with default props', () => {
    const wrapper = mountSelectableCard();
    it('renders without crashing', () => {
      expect(wrapper.find('SelectableCard').length).toBe(1);
    });
  });

  describe('with a help button', () => {
    const mockHelpClick = jest.fn();
    const wrapper = mountSelectableCard({ onHelpClick: mockHelpClick });

    it('renders the help icon', () => {
      expect(wrapper.find('[data-testid="helpButton"]').exists()).toBe(true);
    });

    it('calls the help handler', () => {
      wrapper.find('button[data-testid="helpButton"]').simulate('click');
      expect(mockHelpClick).toHaveBeenCalled();
    });
  });
});

import React from 'react';
import { mount } from 'enzyme';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

import PrimeSimulatorAvailableMoves from './AvailableMovesQueue';

const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
}));

describe('getPrimeAvailableMoves', () => {
  const queryClient = new QueryClient();
  const wrapper = mount(
    <QueryClientProvider client={queryClient}>
      <PrimeSimulatorAvailableMoves />
    </QueryClientProvider>,
  );

  it('renders the page', () => {
    const filterInput = wrapper.find({ 'data-testid': 'prime-date-filter-input' });
    const filterButton = wrapper.find({ 'data-testid': 'prime-date-filter-button' });

    expect(filterInput);
    expect(filterButton);

    filterInput.simulate('change', { value: '2023-13-299' });

    wrapper
      .find('input')
      .at(0)
      .simulate('change', { target: { name: 'width', value: '2023-13-28' } });
    filterButton.simulate('click');
    expect(wrapper.text('Enter a valid date.'));
  });
});

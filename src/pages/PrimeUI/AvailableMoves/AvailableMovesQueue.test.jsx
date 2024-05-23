import React from 'react';
import { mount } from 'enzyme';

import PrimeSimulatorAvailableMoves from './AvailableMovesQueue';

import { MockProviders } from 'testUtils';

jest.mock('hooks/queries', () => ({
  useUserQueries: () => {
    return {
      isLoading: false,
      isError: false,
      data: {},
    };
  },
  usePrimeSimulatorAvailableMovesQueries: () => {
    return {
      queueResult: {
        perPage: 20,
        page: 1,
        totalCount: 1,
        data: [
          {
            availableToPrimeAt: '2023-12-07T18:21:34.897Z',
            createdAt: '2023-12-07T18:21:34.898Z',
            eTag: 'MjAyMy0xMi0wN1QxODoyMTozNC44OTg3NTda',
            id: '6b3df50c-9b7a-4d95-92d0-796563e50fa8',
            moveCode: 'SidDLH',
            orderID: '07ba4203-8b60-4171-92d0-7d2454bac1b5',
            ppmType: 'PARTIAL',
            referenceId: '7789-3211',
            updatedAt: '2023-12-07T18:21:34.898Z',
          },
        ],
      },
      isLoading: false,
      isError: false,
    };
  },
}));

describe('MoveQueue', () => {
  const wrapper = mount(
    <MockProviders>
      <PrimeSimulatorAvailableMoves />
    </MockProviders>,
  );

  it('should render the h1', () => {
    expect(wrapper.find('h1').text()).toBe('Moves available to Prime (1)');
  });

  it('should render the table', () => {
    expect(wrapper.find('Table').exists()).toBe(true);
  });

  it('should format the column data', () => {
    const moves = wrapper.find('tbody tr');

    const firstMove = moves.at(0);
    expect(firstMove.find({ 'data-testid': 'id-0' }).text()).toBe('6b3df50c-9b7a-4d95-92d0-796563e50fa8');
    expect(firstMove.find({ 'data-testid': 'moveCode-0' }).text()).toBe('SidDLH');
    expect(firstMove.find({ 'data-testid': 'createdAt-0' }).text()).toBe('07 Dec 2023, 06:21 pm');
    expect(firstMove.find({ 'data-testid': 'updatedAt-0' }).text()).toBe('07 Dec 2023, 06:21 pm');
    expect(firstMove.find({ 'data-testid': 'orderID-0' }).text()).toBe('07ba4203-8b60-4171-92d0-7d2454bac1b5');
    expect(firstMove.find({ 'data-testid': 'ppmType-0' }).text()).toBe('PARTIAL');
    expect(firstMove.find({ 'data-testid': 'referenceId-0' }).text()).toBe('7789-3211');
    expect(firstMove.find({ 'data-testid': 'availableToPrimeAt-0' }).text()).toBe('07 Dec 2023, 06:21 pm');
  });

  it('should render the pagination component', () => {
    expect(wrapper.find({ 'data-testid': 'pagination' }).exists()).toBe(true);
  });

  it('applies the sort to the status column in descending direction', () => {
    expect(wrapper.find({ 'data-testid': 'availableToPrimeAt' }).at(0).hasClass('sortAscending')).toBe(true);
  });

  it('toggles the sort direction when clicked', () => {
    const statusHeading = wrapper.find({ 'data-testid': 'availableToPrimeAt' }).at(0);

    statusHeading.simulate('click');
    wrapper.update();

    expect(wrapper.find({ 'data-testid': 'availableToPrimeAt' }).at(0).hasClass('sortDescending')).toBe(true);

    statusHeading.simulate('click');
    wrapper.update();

    // no sort direction should be applied
    expect(wrapper.find({ 'data-testid': 'availableToPrimeAt' }).at(0).hasClass('sortAscending')).toBe(false);
    expect(wrapper.find({ 'data-testid': 'availableToPrimeAt' }).at(0).hasClass('sortDescending')).toBe(false);

    const nameHeading = wrapper.find({ 'data-testid': 'orderID' }).at(0);
    nameHeading.simulate('click');
    wrapper.update();

    expect(wrapper.find({ 'data-testid': 'orderID' }).at(0).hasClass('sortAscending')).toBe(true);
  });
});

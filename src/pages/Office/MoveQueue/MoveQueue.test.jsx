import React from 'react';
import Select from 'react-select';
import { mount } from 'enzyme';

import MoveQueue from './MoveQueue';

import { MockProviders } from 'testUtils';

jest.mock('hooks/queries', () => ({
  useUserQueries: () => {
    return {
      isLoading: false,
      isError: false,
      data: {
        office_user: { transportation_office: { gbloc: 'TEST' } },
      },
    };
  },
  useMovesQueueQueries: () => {
    return {
      isLoading: false,
      isError: false,
      queueResult: {
        totalCount: 2,
        data: [
          {
            id: 'move1',
            customer: {
              agency: 'AIR_FORCE',
              first_name: 'test first',
              last_name: 'test last',
              dodID: '555555555',
            },
            locator: 'AB5P',
            departmentIndicator: 'ARMY',
            shipmentsCount: 2,
            status: 'SUBMITTED',
            originDutyLocation: {
              name: 'Area 51',
            },
            originGBLOC: 'EEEE',
            requestedMoveDate: '2023-02-10',
          },
          {
            id: 'move2',
            customer: {
              agency: 'MARINES',
              first_name: 'test another first',
              last_name: 'test another last',
              dodID: '4444444444',
            },
            locator: 'T12A',
            departmentIndicator: 'NAVY_AND_MARINES',
            shipmentsCount: 1,
            status: 'APPROVED',
            originDutyLocation: {
              name: 'Los Alamos',
            },
            originGBLOC: 'EEEE',
            requestedMoveDate: '2023-02-12',
          },
        ],
      },
    };
  },
}));

describe('MoveQueue', () => {
  const wrapper = mount(
    <MockProviders initialEntries={['moves/queue']}>
      <MoveQueue />
    </MockProviders>,
  );

  it('should render the h1', () => {
    expect(wrapper.find('h1').text()).toBe('All moves (2)');
  });

  it('should render the table', () => {
    expect(wrapper.find('Table').exists()).toBe(true);
  });

  it('should format the column data', () => {
    const moves = wrapper.find('tbody tr');

    const firstMove = moves.at(0);
    expect(firstMove.find({ 'data-testid': 'lastName-0' }).text()).toBe('test last, test first');
    expect(firstMove.find({ 'data-testid': 'dodID-0' }).text()).toBe('555555555');
    expect(firstMove.find({ 'data-testid': 'status-0' }).text()).toBe('New move');
    expect(firstMove.find({ 'data-testid': 'locator-0' }).text()).toBe('AB5P');
    expect(firstMove.find({ 'data-testid': 'branch-0' }).text()).toBe('Air Force');
    expect(firstMove.find({ 'data-testid': 'shipmentsCount-0' }).text()).toBe('2');
    expect(firstMove.find({ 'data-testid': 'originDutyLocation-0' }).text()).toBe('Area 51');
    expect(firstMove.find({ 'data-testid': 'originGBLOC-0' }).text()).toBe('EEEE');
    expect(firstMove.find({ 'data-testid': 'requestedMoveDate-0' }).text()).toBe('10 Feb 2023');

    const secondMove = moves.at(1);
    expect(secondMove.find({ 'data-testid': 'lastName-1' }).text()).toBe('test another last, test another first');
    expect(secondMove.find({ 'data-testid': 'dodID-1' }).text()).toBe('4444444444');
    expect(secondMove.find({ 'data-testid': 'status-1' }).text()).toBe('Move approved');
    expect(secondMove.find({ 'data-testid': 'locator-1' }).text()).toBe('T12A');
    expect(secondMove.find({ 'data-testid': 'branch-1' }).text()).toBe('Marine Corps');
    expect(secondMove.find({ 'data-testid': 'shipmentsCount-1' }).text()).toBe('1');
    expect(secondMove.find({ 'data-testid': 'originDutyLocation-1' }).text()).toBe('Los Alamos');
    expect(secondMove.find({ 'data-testid': 'originGBLOC-1' }).text()).toBe('EEEE');
    expect(secondMove.find({ 'data-testid': 'requestedMoveDate-1' }).text()).toBe('12 Feb 2023');
  });

  it('should render the pagination component', () => {
    expect(wrapper.find({ 'data-testid': 'pagination' }).exists()).toBe(true);
  });

  it('applies the sort to the status column in descending direction', () => {
    expect(wrapper.find({ 'data-testid': 'status' }).at(0).hasClass('sortAscending')).toBe(true);
  });

  it('toggles the sort direction when clicked', () => {
    const statusHeading = wrapper.find({ 'data-testid': 'status' }).at(0);

    statusHeading.simulate('click');
    wrapper.update();

    expect(wrapper.find({ 'data-testid': 'status' }).at(0).hasClass('sortDescending')).toBe(true);

    statusHeading.simulate('click');
    wrapper.update();

    // no sort direction should be applied
    expect(wrapper.find({ 'data-testid': 'status' }).at(0).hasClass('sortAscending')).toBe(false);
    expect(wrapper.find({ 'data-testid': 'status' }).at(0).hasClass('sortDescending')).toBe(false);

    const nameHeading = wrapper.find({ 'data-testid': 'lastName' }).at(0);
    nameHeading.simulate('click');
    wrapper.update();

    expect(wrapper.find({ 'data-testid': 'lastName' }).at(0).hasClass('sortAscending')).toBe(true);
  });

  it('filters the queue', () => {
    const input = wrapper.find(Select).at(0).find('input');
    input.simulate('keyDown', { key: 'ArrowDown', keyCode: 40 });
    input.simulate('keyDown', { key: 'Enter', keyCode: 13 });

    wrapper.update();
    expect(wrapper.find('[data-testid="multi-value-container"]').text()).toEqual('New move');
  });
});

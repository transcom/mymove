import React from 'react';
import Select from 'react-select';
import { mount } from 'enzyme';
import * as reactRouterDom from 'react-router-dom';
import { render, screen, waitFor } from '@testing-library/react';

import MoveQueue from './MoveQueue';

import { MockProviders } from 'testUtils';
import { MOVE_STATUS_OPTIONS, BRANCH_OPTIONS } from 'constants/queues';
import { generalRoutes, tooRoutes } from 'constants/routes';
import { isBooleanFlagEnabled } from 'utils/featureFlags';

jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'), // this line preserves the non-hook exports
  useParams: jest.fn(), // mock useParams
  useNavigate: jest.fn(), // mock useNavigate if needed
}));
jest.setTimeout(60000);

jest.mock('utils/featureFlags', () => ({
  ...jest.requireActual('utils/featureFlags'),
  isBooleanFlagEnabled: jest.fn().mockImplementation(() => Promise.resolve()),
}));

const moveData = [
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
    appearedInTooAt: '2023-02-10T00:00:00.000Z',
    lockExpiresAt: '2099-02-10T00:00:00.000Z',
    lockedByOfficeUserID: '2744435d-7ba8-4cc5-bae5-f302c72c966e',
    assignedTo: {
      officeUserId: 'exampleId2',
      firstName: 'John',
      lastName: 'Denver',
    },
    availableOfficeUsers: [
      {
        officeUserId: 'exampleId1',
        firstName: 'Jimmy',
        lastName: 'John',
      },
      {
        officeUserId: 'exampleId2',
        firstName: 'John',
        lastName: 'Denver',
      },
    ],
  },
  {
    id: 'move2',
    customer: {
      agency: 'COAST_GUARD',
      first_name: 'test another first',
      last_name: 'test another last',
      dodID: '4444444444',
      emplid: '4589652',
    },
    locator: 'T12A',
    departmentIndicator: 'COAST_GUARD',
    shipmentsCount: 1,
    status: 'APPROVED',
    originDutyLocation: {
      name: 'Los Alamos',
    },
    originGBLOC: 'EEEE',
    requestedMoveDate: '2023-02-12',
    appearedInTooAt: '2023-02-12T00:00:00.000Z',
    assignedTo: {
      officeUserId: 'exampleId2',
      firstName: 'John',
      lastName: 'Denver',
    },
    availableOfficeUsers: [
      {
        officeUserId: 'exampleId1',
        firstName: 'Jimmy',
        lastName: 'John',
      },
      {
        officeUserId: 'exampleId2',
        firstName: 'John',
        lastName: 'Denver',
      },
    ],
  },
  {
    id: 'move3',
    customer: {
      agency: 'Marine Corps',
      first_name: 'will',
      last_name: 'robinson',
      dodID: '6666666666',
    },
    locator: 'PREP',
    departmentIndicator: 'MARINES',
    shipmentsCount: 1,
    status: 'SUBMITTED',
    originDutyLocation: {
      name: 'Area 52',
    },
    originGBLOC: 'EEEE',
    requestedMoveDate: '2023-03-12',
    appearedInTooAt: '2023-03-12T00:00:00.000Z',
    lockExpiresAt: '2099-03-12T00:00:00.000Z',
    lockedByOfficeUserID: '2744435d-7ba8-4cc5-bae5-f302c72c966e',
    assignedTo: {
      officeUserId: 'exampleId1',
      firstName: 'Jimmy',
      lastName: 'John',
    },
    availableOfficeUsers: [
      {
        officeUserId: 'exampleId1',
        firstName: 'Jimmy',
        lastName: 'John',
      },
      {
        officeUserId: 'exampleId2',
        firstName: 'John',
        lastName: 'Denver',
      },
    ],
  },
];

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
        totalCount: 3,
        data: moveData,
      },
    };
  },
}));

const GetMountedComponent = (queueTypeToMount) => {
  reactRouterDom.useParams.mockReturnValue({ queueType: queueTypeToMount });
  const wrapper = mount(
    <MockProviders>
      <MoveQueue isQueueManagementFFEnabled />
    </MockProviders>,
  );
  return wrapper;
};
const SEARCH_OPTIONS = ['Move Code', 'DoD ID', 'Customer Name', 'Payment Request Number'];
describe('MoveQueue', () => {
  afterEach(() => {
    jest.restoreAllMocks();
  });

  it('should render the h1', () => {
    expect(GetMountedComponent(tooRoutes.MOVE_QUEUE).find('h1').text()).toBe('All moves (3)');
  });

  it('should render the table', () => {
    expect(GetMountedComponent(tooRoutes.MOVE_QUEUE).find('Table').exists()).toBe(true);
  });

  it('should format the column data', () => {
    let currentIndex = 0;
    let currentMove;
    const moves = GetMountedComponent(tooRoutes.MOVE_QUEUE).find('tbody tr');

    currentMove = moves.at(currentIndex);
    expect(currentMove.find({ 'data-testid': `lastName-${currentIndex}` }).text()).toBe(
      `${moveData[currentIndex].customer.last_name}, ${moveData[currentIndex].customer.first_name}`,
    );
    expect(currentMove.find({ 'data-testid': `dodID-${currentIndex}` }).text()).toBe(
      moveData[currentIndex].customer.dodID,
    );
    expect(currentMove.find({ 'data-testid': `status-${currentIndex}` }).text()).toBe('New move');
    expect(currentMove.find({ 'data-testid': `locator-${currentIndex}` }).text()).toBe(moveData[currentIndex].locator);
    expect(currentMove.find({ 'data-testid': `branch-${currentIndex}` }).text()).toBe(
      BRANCH_OPTIONS.find((value) => value.value === moveData[currentIndex].customer.agency).label,
    );
    expect(currentMove.find({ 'data-testid': `shipmentsCount-${currentIndex}` }).text()).toBe(
      moveData[currentIndex].shipmentsCount.toString(),
    );
    expect(currentMove.find({ 'data-testid': `originDutyLocation-${currentIndex}` }).text()).toBe(
      moveData[currentIndex].originDutyLocation.name,
    );
    expect(currentMove.find({ 'data-testid': `originGBLOC-${currentIndex}` }).text()).toBe(
      moveData[currentIndex].originGBLOC,
    );
    expect(currentMove.find({ 'data-testid': `requestedMoveDate-${currentIndex}` }).text()).toBe('10 Feb 2023');
    expect(currentMove.find({ 'data-testid': `appearedInTooAt-${currentIndex}` }).text()).toBe('10 Feb 2023');

    currentIndex += 1;
    currentMove = moves.at(currentIndex);
    expect(currentMove.find({ 'data-testid': `lastName-${currentIndex}` }).text()).toBe(
      'test another last, test another first',
    );
    expect(currentMove.find({ 'data-testid': `lastName-${currentIndex}` }).text()).toBe(
      `${moveData[currentIndex].customer.last_name}, ${moveData[currentIndex].customer.first_name}`,
    );
    expect(currentMove.find({ 'data-testid': `dodID-${currentIndex}` }).text()).toBe(
      moveData[currentIndex].customer.dodID,
    );
    expect(currentMove.find({ 'data-testid': `emplid-${currentIndex}` }).text()).toBe(
      moveData[currentIndex].customer.emplid,
    );
    expect(currentMove.find({ 'data-testid': `status-${currentIndex}` }).text()).toBe('Move approved');
    expect(currentMove.find({ 'data-testid': `locator-${currentIndex}` }).text()).toBe(moveData[currentIndex].locator);
    expect(currentMove.find({ 'data-testid': `branch-${currentIndex}` }).text()).toBe(
      BRANCH_OPTIONS.find((value) => value.value === moveData[currentIndex].customer.agency).label,
    );
    expect(currentMove.find({ 'data-testid': `shipmentsCount-${currentIndex}` }).text()).toBe(
      moveData[currentIndex].shipmentsCount.toString(),
    );
    expect(currentMove.find({ 'data-testid': `originDutyLocation-${currentIndex}` }).text()).toBe(
      moveData[currentIndex].originDutyLocation.name,
    );
    expect(currentMove.find({ 'data-testid': `originGBLOC-${currentIndex}` }).text()).toBe(
      moveData[currentIndex].originGBLOC,
    );
    expect(currentMove.find({ 'data-testid': `requestedMoveDate-${currentIndex}` }).text()).toBe('12 Feb 2023');
    expect(currentMove.find({ 'data-testid': `appearedInTooAt-${currentIndex}` }).text()).toBe('12 Feb 2023');

    currentIndex += 1;
    currentMove = moves.at(currentIndex);
    expect(currentMove.find({ 'data-testid': `lastName-${currentIndex}` }).text()).toBe(
      `${moveData[currentIndex].customer.last_name}, ${moveData[currentIndex].customer.first_name}`,
    );
    expect(currentMove.find({ 'data-testid': `dodID-${currentIndex}` }).text()).toBe(
      moveData[currentIndex].customer.dodID,
    );
    expect(currentMove.find({ 'data-testid': `status-${currentIndex}` }).text()).toBe('New move');
    expect(currentMove.find({ 'data-testid': `locator-${currentIndex}` }).text()).toBe(moveData[currentIndex].locator);
    expect(currentMove.find({ 'data-testid': `branch-${currentIndex}` }).text()).toBe(
      moveData[currentIndex].customer.agency.toString(),
    );
    expect(currentMove.find({ 'data-testid': `shipmentsCount-${currentIndex}` }).text()).toBe(
      moveData[currentIndex].shipmentsCount.toString(),
    );
    expect(currentMove.find({ 'data-testid': `originDutyLocation-${currentIndex}` }).text()).toBe(
      moveData[currentIndex].originDutyLocation.name,
    );
    expect(currentMove.find({ 'data-testid': `originGBLOC-${currentIndex}` }).text()).toBe(
      moveData[currentIndex].originGBLOC,
    );
    expect(currentMove.find({ 'data-testid': `requestedMoveDate-${currentIndex}` }).text()).toBe('12 Mar 2023');
    expect(currentMove.find({ 'data-testid': `appearedInTooAt-${currentIndex}` }).text()).toBe('12 Mar 2023');
    expect(currentMove.find({ 'data-testid': `assignedTo-${currentIndex}` }).text()).toBe('John, Jimmy');
  });

  it('should render the pagination component', () => {
    expect(GetMountedComponent(tooRoutes.MOVE_QUEUE).find({ 'data-testid': 'pagination' }).exists()).toBe(true);
  });

  it('applies the sort to the status column in descending direction', () => {
    expect(
      GetMountedComponent(tooRoutes.MOVE_QUEUE).find({ 'data-testid': 'status' }).at(0).hasClass('sortAscending'),
    ).toBe(true);
  });

  it('toggles the sort direction when clicked', () => {
    const wrapper = GetMountedComponent(tooRoutes.MOVE_QUEUE);
    const statusHeading = wrapper.find({ 'data-testid': 'status' }).at(0);

    statusHeading.simulate('click');
    GetMountedComponent(tooRoutes.MOVE_QUEUE).update();
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
    const wrapper = GetMountedComponent(tooRoutes.MOVE_QUEUE);
    const input = wrapper.find(Select).at(0).find('input');
    input.simulate('keyDown', { key: 'ArrowDown', keyCode: 40 });
    input.simulate('keyDown', { key: 'Enter', keyCode: 13 });

    wrapper.update();
    expect(wrapper.find('[data-testid="multi-value-container"]').text()).toEqual('New move');
  });
  it('renders Search and Move Queue tabs', () => {
    reactRouterDom.useParams.mockReturnValue({ queueType: generalRoutes.QUEUE_SEARCH_PATH });
    render(
      <reactRouterDom.BrowserRouter>
        <MoveQueue />
      </reactRouterDom.BrowserRouter>,
    );
    expect(screen.getByTestId('closeout-tab-link')).toBeInTheDocument();
    expect(screen.getByTestId('search-tab-link')).toBeInTheDocument();
    expect(screen.getByText('Task Order Queue', { selector: 'span' })).toBeInTheDocument();
    expect(screen.getByText('Search', { selector: 'span' })).toBeInTheDocument();
  });
  it('renders TableQueue when Search tab is selected', () => {
    reactRouterDom.useParams.mockReturnValue({ queueType: tooRoutes.MOVE_QUEUE });
    render(
      <reactRouterDom.BrowserRouter>
        <MoveQueue />
      </reactRouterDom.BrowserRouter>,
    );
    expect(screen.queryByTestId('table-queue')).toBeInTheDocument();
    expect(screen.queryByTestId('move-search')).not.toBeInTheDocument();
  });
  it('has all options for searches', async () => {
    reactRouterDom.useParams.mockReturnValue({ queueType: generalRoutes.QUEUE_SEARCH_PATH });
    render(
      <reactRouterDom.BrowserRouter>
        <MoveQueue />
      </reactRouterDom.BrowserRouter>,
    );
    SEARCH_OPTIONS.forEach((option) => expect(screen.findByLabelText(option)));
  });
  it('Has all status options for move search', async () => {
    reactRouterDom.useParams.mockReturnValue({ queueType: generalRoutes.QUEUE_SEARCH_PATH });
    render(
      <reactRouterDom.BrowserRouter>
        <MoveQueue />
      </reactRouterDom.BrowserRouter>,
    );
    MOVE_STATUS_OPTIONS.forEach((option) => expect(screen.findByLabelText(option)));
  });

  it('Has all status options for move queue', async () => {
    reactRouterDom.useParams.mockReturnValue({ queueType: tooRoutes.MOVE_QUEUE });
    render(
      <reactRouterDom.BrowserRouter>
        <MoveQueue />
      </reactRouterDom.BrowserRouter>,
    );
    MOVE_STATUS_OPTIONS.forEach((option) => expect(screen.findByLabelText(option)));
  });
  it('renders a 404 if a bad route is provided', async () => {
    reactRouterDom.useParams.mockReturnValue({ queueType: 'BadRoute' });
    render(
      <reactRouterDom.BrowserRouter>
        <MoveQueue />
      </reactRouterDom.BrowserRouter>,
    );
    await expect(screen.getByText('Error - 404')).toBeInTheDocument();
    await expect(screen.getByText("We can't find the page you're looking for")).toBeInTheDocument();
  });
  it('renders a lock icon when move lock flag is on', async () => {
    isBooleanFlagEnabled.mockResolvedValue(true);
    reactRouterDom.useParams.mockReturnValue({ queueType: tooRoutes.MOVE_QUEUE });
    render(
      <reactRouterDom.BrowserRouter>
        <MoveQueue />
      </reactRouterDom.BrowserRouter>,
    );
    await waitFor(() => {
      const lockIcon = screen.queryAllByTestId('lock-icon')[0];
      expect(lockIcon).toBeInTheDocument();
    });
  });
  it('does NOT render a lock icon when move lock flag is off', async () => {
    isBooleanFlagEnabled.mockResolvedValue(false);
    reactRouterDom.useParams.mockReturnValue({ queueType: tooRoutes.MOVE_QUEUE });
    render(
      <reactRouterDom.BrowserRouter>
        <MoveQueue />
      </reactRouterDom.BrowserRouter>,
    );
    await await waitFor(() => {
      const lockIcon = screen.queryByTestId('lock-icon');
      expect(lockIcon).not.toBeInTheDocument();
    });
  });
  it('renders an assigned column when the queue management flag is on', async () => {
    reactRouterDom.useParams.mockReturnValue({ queueType: tooRoutes.MOVE_QUEUE });
    render(
      <reactRouterDom.BrowserRouter>
        <MoveQueue isQueueManagementFFEnabled />
      </reactRouterDom.BrowserRouter>,
    );
    await waitFor(() => {
      const assignedSelect = screen.queryAllByTestId('assigned-col')[0];
      expect(assignedSelect).toBeInTheDocument();
    });
  });
  it('renders an assigned column when the queue management flag is off', async () => {
    reactRouterDom.useParams.mockReturnValue({ queueType: tooRoutes.MOVE_QUEUE });
    render(
      <reactRouterDom.BrowserRouter>
        <MoveQueue isQueueManagementFFEnabled={false} />
      </reactRouterDom.BrowserRouter>,
    );
    await waitFor(() => {
      const assignedSelect = screen.queryByTestId('assigned-col');
      expect(assignedSelect).not.toBeInTheDocument();
    });
  });
});

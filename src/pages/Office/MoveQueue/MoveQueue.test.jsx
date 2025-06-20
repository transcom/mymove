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
import { APPROVAL_REQUEST_TYPES } from 'constants/approvalRequestTypes';

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

beforeEach(() => {
  jest.clearAllMocks();
});

const moveData = [
  {
    id: 'move1',
    customer: {
      agency: 'AIR_FORCE',
      first_name: 'test first',
      last_name: 'test last',
      edipi: '555555555',
    },
    locator: 'AB5P',
    departmentIndicator: 'ARMY',
    shipmentsCount: 2,
    status: 'APPROVALS REQUESTED',
    originDutyLocation: {
      name: 'Area 51',
    },
    destinationDutyLocation: {
      name: 'Area 52',
    },
    originGBLOC: 'EEEE',
    counselingOffice: '67592323-fc7e-4b35-83a7-57faa53b7acf',
    requestedMoveDates: '10 Feb 2023, 10 Mar 2023',
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
      edipi: '4444444444',
      emplid: '4589652',
    },
    locator: 'T12A',
    departmentIndicator: 'COAST_GUARD',
    shipmentsCount: 1,
    status: 'APPROVED',
    originDutyLocation: {
      name: 'Los Alamos',
    },
    destinationDutyLocation: {
      name: 'Area 52',
    },
    originGBLOC: 'EEEE',
    counselingOffice: '67592323-fc7e-4b35-83a7-57faa53b7acf',
    requestedMoveDates: '12 Feb 2023',
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
      edipi: '6666666666',
    },
    locator: 'PREP',
    departmentIndicator: 'MARINES',
    shipmentsCount: 1,
    status: 'SUBMITTED',
    originDutyLocation: {
      name: 'Area 52',
    },
    destinationDutyLocation: {
      name: 'Area 52',
    },
    originGBLOC: 'EEEE',
    counselingOffice: '67592323-fc7e-4b35-83a7-57faa53b7acf',
    requestedMoveDates: '12 Mar 2023',
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
  useDestinationRequestsQueueQueries: () => {
    return {
      isLoading: false,
      isError: false,
      queueResult: {
        totalCount: 3,
        data: moveData,
      },
    };
  },
  useBulkAssignmentQueries: () => {
    return {
      availableOfficeUsers: [
        {
          firstName: 'John',
          lastName: 'Snow',
          officeUserId: '123',
          workload: 0,
        },
        {
          firstName: 'Jane',
          lastName: 'Doe',
          officeUserId: '456',
          workload: 1,
        },
        {
          firstName: 'Jimmy',
          lastName: 'Page',
          officeUserId: '789',
          workload: 2,
        },
      ],
      bulkAssignmentMoveIDs: ['1', '2', '3'],
    };
  },
}));

const GetMountedComponent = (queueTypeToMount) => {
  reactRouterDom.useParams.mockReturnValue({ queueType: queueTypeToMount });
  const wrapper = mount(
    <MockProviders>
      <MoveQueue isQueueManagementFFEnabled isApprovalRequestTypeFFEnabled />
    </MockProviders>,
  );
  return wrapper;
};

moveData[0].approvalRequestTypes = [
  APPROVAL_REQUEST_TYPES.EXCESS_WEIGHT,
  APPROVAL_REQUEST_TYPES.IBHF,
  APPROVAL_REQUEST_TYPES.IOASIT,
  APPROVAL_REQUEST_TYPES.DDSHUT,
];
moveData[1].approvalRequestTypes = [APPROVAL_REQUEST_TYPES.IOASIT, APPROVAL_REQUEST_TYPES.AMENDED_ORDERS];
moveData[2].approvalRequestTypes = [APPROVAL_REQUEST_TYPES.IOASIT, APPROVAL_REQUEST_TYPES.NEW_SHIPMENT];

const SEARCH_OPTIONS = ['Move Code', 'DoD ID', 'Customer Name', 'Payment Request Number'];
describe('MoveQueue & DestinationRequestsQueue', () => {
  afterEach(() => {
    jest.restoreAllMocks();
  });

  it('should render the h1 on each queue', () => {
    expect(GetMountedComponent(tooRoutes.MOVE_QUEUE).find('h1').text()).toBe('All moves (3)');
    expect(GetMountedComponent(tooRoutes.DESTINATION_REQUESTS_QUEUE).find('h1').text()).toBe(
      'Destination requests (3)',
    );
  });

  it('should render the table on each queue', () => {
    expect(GetMountedComponent(tooRoutes.MOVE_QUEUE).find('Table').exists()).toBe(true);
    expect(GetMountedComponent(tooRoutes.DESTINATION_REQUESTS_QUEUE).find('Table').exists()).toBe(true);
  });

  it('should format the column data - MoveQueue', () => {
    let currentIndex = 0;
    let currentMove;
    const moves = GetMountedComponent(tooRoutes.MOVE_QUEUE).find('tbody tr');

    currentMove = moves.at(currentIndex);
    expect(currentMove.find({ 'data-testid': `customerName-${currentIndex}` }).text()).toBe(
      `${moveData[currentIndex].customer.last_name}, ${moveData[currentIndex].customer.first_name}`,
    );
    expect(currentMove.find({ 'data-testid': `edipi-${currentIndex}` }).text()).toBe(
      moveData[currentIndex].customer.edipi,
    );
    expect(currentMove.find({ 'data-testid': `status-${currentIndex}` }).text()).toBe('Approvals requested');
    expect(currentMove.find({ 'data-testid': `approvalRequestTypes-${currentIndex}` }).text()).toBe(
      'Boat, Excess weight, SIT',
    );
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
    expect(currentMove.find({ 'data-testid': `counselingOffice-${currentIndex}` }).text()).toBe(
      moveData[currentIndex].counselingOffice,
    );
    expect(currentMove.find({ 'data-testid': `requestedMoveDate-${currentIndex}` }).text()).toBe(
      '10 Feb 2023, 10 Mar 2023',
    );
    expect(currentMove.find({ 'data-testid': `appearedInTooAt-${currentIndex}` }).text()).toBe('10 Feb 2023');

    currentIndex += 1;
    currentMove = moves.at(currentIndex);
    expect(currentMove.find({ 'data-testid': `customerName-${currentIndex}` }).text()).toBe(
      'test another last, test another first',
    );
    expect(currentMove.find({ 'data-testid': `customerName-${currentIndex}` }).text()).toBe(
      `${moveData[currentIndex].customer.last_name}, ${moveData[currentIndex].customer.first_name}`,
    );
    expect(currentMove.find({ 'data-testid': `edipi-${currentIndex}` }).text()).toBe(
      moveData[currentIndex].customer.edipi,
    );
    expect(currentMove.find({ 'data-testid': `emplid-${currentIndex}` }).text()).toBe(
      moveData[currentIndex].customer.emplid,
    );
    expect(currentMove.find({ 'data-testid': `status-${currentIndex}` }).text()).toBe('Move approved');
    expect(currentMove.find({ 'data-testid': `approvalRequestTypes-${currentIndex}` }).text()).toBe('');
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
    expect(currentMove.find({ 'data-testid': `counselingOffice-${currentIndex}` }).text()).toBe(
      moveData[currentIndex].counselingOffice,
    );
    expect(currentMove.find({ 'data-testid': `requestedMoveDate-${currentIndex}` }).text()).toBe('12 Feb 2023');
    expect(currentMove.find({ 'data-testid': `appearedInTooAt-${currentIndex}` }).text()).toBe('12 Feb 2023');

    currentIndex += 1;
    currentMove = moves.at(currentIndex);
    expect(currentMove.find({ 'data-testid': `customerName-${currentIndex}` }).text()).toBe(
      `${moveData[currentIndex].customer.last_name}, ${moveData[currentIndex].customer.first_name}`,
    );
    expect(currentMove.find({ 'data-testid': `edipi-${currentIndex}` }).text()).toBe(
      moveData[currentIndex].customer.edipi,
    );
    expect(currentMove.find({ 'data-testid': `status-${currentIndex}` }).text()).toBe('New move');
    expect(currentMove.find({ 'data-testid': `approvalRequestTypes-${currentIndex}` }).text()).toBe('');
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
    expect(currentMove.find({ 'data-testid': `counselingOffice-${currentIndex}` }).text()).toBe(
      moveData[currentIndex].counselingOffice,
    );
    expect(currentMove.find({ 'data-testid': `requestedMoveDate-${currentIndex}` }).text()).toBe('12 Mar 2023');
    expect(currentMove.find({ 'data-testid': `appearedInTooAt-${currentIndex}` }).text()).toBe('12 Mar 2023');
    expect(currentMove.find({ 'data-testid': `assignedTo-${currentIndex}` }).text()).toBe('John, Jimmy');
  });

  it('should format the column data - DestinationRequestsQueue', () => {
    let currentIndex = 0;
    let currentMove;
    const moves = GetMountedComponent(tooRoutes.DESTINATION_REQUESTS_QUEUE).find('tbody tr');

    currentMove = moves.at(currentIndex);
    expect(currentMove.find({ 'data-testid': `customerName-${currentIndex}` }).text()).toBe(
      `${moveData[currentIndex].customer.last_name}, ${moveData[currentIndex].customer.first_name}`,
    );
    expect(currentMove.find({ 'data-testid': `edipi-${currentIndex}` }).text()).toBe(
      moveData[currentIndex].customer.edipi,
    );
    expect(currentMove.find({ 'data-testid': `status-${currentIndex}` }).text()).toBe('Approvals requested');
    expect(currentMove.find({ 'data-testid': `approvalRequestTypes-${currentIndex}` }).text()).toBe('Shuttle');
    expect(currentMove.find({ 'data-testid': `locator-${currentIndex}` }).text()).toBe(moveData[currentIndex].locator);
    expect(currentMove.find({ 'data-testid': `branch-${currentIndex}` }).text()).toBe(
      BRANCH_OPTIONS.find((value) => value.value === moveData[currentIndex].customer.agency).label,
    );
    expect(currentMove.find({ 'data-testid': `shipmentsCount-${currentIndex}` }).text()).toBe(
      moveData[currentIndex].shipmentsCount.toString(),
    );
    expect(currentMove.find({ 'data-testid': `destinationDutyLocation-${currentIndex}` }).text()).toBe(
      moveData[currentIndex].destinationDutyLocation.name,
    );
    expect(currentMove.find({ 'data-testid': `counselingOffice-${currentIndex}` }).text()).toBe(
      moveData[currentIndex].counselingOffice,
    );
    expect(currentMove.find({ 'data-testid': `requestedMoveDate-${currentIndex}` }).text()).toBe(
      '10 Feb 2023, 10 Mar 2023',
    );
    expect(currentMove.find({ 'data-testid': `appearedInTooAt-${currentIndex}` }).text()).toBe('10 Feb 2023');

    currentIndex += 1;
    currentMove = moves.at(currentIndex);
    expect(currentMove.find({ 'data-testid': `customerName-${currentIndex}` }).text()).toBe(
      'test another last, test another first',
    );
    expect(currentMove.find({ 'data-testid': `customerName-${currentIndex}` }).text()).toBe(
      `${moveData[currentIndex].customer.last_name}, ${moveData[currentIndex].customer.first_name}`,
    );
    expect(currentMove.find({ 'data-testid': `edipi-${currentIndex}` }).text()).toBe(
      moveData[currentIndex].customer.edipi,
    );
    expect(currentMove.find({ 'data-testid': `emplid-${currentIndex}` }).text()).toBe(
      moveData[currentIndex].customer.emplid,
    );
    expect(currentMove.find({ 'data-testid': `status-${currentIndex}` }).text()).toBe('Move approved');
    expect(currentMove.find({ 'data-testid': `approvalRequestTypes-${currentIndex}` }).text()).toBe('');
    expect(currentMove.find({ 'data-testid': `locator-${currentIndex}` }).text()).toBe(moveData[currentIndex].locator);
    expect(currentMove.find({ 'data-testid': `branch-${currentIndex}` }).text()).toBe(
      BRANCH_OPTIONS.find((value) => value.value === moveData[currentIndex].customer.agency).label,
    );
    expect(currentMove.find({ 'data-testid': `shipmentsCount-${currentIndex}` }).text()).toBe(
      moveData[currentIndex].shipmentsCount.toString(),
    );
    expect(currentMove.find({ 'data-testid': `counselingOffice-${currentIndex}` }).text()).toBe(
      moveData[currentIndex].counselingOffice,
    );
    expect(currentMove.find({ 'data-testid': `appearedInTooAt-${currentIndex}` }).text()).toBe('12 Feb 2023');

    currentIndex += 1;
    currentMove = moves.at(currentIndex);
    expect(currentMove.find({ 'data-testid': `customerName-${currentIndex}` }).text()).toBe(
      `${moveData[currentIndex].customer.last_name}, ${moveData[currentIndex].customer.first_name}`,
    );
    expect(currentMove.find({ 'data-testid': `edipi-${currentIndex}` }).text()).toBe(
      moveData[currentIndex].customer.edipi,
    );
    expect(currentMove.find({ 'data-testid': `status-${currentIndex}` }).text()).toBe('New move');
    expect(currentMove.find({ 'data-testid': `approvalRequestTypes-${currentIndex}` }).text()).toBe('');
    expect(currentMove.find({ 'data-testid': `locator-${currentIndex}` }).text()).toBe(moveData[currentIndex].locator);
    expect(currentMove.find({ 'data-testid': `branch-${currentIndex}` }).text()).toBe(
      moveData[currentIndex].customer.agency.toString(),
    );
    expect(currentMove.find({ 'data-testid': `shipmentsCount-${currentIndex}` }).text()).toBe(
      moveData[currentIndex].shipmentsCount.toString(),
    );
    expect(currentMove.find({ 'data-testid': `destinationDutyLocation-${currentIndex}` }).text()).toBe(
      moveData[currentIndex].destinationDutyLocation.name,
    );
    expect(currentMove.find({ 'data-testid': `counselingOffice-${currentIndex}` }).text()).toBe(
      moveData[currentIndex].counselingOffice,
    );
    expect(currentMove.find({ 'data-testid': `requestedMoveDate-${currentIndex}` }).text()).toBe('12 Mar 2023');
    expect(currentMove.find({ 'data-testid': `appearedInTooAt-${currentIndex}` }).text()).toBe('12 Mar 2023');
    expect(currentMove.find({ 'data-testid': `assignedTo-${currentIndex}` }).text()).toBe('John, Jimmy');
  });

  it('should render the pagination component on both queues', () => {
    expect(GetMountedComponent(tooRoutes.MOVE_QUEUE).find({ 'data-testid': 'pagination' }).exists()).toBe(true);
    expect(
      GetMountedComponent(tooRoutes.DESTINATION_REQUESTS_QUEUE).find({ 'data-testid': 'pagination' }).exists(),
    ).toBe(true);
  });

  it('applies the sort to the status column in descending direction on both queues', () => {
    expect(
      GetMountedComponent(tooRoutes.MOVE_QUEUE).find({ 'data-testid': 'status' }).at(0).hasClass('sortAscending'),
    ).toBe(false);
    expect(
      GetMountedComponent(tooRoutes.DESTINATION_REQUESTS_QUEUE)
        .find({ 'data-testid': 'status' })
        .at(0)
        .hasClass('sortAscending'),
    ).toBe(true);
  });

  it('toggles the sort direction when clicked - MovesQueue', () => {
    const wrapper = GetMountedComponent(tooRoutes.MOVE_QUEUE);
    const statusHeading = wrapper.find({ 'data-testid': 'status' }).at(0);
    statusHeading.simulate('click');
    wrapper.update();

    expect(wrapper.find({ 'data-testid': 'status' }).at(0).hasClass('sortDescending')).toBe(false);

    statusHeading.simulate('click');
    wrapper.update();

    // asc should be applied
    expect(wrapper.find({ 'data-testid': 'status' }).at(0).hasClass('sortAscending')).toBe(true);
    expect(wrapper.find({ 'data-testid': 'status' }).at(0).hasClass('sortDescending')).toBe(false);

    const nameHeading = wrapper.find({ 'data-testid': 'customerName' }).at(0);
    nameHeading.simulate('click');
    wrapper.update();

    expect(wrapper.find({ 'data-testid': 'customerName' }).at(0).hasClass('sortAscending')).toBe(true);

    const assignedHeading = wrapper.find({ 'data-testid': 'assignedTo' }).at(0);
    assignedHeading?.simulate('click');
    wrapper.update();

    expect(wrapper.find({ 'data-testid': 'assignedTo' }).at(0).hasClass('sortAscending')).toBe(true);
  });
  it('toggles the sort direction when clicked - DestinationRequestsQueue', () => {
    const wrapper = GetMountedComponent(tooRoutes.DESTINATION_REQUESTS_QUEUE);

    const nameHeading = wrapper.find({ 'data-testid': 'customerName' }).at(0);
    nameHeading.simulate('click');
    wrapper.update();

    expect(wrapper.find({ 'data-testid': 'customerName' }).at(0).hasClass('sortAscending')).toBe(true);

    const assignedHeading = wrapper.find({ 'data-testid': 'assignedTo' }).at(0);
    assignedHeading?.simulate('click');
    wrapper.update();

    expect(wrapper.find({ 'data-testid': 'assignedTo' }).at(0).hasClass('sortAscending')).toBe(true);
  });

  it('filters the queue - MovesQueue', () => {
    const wrapper = GetMountedComponent(tooRoutes.MOVE_QUEUE);
    const input = wrapper.find(Select).at(0).find('input');
    input.simulate('keyDown', { key: 'ArrowDown', keyCode: 40 });
    input.simulate('keyDown', { key: 'Enter', keyCode: 13 });

    wrapper.update();
    expect(wrapper.find('[data-testid="multi-value-container"]').text()).toEqual('New move');
  });

  it('renders Search, Destination Queue and Move Queue tabs', () => {
    reactRouterDom.useParams.mockReturnValue({ queueType: generalRoutes.QUEUE_SEARCH_PATH });
    render(
      <MockProviders>
        <MoveQueue />
      </MockProviders>,
    );
    expect(screen.getByTestId('task-orders-tab-link')).toBeInTheDocument();
    expect(screen.getByTestId('search-tab-link')).toBeInTheDocument();
    expect(screen.getByText('Task Order Queue', { selector: 'span' })).toBeInTheDocument();
    expect(screen.getByText('Destination Requests Queue', { selector: 'span' })).toBeInTheDocument();
    expect(screen.getByText('Search', { selector: 'span' })).toBeInTheDocument();
  });
  it('renders TableQueue when Search tab is selected - MoveQueue', () => {
    reactRouterDom.useParams.mockReturnValue({ queueType: tooRoutes.MOVE_QUEUE });
    render(
      <MockProviders>
        <MoveQueue />
      </MockProviders>,
    );
    expect(screen.queryByTestId('table-queue')).toBeInTheDocument();
    expect(screen.queryByTestId('move-search')).not.toBeInTheDocument();
  });
  it('renders TableQueue when Search tab is selected - DestinationRequestsQueue', () => {
    reactRouterDom.useParams.mockReturnValue({ queueType: tooRoutes.DESTINATION_REQUESTS_QUEUE });
    render(
      <MockProviders>
        <MoveQueue />
      </MockProviders>,
    );
    expect(screen.queryByTestId('table-queue')).toBeInTheDocument();
    expect(screen.queryByTestId('move-search')).not.toBeInTheDocument();
  });
  it('has all options for searches', async () => {
    reactRouterDom.useParams.mockReturnValue({ queueType: generalRoutes.QUEUE_SEARCH_PATH });
    render(
      <MockProviders>
        <MoveQueue />
      </MockProviders>,
    );
    SEARCH_OPTIONS.forEach((option) => expect(screen.findByLabelText(option)));
  });
  it('Has all status options for move search', async () => {
    reactRouterDom.useParams.mockReturnValue({ queueType: generalRoutes.QUEUE_SEARCH_PATH });
    render(
      <MockProviders>
        <MoveQueue />
      </MockProviders>,
    );
    MOVE_STATUS_OPTIONS.forEach((option) => expect(screen.findByLabelText(option)));
  });

  it('Has all status options for move queue', async () => {
    reactRouterDom.useParams.mockReturnValue({ queueType: tooRoutes.MOVE_QUEUE });
    render(
      <MockProviders>
        <MoveQueue />
      </MockProviders>,
    );
    MOVE_STATUS_OPTIONS.forEach((option) => expect(screen.findByLabelText(option)));
  });
  it('Has all status options for destination requests queue', async () => {
    reactRouterDom.useParams.mockReturnValue({ queueType: tooRoutes.DESTINATION_REQUESTS_QUEUE });
    render(
      <MockProviders>
        <MoveQueue />
      </MockProviders>,
    );
    MOVE_STATUS_OPTIONS.forEach((option) => expect(screen.findByLabelText(option)));
  });

  it('renders a 404 if a bad route is provided', async () => {
    reactRouterDom.useParams.mockReturnValue({ queueType: 'BadRoute' });
    render(
      <MockProviders>
        <MoveQueue />
      </MockProviders>,
    );
    await expect(screen.getByText('Error - 404')).toBeInTheDocument();
    await expect(screen.getByText("We can't find the page you're looking for")).toBeInTheDocument();
  });

  it('renders a lock icon when move lock flag is on - MoveQueue', async () => {
    isBooleanFlagEnabled.mockResolvedValue(true);
    reactRouterDom.useParams.mockReturnValue({ queueType: tooRoutes.MOVE_QUEUE });
    render(
      <MockProviders>
        <MoveQueue />
      </MockProviders>,
    );
    await waitFor(() => {
      const lockIcon = screen.queryAllByTestId('lock-icon')[0];
      expect(lockIcon).toBeInTheDocument();
    });
  });
  it('renders a lock icon when move lock flag is on - DestinationRequestsQueue', async () => {
    isBooleanFlagEnabled.mockResolvedValue(true);
    reactRouterDom.useParams.mockReturnValue({ queueType: tooRoutes.DESTINATION_REQUESTS_QUEUE });
    render(
      <MockProviders>
        <MoveQueue />
      </MockProviders>,
    );
    await waitFor(() => {
      const lockIcon = screen.queryAllByTestId('lock-icon')[0];
      expect(lockIcon).toBeInTheDocument();
    });
  });

  it('does NOT render a lock icon when move lock flag is off - MoveQueue', async () => {
    isBooleanFlagEnabled.mockResolvedValue(false);
    reactRouterDom.useParams.mockReturnValue({ queueType: tooRoutes.MOVE_QUEUE });
    render(
      <MockProviders>
        <MoveQueue />
      </MockProviders>,
    );
    await await waitFor(() => {
      const lockIcon = screen.queryByTestId('lock-icon');
      expect(lockIcon).not.toBeInTheDocument();
    });
  });
  it('does NOT render a lock icon when move lock flag is off - DestinationRequestsQueue', async () => {
    isBooleanFlagEnabled.mockResolvedValue(false);
    reactRouterDom.useParams.mockReturnValue({ queueType: tooRoutes.DESTINATION_REQUESTS_QUEUE });
    render(
      <MockProviders>
        <MoveQueue />
      </MockProviders>,
    );
    await await waitFor(() => {
      const lockIcon = screen.queryByTestId('lock-icon');
      expect(lockIcon).not.toBeInTheDocument();
    });
  });

  it('renders an assigned column when the queue management flag is on - MoveQueue', async () => {
    reactRouterDom.useParams.mockReturnValue({ queueType: tooRoutes.MOVE_QUEUE });
    render(
      <MockProviders>
        <MoveQueue isQueueManagementFFEnabled />
      </MockProviders>,
    );
    await waitFor(() => {
      const assignedSelect = screen.queryAllByTestId('assigned-col')[0];
      expect(assignedSelect).toBeInTheDocument();
    });
  });
  it('renders an assigned column when the queue management flag is on - DestinationRequestsQueue', async () => {
    reactRouterDom.useParams.mockReturnValue({ queueType: tooRoutes.DESTINATION_REQUESTS_QUEUE });
    render(
      <MockProviders>
        <MoveQueue isQueueManagementFFEnabled />
      </MockProviders>,
    );
    await waitFor(() => {
      const assignedSelect = screen.queryAllByTestId('assigned-col')[0];
      expect(assignedSelect).toBeInTheDocument();
    });
  });

  it('renders an assigned column when the queue management flag is off - MoveQueue', async () => {
    reactRouterDom.useParams.mockReturnValue({ queueType: tooRoutes.MOVE_QUEUE });
    render(
      <MockProviders>
        <MoveQueue isQueueManagementFFEnabled={false} />
      </MockProviders>,
    );
    await waitFor(() => {
      const assignedSelect = screen.queryByTestId('assigned-col');
      expect(assignedSelect).not.toBeInTheDocument();
    });
  });
  it('renders an assigned column when the queue management flag is off - DestinationRequestsQueue', async () => {
    reactRouterDom.useParams.mockReturnValue({ queueType: tooRoutes.DESTINATION_REQUESTS_QUEUE });
    render(
      <MockProviders>
        <MoveQueue isQueueManagementFFEnabled={false} />
      </MockProviders>,
    );
    await waitFor(() => {
      const assignedSelect = screen.queryByTestId('assigned-col');
      expect(assignedSelect).not.toBeInTheDocument();
    });
  });
  it('renders an approval request type column when the feature flag is on - MoveQueue', async () => {
    reactRouterDom.useParams.mockReturnValue({ queueType: tooRoutes.MOVE_QUEUE });
    render(
      <MockProviders>
        <MoveQueue isApprovalRequestTypeFFEnabled />
      </MockProviders>,
    );
    await waitFor(() => {
      const approvalRequestColumn = screen.queryAllByTestId('approvalRequestTypes');
      expect(approvalRequestColumn).not.toHaveLength(0);
    });
  });
  it('renders an approval request type column when the feature flag is on - DestinationRequestsQueue', async () => {
    reactRouterDom.useParams.mockReturnValue({ queueType: tooRoutes.DESTINATION_REQUESTS_QUEUE });
    render(
      <MockProviders>
        <MoveQueue isApprovalRequestTypeFFEnabled />
      </MockProviders>,
    );
    await waitFor(() => {
      const approvalRequestColumn = screen.queryAllByTestId('approvalRequestTypes');
      expect(approvalRequestColumn).not.toHaveLength(0);
    });
  });
  it('does not render an approval request type column when the feature flag is off - MoveQueue', async () => {
    reactRouterDom.useParams.mockReturnValue({ queueType: tooRoutes.MOVE_QUEUE });
    render(
      <MockProviders>
        <MoveQueue isApprovalRequestTypeFFEnable={false} />
      </MockProviders>,
    );
    await waitFor(() => {
      const approvalRequestColumn = screen.queryByTestId('approvalRequestTypes');
      expect(approvalRequestColumn).not.toBeInTheDocument();
    });
  });
  it('does not render an approval request type column when the feature flag is off - DestinationRequestsQueue', async () => {
    reactRouterDom.useParams.mockReturnValue({ queueType: tooRoutes.DESTINATION_REQUESTS_QUEUE });
    render(
      <MockProviders>
        <MoveQueue isApprovalRequestTypeFFEnable={false} />
      </MockProviders>,
    );
    await waitFor(() => {
      const approvalRequestColumn = screen.queryAllByTestId('approvalRequestTypes');
      expect(approvalRequestColumn).toHaveLength(0);
    });
  });
});

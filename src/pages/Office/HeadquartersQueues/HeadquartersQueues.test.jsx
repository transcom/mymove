import React from 'react';
import Select from 'react-select';
import { mount } from 'enzyme';
import * as reactRouterDom from 'react-router-dom';
import { render, screen, waitFor } from '@testing-library/react';

import HeadquartersQueue from './HeadquartersQueues';

import { MockProviders } from 'testUtils';
import { MOVE_STATUS_OPTIONS } from 'constants/queues';
import { generalRoutes, hqRoutes } from 'constants/routes';
import { isBooleanFlagEnabled } from 'utils/featureFlags';
import SelectedGblocProvider from 'components/Office/GblocSwitcher/SelectedGblocProvider';

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
              edipi: '555555555',
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
          },
          {
            id: 'move2',
            customer: {
              agency: 'MARINES',
              first_name: 'test another first',
              last_name: 'test another last',
              edipi: '4444444444',
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
            appearedInTooAt: '2023-02-12T00:00:00.000Z',
          },
        ],
      },
    };
  },
  useServicesCounselingQueueQueries: () => {
    return {
      isLoading: false,
      isError: false,
      queueResult: {
        totalCount: 2,
        data: [
          {
            id: 'move1',
            customer: {
              agency: 'ARMY',
              first_name: 'test first',
              last_name: 'test last',
              edipi: '555555555',
            },
            locator: 'AB5PC',
            requestedMoveDate: '2021-03-01T00:00:00.000Z',
            submittedAt: '2021-01-31T00:00:00.000Z',
            status: 'NEEDS SERVICE COUNSELING',
            originDutyLocation: {
              name: 'Area 51',
            },
            originGBLOC: 'LKNQ',
          },
          {
            id: 'move2',
            customer: {
              agency: 'COAST_GUARD',
              first_name: 'test another first',
              last_name: 'test another last',
              edipi: '4444444444',
            },
            locator: 'T12AR',
            requestedMoveDate: '2021-04-15T00:00:00.000Z',
            submittedAt: '2021-01-01T00:00:00.000Z',
            status: 'NEEDS SERVICE COUNSELING',
            originDutyLocation: {
              name: 'Los Alamos',
            },
            originGBLOC: 'LKNQ',
          },
          {
            id: 'move3',
            customer: {
              agency: 'MARINES',
              first_name: 'test third first',
              last_name: 'test third last',
              edipi: '4444444444',
            },
            locator: 'T12MP',
            requestedMoveDate: '2021-04-15T00:00:00.000Z',
            submittedAt: '2021-01-01T00:00:00.000Z',
            status: 'NEEDS SERVICE COUNSELING',
            originDutyLocation: {
              name: 'Denver, 80136',
            },
            originGBLOC: 'LKNQ',
          },
        ],
      },
    };
  },
  usePaymentRequestQueueQueries: () => {
    return {
      isLoading: false,
      isError: false,
      queueResult: {
        totalCount: 0,
        data: [],
      },
    };
  },
  useServicesCounselingQueuePPMQueries: () => {
    return {
      isLoading: false,
      isError: false,
      queueResult: {
        totalCount: 0,
        data: [],
      },
    };
  },
}));

const GetMountedComponent = (queueTypeToMount) => {
  reactRouterDom.useParams.mockReturnValue({ queueType: queueTypeToMount });
  const wrapper = mount(
    <MockProviders>
      <HeadquartersQueue />
    </MockProviders>,
  );
  return wrapper;
};

const SEARCH_OPTIONS = ['Move Code', 'edipi', 'Customer Name', 'Payment Request Number'];

describe('HeadquartersQueue', () => {
  afterEach(() => {
    jest.restoreAllMocks();
  });

  it('should render the h1', () => {
    expect(GetMountedComponent(hqRoutes.MOVE_QUEUE).find('h1').text()).toBe('All moves (2)');
  });

  it('should render the table', () => {
    expect(GetMountedComponent(hqRoutes.MOVE_QUEUE).find('Table').exists()).toBe(true);
  });

  it('should format the column data', () => {
    const moves = GetMountedComponent(hqRoutes.MOVE_QUEUE).find('tbody tr');

    const firstMove = moves.at(0);
    expect(firstMove.find({ 'data-testid': 'lastName-0' }).text()).toBe('test last, test first');
    expect(firstMove.find({ 'data-testid': 'edipi-0' }).text()).toBe('555555555');
    expect(firstMove.find({ 'data-testid': 'status-0' }).text()).toBe('New move');
    expect(firstMove.find({ 'data-testid': 'locator-0' }).text()).toBe('AB5P');
    expect(firstMove.find({ 'data-testid': 'branch-0' }).text()).toBe('Air Force');
    expect(firstMove.find({ 'data-testid': 'shipmentsCount-0' }).text()).toBe('2');
    expect(firstMove.find({ 'data-testid': 'originDutyLocation-0' }).text()).toBe('Area 51');
    expect(firstMove.find({ 'data-testid': 'originGBLOC-0' }).text()).toBe('EEEE');
    expect(firstMove.find({ 'data-testid': 'requestedMoveDate-0' }).text()).toBe('10 Feb 2023');
    expect(firstMove.find({ 'data-testid': 'appearedInTooAt-0' }).text()).toBe('10 Feb 2023');

    const secondMove = moves.at(1);
    expect(secondMove.find({ 'data-testid': 'lastName-1' }).text()).toBe('test another last, test another first');
    expect(secondMove.find({ 'data-testid': 'edipi-1' }).text()).toBe('4444444444');
    expect(secondMove.find({ 'data-testid': 'status-1' }).text()).toBe('Move approved');
    expect(secondMove.find({ 'data-testid': 'locator-1' }).text()).toBe('T12A');
    expect(secondMove.find({ 'data-testid': 'branch-1' }).text()).toBe('Marine Corps');
    expect(secondMove.find({ 'data-testid': 'shipmentsCount-1' }).text()).toBe('1');
    expect(secondMove.find({ 'data-testid': 'originDutyLocation-1' }).text()).toBe('Los Alamos');
    expect(secondMove.find({ 'data-testid': 'originGBLOC-1' }).text()).toBe('EEEE');
    expect(secondMove.find({ 'data-testid': 'requestedMoveDate-1' }).text()).toBe('12 Feb 2023');
    expect(secondMove.find({ 'data-testid': 'appearedInTooAt-1' }).text()).toBe('12 Feb 2023');
  });

  it('should render the pagination component', () => {
    expect(GetMountedComponent(hqRoutes.MOVE_QUEUE).find({ 'data-testid': 'pagination' }).exists()).toBe(true);
  });

  it('applies the sort to the status column in descending direction', () => {
    expect(
      GetMountedComponent(hqRoutes.MOVE_QUEUE).find({ 'data-testid': 'status' }).at(0).hasClass('sortAscending'),
    ).toBe(true);
  });

  it('toggles the sort direction when clicked', () => {
    const wrapper = GetMountedComponent(hqRoutes.MOVE_QUEUE);
    const statusHeading = wrapper.find({ 'data-testid': 'status' }).at(0);

    statusHeading.simulate('click');
    GetMountedComponent(hqRoutes.MOVE_QUEUE).update();
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
    const wrapper = GetMountedComponent(hqRoutes.MOVE_QUEUE);
    const input = wrapper.find(Select).at(0).find('input');
    input.simulate('keyDown', { key: 'ArrowDown', keyCode: 40 });
    input.simulate('keyDown', { key: 'Enter', keyCode: 13 });

    wrapper.update();
    expect(wrapper.find('[data-testid="multi-value-container"]').text()).toEqual('New move');
  });

  it('renders navigation tabs', () => {
    reactRouterDom.useParams.mockReturnValue({ queueType: generalRoutes.QUEUE_SEARCH_PATH });
    render(
      <reactRouterDom.BrowserRouter>
        <SelectedGblocProvider>
          <HeadquartersQueue />
        </SelectedGblocProvider>
      </reactRouterDom.BrowserRouter>,
    );
    expect(screen.getByTestId('task-order-queue-tab-link')).toBeInTheDocument();
    expect(screen.getByTestId('payment-request-queue-tab-link')).toBeInTheDocument();
    expect(screen.getByTestId('counseling-queue-tab-link')).toBeInTheDocument();
    expect(screen.getByTestId('closeout-queue-tab-link')).toBeInTheDocument();
    expect(screen.getByTestId('move-search-tab-link')).toBeInTheDocument();
    expect(screen.getByTestId('customer-search-tab-link')).toBeInTheDocument();

    expect(screen.getByText('Task Order Queue', { selector: 'span' })).toBeInTheDocument();
    expect(screen.getByText('Payment Request Queue', { selector: 'span' })).toBeInTheDocument();
    expect(screen.getByText('Counseling Queue', { selector: 'span' })).toBeInTheDocument();
    expect(screen.getByText('PPM Closeout Queue', { selector: 'span' })).toBeInTheDocument();
    expect(screen.getByText('Move Search', { selector: 'span' })).toBeInTheDocument();
    expect(screen.getByText('Customer Search', { selector: 'span' })).toBeInTheDocument();
  });

  it('renders Move Search when Move Search is selected', () => {
    reactRouterDom.useParams.mockReturnValue({ queueType: generalRoutes.QUEUE_SEARCH_PATH });
    render(
      <reactRouterDom.BrowserRouter>
        <SelectedGblocProvider>
          <HeadquartersQueue />
        </SelectedGblocProvider>
      </reactRouterDom.BrowserRouter>,
    );
    expect(screen.queryByTestId('table-queue')).not.toBeInTheDocument();
    expect(screen.queryByTestId('move-search-form')).not.toBeInTheDocument();
    expect(screen.getByText('Search for a Move', { selector: 'h1' })).toBeInTheDocument();
  });

  it('renders Counseling Queue when Counseling Queue is selected', () => {
    reactRouterDom.useParams.mockReturnValue({ queueType: hqRoutes.COUNSELING_QUEUE });
    render(
      <reactRouterDom.BrowserRouter>
        <SelectedGblocProvider>
          <HeadquartersQueue />
        </SelectedGblocProvider>
      </reactRouterDom.BrowserRouter>,
    );
    expect(screen.queryByTestId('table-queue')).toBeInTheDocument();
    expect(screen.getByText(/Moves/, { selector: 'h1' })).toBeInTheDocument();
  });

  it('renders PPM Closeout Queue when PPM Closeout Queue is selected', () => {
    reactRouterDom.useParams.mockReturnValue({ queueType: hqRoutes.CLOSEOUT_QUEUE });
    render(
      <reactRouterDom.BrowserRouter>
        <SelectedGblocProvider>
          <HeadquartersQueue />
        </SelectedGblocProvider>
      </reactRouterDom.BrowserRouter>,
    );
    expect(screen.queryByTestId('table-queue')).toBeInTheDocument();
    expect(screen.getByText('Closeout initiated', { selector: 'th' })).toBeInTheDocument();
    expect(screen.getByText(/Moves/, { selector: 'h1' })).toBeInTheDocument();
  });

  it('renders Payment Requests Queue when Payment Requests Queue is selected', () => {
    reactRouterDom.useParams.mockReturnValue({ queueType: hqRoutes.PAYMENT_REQUEST_QUEUE });
    render(
      <reactRouterDom.BrowserRouter>
        <SelectedGblocProvider>
          <HeadquartersQueue />
        </SelectedGblocProvider>
      </reactRouterDom.BrowserRouter>,
    );
    expect(screen.queryByTestId('table-queue')).toBeInTheDocument();
    expect(screen.getByText(/Payment requests/, { selector: 'h1' })).toBeInTheDocument();
  });

  it('has all options for searches', async () => {
    reactRouterDom.useParams.mockReturnValue({ queueType: generalRoutes.QUEUE_SEARCH_PATH });
    render(
      <reactRouterDom.BrowserRouter>
        <SelectedGblocProvider>
          <HeadquartersQueue />
        </SelectedGblocProvider>
      </reactRouterDom.BrowserRouter>,
    );
    SEARCH_OPTIONS.forEach((option) => expect(screen.findByLabelText(option)));
  });

  it('Has all status options for move search', async () => {
    reactRouterDom.useParams.mockReturnValue({ queueType: generalRoutes.QUEUE_SEARCH_PATH });
    render(
      <reactRouterDom.BrowserRouter>
        <SelectedGblocProvider>
          <HeadquartersQueue />
        </SelectedGblocProvider>
      </reactRouterDom.BrowserRouter>,
    );
    MOVE_STATUS_OPTIONS.forEach((option) => expect(screen.findByLabelText(option)));
  });

  it('Has all status options for move queue', async () => {
    reactRouterDom.useParams.mockReturnValue({ queueType: hqRoutes.MOVE_QUEUE });
    render(
      <reactRouterDom.BrowserRouter>
        <SelectedGblocProvider>
          <HeadquartersQueue />
        </SelectedGblocProvider>
      </reactRouterDom.BrowserRouter>,
    );
    MOVE_STATUS_OPTIONS.forEach((option) => expect(screen.findByLabelText(option)));
  });

  it('renders a 404 if a bad route is provided', async () => {
    reactRouterDom.useParams.mockReturnValue({ queueType: 'BadRoute' });
    render(
      <reactRouterDom.BrowserRouter>
        <SelectedGblocProvider>
          <HeadquartersQueue />
        </SelectedGblocProvider>
      </reactRouterDom.BrowserRouter>,
    );
    await expect(screen.getByText('Error - 404')).toBeInTheDocument();
    await expect(screen.getByText("We can't find the page you're looking for")).toBeInTheDocument();
  });

  it('renders a lock icon when move lock flag is on', async () => {
    isBooleanFlagEnabled.mockResolvedValue(true);
    reactRouterDom.useParams.mockReturnValue({ queueType: hqRoutes.MOVE_QUEUE });
    render(
      <reactRouterDom.BrowserRouter>
        <SelectedGblocProvider>
          <HeadquartersQueue />
        </SelectedGblocProvider>
      </reactRouterDom.BrowserRouter>,
    );
    await waitFor(() => {
      const lockIcon = screen.queryByTestId('lock-icon');
      expect(lockIcon).toBeInTheDocument();
    });
  });

  it('does NOT render a lock icon when move lock flag is off', async () => {
    isBooleanFlagEnabled.mockResolvedValue(false);
    reactRouterDom.useParams.mockReturnValue({ queueType: hqRoutes.MOVE_QUEUE });
    render(
      <reactRouterDom.BrowserRouter>
        <SelectedGblocProvider>
          <HeadquartersQueue />
        </SelectedGblocProvider>
      </reactRouterDom.BrowserRouter>,
    );
    await await waitFor(() => {
      const lockIcon = screen.queryByTestId('lock-icon');
      expect(lockIcon).not.toBeInTheDocument();
    });
  });
});

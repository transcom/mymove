import React from 'react';
import { mount } from 'enzyme';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import { OFFICE_TABLE_QUEUE_SESSION_STORAGE_ID } from '../../../components/Table/utils';

import ServicesCounselingQueue from './ServicesCounselingQueue';

import { useUserQueries, useServicesCounselingQueueQueries, useServicesCounselingQueuePPMQueries } from 'hooks/queries';
import { MockProviders, MockRouterProvider } from 'testUtils';
import { MOVE_STATUSES } from 'shared/constants';
import SERVICE_MEMBER_AGENCIES from 'content/serviceMemberAgencies';
import { servicesCounselingRoutes } from 'constants/routes';
import MoveSearchForm from 'components/MoveSearchForm/MoveSearchForm';
import { isBooleanFlagEnabled } from 'utils/featureFlags';

jest.mock('hooks/queries', () => ({
  useUserQueries: jest.fn(),
  useServicesCounselingQueueQueries: jest.fn(),
  useServicesCounselingQueuePPMQueries: jest.fn(),
}));

jest.mock('utils/featureFlags', () => ({
  ...jest.requireActual('utils/featureFlags'),
  isBooleanFlagEnabled: jest.fn().mockImplementation(() => Promise.resolve(false)),
}));

const localStorageMock = (() => {
  let store = {};

  return {
    getItem(key) {
      return store[key] || null;
    },
    setItem(key, value) {
      store[key] = value;
    },
    removeItem(key) {
      delete store[key];
    },
    clear() {
      store = {};
    },
  };
})();

Object.defineProperty(window, 'sessionStorage', {
  value: localStorageMock,
});

const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  Navigate: (props) => {
    mockNavigate(props?.to);
    return null;
  },
}));
const pagePath = '/:queueType/*';
const serviceCounselorUser = {
  isLoading: false,
  isError: false,
  data: {
    office_user: { transportation_office: { gbloc: 'LKNQ' } },
  },
};

const serviceCounselorUserForCloseout = {
  isLoading: false,
  isError: false,
  data: {
    office_user: { transportation_office: { gbloc: 'TVCB' } },
  },
};

const emptyServiceCounselingMoves = {
  isLoading: false,
  isError: false,
  queueResult: {
    totalCount: 0,
    data: [],
  },
};

const needsCounselingMoves = {
  isLoading: false,
  isError: false,
  queueResult: {
    totalCount: 3,
    data: [
      {
        id: 'move1',
        customer: {
          agency: SERVICE_MEMBER_AGENCIES.ARMY,
          first_name: 'test first',
          last_name: 'test last',
          eipid: '555555555',
        },
        locator: 'AB5PC',
        requestedMoveDate: '2021-03-01T00:00:00.000Z',
        submittedAt: '2021-01-31T00:00:00.000Z',
        status: MOVE_STATUSES.NEEDS_SERVICE_COUNSELING,
        originDutyLocation: {
          name: 'Area 51',
        },
        originGBLOC: 'LKNQ',
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
      {
        id: 'move2',
        customer: {
          agency: SERVICE_MEMBER_AGENCIES.COAST_GUARD,
          first_name: 'test another first',
          last_name: 'test another last',
          eipid: '4444444444',
          emplid: '4521567',
        },
        locator: 'T12AR',
        requestedMoveDate: '2021-04-15T00:00:00.000Z',
        submittedAt: '2021-01-01T00:00:00.000Z',
        status: MOVE_STATUSES.NEEDS_SERVICE_COUNSELING,
        originDutyLocation: {
          name: 'Los Alamos',
        },
        originGBLOC: 'LKNQ',
        counselingOffice: '',
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
          agency: SERVICE_MEMBER_AGENCIES.MARINES,
          first_name: 'test third first',
          last_name: 'test third last',
          eipid: '4444444444',
        },
        locator: 'T12MP',
        requestedMoveDate: '2021-04-15T00:00:00.000Z',
        submittedAt: '2021-01-01T00:00:00.000Z',
        status: MOVE_STATUSES.NEEDS_SERVICE_COUNSELING,
        originDutyLocation: {
          name: 'Denver, 80136',
        },
        originGBLOC: 'LKNQ',
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
    ],
  },
};

const serviceCounselingCompletedMoves = {
  isLoading: false,
  isError: false,
  queueResult: {
    totalCount: 2,
    data: [
      {
        id: 'move1',
        customer: {
          agency: SERVICE_MEMBER_AGENCIES.ARMY,
          first_name: 'test first',
          last_name: 'test last',
          eipid: '555555555',
        },
        locator: 'AB5PC',
        requestedMoveDate: '2021-03-01T00:00:00.000Z',
        submittedAt: '2021-01-31T00:00:00.000Z',
        status: MOVE_STATUSES.SERVICE_COUNSELING_COMPLETED,
        originDutyLocation: {
          name: 'Area 51',
        },
        originGBLOC: 'LKNQ',
        assignedTo: {
          id: 'exampleId1',
          firstname: 'Jimmy',
          lastname: 'John',
        },
      },
      {
        id: 'move2',
        customer: {
          agency: SERVICE_MEMBER_AGENCIES.COAST_GUARD,
          first_name: 'test another first',
          last_name: 'test another last',
          eipid: '4444444444',
        },
        locator: 'T12AR',
        requestedMoveDate: '2021-04-15T00:00:00.000Z',
        submittedAt: '2021-01-01T00:00:00.000Z',
        status: MOVE_STATUSES.SERVICE_COUNSELING_COMPLETED,
        originDutyLocation: {
          name: 'Los Alamos',
        },
        originGBLOC: 'LKNQ',
        counselingOffice: '67592323-fc7e-4b35-83a7-57faa53b7acf',
        assignedTo: {
          id: 'exampleId1',
          firstname: 'Jimmy',
          lastname: 'John',
        },
      },
    ],
  },
};

afterEach(() => {
  jest.resetAllMocks();
});

describe('ServicesCounselingQueue', () => {
  describe('no moves in service counseling statuses', () => {
    useUserQueries.mockReturnValue(serviceCounselorUser);
    useServicesCounselingQueueQueries.mockReturnValue(emptyServiceCounselingMoves);
    const wrapper = mount(
      <MockRouterProvider path={pagePath} params={{ queueType: 'counseling' }}>
        <ServicesCounselingQueue isQueueManagementFFEnabled />
      </MockRouterProvider>,
    );

    it('displays move header with count', () => {
      expect(wrapper.find('h1').text()).toBe('Moves (0)');
    });

    it('renders the table', () => {
      expect(wrapper.find('Table').exists()).toBe(true);
    });

    it('no move rows are rendered', () => {
      expect(wrapper.find('tbody tr').length).toBe(0);
    });
  });

  describe('Service Counselor', () => {
    useUserQueries.mockReturnValue(serviceCounselorUser);
    useServicesCounselingQueueQueries.mockReturnValue(needsCounselingMoves);
    isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(true));
    const wrapper = mount(
      <MockRouterProvider path={pagePath} params={{ queueType: 'counseling' }}>
        <ServicesCounselingQueue isQueueManagementFFEnabled />
      </MockRouterProvider>,
    );
    render(
      <MockRouterProvider path={pagePath} params={{ queueType: 'counseling' }}>
        <ServicesCounselingQueue isQueueManagementFFEnabled />
      </MockRouterProvider>,
    );
    it('displays move header with needs service counseling count', () => {
      expect(wrapper.find('h1').text()).toBe('Moves (3)');
    });

    it('renders the table', () => {
      expect(wrapper.find('Table').exists()).toBe(true);
    });

    it('renders the pagination component', () => {
      expect(wrapper.find({ 'data-testid': 'pagination' }).exists()).toBe(true);
    });

    it('formats the move data in rows', () => {
      const moves = wrapper.find('tbody tr');
      const firstMove = moves.at(0);
      expect(firstMove.find('td.customerName').text()).toBe('test last, test first');
      expect(firstMove.find('td.eipid').text()).toBe('555555555');
      expect(firstMove.find('td.locator').text()).toBe('AB5PC');
      expect(firstMove.find('td.status').text()).toBe('Needs counseling');
      expect(firstMove.find('td.requestedMoveDate').text()).toBe('01 Mar 2021');
      expect(firstMove.find('td.submittedAt').text()).toBe('31 Jan 2021');
      expect(firstMove.find('td.branch').text()).toBe('Army');
      expect(firstMove.find('td.originGBLOC').text()).toBe('LKNQ');
      expect(firstMove.find('td.originDutyLocation').text()).toBe('Area 51');
      expect(firstMove.find('td.assignedTo').text()).toBe('John, Jimmy');

      const secondMove = moves.at(1);
      expect(secondMove.find('td.customerName').text()).toBe('test another last, test another first');
      expect(secondMove.find('td.eipid').text()).toBe('4444444444');
      expect(secondMove.find('td.emplid').text()).toBe('4521567');
      expect(secondMove.find('td.locator').text()).toBe('T12AR');
      expect(secondMove.find('td.status').text()).toBe('Needs counseling');
      expect(secondMove.find('td.requestedMoveDate').text()).toBe('15 Apr 2021');
      expect(secondMove.find('td.submittedAt').text()).toBe('01 Jan 2021');
      expect(secondMove.find('td.branch').text()).toBe('Coast Guard');
      expect(secondMove.find('td.originGBLOC').text()).toBe('LKNQ');
      expect(secondMove.find('td.originDutyLocation').text()).toBe('Los Alamos');
      expect(secondMove.find('td.assignedTo').text()).toBe('Denver, John');

      const thirdMove = moves.at(2);
      expect(thirdMove.find('td.customerName').text()).toBe('test third last, test third first');
      expect(thirdMove.find('td.eipid').text()).toBe('4444444444');
      expect(thirdMove.find('td.locator').text()).toBe('T12MP');
      expect(thirdMove.find('td.status').text()).toBe('Needs counseling');
      expect(thirdMove.find('td.requestedMoveDate').text()).toBe('15 Apr 2021');
      expect(thirdMove.find('td.submittedAt').text()).toBe('01 Jan 2021');
      expect(thirdMove.find('td.branch').text()).toBe('Marine Corps');
      expect(thirdMove.find('td.originGBLOC').text()).toBe('LKNQ');
      expect(thirdMove.find('td.originDutyLocation').text()).toBe('Denver, 80136');
      expect(thirdMove.find('td.assignedTo').text()).toBe('John, Jimmy');
    });

    it('sorts by submitted at date ascending by default', () => {
      expect(wrapper.find('th[data-testid="submittedAt"][role="columnheader"]').hasClass('sortAscending')).toBe(true);
    });

    it('allows sorting on certain columns', () => {
      expect(wrapper.find('th[data-testid="customerName"][role="columnheader"]').prop('onClick')).not.toBe(undefined);
      expect(wrapper.find('th[data-testid="dodID"][role="columnheader"]').prop('onClick')).not.toBe(undefined);
      expect(wrapper.find('th[data-testid="emplid"][role="columnheader"]').prop('onClick')).not.toBe(undefined);
      expect(wrapper.find('th[data-testid="locator"][role="columnheader"]').prop('onClick')).not.toBe(undefined);
      expect(wrapper.find('th[data-testid="requestedMoveDate"][role="columnheader"]').prop('onClick')).not.toBe(
        undefined,
      );
      expect(wrapper.find('th[data-testid="submittedAt"][role="columnheader"]').prop('onClick')).not.toBe(undefined);
      expect(wrapper.find('th[data-testid="branch"][role="columnheader"]').prop('onClick')).not.toBe(undefined);
      expect(wrapper.find('th[data-testid="originDutyLocation"][role="columnheader"]').prop('onClick')).not.toBe(
        undefined,
      );
      expect(wrapper.find('th[data-testid="assignedTo"][role="columnheader"]').prop('onClick')).not.toBe(undefined);
    });

    it('disables sort by for origin GBLOC and status columns', () => {
      expect(wrapper.find('th[data-testid="originGBLOC"][role="columnheader"]').prop('onClick')).toBe(undefined);
      expect(wrapper.find('th[data-testid="status"][role="columnheader"]').prop('onClick')).toBe(undefined);
    });

    it('omits filter input element for origin GBLOC column', () => {
      expect(wrapper.find('th[data-testid="originGBLOC"] input').exists()).toBe(false);
    });
  });

  describe('verify cached filters are displayed in respective filter column header on page reload -  Service Counselor', () => {
    window.sessionStorage.setItem(
      OFFICE_TABLE_QUEUE_SESSION_STORAGE_ID,
      '{"counseling":{"filters":[{"id":"customerName","value":"Spacemen"},{"id":"edipi","value":"7232607949"},{"id":"locator","value":"PPMADD"},{"id":"requestedMoveDate","value":"2024-06-21"},{"id":"submittedAt","value":"2024-06-20T04:00:00+00:00"},{"id":"branch","value":"ARMY"},{"id":"originDutyLocation","value":"12345"}], "sortParam":[{"id":"customerName","desc":false}], "page":3,"pageSize":10}}',
    );
    useUserQueries.mockReturnValue(serviceCounselorUser);

    const moves = JSON.parse(JSON.stringify(needsCounselingMoves));

    for (let i = 0; i < 30; i += 1) {
      moves.queueResult.data.push({
        id: `move${moves.queueResult.data.length}${1}`,
        customer: {
          agency: SERVICE_MEMBER_AGENCIES.ARMY,
          first_name: 'test first',
          last_name: 'test last',
          eipid: '555555555',
        },
        locator: 'AB5PC',
        requestedMoveDate: '2021-03-01T00:00:00.000Z',
        submittedAt: '2021-01-31T00:00:00.000Z',
        status: MOVE_STATUSES.NEEDS_SERVICE_COUNSELING,
        originDutyLocation: {
          name: 'Area 51',
        },
        originGBLOC: 'LKNQ',
      });
    }
    moves.queueResult.totalCount = moves.queueResult.data.length;

    useServicesCounselingQueueQueries.mockReturnValue(moves);
    const wrapper = mount(
      <MockProviders path={pagePath} params={{ queueType: 'counseling' }}>
        <ServicesCounselingQueue />
      </MockProviders>,
    );

    // Verify controls are using cached data on load.
    // If any of these fail check setup data window.sessionStorage.setItem()
    expect(wrapper.find('th[data-testid="customerName"] input').instance().value).toBe('Spacemen');
    expect(wrapper.find('th[data-testid="edipi"] input').instance().value).toBe('7232607949');
    expect(wrapper.find('th[data-testid="locator"] input').instance().value).toBe('PPMADD');
    expect(wrapper.find('th[data-testid="requestedMoveDate"] input').instance().value).toBe('21 Jun 2024');
    expect(wrapper.find('th[data-testid="submittedAt"] input').instance().value).toBe('20 Jun 2024');
    expect(wrapper.find('th[data-testid="originDutyLocation"] input').instance().value).toBe('12345');
    expect(wrapper.find('th[data-testid="branch"] select').instance().value).toBe('ARMY');
    expect(wrapper.find('[data-testid="pagination"] select[id="table-rows-per-page"]').instance().value).toBe('10');
    expect(wrapper.find('[data-testid="pagination"] select[id="table-pagination"]').instance().value).toBe('2');
    expect(wrapper.find('th[data-testid="customerName"][role="columnheader"]').instance().className).toBe(
      'sortAscending',
    );
  });

  describe('filter sessionStorage filters - no cache-  Service Counselor', () => {
    window.sessionStorage.clear();
    useUserQueries.mockReturnValue(serviceCounselorUser);
    useServicesCounselingQueueQueries.mockReturnValue(needsCounselingMoves);
    const wrapper = mount(
      <MockProviders path={pagePath} params={{ queueType: 'counseling' }}>
        <ServicesCounselingQueue />
      </MockProviders>,
    );
    expect(wrapper.find('th[data-testid="customerName"] input').instance().value).toBe('');
    expect(wrapper.find('th[data-testid="edipi"] input').instance().value).toBe('');
    expect(wrapper.find('th[data-testid="locator"] input').instance().value).toBe('');
    expect(wrapper.find('th[data-testid="requestedMoveDate"] input').instance().value).toBe('');
    expect(wrapper.find('th[data-testid="submittedAt"] input').instance().value).toBe('');
    expect(wrapper.find('th[data-testid="originDutyLocation"] input').instance().value).toBe('');
    expect(wrapper.find('th[data-testid="branch"] select').instance().value).toBe('');
  });

  describe('service counseling tab routing', () => {
    it.each([
      ['counseling', servicesCounselingRoutes.BASE_QUEUE_COUNSELING_PATH, serviceCounselorUser],
      ['closeout', servicesCounselingRoutes.BASE_QUEUE_CLOSEOUT_PATH, serviceCounselorUserForCloseout],
    ])(
      'a %s user accessing the SC queue default path gets redirected appropriately to %s',
      (userDescription, expectedPath, user) => {
        //  ['closeout', servicesCounselingRoutes.DEFAULT_QUEUE_PATH, false, serviceCounselorUserForCloseout],

        useUserQueries.mockReturnValue(user);
        useServicesCounselingQueueQueries.mockReturnValue(serviceCounselingCompletedMoves);
        useServicesCounselingQueuePPMQueries.mockReturnValue(emptyServiceCounselingMoves);
        render(
          <MockProviders>
            <ServicesCounselingQueue />
          </MockProviders>,
        );

        expect(mockNavigate).toHaveBeenCalledTimes(1);
        expect(mockNavigate).toHaveBeenCalledWith(expectedPath);
      },
    );

    it.each([
      ['counselor', servicesCounselingRoutes.QUEUE_COUNSELING_PATH, 'counseling', serviceCounselorUser],
      ['counselor', servicesCounselingRoutes.QUEUE_CLOSEOUT_PATH, 'closeout', serviceCounselorUser],
      ['closeout', servicesCounselingRoutes.QUEUE_COUNSELING_PATH, 'counseling', serviceCounselorUserForCloseout],
      ['closeout', servicesCounselingRoutes.QUEUE_CLOSEOUT_PATH, 'closeout', serviceCounselorUserForCloseout],
    ])('a %s user accessing path "%s"', async (userDescription, queueType, showsCounselingTab, user) => {
      isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(true));
      useUserQueries.mockReturnValue(user);
      useServicesCounselingQueueQueries.mockReturnValue(serviceCounselingCompletedMoves);
      useServicesCounselingQueuePPMQueries.mockReturnValue(emptyServiceCounselingMoves);
      render(
        <MockProviders path={pagePath} params={{ queueType }}>
          <ServicesCounselingQueue isQueueManagementFFEnabled />
        </MockProviders>,
      );

      await waitFor(() => {
        if (showsCounselingTab === 'counseling') {
          // Make sure "Counseling" is the active tab.
          const counselingActive = screen.getByText('Counseling Queue', { selector: '.usa-current .tab-title' });
          expect(counselingActive).toBeInTheDocument();

          // Check for the "Counseling" columns.
          expect(screen.getByText(/Status/)).toBeInTheDocument();
          expect(screen.getAllByText(/Requested move date/)[0]).toBeInTheDocument();
          expect(screen.getAllByText(/Date submitted/)[0]).toBeInTheDocument();
          expect(screen.getByText(/Origin GBLOC/)).toBeInTheDocument();
          expect(screen.getByText(/Assigned/)).toBeInTheDocument();
        } else if (showsCounselingTab === 'closeout') {
          // Make sure "PPM Closeout" is the active tab.
          const ppmCloseoutActive = screen.getByText('PPM Closeout Queue', { selector: '.usa-current .tab-title' });
          expect(ppmCloseoutActive).toBeInTheDocument();

          // Check for the "PPM Closeout" columns.
          expect(screen.getByText(/Closeout initiated/)).toBeInTheDocument();
          expect(screen.getByText(/PPM closeout location/)).toBeInTheDocument();
          expect(screen.getByText(/Full or partial PPM/)).toBeInTheDocument();
          expect(screen.getByText(/Destination duty location/)).toBeInTheDocument();
          expect(screen.getByText(/Status/)).toBeInTheDocument();
        } else {
          // Check for the "Search" tab
          const searchActive = screen.getByText('Search', { selector: '.usa-current .tab-title' });
          expect(searchActive).toBeInTheDocument();
          expect(MoveSearchForm).toBeInTheDocument();
          userEvent.type(screen.getByLabelText('Search'), 'Joe');
          const addCustomer = screen.getByText('Add Customer', { selector: '.usa-current .tab-title' });
          expect(addCustomer).toBeInTheDocument();
        }
      });
    });
  });
});

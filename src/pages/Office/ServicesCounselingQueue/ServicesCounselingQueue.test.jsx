import React from 'react';
import { mount } from 'enzyme';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import ServicesCounselingQueue from './ServicesCounselingQueue';

import { useUserQueries, useServicesCounselingQueueQueries, useServicesCounselingQueuePPMQueries } from 'hooks/queries';
import { MockProviders, MockRouterProvider } from 'testUtils';
import { MOVE_STATUSES } from 'shared/constants';
import SERVICE_MEMBER_AGENCIES from 'content/serviceMemberAgencies';
import { servicesCounselingRoutes } from 'constants/routes';
import MoveSearchForm from 'components/MoveSearchForm/MoveSearchForm';

jest.mock('hooks/queries', () => ({
  useUserQueries: jest.fn(),
  useServicesCounselingQueueQueries: jest.fn(),
  useServicesCounselingQueuePPMQueries: jest.fn(),
}));

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
          dodID: '555555555',
        },
        locator: 'AB5PC',
        requestedMoveDate: '2021-03-01T00:00:00.000Z',
        submittedAt: '2021-01-31T00:00:00.000Z',
        status: MOVE_STATUSES.NEEDS_SERVICE_COUNSELING,
        originDutyLocation: {
          name: 'Area 51',
        },
        originGBLOC: 'LKNQ',
      },
      {
        id: 'move2',
        customer: {
          agency: SERVICE_MEMBER_AGENCIES.COAST_GUARD,
          first_name: 'test another first',
          last_name: 'test another last',
          dodID: '4444444444',
        },
        locator: 'T12AR',
        requestedMoveDate: '2021-04-15T00:00:00.000Z',
        submittedAt: '2021-01-01T00:00:00.000Z',
        status: MOVE_STATUSES.NEEDS_SERVICE_COUNSELING,
        originDutyLocation: {
          name: 'Los Alamos',
        },
        originGBLOC: 'LKNQ',
      },
      {
        id: 'move3',
        customer: {
          agency: SERVICE_MEMBER_AGENCIES.MARINES,
          first_name: 'test third first',
          last_name: 'test third last',
          dodID: '4444444444',
        },
        locator: 'T12MP',
        requestedMoveDate: '2021-04-15T00:00:00.000Z',
        submittedAt: '2021-01-01T00:00:00.000Z',
        status: MOVE_STATUSES.NEEDS_SERVICE_COUNSELING,
        originDutyLocation: {
          name: 'Denver, 80136',
        },
        originGBLOC: 'LKNQ',
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
          dodID: '555555555',
        },
        locator: 'AB5PC',
        requestedMoveDate: '2021-03-01T00:00:00.000Z',
        submittedAt: '2021-01-31T00:00:00.000Z',
        status: MOVE_STATUSES.SERVICE_COUNSELING_COMPLETED,
        originDutyLocation: {
          name: 'Area 51',
        },
        originGBLOC: 'LKNQ',
      },
      {
        id: 'move2',
        customer: {
          agency: SERVICE_MEMBER_AGENCIES.COAST_GUARD,
          first_name: 'test another first',
          last_name: 'test another last',
          dodID: '4444444444',
        },
        locator: 'T12AR',
        requestedMoveDate: '2021-04-15T00:00:00.000Z',
        submittedAt: '2021-01-01T00:00:00.000Z',
        status: MOVE_STATUSES.SERVICE_COUNSELING_COMPLETED,
        originDutyLocation: {
          name: 'Los Alamos',
        },
        originGBLOC: 'LKNQ',
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
        <ServicesCounselingQueue />
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
    const wrapper = mount(
      <MockRouterProvider path={pagePath} params={{ queueType: 'counseling' }}>
        <ServicesCounselingQueue />
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
      expect(firstMove.find('td.lastName').text()).toBe('test last, test first');
      expect(firstMove.find('td.dodID').text()).toBe('555555555');
      expect(firstMove.find('td.locator').text()).toBe('AB5PC');
      expect(firstMove.find('td.status').text()).toBe('Needs counseling');
      expect(firstMove.find('td.requestedMoveDate').text()).toBe('01 Mar 2021');
      expect(firstMove.find('td.submittedAt').text()).toBe('31 Jan 2021');
      expect(firstMove.find('td.branch').text()).toBe('Army');
      expect(firstMove.find('td.originGBLOC').text()).toBe('LKNQ');
      expect(firstMove.find('td.originDutyLocation').text()).toBe('Area 51');

      const secondMove = moves.at(1);
      expect(secondMove.find('td.lastName').text()).toBe('test another last, test another first');
      expect(secondMove.find('td.dodID').text()).toBe('4444444444');
      expect(secondMove.find('td.locator').text()).toBe('T12AR');
      expect(secondMove.find('td.status').text()).toBe('Needs counseling');
      expect(secondMove.find('td.requestedMoveDate').text()).toBe('15 Apr 2021');
      expect(secondMove.find('td.submittedAt').text()).toBe('01 Jan 2021');
      expect(secondMove.find('td.branch').text()).toBe('Coast Guard');
      expect(secondMove.find('td.originGBLOC').text()).toBe('LKNQ');
      expect(secondMove.find('td.originDutyLocation').text()).toBe('Los Alamos');

      const thirdMove = moves.at(2);
      expect(thirdMove.find('td.lastName').text()).toBe('test third last, test third first');
      expect(thirdMove.find('td.dodID').text()).toBe('4444444444');
      expect(thirdMove.find('td.locator').text()).toBe('T12MP');
      expect(thirdMove.find('td.status').text()).toBe('Needs counseling');
      expect(thirdMove.find('td.requestedMoveDate').text()).toBe('15 Apr 2021');
      expect(thirdMove.find('td.submittedAt').text()).toBe('01 Jan 2021');
      expect(thirdMove.find('td.branch').text()).toBe('Marine Corps');
      expect(thirdMove.find('td.originGBLOC').text()).toBe('LKNQ');
      expect(thirdMove.find('td.originDutyLocation').text()).toBe('Denver, 80136');
    });

    it('sorts by submitted at date ascending by default', () => {
      expect(wrapper.find('th[data-testid="submittedAt"][role="columnheader"]').hasClass('sortAscending')).toBe(true);
    });

    it('allows sorting on certain columns', () => {
      expect(wrapper.find('th[data-testid="lastName"][role="columnheader"]').prop('onClick')).not.toBe(undefined);
      expect(wrapper.find('th[data-testid="dodID"][role="columnheader"]').prop('onClick')).not.toBe(undefined);
      expect(wrapper.find('th[data-testid="locator"][role="columnheader"]').prop('onClick')).not.toBe(undefined);
      expect(wrapper.find('th[data-testid="status"][role="columnheader"]').prop('onClick')).not.toBe(undefined);
      expect(wrapper.find('th[data-testid="requestedMoveDate"][role="columnheader"]').prop('onClick')).not.toBe(
        undefined,
      );
      expect(wrapper.find('th[data-testid="submittedAt"][role="columnheader"]').prop('onClick')).not.toBe(undefined);
      expect(wrapper.find('th[data-testid="branch"][role="columnheader"]').prop('onClick')).not.toBe(undefined);
      expect(wrapper.find('th[data-testid="originDutyLocation"][role="columnheader"]').prop('onClick')).not.toBe(
        undefined,
      );
    });

    it('disables sort by for origin GBLOC column', () => {
      expect(wrapper.find('th[data-testid="originGBLOC"][role="columnheader"]').prop('onClick')).toBe(undefined);
    });

    it('omits filter input element for origin GBLOC column', () => {
      expect(wrapper.find('th[data-testid="originGBLOC"] input').exists()).toBe(false);
    });
  });

  describe('service counseling completed moves', () => {
    useUserQueries.mockReturnValue(serviceCounselorUser);
    useServicesCounselingQueueQueries.mockReturnValue(serviceCounselingCompletedMoves);
    const wrapper = mount(
      <MockRouterProvider path={pagePath} params={{ queueType: 'counseling' }}>
        <ServicesCounselingQueue />
      </MockRouterProvider>,
    );

    it('displays move header with needs service counseling count', () => {
      expect(wrapper.find('h1').text()).toBe('Moves (2)');
    });

    it('should render the table', () => {
      expect(wrapper.find('Table').exists()).toBe(true);
    });

    it('formats the move data in rows', () => {
      const moves = wrapper.find('tbody tr');
      const firstMove = moves.at(0);
      expect(firstMove.find('td.status').text()).toBe('Service counseling completed');

      const secondMove = moves.at(1);
      expect(secondMove.find('td.status').text()).toBe('Service counseling completed');
    });
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
    ])('a %s user accessing path "%s"', (userDescription, queueType, showsCounselingTab, user) => {
      useUserQueries.mockReturnValue(user);
      useServicesCounselingQueueQueries.mockReturnValue(serviceCounselingCompletedMoves);
      useServicesCounselingQueuePPMQueries.mockReturnValue(emptyServiceCounselingMoves);
      render(
        <MockProviders path={pagePath} params={{ queueType }}>
          <ServicesCounselingQueue />
        </MockProviders>,
      );

      if (showsCounselingTab === 'counseling') {
        // Make sure "Counseling" is the active tab.
        const counselingActive = screen.getByText('Counseling Queue', { selector: '.usa-current .tab-title' });
        expect(counselingActive).toBeInTheDocument();

        // Check for the "Counseling" columns.
        expect(screen.getByText(/Status/)).toBeInTheDocument();
        expect(screen.getByText(/Requested move date/)).toBeInTheDocument();
        expect(screen.getByText(/Date submitted/)).toBeInTheDocument();
        expect(screen.getByText(/Origin GBLOC/)).toBeInTheDocument();
      } else if (showsCounselingTab === 'closeout') {
        // Make sure "PPM Closeout" is the active tab.
        const ppmCloseoutActive = screen.getByText('PPM Closeout Queue', { selector: '.usa-current .tab-title' });
        expect(ppmCloseoutActive).toBeInTheDocument();

        // Check for the "PPM Closeout" columns.
        expect(screen.getByText(/Closeout initiated/)).toBeInTheDocument();
        expect(screen.getByText(/PPM closeout location/)).toBeInTheDocument();
        expect(screen.getByText(/Full or partial PPM/)).toBeInTheDocument();
        expect(screen.getByText(/Destination duty location/)).toBeInTheDocument();
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

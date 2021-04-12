import React from 'react';
import { mount } from 'enzyme';

import ServicesCounselingQueue from './ServicesCounselingQueue';

import { useUserQueries, useServicesCounselingQueueQueries } from 'hooks/queries';
import { MockProviders } from 'testUtils';
import { MOVE_STATUSES } from 'shared/constants';
import SERVICE_MEMBER_AGENCIES from 'content/serviceMemberAgencies';

jest.mock('hooks/queries', () => ({
  useUserQueries: jest.fn(),
  useServicesCounselingQueueQueries: jest.fn(),
}));

const serviceCounselorUser = {
  isLoading: false,
  isError: false,
  data: {
    office_user: { transportation_office: { gbloc: 'LKNQ' } },
  },
};

const marineCorpsUser = {
  isLoading: false,
  isError: false,
  data: {
    office_user: { transportation_office: { gbloc: 'USMC' } },
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
        status: MOVE_STATUSES.NEEDS_SERVICE_COUNSELING,
        destinationDutyStation: {
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
        destinationDutyStation: {
          name: 'Los Alamos',
        },
        originGBLOC: 'LKNQ',
      },
    ],
  },
};

const marineCorpsNeedsCounselingMoves = {
  isLoading: false,
  isError: false,
  queueResult: {
    totalCount: 2,
    data: [
      {
        id: 'move1',
        customer: {
          agency: SERVICE_MEMBER_AGENCIES.MARINES,
          first_name: 'test first',
          last_name: 'test last',
          dodID: '555555555',
        },
        locator: 'AB5PC',
        requestedMoveDate: '2021-03-01T00:00:00.000Z',
        submittedAt: '2021-01-31T00:00:00.000Z',
        status: MOVE_STATUSES.NEEDS_SERVICE_COUNSELING,
        destinationDutyStation: {
          name: 'Area 51',
        },
        originGBLOC: 'LKNQ',
      },
      {
        id: 'move2',
        customer: {
          agency: SERVICE_MEMBER_AGENCIES.MARINES,
          first_name: 'test another first',
          last_name: 'test another last',
          dodID: '4444444444',
        },
        locator: 'T12AR',
        requestedMoveDate: '2021-04-15T00:00:00.000Z',
        submittedAt: '2021-01-01T00:00:00.000Z',
        status: MOVE_STATUSES.NEEDS_SERVICE_COUNSELING,
        destinationDutyStation: {
          name: 'Los Alamos',
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
        destinationDutyStation: {
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
        destinationDutyStation: {
          name: 'Los Alamos',
        },
        originGBLOC: 'LKNQ',
      },
    ],
  },
};

describe('ServicesCounselingQueue', () => {
  describe('no moves in service counseling statuses', () => {
    useUserQueries.mockImplementation(() => serviceCounselorUser);
    useServicesCounselingQueueQueries.mockImplementation(() => emptyServiceCounselingMoves);
    const wrapper = mount(
      <MockProviders initialEntries={['counseling/queue']}>
        <ServicesCounselingQueue />
      </MockProviders>,
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

  describe('non-USMC Service Counselor', () => {
    useUserQueries.mockImplementation(() => serviceCounselorUser);
    useServicesCounselingQueueQueries.mockImplementation(() => needsCounselingMoves);
    const wrapper = mount(
      <MockProviders initialEntries={['counseling/queue']}>
        <ServicesCounselingQueue />
      </MockProviders>,
    );

    it('displays move header with needs service counseling count', () => {
      expect(wrapper.find('h1').text()).toBe('Moves (2)');
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
      expect(firstMove.find('td.destinationDutyStation').text()).toBe('Area 51');

      const secondMove = moves.at(1);
      expect(secondMove.find('td.lastName').text()).toBe('test another last, test another first');
      expect(secondMove.find('td.dodID').text()).toBe('4444444444');
      expect(secondMove.find('td.locator').text()).toBe('T12AR');
      expect(secondMove.find('td.status').text()).toBe('Needs counseling');
      expect(secondMove.find('td.requestedMoveDate').text()).toBe('15 Apr 2021');
      expect(secondMove.find('td.submittedAt').text()).toBe('01 Jan 2021');
      expect(secondMove.find('td.branch').text()).toBe('Coast Guard');
      expect(secondMove.find('td.originGBLOC').text()).toBe('LKNQ');
      expect(secondMove.find('td.destinationDutyStation').text()).toBe('Los Alamos');
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
      expect(wrapper.find('th[data-testid="destinationDutyStation"][role="columnheader"]').prop('onClick')).not.toBe(
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

  describe('USMC Service Counselor', () => {
    useUserQueries.mockImplementation(() => marineCorpsUser);
    useServicesCounselingQueueQueries.mockImplementation(() => marineCorpsNeedsCounselingMoves);
    const wrapper = mount(
      <MockProviders initialEntries={['counseling/queue']}>
        <ServicesCounselingQueue />
      </MockProviders>,
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
      expect(firstMove.find('td.lastName').text()).toBe('test last, test first');
      expect(firstMove.find('td.dodID').text()).toBe('555555555');
      expect(firstMove.find('td.locator').text()).toBe('AB5PC');
      expect(firstMove.find('td.status').text()).toBe('Needs counseling');
      expect(firstMove.find('td.requestedMoveDate').text()).toBe('01 Mar 2021');
      expect(firstMove.find('td.submittedAt').text()).toBe('31 Jan 2021');
      expect(firstMove.find('td.branch').text()).toBe('Marine Corps');
      expect(firstMove.find('td.originGBLOC').text()).toBe('LKNQ');
      expect(firstMove.find('td.destinationDutyStation').text()).toBe('Area 51');

      const secondMove = moves.at(1);
      expect(secondMove.find('td.lastName').text()).toBe('test another last, test another first');
      expect(secondMove.find('td.dodID').text()).toBe('4444444444');
      expect(secondMove.find('td.locator').text()).toBe('T12AR');
      expect(secondMove.find('td.status').text()).toBe('Needs counseling');
      expect(secondMove.find('td.requestedMoveDate').text()).toBe('15 Apr 2021');
      expect(secondMove.find('td.submittedAt').text()).toBe('01 Jan 2021');
      expect(secondMove.find('td.branch').text()).toBe('Marine Corps');
      expect(secondMove.find('td.originGBLOC').text()).toBe('LKNQ');
      expect(secondMove.find('td.destinationDutyStation').text()).toBe('Los Alamos');
    });

    it('allows sorting on certain columns', () => {
      expect(wrapper.find('th[data-testid="lastName"][role="columnheader"]').prop('onClick')).toBeDefined();
      expect(wrapper.find('th[data-testid="dodID"][role="columnheader"]').prop('onClick')).toBeDefined();
      expect(wrapper.find('th[data-testid="locator"][role="columnheader"]').prop('onClick')).toBeDefined();
      expect(wrapper.find('th[data-testid="status"][role="columnheader"]').prop('onClick')).toBeDefined();
      expect(wrapper.find('th[data-testid="requestedMoveDate"][role="columnheader"]').prop('onClick')).toBeDefined();
      expect(wrapper.find('th[data-testid="submittedAt"][role="columnheader"]').prop('onClick')).toBeDefined();
      expect(wrapper.find('th[data-testid="originGBLOC"][role="columnheader"]').prop('onClick')).toBeDefined();
      expect(
        wrapper.find('th[data-testid="destinationDutyStation"][role="columnheader"]').prop('onClick'),
      ).toBeDefined();
    });

    it('disables sorting on the branch column', () => {
      expect(wrapper.find('th[data-testid="branch"][role="columnheader"]').prop('onClick')).toBe(undefined);
    });

    it('omits select input element for branch column', () => {
      expect(wrapper.find('th[data-testid="branch"] select').exists()).toBe(false);
    });

    it('includes filter input on Origin GBLOC column', () => {
      expect(wrapper.find('th[data-testid="originGBLOC"] input').exists()).toBe(true);
    });
  });

  describe('service counseling completed moves', () => {
    useUserQueries.mockImplementation(() => serviceCounselorUser);
    useServicesCounselingQueueQueries.mockImplementation(() => serviceCounselingCompletedMoves);
    const wrapper = mount(
      <MockProviders initialEntries={['counseling/queue']}>
        <ServicesCounselingQueue />
      </MockProviders>,
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
});

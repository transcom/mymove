import React from 'react';
import { mount } from 'enzyme';

import ServicesCounselingMoveInfo from './ServicesCounselingMoveInfo';

import { MockProviders } from 'testUtils';

const testMoveCode = '1A5PM3';

jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useParams: jest.fn().mockReturnValue({ moveCode: '1A5PM3' }),
}));

jest.mock('hooks/queries', () => ({
  ...jest.requireActual('hooks/queries'),
  useTXOMoveInfoQueries: () => {
    return {
      customerData: { id: '2468', last_name: 'Kerry', first_name: 'Smith', dodID: '999999999' },
      order: {
        id: '4321',
        customerID: '2468',
        uploaded_order_id: '2',
        departmentIndicator: 'Navy',
        grade: 'E-6',
        originDutyStation: {
          name: 'JBSA Lackland',
        },
        destinationDutyStation: {
          name: 'JB Lewis-McChord',
        },
        report_by_date: '2018-08-01',
      },
      isLoading: false,
      isError: false,
      isSuccess: true,
    };
  },
}));

describe('Services Counseling Move Info Container', () => {
  it('should render the move tab container', () => {
    const wrapper = mount(
      <MockProviders initialEntries={[`/counseling/moves/${testMoveCode}/details`]}>
        <ServicesCounselingMoveInfo />
      </MockProviders>,
    );

    expect(wrapper.find('CustomerHeader').exists()).toBe(true);
  });

  describe('routing', () => {
    it('should handle the Services Counseling Move Details route', () => {
      const wrapper = mount(
        <MockProviders initialEntries={[`/counseling/moves/${testMoveCode}/details`]}>
          <ServicesCounselingMoveInfo />
        </MockProviders>,
      );

      expect(wrapper.find('ServicesCounselingMoveDetails')).toHaveLength(1);
    });

    it('should redirect from move info root to the Services Counseling Move Details route', () => {
      const wrapper = mount(
        <MockProviders initialEntries={[`/counseling/moves/${testMoveCode}`]}>
          <ServicesCounselingMoveInfo />
        </MockProviders>,
      );

      const renderedRoute = wrapper.find('ServicesCounselingMoveDetails');
      expect(renderedRoute).toHaveLength(1);
    });
  });
});

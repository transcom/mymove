import 'raf/polyfill';
import React from 'react';

import moment from 'moment';
import { shallow, mount } from 'enzyme';

import ConnectedLanding, { Landing } from './index';
import { MoveSummary } from './MoveSummary';
import PpmAlert from './PpmAlert';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { MockProviders } from 'testUtils';

describe('HomePage tests', () => {
  let wrapper;

  const minProps = {
    showLoggedInUser: () => {},
    context: {
      flags: {
        hhgFlow: false,
      },
    },
  };

  describe('when not loggedIn', () => {
    it('renders without crashing', () => {
      const div = document.createElement('div');
      wrapper = shallow(<Landing isLoggedIn={false} {...minProps} />, div);
      expect(wrapper.find('.grid-container').length).toEqual(1);
    });
  });

  describe('When loggedIn', () => {
    let service_member = { id: 'foo' };
    it('renders without crashing', () => {
      const div = document.createElement('div');
      wrapper = shallow(<Landing isLoggedIn={true} {...minProps} />, div);
      expect(wrapper.find('.grid-container').length).toEqual(1);
    });

    describe('if the service member is not created', () => {
      it('creates the service member', () => {
        const testProps = {
          createServiceMember: jest.fn(() => Promise.resolve()),
          showLoggedInUser: jest.fn(),
          push: jest.fn(),
        };

        const wrapper = mount(
          <Landing isLoggedIn={true} loggedInUserSuccess={false} serviceMember={{}} {...testProps} />,
        );

        wrapper.setProps({ loggedInUserSuccess: true });
        expect(testProps.createServiceMember).toHaveBeenCalledTimes(1);
        expect(testProps.showLoggedInUser).toHaveBeenCalledTimes(1);
      });
    });

    describe('When the user has never logged in before', () => {
      it('redirects to enter profile page', () => {
        const mockPush = jest.fn();
        const wrapper = mount(
          <Landing
            isLoggedIn={true}
            serviceMember={{}}
            createdServiceMemberIsLoading={false}
            loggedInUserSuccess={true}
            isProfileComplete={false}
            push={mockPush}
            {...minProps}
          />,
        );

        const resumeMoveFn = jest.spyOn(wrapper.instance(), 'resumeMove');
        wrapper.setProps({ serviceMember: service_member });
        expect(resumeMoveFn).toHaveBeenCalledTimes(1);
      });
    });

    describe('When the user profile has started but is not complete', () => {
      it('MoveSummary does not render', () => {
        const div = document.createElement('div');
        wrapper = shallow(
          <Landing
            serviceMember={service_member}
            isLoggedIn={true}
            loggedInUserSuccess={true}
            isProfileComplete={false}
            {...minProps}
          />,
          div,
        );
        expect(wrapper.find('.grid-container').length).toEqual(1);
        expect(wrapper.find(MoveSummary).length).toEqual(0);
      });
    });

    describe('When orders have been entered but the move is not complete', () => {
      const futureFortNight = moment().add(14, 'day');
      const orders = {
        orders_type: 'foo',
        issue_date: '2019-01-01',
        report_by_date: '2019-02-01',
        new_duty_station: { id: 'something' },
      };
      const ppmObj = {
        original_move_date: futureFortNight,
        weight_estimate: '10000',
        estimated_incentive: '$24665.59 - 27261.97',
      };

      describe('When a ppm only move is submitted', () => {
        it('renders the ppm only alert', () => {
          const moveObj = { selected_move_type: SHIPMENT_OPTIONS.PPM, status: 'SUBMITTED', moveSubmitSuccess: true };
          wrapper = shallow(
            <Landing
              move={moveObj}
              moveSubmitSuccess={moveObj.moveSubmitSuccess}
              serviceMember={service_member}
              ppm={ppmObj}
              orders={orders}
              isLoggedIn={true}
              loggedInUserSuccess={true}
              isProfileComplete={true}
              {...minProps}
            />,
          );
          const ppmAlert = wrapper.find(PpmAlert).shallow();

          expect(ppmAlert.length).toEqual(1);
          expect(ppmAlert.props().heading).toEqual('Congrats - your move is submitted!');
        });
      });
    });
  });
});

describe('ConnectedLanding', () => {
  const initialState = {
    entities: {
      user: {},
      orders: {},
      mtoShipments: {},
      backupContacts: {},
      personallyProcuredMoves: {},
    },
  };

  it('renders while loading', () => {
    const wrapper = mount(
      <MockProviders initialState={initialState}>
        <ConnectedLanding />
      </MockProviders>,
    );

    console.log(wrapper.debug());
    expect(wrapper.exists()).toBe(true);
  });
});

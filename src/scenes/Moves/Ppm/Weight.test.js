import React from 'react';
import { shallow, mount } from 'enzyme';
import moment from 'moment';

import ConnectedPpmWeight, { PpmWeight } from './Weight';
import { MockProviders } from 'testUtils';

describe('Weight', () => {
  const moveDate = moment().add(7, 'day');
  const minProps = {
    selectedWeightInfo: { min: 0, max: 0 },
    hasLoadSuccess: true,
    entitlement: {
      weight: 0,
      pro_gear: 0,
      pro_gear_spouse: 0,
    },
    currentPPM: {
      original_move_date: moveDate,
      pickup_postal_code: '00000',
    },
    tempCurrentPPM: {
      original_move_date: moveDate,
      pickup_postal_code: '00000',
    },
    orders: { id: 1 },
    getPpmWeightEstimate: () =>
      new Promise((approve, reject) => {
        approve({});
      }),
    match: { params: { moveId: 'some id' } },
    loadPPMs: jest.fn(),
    fetchLatestOrders: jest.fn(),
  };
  it('Component renders', () => {
    expect(shallow(<PpmWeight {...minProps} />).length).toEqual(1);
  });

  describe('Test estimate icon', () => {
    let wrapper;
    describe('Move under 500 lbs', () => {
      it('Should show car icon for 499 lbs', () => {
        wrapper = shallow(<PpmWeight {...minProps} currentPPM={{ weight_estimate: 499 }} />);
        expect(wrapper.find({ 'data-testid': 'vehicleIcon' }).prop('src')).toEqual('car-gray.svg');
      });
    });
    describe('Move between 500 lbs and 1499 lbs', () => {
      it('Should show trailer icon for 500 lbs', () => {
        wrapper = shallow(<PpmWeight {...minProps} currentPPM={{ weight_estimate: 500 }} />);
        expect(wrapper.find({ 'data-testid': 'vehicleIcon' }).prop('src')).toEqual('trailer-gray.svg');
      });
      it('Should show trailer icon for 1499 lbs', () => {
        wrapper = shallow(<PpmWeight {...minProps} currentPPM={{ weight_estimate: 1499 }} />);
        expect(wrapper.find({ 'data-testid': 'vehicleIcon' }).prop('src')).toEqual('trailer-gray.svg');
      });
    });
    describe('Move 1500 lbs or greater', () => {
      it('Should show truck icon for 1500 lbs', () => {
        wrapper = shallow(<PpmWeight {...minProps} currentPPM={{ weight_estimate: 1500 }} />);
        expect(wrapper.find({ 'data-testid': 'vehicleIcon' }).prop('src')).toEqual('truck-gray.svg');
      });
    });
  });

  describe('Test estimate text', () => {
    let wrapper;

    describe('Move under 500 lbs', () => {
      it('Should show text for 499 lbs', () => {
        wrapper = shallow(<PpmWeight {...minProps} currentPPM={{ weight_estimate: 499 }} />);
        expect(wrapper.find({ 'data-testid': 'estimateText' }).text()).toEqual('Just a few things. One trip in a car.');
      });
    });
    describe('Move between 500 lbs and 1000 lbs', () => {
      it('Should show text for 500 lbs', () => {
        wrapper = shallow(<PpmWeight {...minProps} currentPPM={{ weight_estimate: 500 }} />);
        expect(wrapper.find({ 'data-testid': 'estimateText' }).text()).toEqual(
          'Studio apartment, minimal stuff. A large car, a pickup, a van, or a car with trailer.',
        );
      });
      it('Should show text for 999 lbs', () => {
        wrapper = shallow(<PpmWeight {...minProps} currentPPM={{ weight_estimate: 999 }} />);
        expect(wrapper.find({ 'data-testid': 'estimateText' }).text()).toEqual(
          'Studio apartment, minimal stuff. A large car, a pickup, a van, or a car with trailer.',
        );
      });
    });
    describe('Move between 1000 lbs and 2000 lbs', () => {
      it('Should show text for 1000 lbs', () => {
        wrapper = shallow(<PpmWeight {...minProps} currentPPM={{ weight_estimate: 1000 }} />);
        expect(wrapper.find({ 'data-testid': 'estimateText' }).text()).toEqual(
          '1-2 rooms, light furniture. A pickup, a van, or a car with a small or medium trailer.',
        );
      });
      it('Should show text for 1999 lbs', () => {
        wrapper = shallow(<PpmWeight {...minProps} currentPPM={{ weight_estimate: 1999 }} />);
        expect(wrapper.find({ 'data-testid': 'estimateText' }).text()).toEqual(
          '1-2 rooms, light furniture. A pickup, a van, or a car with a small or medium trailer.',
        );
      });
    });
    describe('Move between 2000 lbs and 3000 lbs', () => {
      it('Should show text for 2000 lbs', () => {
        wrapper = shallow(<PpmWeight {...minProps} currentPPM={{ weight_estimate: 2000 }} />);
        expect(wrapper.find({ 'data-testid': 'estimateText' }).text()).toEqual(
          '2-3 rooms, some bulky items. Cargo van, small or medium moving truck, medium or large cargo trailer.',
        );
      });
      it('Should show text for 2999 lbs', () => {
        wrapper = shallow(<PpmWeight {...minProps} currentPPM={{ weight_estimate: 2999 }} />);
        expect(wrapper.find({ 'data-testid': 'estimateText' }).text()).toEqual(
          '2-3 rooms, some bulky items. Cargo van, small or medium moving truck, medium or large cargo trailer.',
        );
      });
    });
    describe('Move between 3000 lbs and 4000 lbs', () => {
      it('Should show text for 3000 lbs', () => {
        wrapper = shallow(<PpmWeight {...minProps} currentPPM={{ weight_estimate: 3000 }} />);
        expect(wrapper.find({ 'data-testid': 'estimateText' }).text()).toEqual(
          '3-4 rooms. Small to medium moving truck, or a couple of trips.',
        );
      });
      it('Should show text for 3999 lbs', () => {
        wrapper = shallow(<PpmWeight {...minProps} currentPPM={{ weight_estimate: 3999 }} />);
        expect(wrapper.find({ 'data-testid': 'estimateText' }).text()).toEqual(
          '3-4 rooms. Small to medium moving truck, or a couple of trips.',
        );
      });
    });
    describe('Move between 4000 lbs and 5000 lbs', () => {
      it('Should show text for 4000 lbs', () => {
        wrapper = shallow(<PpmWeight {...minProps} currentPPM={{ weight_estimate: 4000 }} />);
        expect(wrapper.find({ 'data-testid': 'estimateText' }).text()).toEqual(
          '4+ rooms, or just a lot of large, heavy things. Medium or large moving truck, or multiple trips.',
        );
      });
      it('Should show text for 4999 lbs', () => {
        wrapper = shallow(<PpmWeight {...minProps} currentPPM={{ weight_estimate: 4999 }} />);
        expect(wrapper.find({ 'data-testid': 'estimateText' }).text()).toEqual(
          '4+ rooms, or just a lot of large, heavy things. Medium or large moving truck, or multiple trips.',
        );
      });
    });
    describe('Move between 5000 lbs and 6000 lbs', () => {
      it('Should show text for 5000 lbs', () => {
        wrapper = shallow(<PpmWeight {...minProps} currentPPM={{ weight_estimate: 5000 }} />);
        expect(wrapper.find({ 'data-testid': 'estimateText' }).text()).toEqual(
          'Many rooms, many things, lots of them heavy. Medium or large moving truck, or multiple trips.',
        );
      });
      it('Should show text for 5999 lbs', () => {
        wrapper = shallow(<PpmWeight {...minProps} currentPPM={{ weight_estimate: 5999 }} />);
        expect(wrapper.find({ 'data-testid': 'estimateText' }).text()).toEqual(
          'Many rooms, many things, lots of them heavy. Medium or large moving truck, or multiple trips.',
        );
      });
    });
    describe('Move between 6000 lbs and 7000 lbs', () => {
      it('Should show text for 6000 lbs', () => {
        wrapper = shallow(<PpmWeight {...minProps} currentPPM={{ weight_estimate: 6000 }} />);
        expect(wrapper.find({ 'data-testid': 'estimateText' }).text()).toEqual(
          'Large house, a lot of things. The biggest rentable moving trucks, or multiple trips or vehicles.',
        );
      });
      it('Should show text for 6999 lbs', () => {
        wrapper = shallow(<PpmWeight {...minProps} currentPPM={{ weight_estimate: 6999 }} />);
        expect(wrapper.find({ 'data-testid': 'estimateText' }).text()).toEqual(
          'Large house, a lot of things. The biggest rentable moving trucks, or multiple trips or vehicles.',
        );
      });
    });
    describe('Move 7000 lbs or over', () => {
      it('Should show text for 7000 lbs', () => {
        wrapper = shallow(<PpmWeight {...minProps} currentPPM={{ weight_estimate: 7000 }} />);
        expect(wrapper.find({ 'data-testid': 'estimateText' }).text()).toEqual(
          'A large house or small palace, many heavy or bulky items. Multiple trips using large vehicles, or hire professional movers.',
        );
      });
    });
  });

  describe('Incentive estimate errors', () => {
    let wrapper;
    const iconAndTextProps = {
      currentPPM: {},
      orders: { id: 1 },
      getPpmWeightEstimate: () =>
        new Promise((approve, reject) => {
          reject({});
        }),
    };
    it('Should not show an estimate error', () => {
      wrapper = shallow(<PpmWeight {...minProps} {...iconAndTextProps} />);
      expect(wrapper.find('.error-message').exists()).toBe(false);
      expect(wrapper.find('ReduxForm').props().readyToSubmit).toEqual(true);
    });
    describe('Short Haul Error', () => {
      it('Should show short haul error and next button disabled', () => {
        wrapper = shallow(<PpmWeight {...minProps} {...iconAndTextProps} rateEngineError={{ statusCode: 409 }} />);
        expect(wrapper.find('.error-message').exists()).toBe(true);
        expect(wrapper.find('Alert').dive().text()).toMatch(
          /MilMove does not presently support short-haul PPM moves. Please contact your PPPO./,
        );
        expect(wrapper.find('ReduxForm').props().readyToSubmit).toEqual(false); // next button should be disabled
      });
    });
    describe('No rate data error', () => {
      it('Should show estimate error and next button not disabled', () => {
        wrapper = shallow(<PpmWeight {...minProps} {...iconAndTextProps} rateEngineError={{ statusCode: 404 }} />);
        expect(wrapper.find('.error-message').exists()).toBe(true);
        expect(wrapper.find('Alert').dive().text()).toMatch(
          /There was an issue retrieving an estimate for your incentive./,
        );
        expect(wrapper.find('ReduxForm').props().readyToSubmit).toEqual(true);
      });
    });
    it('Should show estimate not retrieved error', () => {
      wrapper = shallow(<PpmWeight {...minProps} {...iconAndTextProps} />);
      wrapper.setState({ hasEstimateError: true });
      expect(wrapper.find('.error-message').exists()).toBe(true);
      expect(wrapper.find('Alert').dive().text()).toMatch(
        /There was an issue retrieving an estimate for your incentive./,
      );
      expect(wrapper.find('ReduxForm').props().readyToSubmit).toEqual(true);
    });
  });
});

describe('ConnectedPpmWeight', () => {
  it('renders without errors', () => {
    const initialState = {
      entities: {
        user: {
          userId123: {
            id: 'userId123',
            service_member: 'testServiceMemberId456',
          },
        },
        serviceMembers: {
          testServiceMemberId456: {
            id: 'testServiceMemberId456',
            first_name: 'Frida',
            current_station: {
              address: {
                postal_code: '90813',
              },
            },
            weight_allotment: {
              total_weight_self: 2000,
              total_weight_self_plus_dependents: 2000,
              pro_gear_weight: 500,
              pro_gear_weight_spouse: 0,
            },
          },
        },
        orders: {
          testOrder123: {
            service_member_id: 'testServiceMemberId456',
            has_dependents: false,
            spouse_has_pro_gear: false,
          },
        },
        mtoShipments: {},
        backupContacts: {},
        personallyProcuredMoves: {
          testPPM345: {
            id: 'testPPM345',
            moveId: 'testMove890',
            status: 'DRAFT',
            original_move_date: moment().add(7, 'day'),
            pickup_postal_code: '80913',
          },
        },
      },
      moves: {
        currentMove: {
          id: 'testMove890',
        },
      },
      user: {
        userInfo: {
          service_member: {
            id: 'testServiceMemberId456',
          },
        },
      },
      serviceMember: {
        currentServiceMember: {
          id: 'testServiceMemberId456',
          first_name: 'Frida',
          current_station: {
            address: {
              postal_code: '90813',
            },
          },
          weight_allotment: {
            total_weight_self: 2500,
            total_weight_self_plus_dependents: 2000,
            pro_gear_weight: 500,
            pro_gear_weight_spouse: 0,
          },
        },
      },
    };

    const wrapper = mount(
      <MockProviders initialState={initialState}>
        <ConnectedPpmWeight
          pages={[]}
          pageKey="/moves/:moveId/ppm-incentive"
          match={{ params: { moveId: 'testMove890' } }}
        />
      </MockProviders>,
    );

    expect(wrapper.find('[data-testid="weight-entitlement"]').text()).toContain('2,500');
  });
});

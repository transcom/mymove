import React from 'react';
import { PpmWeight } from './Weight';
import { shallow } from 'enzyme';

describe('Weight', () => {
  const minProps = {
    selectedWeightInfo: { min: 0, max: 0 },
    hasLoadSuccess: true,
    entitlement: {
      weight: 0,
      pro_gear: 0,
      pro_gear_spouse: 0,
    },
  };
  it('Component renders', () => {
    expect(shallow(<PpmWeight {...minProps} />).length).toEqual(1);
  });

  describe('Test Estimate icon', () => {
    let wrapper;
    const iconAndTextProps = {
      currentPpm: {},
      orders: { id: 1 },
      getPpmWeightEstimate: jest.fn(),
    };
    describe('Move under 500 lbs', () => {
      it('Should show car icon', () => {
        wrapper = shallow(<PpmWeight {...minProps} {...iconAndTextProps} currentWeight={499} />);
        expect(wrapper.find({ 'data-cy': 'vehicleIcon' }).prop('src')).toEqual('car-gray.svg');
        expect(wrapper.find({ 'data-cy': 'estimateText' }).text()).toEqual('Just a few things. One trip in a car.');
      });
    });
    describe('Move between 500 lbs and 1499 lbs', () => {
      it('Should show correct icon and text for move size', () => {
        wrapper = shallow(<PpmWeight {...minProps} {...iconAndTextProps} currentWeight={500} />);
        expect(wrapper.find({ 'data-cy': 'vehicleIcon' }).prop('src')).toEqual('trailer-gray.svg');
      });
      it('Should show correct icon and text for move size', () => {
        wrapper = shallow(<PpmWeight {...minProps} {...iconAndTextProps} currentWeight={1499} />);
        expect(wrapper.find({ 'data-cy': 'vehicleIcon' }).prop('src')).toEqual('trailer-gray.svg');
      });
    });
    describe('Move 1500 lbs or greater', () => {
      it('Should show correct icon and text for move size', () => {
        wrapper = shallow(<PpmWeight {...minProps} {...iconAndTextProps} currentWeight={1500} />);
        expect(wrapper.find({ 'data-cy': 'vehicleIcon' }).prop('src')).toEqual('truck-gray.svg');
      });
    });
  });

  describe('Incentive estimate errors', () => {
    let wrapper;
    it('Should not show an estimate error', () => {
      wrapper = shallow(<PpmWeight {...minProps} />);
      expect(wrapper.find('.error-message').exists()).toBe(false);
      expect(wrapper.find('ReduxForm').props().readyToSubmit).toEqual(true);
    });
    describe('Short Haul Error', () => {
      it('Should show short haul error and next button disabled', () => {
        wrapper = shallow(<PpmWeight {...minProps} rateEngineError={{ statusCode: 409 }} />);
        expect(wrapper.find('.error-message').exists()).toBe(true);
        expect(
          wrapper
            .find('Alert')
            .dive()
            .text(),
        ).toMatch(/MilMove does not presently support short-haul PPM moves. Please contact your PPPO./);
        expect(wrapper.find('ReduxForm').props().readyToSubmit).toEqual(false); // next button should be disabled
      });
    });
    describe('No rate data error', () => {
      it('Should show estimate error and next button not disabled', () => {
        wrapper = shallow(<PpmWeight {...minProps} rateEngineError={{ statusCode: 404 }} />);
        expect(wrapper.find('.error-message').exists()).toBe(true);
        expect(
          wrapper
            .find('Alert')
            .dive()
            .text(),
        ).toMatch(/There was an issue retrieving an estimate for your incentive./);
        expect(wrapper.find('ReduxForm').props().readyToSubmit).toEqual(true);
      });
    });
    it('Should show estimate not retrieved error', () => {
      wrapper = shallow(<PpmWeight {...minProps} hasEstimateError={true} />);
      expect(wrapper.find('.error-message').exists()).toBe(true);
      expect(
        wrapper
          .find('Alert')
          .dive()
          .text(),
      ).toMatch(/There was an issue retrieving an estimate for your incentive./);
      expect(wrapper.find('ReduxForm').props().readyToSubmit).toEqual(true);
    });
  });
});

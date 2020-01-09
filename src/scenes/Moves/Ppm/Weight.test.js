import React from 'react';
import { PpmWeight } from './Weight';
import { shallow } from 'enzyme';

describe('Weight', () => {
  describe('Incentive estimate errors', () => {
    let wrapper;
    const props = {
      selectedWeightInfo: { min: 500, max: 1000 },
      hasLoadSuccess: true,
    };
    beforeEach(() => {
      props.rateEngineError = undefined;
      props.hasEstimateError = undefined;
    });
    it('Should not show an estimate error', () => {
      wrapper = shallow(<PpmWeight {...props} />);
      expect(wrapper.find('.error-message').exists()).toBe(false);
    });
    it('Should show short haul error', () => {
      props.rateEngineError = true;
      wrapper = shallow(<PpmWeight {...props} />);
      expect(wrapper.find('.error-message').exists()).toBe(true);
      expect(
        wrapper
          .find('Alert')
          .dive()
          .text(),
      ).toMatch(/MilMove does not presently support short-haul PPM moves. Please contact your PPPO./);
    });
    it('Should show estimate not retrieved error', () => {
      props.hasEstimateError = true;
      wrapper = shallow(<PpmWeight {...props} />);
      expect(wrapper.find('.error-message').exists()).toBe(true);
      expect(
        wrapper
          .find('Alert')
          .dive()
          .text(),
      ).toMatch(/There was an issue retrieving an estimate for your incentive./);
    });
  });
});

import React from 'react';
import { PpmWeight } from './Weight';
import { shallow } from 'enzyme';

describe('Weight', () => {
  const minProps = {
    selectedWeightInfo: { min: 0, max: 0 },
    hasLoadSuccess: true,
  };
  it('Renders', () => {
    expect(shallow(<PpmWeight {...minProps} />).length).toEqual(1);
  });
  describe('Incentive estimate errors', () => {
    let wrapper;
    it('Should not show an estimate error', () => {
      wrapper = shallow(<PpmWeight {...minProps} />);
      expect(wrapper.find('.error-message').exists()).toBe(false);
    });
    it('Should show short haul error', () => {
      wrapper = shallow(<PpmWeight {...minProps} rateEngineError={true} />);
      expect(wrapper.find('.error-message').exists()).toBe(true);
      expect(
        wrapper
          .find('Alert')
          .dive()
          .text(),
      ).toMatch(/MilMove does not presently support short-haul PPM moves. Please contact your PPPO./);
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
    });
  });
});

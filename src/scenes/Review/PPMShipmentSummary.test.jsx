import React from 'react';
import { PPMShipmentSummary } from './PPMShipmentSummary';
import { shallow } from 'enzyme';

describe('Review -> Ppm Shipment Summary', () => {
  const ppmEst = {
    hasEstimateError: false,
    hasEstimateSuccess: false,
    rateEngineError: null,
    originDutyLocationZip: '',
    incentive_estimate_min: 0,
    incentive_estimate_max: 0,
  };
  const minProps = {
    ppm: {},
    orders: {},
    movePath: '',
    ppmEstimate: {
      ...ppmEst,
    },
  };
  it('Component renders', () => {
    expect(shallow(<PPMShipmentSummary {...minProps} />).length).toEqual(1);
  });

  describe('Incentive estimate', () => {
    let wrapper;
    it('Should show estimate min & max', () => {
      wrapper = shallow(
        <PPMShipmentSummary
          {...minProps}
          ppmEstimate={{ ...ppmEst, incentive_estimate_min: 500, incentive_estimate_max: 1055 }}
        />,
      );
      expect(wrapper.find({ 'data-testid': 'estimate' }).exists()).toBe(true);
      expect(wrapper.find({ 'data-testid': 'estimate' }).text()).toEqual(' $5.00 - 10.55');
    });
    it('Should show short haul error', () => {
      wrapper = shallow(
        <PPMShipmentSummary {...minProps} ppmEstimate={{ ...ppmEst, rateEngineError: { statusCode: 409 } }} />,
      );
      expect(wrapper.find({ 'data-testid': 'estimateError' }).exists()).toBe(true);
      expect(wrapper.find({ 'data-testid': 'estimateError' }).text()).toMatch(
        /MilMove does not presently support short-haul PPM moves. Please contact your PPPO./,
      );
    });
    it('Should show estimate not ready error', () => {
      wrapper = shallow(
        <PPMShipmentSummary {...minProps} ppmEstimate={{ ...ppmEst, rateEngineError: { statusCode: 404 } }} />,
      );
      expect(wrapper.find({ 'data-testid': 'estimateError' }).exists()).toBe(true);
      expect(wrapper.find({ 'data-testid': 'estimateError' }).text()).toMatch(/Not ready yet/);
    });
    it('Should show estimate not ready error', () => {
      wrapper = shallow(<PPMShipmentSummary {...minProps} ppmEstimate={{ ...ppmEst, hasEstimateError: true }} />);
      expect(wrapper.find({ 'data-testid': 'estimateError' }).exists()).toBe(true);
      expect(wrapper.find({ 'data-testid': 'estimateError' }).text()).toMatch(/Not ready yet/);
    });
  });
});

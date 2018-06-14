import React from 'react';
import { shallow } from 'enzyme';
import { MoveSummary } from './MoveSummary';
import moment from 'moment';

describe('MoveSummary', () => {
  let wrapper, div;
  const editMoveFn = jest.fn();
  const resumeMoveFn = jest.fn();
  const entitlementObj = { sum: '10000' };
  const serviceMember = { current_station: { name: 'Ft Carson' } };
  const ordersObj = {};
  const getStepHtml = (
    entitlementObj,
    serviceMember,
    ordersObj,
    moveObj,
    ppmObj,
    editMoveFn,
    resumeMoveFn,
  ) => {
    const wrapper = shallow(
      <MoveSummary
        entitlement={entitlementObj}
        profile={serviceMember}
        orders={ordersObj}
        move={moveObj}
        ppm={ppmObj}
        editMove={editMoveFn}
        resumeMove={resumeMoveFn}
      />,
      div,
    );
    const stepText = wrapper
      .find('.step')
      .find('div.title')
      .first()
      .html();
    return stepText;
  };
  describe('when a move is in submitted state', () => {
    it('renders submitted content', () => {
      const moveObj = { status: 'SUBMITTED' };
      const futureFortNight = moment().add(14, 'day');
      const ppmObj = {
        planned_move_date: futureFortNight,
        weight_estimate: '10000',
        estimated_incentive: '$24665.59 - 27261.97',
      };
      expect(
        getStepHtml(
          entitlementObj,
          serviceMember,
          ordersObj,
          moveObj,
          ppmObj,
          editMoveFn,
          resumeMoveFn,
        ),
      ).toEqual('<div class="title">Next Step: Awaiting approval</div>');
    });
  });
  describe('when a move is in approved state', () => {
    it('renders submitted content', () => {
      const moveObj = { status: 'APPROVED' };
      const futureFortNight = moment().add(14, 'day');
      const ppmObj = {
        planned_move_date: futureFortNight,
        weight_estimate: '10000',
        estimated_incentive: '$24665.59 - 27261.97',
      };
      expect(
        getStepHtml(
          entitlementObj,
          serviceMember,
          ordersObj,
          moveObj,
          ppmObj,
          editMoveFn,
          resumeMoveFn,
        ),
      ).toEqual('<div class="title">Next Step: Get ready to move</div>');
    });
  });
  describe('when a move is in in progress state', () => {
    it('renders submitted content', () => {
      const moveObj = { status: 'APPROVED' };
      const pastFortNight = moment().subtract(14, 'day');
      const ppmObj = {
        planned_move_date: pastFortNight,
        weight_estimate: '10000',
        estimated_incentive: '$24665.59 - 27261.97',
      };
      expect(
        getStepHtml(
          entitlementObj,
          serviceMember,
          ordersObj,
          moveObj,
          ppmObj,
          editMoveFn,
          resumeMoveFn,
        ),
      ).toEqual('<div class="title">Next Step: Request payment</div>');
    });
  });
});

import React from 'react';
import { shallow } from 'enzyme';
import {
  MoveSummary,
  CanceledMoveSummary,
  ApprovedMoveSummary,
  DraftMoveSummary,
  SubmittedMoveSummary,
} from './MoveSummary';
import moment from 'moment';

describe('MoveSummary', () => {
  let wrapper, div;
  const editMoveFn = jest.fn();
  const resumeMoveFn = jest.fn();
  const entitlementObj = { sum: '10000' };
  const serviceMember = { current_station: { name: 'Ft Carson' } };
  const ordersObj = {};
  const getShallowRender = (
    entitlementObj,
    serviceMember,
    ordersObj,
    moveObj,
    ppmObj,
    editMoveFn,
    resumeMoveFn,
  ) => {
    return shallow(
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
  };
  describe('when a move is in canceled state', () => {
    it('renders submitted content', () => {
      const moveObj = { status: 'CANCELED' };
      const futureFortNight = moment().add(14, 'day');
      const ppmObj = {
        planned_move_date: futureFortNight,
        weight_estimate: '10000',
        estimated_incentive: '$24665.59 - 27261.97',
      };
      const subComponent = getShallowRender(
        entitlementObj,
        serviceMember,
        ordersObj,
        moveObj,
        ppmObj,
        editMoveFn,
        resumeMoveFn,
      ).find(CanceledMoveSummary);
      expect(subComponent).not.toBeNull();
      expect(
        subComponent
          .dive()
          .find('h2')
          .html(),
      ).toEqual('<h2>New move</h2>');
    });
  });
  describe('when a move is in submitted state', () => {
    it('renders submitted content', () => {
      const moveObj = { status: 'SUBMITTED' };
      const futureFortNight = moment().add(14, 'day');
      const ppmObj = {
        planned_move_date: futureFortNight,
        weight_estimate: '10000',
        estimated_incentive: '$24665.59 - 27261.97',
      };
      const subComponent = getShallowRender(
        entitlementObj,
        serviceMember,
        ordersObj,
        moveObj,
        ppmObj,
        editMoveFn,
        resumeMoveFn,
      ).find(SubmittedMoveSummary);
      expect(subComponent).not.toBeNull();
      expect(
        subComponent
          .dive()
          .find('.step')
          .find('div.title')
          .first()
          .html(),
      ).toEqual('<div class="title">Next Step: Awaiting approval</div>');
    });
  });
  describe('when a move is in approved state but ppm is submitted state', () => {
    it('renders submitted rather than approved content', () => {
      const moveObj = { status: 'APPROVED' };
      const futureFortNight = moment().add(14, 'day');
      const ppmObj = {
        planned_move_date: futureFortNight,
        weight_estimate: '10000',
        estimated_incentive: '$24665.59 - 27261.97',
        status: 'SUBMITTED',
      };
      const subComponent = getShallowRender(
        entitlementObj,
        serviceMember,
        ordersObj,
        moveObj,
        ppmObj,
        editMoveFn,
        resumeMoveFn,
      ).find(SubmittedMoveSummary);
      expect(subComponent).not.toBeNull();
      expect(
        subComponent
          .dive()
          .find('.step')
          .find('div.title')
          .first()
          .html(),
      ).toEqual('<div class="title">Next Step: Awaiting approval</div>');
    });
  });
  describe('when a move and ppm are in approved state', () => {
    it('renders approved content', () => {
      const moveObj = { status: 'APPROVED' };
      const futureFortNight = moment().add(14, 'day');
      const ppmObj = {
        planned_move_date: futureFortNight,
        weight_estimate: '10000',
        estimated_incentive: '$24665.59 - 27261.97',
        status: 'APPROVED',
      };
      const subComponent = getShallowRender(
        entitlementObj,
        serviceMember,
        ordersObj,
        moveObj,
        ppmObj,
        editMoveFn,
        resumeMoveFn,
      ).find(ApprovedMoveSummary);
      expect(subComponent).not.toBeNull();
      // expect(
      //   subComponent
      //     .dive()
      //     .find('.step')
      //     .find('div.title')
      //     .first()
      //     .html(),
      // ).toEqual('<div class="title">Next Step: Get ready to move</div>');
    });
  });
  describe('when a move is in in progress state', () => {
    it('renders in progress content', () => {
      const moveObj = { status: 'APPROVED' };
      const pastFortNight = moment().subtract(14, 'day');
      const ppmObj = {
        planned_move_date: pastFortNight,
        weight_estimate: '10000',
        estimated_incentive: '$24665.59 - 27261.97',
      };
      const subComponent = getShallowRender(
        entitlementObj,
        serviceMember,
        ordersObj,
        moveObj,
        ppmObj,
        editMoveFn,
        resumeMoveFn,
      ).find(ApprovedMoveSummary);
      expect(subComponent).not.toBeNull();
      //   expect(
      //     subComponent
      //       .dive()
      //       .find('.step')
      //       .find('div.title')
      //       .first()
      //       .html(),
      //   ).toEqual('<div class="title">Next Step: Request payment</div>');
    });
  });
});

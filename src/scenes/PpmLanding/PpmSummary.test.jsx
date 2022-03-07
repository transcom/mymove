import React from 'react';
import { shallow } from 'enzyme';
import { PpmSummaryComponent } from './PpmSummary';
import moment from 'moment';
import CanceledMoveSummary from 'scenes/PpmLanding/MoveSummary/CanceledMoveSummary';
import DraftMoveSummary from 'scenes/PpmLanding/MoveSummary/DraftMoveSummary';
import SubmittedPpmMoveSummary from 'scenes/PpmLanding/MoveSummary/SubmittedPpmMoveSummary';
import { SHIPMENT_OPTIONS } from 'shared/constants';

describe('PpmSummaryComponent', () => {
  const editMoveFn = jest.fn();
  const resumeMoveFn = jest.fn();
  const entitlementObj = { sum: '10000' };
  const serviceMember = { current_location: { name: 'Ft Carson' } };
  const ordersObj = {};
  const getMoveDocumentsForMove = jest.fn(() => ({ then: () => {} }));
  const getShallowRender = (entitlementObj, serviceMember, ordersObj, moveObj, ppmObj, editMoveFn, resumeMoveFn) => {
    return shallow(
      <PpmSummaryComponent
        entitlement={entitlementObj}
        profile={serviceMember}
        orders={ordersObj}
        move={moveObj}
        ppm={ppmObj}
        editMove={editMoveFn}
        resumeMove={resumeMoveFn}
        getMoveDocumentsForMove={getMoveDocumentsForMove}
      />,
    );
  };

  describe('when a ppm move is in a draft state', () => {
    it('renders resume setup content', () => {
      const moveObj = { selected_move_type: SHIPMENT_OPTIONS.PPM, status: 'DRAFT' };
      const futureFortNight = moment().add(14, 'day');
      const ppmObj = {
        original_move_date: futureFortNight,
        weight_estimate: '10000',
        estimated_incentive: '$24665.59 - 27261.97',
        status: 'CANCELED',
      };
      const subComponent = getShallowRender(
        entitlementObj,
        serviceMember,
        ordersObj,
        moveObj,
        ppmObj,
        editMoveFn,
        resumeMoveFn,
      );
      expect(subComponent.find(DraftMoveSummary).length).toBe(1);
      expect(subComponent.find(DraftMoveSummary).dive().find('.step').find('.title').html()).toEqual(
        '<div class="title">Next Step: Finish setting up your move</div>',
      );
    });
  });
  // PPM
  describe('when a ppm move is in canceled state', () => {
    it('renders cancel content', () => {
      const moveObj = { selected_move_type: SHIPMENT_OPTIONS.PPM, status: 'CANCELED' };
      const futureFortNight = moment().add(14, 'day');
      const ppmObj = {
        original_move_date: futureFortNight,
        weight_estimate: '10000',
        estimated_incentive: '$24665.59 - 27261.97',
        status: 'CANCELED',
      };
      const subComponent = getShallowRender(
        entitlementObj,
        serviceMember,
        ordersObj,
        moveObj,
        ppmObj,
        editMoveFn,
        resumeMoveFn,
      );
      expect(subComponent.find(CanceledMoveSummary).length).toBe(1);
      expect(subComponent.find(CanceledMoveSummary).dive().find('h1').html()).toEqual('<h1>New move</h1>');
    });
  });
  describe('when a move with a ppm is in submitted state', () => {
    it('renders submitted content', () => {
      const moveObj = { selected_move_type: SHIPMENT_OPTIONS.PPM, status: 'SUBMITTED' };
      const futureFortNight = moment().add(14, 'day');
      const ppmObj = {
        original_move_date: futureFortNight,
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
      ).find(SubmittedPpmMoveSummary);
      expect(subComponent.find(SubmittedPpmMoveSummary).length).toBe(1);
      expect(subComponent.find(SubmittedPpmMoveSummary).dive().find('.step').find('div.title').first().html()).toEqual(
        '<div class="title">Next Step: Wait for approval &amp; get ready</div>',
      );
    });
  });

  describe('when a move is in approved state but ppm is submitted state', () => {
    it('renders submitted rather than approved content', () => {
      const moveObj = { selected_move_type: SHIPMENT_OPTIONS.PPM, status: 'APPROVED' };
      const futureFortNight = moment().add(14, 'day');
      const ppmObj = {
        original_move_date: futureFortNight,
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
      ).find(SubmittedPpmMoveSummary);
      expect(subComponent.find(SubmittedPpmMoveSummary).length).toBe(1);
      expect(subComponent.find(SubmittedPpmMoveSummary).dive().find('.step').find('div.title').first().html()).toEqual(
        '<div class="title">Next Step: Wait for approval &amp; get ready</div>',
      );
    });
  });
  describe('when a move and ppm are in approved state', () => {
    it('renders approved content', () => {
      const moveObj = { status: 'APPROVED' };
      const futureFortNight = moment().add(14, 'day');
      const ppmObj = {
        original_move_date: futureFortNight,
        weight_estimate: '10000',
        estimated_incentive: '$24665.59 - 27261.97',
        status: 'APPROVED',
      };
      const component = getShallowRender(
        entitlementObj,
        serviceMember,
        ordersObj,
        moveObj,
        ppmObj,
        editMoveFn,
        resumeMoveFn,
      );
      const ppmSummary = component.find('Connect(ApprovedMoveSummary)');
      expect(ppmSummary.exists()).toBe(true);
    });
  });
  describe('when a move with a ppm is in in progress state', () => {
    it('renders in progress content', () => {
      const moveObj = { status: 'APPROVED' };
      const pastFortNight = moment().subtract(14, 'day');
      const ppmObj = {
        original_move_date: pastFortNight,
        weight_estimate: '10000',
        estimated_incentive: '$24665.59 - 27261.97',
      };
      const component = getShallowRender(
        entitlementObj,
        serviceMember,
        ordersObj,
        moveObj,
        ppmObj,
        editMoveFn,
        resumeMoveFn,
      );
      const ppmSummary = component.find(SubmittedPpmMoveSummary);
      expect(ppmSummary.exists()).toBe(true);
    });
  });
});

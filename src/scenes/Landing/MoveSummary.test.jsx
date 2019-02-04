import React from 'react';
import { shallow } from 'enzyme';
import {
  MoveSummary,
  CanceledMoveSummary,
  ApprovedMoveSummary,
  SubmittedPpmMoveSummary,
  SubmittedHhgMoveSummary,
} from './MoveSummary';
import moment from 'moment';

describe('MoveSummary', () => {
  const editMoveFn = jest.fn();
  const resumeMoveFn = jest.fn();
  const addPPMShipmentFn = jest.fn();
  const entitlementObj = { sum: '10000' };
  const serviceMember = { current_station: { name: 'Ft Carson' } };
  const ordersObj = {};
  const getShallowRender = (
    entitlementObj,
    serviceMember,
    ordersObj,
    moveObj,
    ppmObj,
    hhgObj,
    editMoveFn,
    resumeMoveFn,
    addPPMShipmentFn,
  ) => {
    return shallow(
      <MoveSummary
        entitlement={entitlementObj}
        profile={serviceMember}
        orders={ordersObj}
        move={moveObj}
        ppm={ppmObj}
        shipment={hhgObj}
        editMove={editMoveFn}
        resumeMove={resumeMoveFn}
        addPPMShipment={addPPMShipmentFn}
      />,
    );
  };

  // PPM
  describe('when a ppm move is in canceled state', () => {
    it('renders cancel content', () => {
      const moveObj = { selected_move_type: 'PPM', status: 'CANCELED' };
      const futureFortNight = moment().add(14, 'day');
      const ppmObj = {
        planned_move_date: futureFortNight,
        weight_estimate: '10000',
        estimated_incentive: '$24665.59 - 27261.97',
        status: 'CANCELED',
      };
      const hhgObj = {};
      const subComponent = getShallowRender(
        entitlementObj,
        serviceMember,
        ordersObj,
        moveObj,
        ppmObj,
        hhgObj,
        editMoveFn,
        resumeMoveFn,
      );
      expect(subComponent.find(CanceledMoveSummary).length).toBe(1);
      expect(
        subComponent
          .find(CanceledMoveSummary)
          .dive()
          .find('h2')
          .html(),
      ).toEqual('<h2>New move</h2>');
    });
  });
  describe('when a move with a ppm is in submitted state', () => {
    it('renders submitted content', () => {
      const moveObj = { selected_move_type: 'PPM', status: 'SUBMITTED' };
      const futureFortNight = moment().add(14, 'day');
      const ppmObj = {
        planned_move_date: futureFortNight,
        weight_estimate: '10000',
        estimated_incentive: '$24665.59 - 27261.97',
      };
      const hhgObj = {};
      const subComponent = getShallowRender(
        entitlementObj,
        serviceMember,
        ordersObj,
        moveObj,
        ppmObj,
        hhgObj,
        editMoveFn,
        resumeMoveFn,
      ).find(SubmittedPpmMoveSummary);
      expect(subComponent.find(SubmittedPpmMoveSummary).length).toBe(1);
      expect(
        subComponent
          .find(SubmittedPpmMoveSummary)
          .dive()
          .find('.step')
          .find('div.title')
          .first()
          .html(),
      ).toEqual('<div class="title">Next Step: Wait for approval</div>');
    });
  });
  describe('when a move is in approved state but ppm is submitted state', () => {
    it('renders submitted rather than approved content', () => {
      const moveObj = { selected_move_type: 'PPM', status: 'APPROVED' };
      const futureFortNight = moment().add(14, 'day');
      const ppmObj = {
        planned_move_date: futureFortNight,
        weight_estimate: '10000',
        estimated_incentive: '$24665.59 - 27261.97',
        status: 'SUBMITTED',
      };
      const hhgObj = {};
      const subComponent = getShallowRender(
        entitlementObj,
        serviceMember,
        ordersObj,
        moveObj,
        ppmObj,
        hhgObj,
        editMoveFn,
        resumeMoveFn,
      ).find(SubmittedPpmMoveSummary);
      expect(subComponent.find(SubmittedPpmMoveSummary).length).toBe(1);
      expect(
        subComponent
          .find(SubmittedPpmMoveSummary)
          .dive()
          .find('.step')
          .find('div.title')
          .first()
          .html(),
      ).toEqual('<div class="title">Next Step: Wait for approval</div>');
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
      const hhgObj = {};
      const subComponent = getShallowRender(
        entitlementObj,
        serviceMember,
        ordersObj,
        moveObj,
        ppmObj,
        hhgObj,
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
  describe('when a move with a ppm is in in progress state', () => {
    it('renders in progress content', () => {
      const moveObj = { status: 'APPROVED' };
      const pastFortNight = moment().subtract(14, 'day');
      const ppmObj = {
        planned_move_date: pastFortNight,
        weight_estimate: '10000',
        estimated_incentive: '$24665.59 - 27261.97',
      };
      const hhgObj = {};
      const subComponent = getShallowRender(
        entitlementObj,
        serviceMember,
        ordersObj,
        moveObj,
        ppmObj,
        hhgObj,
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

  // HHG Shipment
  describe('when an hhg move is in canceled state', () => {
    it('renders cancel content', () => {
      const moveObj = { selected_move_type: 'HHG', status: 'CANCELED' };
      const futureFortNight = moment().add(14, 'day');
      const ppmObj = {};
      const hhgObj = {
        requested_move_date: futureFortNight,
        weight_estimate: '10000',
        status: 'CANCELED',
      };
      const subComponent = getShallowRender(
        entitlementObj,
        serviceMember,
        ordersObj,
        moveObj,
        ppmObj,
        hhgObj,
        editMoveFn,
        resumeMoveFn,
      );
      expect(subComponent.find(CanceledMoveSummary).length).toBe(1);
      expect(
        subComponent
          .find(CanceledMoveSummary)
          .dive()
          .find('h2')
          .html(),
      ).toEqual('<h2>New move</h2>');
    });
  });
  describe('when a move and its hhg are in submitted state', () => {
    it('renders submitted content', () => {
      const moveObj = { selected_move_type: 'HHG', status: 'SUBMITTED' };
      const futureFortNight = moment().add(14, 'day');
      const ppmObj = {};
      const hhgObj = {
        planned_move_date: futureFortNight,
        weight_estimate: '10000',
        status: 'SUBMITTED',
      };
      const subComponent = getShallowRender(
        entitlementObj,
        serviceMember,
        ordersObj,
        moveObj,
        ppmObj,
        hhgObj,
        editMoveFn,
        resumeMoveFn,
      ).find(SubmittedHhgMoveSummary);

      expect(subComponent.find(SubmittedHhgMoveSummary).length).toBe(1);
      expect(
        subComponent
          .find(SubmittedHhgMoveSummary)
          .dive()
          .find('.step')
          .find('div.title')
          .first()
          .html(),
      ).toEqual('<div class="title">Next Step: Prepare for move</div>');
    });
  });
  describe('when it is past ten days after a move is submitted', () => {
    it('renders ten days after book date content', () => {
      const pastFortNight = moment().add(-20, 'day');
      const moveObj = { selected_move_type: 'HHG', status: 'SUBMITTED' };
      const ppmObj = {};
      const hhgObj = { status: 'SUBMITTED', book_date: pastFortNight };
      const subComponent = getShallowRender(
        entitlementObj,
        serviceMember,
        ordersObj,
        moveObj,
        ppmObj,
        hhgObj,
        editMoveFn,
        resumeMoveFn,
        addPPMShipmentFn,
      ).find(SubmittedHhgMoveSummary);
      expect(subComponent.find(SubmittedHhgMoveSummary).length).toBe(1);
      expect(
        subComponent
          .find(SubmittedHhgMoveSummary)
          .dive()
          .find('.step')
          .find('div.title')
          .first()
          .html(),
      ).toEqual('<div class="title">Next step: Read pre-move tips</div>');
    });
  });
  describe('when an hhg is in awarded state', () => {
    it('renders submitted content', () => {
      const moveObj = { selected_move_type: 'HHG' };
      const ppmObj = {};
      const hhgObj = { status: 'AWARDED' };
      const subComponent = getShallowRender(
        entitlementObj,
        serviceMember,
        ordersObj,
        moveObj,
        ppmObj,
        hhgObj,
        editMoveFn,
        resumeMoveFn,
      ).find(SubmittedHhgMoveSummary);
      expect(subComponent.find(SubmittedHhgMoveSummary).length).toBe(1);
      expect(
        subComponent
          .dive()
          .find('.step')
          .find('div.title')
          .first()
          .html(),
      ).toEqual('<div class="title">Next Step: Prepare for move</div>');
    });
  });
  describe('when an hhg is in accepted state', () => {
    it('renders submitted content', () => {
      const moveObj = { selected_move_type: 'HHG' };
      const ppmObj = {};
      const hhgObj = { status: 'ACCEPTED' };
      const subComponent = getShallowRender(
        entitlementObj,
        serviceMember,
        ordersObj,
        moveObj,
        ppmObj,
        hhgObj,
        editMoveFn,
        resumeMoveFn,
      ).find(SubmittedHhgMoveSummary);
      expect(subComponent.find(SubmittedHhgMoveSummary).length).toBe(1);
      expect(
        subComponent
          .dive()
          .find('.step')
          .find('div.title')
          .first()
          .html(),
      ).toEqual('<div class="title">Next Step: Prepare for move</div>');
    });
  });
  describe('when an hhg is in approved state', () => {
    it('renders submitted content', () => {
      const moveObj = { selected_move_type: 'HHG' };
      const ppmObj = {};
      const hhgObj = { status: 'APPROVED' };
      const subComponent = getShallowRender(
        entitlementObj,
        serviceMember,
        ordersObj,
        moveObj,
        ppmObj,
        hhgObj,
        editMoveFn,
        resumeMoveFn,
      ).find(SubmittedHhgMoveSummary);
      expect(subComponent.find(SubmittedHhgMoveSummary).length).toBe(1);
      expect(
        subComponent
          .dive()
          .find('.step')
          .find('div.title')
          .first()
          .html(),
      ).toEqual('<div class="title">Next Step: Prepare for move</div>');
    });
  });
  describe('when an hhg is in in_transit state', () => {
    it('renders submitted content', () => {
      const moveObj = { selected_move_type: 'HHG' };
      const ppmObj = {};
      const hhgObj = { status: 'IN_TRANSIT' };
      const subComponent = getShallowRender(
        entitlementObj,
        serviceMember,
        ordersObj,
        moveObj,
        ppmObj,
        hhgObj,
        editMoveFn,
        resumeMoveFn,
      ).find(SubmittedHhgMoveSummary);
      expect(subComponent.find(SubmittedHhgMoveSummary).length).toBe(1);
      expect(
        subComponent
          .dive()
          .find('.step')
          .find('div.title')
          .first()
          .html(),
      ).toEqual('<div class="title">Next Step: Prepare for move</div>');
    });
  });
  describe('when an hhg is in delivered state', () => {
    it('renders submitted content', () => {
      const moveObj = { selected_move_type: 'HHG' };
      const ppmObj = {};
      const hhgObj = { status: 'DELIVERED' };
      const subComponent = getShallowRender(
        entitlementObj,
        serviceMember,
        ordersObj,
        moveObj,
        ppmObj,
        hhgObj,
        editMoveFn,
        resumeMoveFn,
      ).find(SubmittedHhgMoveSummary);
      expect(subComponent.find(SubmittedHhgMoveSummary).length).toBe(1);
      expect(
        subComponent
          .dive()
          .find('.step')
          .find('div.title')
          .first()
          .html(),
      ).toEqual('<div class="title">Next Step: Survey</div>');
    });
  });
  describe('when an hhg is in completed state', () => {
    it('renders submitted content', () => {
      const moveObj = { selected_move_type: 'HHG' };
      const ppmObj = {};
      const hhgObj = { status: 'COMPLETED' };
      const subComponent = getShallowRender(
        entitlementObj,
        serviceMember,
        ordersObj,
        moveObj,
        ppmObj,
        hhgObj,
        editMoveFn,
        resumeMoveFn,
      ).find(SubmittedHhgMoveSummary);
      expect(subComponent.find(SubmittedHhgMoveSummary).length).toBe(1);
      expect(
        subComponent
          .dive()
          .find('.step')
          .find('div.title')
          .first()
          .html(),
      ).toEqual('<div class="title">Next Step: Survey</div>');
    });
  });
});

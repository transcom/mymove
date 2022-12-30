import React from 'react';
import WeightTicket from './WeightTicket';
import { mount } from 'enzyme';
import PPMPaymentRequestActionBtns from './PPMPaymentRequestActionBtns';
import { MockProviders } from 'testUtils';

function mountComponents(
  moreWeightTickets = 'Yes',
  formInvalid,
  uploaderWithInvalidState,
  weightTicketSetType = 'CAR',
) {
  const initialState = {
    emptyWeight: 1100,
    fullWeight: 2000,
    weightTicketSetType: weightTicketSetType,
    weightTicketDate: '2019-05-22',
  };
  const params = { moveId: 'someID' };
  const wrapper = mount(
    <MockProviders params={params} initialState={initialState}>
      <WeightTicket />
    </MockProviders>,
  );
  const wt = wrapper.find('WeightTicket');
  if (formInvalid !== undefined) {
    wt.instance().invalid = jest.fn().mockReturnValue(formInvalid);
    wt.instance().uploaderWithInvalidState = jest.fn().mockReturnValue(uploaderWithInvalidState);
  }
  wt.setState({ additionalWeightTickets: moreWeightTickets, ...initialState });
  wt.update();
  return wrapper.find('WeightTicket');
}

describe('Weight tickets page', () => {
  describe('Service member is missing a weight ticket', () => {
    it('renders both the Save buttons are disabled', () => {
      const weightTicket = mountComponents('No', true, true);
      const buttonGroup = weightTicket.find(PPMPaymentRequestActionBtns);
      const finishLater = weightTicket.find('button').at(0);
      const saveAndAdd = weightTicket.find('button').at(1);

      expect(buttonGroup.length).toEqual(1);
      expect(saveAndAdd.props().disabled).toEqual(true);
      expect(finishLater.props().disabled).not.toEqual(true);
    });
  });
  describe('Service member chooses CAR as weight ticket type', () => {
    it('renders vehicle make and model fields', () => {
      const weightTicket = mountComponents('No', true, true, 'CAR');
      const vehicleNickname = weightTicket.find('[data-testid="vehicle_nickname"]');
      const vehicleMake = weightTicket.find('[data-testid="vehicle_make"]');
      const vehicleModel = weightTicket.find('[data-testid="vehicle_model"]');

      expect(vehicleNickname.length).toEqual(0);
      expect(vehicleMake.length).toEqual(1);
      expect(vehicleModel.length).toEqual(1);
    });
  });
  describe('Service member chooses BOX TRUCK as weight ticket type', () => {
    it('renders vehicle nickname field', () => {
      const weightTicket = mountComponents('No', true, true, 'BOX_TRUCK');
      const vehicleNickname = weightTicket.find('[data-testid="vehicle_nickname"]');
      const vehicleMake = weightTicket.find('[data-testid="vehicle_make"]');
      const vehicleModel = weightTicket.find('[data-testid="vehicle_model"]');

      expect(vehicleNickname.length).toEqual(1);
      expect(vehicleMake.length).toEqual(0);
      expect(vehicleModel.length).toEqual(0);
    });
  });
  describe('Service member chooses PROGEAR as weight ticket type', () => {
    it('renders vehicle nickname (progear type) field', () => {
      const weightTicket = mountComponents('No', true, true, 'PRO_GEAR');
      const vehicleNickname = weightTicket.find('[data-testid="vehicle_nickname"]');
      const vehicleMake = weightTicket.find('[data-testid="vehicle_make"]');
      const vehicleModel = weightTicket.find('[data-testid="vehicle_model"]');

      expect(vehicleNickname.length).toEqual(1);
      expect(vehicleMake.length).toEqual(0);
      expect(vehicleModel.length).toEqual(0);
    });
  });
  describe('Service member has uploaded both a weight tickets', () => {
    it('renders both the Save buttons are enabled', () => {
      const weightTicket = mountComponents('No', false, false);
      const buttonGroup = weightTicket.find(PPMPaymentRequestActionBtns);
      const finishLater = weightTicket.find('button').at(0);
      const saveAndAdd = weightTicket.find('button').at(1);

      expect(buttonGroup.length).toEqual(1);
      expect(saveAndAdd.props().disabled).toEqual(false);
      expect(finishLater.props().disabled).not.toEqual(true);
    });
  });
  describe('Service member answers "Yes" that they have more weight tickets', () => {
    it('renders Save and Add Another Button', () => {
      const weightTicket = mountComponents('Yes');
      const buttonGroup = weightTicket.find(PPMPaymentRequestActionBtns);
      expect(buttonGroup.length).toEqual(1);
      expect(buttonGroup.props().nextBtnLabel).toEqual('Save & Add Another');
    });
  });
  describe('Service member answers "No" that they have more weight tickets', () => {
    it('renders Save and Add Continue Button', () => {
      const weightTicket = mountComponents('No');
      const buttonGroup = weightTicket.find(PPMPaymentRequestActionBtns);
      expect(buttonGroup.length).toEqual(1);
      expect(buttonGroup.props().nextBtnLabel).toEqual('Save & Continue');
    });
  });
});

describe('uploaderWithInvalidState', () => {
  it('returns true when there are no uploaders', () => {
    const weightTicket = mountComponents('No');
    const uploaders = {
      emptyWeight: { uploaderRef: {} },
      fullWeight: { uploaderRef: {} },
      trailer: { uploaderRef: {} },
    };
    weightTicket.instance().uploaders = uploaders;
    uploaders.emptyWeight.uploaderRef.isEmpty = jest.fn(() => false);
    uploaders.emptyWeight.isMissingChecked = jest.fn(() => false);
    uploaders.fullWeight.uploaderRef.isEmpty = jest.fn(() => false);
    uploaders.fullWeight.isMissingChecked = jest.fn(() => false);

    expect(weightTicket.instance().uploaderWithInvalidState()).toEqual(false);
  });
  it('returns false when uploaders have at least one file and isMissing is not checked', () => {
    const weightTicket = mountComponents('No');
    const uploaders = {
      emptyWeight: { uploaderRef: {} },
      fullWeight: { uploaderRef: {} },
      trailer: { uploaderRef: {} },
    };
    uploaders.emptyWeight.uploaderRef.isEmpty = jest.fn(() => false);
    uploaders.emptyWeight.isMissingChecked = jest.fn(() => false);
    uploaders.fullWeight.uploaderRef.isEmpty = jest.fn(() => false);
    uploaders.fullWeight.isMissingChecked = jest.fn(() => false);
    weightTicket.instance().uploaders = uploaders;

    expect(weightTicket.instance().uploaderWithInvalidState()).toEqual(false);
  });
  it('returns true when uploaders have at least one file and isMissing is checked', () => {
    const weightTicket = mountComponents('No');
    const uploaders = {
      emptyWeight: { uploaderRef: {} },
      fullWeight: { uploaderRef: {} },
      trailer: { uploaderRef: {} },
    };
    uploaders.emptyWeight.uploaderRef.isEmpty = jest.fn(() => false);
    uploaders.emptyWeight.isMissingChecked = jest.fn(() => false);
    uploaders.fullWeight.uploaderRef.isEmpty = jest.fn(() => false);
    uploaders.fullWeight.isMissingChecked = jest.fn(() => true);
    weightTicket.instance().uploaders = uploaders;

    expect(weightTicket.instance().uploaderWithInvalidState()).toEqual(true);
  });
});

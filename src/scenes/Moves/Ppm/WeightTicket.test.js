import React from 'react';
import WeightTicket from './WeightTicket';
import { mount } from 'enzyme';
import { Provider } from 'react-redux';
import store from 'shared/store';
import MockRouter from 'react-mock-router';
import PPMPaymentRequestActionBtns from './PPMPaymentRequestActionBtns';

function mountComponents(moreWeightTickets = 'Yes', formInvalid, uploderWithInvalidState) {
  const initialValues = {
    empty_weight: 1100,
    full_weight: 2000,
    vehicle_nickname: 'KIRBY',
    weight_ticket_set_type: 'CAR',
    weight_ticket_date: '2019-05-22',
  };
  const match = { params: { moveId: 'someID' } };
  const wrapper = mount(
    <Provider store={store}>
      <MockRouter push={jest.fn()}>
        <WeightTicket match={match} />
      </MockRouter>
    </Provider>,
  );
  const wt = wrapper.find('WeightTicket');
  if (formInvalid !== undefined) {
    wt.instance().invalid = jest.fn().mockReturnValue(formInvalid);
    wt.instance().uploaderWithInvalidState = jest.fn().mockReturnValue(uploderWithInvalidState);
  }
  wt.setState({ additionalWeightTickets: moreWeightTickets, initialValues: initialValues });
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

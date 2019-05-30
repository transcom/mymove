import React from 'react';
import WeightTicket from './WeightTicket';
import { mount } from 'enzyme';
import { Provider } from 'react-redux';
import store from 'shared/store';
import MockRouter from 'react-mock-router';
import PPMPaymentRequestActionBtns from './PPMPaymentRequestActionBtns';

function mountComponents(moreWeightTickets = 'Yes', missingWeightTickets) {
  const initialValues = {
    empty_weight: 1100,
    full_weight: 2000,
    vehicle_nickname: 'HALE',
    vehicle_options: 'CAR',
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
  if (missingWeightTickets !== undefined) {
    wt.instance().isMissingWeightTickets = jest.fn().mockReturnValue(missingWeightTickets);
    wt.instance().formIsComplete = jest.fn().mockReturnValue(true);
  }
  wt.setState({ additionalWeightTickets: moreWeightTickets, initialValues: initialValues });
  wt.update();
  return wrapper.find('WeightTicket');
}

describe('Weight tickets page', () => {
  describe('Service member is missing a weight ticket', () => {
    it('renders both the Save buttons are disabled', () => {
      const weightTicket = mountComponents('No', true);
      const buttonGroup = weightTicket.find(PPMPaymentRequestActionBtns);
      const cancel = weightTicket.find('button').at(0);
      const saveForLater = weightTicket.find('button').at(1);
      const saveAnd = weightTicket.find('button').at(2);

      expect(buttonGroup.length).toEqual(1);
      expect(cancel.props().disabled).not.toEqual(true);
      expect(saveAnd.props().disabled).toEqual(true);
      expect(saveForLater.props().disabled).toEqual(true);
    });
  });
  describe('Service member has uploaded both a weight tickets', () => {
    it('renders both the Save buttons are enabled', () => {
      const weightTicket = mountComponents('No', false);
      const buttonGroup = weightTicket.find(PPMPaymentRequestActionBtns);
      const cancel = weightTicket.find('button').at(0);
      const saveForLater = weightTicket.find('button').at(1);
      const saveAnd = weightTicket.find('button').at(2);

      expect(buttonGroup.length).toEqual(1);
      expect(cancel.props().disabled).not.toEqual(true);
      expect(saveAnd.props().disabled).toEqual(false);
      expect(saveForLater.props().disabled).toEqual(false);
    });
  });

  describe('Service member hasnt provided an Empty Weight weight ticket', () => {
    it('renders both the Save buttons are disabled', () => {
      const weightTicket = mountComponents('No');
      const buttonGroup = weightTicket.find(PPMPaymentRequestActionBtns);
      const cancel = weightTicket.find('button').at(0);
      const saveForLater = weightTicket.find('button').at(1);
      const saveAnd = weightTicket.find('button').at(2);

      expect(buttonGroup.length).toEqual(1);
      expect(cancel.props().disabled).not.toEqual(true);
      expect(saveAnd.props().disabled).toEqual(true);
      expect(saveForLater.props().disabled).toEqual(true);
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

describe('missingWeightTickets', () => {
  it('returns true when there are no uploaders', () => {
    const weightTicket = mountComponents('No');
    weightTicket.instance().uploaders = {};

    expect(weightTicket.instance().isMissingWeightTickets()).toEqual(true);
  });
  it('returns false when both uploaders have at least one file', () => {
    const weightTicket = mountComponents('No');
    const uploaders = { one: {}, two: {} };
    uploaders['one'].isEmpty = jest.fn(() => false);
    uploaders['two'].isEmpty = jest.fn(() => false);
    weightTicket.instance().uploaders = uploaders;

    expect(weightTicket.instance().isMissingWeightTickets()).toEqual(false);
  });
  it('returns true when one uploaders do not have at least one file', () => {
    const weightTicket = mountComponents('No');
    const uploaders = { one: {}, two: {} };
    uploaders['one'].isEmpty = jest.fn(() => false);
    uploaders['two'].isEmpty = jest.fn(() => true);
    weightTicket.instance().uploaders = uploaders;

    expect(weightTicket.instance().isMissingWeightTickets()).toEqual(true);
  });
});

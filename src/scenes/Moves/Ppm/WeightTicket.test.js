import React from 'react';
import WeightTicket from './WeightTicket';
import { mount } from 'enzyme';
import { Provider } from 'react-redux';
import store from 'shared/store';
import MockRouter from 'react-mock-router';
import PPMPaymentRequestActionBtns from './PPMPaymentRequestActionBtns';

function mountComponents(moreWeightTickets = 'Yes') {
  const wrapper = mount(
    <Provider store={store}>
      <MockRouter push={jest.fn()}>
        <WeightTicket />
      </MockRouter>
    </Provider>,
  );
  const wt = wrapper.find('WeightTicket');
  wt.setState({ additionalWeightTickets: moreWeightTickets });
  wt.update();
  return wrapper.find('WeightTicket');
}

describe('Weight tickets page', () => {
  describe('Service member answers "Yes" that they have more weight tickets', () => {
    it('renders Save and Add Another Button', () => {
      const weightTicket = mountComponents('Yes');
      const button = weightTicket.find(PPMPaymentRequestActionBtns);
      expect(button.length).toEqual(1);
      expect(button.props().nextBtnLabel).toEqual('Save & Add Another');
    });
  });
  describe('Service member answers "No" that they have more weight tickets', () => {
    it('renders Save and Add Continue Button', () => {
      const weightTicket = mountComponents('No');
      const button = weightTicket.find(PPMPaymentRequestActionBtns);
      expect(button.length).toEqual(1);
      expect(button.props().nextBtnLabel).toEqual('Save & Continue');
    });
  });
});

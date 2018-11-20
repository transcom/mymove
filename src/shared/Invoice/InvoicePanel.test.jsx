import React from 'react';
import { mount } from 'enzyme';

import { InvoicePanel } from './InvoicePanel';
import { isOfficeSite } from 'shared/constants.js';
import * as CONSTANTS from 'shared/constants.js';

describe('InvoicePanel tests', () => {
  let wrapper,
    shipmentState = null;
  const shipmentLineItems = [''];
  beforeEach(() => {
    CONSTANTS.isOfficeSite = true;
    shipmentState = 'DELIVERED';
    wrapper = mount(<InvoicePanel shipmentLineItems={shipmentLineItems} shipmentState={shipmentState} />);
  });

  describe('When no items exist', () => {
    it('renders without crashing', () => {
      expect(wrapper.find('.empty-content').length).toEqual(1);
    });
  });

  describe('Approve Payment button shows on delivered state and office app', () => {
    it('renders enabled "Approve Payment" button', () => {
      expect(isOfficeSite).toBe(true);
      expect(wrapper.props().shipmentState).toBe('DELIVERED');

      wrapper.update();
      expect(
        wrapper
          .children()
          .containsMatchingElement(<button className="button button-secondary">Approve Payment</button>),
      ).toBeTruthy();
    });

    it('renders disabled "Approve Payment" button', () => {
      expect(isOfficeSite).toBe(true);
      expect(wrapper.props().shipmentState).toBe('DELIVERED');

      wrapper.update();
      expect(wrapper.find('button').prop('disabled')).toBeTruthy();
    });
  });
});

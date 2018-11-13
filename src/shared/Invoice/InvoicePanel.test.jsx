import React from 'react';
import { mount } from 'enzyme';

import { InvoicePanel } from './InvoicePanel';

describe('InvoicePanel tests', () => {
  let wrapper;
  const shipmentLineItems = [''];
  beforeEach(() => {
    wrapper = mount(<InvoicePanel shipmentLineItems={shipmentLineItems} />);
  });

  describe('When no items exist', () => {
    it('renders without crashing', () => {
      expect(wrapper.find('.empty-content').length).toEqual(1);
    });
  });
});

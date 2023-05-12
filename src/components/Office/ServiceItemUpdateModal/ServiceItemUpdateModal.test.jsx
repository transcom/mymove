import React from 'react';
import { render } from '@testing-library/react';

import ServiceItemUpdateModal from './ServiceItemUpdateModal';

describe('ServiceItemUpdateModal', () => {
  describe('when a editing an address change to a service items', () => {
    render(<ServiceItemUpdateModal />);
    it('renders modal with the correct content', async () => {});
    it('the form is editing as expected', async () => {});
    it('shows error messages appear when form validations are not met', async () => {});
    it('when the save button is pressed, the onSave handler is called', async () => {});
    it('when the cancel button is pressed, the onCancel handler is called', async () => {});
  });
  describe('when a reviewing an address change to a service items', () => {
    render(<ServiceItemUpdateModal />);
    it('renders modal with the correct content', async () => {});
    it('the form is editing as expected', async () => {});
    it('shows error messages appear when form validations are not met', async () => {});
    it('when the save button is pressed, the onSave handler is called', async () => {});
    it('when the cancel button is pressed, the onCancel handler is called', async () => {});
  });
});

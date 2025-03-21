import React from 'react';
import { render, waitFor, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { Provider } from 'react-redux';

import { EditFacilityInfoModal } from './EditFacilityInfoModal';

import { configureStore } from 'shared/store';

describe('EditFacilityInfoModal', () => {
  const storageFacility = {
    facilityName: 'My Facility',
    phone: '915-555-2942',
    email: 'my@email.com',
    address: {
      streetAddress1: '123 Fake Street',
      streetAddress2: '',
      city: 'Pasadena',
      state: 'CA',
      postalCode: '90210',
    },
    lotNumber: '11232',
  };
  const incompleteStorageFacility = {
    facilityName: 'My Facility',
    phone: '915-555-2942',
    email: 'my@email.com',
    address: {
      streetAddress1: '',
      streetAddress2: '',
      city: 'Pasadena',
      state: 'CA',
      postalCode: '90210',
    },
    lotNumber: '11232',
  };

  it('calls onSubmit prop on save button click when the form has valid data', async () => {
    const mockOnSubmit = jest.fn();
    const mockStore = configureStore({});
    render(
      <Provider store={mockStore.store}>
        <EditFacilityInfoModal
          onClose={() => {}}
          onSubmit={mockOnSubmit}
          storageFacility={storageFacility}
          serviceOrderNumber="12345"
          shipmentType="HHG_INTO_NTS"
        />
      </Provider>,
    );
    const submitBtn = screen.getByRole('button', { name: 'Save' });

    await userEvent.click(submitBtn);

    await waitFor(() => {
      expect(mockOnSubmit).toHaveBeenCalled();
    });
  });

  it('calls onSubmit prop on save button click when valid data is entered', async () => {
    const mockOnSubmit = jest.fn();
    const mockStore = configureStore({});
    render(
      <Provider store={mockStore.store}>
        <EditFacilityInfoModal
          onClose={() => {}}
          onSubmit={mockOnSubmit}
          storageFacility={incompleteStorageFacility}
          serviceOrderNumber="12345"
          shipmentType="HHG_INTO_NTS"
        />
      </Provider>,
    );
    const addressInput = screen.getByLabelText(/Address 1/);
    const submitBtn = screen.getByRole('button', { name: 'Save' });

    await userEvent.type(addressInput, '123 Fake Street');
    await userEvent.click(submitBtn);

    await waitFor(() => {
      expect(mockOnSubmit).toHaveBeenCalled();
    });
  });

  it('does not allow saving with incomplete form data', async () => {
    const mockStore = configureStore({});
    render(
      <Provider store={mockStore.store}>
        <EditFacilityInfoModal
          onClose={() => {}}
          onSubmit={() => {}}
          storageFacility={incompleteStorageFacility}
          serviceOrderNumber="12345"
          shipmentType="HHG_INTO_NTS"
        />
      </Provider>,
    );
    const submitBtn = screen.getByRole('button', { name: 'Save' });
    await waitFor(() => {
      expect(submitBtn).toBeDisabled();
    });
  });

  it('calls onclose prop on modal close', async () => {
    const mockClose = jest.fn();
    const mockStore = configureStore({});
    render(
      <Provider store={mockStore.store}>
        <EditFacilityInfoModal
          onClose={mockClose}
          onSubmit={() => {}}
          storageFacility={storageFacility}
          serviceOrderNumber="12345"
          shipmentType="HHG_INTO_NTS"
        />
      </Provider>,
    );
    const closeBtn = screen.getByRole('button', { name: 'Cancel' });

    await userEvent.click(closeBtn);

    await waitFor(() => {
      expect(mockClose).toHaveBeenCalled();
    });
  });
});

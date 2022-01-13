import React from 'react';
import { render, waitFor, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import { EditFacilityInfoModal } from './EditFacilityInfoModal';

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
    render(
      <EditFacilityInfoModal
        onClose={() => {}}
        onSubmit={mockOnSubmit}
        storageFacility={storageFacility}
        serviceOrderNumber="12345"
        shipmentType="HHG_INTO_NTS_DOMESTIC"
      />,
    );
    const submitBtn = screen.getByRole('button', { name: 'Save' });

    userEvent.click(submitBtn);

    await waitFor(() => {
      expect(mockOnSubmit).toHaveBeenCalled();
    });
  });

  it('calls onSubmit prop on save button click when valid data is entered', async () => {
    const mockOnSubmit = jest.fn();
    render(
      <EditFacilityInfoModal
        onClose={() => {}}
        onSubmit={mockOnSubmit}
        storageFacility={incompleteStorageFacility}
        serviceOrderNumber="12345"
        shipmentType="HHG_INTO_NTS_DOMESTIC"
      />,
    );
    const addressInput = screen.getByLabelText('Address 1');
    const submitBtn = screen.getByRole('button', { name: 'Save' });

    userEvent.type(addressInput, '123 Fake Street');
    userEvent.click(submitBtn);

    await waitFor(() => {
      expect(mockOnSubmit).toHaveBeenCalled();
    });
  });

  it('does not allow saving with incomplete form data', async () => {
    render(
      <EditFacilityInfoModal
        onClose={() => {}}
        onSubmit={() => {}}
        storageFacility={incompleteStorageFacility}
        serviceOrderNumber="12345"
        shipmentType="HHG_INTO_NTS_DOMESTIC"
      />,
    );
    const submitBtn = screen.getByRole('button', { name: 'Save' });
    await waitFor(() => {
      expect(submitBtn).toBeDisabled();
    });
  });

  it('calls onclose prop on modal close', async () => {
    const mockClose = jest.fn();
    render(
      <EditFacilityInfoModal
        onClose={mockClose}
        onSubmit={() => {}}
        storageFacility={storageFacility}
        serviceOrderNumber="12345"
        shipmentType="HHG_INTO_NTS_DOMESTIC"
      />,
    );
    const closeBtn = screen.getByRole('button', { name: 'Cancel' });

    userEvent.click(closeBtn);

    await waitFor(() => {
      expect(mockClose).toHaveBeenCalled();
    });
  });
});

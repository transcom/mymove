import React from 'react';
import { render, screen, within, act, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import { ServiceItemUpdateModal } from './ServiceItemUpdateModal';
import {
  domesticDestinationSitServiceItem,
  dddSitWithAddressUpdate,
  newAddress,
} from './ServiceItemUpdateModalTestParams';
import EditSitAddressChangeForm from './EditSitAddressChangeForm';

const defaultValues = {
  closeModal: () => {},
  onSave: () => {},
};

describe('ServiceItemUpdateModal', () => {
  describe('renders base component with the shared components', () => {
    it('renders modal with the correct content', async () => {
      render(
        <ServiceItemUpdateModal
          title="Title for the modal"
          serviceItem={domesticDestinationSitServiceItem}
          {...defaultValues}
        />,
      );
      const officeRemarks = screen.getByLabelText('Office remarks');
      const sitDetailsTable = await screen.findByTestId('sitAddressUpdateDetailTable');
      const saveButton = screen.getByRole('button', { name: 'Save' });
      const cancelButton = screen.getByRole('button', { name: 'Cancel' });
      expect(screen.getByText('HHG')).toBeInTheDocument();
      expect(screen.getByText('Title for the modal')).toBeInTheDocument();
      expect(within(sitDetailsTable).getByText('Customer contact 1:')).toBeInTheDocument();
      expect(within(sitDetailsTable).getByText('1200Z')).toBeInTheDocument();
      expect(within(sitDetailsTable).getByText('First available delivery date 1:')).toBeInTheDocument();
      expect(within(sitDetailsTable).getByText('15 Sep 2020')).toBeInTheDocument();
      expect(within(sitDetailsTable).getByText('Customer contact 2:')).toBeInTheDocument();
      expect(within(sitDetailsTable).getByText('2300Z')).toBeInTheDocument();
      expect(within(sitDetailsTable).getByText('First available delivery date 2:')).toBeInTheDocument();
      expect(within(sitDetailsTable).getByText('21 Sep 2020')).toBeInTheDocument();
      expect(within(sitDetailsTable).getByText('Reason:')).toBeInTheDocument();
      expect(screen.getByText("Customer's housing at base is not ready")).toBeInTheDocument();
      expect(screen.getByText('SIT delivery address')).toBeInTheDocument();
      expect(officeRemarks).toBeInTheDocument();
      expect(saveButton).toBeInTheDocument();
      expect(cancelButton).toBeInTheDocument();
    });
    it('when the save button is pressed, the onSave handler is called', async () => {});
    it('when the cancel button is pressed, the onCancel handler is called', async () => {});
  });
  describe('when a editing an address change to a service items', () => {
    it('renders modal with the correct content', async () => {
      render(
        <ServiceItemUpdateModal
          title="Edit service item"
          serviceItem={dddSitWithAddressUpdate}
          {...defaultValues}
          initialValues={{ newAddress }}
        >
          <EditSitAddressChangeForm initialAddress={newAddress} />
        </ServiceItemUpdateModal>,
      );
      expect(screen.getByText('Edit service item')).toBeInTheDocument();
      expect(screen.getByText('Final SIT delivery')).toBeInTheDocument();
      expect(screen.getByText('Initial SIT delivery address')).toBeInTheDocument();
      expect(screen.queryByText('555 Fakest Dr,')).toBeInTheDocument();
      expect(screen.queryByText('Unit 133,')).toBeInTheDocument();
      expect(screen.queryByText('Alexandria, VA 12867')).toBeInTheDocument();
      const form = screen.getByTestId('editAddressForm');
      expect(form).toBeInTheDocument();
    });
    it('the form is editing and submits as expected', async () => {
      const mockOnSubmit = jest.fn();
      render(
        <ServiceItemUpdateModal
          title="Edit service item"
          serviceItem={dddSitWithAddressUpdate}
          onSave={mockOnSubmit}
          closeModal={() => {}}
          initialValues={{ newAddress }}
        >
          <EditSitAddressChangeForm initialAddress={newAddress} />
        </ServiceItemUpdateModal>,
      );
      const address1 = screen.getByLabelText('Address 1');
      const officeRemarksInput = screen.getByLabelText('Office remarks');
      const submitBtn = screen.getByRole('button', { name: 'Save' });
      await act(() => userEvent.clear(address1));
      await act(() => userEvent.type(address1, '123 Fake Street'));
      await act(() => userEvent.type(officeRemarksInput, 'Approved!'));
      await act(() => userEvent.click(submitBtn));
      await waitFor(() => {
        expect(mockOnSubmit).toHaveBeenCalled();
        expect(mockOnSubmit).toHaveBeenCalledWith('abc123', {
          officeRemarks: 'Approved!',
          newAddress: {
            city: 'Alexandria',
            state: 'VA',
            postalCode: '12867',
            streetAddress1: '123 Fake Street',
            streetAddress2: 'Unit 133',
            streetAddress3: '',
            country: 'USA',
          },
        });
      });
    });
    it('Save button is disabled if form validations are not met', async () => {
      render(
        <ServiceItemUpdateModal
          title="Edit service item"
          serviceItem={dddSitWithAddressUpdate}
          {...defaultValues}
          initialValues={{ newAddress }}
        >
          <EditSitAddressChangeForm initialAddress={newAddress} />
        </ServiceItemUpdateModal>,
      );
      const officeRemarksInput = screen.getByLabelText('Office remarks');
      const submitBtn = screen.getByRole('button', { name: 'Save' });
      // Testing Office remarks validation.
      await act(() => userEvent.clear(officeRemarksInput));
      await waitFor(() => {
        expect(submitBtn).toBeDisabled();
      });
    });
    it('when the cancel button is pressed, the onCancel handler is called', async () => {
      const mockOnClose = jest.fn();
      render(
        <ServiceItemUpdateModal
          title="Edit service item"
          serviceItem={dddSitWithAddressUpdate}
          onSave={() => {}}
          closeModal={mockOnClose}
          initialValues={{ newAddress }}
        >
          <EditSitAddressChangeForm initialAddress={newAddress} />
        </ServiceItemUpdateModal>,
      );
      const cancelButton = screen.getByRole('button', { name: 'Cancel' });
      await act(() => userEvent.click(cancelButton));
      await waitFor(() => {
        expect(mockOnClose).toHaveBeenCalled();
      });
    });
  });
  // describe('when a reviewing an address change to a service items', () => {
  //   render(<ServiceItemUpdateModal />);
  //   it('renders modal with the correct content', async () => {});
  //   it('the form is editing as expected', async () => {});
  //   it('shows error messages appear when form validations are not met', async () => {});
  //   it('when the save button is pressed, the onSave handler is called', async () => {});
  //   it('when the cancel button is pressed, the onCancel handler is called', async () => {});
  // });
});

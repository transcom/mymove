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
import ReviewSitAddressChange from './ReviewSitAddressChange';

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

  describe('When a TOO reviews a SIT address update request', () => {
    it('clicking review request renders the associated modal with the correct content', async () => {
      render(
        <ServiceItemUpdateModal
          title="Review request: service item update"
          {...defaultValues}
          serviceItem={{ dddSitWithAddressUpdate }}
        >
          <ReviewSitAddressChange sitAddressUpdate={dddSitWithAddressUpdate.sitAddressUpdates[0]} />
        </ServiceItemUpdateModal>,
      );

      // Checking modal header
      expect(screen.getByText('Review request: service item update')).toBeInTheDocument();
      expect(screen.getByTestId('sitAddressUpdateTag')).toHaveTextContent('UPDATE REQUESTED');

      // Check for yellow alert
      expect(screen.getByTestId('distanceAlert')).toBeInTheDocument();
      expect(
        screen.getByText(
          'Requested final SIT delivery address is 500 miles from the initial SIT delivery address. Approvals over 50 miles will result in updated pricing for this shipment.',
        ),
      ).toBeInTheDocument();

      // Check for address section
      expect(screen.getByText('Initial SIT delivery address')).toBeInTheDocument();
      expect(screen.getByText('345 Faker Rd, Richmond, VA 12508')).toBeInTheDocument();

      expect(screen.getByText('Requested final SIT delivery address')).toBeInTheDocument();
      expect(screen.getByText('555 Fakest Dr, Unit 133, Alexandria, VA 12867')).toBeInTheDocument();

      // Check for remarks section
      expect(screen.getByText('Update request details')).toBeInTheDocument();
      expect(screen.getByText('Contractor remarks:')).toBeInTheDocument();
      expect(screen.getByText('Customer wishes to be closer to family')).toBeInTheDocument();

      // Check for radio button section
      expect(screen.getByText('Review Request')).toBeInTheDocument();
      expect(screen.getByText('Approve address change?')).toBeInTheDocument();
      const form = screen.getByTestId('reviewSITAddressUpdateForm');
      expect(form).toBeInTheDocument();
    });

    it('the form is editing and submits as expected for yes button', async () => {
      const mockOnSubmit = jest.fn();
      render(
        <ServiceItemUpdateModal
          title="Review request: service item update"
          onSave={mockOnSubmit}
          closeModal={() => {}}
          serviceItem={dddSitWithAddressUpdate}
        >
          <ReviewSitAddressChange sitAddressUpdate={dddSitWithAddressUpdate.sitAddressUpdates[0]} />
        </ServiceItemUpdateModal>,
      );
      const approveSITAddressUpdateBtn = screen.getByRole('radio', { name: /yes/i });
      const officeRemarksInput = screen.getByLabelText('Office remarks');
      const submitBtn = screen.getByRole('button', { name: 'Save' });
      await act(() => userEvent.click(approveSITAddressUpdateBtn));
      await act(() => userEvent.type(officeRemarksInput, 'Approved!'));
      await act(() => userEvent.click(submitBtn));
      await waitFor(() => {
        expect(mockOnSubmit).toHaveBeenCalled();
        expect(mockOnSubmit).toHaveBeenCalledWith('abc123', {
          officeRemarks: 'Approved!',
          sitAddressUpdate: 'YES',
        });
      });
    });
    it('the form is editing and submits as expected for no button', async () => {
      const mockOnSubmit = jest.fn();
      render(
        <ServiceItemUpdateModal
          title="Review request: service item update"
          onSave={mockOnSubmit}
          closeModal={() => {}}
          serviceItem={dddSitWithAddressUpdate}
        >
          <ReviewSitAddressChange sitAddressUpdate={dddSitWithAddressUpdate.sitAddressUpdates[0]} />
        </ServiceItemUpdateModal>,
      );
      const approveSITAddressUpdateBtn = screen.getByRole('radio', { name: /no/i });
      const officeRemarksInput = screen.getByLabelText('Office remarks');
      const submitBtn = screen.getByRole('button', { name: 'Save' });
      await act(() => userEvent.click(approveSITAddressUpdateBtn));
      await act(() => userEvent.type(officeRemarksInput, 'Rejected!'));
      await act(() => userEvent.click(submitBtn));
      await waitFor(() => {
        expect(mockOnSubmit).toHaveBeenCalled();
        expect(mockOnSubmit).toHaveBeenCalledWith('abc123', {
          officeRemarks: 'Rejected!',
          sitAddressUpdate: 'NO',
        });
      });
    });
    it('when the cancel button is pressed, the onCancel handler is called', async () => {
      const mockOnClose = jest.fn();
      render(
        <ServiceItemUpdateModal
          title="Review request: service item update"
          onSave={() => {}}
          closeModal={mockOnClose}
          serviceItem={{ dddSitWithAddressUpdate }}
        >
          <ReviewSitAddressChange sitAddressUpdate={dddSitWithAddressUpdate.sitAddressUpdates[0]} />
        </ServiceItemUpdateModal>,
      );
      const cancelButton = screen.getByRole('button', { name: 'Cancel' });
      await act(() => userEvent.click(cancelButton));
      await waitFor(() => {
        expect(mockOnClose).toHaveBeenCalled();
      });
    });

    it('Save button is disabled if form validations are not met', async () => {
      render(
        <ServiceItemUpdateModal
          title="Review request: service item update"
          {...defaultValues}
          serviceItem={{ dddSitWithAddressUpdate }}
        >
          <ReviewSitAddressChange sitAddressUpdate={dddSitWithAddressUpdate.sitAddressUpdates[0]} />
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
  });
});

import React from 'react';
import { render, screen, waitFor, act } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { Provider } from 'react-redux';

import EditPPMHeaderSummaryModal from './EditPPMHeaderSummaryModal';

import { configureStore } from 'shared/store';

let onClose;
let onSubmit;
beforeEach(() => {
  onClose = jest.fn();
  onSubmit = jest.fn();
});

afterEach(() => {
  jest.clearAllMocks();
});

describe('EditPPMHeaderSummaryModal', () => {
  const sectionInfo = {
    actualMoveDate: '2022-01-01',
    advanceAmountReceived: 50000,
    destinationAddressObj: {
      city: 'Fairfield',
      country: 'US',
      id: '672ff379-f6e3-48b4-a87d-796713f8f997',
      postalCode: '94535',
      state: 'CA',
      streetAddress1: '987 Any Avenue',
      streetAddress2: 'P.O. Box 9876',
      streetAddress3: 'c/o Some Person',
      county: 'SOLANO',
    },
    pickupAddressObj: {
      city: 'Beverly Hills',
      country: 'US',
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi4zODQ3Njla',
      id: '1686751b-ab36-43cf-b3c9-c0f467d13c19',
      postalCode: '90210',
      state: 'CA',
      streetAddress1: '123 Any Street',
      streetAddress2: 'P.O. Box 12345',
      streetAddress3: 'c/o Some Person',
      county: 'LOS ANGELES',
    },
  };

  it('renders the component', async () => {
    await act(async () => {
      render(
        <EditPPMHeaderSummaryModal
          sectionType="shipmentInfo"
          sectionInfo={sectionInfo}
          onClose={onClose}
          onSubmit={onSubmit}
          editItemName="actualMoveDate"
        />,
      );
    });

    expect(await screen.findByRole('heading', { level: 3, name: 'Edit Shipment Info' })).toBeInTheDocument();
    expect(screen.getByLabelText('Actual move start date')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Save' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Cancel' })).toBeInTheDocument();
    expect(screen.getByLabelText('Close')).toBeInstanceOf(HTMLButtonElement);
  });

  it('renders pickup address', async () => {
    const mockStore = configureStore({});

    await act(async () => {
      render(
        <Provider store={mockStore.store}>
          <EditPPMHeaderSummaryModal
            sectionType="shipmentInfo"
            sectionInfo={sectionInfo}
            onClose={onClose}
            onSubmit={onSubmit}
            editItemName="pickupAddress"
          />
        </Provider>,
      );
    });

    expect(await screen.findByRole('heading', { level: 3, name: 'Edit Shipment Info' })).toBeInTheDocument();
    expect(screen.getByText('Pickup Address')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Save' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Cancel' })).toBeInTheDocument();
    expect(screen.getByLabelText('Close')).toBeInstanceOf(HTMLButtonElement);
  });

  it('renders destination address', async () => {
    const mockStore = configureStore({});

    await act(async () => {
      render(
        <Provider store={mockStore.store}>
          <EditPPMHeaderSummaryModal
            sectionType="shipmentInfo"
            sectionInfo={sectionInfo}
            onClose={onClose}
            onSubmit={onSubmit}
            editItemName="destinationAddress"
          />
        </Provider>,
      );
    });

    expect(await screen.findByRole('heading', { level: 3, name: 'Edit Shipment Info' })).toBeInTheDocument();
    expect(screen.getByText('Destination Address')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Save' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Cancel' })).toBeInTheDocument();
    expect(screen.getByLabelText('Close')).toBeInstanceOf(HTMLButtonElement);
  });

  it('renders actual expense reimbursement', async () => {
    await act(async () => {
      render(
        <EditPPMHeaderSummaryModal
          sectionType="shipmentInfo"
          sectionInfo={sectionInfo}
          onClose={onClose}
          onSubmit={onSubmit}
          editItemName="isActualExpenseReimbursement"
        />,
      );
    });

    expect(await screen.findByRole('heading', { level: 3, name: 'Edit Shipment Info' })).toBeInTheDocument();
    expect(screen.getByText('Is this PPM an Actual Expense Reimbursement?')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Save' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Cancel' })).toBeInTheDocument();
    expect(screen.getByLabelText('Close')).toBeInstanceOf(HTMLButtonElement);
  });

  it('closes the modal when close icon is clicked', async () => {
    await act(async () => {
      render(
        <EditPPMHeaderSummaryModal
          sectionType="shipmentInfo"
          sectionInfo={sectionInfo}
          onClose={onClose}
          onSubmit={onSubmit}
          editItemName="actualMoveDate"
        />,
      );
    });

    await act(async () => {
      await userEvent.click(await screen.getByLabelText('Close'));
    });

    await waitFor(() => {
      expect(onClose).toHaveBeenCalledTimes(1);
    });
  });

  it('closes the modal when the cancel button is clicked', async () => {
    await act(async () => {
      render(
        <EditPPMHeaderSummaryModal
          sectionType="shipmentInfo"
          sectionInfo={sectionInfo}
          onClose={onClose}
          onSubmit={onSubmit}
          editItemName="actualMoveDate"
        />,
      );
    });

    await act(async () => {
      await userEvent.click(await screen.getByRole('button', { name: 'Cancel' }));
    });

    await waitFor(() => {
      expect(onClose).toHaveBeenCalledTimes(1);
    });
  });

  it('calls the submit function when submit button is clicked', async () => {
    await act(async () => {
      render(
        <EditPPMHeaderSummaryModal
          sectionType="shipmentInfo"
          sectionInfo={sectionInfo}
          onClose={onClose}
          onSubmit={onSubmit}
          editItemName="actualMoveDate"
        />,
      );
    });

    await act(async () => {
      await userEvent.click(await screen.getByRole('button', { name: 'Save' }));
    });

    await waitFor(() => {
      expect(onSubmit).toHaveBeenCalled();
    });
  });

  it('displays required validation error when actual move date is empty', async () => {
    await act(async () => {
      render(
        <EditPPMHeaderSummaryModal
          sectionType="shipmentInfo"
          sectionInfo={{ ...sectionInfo, actualMoveDate: '' }}
          onClose={onClose}
          onSubmit={onSubmit}
          editItemName="actualMoveDate"
        />,
      );
    });

    await act(async () => {
      await userEvent.clear(await screen.getByLabelText('Actual move start date'));
      await userEvent.click(await screen.getByRole('button', { name: 'Save' }));
    });

    expect(await screen.findByText('Required')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Save' })).toHaveAttribute('disabled');
  });

  it('displays required validation error when advance amount received is empty', async () => {
    await act(async () => {
      render(
        <EditPPMHeaderSummaryModal
          sectionType="incentives"
          sectionInfo={{ advanceAmountReceived: '' }}
          onClose={onClose}
          onSubmit={onSubmit}
          editItemName="advanceAmountReceived"
        />,
      );
    });

    await act(async () => {
      await userEvent.clear(await screen.getByLabelText('Advance received'));
      await userEvent.click(await screen.getByRole('button', { name: 'Save' }));
    });

    expect(await screen.findByText('Required')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Save' })).toHaveAttribute('disabled');
  });
});

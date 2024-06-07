import React from 'react';
import { render, screen, waitFor, act } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import EditPPMHeaderSummaryModal from './EditPPMHeaderSummaryModal';

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
          sectionInfo={sectionInfo}
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

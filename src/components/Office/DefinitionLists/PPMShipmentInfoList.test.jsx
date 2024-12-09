import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import PPMShipmentInfoList from './PPMShipmentInfoList';

import affiliation from 'content/serviceMemberAgencies';
import { MockProviders } from 'testUtils';
import { permissionTypes } from 'constants/permissions';
import { ADVANCE_STATUSES } from 'constants/ppms';
import { ppmShipmentStatuses } from 'constants/shipments';
import { downloadPPMAOAPacket, downloadPPMPaymentPacket } from 'services/ghcApi';

jest.mock('services/ghcApi', () => ({
  ...jest.requireActual('services/ghcApi'),
  downloadPPMAOAPacket: jest.fn(),
  downloadPPMPaymentPacket: jest.fn(),
}));

afterEach(() => {
  jest.resetAllMocks();
});

const renderWithPermissions = (shipment) => {
  render(
    <MockProviders permissions={[permissionTypes.viewCloseoutOffice]}>
      <PPMShipmentInfoList isExpanded shipment={shipment} />
    </MockProviders>,
  );
};

describe('PPMShipmentInfoList', () => {
  it('renders closeout display for Marines', () => {
    renderWithPermissions({ agency: affiliation.MARINES });
    expect(screen.getByTestId('closeout')).toBeInTheDocument();
    expect(screen.getByTestId('closeout').textContent).toEqual('TVCB');
    expect(screen.getByText('Closeout office')).toBeInTheDocument();
  });

  it('renders closeout display for Navy', () => {
    renderWithPermissions({ agency: affiliation.NAVY });
    expect(screen.getByTestId('closeout')).toBeInTheDocument();
    expect(screen.getByTestId('closeout').textContent).toEqual('NAVY');
    expect(screen.getByText('Closeout office')).toBeInTheDocument();
  });

  it('renders closeout display Coast guard', () => {
    renderWithPermissions({ agency: affiliation.COAST_GUARD });
    expect(screen.getByTestId('closeout')).toBeInTheDocument();
    expect(screen.getByTestId('closeout').textContent).toEqual('USCG');
    expect(screen.getByText('Closeout office')).toBeInTheDocument();
  });

  it('renders closeout display for Army and Air Force', () => {
    renderWithPermissions({ closeoutOffice: 'Test office' });
    expect(screen.getByTestId('closeout')).toBeInTheDocument();
    expect(screen.getByTestId('closeout').textContent).toEqual('Test office');
    expect(screen.getByText('Closeout office')).toBeInTheDocument();
  });

  it('renders closeout display when there is no closeout office', () => {
    renderWithPermissions({ closeoutOffice: '—' });
    expect(screen.getByTestId('closeout')).toBeInTheDocument();
    expect(screen.getByTestId('closeout').textContent).toEqual('—');
    expect(screen.getByText('Closeout office')).toBeInTheDocument();
  });

  it('renders estimated and max incentives', () => {
    renderWithPermissions({ ppmShipment: { estimatedIncentive: 100000, maxIncentive: 200000 } });
    expect(screen.getByTestId('estimatedIncentive')).toBeInTheDocument();
    expect(screen.getByText('Estimated Incentive')).toBeInTheDocument();

    expect(screen.getByTestId('maxIncentive')).toBeInTheDocument();
    expect(screen.getByText('Max Incentive')).toBeInTheDocument();
  });

  it('PPM Download AOA Paperwork - success with Approved', async () => {
    const mockResponse = {
      ok: true,
      headers: {
        'content-disposition': 'filename="test.pdf"',
      },
      status: 200,
      data: null,
    };
    downloadPPMAOAPacket.mockImplementation(() => Promise.resolve(mockResponse));

    renderWithPermissions({ ppmShipment: { advanceStatus: ADVANCE_STATUSES.APPROVED.apiValue } });

    expect(screen.getByText('Download AOA Paperwork (PDF)', { exact: false })).toBeInTheDocument();

    const downloadAOAButton = screen.getByText('Download AOA Paperwork (PDF)');
    expect(downloadAOAButton).toBeInTheDocument();

    await userEvent.click(downloadAOAButton);

    await waitFor(() => {
      expect(downloadPPMAOAPacket).toHaveBeenCalledTimes(1);
    });
  });

  it('PPM Download AOA Paperwork - success with Edited', async () => {
    const mockResponse = {
      ok: true,
      headers: {
        'content-disposition': 'filename="test.pdf"',
      },
      status: 200,
      data: null,
    };
    downloadPPMAOAPacket.mockImplementation(() => Promise.resolve(mockResponse));

    renderWithPermissions({ ppmShipment: { advanceStatus: ADVANCE_STATUSES.EDITED.apiValue } });

    expect(screen.getByText('Download AOA Paperwork (PDF)', { exact: false })).toBeInTheDocument();

    const downloadAOAButton = screen.getByText('Download AOA Paperwork (PDF)');
    expect(downloadAOAButton).toBeInTheDocument();

    await userEvent.click(downloadAOAButton);

    await waitFor(() => {
      expect(downloadPPMAOAPacket).toHaveBeenCalledTimes(1);
    });
  });

  it('PPM Download AOA Paperwork - failure', async () => {
    downloadPPMAOAPacket.mockRejectedValue({
      response: { body: { title: 'Error title', detail: 'Error detail' } },
    });

    const shipment = { ppmShipment: { advanceStatus: ADVANCE_STATUSES.APPROVED.apiValue } };
    const onErrorHandler = jest.fn();

    render(
      <MockProviders permissions={[permissionTypes.viewCloseoutOffice]}>
        <PPMShipmentInfoList shipment={shipment} onErrorModalToggle={onErrorHandler} />
      </MockProviders>,
    );

    expect(screen.getByText('Download AOA Paperwork (PDF)')).toBeInTheDocument();

    const downloadAOAButton = screen.getByText('Download AOA Paperwork (PDF)');
    expect(downloadAOAButton).toBeInTheDocument();
    await userEvent.click(downloadAOAButton);

    await waitFor(() => {
      expect(downloadPPMAOAPacket).toHaveBeenCalledTimes(1);
      expect(onErrorHandler).toHaveBeenCalledTimes(1);
    });
  });

  it('PPM Download Payment Paperwork - success', async () => {
    const mockResponse = {
      ok: true,
      headers: {
        'content-disposition': 'filename="test.pdf"',
      },
      status: 200,
      data: null,
    };
    downloadPPMPaymentPacket.mockImplementation(() => Promise.resolve(mockResponse));

    renderWithPermissions({ ppmShipment: { status: ppmShipmentStatuses.CLOSEOUT_COMPLETE } });

    expect(screen.getByText('Download Payment Packet (PDF)', { exact: false })).toBeInTheDocument();

    const downloadPaymentButton = screen.getByText('Download Payment Packet (PDF)');
    expect(downloadPaymentButton).toBeInTheDocument();

    await userEvent.click(downloadPaymentButton);

    await waitFor(() => {
      expect(downloadPPMPaymentPacket).toHaveBeenCalledTimes(1);
    });
  });

  it('PPM Download Payment Packet - failure', async () => {
    downloadPPMPaymentPacket.mockRejectedValue({
      response: { body: { title: 'Error title', detail: 'Error detail' } },
    });

    const shipment = { ppmShipment: { status: ppmShipmentStatuses.CLOSEOUT_COMPLETE } };
    const onErrorHandler = jest.fn();

    render(
      <MockProviders permissions={[permissionTypes.viewCloseoutOffice]}>
        <PPMShipmentInfoList shipment={shipment} onErrorModalToggle={onErrorHandler} />
      </MockProviders>,
    );

    expect(screen.getByText('Download Payment Packet (PDF)')).toBeInTheDocument();

    const downloadPaymentButton = screen.getByText('Download Payment Packet (PDF)');
    expect(downloadPaymentButton).toBeInTheDocument();
    await userEvent.click(downloadPaymentButton);

    await waitFor(() => {
      expect(downloadPPMPaymentPacket).toHaveBeenCalledTimes(1);
      expect(onErrorHandler).toHaveBeenCalledTimes(1);
    });
  });

  it('renders actual move date', () => {
    renderWithPermissions({
      ppmShipment: { expectedDepartureDate: '2024-07-20T09:48:21.573Z', actualMoveDate: '2024-07-22T09:48:21.573Z' },
    });
    expect(screen.getByTestId('actualDepartureDate')).toBeInTheDocument();
    expect(screen.getByText('Actual Departure date')).toBeInTheDocument();
    expect(screen.getByText('22 Jul 2024')).toBeInTheDocument();
  });

  it('renders estimated move date', () => {
    renderWithPermissions({ ppmShipment: { expectedDepartureDate: '2024-07-20T09:48:21.573Z' } });
    expect(screen.getByTestId('expectedDepartureDate')).toBeInTheDocument();
    expect(screen.getByText('Estimated Departure date')).toBeInTheDocument();
    expect(screen.getByText('20 Jul 2024')).toBeInTheDocument();
  });
});

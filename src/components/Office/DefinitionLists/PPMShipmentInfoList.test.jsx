import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';

import PPMShipmentInfoList from './PPMShipmentInfoList';

import affiliation from 'content/serviceMemberAgencies';
import { MockProviders } from 'testUtils';
import { permissionTypes } from 'constants/permissions';
import { ADVANCE_STATUSES } from 'constants/ppms';
import { downloadPPMAOAPacket } from 'services/ghcApi';
import { downloadPPMAOAPacketOnSuccessHandler } from 'utils/download';
import userEvent from '@testing-library/user-event';

jest.mock('services/ghcApi', () => ({
  ...jest.requireActual('services/ghcApi'),
  downloadPPMAOAPacket: jest.fn(),
}));

jest.mock('utils/download', () => ({
  ...jest.requireActual('utils/download'),
  downloadPPMAOAPacketOnSuccessHandler: jest.fn(),
}));

afterEach(() => {
  jest.resetAllMocks();
});

const renderWithPermissions = (shipment) => {
  render(
    <MockProviders permissions={[permissionTypes.viewCloseoutOffice]}>
      <PPMShipmentInfoList shipment={shipment} />
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

  it('PPM Download AOA Paperwork - success', async () => {
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
      expect(downloadPPMAOAPacketOnSuccessHandler).toHaveBeenCalledTimes(1);
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
      expect(downloadPPMAOAPacketOnSuccessHandler).toHaveBeenCalledTimes(0);
      expect(onErrorHandler).toHaveBeenCalledTimes(1);
    });
  });
});

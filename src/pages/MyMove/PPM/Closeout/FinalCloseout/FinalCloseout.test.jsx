import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { useParams } from 'react-router-dom';
import { useDispatch } from 'react-redux';

import { generalRoutes } from 'constants/routes';
import FinalCloseout from 'pages/MyMove/PPM/Closeout/FinalCloseout/FinalCloseout';
import { selectMTOShipmentById } from 'store/entities/selectors';
import { updateMTOShipment } from 'store/entities/actions';
import { MockProviders, setUpProvidersWithHistory } from 'testUtils';
import { createPPMShipmentWithFinalIncentive } from 'utils/test/factories/ppmShipment';
import { submitPPMShipmentSignedCertification } from 'services/internalApi';
import { ppmSubmissionCertificationText } from 'scenes/Legalese/legaleseText';
import { formatDateForSwagger } from 'shared/dates';

jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useParams: jest.fn(),
}));

jest.mock('react-redux', () => ({
  ...jest.requireActual('react-redux'),
  useDispatch: jest.fn(),
}));

jest.mock('services/internalApi', () => ({
  ...jest.requireActual('services/internalApi'),
  getResponseError: jest.fn(),
  submitPPMShipmentSignedCertification: jest.fn(),
}));

jest.mock('store/entities/selectors', () => ({
  ...jest.requireActual('store/entities/selectors'),
  selectMTOShipmentById: jest.fn(),
}));

jest.mock('store/entities/actions', () => ({
  ...jest.requireActual('store/entities/actions'),
  updateMTOShipment: jest.fn(),
}));

beforeEach(() => {
  jest.clearAllMocks();
});

describe('Final Closeout page', () => {
  const setUpMocks = (mtoShipment) => {
    const shipment = mtoShipment || createPPMShipmentWithFinalIncentive();

    useParams.mockReturnValue({ mtoShipmentId: shipment.id });

    selectMTOShipmentById.mockReturnValueOnce(shipment);

    return shipment;
  };

  it('loads the selected shipment from redux', () => {
    const shipment = setUpMocks();

    render(<FinalCloseout />, { wrapper: MockProviders });

    expect(selectMTOShipmentById).toHaveBeenCalledWith(expect.anything(), shipment.id);
  });

  it('renders the page headings', () => {
    setUpMocks();

    render(<FinalCloseout />, { wrapper: MockProviders });

    expect(screen.getByTestId('tag')).toHaveTextContent('PPM');

    expect(screen.getByRole('heading', { level: 1, name: 'Complete PPM' })).toBeInTheDocument();
    expect(screen.getByRole('heading', { level: 2, name: /Your final estimated incentive: \$/ })).toBeInTheDocument();
  });

  it('routes to the home page when the return to homepage link is clicked', async () => {
    setUpMocks();

    const { memoryHistory, mockProviderWithHistory } = setUpProvidersWithHistory();

    render(<FinalCloseout />, { wrapper: mockProviderWithHistory });

    await userEvent.click(screen.getByText('Return To Homepage'));

    expect(memoryHistory.location.pathname).toEqual(generalRoutes.HOME_PATH);
  });

  it('submits the ppm signed certification', async () => {
    const mockShipment = setUpMocks();
    submitPPMShipmentSignedCertification.mockResolvedValueOnce(mockShipment.ppmShipment);

    const mockDispatch = jest.fn();
    useDispatch.mockImplementation(mockDispatch);

    const { memoryHistory, mockProviderWithHistory } = setUpProvidersWithHistory();

    render(<FinalCloseout />, { wrapper: mockProviderWithHistory });

    await userEvent.type(screen.getByRole('textbox', { name: 'Signature' }), 'Grace Griffin');
    await userEvent.click(screen.getByRole('button', { name: 'Submit PPM Documentation' }));

    await waitFor(() =>
      expect(submitPPMShipmentSignedCertification).toHaveBeenCalledWith(mockShipment.ppmShipment.id, {
        certification_text: ppmSubmissionCertificationText,
        signature: 'Grace Griffin',
        date: formatDateForSwagger(new Date()),
      }),
    );

    expect(updateMTOShipment).toHaveBeenCalledWith(mockShipment);

    expect(memoryHistory.location.pathname).toEqual(generalRoutes.HOME_PATH);
  });
});

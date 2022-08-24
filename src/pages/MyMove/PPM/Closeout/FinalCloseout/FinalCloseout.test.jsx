import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { useParams } from 'react-router-dom';

import { generalRoutes } from 'constants/routes';
import FinalCloseout from 'pages/MyMove/PPM/Closeout/FinalCloseout/FinalCloseout';
// import { getResponseError, patchMTOShipment } from 'services/internalApi';
import { selectMTOShipmentById } from 'store/entities/selectors';
import { MockProviders, setUpProvidersWithHistory } from 'testUtils';
import { createPPMShipmentWithFinalIncentive } from 'utils/test/factories/ppmShipment';

jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useParams: jest.fn(),
}));

jest.mock('services/internalApi', () => ({
  ...jest.requireActual('services/internalApi'),
  getResponseError: jest.fn(),
  patchMTOShipment: jest.fn(),
}));

jest.mock('store/entities/selectors', () => ({
  ...jest.requireActual('store/entities/selectors'),
  selectMTOShipmentById: jest.fn(),
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

  it('routes to the home page when the finish later link is clicked', async () => {
    setUpMocks();

    const { memoryHistory, mockProviderWithHistory } = setUpProvidersWithHistory();

    render(<FinalCloseout />, { wrapper: mockProviderWithHistory });

    await userEvent.click(screen.getByText('Finish Later'));

    expect(memoryHistory.location.pathname).toEqual(generalRoutes.HOME_PATH);
  });
});

import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { generatePath } from 'react-router-dom';

import { selectMTOShipmentById } from 'store/entities/selectors';
import EstimatedIncentive from 'pages/MyMove/PPM/Booking/EstimatedIncentive/EstimatedIncentive';
import { MockProviders } from 'testUtils';
import { customerRoutes } from 'constants/routes';

const mockMoveId = 'move123';
const mockShipmentId = 'shipment123';

const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
}));
const mockRoutingParams = {
  moveId: mockMoveId,
  mtoShipmentId: mockShipmentId,
};
const mockRoutingConfig = {
  path: customerRoutes.SHIPMENT_PPM_ESTIMATED_INCENTIVE_PATH,
  params: mockRoutingParams,
};

jest.mock('store/entities/selectors', () => ({
  selectMTOShipmentById: jest.fn(),
}));

beforeEach(() => {
  jest.clearAllMocks();
});

const shipmentEntity = {
  id: 'shipment123',
  ppmShipment: {
    pickupPostalCode: '10001',
    destinationPostalCode: '10002',
    expectedDepartureDate: '2022-04-01',
    estimatedWeight: 4567,
    estimatedIncentive: 789000,
  },
};

describe('EstimatedIncentive component', () => {
  it('loads the selected shipment from redux', () => {
    selectMTOShipmentById.mockReturnValue(shipmentEntity);

    render(
      <MockProviders {...mockRoutingConfig}>
        <EstimatedIncentive />
      </MockProviders>,
    );
    expect(selectMTOShipmentById).toHaveBeenCalledWith(expect.anything(), mockShipmentId);
  });

  it('renders the shipment tag and page title', () => {
    render(
      <MockProviders {...mockRoutingConfig}>
        <EstimatedIncentive />
      </MockProviders>,
    );
    expect(screen.getByTestId('tag')).toHaveTextContent('PPM');
    expect(screen.getByRole('heading', { level: 1 })).toHaveTextContent('Estimated incentive');
    expect(screen.getByRole('heading', { level: 2 })).toHaveTextContent('$7,890 is your estimated incentive');
  });

  it('renders the buttons and navigates to previous and next routes', async () => {
    render(
      <MockProviders {...mockRoutingConfig}>
        <EstimatedIncentive />
      </MockProviders>,
    );

    const shipmentInfo = {
      moveId: mockMoveId,
      mtoShipmentId: mockShipmentId,
    };

    await userEvent.click(screen.getByRole('button', { name: 'Back' }));
    expect(mockNavigate).toHaveBeenCalledWith(
      generatePath(customerRoutes.SHIPMENT_PPM_ESTIMATED_WEIGHT_PATH, shipmentInfo),
    );

    await userEvent.click(screen.getByRole('button', { name: 'Next' }));
    expect(mockNavigate).toHaveBeenCalledWith(generatePath(customerRoutes.SHIPMENT_PPM_ADVANCES_PATH, shipmentInfo));
  });
});

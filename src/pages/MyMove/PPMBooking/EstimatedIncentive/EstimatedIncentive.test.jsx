import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { generatePath } from 'react-router';

import { selectMTOShipmentById } from 'store/entities/selectors';
import EstimatedIncentive from 'pages/MyMove/PPMBooking/EstimatedIncentive/EstimatedIncentive';
import { MockProviders } from 'testUtils';
import { customerRoutes } from 'constants/routes';

const mockPush = jest.fn();
const mockBack = jest.fn();
const mockMoveId = 'move123';
const mockShipmentId = 'shipment123';

jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useHistory: () => ({
    push: mockPush,
    goBack: mockBack,
  }),
  useParams: () => ({
    moveId: mockMoveId,
    mtoShipmentId: mockShipmentId,
  }),
}));

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
      <MockProviders>
        <EstimatedIncentive />
      </MockProviders>,
    );
    expect(selectMTOShipmentById).toHaveBeenCalledWith(expect.anything(), mockShipmentId);
  });

  it('renders the shipment tag and page title', () => {
    render(
      <MockProviders>
        <EstimatedIncentive />
      </MockProviders>,
    );
    expect(screen.getByTestId('tag')).toHaveTextContent('PPM');
    expect(screen.getByRole('heading', { level: 1 })).toHaveTextContent('Estimated incentive');
    expect(screen.getByRole('heading', { level: 2 })).toHaveTextContent('$7,890 is your estimated incentive');
  });

  it('renders the buttons and navigates to previous and next routes', () => {
    render(
      <MockProviders>
        <EstimatedIncentive />
      </MockProviders>,
    );

    userEvent.click(screen.getByRole('button', { name: 'Back' }));

    expect(mockBack).toHaveBeenCalled();

    userEvent.click(screen.getByRole('button', { name: 'Next' }));

    expect(mockPush).toHaveBeenCalledWith(
      generatePath(customerRoutes.SHIPMENT_PPM_ADVANCES_PATH, {
        moveId: mockMoveId,
        mtoShipmentId: mockShipmentId,
      }),
    );
  });
});

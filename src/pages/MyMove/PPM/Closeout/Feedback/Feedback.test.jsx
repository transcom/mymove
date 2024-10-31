import React from 'react';
import { render, screen } from '@testing-library/react';
import { v4 } from 'uuid';

import Feedback, { GetTripWeight, FormatRow } from './Feedback';

import { MockProviders } from 'testUtils';
import { selectMTOShipmentById } from 'store/entities/selectors';
import { customerRoutes } from 'constants/routes';

const mockMoveId = v4();
const mockMTOShipmentId = v4();

const mockRoutingConfig = {
  path: customerRoutes.SHIPMENT_PPM_FEEDBACK_PATH,
  params: {
    moveId: mockMoveId,
    mtoShipmentId: mockMTOShipmentId,
  },
};

const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
}));

const mockMTOShipment = {
  ppmShipment: {
    actualDestinationPostalCode: '20889',
    actualMoveDate: '2024-05-08',
    actualPickupPostalCode: '59402',
    movingExpenses: [],
    proGearWeightTickets: [],
    w2Address: {
      city: 'Missoula',
      county: 'MISSOULA',
      id: '44fdfd2c-215c-48a0-8d41-065dbe38885b',
      postalCode: '59801',
      state: 'MT',
      streetAddress1: '422 Dearborn Ave',
    },
    weightTickets: [
      {
        emptyWeight: 1999,
        fullWeight: 5844,
      },
    ],
  },
};

jest.mock('store/entities/selectors', () => ({
  ...jest.requireActual('store/entities/selectors'),
  selectMTOShipmentById: jest.fn(() => mockMTOShipment),
}));

beforeEach(() => {
  jest.clearAllMocks();
});

const renderFeedbackPage = () => {
  return render(
    <MockProviders {...mockRoutingConfig}>
      <Feedback />
    </MockProviders>,
  );
};

describe('Feedback page', () => {
  it('displays PPM details', () => {
    renderFeedbackPage();

    expect(selectMTOShipmentById).toHaveBeenCalledWith(expect.anything(), mockMTOShipmentId);
    expect(screen.getByText('About Your PPM')).toBeInTheDocument();
    expect(screen.getByText('Departure Date: 08 May 2024')).toBeInTheDocument();
    expect(screen.getByText('Starting ZIP: 59402')).toBeInTheDocument();
    expect(screen.getByText('Ending ZIP: 20889')).toBeInTheDocument();
    expect(screen.getByText('Advance: No')).toBeInTheDocument();
    expect(screen.getByTestId('w-2Address')).toHaveTextContent('W-2 address: 422 Dearborn AveMissoula, MT 59801');
  });

  it('formats and displays trip weight', () => {
    renderFeedbackPage();

    expect(screen.getByText('Trip weight:')).toBeInTheDocument();
    expect(screen.getByText('3,845 lbs')).toBeInTheDocument();
  });

  it('does not display pro-gear if no pro-gear documents are present', () => {
    renderFeedbackPage();

    expect(screen.queryByTestId('pro-gear-items')).not.toBeInTheDocument();
  });

  it('does not display expenses if no expense documents are present', () => {
    renderFeedbackPage();

    expect(screen.queryByTestId('expenses-items')).not.toBeInTheDocument();
  });

  it('returns correct trip weight', () => {
    const doc = { fullWeight: 5844, emptyWeight: 1999 };
    expect(GetTripWeight(doc)).toBe(3845);
  });

  it('formats row correctly', () => {
    const row = { value: 1000, format: (val) => `$${val}` };
    const formattedRow = FormatRow(row);
    expect(formattedRow.value).toBe('$1000');
  });

  it('displays loading placeholder when mtoShipment is not present', () => {
    selectMTOShipmentById.mockReturnValueOnce(null);
    renderFeedbackPage();
    expect(screen.getByText('Loading...')).toBeInTheDocument();
  });
});

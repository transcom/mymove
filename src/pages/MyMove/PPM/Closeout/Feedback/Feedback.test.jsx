import React from 'react';
import { render, screen } from '@testing-library/react';
import { v4 } from 'uuid';

import Feedback from './Feedback';

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

const mockMTOShipmentWithAdvance = {
  ppmShipment: {
    actualDestinationPostalCode: '20889',
    actualMoveDate: '2024-05-08',
    actualPickupPostalCode: '59402',
    hasReceivedAdvance: true,
    advanceAmountReceived: 100000,
    movingExpenses: [{ id: 'exp1', amount: 5000 }],
    proGearWeightTickets: [{ id: 'pg1', weight: 75 }],
    w2Address: {
      city: 'Missoula',
      county: 'MISSOULA',
      id: '44fdfd2c-215c-48a0-8d41-065dbe38885b',
      postalCode: '59801',
      state: 'MT',
      streetAddress1: '422 Dearborn Ave',
    },
    weightTickets: [{ id: 'wt1', fullWeight: 3000, emptyWeight: 1500 }],
  },
};

jest.mock('store/entities/selectors', () => ({
  ...jest.requireActual('store/entities/selectors'),
  selectMTOShipmentById: jest.fn(() => mockMTOShipment),
}));

beforeEach(() => {
  jest.clearAllMocks();
});

const renderFeedbackPage = (mockData) => {
  // Set the mock selector to return the specified mock data
  selectMTOShipmentById.mockReturnValue(mockData);
  return render(
    <MockProviders {...mockRoutingConfig}>
      <Feedback />
    </MockProviders>,
  );
};

describe('Feedback page', () => {
  it('displays PPM details', () => {
    renderFeedbackPage(mockMTOShipment);

    expect(selectMTOShipmentById).toHaveBeenCalledWith(expect.anything(), mockMTOShipmentId);
    expect(screen.getByText('About Your PPM')).toBeInTheDocument();
    expect(screen.getByText('Departure Date: 08 May 2024')).toBeInTheDocument();
    expect(screen.getByText('Starting ZIP: 59402')).toBeInTheDocument();
    expect(screen.getByText('Ending ZIP: 20889')).toBeInTheDocument();
    expect(screen.getByText('Advance: No')).toBeInTheDocument();
    expect(screen.getByTestId('w-2Address')).toHaveTextContent('W-2 address: 422 Dearborn AveMissoula, MT 59801');
  });

  it('formats and diplays trip weight', () => {
    renderFeedbackPage(mockMTOShipment);

    expect(screen.getByText('Trip weight:')).toBeInTheDocument();
    expect(screen.getByText('3,845 lbs')).toBeInTheDocument();
  });

  it('does not display pro-gear if no pro-gear documents are present', () => {
    renderFeedbackPage(mockMTOShipment);

    expect(screen.queryByTestId('pro-gear-items')).not.toBeInTheDocument();
  });

  it('does not display expenses if no expense documents are present', () => {
    renderFeedbackPage(mockMTOShipment);

    expect(screen.queryByTestId('expenses-items')).not.toBeInTheDocument();
  });
});

describe('Additional Feedback page tests for improved coverage', () => {
  it('displays weight moved section if weight tickets are present', () => {
    renderFeedbackPage(mockMTOShipmentWithAdvance);

    const weightMovedHeading = screen.getByText('Weight Moved');
    const weightMovedValue = weightMovedHeading.closest('.headingContent').querySelector('span');

    expect(weightMovedValue).toHaveTextContent('1,500 lbs');
  });

  it('displays pro-gear section if pro-gear weight tickets are present', () => {
    renderFeedbackPage(mockMTOShipmentWithAdvance);

    expect(screen.getByTestId('pro-gear-items')).toBeInTheDocument();
    expect(screen.getByText('Pro-gear')).toBeInTheDocument();
    expect(screen.getByText('75 lbs')).toBeInTheDocument();
  });

  it('displays expenses section if moving expenses are present', () => {
    renderFeedbackPage(mockMTOShipmentWithAdvance);

    expect(screen.getByTestId('expenses-items')).toBeInTheDocument();
    expect(screen.getByText('Expenses')).toBeInTheDocument();
    expect(screen.getByText('- $50.00')).toBeInTheDocument();
  });
});

import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import { v4 } from 'uuid';
import { createMemoryHistory } from 'history';
import userEvent from '@testing-library/user-event';
import { generatePath } from 'react-router-dom';

import { MockProviders } from 'testUtils';
import { selectMTOShipmentById } from 'store/entities/selectors';
import Review from 'pages/MyMove/PPM/Closeout/Review/Review';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { customerRoutes } from 'constants/routes';
import { deleteWeightTicket } from 'services/internalApi';

const mockMoveId = v4();
const mockMTOShipmentId = v4();
const mockPPMShipmentId = v4();
const mockWeightTicketId = v4();

const mockMTOShipment = {
  id: mockMTOShipmentId,
  shipmentType: SHIPMENT_OPTIONS.PPM,
  ppmShipment: {
    id: mockPPMShipmentId,
    actualMoveDate: '2022-05-01',
    actualPickupPostalCode: '10003',
    actualDestinationPostalCode: '10004',
    advanceReceived: true,
    advanceAmountReceived: '6000000',
    pickupPostalCode: '10001',
    destinationPostalCode: '10002',
    expectedDepartureDate: '2022-04-30',
    advanceRequested: true,
    advanceAmountRequested: 598700,
    estimatedWeight: 4000,
    estimatedIncentive: 1000000,
    sitExpected: false,
    hasProGear: false,
    proGearWeight: null,
    spouseProGearWeight: null,
  },
  eTag: 'dGVzdGluZzIzNDQzMjQ',
};

const mockMTOShipmentWithWeightTicket = {
  id: mockMTOShipmentId,
  shipmentType: SHIPMENT_OPTIONS.PPM,
  ppmShipment: {
    id: mockPPMShipmentId,
    actualMoveDate: '2022-05-01',
    actualPickupPostalCode: '10003',
    actualDestinationPostalCode: '10004',
    advanceReceived: true,
    advanceAmountReceived: '6000000',
    pickupPostalCode: '10001',
    destinationPostalCode: '10002',
    expectedDepartureDate: '2022-04-30',
    advanceRequested: true,
    advanceAmountRequested: 598700,
    estimatedWeight: 4000,
    estimatedIncentive: 1000000,
    sitExpected: false,
    hasProGear: false,
    proGearWeight: null,
    spouseProGearWeight: null,
    weightTickets: [
      {
        id: mockWeightTicketId,
        vehicleDescription: 'DMC Delorean',
        emptyWeight: 2500,
        fullWeight: 3500,
        eTag: 'eTag value',
      },
    ],
  },
  eTag: 'dGVzdGluZzIzNDQzMjQ',
};

const mockMTOShipmentWithIncompleteWeightTicket = {
  id: mockMTOShipmentId,
  shipmentType: SHIPMENT_OPTIONS.PPM,
  ppmShipment: {
    id: mockPPMShipmentId,
    actualMoveDate: '2022-05-01',
    actualPickupPostalCode: '10003',
    actualDestinationPostalCode: '10004',
    advanceReceived: true,
    advanceAmountReceived: '6000000',
    pickupPostalCode: '10001',
    destinationPostalCode: '10002',
    expectedDepartureDate: '2022-04-30',
    advanceRequested: true,
    advanceAmountRequested: 598700,
    estimatedWeight: 4000,
    estimatedIncentive: 1000000,
    sitExpected: false,
    hasProGear: false,
    proGearWeight: null,
    spouseProGearWeight: null,
    weightTickets: [
      {
        id: mockWeightTicketId,
        eTag: 'eTag value',
      },
    ],
  },
  eTag: 'dGVzdGluZzIzNDQzMjQ',
};

const mockPush = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useHistory: () => ({
    push: mockPush,
  }),
  useParams: () => ({
    moveId: mockMoveId,
    mtoShipmentId: mockMTOShipmentId,
  }),
}));

jest.mock('services/internalApi', () => ({
  ...jest.requireActual('services/internalApi'),
  deleteWeightTicket: jest.fn(() => {}),
}));

jest.mock('store/entities/selectors', () => ({
  ...jest.requireActual('store/entities/selectors'),
  selectMTOShipmentById: jest.fn(() => mockMTOShipment),
}));

beforeEach(() => {
  jest.clearAllMocks();
});

describe('About page', () => {
  it('loads the selected shipment from redux', () => {
    render(<Review />, { wrapper: MockProviders });

    expect(selectMTOShipmentById).toHaveBeenCalledWith(expect.anything(), mockMTOShipmentId);
  });

  it('renders the page headings', () => {
    render(<Review />, { wrapper: MockProviders });

    expect(screen.getByTestId('tag')).toHaveTextContent('PPM');
    expect(screen.getByRole('heading', { level: 1 })).toHaveTextContent('Review');
    expect(screen.getAllByRole('heading', { level: 2 })[0]).toHaveTextContent('About Your PPM');
    expect(screen.getAllByRole('heading', { level: 2 })[1]).toHaveTextContent('Documents');
    expect(screen.getAllByRole('heading', { level: 3 })[0]).toHaveTextContent('Weight moved');
    expect(screen.getAllByRole('heading', { level: 3 })[1]).toHaveTextContent('Pro-gear');
    expect(screen.getAllByRole('heading', { level: 3 })[2]).toHaveTextContent('Expenses');
  });

  it('renders the empty message when there are not weight tickets', () => {
    render(<Review />, { wrapper: MockProviders });

    expect(
      screen.getByText('No weight tickets uploaded. Add at least one set of weight tickets to request payment.'),
    ).toBeInTheDocument();
  });

  it('routes to the edit about your ppm page when the edit link is clicked', async () => {
    const editAboutYourPPM = generatePath(customerRoutes.SHIPMENT_PPM_ABOUT_PATH, {
      moveId: mockMoveId,
      mtoShipmentId: mockMTOShipmentId,
    });
    const memoryHistory = createMemoryHistory();
    const mockProviderWithHistory = ({ children }) => <MockProviders history={memoryHistory}>{children}</MockProviders>;
    render(<Review />, { wrapper: mockProviderWithHistory });

    userEvent.click(screen.getAllByText('Edit')[0]);

    await waitFor(() => {
      expect(memoryHistory.location.pathname).toEqual(editAboutYourPPM);
    });
  });

  it('routes to the add weight ticket page when the add link is clicked', async () => {
    const newWeightTicket = generatePath(customerRoutes.SHIPMENT_PPM_WEIGHT_TICKETS_PATH, {
      moveId: mockMoveId,
      mtoShipmentId: mockMTOShipmentId,
    });
    const memoryHistory = createMemoryHistory();
    const mockProviderWithHistory = ({ children }) => <MockProviders history={memoryHistory}>{children}</MockProviders>;
    render(<Review />, { wrapper: mockProviderWithHistory });

    userEvent.click(screen.getByText('Add More Weight'));

    await waitFor(() => {
      expect(memoryHistory.location.pathname).toEqual(newWeightTicket);
    });
  });

  it('routes to the edit weight ticket page when the edit link is clicked', async () => {
    selectMTOShipmentById.mockImplementationOnce(() => mockMTOShipmentWithWeightTicket);
    const editWeightTicket = generatePath(customerRoutes.SHIPMENT_PPM_WEIGHT_TICKETS_EDIT_PATH, {
      moveId: mockMoveId,
      mtoShipmentId: mockMTOShipmentId,
      weightTicketId: mockWeightTicketId,
    });
    const memoryHistory = createMemoryHistory();
    const mockProviderWithHistory = ({ children }) => <MockProviders history={memoryHistory}>{children}</MockProviders>;
    render(<Review />, { wrapper: mockProviderWithHistory });

    userEvent.click(screen.getAllByText('Edit')[1]);

    await waitFor(() => {
      expect(memoryHistory.location.pathname).toEqual(editWeightTicket);
    });
  });

  it('routes to the home page when the finish later link is clicked', async () => {
    const memoryHistory = createMemoryHistory();
    const mockProviderWithHistory = ({ children }) => <MockProviders history={memoryHistory}>{children}</MockProviders>;
    render(<Review />, { wrapper: mockProviderWithHistory });

    userEvent.click(screen.getByText('Finish Later'));

    await waitFor(() => {
      expect(memoryHistory.location.pathname).toEqual('/');
    });
  });

  it('routes to the complete page when the save and continue link is clicked', async () => {
    selectMTOShipmentById.mockImplementationOnce(() => mockMTOShipmentWithWeightTicket);
    const completePath = generatePath(customerRoutes.SHIPMENT_PPM_COMPLETE_PATH, {
      moveId: mockMoveId,
      mtoShipmentId: mockMTOShipmentId,
    });

    const memoryHistory = createMemoryHistory();
    const mockProviderWithHistory = ({ children }) => <MockProviders history={memoryHistory}>{children}</MockProviders>;
    render(<Review />, { wrapper: mockProviderWithHistory });

    userEvent.click(screen.getByText('Save & Continue'));

    await waitFor(() => {
      expect(memoryHistory.location.pathname).toEqual(completePath);
    });
  });

  it('disables the save and continue link when there are no weight tickets', async () => {
    const memoryHistory = createMemoryHistory();
    const mockProviderWithHistory = ({ children }) => <MockProviders history={memoryHistory}>{children}</MockProviders>;
    render(<Review />, { wrapper: mockProviderWithHistory });

    expect(screen.getByText('Save & Continue')).toHaveClass('usa-button--disabled');
    expect(screen.getByText('Save & Continue')).toHaveAttribute('aria-disabled', 'true');
  });

  it('disables the save and continue link when there is an incomplete weight ticket', async () => {
    selectMTOShipmentById.mockImplementationOnce(() => mockMTOShipmentWithIncompleteWeightTicket);
    const memoryHistory = createMemoryHistory();
    const mockProviderWithHistory = ({ children }) => <MockProviders history={memoryHistory}>{children}</MockProviders>;
    render(<Review />, { wrapper: mockProviderWithHistory });

    expect(screen.getByText('Save & Continue')).toHaveClass('usa-button--disabled');
    expect(screen.getByText('Save & Continue')).toHaveAttribute('aria-disabled', 'true');
  });

  it('displays the delete confirmation modal when the delete button is clicked', async () => {
    selectMTOShipmentById.mockImplementationOnce(() => mockMTOShipmentWithWeightTicket);
    render(<Review />, { wrapper: MockProviders });

    userEvent.click(screen.getAllByRole('button', { name: 'Delete' })[0]);

    await waitFor(() => {
      expect(screen.getByRole('heading', { level: 3, name: 'Delete this?' })).toBeInTheDocument();
    });

    userEvent.click(screen.getByRole('button', { name: 'No, Keep It' }));

    expect(screen.queryByRole('heading', { level: 3, name: 'Delete this?' })).not.toBeInTheDocument();
  });

  it('calls the delete weight ticket api when confirm is clicked', async () => {
    selectMTOShipmentById.mockImplementationOnce(() => mockMTOShipmentWithWeightTicket);
    const mockDeleteWeightTicket = jest.fn().mockResolvedValue({});
    deleteWeightTicket.mockImplementationOnce(mockDeleteWeightTicket);
    render(<Review />, { wrapper: MockProviders });

    userEvent.click(screen.getAllByRole('button', { name: 'Delete' })[0]);

    await waitFor(() => {
      expect(screen.getByRole('heading', { level: 3, name: 'Delete this?' })).toBeInTheDocument();
    });

    userEvent.click(screen.getByRole('button', { name: 'Yes, Delete' }));

    await waitFor(() => {
      expect(mockDeleteWeightTicket).toHaveBeenCalledWith(mockWeightTicketId, 'eTag value');
    });
  });
});

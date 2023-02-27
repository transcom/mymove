import React from 'react';
import { screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import ConnectedReview from 'pages/MyMove/Review/Review';
import { renderWithProviders } from 'testUtils';
import MOVE_STATUSES from 'constants/moves';
import { selectCurrentMove, selectMTOShipmentsForCurrentMove } from 'store/entities/selectors';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { customerRoutes } from 'constants/routes';

// Mock the summary part of the review page since we're just testing the
// navigation portion.
jest.mock('components/Customer/Review/Summary/Summary', () => 'summary');

// Explicitly setup navigate mock so we can verify it was called with correct pathing in tests
const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
}));

jest.mock('store/entities/selectors', () => ({
  ...jest.requireActual('store/entities/selectors'),
  selectCurrentMove: jest.fn(),
  selectMTOShipmentsForCurrentMove: jest.fn(),
}));

afterEach(jest.resetAllMocks);

describe('Review page', () => {
  const testDraftMove = {
    status: MOVE_STATUSES.DRAFT,
  };

  const testSubmittedMove = {
    status: MOVE_STATUSES.SUBMITTED,
  };

  const testIncompletePPMShipment = {
    id: '1',
    shipmentType: SHIPMENT_OPTIONS.PPM,
    ppmShipment: {
      id: '2',
      expectedDepartureDate: '2022-04-01',
      pickupPostalCode: '90210',
      destinationPostalCode: '10001',
      sitExpected: false,
    },
  };

  const testCompletePPMShipment = {
    id: '1',
    shipmentType: SHIPMENT_OPTIONS.PPM,
    ppmShipment: {
      id: '2',
      expectedDepartureDate: '2022-04-01',
      pickupPostalCode: '90210',
      destinationPostalCode: '10001',
      sitExpected: false,
      estimatedWeight: 5999,
      hasProGear: false,
      hasRequestedAdvance: false,
    },
  };

  const testHHGShipment = {
    id: '3',
    shipmentType: SHIPMENT_OPTIONS.HHG,
    requestedPickupDate: '2022-04-01',
    pickupLocation: {
      streetAddress1: '17 8th St',
      city: 'New York',
      state: 'NY',
      postalCode: '11111',
    },
    requestedDeliveryDate: '2022-05-01',
    destinationZIP: '73523',
  };

  const mockParams = { moveId: '3a8c9f4f-7344-4f18-9ab5-0de3ef57b901' };
  const mockPath = customerRoutes.MOVE_REVIEW_PATH;
  const mockRoutingOptions = { path: mockPath, params: mockParams };

  const testFlashState = {
    flash: {
      flashMessage: {
        type: 'SET_FLASH_MESSAGE',
        title: 'Details saved',
        messageType: 'success',
        message: 'Review your info and submit your move request now, or come back and finish later.',
        key: 'PPM_ONBOARDING_SUBMIT_SUCCESS',
        slim: false,
      },
    },
  };

  it('renders the Review Page', async () => {
    selectCurrentMove.mockImplementation(() => testDraftMove);
    renderWithProviders(<ConnectedReview />, mockRoutingOptions);

    await screen.findByRole('heading', { level: 1, name: 'Review your details' });
  });

  it('Finish Later button goes back to the home page', async () => {
    selectCurrentMove.mockImplementation(() => testDraftMove);

    renderWithProviders(<ConnectedReview />, mockRoutingOptions);

    const backButton = screen.getByRole('button', { name: 'Finish later' });

    expect(backButton).toBeInTheDocument();

    await userEvent.click(backButton);

    screen.debug();

    expect(mockNavigate).toHaveBeenCalledWith('/');
  });

  it('next button goes to the Agreement page when move is in DRAFT status', async () => {
    selectCurrentMove.mockImplementation(() => testDraftMove);
    selectMTOShipmentsForCurrentMove.mockImplementation(() => [testCompletePPMShipment, testHHGShipment]);

    renderWithProviders(<ConnectedReview />, mockRoutingOptions);

    const submitButton = await screen.findByRole('button', { name: 'Next' });

    expect(submitButton).toBeInTheDocument();

    await userEvent.click(submitButton);

    expect(mockNavigate).toHaveBeenCalledWith(`/moves/${mockParams.moveId}/agreement`);
  });

  it('next button goes to the Agreement page when move is in DRAFT status with only HHG shipment', async () => {
    selectCurrentMove.mockImplementation(() => testDraftMove);
    selectMTOShipmentsForCurrentMove.mockImplementation(() => [testHHGShipment]);

    renderWithProviders(<ConnectedReview />, mockRoutingOptions);

    const submitButton = await screen.findByRole('button', { name: 'Next' });

    expect(submitButton).toBeInTheDocument();

    await userEvent.click(submitButton);

    expect(mockNavigate).toHaveBeenCalledWith(`/moves/${mockParams.moveId}/agreement`);
  });

  it('next button is disabled when a PPM shipment is in an incomplete state', async () => {
    selectCurrentMove.mockImplementation(() => testDraftMove);
    selectMTOShipmentsForCurrentMove.mockImplementation(() => [testIncompletePPMShipment, testHHGShipment]);

    renderWithProviders(<ConnectedReview />, mockRoutingOptions);

    const submitButton = await screen.findByRole('button', { name: 'Next' });

    expect(submitButton).toBeDisabled();
  });

  it('next button is disabled when a there are no shipments', async () => {
    selectCurrentMove.mockImplementation(() => testDraftMove);
    selectMTOShipmentsForCurrentMove.mockImplementation(() => []);

    renderWithProviders(<ConnectedReview />, mockRoutingOptions);

    const submitButton = await screen.findByRole('button', { name: 'Next' });

    expect(submitButton).toBeDisabled();
  });

  it('return home button is displayed when move has been submitted', async () => {
    selectCurrentMove.mockImplementation(() => testSubmittedMove);

    renderWithProviders(<ConnectedReview />, mockRoutingOptions);

    expect(screen.queryByRole('button', { name: 'Next' })).not.toBeInTheDocument();

    expect(screen.getByRole('button', { name: 'Return home' })).toBeInTheDocument();
  });

  it('renders the success alert flash message', async () => {
    renderWithProviders(<ConnectedReview />, { ...mockRoutingOptions, initialState: testFlashState });

    expect(screen.getByRole('heading', { level: 4 })).toHaveTextContent('Details saved');
    expect(
      screen.getByText('Review your info and submit your move request now, or come back and finish later.'),
    ).toBeInTheDocument();
  });
});

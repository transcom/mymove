import React from 'react';
import { render, waitFor, screen, within, act } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { generatePath } from 'react-router-dom';
import { v4 } from 'uuid';
import moment from 'moment';

import About from 'pages/Office/PPM/Closeout/About/About';
import { servicesCounselingRoutes } from 'constants/routes';
import { updateMTOShipment } from 'services/ghcApi';
import { MockProviders } from 'testUtils';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { shipmentStatuses } from 'constants/shipments';
import { usePPMShipmentAndDocsOnlyQueries } from 'hooks/queries';
import { isBooleanFlagEnabled } from 'utils/featureFlags';

const mockMoveId = v4();
const mockMTOShipmentId = v4();
const mockPPMShipmentId = v4();

const mockRoutingConfig = {
  path: servicesCounselingRoutes.BASE_SHIPMENT_PPM_ABOUT_PATH,
  params: {
    moveCode: mockMoveId,
    shipmentId: mockMTOShipmentId,
  },
};

jest.mock('utils/featureFlags', () => ({
  ...jest.requireActual('utils/featureFlags'),
  isBooleanFlagEnabled: jest.fn().mockImplementation(() => Promise.resolve()),
}));

const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
}));

jest.mock('services/ghcApi', () => ({
  ...jest.requireActual('services/ghcApi'),
  updateMTOShipment: jest.fn(),
}));

jest.mock('hooks/queries', () => ({
  usePPMShipmentAndDocsOnlyQueries: jest.fn(),
}));

const mtoShipmentCreatedDate = new Date();
const ppmShipmentCreatedDate = moment(mtoShipmentCreatedDate).add(5, 'seconds');
const approvedDate = moment(ppmShipmentCreatedDate).add(2, 'days');

const mockMTOShipment = {
  id: mockMTOShipmentId,
  shipmentType: SHIPMENT_OPTIONS.PPM,
  status: shipmentStatuses.APPROVED,
  moveTaskOrderId: mockMoveId,
  ppmShipment: {
    id: mockPPMShipmentId,
    shipmentId: mockMTOShipmentId,
    weightTickets: [],
    pickupAddress: {
      streetAddress1: '812 S 129th St',
      streetAddress2: '#123',
      streetAddress3: '',
      city: 'San Antonio',
      state: 'TX',
      postalCode: '78234',
      usPostRegionCitiesID: '',
      country: {
        code: 'US',
        name: 'UNITED STATES',
        id: '791899e6-cd77-46f2-981b-176ecb8d7098',
      },
      countryID: '791899e6-cd77-46f2-981b-176ecb8d7098',
    },
    destinationAddress: {
      streetAddress1: '441 SW Rio de la Plata Drive',
      streetAddress2: '#124',
      streetAddress3: 'Some Person',
      city: 'Tacoma',
      state: 'WA',
      postalCode: '98421',
      usPostRegionCitiesID: '',
      country: {
        code: 'US',
        name: 'UNITED STATES',
        id: '791899e6-cd77-46f2-981b-176ecb8d7098',
      },
      countryID: '791899e6-cd77-46f2-981b-176ecb8d7098',
    },
    w2Address: {
      streetAddress1: '11 NE Elm Road',
      streetAddress2: '',
      streetAddress3: '',
      city: 'Jacksonville',
      state: 'FL',
      postalCode: '32217',
      usPostRegionCitiesID: '',
      country: {
        code: 'US',
        name: 'UNITED STATES',
        id: '791899e6-cd77-46f2-981b-176ecb8d7098',
      },
      countryID: '791899e6-cd77-46f2-981b-176ecb8d7098',
    },
  },
  createdAt: mtoShipmentCreatedDate.toISOString(),
  updatedAt: approvedDate.toISOString(),
  eTag: window.btoa(approvedDate.toISOString()),
};

const partialPayload = {
  actualMoveDate: '2022-05-31',
  pickupAddress: {
    streetAddress1: '812 S 129th St',
    streetAddress2: '#123',
    streetAddress3: 'Some Person',
    city: 'San Antonio',
    state: 'TX',
    postalCode: '78234',
    usPostRegionCitiesID: '',
    country: {
      code: 'US',
      name: 'UNITED STATES',
      id: '791899e6-cd77-46f2-981b-176ecb8d7098',
    },
    countryID: '791899e6-cd77-46f2-981b-176ecb8d7098',
  },
  destinationAddress: {
    streetAddress1: '441 SW Rio de la Plata Drive',
    streetAddress2: '#124',
    streetAddress3: 'Some Person',
    city: 'Tacoma',
    state: 'WA',
    postalCode: '98421',
    usPostRegionCitiesID: '',
    country: {
      code: 'US',
      name: 'UNITED STATES',
      id: '791899e6-cd77-46f2-981b-176ecb8d7098',
    },
    countryID: '791899e6-cd77-46f2-981b-176ecb8d7098',
  },
  secondaryPickupAddress: null,
  secondaryDestinationAddress: null,
  hasSecondaryPickupAddress: false,
  hasSecondaryDestinationAddress: false,
  hasReceivedAdvance: true,
  advanceAmountReceived: 750000,
  w2Address: {
    streetAddress1: '11 NE Elm Road',
    streetAddress2: '',
    streetAddress3: '',
    city: 'Jacksonville',
    state: 'FL',
    postalCode: '32217',
    usPostRegionCitiesID: '',
    country: {
      code: 'US',
      name: 'UNITED STATES',
      id: '791899e6-cd77-46f2-981b-176ecb8d7098',
    },
    countryID: '791899e6-cd77-46f2-981b-176ecb8d7098',
  },
};

const mockMTOShipmentResponse = {
  ...mockMTOShipment,
  ppmShipment: {
    ...mockMTOShipment.ppmShipment,
    ...partialPayload,
  },
};

beforeEach(() => {
  jest.clearAllMocks();
});

// const movePath = generatePath(-1);
const weightTicketsPath = generatePath(servicesCounselingRoutes.BASE_SHIPMENT_PPM_REVIEW_PATH, {
  moveCode: mockMoveId,
  shipmentId: mockMTOShipmentId,
});

const fillOutAdvanceSections = async (form) => {
  await act(async () => {
    await userEvent.click(within(form).getAllByLabelText('Yes')[2]);
  });

  within(form).getByLabelText('How much did you receive?').focus();
  await act(async () => {
    await userEvent.paste('7500');
  });
};

const renderAboutPage = () => {
  return render(
    <MockProviders {...mockRoutingConfig}>
      <About />
    </MockProviders>,
  );
};

describe('About page', () => {
  it('renders the page Content', async () => {
    usePPMShipmentAndDocsOnlyQueries.mockReturnValue({
      isLoading: false,
      mtoShipment: mockMTOShipmentResponse,
      error: null,
    });
    renderAboutPage();

    await waitFor(() => {
      expect(screen.getByTestId('tag')).toHaveTextContent('PPM');
    });

    expect(screen.getByRole('heading', { level: 1 })).toHaveTextContent('About your PPM');
  });

  it('routes back to home when return to homepage is clicked', async () => {
    usePPMShipmentAndDocsOnlyQueries.mockReturnValue({
      isLoading: false,
      mtoShipment: mockMTOShipmentResponse,
      error: null,
    });
    renderAboutPage();

    await act(async () => {
      await userEvent.click(screen.getByRole('button', { name: 'Cancel' }));
    });

    expect(mockNavigate).toHaveBeenCalledWith(-1);
  });

  it('calls the patch shipment with the appropriate payload', async () => {
    usePPMShipmentAndDocsOnlyQueries.mockReturnValue({
      isLoading: false,
      mtoShipment: mockMTOShipmentResponse,
      error: null,
    });
    renderAboutPage();

    updateMTOShipment.mockImplementation(() => Promise.resolve({}));

    let form;
    await waitFor(() => {
      form = screen.getByTestId('aboutForm');
    });

    await fillOutAdvanceSections(form);

    await act(async () => {
      await userEvent.click(within(form).getByRole('button', { name: 'Save & Continue' }));
    });

    await waitFor(() => {
      expect(mockNavigate).toHaveBeenCalledWith(weightTicketsPath);
    });
  });

  it('displays an error when the patch shipment API fails', async () => {
    isBooleanFlagEnabled.mockResolvedValue(true);
    const mockErrorMsg = {
      response: {
        body: {
          message: 'Error Updating',
        },
      },
    };
    usePPMShipmentAndDocsOnlyQueries.mockReturnValue({
      isLoading: false,
      mtoShipment: mockMTOShipmentResponse,
      error: null,
    });

    renderAboutPage();

    updateMTOShipment.mockImplementation(() => Promise.reject(mockErrorMsg));
    let form;
    await waitFor(() => {
      form = screen.getByTestId('aboutForm');
    });

    expect(within(form).getByRole('button', { name: 'Save & Continue' })).toBeEnabled();
    await act(async () => {
      await userEvent.click(within(form).getByRole('button', { name: 'Save & Continue' }));
    });

    await waitFor(() => {
      expect(updateMTOShipment).toHaveBeenCalled();
    });

    await waitFor(() => {
      expect(screen.getByText(mockErrorMsg.response.body.message)).toBeInTheDocument();
    });
  });

  it('expect loadingPlaceholder when mtoShipment is falsy', async () => {
    usePPMShipmentAndDocsOnlyQueries.mockReturnValue({
      isLoading: false,
      mtoShipment: null,
      error: null,
    });

    renderAboutPage();
    expect(screen.getByRole('heading', { level: 2 })).toHaveTextContent('Loading, please wait...');
  });
});

import React from 'react';
import { render, waitFor, screen, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { generatePath } from 'react-router-dom';
import { v4 } from 'uuid';
import moment from 'moment';

import About from 'pages/MyMove/PPM/Closeout/About/About';
import { selectMTOShipmentById } from 'store/entities/selectors';
import { customerRoutes } from 'constants/routes';
import { getResponseError, patchMTOShipment, getMTOShipmentsForMove } from 'services/internalApi';
import { updateMTOShipment } from 'store/entities/actions';
import { MockProviders } from 'testUtils';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { shipmentStatuses } from 'constants/shipments';
import { shipment } from 'shared/Entities/schema';

const mockMoveId = v4();
const mockMTOShipmentId = v4();
const mockPPMShipmentId = v4();

const mockRoutingConfig = {
  path: customerRoutes.SHIPMENT_PPM_ABOUT_PATH,
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

jest.mock('services/internalApi', () => ({
  ...jest.requireActual('services/internalApi'),
  patchMTOShipment: jest.fn(),
  getResponseError: jest.fn(),
  getMTOShipmentsForMove: jest.fn(),
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
  },
  createdAt: mtoShipmentCreatedDate.toISOString(),
  updatedAt: approvedDate.toISOString(),
  eTag: window.btoa(approvedDate.toISOString()),
};

const partialPayload = {
  actualMoveDate: '2022-05-31',
  actualPickupPostalCode: '78234',
  actualDestinationPostalCode: '98421',
  pickupAddress: {
    streetAddress1: '812 S 129th St',
    streetAddress2: '#123',
    streetAddress3: 'Some Person',
    city: 'San Antonio',
    state: 'TX',
    postalCode: '78234',
  },
  destinationAddress: {
    streetAddress1: '441 SW Rio de la Plata Drive',
    streetAddress2: '#124',
    streetAddress3: 'Some Person',
    city: 'Tacoma',
    state: 'WA',
    postalCode: '98421',
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
  },
};

const mockPayload = {
  shipmentType: SHIPMENT_OPTIONS.PPM,
  ppmShipment: {
    id: mockPPMShipmentId,
    ...partialPayload,
  },
};

const mockMTOShipmentResponse = {
  ...mockMTOShipment,
  ppmShipment: {
    ...mockMTOShipment.ppmShipment,
    ...partialPayload,
  },
};

jest.mock('store/entities/selectors', () => ({
  ...jest.requireActual('store/entities/selectors'),
  selectMTOShipmentById: jest.fn(() => mockMTOShipment),
}));

const mockDispatch = jest.fn();
jest.mock('react-redux', () => ({
  ...jest.requireActual('react-redux'),
  useDispatch: jest.fn(() => mockDispatch),
}));

beforeEach(() => {
  jest.clearAllMocks();
});

const movePath = generatePath(customerRoutes.MOVE_HOME_PAGE);
const weightTicketsPath = generatePath(customerRoutes.SHIPMENT_PPM_WEIGHT_TICKETS_PATH, {
  moveId: mockMoveId,
  mtoShipmentId: mockMTOShipmentId,
});

const fillOutBasicForm = async (form) => {
  within(form).getByLabelText('When did you leave your origin?').focus();
  await userEvent.paste('2022-05-31');

  within(form)
    .getAllByLabelText(/Address 1/)[0]
    .focus();
  await userEvent.paste('812 S 129th St');

  within(form)
    .getAllByLabelText(/Address 2/)[0]
    .focus();
  await userEvent.paste('#123');

  within(form)
    .getAllByLabelText(/Address 3/)[0]
    .focus();
  await userEvent.paste('Some Person');

  within(form).getAllByLabelText(/City/)[0].focus();
  await userEvent.paste('San Antonio');

  within(form).getAllByLabelText(/State/)[0].focus();
  await userEvent.selectOptions(within(form).getAllByLabelText(/State/)[0], 'TX');

  within(form).getAllByLabelText(/ZIP/)[0].focus();
  await userEvent.paste('78234');

  within(form)
    .getAllByLabelText(/Address 1/)[1]
    .focus();
  await userEvent.paste('441 SW Rio de la Plata Drive');

  within(form)
    .getAllByLabelText(/Address 2/)[1]
    .focus();
  await userEvent.paste('#124');

  within(form)
    .getAllByLabelText(/Address 3/)[1]
    .focus();
  await userEvent.paste('Some Person');

  within(form).getAllByLabelText(/City/)[1].focus();
  await userEvent.paste('Tacoma');

  within(form).getAllByLabelText(/State/)[1].focus();
  await userEvent.selectOptions(within(form).getAllByLabelText(/State/)[1], 'WA');

  within(form).getAllByLabelText(/ZIP/)[1].focus();
  await userEvent.paste('98421');

  within(form)
    .getAllByLabelText(/Address 1/)[2]
    .focus();
  await userEvent.paste('11 NE Elm Road');

  within(form).getAllByLabelText(/City/)[2].focus();
  await userEvent.paste('Jacksonville');

  within(form).getAllByLabelText(/State/)[2].focus();
  await userEvent.selectOptions(within(form).getAllByLabelText(/State/)[2], 'FL');

  within(form).getAllByLabelText(/ZIP/)[2].focus();
  await userEvent.paste('32217');
};

const fillOutAdvanceSections = async (form) => {
  await userEvent.click(within(form).getAllByLabelText('Yes')[2]);

  within(form).getByLabelText('How much did you receive?').focus();
  await userEvent.paste('7500');
};

const renderAboutPage = () => {
  return render(
    <MockProviders {...mockRoutingConfig}>
      <About />
    </MockProviders>,
  );
};

describe('About page', () => {
  it('loads the selected shipment from redux', () => {
    getMTOShipmentsForMove.mockResolvedValueOnce(shipment);
    renderAboutPage();

    expect(selectMTOShipmentById).toHaveBeenCalledWith(expect.anything(), mockMTOShipmentId);
  });

  it('renders the page Content', async () => {
    await getMTOShipmentsForMove.mockResolvedValueOnce(shipment);
    renderAboutPage();

    await waitFor(() => {
      expect(screen.getByTestId('tag')).toHaveTextContent('PPM');
    });

    expect(screen.getByRole('heading', { level: 1 })).toHaveTextContent('About your PPM');
    expect(screen.getByText('Finish moving this PPM before you start documenting it.')).toBeInTheDocument();
    const headings = screen.getAllByRole('heading', { level: 2 });
    expect(headings[0]).toHaveTextContent('How to complete your PPM');
    expect(headings[1]).toHaveTextContent('About your final payment');

    // renders form content
    expect(headings[2]).toHaveTextContent('Departure date');
    expect(headings[3]).toHaveTextContent('Locations');
    expect(headings[4]).toHaveTextContent('Advance (AOA)');
    expect(headings[5]).toHaveTextContent('W-2 address');
  });

  it('routes back to home when return to homepage is clicked', async () => {
    await getMTOShipmentsForMove.mockResolvedValueOnce(shipment);
    renderAboutPage();

    await waitFor(async () => {
      await userEvent.click(screen.getByRole('button', { name: 'Return To Homepage' }));
    });

    expect(mockNavigate).toHaveBeenCalledWith(movePath);
  });

  it('calls the patch shipment with the appropriate payload', async () => {
    patchMTOShipment.mockResolvedValueOnce(mockMTOShipmentResponse);
    getMTOShipmentsForMove.mockResolvedValueOnce(shipment);

    renderAboutPage();

    let form;
    await waitFor(() => {
      form = screen.getByTestId('aboutForm');
    });

    await fillOutBasicForm(form);
    await fillOutAdvanceSections(form);

    await userEvent.click(within(form).getByRole('button', { name: 'Save & Continue' }));
    await waitFor(() => {
      expect(patchMTOShipment).toHaveBeenCalledWith(mockMTOShipmentId, mockPayload, mockMTOShipment.eTag);
    });

    expect(mockDispatch).toHaveBeenCalledWith(updateMTOShipment(mockMTOShipmentResponse));
    expect(mockNavigate).toHaveBeenCalledWith(weightTicketsPath);
  });

  it('displays an error when the patch shipment API fails', async () => {
    const mockErrorMsg = 'Error Updating';
    await getMTOShipmentsForMove.mockResolvedValueOnce(shipment);
    await patchMTOShipment.mockRejectedValue(mockErrorMsg);
    await getResponseError.mockReturnValue(mockErrorMsg);

    renderAboutPage();

    let form;
    await waitFor(() => {
      form = screen.getByTestId('aboutForm');
    });

    await fillOutBasicForm(form);

    expect(within(form).getByRole('button', { name: 'Save & Continue' })).toBeEnabled();
    await userEvent.click(within(form).getByRole('button', { name: 'Save & Continue' }));
    const payload = {
      ...mockPayload,
      ppmShipment: {
        ...mockPayload.ppmShipment,
        hasReceivedAdvance: false,
        advanceAmountReceived: null,
      },
    };
    await waitFor(() => {
      expect(patchMTOShipment).toHaveBeenCalledWith(mockMTOShipmentId, payload, mockMTOShipment.eTag);
    });

    await waitFor(() => {
      expect(screen.getByText(mockErrorMsg)).toBeInTheDocument();
    });
  });

  it('expect loadingPlaceholder when mtoShipment is falsy', () => {
    selectMTOShipmentById.mockReturnValue(null);
    getMTOShipmentsForMove.mockResolvedValueOnce(shipment);

    renderAboutPage();
    expect(screen.getByRole('heading', { level: 2 })).toHaveTextContent('Loading, please wait...');
  });
});

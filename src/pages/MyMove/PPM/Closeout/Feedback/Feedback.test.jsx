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
        adjustedNetWeight: null,
        createdAt: '2024-10-30T16:09:54.526Z',
        eTag: 'MjAyNC0xMC0zMFQxNjoxNTowNi40NDA0NzZa',
        emptyDocument: {
          id: '3441bab5-43c6-48be-9819-f6e8ac6abbef',
          service_member_id: '638ae45c-cac0-42f5-9058-b8847e55ac29',
          uploads: [
            {
              bytes: 72540,
              contentType: 'image/png',
              createdAt: '2024-10-30T16:10:11.693Z',
              filename: 'Screenshot 2024-10-29 at 10.16.08 AM.png-20241030101011',
              id: '47298bc3-9ab8-45a7-8430-8fc4547a36b1',
              status: 'PROCESSING',
              updatedAt: '2024-10-30T16:10:11.693Z',
              uploadType: 'USER',
              url: '/storage/user/21142e5f-f599-4658-bd5b-82c0f5392ac8/uploads/47298bc3-9ab8-45a7-8430-8fc4547a36b1?contentType=image%2Fpng\u0026filename=Screenshot+2024-10-29+at+10.16.08%E2%80%AFAM.png-20241030101011',
            },
          ],
        },
        emptyDocumentId: '3441bab5-43c6-48be-9819-f6e8ac6abbef',
        emptyWeight: 1999,
        fullDocument: {
          id: '7528d4f7-35d7-4f5c-9ceb-1c6d477763e7',
          service_member_id: '638ae45c-cac0-42f5-9058-b8847e55ac29',
          uploads: [
            {
              bytes: 324015,
              contentType: 'image/png',
              createdAt: '2024-10-30T16:10:18.163Z',
              filename: 'Screenshot 2024-09-26 at 12.27.12 PM.png-20241030101018',
              id: '543996cf-1585-4790-9437-55e74bbd130b',
              status: 'PROCESSING',
              updatedAt: '2024-10-30T16:10:18.163Z',
              uploadType: 'USER',
              url: '/storage/user/21142e5f-f599-4658-bd5b-82c0f5392ac8/uploads/543996cf-1585-4790-9437-55e74bbd130b?contentType=image%2Fpng\u0026filename=Screenshot+2024-09-26+at+12.27.12%E2%80%AFPM.png-20241030101018',
            },
          ],
        },
        fullDocumentId: '7528d4f7-35d7-4f5c-9ceb-1c6d477763e7',
        fullWeight: 5844,
        id: 'a9d11921-5c60-4e41-a45c-8a978d1bf670',
        missingEmptyWeightTicket: false,
        missingFullWeightTicket: false,
        netWeightRemarks: null,
        ownsTrailer: false,
        ppmShipmentId: 'bca99262-41f7-49e8-b24a-12a4b3958b9e',
        proofOfTrailerOwnershipDocument: {
          id: 'e4b31f17-4804-4d71-b133-349c759bbaa6',
          service_member_id: '638ae45c-cac0-42f5-9058-b8847e55ac29',
          uploads: [],
        },
        proofOfTrailerOwnershipDocumentId: 'e4b31f17-4804-4d71-b133-349c759bbaa6',
        reason: 'asdf',
        status: 'REJECTED',
        submittedEmptyWeight: 1999,
        submittedFullWeight: 5844,
        submittedOwnsTrailer: false,
        submittedTrailerMeetsCriteria: false,
        trailerMeetsCriteria: false,
        updatedAt: '2024-10-30T16:15:06.440Z',
        vehicleDescription: '2023 F-150',
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

  it('formats and diplays trip weight', () => {
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
});

import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { v4 as uuidv4 } from 'uuid';
import { generatePath } from 'react-router';

import FinalCloseout from 'pages/Office/PPM/Closeout/FinalCloseout/FinalCloseout';
import { updateMTOShipment } from 'store/entities/actions';
import { MockProviders } from 'testUtils';
import { ppmSubmissionCertificationText } from 'scenes/Legalese/legaleseText';
import { formatDateForSwagger } from 'shared/dates';
import { servicesCounselingRoutes } from 'constants/routes';
import { getMove, getMTOShipments, submitPPMShipmentSignedCertification } from 'services/ghcApi';

const testMove = {
  additionalDocuments: {
    id: 'c43ae36e-4e15-4cb3-865a-e4dccffa0df7',
    service_member_id: 'dfdd3e21-3988-4104-a5c2-06b195f9b7f0',
    uploads: [
      {
        bytes: 120653,
        contentType: 'application/pdf',
        createdAt: '2024-05-29T19:14:39.108Z',
        filename: '9380-Statement-20240430.pdf',
        id: 'c3c0cda9-a77e-4b8b-8b8b-67ccadc3c862',
        status: 'PROCESSING',
        updatedAt: '2024-05-29T19:14:39.108Z',
        url: '/storage/user/accf760b-2e3d-4af8-a59b-c10b591dcc15/uploads/c3c0cda9-a77e-4b8b-8b8b-67ccadc3c862?contentType=application%2Fpdf',
      },
      {
        bytes: 307051,
        contentType: 'image/png',
        createdAt: '2024-05-30T04:23:27.241Z',
        filename: 'Screenshot 2024-05-16 at 3.33.52 PM.png',
        id: '70a35ab0-a3f5-44a3-8702-0bb7d0c568c8',
        status: 'PROCESSING',
        updatedAt: '2024-05-30T04:23:27.241Z',
        url: '/storage/user/accf760b-2e3d-4af8-a59b-c10b591dcc15/uploads/70a35ab0-a3f5-44a3-8702-0bb7d0c568c8?contentType=image%2Fpng',
      },
      {
        bytes: 82301,
        contentType: 'image/png',
        createdAt: '2024-05-30T04:33:10.622Z',
        filename: 'Screenshot 2024-05-17 at 1.09.21 PM.png',
        id: 'b11c0130-2403-4287-b464-4c5ac17797b3',
        status: 'PROCESSING',
        updatedAt: '2024-05-30T04:33:10.622Z',
        url: '/storage/user/accf760b-2e3d-4af8-a59b-c10b591dcc15/uploads/b11c0130-2403-4287-b464-4c5ac17797b3?contentType=image%2Fpng',
      },
    ],
  },
  created_at: '2024-05-29T18:46:17.808Z',
  eTag: 'MjAyNC0wNS0yOVQxOToxNDozOS4xMDQyNzJa',
  id: '43a369e8-5fa3-4a13-9d9a-36d86731c1da',
  locator: '988HDJ',
  mto_shipments: ['c93bf4d1-1470-4c50-b2b6-f736abd2986a'],
  orders_id: '69967de3-3d9d-4e73-a497-f401884393bf',
  primeCounselingCompletedAt: '0001-01-01T00:00:00.000Z',
  service_member_id: 'dfdd3e21-3988-4104-a5c2-06b195f9b7f0',
  status: 'NEEDS SERVICE COUNSELING',
  submitted_at: '2024-05-29T18:47:26.360Z',
  updated_at: '2024-05-29T19:14:39.104Z',
};

const shipmentID = uuidv4();
const response = {
  mtoShipments: {
    [shipmentID]: {
      actualProGearWeight: null,
      actualSpouseProGearWeight: null,
      createdAt: '2025-03-25T14:33:35.101Z',
      destinationSitAuthEndDate: '0001-01-01T00:00:00.000Z',
      distance: 18,
      eTag: 'MjAyNS0wMy0yNVQxNTo1NjowNS42NTg5NDla',
      hasSecondaryDeliveryAddress: false,
      hasSecondaryPickupAddress: false,
      hasTertiaryDeliveryAddress: false,
      hasTertiaryPickupAddress: false,
      id: shipmentID,
      marketCode: 'd',
      moveTaskOrderID: '9cd0fda6-7a5b-462a-985b-37cbdef70e68',
      originSitAuthEndDate: '0001-01-01T00:00:00.000Z',
      ppmShipment: {
        actualDestinationPostalCode: '90247',
        actualMoveDate: '2025-03-17',
        actualPickupPostalCode: '90210',
        advanceAmountReceived: null,
        advanceAmountRequested: null,
        allowableWeight: null,
        approvedAt: '2025-03-25T15:54:35.870Z',
        createdAt: '2025-03-25T14:33:35.157Z',
        destinationAddress: {
          city: 'GARDENA',
          county: 'LOS ANGELES',
          eTag: 'MjAyNS0wMy0yNVQxNTo1NjowNC41ODY0MDla',
          id: 'a2248b43-1b89-4492-902b-d55b0688c502',
          isOconus: false,
          postalCode: '90247',
          state: 'CA',
          streetAddress1: '123 Any Street',
          streetAddress2: 'P.O. Box 12345',
          streetAddress3: 'c/o Some Person',
          usPostRegionCitiesID: '4e704eb8-741a-4e71-b247-217834247025',
        },
        eTag: 'MjAyNS0wMy0yNVQxNTo1NjowNS42NTUwNzFa',
        estimatedIncentive: 1714921,
        estimatedWeight: 1600,
        expectedDepartureDate: '2025-03-28',
        finalIncentive: 115225,
        hasProGear: false,
        hasReceivedAdvance: false,
        hasRequestedAdvance: false,
        hasSecondaryDestinationAddress: false,
        hasSecondaryPickupAddress: false,
        hasTertiaryDestinationAddress: false,
        hasTertiaryPickupAddress: false,
        id: 'b459e544-3654-430d-92d5-296dfca4009e',
        isActualExpenseReimbursement: false,
        maxIncentive: 36935085,
        movingExpenses: null,
        pickupAddress: {
          city: 'BEVERLY HILLS',
          county: 'LOS ANGELES',
          eTag: 'MjAyNS0wMy0yNVQxNTo1NjowNC41ODE2Mzla',
          id: 'b73bdf2d-2cc9-43a2-a18f-2859e5e3fb50',
          isOconus: false,
          postalCode: '90210',
          state: 'CA',
          streetAddress1: '123 Any Street',
          streetAddress2: 'P.O. Box 12345',
          streetAddress3: 'c/o Some Person',
          usPostRegionCitiesID: '3b9f0ae6-3b2b-44a6-9fcd-8ead346648c4',
        },
        ppmType: 'INCENTIVE_BASED',
        proGearWeight: null,
        proGearWeightTickets: null,
        reviewedAt: null,
        shipmentId: '82ccbd7e-4b87-4cb6-9d19-4f810096c42e',
        sitEstimatedCost: null,
        sitEstimatedDepartureDate: null,
        sitEstimatedEntryDate: null,
        sitEstimatedWeight: null,
        sitExpected: false,
        spouseProGearWeight: null,
        status: 'WAITING_ON_CUSTOMER',
        submittedAt: null,
        updatedAt: '2025-03-25T15:56:05.655Z',
        w2Address: {
          city: 'BEVERLY HILLS',
          county: 'LOS ANGELES',
          eTag: 'MjAyNS0wMy0yNVQxNTo1NjowNC41Nzc1Mzda',
          id: 'e744ea78-2f88-41be-af0d-b67a7757c619',
          isOconus: false,
          postalCode: '90212',
          state: 'CA',
          streetAddress1: '123 Any Street',
          streetAddress2: 'P.O. Box 12345',
          streetAddress3: 'c/o Some Person',
          usPostRegionCitiesID: 'dfbd8c2d-c92d-465d-a5ac-9fda39401d00',
        },
        weightTickets: [
          {
            adjustedNetWeight: null,
            createdAt: '2025-03-25T15:55:48.798Z',
            eTag: 'MjAyNS0wMy0yNVQxNTo1NjowNS42NjAzMTla',
            emptyDocument: {
              id: 'c01bab65-3026-4866-b063-4c2fc6dcd5af',
              service_member_id: 'a3e6390b-4d3e-467f-b515-aa618b3d703d',
              uploads: [
                {
                  bytes: 79011,
                  contentType: 'image/png',
                  createdAt: '2025-03-25T15:56:01.120Z',
                  filename: 'Screenshot 2025-01-17 at 11.37.10 AM.png-20250325115601',
                  id: '991f8548-cd46-4209-bc42-22166a844c11',
                  status: 'CLEAN',
                  updatedAt: '2025-03-25T15:56:01.120Z',
                  uploadType: 'USER',
                  url: '/storage/user/99bad872-00f5-4386-b546-ac9eac9b21ea/uploads/991f8548-cd46-4209-bc42-22166a844c11?contentType=image%2Fpng&filename=Screenshot+2025-01-17+at+11.37.10%E2%80%AFAM.png-20250325115601',
                },
              ],
            },
            emptyDocumentId: 'c01bab65-3026-4866-b063-4c2fc6dcd5af',
            emptyWeight: 1000,
            fullDocument: {
              id: 'c81d02be-cc86-4b97-a6e3-560fd213df3a',
              service_member_id: 'a3e6390b-4d3e-467f-b515-aa618b3d703d',
              uploads: [
                {
                  bytes: 121228,
                  contentType: 'image/png',
                  createdAt: '2025-03-25T15:56:03.467Z',
                  filename: 'Screenshot 2025-02-21 at 10.32.51 AM.png-20250325115603',
                  id: 'ba2eb17a-9eca-478a-bc33-48fe81850123',
                  status: 'CLEAN',
                  updatedAt: '2025-03-25T15:56:03.467Z',
                  uploadType: 'USER',
                  url: '/storage/user/99bad872-00f5-4386-b546-ac9eac9b21ea/uploads/ba2eb17a-9eca-478a-bc33-48fe81850123?contentType=image%2Fpng&filename=Screenshot+2025-02-21+at+10.32.51%E2%80%AFAM.png-20250325115603',
                },
              ],
            },
            fullDocumentId: 'c81d02be-cc86-4b97-a6e3-560fd213df3a',
            fullWeight: 2000,
            id: '63e81ade-8d9c-43e1-bc1d-e032bf9eb47d',
            missingEmptyWeightTicket: false,
            missingFullWeightTicket: false,
            netWeightRemarks: null,
            ownsTrailer: false,
            ppmShipmentId: 'b459e544-3654-430d-92d5-296dfca4009e',
            proofOfTrailerOwnershipDocument: {
              id: 'd52e3f8b-746e-4d84-9d00-93a8740e659b',
              service_member_id: 'a3e6390b-4d3e-467f-b515-aa618b3d703d',
              uploads: [],
            },
            proofOfTrailerOwnershipDocumentId: 'd52e3f8b-746e-4d84-9d00-93a8740e659b',
            reason: null,
            status: null,
            submittedEmptyWeight: null,
            submittedFullWeight: null,
            submittedOwnsTrailer: null,
            submittedTrailerMeetsCriteria: null,
            trailerMeetsCriteria: false,
            updatedAt: '2025-03-25T15:56:05.660Z',
            vehicleDescription: 'test',
          },
        ],
      },
      shipmentLocator: '8GHTH4-01',
      shipmentType: 'PPM',
      sitDaysAllowance: 90,
      status: 'APPROVED',
      updatedAt: '2025-03-25T15:56:05.658Z',
    },
  },
};

const testMoveCode = '1A5PM3';

const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
}));

const mockDispatch = jest.fn();
jest.mock('react-redux', () => ({
  ...jest.requireActual('react-redux'),
  useDispatch: jest.fn().mockImplementation(() => mockDispatch),
}));

jest.mock('services/ghcApi', () => ({
  ...jest.requireActual('services/internalApi'),
  getMTOShipments: jest.fn().mockImplementation(() => Promise.resolve()),
  getMove: jest.fn().mockImplementation(() => Promise.resolve()),
  submitPPMShipmentSignedCertification: jest.fn().mockImplementation(() => Promise.resolve()),
}));

jest.mock('store/entities/selectors', () => ({
  ...jest.requireActual('store/entities/selectors'),
  selectMTOShipmentById: jest.fn(),
}));

jest.mock('store/entities/actions', () => ({
  ...jest.requireActual('store/entities/actions'),
  updateMTOShipment: jest.fn(),
}));

updateMTOShipment.mockImplementation(() => Promise.resolve({}));

beforeEach(() => {
  jest.clearAllMocks();
});

const ppmReviewPath = generatePath(servicesCounselingRoutes.BASE_SHIPMENT_PPM_REVIEW_PATH, {
  moveCode: testMoveCode,
  shipmentId: shipmentID,
});

describe('Final Closeout page', () => {
  getMove.mockResolvedValue(testMove);
  getMTOShipments.mockResolvedValue(response);

  it('loads the selected shipment from redux', async () => {
    const mockRoutingConfig = {
      path: servicesCounselingRoutes.BASE_SHIPMENT_PPM_COMPLETE_PATH,
      params: { moveCode: testMoveCode, shipmentId: response.mtoShipments[shipmentID].id },
    };

    render(
      <MockProviders {...mockRoutingConfig}>
        <FinalCloseout />
      </MockProviders>,
    );

    await waitFor(() => {
      expect(screen.getByTestId('scCompletePPMHeader')).toBeInTheDocument();
      expect(screen.queryByTestId('loading-placeholder')).toBeNull();
    });
  });

  it('renders the page headings', async () => {
    const mockRoutingConfig = {
      path: servicesCounselingRoutes.BASE_SHIPMENT_PPM_COMPLETE_PATH,
      params: { moveCode: testMoveCode, shipmentId: response.mtoShipments[shipmentID].id },
    };

    render(
      <MockProviders {...mockRoutingConfig}>
        <FinalCloseout />
      </MockProviders>,
    );

    await waitFor(() => {
      expect(screen.getByTestId('tag')).toHaveTextContent('PPM');
    });

    expect(screen.getByRole('heading', { level: 1, name: 'Complete PPM' })).toBeInTheDocument();
    expect(screen.getByRole('heading', { level: 2, name: /Your final estimated incentive: \$/ })).toBeInTheDocument();
  });

  it('routes to the home page when the return to homepage link is clicked', async () => {
    const mockRoutingConfig = {
      path: servicesCounselingRoutes.BASE_SHIPMENT_PPM_COMPLETE_PATH,
      params: { moveCode: testMoveCode, shipmentId: response.mtoShipments[shipmentID].id },
    };

    render(
      <MockProviders {...mockRoutingConfig}>
        <FinalCloseout />
      </MockProviders>,
    );

    await waitFor(async () => {
      await userEvent.click(screen.getByText('Return To Homepage'));
    });

    expect(mockNavigate).toHaveBeenCalledWith(ppmReviewPath);
  });

  it('submits the ppm signed certification', async () => {
    submitPPMShipmentSignedCertification.mockResolvedValueOnce(response.mtoShipments[0].ppmShipment);

    const mockRoutingConfig = {
      path: servicesCounselingRoutes.BASE_SHIPMENT_PPM_COMPLETE_PATH,
      params: { moveCode: testMoveCode, shipmentId: response.mtoShipments[shipmentID].id },
    };

    render(
      <MockProviders {...mockRoutingConfig}>
        <FinalCloseout />
      </MockProviders>,
    );
    await waitFor(() => {
      expect(screen.getByTestId('tag')).toHaveTextContent('PPM');
    });

    await userEvent.type(screen.getByRole('textbox', { name: 'Signature' }), 'Grace Griffin');
    await userEvent.click(screen.getByRole('button', { name: 'Submit PPM Documentation' }));

    await waitFor(() =>
      expect(submitPPMShipmentSignedCertification).toHaveBeenCalledWith(response.mtoShipments[0].ppmShipment.id, {
        certification_text: ppmSubmissionCertificationText,
        signature: 'Grace Griffin',
        date: formatDateForSwagger(new Date()),
      }),
    );

    expect(updateMTOShipment).toHaveBeenCalledWith(response.mtoShipments[shipmentID]);
    expect(mockDispatch).toHaveBeenCalledTimes(2);

    expect(mockNavigate).toHaveBeenCalledWith(servicesCounselingRoutes.BASE_MOVE_VIEW_PATH);
  });
});

import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { v4 as uuidv4 } from 'uuid';
import { generatePath } from 'react-router';

import FinalCloseout from 'pages/Office/PPM/Closeout/FinalCloseout/FinalCloseout';
import { updateMTOShipment } from 'store/entities/actions';
import { MockProviders } from 'testUtils';
import { servicesCounselingRoutes } from 'constants/routes';
import { submitPPMShipmentSignedCertification } from 'services/ghcApi';
import { useEditShipmentQueries } from 'hooks/queries';

const shipmentID = uuidv4();

const useEditShipmentQueriesReturnValue = {
  move: {
    id: '9c7b255c-2981-4bf8-839f-61c7458e2b4d',
    ordersId: '1',
    status: 'NEEDS SERVICE COUNSELING',
    closeoutOffice: {
      name: 'Altus AFB',
    },
  },
  order: {
    agency: 'ARMY',
    id: '1',
    originDutyLocation: {
      address: {
        streetAddress1: '',
        city: 'Fort Knox',
        state: 'KY',
        postalCode: '40121',
      },
    },
    destinationDutyLocation: {
      address: {
        streetAddress1: '',
        city: 'Fort Irwin',
        state: 'CA',
        postalCode: '92310',
      },
    },
    customer: {
      agency: 'ARMY',
      backup_contact: {
        email: 'email@example.com',
        name: 'name',
        phone: '555-555-5555',
      },
      current_address: {
        city: 'Beverly Hills',
        country: 'US',
        eTag: 'MjAyMS0wMS0yMVQxNTo0MTozNS41Mzg0Njha',
        id: '3a5f7cf2-6193-4eb3-a244-14d21ca05d7b',
        postalCode: '90210',
        state: 'CA',
        streetAddress1: '123 Any Street',
        streetAddress2: 'P.O. Box 12345',
        streetAddress3: 'c/o Some Person',
      },
      dodID: '6833908165',
      eTag: 'MjAyMS0wMS0yMVQxNTo0MTozNS41NjAzNTJa',
      email: 'combo@ppm.hhg',
      first_name: 'Submitted',
      id: 'f6bd793f-7042-4523-aa30-34946e7339c9',
      last_name: 'Ppmhhg',
      phone: '555-555-5555',
    },
    entitlement: {
      authorizedWeight: 8000,
      dependentsAuthorized: true,
      eTag: 'MjAyMS0wMS0yMVQxNTo0MTozNS41NzgwMzda',
      id: 'e0fefe58-0710-40db-917b-5b96567bc2a8',
      nonTemporaryStorage: true,
      privatelyOwnedVehicle: true,
      proGearWeight: 2000,
      proGearWeightSpouse: 500,
      storageInTransit: 2,
      totalDependents: 1,
      totalWeight: 8000,
    },
    order_number: 'ORDER3',
    order_type: 'PERMANENT_CHANGE_OF_STATION',
    order_type_detail: 'HHG_PERMITTED',
    tac: '9999',
  },
  mtoShipments: [
    {
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
  ],
  isLoading: false,
  isError: false,
  isSuccess: true,
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

jest.mock('hooks/queries', () => ({
  usePPMShipmentDocsQueries: jest.fn(),
  useEditShipmentQueries: jest.fn(),
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
  useEditShipmentQueries.mockReturnValue(useEditShipmentQueriesReturnValue);

  it('loads the selected shipment from redux', async () => {
    const mockRoutingConfig = {
      path: servicesCounselingRoutes.BASE_SHIPMENT_PPM_COMPLETE_PATH,
      params: { moveCode: testMoveCode, shipmentId: shipmentID },
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

  it('renders the page headings and closeout office name', async () => {
    const mockRoutingConfig = {
      path: servicesCounselingRoutes.BASE_SHIPMENT_PPM_COMPLETE_PATH,
      params: { moveCode: testMoveCode, shipmentId: shipmentID },
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
    expect(screen.getByText(useEditShipmentQueriesReturnValue.move.closeoutOffice.name, { exact: false }));
  });

  it('routes to the home page when the return to homepage link is clicked', async () => {
    const mockRoutingConfig = {
      path: servicesCounselingRoutes.BASE_SHIPMENT_PPM_COMPLETE_PATH,
      params: { moveCode: testMoveCode, shipmentId: shipmentID },
    };

    render(
      <MockProviders {...mockRoutingConfig}>
        <FinalCloseout />
      </MockProviders>,
    );

    await waitFor(async () => {
      await userEvent.click(screen.getByText('Back'));
    });

    expect(mockNavigate).toHaveBeenCalledWith(ppmReviewPath);
  });

  it('submits the ppm signed certification', async () => {
    submitPPMShipmentSignedCertification.mockResolvedValueOnce(
      useEditShipmentQueriesReturnValue.mtoShipments[0].ppmShipment,
    );

    const mockRoutingConfig = {
      path: servicesCounselingRoutes.BASE_SHIPMENT_PPM_COMPLETE_PATH,
      params: { moveCode: testMoveCode, shipmentId: shipmentID },
    };

    render(
      <MockProviders {...mockRoutingConfig}>
        <FinalCloseout />
      </MockProviders>,
    );
    await waitFor(() => {
      expect(screen.getByTestId('tag')).toHaveTextContent('PPM');
    });

    await userEvent.click(screen.getByRole('button', { name: 'Submit PPM Documentation' }));

    await waitFor(() =>
      expect(submitPPMShipmentSignedCertification).toHaveBeenCalledWith(
        useEditShipmentQueriesReturnValue.mtoShipments[0].ppmShipment.id,
      ),
    );
    expect(mockNavigate).toHaveBeenCalledWith(`/counseling/moves/${testMoveCode}/details`);
  });
});

import React from 'react';
import { waitFor, screen, fireEvent, act, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import HeaderSection from './HeaderSection';

import { useEditShipmentQueries, usePPMShipmentDocsQueries } from 'hooks/queries';
import { renderWithProviders } from 'testUtils';

beforeEach(() => {
  jest.clearAllMocks();
});

const routingParams = { moveCode: 'move123', shipmentId: 'shipment123' };
const mockRoutingConfig = {
  params: routingParams,
};

jest.mock('hooks/queries', () => ({
  usePPMShipmentDocsQueries: jest.fn(),
  useEditShipmentQueries: jest.fn(),
}));

const mockUpdateMTOShipment = jest.fn();
jest.mock('services/ghcApi', () => ({
  ...jest.requireActual('services/ghcApi'),
  updateMTOShipment: (options) => mockUpdateMTOShipment(options),
}));

const usePPMShipmentDocsQueriesReturnValue = {
  mtoShipment: {
    id: 'shipment123',
    moveTaskOrderID: 'move123',
    eTag: 'etag123',
    ppmShipment: {
      id: 'ppm123',
      actualMoveDate: '2022-01-12',
      isActualExpenseReimbursement: false,
    },
    pickupAddress: {
      city: 'Beverly Hills',
      country: 'US',
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi4zODQ3Njla',
      id: '1686751b-ab36-43cf-b3c9-c0f467d13c19',
      postalCode: '90210',
      state: 'CA',
      streetAddress1: '123 Any Street',
      streetAddress2: 'P.O. Box 12345',
      streetAddress3: 'c/o Some Person',
    },
    destinationAddress: {
      city: 'Fairfield',
      country: 'US',
      id: '672ff379-f6e3-48b4-a87d-796713f8f997',
      postalCode: '94535',
      state: 'CA',
      streetAddress1: '987 Any Avenue',
      streetAddress2: 'P.O. Box 9876',
      streetAddress3: 'c/o Some Person',
    },
  },
  refetchMTOShipment: jest.fn(), // Mock the refetch function
  isFetching: false,
};

const useEditShipmentQueriesReturnValue = {
  move: {
    id: '9c7b255c-2981-4bf8-839f-61c7458e2b4d',
    ordersId: '1',
    status: 'NEEDS SERVICE COUNSELING',
  },
  order: {
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
      customerRemarks: 'please treat gently',
      destinationAddress: {
        city: 'Fairfield',
        country: 'US',
        id: '672ff379-f6e3-48b4-a87d-796713f8f997',
        postalCode: '94535',
        state: 'CA',
        streetAddress1: '987 Any Avenue',
        streetAddress2: 'P.O. Box 9876',
        streetAddress3: 'c/o Some Person',
      },
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi40MDQwMzFa',
      id: 'shipment123',
      moveTaskOrderID: '9c7b255c-2981-4bf8-839f-61c7458e2b4d',
      pickupAddress: {
        city: 'Beverly Hills',
        country: 'US',
        eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi4zODQ3Njla',
        id: '1686751b-ab36-43cf-b3c9-c0f467d13c19',
        postalCode: '90210',
        state: 'CA',
        streetAddress1: '123 Any Street',
        streetAddress2: 'P.O. Box 12345',
        streetAddress3: 'c/o Some Person',
      },
      requestedPickupDate: '2018-03-15',
      scheduledPickupDate: '2018-03-16',
      requestedDeliveryDate: '2018-04-15',
      scheduledDeliveryDate: '2014-04-16',
      shipmentType: 'HHG',
      status: 'SUBMITTED',
      updatedAt: '2020-06-10T15:58:02.404031Z',
    },
  ],
  isLoading: false,
  isError: false,
  isSuccess: true,
  refetchMTOShipment: jest.fn(),
  isFetching: false,
};

const ppmShipmentInfoProps = {
  sectionInfo: {
    type: 'shipmentInfo',
    plannedMoveDate: '2020-03-15',
    actualMoveDate: '2022-01-12',
    pickupAddress: '812 S 129th St, #123, San Antonio, TX 78234',
    destinationAddress: '456 Oak Ln., #123, Oakland, CA 94611',
    miles: 513,
    estimatedWeight: 4000,
    actualWeight: 4200,
    allowableWeight: 4300,
  },
  setUpdatedItemName: jest.fn(),
  setIsSubmitting: jest.fn(),
};

const incentivesProps = {
  sectionInfo: {
    type: 'incentives',
    isAdvanceRequested: true,
    isAdvanceReceived: true,
    advanceAmountRequested: 598700,
    advanceAmountReceived: 112244,
    grossIncentive: 7231285,
    gcc: 7231285,
    remainingIncentive: 7119041,
  },
  setUpdatedItemName: jest.fn(),
  setIsSubmitting: jest.fn(),
};
const incentivesAdvanceReceivedZeroProps = {
  sectionInfo: {
    type: 'incentives',
    isAdvanceRequested: false,
    isAdvanceReceived: true,
    advanceAmountRequested: 598700,
    advanceAmountReceived: 0,
    grossIncentive: 7231285,
    gcc: 7231285,
    remainingIncentive: 7119041,
  },
  setUpdatedItemName: jest.fn(),
  setIsSubmitting: jest.fn(),
};

const incentiveFactorsProps = {
  sectionInfo: {
    type: 'incentiveFactors',
    haulType: 'Linehaul',
    haulPrice: 6892668,
    haulFSC: -143,
    packPrice: 20000,
    unpackPrice: 10000,
    dop: 15640,
    ddp: 34640,
    sitReimbursement: 30000,
  },
};

const incentiveFactorsShorthaulProps = {
  sectionInfo: {
    type: 'incentiveFactors',
    haulType: 'Shorthaul',
    haulPrice: 6892668,
    haulFSC: -143,
    packPrice: 20000,
    unpackPrice: 10000,
    dop: 15640,
    ddp: 34640,
    sitReimbursement: 30000,
  },
};

const invalidSectionTypeProps = {
  sectionInfo: {
    type: 'someUnknownSectionType',
  },
};

const clickDetailsButton = async (buttonType) => {
  await act(async () => {
    await fireEvent.click(screen.getByTestId(`${buttonType}-showRequestDetailsButton`));
  });
  await waitFor(() => {
    expect(screen.getByText('Hide Details', { exact: false })).toBeInTheDocument();
  });
};

describe('PPMHeaderSummary component', () => {
  describe('displays Shipment Info section', () => {
    it('renders Shipment Info section on load with defaults', async () => {
      usePPMShipmentDocsQueries.mockReturnValue(usePPMShipmentDocsQueriesReturnValue);
      useEditShipmentQueries.mockReturnValue(useEditShipmentQueriesReturnValue);
      await act(async () => {
        renderWithProviders(<HeaderSection {...ppmShipmentInfoProps} />, mockRoutingConfig);
      });

      await waitFor(() => {
        expect(screen.getByRole('heading', { level: 4, name: 'Shipment Info' })).toBeInTheDocument();
      });
      await act(async () => {
        clickDetailsButton('shipmentInfo');
      });

      expect(screen.getByText('Actual Expense Reimbursement')).toBeInTheDocument();
      expect(screen.getByText('Planned Move Start Date')).toBeInTheDocument();
      expect(screen.getByText('15-Mar-2020')).toBeInTheDocument();
      expect(screen.getByText('Actual Move Start Date')).toBeInTheDocument();
      expect(screen.getByText('12-Jan-2022')).toBeInTheDocument();
      expect(screen.getByText('Starting Address')).toBeInTheDocument();
      expect(screen.getByText('812 S 129th St, #123, San Antonio, TX 78234')).toBeInTheDocument();
      expect(screen.getByText('Ending Address')).toBeInTheDocument();
      expect(screen.getByText('456 Oak Ln., #123, Oakland, CA 94611')).toBeInTheDocument();
      expect(screen.getByText('Miles')).toBeInTheDocument();
      expect(screen.getByText('513')).toBeInTheDocument();
      expect(screen.getByText('Estimated Net Weight')).toBeInTheDocument();
      expect(screen.getByText('4,000 lbs')).toBeInTheDocument();
      expect(screen.getByText('Actual Net Weight')).toBeInTheDocument();
      expect(screen.getByText('4,200 lbs')).toBeInTheDocument();
      expect(screen.getByText('Allowable Weight')).toBeInTheDocument();
      expect(screen.getByText('4,300 lbs')).toBeInTheDocument();
    });
  });

  describe('displays "Incentives/Costs" section', () => {
    it('renders "Incentives/Costs" section on load with correct prop values', async () => {
      usePPMShipmentDocsQueries.mockReturnValue(usePPMShipmentDocsQueriesReturnValue);
      useEditShipmentQueries.mockReturnValue(useEditShipmentQueriesReturnValue);
      await act(async () => {
        renderWithProviders(<HeaderSection {...incentivesProps} />, mockRoutingConfig);
      });
      await act(async () => {
        clickDetailsButton('incentives');
      });

      expect(screen.getByText('Government Constructed Cost (GCC)')).toBeInTheDocument();
      expect(screen.getByTestId('gcc')).toHaveTextContent('$72,312.85');
      expect(screen.getByText('Gross Incentive')).toBeInTheDocument();
      expect(screen.getByTestId('grossIncentive')).toHaveTextContent('$72,312.85');
      expect(screen.getByText('Advance Requested')).toBeInTheDocument();
      expect(screen.getByTestId('advanceRequested')).toHaveTextContent('$5,987.00');
      expect(screen.getByText('Advance Received')).toBeInTheDocument();
      expect(screen.getByTestId('advanceReceived')).toHaveTextContent('$1,122.44');
      expect(screen.getByText('Remaining Incentive')).toBeInTheDocument();
      expect(screen.getByTestId('remainingIncentive')).toHaveTextContent('$71,190.41');
    });
  });

  describe('edit items correctly', () => {
    it('edits actual expense reimbursement correctly', async () => {
      usePPMShipmentDocsQueries.mockReturnValue(usePPMShipmentDocsQueriesReturnValue);
      useEditShipmentQueries.mockReturnValue(useEditShipmentQueriesReturnValue);
      await act(async () => {
        renderWithProviders(<HeaderSection {...ppmShipmentInfoProps} />, mockRoutingConfig);
      });
      await act(async () => {
        clickDetailsButton('shipmentInfo');
      });

      const modalButton = within(screen.getByTestId('isActualExpenseReimbursement')).getByTestId('editTextButton');
      await act(() => userEvent.click(modalButton));
      await act(() => userEvent.click(screen.getByText('Yes')));
      await act(() => userEvent.click(screen.getByRole('button', { name: 'Save' })));

      await act(async () => {
        expect(mockUpdateMTOShipment).toHaveBeenCalledTimes(1);
      });

      await act(async () => {
        expect(mockUpdateMTOShipment).toHaveBeenCalledWith({
          moveTaskOrderID: 'move123',
          shipmentID: 'shipment123',
          ifMatchETag: 'etag123',
          body: {
            ppmShipment: {
              isActualExpenseReimbursement: true,
            },
          },
        });
      });
    });

    it('edits actual move date correctly', async () => {
      usePPMShipmentDocsQueries.mockReturnValue(usePPMShipmentDocsQueriesReturnValue);
      useEditShipmentQueries.mockReturnValue(useEditShipmentQueriesReturnValue);
      await act(async () => {
        renderWithProviders(<HeaderSection {...ppmShipmentInfoProps} />, mockRoutingConfig);
      });
      await act(async () => {
        clickDetailsButton('shipmentInfo');
      });

      const modalButton = within(screen.getByTestId('actualMoveDate')).getByTestId('editTextButton');
      await act(() => userEvent.click(modalButton));
      await act(() => userEvent.click(screen.getByRole('button', { name: 'Save' })));

      await act(async () => {
        expect(mockUpdateMTOShipment).toHaveBeenCalledTimes(1);
      });

      await act(async () => {
        expect(mockUpdateMTOShipment).toHaveBeenCalledWith({
          moveTaskOrderID: 'move123',
          shipmentID: 'shipment123',
          ifMatchETag: 'etag123',
          body: {
            ppmShipment: {
              actualMoveDate: '2022-01-20',
            },
          },
        });
      });
    });

    it('edits advance amount received correctly', async () => {
      usePPMShipmentDocsQueries.mockReturnValue(usePPMShipmentDocsQueriesReturnValue);
      useEditShipmentQueries.mockReturnValue(useEditShipmentQueriesReturnValue);
      await act(async () => {
        renderWithProviders(<HeaderSection {...incentivesProps} />, mockRoutingConfig);
      });
      await act(async () => {
        clickDetailsButton('incentives');
      });

      const modalButton = within(screen.getByTestId('advanceReceived')).getByTestId('editTextButton');
      await act(() => userEvent.click(modalButton));
      await act(() => userEvent.click(screen.getByRole('button', { name: 'Save' })));

      await act(async () => {
        expect(mockUpdateMTOShipment).toHaveBeenCalledTimes(1);
      });

      await act(async () => {
        expect(mockUpdateMTOShipment).toHaveBeenCalledWith({
          moveTaskOrderID: 'move123',
          shipmentID: 'shipment123',
          ifMatchETag: 'etag123',
          body: {
            ppmShipment: {
              advanceAmountReceived: 112200,
              hasReceivedAdvance: true,
            },
          },
        });
      });
    });

    it('if advance amount received is 0', async () => {
      usePPMShipmentDocsQueries.mockReturnValue(usePPMShipmentDocsQueriesReturnValue);
      useEditShipmentQueries.mockReturnValue(useEditShipmentQueriesReturnValue);
      await act(async () => {
        renderWithProviders(<HeaderSection {...incentivesAdvanceReceivedZeroProps} />, mockRoutingConfig);
      });
      await act(async () => {
        clickDetailsButton('incentives');
      });

      const modalButton = within(screen.getByTestId('advanceReceived')).getByTestId('editTextButton');
      await act(() => userEvent.click(modalButton));
      await act(() => userEvent.click(screen.getByRole('button', { name: 'Save' })));

      await act(async () => {
        expect(mockUpdateMTOShipment).toHaveBeenCalledTimes(1);
      });

      await act(async () => {
        expect(mockUpdateMTOShipment).toHaveBeenCalledWith({
          moveTaskOrderID: 'move123',
          shipmentID: 'shipment123',
          ifMatchETag: 'etag123',
          body: {
            ppmShipment: {
              advanceAmountReceived: null,
              hasReceivedAdvance: false,
            },
          },
        });
      });
    });
  });

  describe('displays "Incentive Factors" section', () => {
    it('renders "Incentive Factors" on load with correct prop values', async () => {
      usePPMShipmentDocsQueries.mockReturnValue(usePPMShipmentDocsQueriesReturnValue);
      useEditShipmentQueries.mockReturnValue(useEditShipmentQueriesReturnValue);
      await act(async () => {
        renderWithProviders(<HeaderSection {...incentiveFactorsProps} />, mockRoutingConfig);
      });
      await act(async () => {
        clickDetailsButton('incentiveFactors');
      });

      expect(screen.getByText('Linehaul Price')).toBeInTheDocument();
      expect(screen.getByTestId('haulPrice')).toHaveTextContent('$68,926.68');
      expect(screen.getByText('Linehaul Fuel Rate Adjustment')).toBeInTheDocument();
      expect(screen.getByTestId('haulFSC')).toHaveTextContent('-$1.43');
      expect(screen.getByText('Packing Charge')).toBeInTheDocument();
      expect(screen.getByTestId('packPrice')).toHaveTextContent('$200.00');
      expect(screen.getByText('Unpacking Charge')).toBeInTheDocument();
      expect(screen.getByTestId('unpackPrice')).toHaveTextContent('$100.00');
      expect(screen.getByText('Origin Price')).toBeInTheDocument();
      expect(screen.getByTestId('originPrice')).toHaveTextContent('$156.40');
      expect(screen.getByText('Destination Price')).toBeInTheDocument();
      expect(screen.getByTestId('destinationPrice')).toHaveTextContent('$346.40');
      expect(screen.getByTestId('sitReimbursement')).toHaveTextContent('$300.00');
    });

    it('renders "Shorthaul" in place of linehaul when given a shorthaul type', async () => {
      usePPMShipmentDocsQueries.mockReturnValue(usePPMShipmentDocsQueriesReturnValue);
      useEditShipmentQueries.mockReturnValue(useEditShipmentQueriesReturnValue);
      await act(async () => {
        renderWithProviders(<HeaderSection {...incentiveFactorsShorthaulProps} />, mockRoutingConfig);
      });
      await act(async () => {
        clickDetailsButton('incentiveFactors');
      });

      expect(screen.getByText('Shorthaul Price')).toBeInTheDocument();
      expect(screen.getByTestId('haulPrice')).toHaveTextContent('$68,926.68');
      expect(screen.getByText('Shorthaul Fuel Rate Adjustment')).toBeInTheDocument();
      expect(screen.getByTestId('haulFSC')).toHaveTextContent('-$1.43');
      expect(screen.getByText('Packing Charge')).toBeInTheDocument();
      expect(screen.getByTestId('packPrice')).toHaveTextContent('$200.00');
      expect(screen.getByText('Unpacking Charge')).toBeInTheDocument();
      expect(screen.getByTestId('unpackPrice')).toHaveTextContent('$100.00');
      expect(screen.getByText('Origin Price')).toBeInTheDocument();
      expect(screen.getByTestId('originPrice')).toHaveTextContent('$156.40');
      expect(screen.getByText('Destination Price')).toBeInTheDocument();
      expect(screen.getByTestId('destinationPrice')).toHaveTextContent('$346.40');
      expect(screen.getByTestId('sitReimbursement')).toHaveTextContent('$300.00');
    });
  });

  describe('handles errors correctly', () => {
    it('renders an alert if an unknown section type was passed in', async () => {
      usePPMShipmentDocsQueries.mockReturnValue(usePPMShipmentDocsQueriesReturnValue);
      useEditShipmentQueries.mockReturnValue(useEditShipmentQueriesReturnValue);
      await act(async () => {
        renderWithProviders(<HeaderSection {...invalidSectionTypeProps} />, mockRoutingConfig);
      });

      const alert = screen.getByTestId('alert');
      expect(alert).toBeInTheDocument();
      expect(alert).toHaveTextContent('Error getting section title!');
    });

    it('renders an alert if an unknown section type was passed in and details are expanded', async () => {
      usePPMShipmentDocsQueries.mockReturnValue(usePPMShipmentDocsQueriesReturnValue);
      useEditShipmentQueries.mockReturnValue(useEditShipmentQueriesReturnValue);
      await act(async () => {
        renderWithProviders(<HeaderSection {...invalidSectionTypeProps} />, mockRoutingConfig);
      });
      await act(async () => {
        clickDetailsButton(invalidSectionTypeProps.sectionInfo.type);
      });
      expect(screen.getByText('An error occured while getting section markup!')).toBeInTheDocument();
    });
  });
});

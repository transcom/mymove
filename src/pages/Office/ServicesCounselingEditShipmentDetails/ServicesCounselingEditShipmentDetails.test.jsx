/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import ServicesCounselingEditShipmentDetails from './ServicesCounselingEditShipmentDetails';

import { updateMTOShipment } from 'services/ghcApi';

const mockPush = jest.fn();

jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useLocation: () => ({
    pathname: 'localhost:3000/',
  }),
  useHistory: () => ({
    push: mockPush,
    goBack: jest.fn(),
  }),
  useParams: jest.fn().mockReturnValue({ moveCode: 'move123', shipmentId: 'shipment123' }),
}));

jest.mock('services/ghcApi', () => ({
  ...jest.requireActual('services/ghcApi'),
  updateMTOShipment: jest.fn(),
}));

/*
The order in which these mockImplementationOnce is directly correlated to the
order in which the tests are run. This is because different data should be
returned from the call each time in order to tests the different states that
the component can be in. The default will be to return all the data needed for
the full tests.
*/
jest.mock('hooks/queries', () => ({
  useEditShipmentQueries: jest
    .fn(() => {
      return {
        order: {
          id: '1',
          originDutyStation: {
            address: {
              street_address_1: '',
              city: 'Fort Knox',
              state: 'KY',
              postal_code: '40121',
            },
          },
          destinationDutyStation: {
            address: {
              street_address_1: '',
              city: 'Fort Irwin',
              state: 'CA',
              postal_code: '92310',
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
              postal_code: '90210',
              state: 'CA',
              street_address_1: '123 Any Street',
              street_address_2: 'P.O. Box 12345',
              street_address_3: 'c/o Some Person',
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
              postal_code: '94535',
              state: 'CA',
              street_address_1: '987 Any Avenue',
              street_address_2: 'P.O. Box 9876',
              street_address_3: 'c/o Some Person',
            },
            eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi40MDQwMzFa',
            id: 'shipment123',
            moveTaskOrderID: '9c7b255c-2981-4bf8-839f-61c7458e2b4d',
            pickupAddress: {
              city: 'Beverly Hills',
              country: 'US',
              eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi4zODQ3Njla',
              id: '1686751b-ab36-43cf-b3c9-c0f467d13c19',
              postal_code: '90210',
              state: 'CA',
              street_address_1: '123 Any Street',
              street_address_2: 'P.O. Box 12345',
              street_address_3: 'c/o Some Person',
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
      };
    })
    .mockImplementationOnce(() => {
      return {
        isLoading: true,
        isError: false,
        isSuccess: false,
      };
    })
    .mockImplementationOnce(() => {
      return {
        isLoading: false,
        isError: true,
        isSuccess: false,
      };
    }),
}));

const props = {
  match: {
    path: '',
    isExact: false,
    url: '',
    params: { moveCode: 'move123', shipmentId: 'shipment123' },
  },
};

describe('ServicesCounselingEditShipmentDetails component', () => {
  it('renders the Loading Placeholder when the query is still loading', async () => {
    render(<ServicesCounselingEditShipmentDetails {...props} />);

    const h2 = await screen.getByRole('heading', { name: 'Loading, please wait...', level: 2 });
    expect(h2).toBeInTheDocument();
  });

  it('renders the Something Went Wrong component when the query errors', async () => {
    render(<ServicesCounselingEditShipmentDetails {...props} />);

    const errorMessage = await screen.getByText(/Something went wrong./);
    expect(errorMessage).toBeInTheDocument();
  });

  it('renders the Services Counseling Shipment Form', async () => {
    render(<ServicesCounselingEditShipmentDetails {...props} />);

    const h1 = await screen.getByRole('heading', { name: 'Edit shipment details', level: 1 });
    expect(h1).toBeInTheDocument();
  });

  it('saves the update to the counselor remarks when the save button is clicked', async () => {
    const newCounselorRemarks = 'Counselor remarks';

    const expectedPayload = {
      body: {
        customerRemarks: 'please treat gently',
        destinationAddress: {
          city: 'Fairfield',
          country: 'US',
          postal_code: '94535',
          state: 'CA',
          street_address_1: '987 Any Avenue',
          street_address_2: 'P.O. Box 9876',
          street_address_3: 'c/o Some Person',
        },
        pickupAddress: {
          city: 'Beverly Hills',
          country: 'US',
          eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi4zODQ3Njla',
          postal_code: '90210',
          state: 'CA',
          street_address_1: '123 Any Street',
          street_address_2: 'P.O. Box 12345',
          street_address_3: 'c/o Some Person',
        },
        requestedDeliveryDate: '2018-04-15',
        requestedPickupDate: '2018-03-15',
        shipmentType: 'HHG',
      },
      shipmentID: 'shipment123',
      ifMatchETag: 'MjAyMC0wNi0xMFQxNTo1ODowMi40MDQwMzFa',
      moveTaskOrderID: '9c7b255c-2981-4bf8-839f-61c7458e2b4d',
      normalize: false,
    };

    const patchResponse = {
      ...expectedPayload,
      created_at: '2021-02-08T16:48:04.117Z',
      updated_at: '2021-02-11T16:48:04.117Z',
    };

    updateMTOShipment.mockImplementation(() => Promise.resolve(patchResponse));

    render(<ServicesCounselingEditShipmentDetails {...props} />);

    const counselorRemarks = await screen.findByLabelText('Counselor remarks');

    userEvent.clear(counselorRemarks);

    userEvent.type(counselorRemarks, newCounselorRemarks);

    const saveButton = screen.getByRole('button', { name: 'Save' });

    expect(saveButton).not.toBeDisabled();

    userEvent.click(saveButton);

    await waitFor(() => {
      expect(updateMTOShipment).toHaveBeenCalledWith(expectedPayload);
    });
  });

  it('shows an error if the patchServiceMember API returns an error', async () => {
    updateMTOShipment.mockImplementation(() =>
      // Disable this rule because makeSwaggerRequest does not throw an error if the API call fails
      // eslint-disable-next-line prefer-promise-reject-errors
      Promise.reject({
        message: 'A server error occurred editing the shipment details',
        response: {
          body: {
            detail: 'A server error occurred editing the shipment details',
          },
        },
      }),
    );

    render(<ServicesCounselingEditShipmentDetails {...props} />);

    const saveButton = screen.getByRole('button', { name: 'Save' });

    userEvent.click(saveButton);

    await waitFor(() => {
      expect(updateMTOShipment).toHaveBeenCalled();
    });

    expect(await screen.findByText('failed to update MTO shipment due to server error')).toBeInTheDocument();
    expect(mockPush).not.toHaveBeenCalled();
  });

  it('routes to the profile page when the save button is clicked', async () => {
    updateMTOShipment.mockImplementation(() => Promise.resolve());

    render(<ServicesCounselingEditShipmentDetails {...props} />);

    const saveButton = screen.getByRole('button', { name: 'Save' });

    expect(saveButton).not.toBeDisabled();

    userEvent.click(saveButton);

    await waitFor(() => {
      expect(mockPush).toHaveBeenCalledWith('/counseling/moves/move123/details');
    });
  });
});

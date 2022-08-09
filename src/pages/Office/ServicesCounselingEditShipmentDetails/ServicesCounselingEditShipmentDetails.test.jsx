/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { render, screen, waitFor, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import ServicesCounselingEditShipmentDetails from './ServicesCounselingEditShipmentDetails';

import { updateMTOShipment } from 'services/ghcApi';
import { validatePostalCode } from 'utils/validation';
import { useEditShipmentQueries } from 'hooks/queries';
import { MOVE_STATUSES, SHIPMENT_OPTIONS } from 'shared/constants';

const mockPush = jest.fn();

jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useLocation: () => ({
    pathname: 'localhost:3000/',
  }),
  useHistory: () => ({
    push: mockPush,
  }),
  useParams: jest.fn().mockReturnValue({ moveCode: 'move123', shipmentId: 'shipment123' }),
}));

jest.mock('services/ghcApi', () => ({
  ...jest.requireActual('services/ghcApi'),
  updateMTOShipment: jest.fn(),
}));

jest.mock('hooks/queries', () => ({
  useEditShipmentQueries: jest.fn(),
}));

jest.mock('utils/validation', () => ({
  ...jest.requireActual('utils/validation'),
  validatePostalCode: jest.fn(),
}));

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
      shipmentType: SHIPMENT_OPTIONS.HHG,
      status: 'SUBMITTED',
      updatedAt: '2020-06-10T15:58:02.404031Z',
    },
  ],
  isLoading: false,
  isError: false,
  isSuccess: true,
};

const ppmShipment = {
  id: 'shipment123',
  shipmentType: SHIPMENT_OPTIONS.PPM,
  status: MOVE_STATUSES.SUBMITTED,
  updatedAt: '2020-09-02T21:08:38.392Z',
  ppmShipment: {
    expectedDepartureDate: '2022-06-28',
    actualMoveDate: '2022-05-11',
    pickupPostalCode: '90210',
    secondaryPickupPostalCode: '90002',
    destinationPostalCode: '10108',
    secondaryDestinationPostalCode: '79329',
    sitExpected: false,
    estimatedWeight: 1111,
    netWeight: 3333,
    hasProGear: false,
  },
};

const loadingReturnValue = {
  ...useEditShipmentQueriesReturnValue,
  isLoading: true,
  isError: false,
  isSuccess: false,
};

const errorReturnValue = {
  ...useEditShipmentQueriesReturnValue,
  isLoading: false,
  isError: true,
  isSuccess: false,
};

const props = {
  onUpdate: () => {},
  match: {
    path: '',
    isExact: false,
    url: '',
    params: { moveCode: 'move123', shipmentId: 'shipment123' },
    onUpdate: () => {},
  },
};

describe('ServicesCounselingEditShipmentDetails component', () => {
  describe('check different component states', () => {
    it('renders the Loading Placeholder when the query is still loading', async () => {
      useEditShipmentQueries.mockReturnValue(loadingReturnValue);

      render(<ServicesCounselingEditShipmentDetails {...props} />);

      const h2 = await screen.getByRole('heading', { name: 'Loading, please wait...', level: 2 });
      expect(h2).toBeInTheDocument();
    });

    it('renders the Something Went Wrong component when the query errors', async () => {
      useEditShipmentQueries.mockReturnValue(errorReturnValue);

      render(<ServicesCounselingEditShipmentDetails {...props} />);

      const errorMessage = await screen.getByText(/Something went wrong./);
      expect(errorMessage).toBeInTheDocument();
    });
  });

  it('renders the Services Counseling Shipment Form', async () => {
    useEditShipmentQueries.mockReturnValue(useEditShipmentQueriesReturnValue);
    render(<ServicesCounselingEditShipmentDetails {...props} />);

    const h1 = await screen.getByRole('heading', { name: 'Edit shipment details', level: 1 });
    await waitFor(() => {
      expect(h1).toBeInTheDocument();
    });
  });

  it('calls props.onUpdate with success and routes to move details when the save button is clicked and the shipment update is successful', async () => {
    updateMTOShipment.mockImplementation(() => Promise.resolve({}));
    const onUpdateMock = jest.fn();

    render(<ServicesCounselingEditShipmentDetails {...props} onUpdate={onUpdateMock} />);

    const saveButton = screen.getByRole('button', { name: 'Save' });

    expect(saveButton).not.toBeDisabled();

    await userEvent.click(saveButton);

    await waitFor(() => {
      expect(mockPush).toHaveBeenCalledWith('/counseling/moves/move123/details');
      expect(onUpdateMock).toHaveBeenCalledWith('success');
    });
  });

  it('calls props.onUpdate with error and routes to move details when the save button is clicked and the shipment update is unsuccessful', async () => {
    jest.spyOn(console, 'error').mockImplementation(() => {});
    updateMTOShipment.mockImplementation(() => Promise.reject(new Error('something went wrong')));
    useEditShipmentQueries.mockReturnValue(useEditShipmentQueriesReturnValue);
    const onUpdateMock = jest.fn();

    render(<ServicesCounselingEditShipmentDetails {...props} onUpdate={onUpdateMock} />);

    const saveButton = screen.getByRole('button', { name: 'Save' });

    expect(saveButton).not.toBeDisabled();

    await userEvent.click(saveButton);

    await waitFor(() => {
      expect(mockPush).toHaveBeenCalledWith('/counseling/moves/move123/details');
      expect(onUpdateMock).toHaveBeenCalledWith('error');
    });
  });

  it('routes to the move details page when the cancel button is clicked', async () => {
    useEditShipmentQueries.mockReturnValue(useEditShipmentQueriesReturnValue);
    render(<ServicesCounselingEditShipmentDetails {...props} />);

    const cancelButton = screen.getByRole('button', { name: 'Cancel' });

    expect(cancelButton).not.toBeDisabled();

    await userEvent.click(cancelButton);

    await waitFor(() => {
      expect(mockPush).toHaveBeenCalledWith('/counseling/moves/move123/details');
    });
  });

  describe('editing PPMs', () => {
    const ppmUseEditShipmentQueriesReturnValue = {
      ...useEditShipmentQueriesReturnValue,
      mtoShipments: [{ ...ppmShipment }],
    };

    it('renders the first page of the edit ppm Shipment Form with prefilled values', async () => {
      useEditShipmentQueries.mockReturnValue(ppmUseEditShipmentQueriesReturnValue);
      render(<ServicesCounselingEditShipmentDetails {...props} />);

      expect(await screen.findByTestId('tag')).toHaveTextContent('PPM');
      expect(await screen.getByRole('textbox', { name: 'Planned departure date' })).toHaveValue('28 Jun 2022');
      expect(await screen.findByRole('textbox', { name: 'Origin ZIP' })).toHaveValue(
        ppmShipment.ppmShipment.pickupPostalCode,
      );
      expect(await screen.findByRole('textbox', { name: 'Second origin ZIP' })).toHaveValue(
        ppmShipment.ppmShipment.secondaryPickupPostalCode,
      );
      expect(await screen.findByRole('textbox', { name: 'Destination ZIP' })).toHaveValue(
        ppmShipment.ppmShipment.destinationPostalCode,
      );
      expect(await screen.findByRole('textbox', { name: 'Second destination ZIP' })).toHaveValue(
        ppmShipment.ppmShipment.secondaryDestinationPostalCode,
      );
      expect(await screen.queryByRole('textbox', { name: 'Estimated SIT weight' })).not.toBeInTheDocument();
      expect(await screen.queryByRole('textbox', { name: 'Estimated storage start' })).not.toBeInTheDocument();
      expect(await screen.queryByRole('textbox', { name: 'Estimated storage end' })).not.toBeInTheDocument();
      expect(await screen.findByRole('textbox', { name: 'Estimated PPM weight' })).toHaveValue('1,111');
      expect(await screen.queryByRole('textbox', { name: 'Estimated pro-gear weight' })).not.toBeInTheDocument();
      expect(await screen.queryByRole('textbox', { name: 'Estimated spouse pro-gear weight' })).not.toBeInTheDocument();
      expect(await screen.findByRole('button', { name: 'Save and Continue' })).toBeInTheDocument();
    });

    describe('Check SIT field validations', () => {
      it.each([
        [
          'sitEstimatedWeight',
          {
            sitEstimatedWeight: '{Tab}',
            sitEstimatedEntryDate: '15 Jun 2022',
            sitEstimatedDepartureDate: '25 Jul 2022',
          },
          'Required',
        ],
        [
          'sitEstimatedEntryDate',
          { sitEstimatedWeight: '1234', sitEstimatedEntryDate: 'asdf', sitEstimatedDepartureDate: '25 Jul 2022' },
          'Enter a complete date in DD MMM YYYY format (day, month, year).',
        ],
        [
          'sitEstimatedDepartureDate',
          { sitEstimatedWeight: '1234', sitEstimatedEntryDate: '15 Jun 2022', sitEstimatedDepartureDate: 'asdf' },
          'Enter a complete date in DD MMM YYYY format (day, month, year).',
        ],
      ])('Verify invalid %s field shows validation error', async (field, data, expectedError) => {
        useEditShipmentQueries.mockReturnValue(ppmUseEditShipmentQueriesReturnValue);
        render(<ServicesCounselingEditShipmentDetails {...props} />);

        const sitExpected = document.getElementById('sitExpectedYes').parentElement;
        const sitExpectedYes = within(sitExpected).getByRole('radio', { name: 'Yes' });
        await userEvent.click(sitExpectedYes);

        // The test is dependent on the ordering of these three lines, and I'm not sure why.
        // If either of the estimated storage dates is entered last, the test that puts an invalid value
        // in that field will fail. But if the estimated SIT weight comes last, everything works fine.
        await userEvent.type(screen.getByLabelText('Estimated storage start'), data.sitEstimatedEntryDate);
        await userEvent.type(screen.getByLabelText('Estimated storage end'), data.sitEstimatedDepartureDate);
        await userEvent.type(screen.getByLabelText('Estimated SIT weight'), data.sitEstimatedWeight);
        await userEvent.tab();

        await waitFor(() => {
          const alerts = screen.getAllByRole('alert');
          expect(alerts).toHaveLength(1);
          expect(alerts[0]).toHaveTextContent(expectedError);
        });

        expect(screen.getByRole('button', { name: 'Save and Continue' })).toBeDisabled();
        expect(screen.getByRole('alert').nextElementSibling.firstElementChild).toHaveAttribute('name', field);
      });
    });

    it('Enables Save and Continue button when sit required fields are filled in', async () => {
      useEditShipmentQueries.mockReturnValue(ppmUseEditShipmentQueriesReturnValue);
      render(<ServicesCounselingEditShipmentDetails {...props} />);

      const sitExpected = document.getElementById('sitExpectedYes').parentElement;
      const sitExpectedYes = within(sitExpected).getByRole('radio', { name: 'Yes' });
      await userEvent.click(sitExpectedYes);
      await userEvent.type(screen.getByLabelText('Estimated SIT weight'), '1234');
      await userEvent.type(screen.getByLabelText('Estimated storage start'), '15 Jun 2022');
      await userEvent.type(screen.getByLabelText('Estimated storage end'), '25 Jun 2022');
      await userEvent.tab();

      await waitFor(() => {
        expect(screen.queryByRole('alert')).not.toBeInTheDocument();
        expect(screen.getByRole('button', { name: 'Save and Continue' })).not.toBeDisabled();
      });
    });

    it('calls props.onUpdate with success and routes to Advance page when the save button is clicked and the shipment update is successful', async () => {
      useEditShipmentQueries.mockReturnValue(ppmUseEditShipmentQueriesReturnValue);
      updateMTOShipment.mockImplementation(() => Promise.resolve({}));
      validatePostalCode.mockImplementation(() => Promise.resolve(false));
      const onUpdateMock = jest.fn();

      render(<ServicesCounselingEditShipmentDetails {...props} onUpdate={onUpdateMock} />);

      await waitFor(() => {
        expect(screen.getByLabelText('Estimated PPM weight')).toHaveValue('1,111');
      });

      const saveButton = screen.getByRole('button', { name: 'Save and Continue' });
      expect(saveButton).not.toBeDisabled();

      await userEvent.click(saveButton);
      await waitFor(() => {
        expect(mockPush).toHaveBeenCalledWith('/counseling/moves/move123/shipments/shipment123/advance');
        expect(onUpdateMock).toHaveBeenCalledWith('success');
      });
    });
  });
});

import React from 'react';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import ServicesCounselingMoveAllowances from 'pages/Office/ServicesCounselingMoveAllowances/ServicesCounselingMoveAllowances';
import { MockProviders } from 'testUtils';
import { useOrdersDocumentQueries } from 'hooks/queries';
import { ORDERS_PAY_GRADE_TYPE, ORDERS_TYPE } from 'constants/orders';
import { permissionTypes } from 'constants/permissions';
import { isBooleanFlagEnabled } from 'utils/featureFlags';
import { MOVE_STATUSES } from 'shared/constants';

const mockOriginDutyLocation = {
  address: {
    city: 'Des Moines',
    country: 'US',
    eTag: 'MjAyMC0wOS0xNFQxNzo0MTozOC42OTg1OTha',
    id: '2e26b066-aaca-4563-b284-d7f3f978fb3c',
    postalCode: '50309',
    state: 'IA',
    streetAddress1: '987 Other Avenue',
    streetAddress2: 'P.O. Box 1234',
    streetAddress3: 'c/o Another Person',
  },
  address_id: '2e26b066-aaca-4563-b284-d7f3f978fb3c',
  eTag: 'MjAyMC0wOS0xNFQxNzo0MTozOC43MDcxOTVa',
  id: 'a3ec2bdd-aa0a-434a-ba58-34c85f047704',
  name: 'XBc1KNi3pA',
};

const mockOconusOriginDutyLocation = {
  address: {
    city: 'Des Moines',
    country: 'US',
    eTag: 'MjAyMC0wOS0xNFQxNzo0MTozOC42OTg1OTha',
    id: '2e26b066-aaca-4563-b284-d7f3f978fb3c',
    postalCode: '99702',
    state: 'AK',
    streetAddress1: '987 Other Avenue',
    streetAddress2: 'P.O. Box 1234',
    streetAddress3: 'c/o Another Person',
    isOconus: true,
  },
  address_id: '2e26b066-aaca-4563-b284-d7f3f978fb3c',
  eTag: 'MjAyMC0wOS0xNFQxNzo0MTozOC43MDcxOTVa',
  id: 'a3ec2bdd-aa0a-434a-ba58-34c85f047704',
  name: 'XBc1KNi3pA',
};

const mockDestinationDutyLocation = {
  address: {
    city: 'Augusta',
    country: 'United States',
    eTag: 'MjAyMC0wOS0xNFQxNzo0MDo0OC44OTM3MDVa',
    id: '5ac95be8-0230-47ea-90b4-b0f6f60de364',
    postalCode: '30813',
    state: 'GA',
    streetAddress1: 'Fort Gordon',
  },
  address_id: '5ac95be8-0230-47ea-90b4-b0f6f60de364',
  eTag: 'MjAyMC0wOS0xNFQxNzo0MDo0OC44OTM3MDVa',
  id: '2d5ada83-e09a-47f8-8de6-83ec51694a86',
  name: 'Fort Gordon',
};

const mockOconusDestinationDutyLocation = {
  address: {
    city: 'Augusta',
    country: 'United States',
    eTag: 'MjAyMC0wOS0xNFQxNzo0MDo0OC44OTM3MDVa',
    id: '5ac95be8-0230-47ea-90b4-b0f6f60de364',
    postalCode: '99702',
    state: 'AK',
    streetAddress1: 'Fort Gordon',
    isOconus: true,
  },
  address_id: '5ac95be8-0230-47ea-90b4-b0f6f60de364',
  eTag: 'MjAyMC0wOS0xNFQxNzo0MDo0OC44OTM3MDVa',
  id: '2d5ada83-e09a-47f8-8de6-83ec51694a86',
  name: 'Fort Gordon',
};

jest.mock('hooks/queries', () => ({
  useOrdersDocumentQueries: jest.fn(),
}));

jest.mock('utils/featureFlags', () => ({
  ...jest.requireActual('utils/featureFlags'),
  isBooleanFlagEnabled: jest.fn().mockImplementation(() => Promise.resolve(false)),
}));

const useOrdersDocumentQueriesReturnValue = {
  orders: {
    1: {
      agency: 'ARMY',
      customerID: '6ac40a00-e762-4f5f-b08d-3ea72a8e4b63',
      date_issued: '2018-03-15',
      department_indicator: 'AIR_AND_SPACE_FORCE',
      destinationDutyLocation: mockDestinationDutyLocation,
      eTag: 'MjAyMC0wOS0xNFQxNzo0MTozOC43MTE0Nlo=',
      entitlement: {
        authorizedWeight: 5000,
        eTag: 'MjAyMC0wOS0xNFQxNzo0MTozOC42ODAwOVo=',
        id: '0dbc9029-dfc5-4368-bc6b-dfc95f5fe317',
        nonTemporaryStorage: true,
        privatelyOwnedVehicle: true,
        proGearWeight: 2000,
        proGearWeightSpouse: 500,
        gunSafeWeight: 400,
        requiredMedicalEquipmentWeight: 1000,
        organizationalClothingAndIndividualEquipment: true,
        storageInTransit: 2,
        totalDependents: 1,
        totalWeight: 5000,
        weightRestriction: 500,
        ubWeightRestriction: 400,
      },
      first_name: 'Leo',
      grade: 'E_1',
      id: '1',
      last_name: 'Spacemen',
      order_number: 'ORDER3',
      order_type: 'PERMANENT_CHANGE_OF_STATION',
      order_type_detail: 'HHG_PERMITTED',
      originDutyLocation: mockOriginDutyLocation,
      report_by_date: '2018-08-01',
      tac: 'F8E1',
      sac: 'E2P3',
    },
  },
  move: {
    status: MOVE_STATUSES.NEEDS_SERVICE_COUNSELING,
  },
};

const useCivilianTDYOrdersDocumentQueriesReturnValue = {
  orders: {
    1: {
      agency: 'ARMY',
      customerID: '6ac40a00-e762-4f5f-b08d-3ea72a8e4b63',
      date_issued: '2018-03-15',
      department_indicator: 'AIR_AND_SPACE_FORCE',
      destinationDutyLocation: mockOconusDestinationDutyLocation,
      eTag: 'MjAyMC0wOS0xNFQxNzo0MTozOC43MTE0Nlo=',
      entitlement: {
        authorizedWeight: 5000,
        eTag: 'MjAyMC0wOS0xNFQxNzo0MTozOC42ODAwOVo=',
        id: '0dbc9029-dfc5-4368-bc6b-dfc95f5fe317',
        nonTemporaryStorage: true,
        privatelyOwnedVehicle: true,
        proGearWeight: 2000,
        proGearWeightSpouse: 500,
        requiredMedicalEquipmentWeight: 1000,
        organizationalClothingAndIndividualEquipment: true,
        storageInTransit: 2,
        totalDependents: 1,
        totalWeight: 5000,
        weightRestriction: 500,
        unaccompaniedBaggageAllowance: 351,
        ubAllowance: 351,
      },
      first_name: 'Leo',
      grade: ORDERS_PAY_GRADE_TYPE.CIVILIAN_EMPLOYEE,
      id: '1',
      last_name: 'Spacemen',
      order_number: 'ORDER3',
      order_type: ORDERS_TYPE.TEMPORARY_DUTY,
      order_type_detail: 'HHG_PERMITTED',
      originDutyLocation: mockOconusOriginDutyLocation,
      report_by_date: '2018-08-01',
      tac: 'F8E1',
      sac: 'E2P3',
    },
  },
  move: {
    status: MOVE_STATUSES.NEEDS_SERVICE_COUNSELING,
  },
};
const editMoveStatuses = [
  MOVE_STATUSES.NEEDS_SERVICE_COUNSELING,
  MOVE_STATUSES.SERVICE_COUNSELING_COMPLETED,
  MOVE_STATUSES.APPROVALS_REQUESTED,
];

const disabledMoveStatuses = [
  MOVE_STATUSES.SUBMITTED,
  MOVE_STATUSES.APPROVED,
  MOVE_STATUSES.CANCELED,
  MOVE_STATUSES.APPROVALS_REQUESTED,
];

const loadingReturnValue = {
  isLoading: true,
  isError: false,
  isSuccess: false,
};

const errorReturnValue = {
  isLoading: false,
  isError: true,
  isSuccess: false,
};

describe('MoveAllowances page', () => {
  describe('check loading and error component states', () => {
    it('renders the Loading Placeholder when the query is still loading', async () => {
      useOrdersDocumentQueries.mockReturnValue(loadingReturnValue);

      render(
        <MockProviders>
          <ServicesCounselingMoveAllowances />
        </MockProviders>,
      );

      const h2 = await screen.getByRole('heading', { name: 'Loading, please wait...', level: 2 });
      expect(h2).toBeInTheDocument();
    });

    it('renders the Something Went Wrong component when the query errors', async () => {
      useOrdersDocumentQueries.mockReturnValue(errorReturnValue);

      render(
        <MockProviders>
          <ServicesCounselingMoveAllowances />
        </MockProviders>,
      );

      const errorMessage = await screen.getByText(/Something went wrong./);
      expect(errorMessage).toBeInTheDocument();
    });
  });

  describe('Basic rendering', () => {
    beforeEach(() => {
      useOrdersDocumentQueries.mockReturnValue(useOrdersDocumentQueriesReturnValue);
    });

    it('renders the sidebar elements', async () => {
      render(
        <MockProviders>
          <ServicesCounselingMoveAllowances />
        </MockProviders>,
      );

      expect(await screen.findByTestId('allowances-header')).toHaveTextContent('View allowances');
      expect(screen.getByTestId('view-orders')).toHaveTextContent('View orders');
      expect(screen.getByTestId('header')).toHaveTextContent('Counseling');
    });

    it('renders displays the allowances in the sidebar form', async () => {
      isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(true));
      await render(
        <MockProviders>
          <ServicesCounselingMoveAllowances />
        </MockProviders>,
      );

      expect(await screen.findByTestId('proGearWeightInput')).toHaveDisplayValue('2,000');
      expect(screen.getByTestId('proGearWeightSpouseInput')).toHaveDisplayValue('500');
      expect(screen.getByTestId('gunSafeWeightInput')).toHaveDisplayValue('400');
      expect(screen.getByTestId('rmeInput')).toHaveDisplayValue('1,000');
      expect(screen.getByTestId('branchInput')).toHaveDisplayValue('Army');
      expect(screen.getByTestId('sitInput')).toHaveDisplayValue('2');

      expect(screen.getByLabelText('OCIE authorized (Army only)')).toBeChecked();

      expect(screen.getByTestId('weightAllowance')).toHaveTextContent('5,000 lbs');
    });

    it('renders displays the allowances in the sidebar form and allows editing with correct permissions', async () => {
      render(
        <MockProviders permissions={[permissionTypes.updateAllowances]}>
          <ServicesCounselingMoveAllowances />
        </MockProviders>,
      );

      expect(await screen.findByTestId('proGearWeightInput')).toHaveDisplayValue('2,000');
      expect(screen.getByTestId('proGearWeightSpouseInput')).toHaveDisplayValue('500');
      expect(screen.getByTestId('rmeInput')).toHaveDisplayValue('1,000');
      expect(screen.getByTestId('branchInput')).toHaveDisplayValue('Army');
      expect(screen.getByTestId('sitInput')).toHaveDisplayValue('2');

      expect(screen.getByLabelText('OCIE authorized (Army only)')).toBeChecked();

      expect(screen.getByTestId('weightAllowance')).toHaveTextContent('5,000 lbs');
      // admin restricted weight location
      const adminWeightCheckbox = await screen.findByTestId('adminWeightLocation');
      expect(adminWeightCheckbox).toBeChecked();
      const weightRestrictionInput = screen.getByTestId('weightRestrictionInput');
      expect(weightRestrictionInput).toHaveValue('500');

      await userEvent.click(weightRestrictionInput);
      await userEvent.clear(weightRestrictionInput);
      await userEvent.type(weightRestrictionInput, '0');
      fireEvent.blur(weightRestrictionInput);

      await waitFor(() => {
        expect(screen.getByText(/Weight restriction must be greater than 0/i)).toBeInTheDocument();
      });

      await userEvent.clear(weightRestrictionInput);

      await waitFor(() => {
        expect(
          screen.getByText(/Weight restriction is required when Admin Restricted Weight Location is enabled/i),
        ).toBeInTheDocument();
      });
      // admin restricted UB weight location
      const adminUBWeightCheckbox = await screen.findByTestId('adminUBWeightLocation');
      expect(adminUBWeightCheckbox).toBeChecked();
      const ubWeightRestrictionInput = screen.getByTestId('ubWeightRestrictionInput');
      expect(ubWeightRestrictionInput).toHaveValue('400');

      await userEvent.click(ubWeightRestrictionInput);
      await userEvent.clear(ubWeightRestrictionInput);
      await userEvent.type(ubWeightRestrictionInput, '0');
      fireEvent.blur(ubWeightRestrictionInput);

      await waitFor(() => {
        expect(screen.getByText(/UB Weight restriction must be greater than 0/i)).toBeInTheDocument();
      });

      await userEvent.clear(ubWeightRestrictionInput);

      await waitFor(() => {
        expect(
          screen.getByText(/UB weight restriction is required when Admin Restricted UB Weight Location is enabled/i),
        ).toBeInTheDocument();
      });
    });

    it('does not render the civilian TDY UB allowance editable field in the sidebar form if not a civilian TDY move', async () => {
      isBooleanFlagEnabled.mockResolvedValue(false);
      useOrdersDocumentQueries.mockReturnValue(useOrdersDocumentQueriesReturnValue);

      render(
        <MockProviders permissions={[permissionTypes.updateAllowances]}>
          <ServicesCounselingMoveAllowances />
        </MockProviders>,
      );

      await waitFor(() => {
        expect(
          screen.queryByLabelText(/If the customer's orders specify a UB weight allowance, enter it here./i),
        ).not.toBeInTheDocument();
      });
    });

    it('renders the civilian TDY UB allowance editable field in the sidebar form if civilian TDY move', async () => {
      isBooleanFlagEnabled.mockResolvedValue(true);
      useOrdersDocumentQueries.mockReturnValue(useCivilianTDYOrdersDocumentQueriesReturnValue);

      render(
        <MockProviders permissions={[permissionTypes.updateAllowances]}>
          <ServicesCounselingMoveAllowances />
        </MockProviders>,
      );

      await waitFor(() => {
        expect(
          screen.queryByLabelText(/If the customer's orders specify a UB weight allowance, enter it here./i),
        ).toBeInTheDocument();
      });
    });
  });

  describe('Conditional disabling', () => {
    beforeEach(() => {
      jest.clearAllMocks();
    });

    it('renders and disables editing with correct statuses', async () => {
      for (let i = 0; i < disabledMoveStatuses.length; i += 1) {
        const orderQueryReturnValues = JSON.parse(JSON.stringify(useOrdersDocumentQueriesReturnValue));
        orderQueryReturnValues.move = {
          id: 123,
          moveCode: 'GLOBAL123',
          ordersId: 1,
          status: disabledMoveStatuses[i],
        };
        useOrdersDocumentQueries.mockReturnValue(orderQueryReturnValues);

        render(
          <MockProviders permissions={[permissionTypes.updateAllowances]}>
            <ServicesCounselingMoveAllowances />
          </MockProviders>,
        );

        const proGearWeightInput = screen.getAllByTestId('proGearWeightInput');
        expect(proGearWeightInput[0]).toBeInTheDocument();
        expect(proGearWeightInput[0]).toBeDisabled();
        const proGearWeightSpouseInput = screen.getAllByTestId('proGearWeightSpouseInput');
        expect(proGearWeightSpouseInput[0]).toBeInTheDocument();
        expect(proGearWeightSpouseInput[0]).toBeDisabled();
        const rmeInput = screen.getAllByTestId('rmeInput');
        expect(rmeInput[0]).toBeInTheDocument();
        expect(rmeInput[0]).toBeDisabled();
        const sitInput = screen.getAllByTestId('sitInput');
        expect(sitInput[0]).toBeInTheDocument();
        expect(sitInput[0]).toBeDisabled();
      }
    });

    it('renders and allows editing with correct statuses', async () => {
      for (let i = 0; i < editMoveStatuses.length; i += 1) {
        const orderQueryReturnValues = JSON.parse(JSON.stringify(useOrdersDocumentQueriesReturnValue));
        orderQueryReturnValues.move = {
          id: 123,
          moveCode: 'GLOBAL123',
          ordersId: 1,
          status: editMoveStatuses[i],
        };
        useOrdersDocumentQueries.mockReturnValue(orderQueryReturnValues);

        render(
          <MockProviders permissions={[permissionTypes.updateAllowances]}>
            <ServicesCounselingMoveAllowances />
          </MockProviders>,
        );

        const proGearWeightInput = screen.getAllByTestId('proGearWeightInput');
        expect(proGearWeightInput[0]).toBeInTheDocument();
        expect(proGearWeightInput[0]).not.toBeDisabled();
        const proGearWeightSpouseInput = screen.getAllByTestId('proGearWeightSpouseInput');
        expect(proGearWeightSpouseInput[0]).toBeInTheDocument();
        expect(proGearWeightSpouseInput[0]).not.toBeDisabled();
        const rmeInput = screen.getAllByTestId('rmeInput');
        expect(rmeInput[0]).toBeInTheDocument();
        expect(rmeInput[0]).not.toBeDisabled();
        const sitInput = screen.getAllByTestId('sitInput');
        expect(sitInput[0]).toBeInTheDocument();
        expect(sitInput[0]).not.toBeDisabled();
      }
    });
  });
});

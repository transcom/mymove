/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { render, screen, waitFor, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import Orders from './Orders';

import { isBooleanFlagEnabled } from 'utils/featureFlags';
import { MockProviders } from 'testUtils';
import { useOrdersDocumentQueries } from 'hooks/queries';
import { permissionTypes } from 'constants/permissions';
import { MOVE_DOCUMENT_TYPE } from 'shared/constants';
import { ORDERS_TYPE, ORDERS_PAY_GRADE_TYPE } from 'constants/orders';

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

const mockLoa = {
  createdAt: '2023-08-03T19:17:10.050Z',
  id: '06254fc3-b763-484c-b555-42855d1ad5cd',
  loaAlltSnID: '123A',
  loaBafID: '1234',
  loaBdgtAcntClsNm: '000000',
  loaBgFyTx: 2006,
  loaBgnDt: '2005-10-01',
  loaDocID: 'HHG12345678900',
  loaDptID: '1',
  loaDscTx: 'PERSONAL PROPERTY - PARANORMAL ACTIVITY DIVISION (OTHER)',
  loaEndDt: '2015-10-01',
  loaEndFyTx: 2016,
  loaHsGdsCd: 'HT',
  loaInstlAcntgActID: '12345',
  loaObjClsID: '22NL',
  loaOpAgncyID: '1A',
  loaPgmElmntID: '00000000',
  loaStatCd: 'U',
  loaSysId: '10003',
  loaTrnsnID: 'B1',
  loaTrsySfxTx: '0000',
  orgGrpDfasCd: 'ZZ',
  updatedAt: '2023-08-03T19:17:38.776Z',
  validHhgProgramCodeForLoa: true,
  validLoaForTac: true,
};

jest.mock('hooks/queries', () => ({
  useOrdersDocumentQueries: jest.fn(),
}));

jest.mock('services/ghcApi', () => ({
  ...jest.requireActual('services/ghcApi'),
  getTacValid: ({ tac }) => {
    return {
      tac,
      isValid: tac === '1111' || tac === '2222',
    };
  },
  getLoa: ({ tacCode }) => {
    // 1111 is our good dummy TAC code
    if (tacCode === '1111') {
      // 200 OK, a LOA was found
      return Promise.resolve(mockLoa);
    }
    if (tacCode === '2222') {
      // 200 OK, but no LOAs were found
      return Promise.resolve(undefined);
    }
    if (tacCode === '3333') {
      // 200 OK, but the LOA found is invalid
      const invalidLoa = { ...mockLoa, validHhgProgramCodeForLoa: false };
      return Promise.resolve(invalidLoa);
    }
    // Default to no LOA
    return Promise.resolve(undefined);
  },
  getPayGradeOptions: jest.fn().mockImplementation(() => {
    const E_1 = 'E-1';
    const E_6 = 'E-6';
    const CIVILIAN_EMPLOYEE = 'CIVILIAN_EMPLOYEE';

    return Promise.resolve({
      body: [
        {
          grade: E_1,
          description: E_1,
        },
        {
          grade: E_6,
          description: E_6,
        },
        {
          description: CIVILIAN_EMPLOYEE,
          grade: CIVILIAN_EMPLOYEE,
        },
      ],
    });
  }),
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
        dependentsAuthorized: true,
        eTag: 'MjAyMC0wOS0xNFQxNzo0MTozOC42ODAwOVo=',
        id: '0dbc9029-dfc5-4368-bc6b-dfc95f5fe317',
        nonTemporaryStorage: true,
        privatelyOwnedVehicle: true,
        proGearWeight: 2000,
        proGearWeightSpouse: 500,
        storageInTransit: 2,
        totalDependents: 1,
        totalWeight: 5000,
      },
      first_name: 'Leo',
      grade: ORDERS_PAY_GRADE_TYPE.E_1,
      id: '1',
      last_name: 'Spacemen',
      order_number: 'ORDER3',
      order_type: 'PERMANENT_CHANGE_OF_STATION',
      order_type_detail: 'HHG_PERMITTED',
      originDutyLocation: mockOriginDutyLocation,
      report_by_date: '2018-08-01',
      tac: 'F8E1',
      sac: 'E2P3',
      ntsTac: '1111',
      ntsSac: '2222',
    },
  },
  move: {
    approved_at: '2018-03-15',
  },
};
const ordersMockProps = {
  files: {
    [MOVE_DOCUMENT_TYPE.ORDERS]: [{ id: 'file-1', name: 'Order File 1' }],
    [MOVE_DOCUMENT_TYPE.AMENDMENTS]: [{ id: 'file-2', name: 'Amended File 1' }],
  },
};

const loadingReturnValue = {
  ...useOrdersDocumentQueriesReturnValue,
  isLoading: true,
  isError: false,
  isSuccess: false,
};

const errorReturnValue = {
  ...useOrdersDocumentQueriesReturnValue,
  isLoading: false,
  isError: true,
  isSuccess: false,
};

describe('Orders page', () => {
  describe('check loading and error component states', () => {
    it('renders the Loading Placeholder when the query is still loading', async () => {
      useOrdersDocumentQueries.mockReturnValueOnce(loadingReturnValue);

      render(
        <MockProviders>
          <Orders {...ordersMockProps} />
        </MockProviders>,
      );

      const h2 = screen.getByRole('heading', { name: 'Loading, please wait...', level: 2 });
      expect(h2).toBeInTheDocument();
    });

    it('renders the Something Went Wrong component when the query errors', async () => {
      isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(true));
      useOrdersDocumentQueries.mockReturnValueOnce(errorReturnValue);

      render(
        <MockProviders>
          <Orders {...ordersMockProps} />
        </MockProviders>,
      );

      const errorMessage = screen.getByText(/Something went wrong./);
      expect(errorMessage).toBeInTheDocument();
    });
  });

  describe('Basic rendering', () => {
    it('renders the sidebar orders detail form', async () => {
      useOrdersDocumentQueries.mockReturnValue(useOrdersDocumentQueriesReturnValue);

      render(
        <MockProviders>
          <Orders {...ordersMockProps} />
        </MockProviders>,
      );

      expect(await screen.findByLabelText('Current duty location *')).toBeInTheDocument();
      expect(screen.getByTestId('ntsTacInput')).toHaveValue('1111');
      expect(screen.getByTestId('ntsSacInput')).toHaveValue('2222');
      expect(screen.getByTestId('payGradeInput')).toHaveDisplayValue(ORDERS_PAY_GRADE_TYPE.E_1);
      expect(screen.getByLabelText('Dependents authorized')).toBeChecked();
    });
  });

  describe('TAC validation', () => {
    it('validates on load', async () => {
      useOrdersDocumentQueries.mockReturnValue(useOrdersDocumentQueriesReturnValue);

      render(
        <MockProviders>
          <Orders {...ordersMockProps} />
        </MockProviders>,
      );

      expect(await screen.findByText(/This TAC does not appear in TGET/)).toBeInTheDocument();
    });

    it('validates on user input', async () => {
      useOrdersDocumentQueries.mockReturnValue(useOrdersDocumentQueriesReturnValue);

      render(
        <MockProviders permissions={[permissionTypes.updateOrders]}>
          <Orders {...ordersMockProps} />
        </MockProviders>,
      );

      const hhgTacInput = screen.getByTestId('hhgTacInput');
      await userEvent.clear(hhgTacInput);
      await userEvent.type(hhgTacInput, '2222');

      await waitFor(() => {
        expect(screen.queryByText(/This TAC does not appear in TGET/)).not.toBeInTheDocument();
      });

      await userEvent.clear(hhgTacInput);
      await userEvent.type(hhgTacInput, '3333');

      await waitFor(() => {
        expect(screen.getByText(/This TAC does not appear in TGET/)).toBeInTheDocument();
      });
    });

    it('validates TAC', async () => {
      useOrdersDocumentQueries.mockReturnValue(useOrdersDocumentQueriesReturnValue);

      render(
        <MockProviders permissions={[permissionTypes.updateOrders]}>
          <Orders {...ordersMockProps} />
        </MockProviders>,
      );

      const hhgTacInput = screen.getByTestId('hhgTacInput');
      await userEvent.clear(hhgTacInput);
      await userEvent.type(hhgTacInput, '****');
      await waitFor(() => {
        // no *
        expect(screen.getByText('TAC cannot contain * or " characters')).toBeInTheDocument();
      });

      await userEvent.clear(hhgTacInput);
      await userEvent.type(hhgTacInput, '""""');
      await waitFor(() => {
        // no "
        expect(screen.getByText('TAC cannot contain * or " characters')).toBeInTheDocument();
      });

      // NTS TAC
      const ntsTacInput = screen.getByTestId('ntsTacInput');
      await userEvent.clear(ntsTacInput);
      await userEvent.type(ntsTacInput, '****');
      await waitFor(() => {
        expect(screen.getByText('TAC cannot contain * or " characters')).toBeInTheDocument();
      });

      await userEvent.clear(ntsTacInput);
      await userEvent.type(ntsTacInput, '""""');
      await waitFor(() => {
        expect(screen.getByText('TAC cannot contain * or " characters')).toBeInTheDocument();
      });
    });

    it('validates SAC', async () => {
      useOrdersDocumentQueries.mockReturnValue(useOrdersDocumentQueriesReturnValue);

      render(
        <MockProviders permissions={[permissionTypes.updateOrders]}>
          <Orders {...ordersMockProps} />
        </MockProviders>,
      );

      // SAC
      const hhgSacInput = screen.getByTestId('hhgSacInput');
      await userEvent.clear(hhgSacInput);
      await userEvent.type(hhgSacInput, '****');
      hhgSacInput.blur();
      await waitFor(() => {
        // no *
        expect(screen.getByText('SAC cannot contain * or " characters')).toBeInTheDocument();
      });

      await userEvent.clear(hhgSacInput);
      await userEvent.type(hhgSacInput, '""""');
      await waitFor(() => {
        // no "
        expect(screen.getByText('SAC cannot contain * or " characters')).toBeInTheDocument();
      });

      // NTS SAC
      const ntsSacInput = screen.getByTestId('ntsSacInput');
      await userEvent.clear(ntsSacInput);
      await userEvent.type(ntsSacInput, '****');
      ntsSacInput.blur();
      await waitFor(() => {
        expect(screen.getByText('NTS SAC cannot contain * or " characters')).toBeInTheDocument();
      });

      await userEvent.clear(ntsSacInput);
      await userEvent.type(ntsSacInput, '""""');
      await waitFor(() => {
        expect(screen.getByText('NTS SAC cannot contain * or " characters')).toBeInTheDocument();
      });
    });

    it('SAC fields can be more than 4 digits', async () => {
      useOrdersDocumentQueries.mockReturnValue(useOrdersDocumentQueriesReturnValue);

      render(
        <MockProviders permissions={[permissionTypes.updateOrders]}>
          <Orders {...ordersMockProps} />
        </MockProviders>,
      );

      // SAC
      const hhgSacInput = screen.getByTestId('hhgSacInput');
      await userEvent.type(hhgSacInput, 'MoreThan4Digits');
      expect(hhgSacInput).toHaveValue('E2P3MoreThan4Digits');

      // NTS SAC
      const ntsSacInput = screen.getByTestId('ntsSacInput');
      await userEvent.type(ntsSacInput, '4DigitsOrMore');
      expect(ntsSacInput).toHaveValue('22224DigitsOrMore');
    });
  });

  describe('LOA validation', () => {
    beforeEach(() => {
      useOrdersDocumentQueries.mockReturnValue(useOrdersDocumentQueriesReturnValue);

      render(
        <MockProviders permissions={[permissionTypes.updateOrders]}>
          <Orders {...ordersMockProps} />
        </MockProviders>,
      );
    });

    it('validates on load', async () => {
      expect(await screen.findByText(/Unable to find a LOA based on the provided details/)).toBeInTheDocument();
    });

    describe('validates on user input', () => {
      it('validates HHG with a valid TAC and no LOA', async () => {
        const hhgTacInput = screen.getByTestId('hhgTacInput');
        await userEvent.clear(hhgTacInput);
        await userEvent.type(hhgTacInput, '2222');

        // TAC is found and valid
        // LOA is NOT found
        await waitFor(() => {
          expect(screen.queryByText(/This TAC does not appear in TGET/)).not.toBeInTheDocument();
          expect(screen.getByText(/Unable to find a LOA based on the provided details/)).toBeInTheDocument();
          expect(
            screen.queryByText(/The LOA identified based on the provided details appears to be invalid/),
          ).not.toBeInTheDocument();
        });
      });
      it('validates NTS with a valid TAC and no LOA', async () => {
        // Empty HHG from having a good useEffect TAC
        const hhgTacInput = screen.getByTestId('hhgTacInput');
        await userEvent.clear(hhgTacInput);
        const ntsTacInput = screen.getByTestId('ntsTacInput');
        await userEvent.clear(ntsTacInput);
        await userEvent.type(ntsTacInput, '2222');

        // TAC is found and valid
        // LOA is NOT found
        await waitFor(() => {
          const loaMissingWarnings = screen.queryAllByText(/Unable to find a LOA based on the provided details/);
          expect(screen.queryByText(/This TAC does not appear in TGET/)).not.toBeInTheDocument(); // TAC should be good
          expect(loaMissingWarnings.length).toBe(2); // Both HHG and NTS LOAs are missing now
          expect(
            screen.queryByText(/The LOA identified based on the provided details appears to be invalid/),
          ).not.toBeInTheDocument();
        });

        // Make HHG good and re-verify that the NTS errors remained
        await userEvent.type(hhgTacInput, '1111');
        await waitFor(() => {
          const loaMissingWarnings = screen.queryAllByText(/Unable to find a LOA based on the provided details/);
          expect(screen.queryByText(/This TAC does not appear in TGET/)).not.toBeInTheDocument(); // TAC should be good
          expect(loaMissingWarnings.length).toBe(1); // Only NTS is missing
          expect(
            screen.queryByText(/The LOA identified based on the provided details appears to be invalid/),
          ).not.toBeInTheDocument();
        });
      });
      it('validates an invalid HHG LOA', async () => {
        const hhgTacInput = screen.getByTestId('hhgTacInput');
        await userEvent.clear(hhgTacInput);
        await userEvent.type(hhgTacInput, '3333');

        // TAC is found and valid
        // LOA is found and NOT valid
        await waitFor(() => {
          const loaInvalidWarnings = screen.queryAllByText(
            /The LOA identified based on the provided details appears to be invalid/,
          );
          const loaMissingWarnings = screen.queryAllByText(/Unable to find a LOA based on the provided details/);
          expect(loaInvalidWarnings.length).toBe(1); // HHG is invalid
          expect(loaMissingWarnings.length).toBe(0); // NTS is valid based on useEffect hook and default passed in TAC
        });
      });
      it('validates an invalid NTS LOA', async () => {
        const ntsTacInput = screen.getByTestId('ntsTacInput');
        await userEvent.clear(ntsTacInput);
        await userEvent.type(ntsTacInput, '3333');

        // TAC is found and valid
        // LOA is found and NOT valid
        await waitFor(() => {
          const loaInvalidWarnings = screen.queryAllByText(
            /The LOA identified based on the provided details appears to be invalid/,
          );
          const loaMissingWarnings = screen.queryAllByText(/Unable to find a LOA based on the provided details/);
          expect(loaInvalidWarnings.length).toBe(1); // NTS is invalid
          expect(loaMissingWarnings.length).toBe(1); // HHG is valid based on useEffect hook and default passed in TAC
        });
      });
    });
  });

  describe('LOA concatenation', () => {
    it('concatenates the LOA string correctly', async () => {
      useOrdersDocumentQueries.mockReturnValue(useOrdersDocumentQueriesReturnValue);

      render(
        <MockProviders permissions={[permissionTypes.updateOrders]}>
          <Orders {...ordersMockProps} />
        </MockProviders>,
      );

      const hhgTacInput = screen.getByTestId('hhgTacInput');
      await userEvent.clear(hhgTacInput);
      await userEvent.type(hhgTacInput, '1111');

      const expectedLongLineOfAccounting =
        '1**20062016*1234*0000**1A*123A**00000000*********22NL***000000*HHG12345678900**12345**B1*';

      const loaTextField = screen.getByTestId('hhgLoaTextField');
      expect(loaTextField).toHaveValue(expectedLongLineOfAccounting);
    });
  });
  describe('LOA concatenation with regex removes extra spaces', () => {
    it('concatenates the LOA string correctly and without extra spaces', async () => {
      let extraSpacesLongLineOfAccounting =
        '1  **20062016*1234 *0000**1A *123A**00000000**  **** ***22NL** *000000*SEE PCS ORDERS* *12345**B1*';
      const expectedLongLineOfAccounting =
        '1**20062016*1234*0000**1A*123A**00000000*********22NL***000000*SEE PCS ORDERS**12345**B1*';

      // preserves spaces in column values such as 'SEE PCS ORDERS'
      // remove any number of spaces following an asterisk in a LOA string
      extraSpacesLongLineOfAccounting = extraSpacesLongLineOfAccounting.replace(/\* +/g, '*');
      // remove any number of spaces preceding an asterisk in a LOA string
      extraSpacesLongLineOfAccounting = extraSpacesLongLineOfAccounting.replace(/ +\*/g, '*');

      expect(extraSpacesLongLineOfAccounting).toEqual(expectedLongLineOfAccounting);
    });
  });
  describe('Manage document permission', () => {
    it('renders manage document component', async () => {
      useOrdersDocumentQueries.mockReturnValue(useOrdersDocumentQueriesReturnValue);

      render(
        <MockProviders permissions={[permissionTypes.updateOrders]}>
          <Orders {...ordersMockProps} />
        </MockProviders>,
      );

      await waitFor(() => {
        expect(screen.queryByText(/Manage Orders/)).toBeInTheDocument();
        expect(screen.queryByText(/Manage Amended Orders/)).toBeInTheDocument();
      });
    });
    it('does not render manage document component', async () => {
      useOrdersDocumentQueries.mockReturnValue(useOrdersDocumentQueriesReturnValue);

      render(
        <MockProviders>
          <Orders {...ordersMockProps} />
        </MockProviders>,
      );

      await waitFor(() => {
        expect(screen.queryByText(/Manage Orders/)).not.toBeInTheDocument();
        expect(screen.queryByText(/Manage Amended Orders/)).not.toBeInTheDocument();
      });
    });
  });

  describe('wounded warrior FF', () => {
    beforeEach(() => {
      jest.resetAllMocks();
    });

    it('wounded warrior FF turned off', async () => {
      isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(false));
      useOrdersDocumentQueries.mockReturnValue(useOrdersDocumentQueriesReturnValue);

      render(
        <MockProviders>
          <Orders {...ordersMockProps} />
        </MockProviders>,
      );

      await waitFor(() => {
        const ordersTypeDropdown = screen.getByLabelText('Orders type *');
        const options = within(ordersTypeDropdown).queryAllByRole('option');
        const hasWoundedWarrior = options.some((option) => option.value === ORDERS_TYPE.WOUNDED_WARRIOR);
        expect(hasWoundedWarrior).toBe(false);
      });
    });

    it('wounded warrior FF turned on', async () => {
      isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(true));
      useOrdersDocumentQueries.mockReturnValue(useOrdersDocumentQueriesReturnValue);

      render(
        <MockProviders>
          <Orders {...ordersMockProps} />
        </MockProviders>,
      );

      await waitFor(() => {
        const ordersTypeDropdown = screen.getByLabelText('Orders type *');
        const options = within(ordersTypeDropdown).queryAllByRole('option');
        const hasWoundedWarrior = options.some((option) => option.value === ORDERS_TYPE.WOUNDED_WARRIOR);
        expect(hasWoundedWarrior).toBe(true);
      });
    });
  });
});

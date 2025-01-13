/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import ServicesCounselingOrders from 'pages/Office/ServicesCounselingOrders/ServicesCounselingOrders';
import { MockProviders } from 'testUtils';
import { useOrdersDocumentQueries } from 'hooks/queries';
import { MOVE_DOCUMENT_TYPE } from 'shared/constants';
import { counselingUpdateOrder, getOrder } from 'services/ghcApi';
import { formatYesNoAPIValue } from 'utils/formatters';
import { ORDERS_TYPE } from 'constants/orders';

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
      isValid: tac === '1111' || tac === '2222' || tac === '3333',
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
  counselingUpdateOrder: jest.fn(),
  getOrder: jest.fn(),
}));

jest.mock('utils/featureFlags', () => ({
  ...jest.requireActual('utils/featureFlags'),
  isBooleanFlagEnabled: jest.fn().mockImplementation(() => Promise.resolve(true)),
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
      ntsTac: '1111',
      ntsSac: 'R6X1',
    },
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
      useOrdersDocumentQueries.mockReturnValue(loadingReturnValue);

      render(
        <MockProviders>
          <ServicesCounselingOrders {...ordersMockProps} />
        </MockProviders>,
      );

      const h2 = await screen.getByRole('heading', { name: 'Loading, please wait...', level: 2 });
      expect(h2).toBeInTheDocument();
    });

    it('renders the Something Went Wrong component when the query errors', async () => {
      useOrdersDocumentQueries.mockReturnValue(errorReturnValue);

      render(
        <MockProviders>
          <ServicesCounselingOrders {...ordersMockProps} />
        </MockProviders>,
      );

      const errorMessage = await screen.getByText(/Something went wrong./);
      expect(errorMessage).toBeInTheDocument();
    });
  });

  describe('Basic rendering', () => {
    it('renders the sidebar orders detail form', async () => {
      useOrdersDocumentQueries.mockReturnValue(useOrdersDocumentQueriesReturnValue);

      render(
        <MockProviders>
          <ServicesCounselingOrders {...ordersMockProps} />
        </MockProviders>,
      );

      expect(await screen.findByLabelText('Current duty location')).toBeInTheDocument();
    });

    it('renders the sidebar elements', async () => {
      useOrdersDocumentQueries.mockReturnValue(useOrdersDocumentQueriesReturnValue);

      render(
        <MockProviders>
          <ServicesCounselingOrders {...ordersMockProps} />
        </MockProviders>,
      );

      expect(await screen.findByTestId('view-orders-header')).toHaveTextContent('View orders');
      expect(screen.getByTestId('view-allowances')).toHaveTextContent('View allowances');
    });

    it('renders each option for orders type dropdown', async () => {
      useOrdersDocumentQueries.mockReturnValue(useOrdersDocumentQueriesReturnValue);

      render(
        <MockProviders>
          <ServicesCounselingOrders {...ordersMockProps} />
        </MockProviders>,
      );

      const ordersTypeDropdown = screen.getByLabelText('Orders type');
      expect(ordersTypeDropdown).toBeInstanceOf(HTMLSelectElement);

      await userEvent.selectOptions(ordersTypeDropdown, 'PERMANENT_CHANGE_OF_STATION');
      expect(ordersTypeDropdown).toHaveValue('PERMANENT_CHANGE_OF_STATION');

      await userEvent.selectOptions(ordersTypeDropdown, 'LOCAL_MOVE');
      expect(ordersTypeDropdown).toHaveValue('LOCAL_MOVE');

      await userEvent.selectOptions(ordersTypeDropdown, 'RETIREMENT');
      expect(ordersTypeDropdown).toHaveValue('RETIREMENT');

      await userEvent.selectOptions(ordersTypeDropdown, 'SEPARATION');
      expect(ordersTypeDropdown).toHaveValue('SEPARATION');

      await userEvent.selectOptions(ordersTypeDropdown, ORDERS_TYPE.EARLY_RETURN_OF_DEPENDENTS);
      expect(ordersTypeDropdown).toHaveValue(ORDERS_TYPE.EARLY_RETURN_OF_DEPENDENTS);

      await userEvent.selectOptions(ordersTypeDropdown, ORDERS_TYPE.STUDENT_TRAVEL);
      expect(ordersTypeDropdown).toHaveValue(ORDERS_TYPE.STUDENT_TRAVEL);
    });

    it('populates initial field values', async () => {
      useOrdersDocumentQueries.mockReturnValue(useOrdersDocumentQueriesReturnValue);

      render(
        <MockProviders>
          <ServicesCounselingOrders {...ordersMockProps} />
        </MockProviders>,
      );

      expect(await screen.findByText(mockOriginDutyLocation.name)).toBeInTheDocument();
      expect(screen.getByText(mockDestinationDutyLocation.name)).toBeInTheDocument();
      expect(screen.getByLabelText('Orders type')).toHaveValue('PERMANENT_CHANGE_OF_STATION');
      expect(screen.getByTestId('hhgTacInput')).toHaveValue('F8E1');
      expect(screen.getByTestId('hhgSacInput')).toHaveValue('E2P3');
      expect(screen.getByTestId('ntsTacInput')).toHaveValue('1111');
      expect(screen.getByTestId('ntsSacInput')).toHaveValue('R6X1');
      expect(screen.getByTestId('payGradeInput')).toHaveValue('E_1');
    });
  });

  it('renders an upload orders button', async () => {
    render(
      <MockProviders>
        <ServicesCounselingOrders {...ordersMockProps} />
      </MockProviders>,
    );

    expect(await screen.findByText('Manage Orders')).toBeInTheDocument();
    expect(await screen.findByText('Manage Amended Orders')).toBeInTheDocument();
  });

  describe('TAC validation', () => {
    it('validates on load', async () => {
      useOrdersDocumentQueries.mockReturnValue(useOrdersDocumentQueriesReturnValue);

      render(
        <MockProviders>
          <ServicesCounselingOrders {...ordersMockProps} />
        </MockProviders>,
      );

      expect(await screen.findByText(/This TAC does not appear in TGET/)).toBeInTheDocument();
    });

    it('validates on user input', async () => {
      useOrdersDocumentQueries.mockReturnValue(useOrdersDocumentQueriesReturnValue);

      render(
        <MockProviders>
          <ServicesCounselingOrders {...ordersMockProps} />
        </MockProviders>,
      );

      const hhgTacInput = screen.getByTestId('hhgTacInput');
      await userEvent.clear(hhgTacInput);
      await userEvent.type(hhgTacInput, '2222');

      await waitFor(() => {
        expect(screen.queryByText(/This TAC does not appear in TGET/)).not.toBeInTheDocument();
      });

      await userEvent.clear(hhgTacInput);
      await userEvent.type(hhgTacInput, '4444');

      await waitFor(() => {
        expect(screen.getByText(/This TAC does not appear in TGET/)).toBeInTheDocument();
      });
    });
  });

  describe('LOA validation', () => {
    beforeEach(() => {
      useOrdersDocumentQueries.mockReturnValue(useOrdersDocumentQueriesReturnValue);

      render(
        <MockProviders>
          <ServicesCounselingOrders {...ordersMockProps} />
        </MockProviders>,
      );
    });

    it('validates on load', async () => {
      // Both TAC and LOA are missing on load (On this test per useOrdersDocumentQueriesReturnValue and the
      // mocked responses)
      expect(await screen.getByText(/This TAC does not appear in TGET/)).toBeInTheDocument();
      expect(await screen.getByText(/Unable to find a LOA based on the provided details/)).toBeInTheDocument();
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
        <MockProviders>
          <ServicesCounselingOrders {...ordersMockProps} />
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

  describe('Order type: STUDENT_TRAVEL', () => {
    beforeEach(() => {
      jest.clearAllMocks();
    });

    it('select STUDENT_TRAVEL', async () => {
      // create a local copy of order return value and set initial values
      const orderQueryReturnValues = JSON.parse(JSON.stringify(useOrdersDocumentQueriesReturnValue));
      orderQueryReturnValues.move = { id: 123, moveCode: 'GLOBAL123', ordersId: 1 };
      orderQueryReturnValues.orders[1].order_type = ORDERS_TYPE.PERMANENT_CHANGE_OF_STATION;
      orderQueryReturnValues.orders[1].has_dependents = formatYesNoAPIValue('no');

      // set return values for mocked functions
      useOrdersDocumentQueries.mockReturnValue(orderQueryReturnValues);
      getOrder.mockResolvedValue(orderQueryReturnValues);

      // render component
      render(
        <MockProviders>
          <ServicesCounselingOrders {...ordersMockProps} />
        </MockProviders>,
      );

      // Select STUDENT_TRAVEL from the dropdown
      const ordersTypeDropdown = await screen.findByLabelText('Orders type');
      await userEvent.selectOptions(ordersTypeDropdown, ORDERS_TYPE.STUDENT_TRAVEL);

      // Submit the form
      const saveButton = screen.getByRole('button', { name: 'Save' });
      await userEvent.click(saveButton);

      // Verify correct values were passed
      await waitFor(() => {
        expect(counselingUpdateOrder).toHaveBeenCalledWith(
          expect.objectContaining({
            body: expect.objectContaining({
              hasDependents: formatYesNoAPIValue('yes'),
              ordersType: ORDERS_TYPE.STUDENT_TRAVEL,
            }),
          }),
        );
      });
    });

    it('De-select STUDENT_TRAVEL', async () => {
      // create a local copy of order return value and set initial values
      const orderQueryReturnValues = JSON.parse(JSON.stringify(useOrdersDocumentQueriesReturnValue));
      orderQueryReturnValues.move = { id: 123, moveCode: 'GLOBAL123', ordersId: 1 };
      orderQueryReturnValues.orders[1].order_type = ORDERS_TYPE.STUDENT_TRAVEL;
      orderQueryReturnValues.orders[1].has_dependents = formatYesNoAPIValue('yes');

      // set return values for mocked functions
      useOrdersDocumentQueries.mockReturnValue(orderQueryReturnValues);
      getOrder.mockResolvedValue(orderQueryReturnValues);

      // render component
      render(
        <MockProviders>
          <ServicesCounselingOrders {...ordersMockProps} />
        </MockProviders>,
      );

      // De-select STUDENT_TRAVEL from the dropdown
      const ordersTypeDropdown = await screen.findByLabelText('Orders type');
      await userEvent.selectOptions(ordersTypeDropdown, ORDERS_TYPE.PERMANENT_CHANGE_OF_STATION);

      // Submit the form
      const saveButton = screen.getByRole('button', { name: 'Save' });
      await userEvent.click(saveButton);

      // Verify correct values were passed
      await waitFor(() => {
        expect(counselingUpdateOrder).toHaveBeenCalledWith(
          expect.objectContaining({
            body: expect.objectContaining({
              hasDependents: formatYesNoAPIValue('yes'),
              ordersType: ORDERS_TYPE.PERMANENT_CHANGE_OF_STATION,
            }),
          }),
        );
      });
    });

    it('select STUDENT_TRAVEL, De-select STUDENT_TRAVEL', async () => {
      // create a local copy of order return value and set initial values
      const orderQueryReturnValues = JSON.parse(JSON.stringify(useOrdersDocumentQueriesReturnValue));
      orderQueryReturnValues.move = { id: 123, moveCode: 'GLOBAL123', ordersId: 1 };
      orderQueryReturnValues.orders[1].order_type = ORDERS_TYPE.PERMANENT_CHANGE_OF_STATION;
      orderQueryReturnValues.orders[1].has_dependents = formatYesNoAPIValue('no');

      // set return values for mocked functions
      useOrdersDocumentQueries.mockReturnValue(orderQueryReturnValues);
      getOrder.mockResolvedValue(orderQueryReturnValues);

      // render component
      render(
        <MockProviders>
          <ServicesCounselingOrders {...ordersMockProps} />
        </MockProviders>,
      );

      // Select EARLY_RETURN_OF_DEPENDENTS and then de-select from the dropdown
      const ordersTypeDropdown = await screen.findByLabelText('Orders type');
      await userEvent.selectOptions(ordersTypeDropdown, ORDERS_TYPE.STUDENT_TRAVEL);
      await userEvent.selectOptions(ordersTypeDropdown, ORDERS_TYPE.LOCAL_MOVE);

      // Submit the form
      const saveButton = screen.getByRole('button', { name: 'Save' });
      await userEvent.click(saveButton);

      // Verify correct values were passed
      await waitFor(() => {
        expect(counselingUpdateOrder).toHaveBeenCalledWith(
          expect.objectContaining({
            body: expect.objectContaining({
              hasDependents: formatYesNoAPIValue('no'),
              ordersType: ORDERS_TYPE.LOCAL_MOVE,
            }),
          }),
        );
      });
    });

    it('select STUDENT_TRAVEL, select EARLY_RETURN_OF_DEPENDENTS', async () => {
      // create a local copy of order return value and set initial values
      const orderQueryReturnValues = JSON.parse(JSON.stringify(useOrdersDocumentQueriesReturnValue));
      orderQueryReturnValues.move = { id: 123, moveCode: 'GLOBAL123', ordersId: 1 };
      orderQueryReturnValues.orders[1].order_type = ORDERS_TYPE.PERMANENT_CHANGE_OF_STATION;
      orderQueryReturnValues.orders[1].has_dependents = formatYesNoAPIValue('no');

      // set return values for mocked functions
      useOrdersDocumentQueries.mockReturnValue(orderQueryReturnValues);
      getOrder.mockResolvedValue(orderQueryReturnValues);

      // render component
      render(
        <MockProviders>
          <ServicesCounselingOrders {...ordersMockProps} />
        </MockProviders>,
      );

      // Select STUDENT_TRAVEL and then select EARLY_RETURN_OF_DEPENDENTS from the dropdown
      const ordersTypeDropdown = await screen.findByLabelText('Orders type');
      await userEvent.selectOptions(ordersTypeDropdown, ORDERS_TYPE.STUDENT_TRAVEL);
      await userEvent.selectOptions(ordersTypeDropdown, ORDERS_TYPE.EARLY_RETURN_OF_DEPENDENTS);

      // Submit the form
      const saveButton = screen.getByRole('button', { name: 'Save' });
      await userEvent.click(saveButton);

      // Verify correct values were passed
      await waitFor(() => {
        expect(counselingUpdateOrder).toHaveBeenCalledWith(
          expect.objectContaining({
            body: expect.objectContaining({
              hasDependents: formatYesNoAPIValue('yes'),
              ordersType: ORDERS_TYPE.EARLY_RETURN_OF_DEPENDENTS,
            }),
          }),
        );
      });
    });
  });

  describe('Order type: EARLY_RETURN_OF_DEPENDENTS', () => {
    beforeEach(() => {
      jest.clearAllMocks();
    });

    it('select EARLY_RETURN_OF_DEPENDENTS', async () => {
      // create a local copy of order return value and set initial values
      const orderQueryReturnValues = JSON.parse(JSON.stringify(useOrdersDocumentQueriesReturnValue));
      orderQueryReturnValues.move = { id: 123, moveCode: 'GLOBAL123', ordersId: 1 };
      orderQueryReturnValues.orders[1].order_type = ORDERS_TYPE.PERMANENT_CHANGE_OF_STATION;
      orderQueryReturnValues.orders[1].has_dependents = formatYesNoAPIValue('no');

      // set return values for mocked functions
      useOrdersDocumentQueries.mockReturnValue(orderQueryReturnValues);
      getOrder.mockResolvedValue(orderQueryReturnValues);

      // render component
      render(
        <MockProviders>
          <ServicesCounselingOrders {...ordersMockProps} />
        </MockProviders>,
      );

      // Select EARLY_RETURN_OF_DEPENDENTS from the dropdown
      const ordersTypeDropdown = await screen.findByLabelText('Orders type');
      await userEvent.selectOptions(ordersTypeDropdown, ORDERS_TYPE.EARLY_RETURN_OF_DEPENDENTS);

      // Submit the form
      const saveButton = screen.getByRole('button', { name: 'Save' });
      await userEvent.click(saveButton);

      // Verify correct values were passed
      await waitFor(() => {
        expect(counselingUpdateOrder).toHaveBeenCalledWith(
          expect.objectContaining({
            body: expect.objectContaining({
              hasDependents: formatYesNoAPIValue('yes'),
              ordersType: ORDERS_TYPE.EARLY_RETURN_OF_DEPENDENTS,
            }),
          }),
        );
      });
    });

    it('De-select EARLY_RETURN_OF_DEPENDENTS', async () => {
      // create a local copy of order return value and set initial values
      const orderQueryReturnValues = JSON.parse(JSON.stringify(useOrdersDocumentQueriesReturnValue));
      orderQueryReturnValues.move = { id: 123, moveCode: 'GLOBAL123', ordersId: 1 };
      orderQueryReturnValues.orders[1].order_type = ORDERS_TYPE.EARLY_RETURN_OF_DEPENDENTS;
      orderQueryReturnValues.orders[1].has_dependents = formatYesNoAPIValue('yes');

      // set return values for mocked functions
      useOrdersDocumentQueries.mockReturnValue(orderQueryReturnValues);
      getOrder.mockResolvedValue(orderQueryReturnValues);

      // render component
      render(
        <MockProviders>
          <ServicesCounselingOrders {...ordersMockProps} />
        </MockProviders>,
      );

      // De-select EARLY_RETURN_OF_DEPENDENTS from the dropdown
      const ordersTypeDropdown = await screen.findByLabelText('Orders type');
      await userEvent.selectOptions(ordersTypeDropdown, ORDERS_TYPE.PERMANENT_CHANGE_OF_STATION);

      // Submit the form
      const saveButton = screen.getByRole('button', { name: 'Save' });
      await userEvent.click(saveButton);

      // Verify correct values were passed
      await waitFor(() => {
        expect(counselingUpdateOrder).toHaveBeenCalledWith(
          expect.objectContaining({
            body: expect.objectContaining({
              hasDependents: formatYesNoAPIValue('yes'),
              ordersType: ORDERS_TYPE.PERMANENT_CHANGE_OF_STATION,
            }),
          }),
        );
      });
    });

    it('select EARLY_RETURN_OF_DEPENDENTS, De-select EARLY_RETURN_OF_DEPENDENTS', async () => {
      // create a local copy of order return value and set initial values
      const orderQueryReturnValues = JSON.parse(JSON.stringify(useOrdersDocumentQueriesReturnValue));
      orderQueryReturnValues.move = { id: 123, moveCode: 'GLOBAL123', ordersId: 1 };
      orderQueryReturnValues.orders[1].order_type = ORDERS_TYPE.PERMANENT_CHANGE_OF_STATION;
      orderQueryReturnValues.orders[1].has_dependents = formatYesNoAPIValue('no');

      // set return values for mocked functions
      useOrdersDocumentQueries.mockReturnValue(orderQueryReturnValues);
      getOrder.mockResolvedValue(orderQueryReturnValues);

      // render component
      render(
        <MockProviders>
          <ServicesCounselingOrders {...ordersMockProps} />
        </MockProviders>,
      );

      // Select EARLY_RETURN_OF_DEPENDENTS and then de-select from the dropdown
      const ordersTypeDropdown = await screen.findByLabelText('Orders type');
      await userEvent.selectOptions(ordersTypeDropdown, ORDERS_TYPE.EARLY_RETURN_OF_DEPENDENTS);
      await userEvent.selectOptions(ordersTypeDropdown, ORDERS_TYPE.LOCAL_MOVE);

      // Submit the form
      const saveButton = screen.getByRole('button', { name: 'Save' });
      await userEvent.click(saveButton);

      // Verify correct values were passed
      await waitFor(() => {
        expect(counselingUpdateOrder).toHaveBeenCalledWith(
          expect.objectContaining({
            body: expect.objectContaining({
              hasDependents: formatYesNoAPIValue('no'),
              ordersType: ORDERS_TYPE.LOCAL_MOVE,
            }),
          }),
        );
      });
    });

    it('select EARLY_RETURN_OF_DEPENDENTS, select STUDENT_TRAVEL', async () => {
      // create a local copy of order return value and set initial values
      const orderQueryReturnValues = JSON.parse(JSON.stringify(useOrdersDocumentQueriesReturnValue));
      orderQueryReturnValues.move = { id: 123, moveCode: 'GLOBAL123', ordersId: 1 };
      orderQueryReturnValues.orders[1].order_type = ORDERS_TYPE.PERMANENT_CHANGE_OF_STATION;
      orderQueryReturnValues.orders[1].has_dependents = formatYesNoAPIValue('no');

      // set return values for mocked functions
      useOrdersDocumentQueries.mockReturnValue(orderQueryReturnValues);
      getOrder.mockResolvedValue(orderQueryReturnValues);

      // render component
      render(
        <MockProviders>
          <ServicesCounselingOrders {...ordersMockProps} />
        </MockProviders>,
      );

      // Select EARLY_RETURN_OF_DEPENDENTS and then select STUDENT_TRAVEL from the dropdown
      const ordersTypeDropdown = await screen.findByLabelText('Orders type');
      await userEvent.selectOptions(ordersTypeDropdown, ORDERS_TYPE.EARLY_RETURN_OF_DEPENDENTS);
      await userEvent.selectOptions(ordersTypeDropdown, ORDERS_TYPE.STUDENT_TRAVEL);

      // Submit the form
      const saveButton = screen.getByRole('button', { name: 'Save' });
      await userEvent.click(saveButton);

      // Verify correct values were passed
      await waitFor(() => {
        expect(counselingUpdateOrder).toHaveBeenCalledWith(
          expect.objectContaining({
            body: expect.objectContaining({
              hasDependents: formatYesNoAPIValue('yes'),
              ordersType: ORDERS_TYPE.STUDENT_TRAVEL,
            }),
          }),
        );
      });
    });
  });
});

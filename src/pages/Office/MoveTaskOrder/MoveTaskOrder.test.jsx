import React from 'react';
import { mount } from 'enzyme';
import { render, screen, within, cleanup } from '@testing-library/react';
import * as reactQuery from '@tanstack/react-query';
import userEvent from '@testing-library/user-event';

import {
  unapprovedMTOQuery,
  approvedMTOWithCanceledShipmentQuery,
  missingWeightQuery,
  someShipmentsApprovedMTOQuery,
  someWeightNotReturned,
  sitExtensionPresent,
  allApprovedMTOQuery,
  lowerReweighsMTOQuery,
  missingSomeWeightQuery,
  noWeightQuery,
  riskOfExcessWeightQuery,
  lowerActualsMTOQuery,
  sitExtensionApproved,
  allApprovedExternalVendorMTOQuery,
  riskOfExcessWeightQueryExternalShipment,
  riskOfExcessWeightQueryExternalUBShipment,
  multiplePaymentRequests,
  moveHistoryTestData,
  actualPPMWeightQuery,
  approvedMTOWithApprovedSitItemsQuery,
} from './moveTaskOrderUnitTestData';

import { MoveTaskOrder } from 'pages/Office/MoveTaskOrder/MoveTaskOrder';
import { useMoveTaskOrderQueries, useMovePaymentRequestsQueries, useGHCGetMoveHistory } from 'hooks/queries';
import { MockProviders } from 'testUtils';
import { permissionTypes } from 'constants/permissions';
import SERVICE_ITEM_STATUS from 'constants/serviceItems';
import { SITStatusOrigin } from 'components/Office/ShipmentSITDisplay/ShipmentSITDisplayTestParams';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { formatWeight } from 'utils/formatters';

jest.mock('hooks/queries', () => ({
  useMoveTaskOrderQueries: jest.fn(),
  useMovePaymentRequestsQueries: jest.fn(),
  useGHCGetMoveHistory: jest.fn(),
}));

const mockPush = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useHistory: () => ({
    push: mockPush,
  }),
  useParams: () => ({
    moveCode: 'TestCode',
  }),
}));

const setUnapprovedShipmentCount = jest.fn();
const setUnapprovedServiceItemCount = jest.fn();
const setExcessWeightRiskCount = jest.fn();
const setUnapprovedSITExtensionCount = jest.fn();

const requiredProps = {
  setMessage: jest.fn(),
};

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

describe('MoveTaskOrder', () => {
  describe('weight display', () => {
    it('displays the weight allowance', async () => {
      render(
        <MockProviders>
          <MoveTaskOrder
            {...requiredProps}
            setUnapprovedShipmentCount={setUnapprovedShipmentCount}
            setUnapprovedServiceItemCount={setUnapprovedServiceItemCount}
            setExcessWeightRiskCount={setExcessWeightRiskCount}
            setUnapprovedSITExtensionCount={setUnapprovedSITExtensionCount}
          />
        </MockProviders>,
      );

      const weightSummaries = await screen.findAllByTestId('weight-display');
      expect(weightSummaries[0]).toHaveTextContent('8,500 lbs');

      const riskOfExcessAlert = await screen.queryByText(/This move is at risk for excess weight./);
      expect(riskOfExcessAlert).toBeFalsy();

      const riskOfExcessTag = await screen.queryByText(/Risk of excess/);
      expect(riskOfExcessTag).toBeFalsy();

      const externalVendorShipmentCount = await screen.queryByText(/1 shipment not moved by GHC prime./);
      expect(externalVendorShipmentCount).toBeFalsy();
    });

    it('displays the max billable weight', async () => {
      render(
        <MockProviders>
          <MoveTaskOrder
            {...requiredProps}
            setUnapprovedShipmentCount={setUnapprovedShipmentCount}
            setUnapprovedServiceItemCount={setUnapprovedServiceItemCount}
            setExcessWeightRiskCount={setExcessWeightRiskCount}
            setUnapprovedSITExtensionCount={setUnapprovedSITExtensionCount}
          />
        </MockProviders>,
      );

      const weightSummaries = await screen.findAllByTestId('weight-display');
      expect(weightSummaries[2]).toHaveTextContent('8,000 lbs');
    });

    it('displays the estimated total weight with all weights not set', async () => {
      useMoveTaskOrderQueries.mockReturnValue(missingWeightQuery);
      useMovePaymentRequestsQueries.mockReturnValue(multiplePaymentRequests);
      useGHCGetMoveHistory.mockReturnValue(moveHistoryTestData);

      render(
        <MockProviders>
          <MoveTaskOrder
            {...requiredProps}
            setUnapprovedShipmentCount={setUnapprovedShipmentCount}
            setUnapprovedServiceItemCount={setUnapprovedServiceItemCount}
            setExcessWeightRiskCount={setExcessWeightRiskCount}
            setUnapprovedSITExtensionCount={setUnapprovedSITExtensionCount}
          />
        </MockProviders>,
      );

      const weightSummaries = await screen.findAllByTestId('weight-display');
      expect(weightSummaries[1]).toHaveTextContent('—');
    });

    it('displays the move weight total with all weights not set', async () => {
      useMoveTaskOrderQueries.mockReturnValue(missingWeightQuery);

      render(
        <MockProviders>
          <MoveTaskOrder
            {...requiredProps}
            setUnapprovedShipmentCount={setUnapprovedShipmentCount}
            setUnapprovedServiceItemCount={setUnapprovedServiceItemCount}
            setExcessWeightRiskCount={setExcessWeightRiskCount}
            setUnapprovedSITExtensionCount={setUnapprovedSITExtensionCount}
          />
        </MockProviders>,
      );

      const weightSummaries = await screen.findAllByTestId('weight-display');
      expect(weightSummaries[3]).toHaveTextContent('—');
    });

    it('displays the estimated total weight with some weights missing', async () => {
      useMoveTaskOrderQueries.mockReturnValue(missingSomeWeightQuery);

      render(
        <MockProviders>
          <MoveTaskOrder
            {...requiredProps}
            setUnapprovedShipmentCount={setUnapprovedShipmentCount}
            setUnapprovedServiceItemCount={setUnapprovedServiceItemCount}
            setExcessWeightRiskCount={setExcessWeightRiskCount}
            setUnapprovedSITExtensionCount={setUnapprovedSITExtensionCount}
          />
        </MockProviders>,
      );

      const weightSummaries = await screen.findAllByTestId('weight-display');
      expect(weightSummaries[1]).toHaveTextContent('125 lbs');
    });

    it('displays the move weight total with some weights missing', async () => {
      useMoveTaskOrderQueries.mockReturnValue(missingSomeWeightQuery);

      render(
        <MockProviders>
          <MoveTaskOrder
            {...requiredProps}
            setUnapprovedShipmentCount={setUnapprovedShipmentCount}
            setUnapprovedServiceItemCount={setUnapprovedServiceItemCount}
            setExcessWeightRiskCount={setExcessWeightRiskCount}
            setUnapprovedSITExtensionCount={setUnapprovedSITExtensionCount}
          />
        </MockProviders>,
      );

      const weightSummaries = await screen.findAllByTestId('weight-display');
      expect(weightSummaries[3]).toHaveTextContent('125 lbs');
    });

    it('displays the estimated total weight with all not sent', async () => {
      useMoveTaskOrderQueries.mockReturnValue(noWeightQuery);

      render(
        <MockProviders>
          <MoveTaskOrder
            {...requiredProps}
            setUnapprovedShipmentCount={setUnapprovedShipmentCount}
            setUnapprovedServiceItemCount={setUnapprovedServiceItemCount}
            setExcessWeightRiskCount={setExcessWeightRiskCount}
            setUnapprovedSITExtensionCount={setUnapprovedSITExtensionCount}
          />
        </MockProviders>,
      );

      const weightSummaries = await screen.findAllByTestId('weight-display');
      expect(weightSummaries[1]).toHaveTextContent('—');
    });

    it('displays the move weight total with all not sent', async () => {
      useMoveTaskOrderQueries.mockReturnValue(noWeightQuery);

      render(
        <MockProviders>
          <MoveTaskOrder
            {...requiredProps}
            setUnapprovedShipmentCount={setUnapprovedShipmentCount}
            setUnapprovedServiceItemCount={setUnapprovedServiceItemCount}
            setExcessWeightRiskCount={setExcessWeightRiskCount}
            setUnapprovedSITExtensionCount={setUnapprovedSITExtensionCount}
          />
        </MockProviders>,
      );

      const weightSummaries = await screen.findAllByTestId('weight-display');
      expect(weightSummaries[3]).toHaveTextContent('—');
    });

    it('displays the estimated total weight with some sent and some not sent', async () => {
      useMoveTaskOrderQueries.mockReturnValue(someWeightNotReturned);

      render(
        <MockProviders>
          <MoveTaskOrder
            {...requiredProps}
            setUnapprovedShipmentCount={setUnapprovedShipmentCount}
            setUnapprovedServiceItemCount={setUnapprovedServiceItemCount}
            setExcessWeightRiskCount={setExcessWeightRiskCount}
            setUnapprovedSITExtensionCount={setUnapprovedSITExtensionCount}
          />
        </MockProviders>,
      );

      const weightSummaries = await screen.findAllByTestId('weight-display');
      expect(weightSummaries[1]).toHaveTextContent('100');
    });

    it('displays the move weight total with some sent and some not sent', async () => {
      useMoveTaskOrderQueries.mockReturnValue(someWeightNotReturned);

      render(
        <MockProviders>
          <MoveTaskOrder
            {...requiredProps}
            setUnapprovedShipmentCount={setUnapprovedShipmentCount}
            setUnapprovedServiceItemCount={setUnapprovedServiceItemCount}
            setExcessWeightRiskCount={setExcessWeightRiskCount}
            setUnapprovedSITExtensionCount={setUnapprovedSITExtensionCount}
          />
        </MockProviders>,
      );

      const weightSummaries = await screen.findAllByTestId('weight-display');
      expect(weightSummaries[3]).toHaveTextContent('101');
    });

    it('displays risk of excess tag', async () => {
      useMoveTaskOrderQueries.mockReturnValue(riskOfExcessWeightQuery);

      render(
        <MockProviders>
          <MoveTaskOrder
            {...requiredProps}
            setUnapprovedShipmentCount={setUnapprovedShipmentCount}
            setUnapprovedServiceItemCount={setUnapprovedServiceItemCount}
            setExcessWeightRiskCount={setExcessWeightRiskCount}
            setUnapprovedSITExtensionCount={setUnapprovedSITExtensionCount}
          />
        </MockProviders>,
      );

      const riskOfExcessTag = await screen.getByText(/Risk of excess/);
      expect(riskOfExcessTag).toBeInTheDocument();
    });

    it('displays risk of excess unaccompanied baggage tag when a move has excess ub shipment weight', async () => {
      useMoveTaskOrderQueries.mockReturnValue(riskOfExcessWeightQueryExternalUBShipment);

      render(
        <MockProviders>
          <MoveTaskOrder
            {...requiredProps}
            setUnapprovedShipmentCount={setUnapprovedShipmentCount}
            setUnapprovedServiceItemCount={setUnapprovedServiceItemCount}
            setExcessWeightRiskCount={setExcessWeightRiskCount}
            setUnapprovedSITExtensionCount={setUnapprovedSITExtensionCount}
          />
        </MockProviders>,
      );

      const acknowledgeUbExcessButton = await screen.findByTestId(/excessUnaccompaniedBaggageWeightAlertButton/);
      expect(acknowledgeUbExcessButton).toBeInTheDocument();
    });

    it('displays the estimated total weight', async () => {
      useMoveTaskOrderQueries.mockReturnValue(allApprovedMTOQuery);

      render(
        <MockProviders>
          <MoveTaskOrder
            {...requiredProps}
            setUnapprovedShipmentCount={setUnapprovedShipmentCount}
            setUnapprovedServiceItemCount={setUnapprovedServiceItemCount}
            setExcessWeightRiskCount={setExcessWeightRiskCount}
            setUnapprovedSITExtensionCount={setUnapprovedSITExtensionCount}
          />
        </MockProviders>,
      );

      const estimatedWeightTotal = await screen.getByText(/350 lbs/);
      expect(estimatedWeightTotal).toBeInTheDocument();
    });

    it('displays the move weight total', async () => {
      useMoveTaskOrderQueries.mockReturnValue(allApprovedMTOQuery);

      render(
        <MockProviders>
          <MoveTaskOrder
            {...requiredProps}
            setUnapprovedShipmentCount={setUnapprovedShipmentCount}
            setUnapprovedServiceItemCount={setUnapprovedServiceItemCount}
            setExcessWeightRiskCount={setExcessWeightRiskCount}
            setUnapprovedSITExtensionCount={setUnapprovedSITExtensionCount}
          />
        </MockProviders>,
      );

      const moveWeightTotal = await screen.getByText(/350 lbs/);
      expect(moveWeightTotal).toBeInTheDocument();
    });

    it('displays the ppm estimated weight and no ppm actual weight', async () => {
      useMoveTaskOrderQueries.mockReturnValue(allApprovedMTOQuery);

      render(
        <MockProviders>
          <MoveTaskOrder
            {...requiredProps}
            setUnapprovedShipmentCount={setUnapprovedShipmentCount}
            setUnapprovedServiceItemCount={setUnapprovedServiceItemCount}
            setExcessWeightRiskCount={setExcessWeightRiskCount}
            setUnapprovedSITExtensionCount={setUnapprovedSITExtensionCount}
          />
        </MockProviders>,
      );

      const weightSummaries = await screen.findAllByTestId('weight-display');

      expect(weightSummaries[4]).toHaveTextContent('2,000 lbs');
      expect(weightSummaries[5]).toHaveTextContent('—');
    });

    it('displays the ppm actual weight (total)', async () => {
      useMoveTaskOrderQueries.mockReturnValue(actualPPMWeightQuery);

      render(
        <MockProviders>
          <MoveTaskOrder
            {...requiredProps}
            setUnapprovedShipmentCount={setUnapprovedShipmentCount}
            setUnapprovedServiceItemCount={setUnapprovedServiceItemCount}
            setExcessWeightRiskCount={setExcessWeightRiskCount}
            setUnapprovedSITExtensionCount={setUnapprovedSITExtensionCount}
          />
        </MockProviders>,
      );

      const weightSummaries = await screen.findAllByTestId('weight-display');

      expect(weightSummaries[4]).toHaveTextContent('2,000 lbs');
      expect(weightSummaries[5]).toHaveTextContent('2,100 lbs');
    });

    it('displays the move weight total using lower reweighs', async () => {
      useMoveTaskOrderQueries.mockReturnValue(lowerReweighsMTOQuery);

      render(
        <MockProviders>
          <MoveTaskOrder
            {...requiredProps}
            setUnapprovedShipmentCount={setUnapprovedShipmentCount}
            setUnapprovedServiceItemCount={setUnapprovedServiceItemCount}
            setExcessWeightRiskCount={setExcessWeightRiskCount}
            setUnapprovedSITExtensionCount={setUnapprovedSITExtensionCount}
          />
        </MockProviders>,
      );

      const moveWeightTotal = await screen.getByText(/247 lbs/);
      expect(moveWeightTotal).toBeInTheDocument();
    });

    it('displays the move weight total using lower actual weights', async () => {
      useMoveTaskOrderQueries.mockReturnValue(lowerActualsMTOQuery);

      render(
        <MockProviders>
          <MoveTaskOrder
            {...requiredProps}
            setUnapprovedShipmentCount={setUnapprovedShipmentCount}
            setUnapprovedServiceItemCount={setUnapprovedServiceItemCount}
            setExcessWeightRiskCount={setExcessWeightRiskCount}
            setUnapprovedSITExtensionCount={setUnapprovedSITExtensionCount}
          />
        </MockProviders>,
      );

      const moveWeightTotal = await screen.getByText(/250 lbs/);
      expect(moveWeightTotal).toBeInTheDocument();
    });

    it('displays the external vendor shipment count and link to move details when there are external vendor shipments', async () => {
      useMoveTaskOrderQueries.mockReturnValue(allApprovedExternalVendorMTOQuery);

      render(
        <MockProviders>
          <MoveTaskOrder
            {...requiredProps}
            setUnapprovedShipmentCount={setUnapprovedShipmentCount}
            setUnapprovedServiceItemCount={setUnapprovedServiceItemCount}
            setExcessWeightRiskCount={setExcessWeightRiskCount}
            setUnapprovedSITExtensionCount={setUnapprovedSITExtensionCount}
          />
        </MockProviders>,
      );

      const externalVendorShipmentCount = await screen.getByText(/1 shipment not moved by GHC prime./);
      expect(externalVendorShipmentCount).toBeInTheDocument();
    });

    it('displays risk of excess tag and external vendor shipment count', async () => {
      useMoveTaskOrderQueries.mockReturnValue(riskOfExcessWeightQueryExternalShipment);

      render(
        <MockProviders>
          <MoveTaskOrder
            {...requiredProps}
            setUnapprovedShipmentCount={setUnapprovedShipmentCount}
            setUnapprovedServiceItemCount={setUnapprovedServiceItemCount}
            setExcessWeightRiskCount={setExcessWeightRiskCount}
            setUnapprovedSITExtensionCount={setUnapprovedSITExtensionCount}
          />
        </MockProviders>,
      );

      const riskOfExcessTag = await screen.getByText(/Risk of excess/);
      expect(riskOfExcessTag).toBeInTheDocument();
      const externalVendorShipmentCount = await screen.getByText(/1 shipment not moved by GHC prime./);
      expect(externalVendorShipmentCount).toBeInTheDocument();
    });
  });

  describe('check loading and error component states', () => {
    it('renders the Loading Placeholder when the query is still loading', async () => {
      useMoveTaskOrderQueries.mockReturnValue(loadingReturnValue);

      render(
        <MockProviders>
          <MoveTaskOrder
            {...requiredProps}
            setUnapprovedShipmentCount={setUnapprovedShipmentCount}
            setUnapprovedServiceItemCount={setUnapprovedServiceItemCount}
            setExcessWeightRiskCount={setExcessWeightRiskCount}
            setUnapprovedSITExtensionCount={setUnapprovedSITExtensionCount}
          />
        </MockProviders>,
      );

      const h2 = await screen.getByRole('heading', { name: 'Loading, please wait...', level: 2 });
      expect(h2).toBeInTheDocument();
    });

    it('renders the Something Went Wrong component when the query errors', async () => {
      useMoveTaskOrderQueries.mockReturnValue(errorReturnValue);

      render(
        <MockProviders>
          <MoveTaskOrder
            {...requiredProps}
            setUnapprovedShipmentCount={setUnapprovedShipmentCount}
            setUnapprovedServiceItemCount={setUnapprovedServiceItemCount}
            setExcessWeightRiskCount={setExcessWeightRiskCount}
            setUnapprovedSITExtensionCount={setUnapprovedSITExtensionCount}
          />
        </MockProviders>,
      );

      const errorMessage = await screen.getByText(/Something went wrong./);
      expect(errorMessage).toBeInTheDocument();
    });
  });

  describe('move is not available to prime', () => {
    useMoveTaskOrderQueries.mockReturnValue(unapprovedMTOQuery);
    const wrapper = mount(
      <MockProviders>
        <MoveTaskOrder
          {...requiredProps}
          setUnapprovedShipmentCount={setUnapprovedShipmentCount}
          setUnapprovedServiceItemCount={setUnapprovedServiceItemCount}
          setExcessWeightRiskCount={setExcessWeightRiskCount}
          setUnapprovedSITExtensionCount={setUnapprovedSITExtensionCount}
        />
      </MockProviders>,
    );

    it('renders the h1', () => {
      expect(wrapper.find({ 'data-testid': 'too-shipment-container' }).exists()).toBe(true);
      expect(wrapper.find('h1').text()).toBe('Move Task Order');
    });

    it('displays empty state message', () => {
      expect(
        wrapper
          .find('[data-testid="too-shipment-container"] p')
          .contains('This move does not have any approved shipments yet.'),
      ).toBe(true);
    });

    it('updates the unapproved shipments tag state', () => {
      expect(setUnapprovedShipmentCount).toHaveBeenCalledWith(2);
    });

    it('updates the unapproved shipments tag state', () => {
      expect(setUnapprovedServiceItemCount).toHaveBeenCalledWith(0);
    });
  });

  describe('SIT entry date update', () => {
    const mockMutateServiceItemSitEntryDate = jest.fn();
    jest.spyOn(reactQuery, 'useMutation').mockImplementation(() => ({
      mutate: mockMutateServiceItemSitEntryDate,
    }));
    beforeEach(() => {
      // Reset the mock before each test
      mockMutateServiceItemSitEntryDate.mockReset();
    });
    afterEach(() => {
      cleanup(); // This will unmount the component after each test
    });

    const renderComponent = () => {
      useMoveTaskOrderQueries.mockReturnValue(approvedMTOWithApprovedSitItemsQuery);
      useMovePaymentRequestsQueries.mockReturnValue({ paymentRequests: [] });
      useGHCGetMoveHistory.mockReturnValue(moveHistoryTestData);
      const isMoveLocked = false;
      render(
        <MockProviders permissions={[permissionTypes.updateMTOServiceItem, permissionTypes.updateMTOPage]}>
          <MoveTaskOrder
            {...requiredProps}
            setUnapprovedShipmentCount={setUnapprovedShipmentCount}
            setUnapprovedServiceItemCount={setUnapprovedServiceItemCount}
            setExcessWeightRiskCount={setExcessWeightRiskCount}
            setUnapprovedSITExtensionCount={setUnapprovedSITExtensionCount}
            isMoveLocked={isMoveLocked}
          />
        </MockProviders>,
      );
    };
    it('shows error message when SIT entry date is invalid', async () => {
      renderComponent();
      // Set up the mock to simulate an error
      mockMutateServiceItemSitEntryDate.mockImplementation((data, options) => {
        options.onError({
          response: {
            status: 422,
            data: JSON.stringify({
              detail:
                'UpdateSitEntryDate failed for service item: the SIT Entry Date (2025-03-05) must be before the SIT Departure Date (2025-02-27)',
            }),
          },
        });
      });
      const approvedServiceItems = await screen.findByTestId('ApprovedServiceItemsTable');
      expect(approvedServiceItems).toBeInTheDocument();
      const spanElement = within(approvedServiceItems).getByText(/Domestic origin 1st day SIT/i);
      expect(spanElement).toBeInTheDocument();
      // Search for the edit button within the approvedServiceItems div
      const editButton = within(approvedServiceItems).getByRole('button', { name: /edit/i });
      expect(editButton).toBeInTheDocument();
      await userEvent.click(editButton);
      const modal = await screen.findByTestId('modal');
      expect(modal).toBeInTheDocument();
      const heading = within(modal).getByRole('heading', { name: /Edit SIT Entry Date/i, level: 2 });
      expect(heading).toBeInTheDocument();
      const formGroups = screen.getAllByTestId('formGroup');
      const sitEntryDateFormGroup = Array.from(formGroups).find(
        (group) =>
          within(group).queryByPlaceholderText('DD MMM YYYY') &&
          within(group).queryByPlaceholderText('DD MMM YYYY').getAttribute('name') === 'sitEntryDate',
      );
      const dateInput = within(sitEntryDateFormGroup).getByPlaceholderText('DD MMM YYYY');
      expect(dateInput).toBeInTheDocument();
      const remarksTextarea = within(modal).getByTestId('officeRemarks');
      expect(remarksTextarea).toBeInTheDocument();
      const saveButton = within(modal).getByRole('button', { name: /Save/ });

      await userEvent.clear(dateInput);
      await userEvent.type(dateInput, '05 Mar 2025');
      await userEvent.type(remarksTextarea, 'Need to update the sit entry date.');
      expect(saveButton).toBeEnabled();
      await userEvent.click(saveButton);

      // Verify that the mutation was called
      expect(mockMutateServiceItemSitEntryDate).toHaveBeenCalled();

      // The modal should close
      expect(screen.queryByTestId('modal')).not.toBeInTheDocument();

      // Verify that the error message is displayed
      const alert = screen.getByTestId('alert');
      expect(alert).toBeInTheDocument();
      expect(alert).toHaveClass('usa-alert--error');
      expect(alert).toHaveTextContent(
        'UpdateSitEntryDate failed for service item: the SIT Entry Date (2025-03-05) must be before the SIT Departure Date (2025-02-27)',
      );
    });

    it('shows success message when SIT entry date is valid', async () => {
      renderComponent();
      // Set up the mock to simulate an error
      mockMutateServiceItemSitEntryDate.mockImplementation((data, options) => {
        options.onSuccess({
          response: {
            status: 200,
            data: JSON.stringify({
              detail: 'SIT entry date updated',
            }),
          },
        });
      });
      const approvedServiceItems = await screen.findByTestId('ApprovedServiceItemsTable');
      expect(approvedServiceItems).toBeInTheDocument();
      const spanElement = within(approvedServiceItems).getByText(/Domestic origin 1st day SIT/i);
      expect(spanElement).toBeInTheDocument();
      // Search for the edit button within the approvedServiceItems div
      const editButton = within(approvedServiceItems).getByRole('button', { name: /edit/i });
      expect(editButton).toBeInTheDocument();
      await userEvent.click(editButton);
      const modal = await screen.findByTestId('modal');
      expect(modal).toBeInTheDocument();
      const heading = within(modal).getByRole('heading', { name: /Edit SIT Entry Date/i, level: 2 });
      expect(heading).toBeInTheDocument();
      const formGroups = screen.getAllByTestId('formGroup');
      const sitEntryDateFormGroup = Array.from(formGroups).find(
        (group) =>
          within(group).queryByPlaceholderText('DD MMM YYYY') &&
          within(group).queryByPlaceholderText('DD MMM YYYY').getAttribute('name') === 'sitEntryDate',
      );
      const dateInput = within(sitEntryDateFormGroup).getByPlaceholderText('DD MMM YYYY');
      expect(dateInput).toBeInTheDocument();
      const remarksTextarea = within(modal).getByTestId('officeRemarks');
      expect(remarksTextarea).toBeInTheDocument();
      const saveButton = within(modal).getByRole('button', { name: /Save/ });

      await userEvent.clear(dateInput);
      await userEvent.type(dateInput, '03 Mar 2024');
      await userEvent.type(remarksTextarea, 'Need to update the sit entry date.');
      expect(saveButton).toBeEnabled();
      await userEvent.click(saveButton);

      // Verify that the mutation was called
      expect(mockMutateServiceItemSitEntryDate).toHaveBeenCalled();

      // The modal should close
      expect(screen.queryByTestId('modal')).not.toBeInTheDocument();

      // Verify that the error message is displayed
      const alert = screen.getByTestId('alert');
      expect(alert).toBeInTheDocument();
      expect(alert).toHaveClass('usa-alert--success');
      expect(alert).toHaveTextContent('SIT entry date updated');
    });
  });

  describe('approved mto with both submitted and approved shipments', () => {
    useMoveTaskOrderQueries.mockReturnValue(someShipmentsApprovedMTOQuery);
    useMovePaymentRequestsQueries.mockReturnValue(multiplePaymentRequests);
    useGHCGetMoveHistory.mockReturnValue(moveHistoryTestData);
    const isMoveLocked = false;
    const wrapper = mount(
      <MockProviders permissions={[permissionTypes.createShipmentCancellation, permissionTypes.updateMTOPage]}>
        <MoveTaskOrder
          {...requiredProps}
          setUnapprovedShipmentCount={setUnapprovedShipmentCount}
          setUnapprovedServiceItemCount={setUnapprovedServiceItemCount}
          setExcessWeightRiskCount={setExcessWeightRiskCount}
          setUnapprovedSITExtensionCount={setUnapprovedSITExtensionCount}
          isMoveLocked={isMoveLocked}
        />
      </MockProviders>,
    );

    it('renders the h1', () => {
      expect(wrapper.find({ 'data-testid': 'too-shipment-container' }).exists()).toBe(true);
      expect(wrapper.find('h1').text()).toBe('Move Task Order');
    });

    it('renders the left nav with shipments', () => {
      expect(wrapper.find('nav').exists()).toBe(true);

      const navLinks = wrapper.find('nav a');
      expect(navLinks.length).toBe(2);
      expect(navLinks.at(1).contains('HHG shipment')).toBe(true);
      expect(navLinks.at(1).prop('href')).toBe('#s-3');
    });

    it('renders the left nav with move weights', () => {
      expect(wrapper.find('nav').exists()).toBe(true);

      const navLinks = wrapper.find('nav a');
      expect(navLinks.length).toBe(2);
      expect(navLinks.at(0).contains('Move weights')).toBe(true);
      expect(navLinks.at(0).prop('href')).toBe('#move-weights');
    });

    it('renders the ShipmentContainer', () => {
      expect(wrapper.find('ShipmentContainer').length).toBe(1);
    });

    it('renders the ShipmentHeading', () => {
      expect(wrapper.find('ShipmentHeading').exists()).toBe(true);
      expect(wrapper.find('h2').at(0).text()).toEqual('Household goods');
      expect(wrapper.find('h4').at(0).text()).toEqual('#');
    });

    it('renders the ImportantShipmentDates', () => {
      expect(wrapper.find('ImportantShipmentDates').exists()).toBe(true);
    });

    it('renders the ShipmentAddresses', () => {
      expect(wrapper.find('ShipmentAddresses').exists()).toBe(true);
    });

    it('renders the ShipmentWeightDetails', () => {
      expect(wrapper.find('ShipmentWeightDetails').exists()).toBe(true);
    });

    it('renders the RequestedServiceItemsTable for requested, approved, and rejected service items', () => {
      const requestedServiceItemsTable = wrapper.find('RequestedServiceItemsTable');
      // There should be 1 of each status table requested, approved, rejected service items
      // Plus approved move-level service items separate from the shipment items
      expect(requestedServiceItemsTable.length).toBe(6);
      expect(requestedServiceItemsTable.at(0).prop('statusForTableType')).toBe(SERVICE_ITEM_STATUS.SUBMITTED);
      expect(requestedServiceItemsTable.at(1).prop('statusForTableType')).toBe(SERVICE_ITEM_STATUS.APPROVED);
      expect(requestedServiceItemsTable.at(2).prop('statusForTableType')).toBe(SERVICE_ITEM_STATUS.REJECTED);
      expect(requestedServiceItemsTable.at(3).prop('statusForTableType')).toBe('Move Task Order Requested');
      expect(requestedServiceItemsTable.at(4).prop('statusForTableType')).toBe('Move Task Order Approved');
      expect(requestedServiceItemsTable.at(5).prop('statusForTableType')).toBe('Move Task Order Rejected');
    });

    it('updates the unapproved shipments tag state', () => {
      expect(setUnapprovedShipmentCount).toHaveBeenCalledWith(0);
    });

    it('updates the unapproved service items tag state', () => {
      expect(setUnapprovedServiceItemCount).toHaveBeenCalledWith(1);
    });
  });

  describe('approved mto with approved shipments', () => {
    useMoveTaskOrderQueries.mockReturnValue(allApprovedMTOQuery);
    const wrapper = mount(
      <MockProviders>
        <MoveTaskOrder
          {...requiredProps}
          setUnapprovedShipmentCount={setUnapprovedShipmentCount}
          setUnapprovedServiceItemCount={setUnapprovedServiceItemCount}
          setExcessWeightRiskCount={setExcessWeightRiskCount}
          setUnapprovedSITExtensionCount={setUnapprovedSITExtensionCount}
        />
      </MockProviders>,
    );

    it('renders the h1', () => {
      expect(wrapper.find({ 'data-testid': 'too-shipment-container' }).exists()).toBe(true);
      expect(wrapper.find('h1').text()).toBe('Move Task Order');
    });

    it('renders the left nav with shipments', () => {
      expect(wrapper.find('nav').exists()).toBe(true);

      const navLinks = wrapper.find('nav a');
      expect(navLinks.at(1).contains('HHG shipment')).toBe(true);
      expect(navLinks.at(1).contains('1'));
      expect(navLinks.at(1).prop('href')).toBe('#s-3');

      expect(navLinks.at(2).contains('NTS shipment')).toBe(true);
      expect(navLinks.at(2).contains('1'));
      expect(navLinks.at(2).prop('href')).toBe('#s-4');

      expect(navLinks.at(3).contains('NTS-release shipment')).toBe(true);
      expect(navLinks.at(3).prop('href')).toBe('#s-5');

      expect(navLinks.at(4).contains('HHG shipment')).toBe(true);
      expect(navLinks.at(4).prop('href')).toBe('#s-6');

      expect(navLinks.at(5).contains('HHG shipment')).toBe(true);
      expect(navLinks.at(5).prop('href')).toBe('#s-7');
    });

    it('renders the ShipmentContainer', () => {
      expect(wrapper.find('ShipmentContainer').length).toBe(6);
    });

    it('renders the ShipmentHeading', () => {
      expect(wrapper.find('ShipmentHeading').exists()).toBe(true);
      expect(wrapper.find('h2').at(0).text()).toEqual('Household goods');
      expect(wrapper.find('h2').at(1).text()).toEqual('Non-temp storage');
    });

    it('renders the ImportantShipmentDates', () => {
      expect(wrapper.find('ImportantShipmentDates').exists()).toBe(true);
    });

    it('renders the ShipmentAddresses', () => {
      expect(wrapper.find('ShipmentAddresses').exists()).toBe(true);
    });

    it('renders the ShipmentWeightDetails', () => {
      expect(wrapper.find('ShipmentWeightDetails').exists()).toBe(true);
    });

    it('renders the RequestedServiceItemsTable for SUBMITTED service item', () => {
      const requestedServiceItemsTable = wrapper.find('RequestedServiceItemsTable');
      // There are no approved or rejected service item tables to display
      expect(requestedServiceItemsTable.length).toBe(2);
      expect(requestedServiceItemsTable.at(0).prop('statusForTableType')).toBe(SERVICE_ITEM_STATUS.SUBMITTED);
      expect(requestedServiceItemsTable.at(1).prop('statusForTableType')).toBe(SERVICE_ITEM_STATUS.SUBMITTED);
    });

    it('updates the unapproved shipments tag state', () => {
      expect(setUnapprovedShipmentCount).toHaveBeenCalledWith(0);
    });

    it('updates the unapproved service items tag state', () => {
      expect(setUnapprovedServiceItemCount).toHaveBeenCalledWith(2);
    });
  });

  describe('approved mto with canceled shipment', () => {
    useMoveTaskOrderQueries.mockReturnValue(approvedMTOWithCanceledShipmentQuery);
    const wrapper = mount(
      <MockProviders>
        <MoveTaskOrder
          {...requiredProps}
          setUnapprovedShipmentCount={setUnapprovedShipmentCount}
          setUnapprovedServiceItemCount={setUnapprovedServiceItemCount}
          setExcessWeightRiskCount={setExcessWeightRiskCount}
          setUnapprovedSITExtensionCount={setUnapprovedSITExtensionCount}
        />
      </MockProviders>,
    );

    it('renders the h1', () => {
      expect(wrapper.find({ 'data-testid': 'too-shipment-container' }).exists()).toBe(true);
      expect(wrapper.find('h1').text()).toBe('Move Task Order');
    });

    it('renders the left nav with shipments', () => {
      expect(wrapper.find('nav').exists()).toBe(true);

      const navLinks = wrapper.find('nav a');
      expect(navLinks.at(1).contains('HHG shipment')).toBe(true);
      expect(navLinks.at(1).contains('1'));
      expect(navLinks.at(1).prop('href')).toBe('#s-3');
    });

    it('renders the ShipmentContainer', () => {
      expect(wrapper.find('ShipmentContainer').length).toBe(1);
    });

    it('renders the ShipmentHeading', () => {
      expect(wrapper.find('ShipmentHeading').exists()).toBe(true);
      expect(wrapper.find('h2').at(0).text()).toEqual('Household goods');
      expect(wrapper.find('span[data-testid="tag"]').at(0).text()).toEqual('canceled');
    });

    it('renders the ImportantShipmentDates', () => {
      expect(wrapper.find('ImportantShipmentDates').exists()).toBe(true);
    });

    it('renders the ShipmentAddresses', () => {
      expect(wrapper.find('ShipmentAddresses').exists()).toBe(true);
    });

    it('renders the ShipmentWeightDetails', () => {
      expect(wrapper.find('ShipmentWeightDetails').exists()).toBe(true);
      expect(wrapper.find('span[data-testid="tag"]').at(1).text()).toEqual('reweigh requested');
    });

    it('renders the RequestedServiceItemsTable for SUBMITTED service item', () => {
      const requestedServiceItemsTable = wrapper.find('RequestedServiceItemsTable');
      // There are no approved or rejected service item tables to display
      expect(requestedServiceItemsTable.length).toBe(1);
      expect(requestedServiceItemsTable.at(0).prop('statusForTableType')).toBe(SERVICE_ITEM_STATUS.SUBMITTED);
    });

    it('updates the unapproved shipments tag state', () => {
      expect(setUnapprovedShipmentCount).toHaveBeenCalledWith(0);
    });

    it('updates the unapproved service items tag state', () => {
      expect(setUnapprovedServiceItemCount).toHaveBeenCalledWith(2);
    });
  });
  describe('SIT extension pending', () => {
    useMoveTaskOrderQueries.mockReturnValue(sitExtensionPresent);
    const wrapper = mount(
      <MockProviders>
        <MoveTaskOrder
          {...requiredProps}
          setUnapprovedShipmentCount={setUnapprovedShipmentCount}
          setUnapprovedServiceItemCount={setUnapprovedServiceItemCount}
          setExcessWeightRiskCount={setExcessWeightRiskCount}
          setUnapprovedSITExtensionCount={setUnapprovedSITExtensionCount}
        />
      </MockProviders>,
    );

    it('updates the unapproved SIT extension count state', () => {
      expect(setUnapprovedSITExtensionCount).toHaveBeenCalledWith(1);
    });

    it('renders the left nav with tag for SIT extension request', () => {
      expect(wrapper.find('nav').exists()).toBe(true);
      const navLinks = wrapper.find('nav a');
      expect(navLinks.at(1).contains('HHG shipment')).toBe(true);
      expect(navLinks.at(1).contains('1'));
    });
  });
  describe('SIT extension approved', () => {
    useMoveTaskOrderQueries.mockReturnValue(sitExtensionApproved);
    const wrapper = mount(
      <MockProviders>
        <MoveTaskOrder
          {...requiredProps}
          setUnapprovedShipmentCount={setUnapprovedShipmentCount}
          setUnapprovedServiceItemCount={setUnapprovedServiceItemCount}
          setExcessWeightRiskCount={setExcessWeightRiskCount}
          setUnapprovedSITExtensionCount={setUnapprovedSITExtensionCount}
        />
      </MockProviders>,
    );

    it('updates the unapproved SIT extension count state (with a zero count)', () => {
      expect(setUnapprovedSITExtensionCount).toHaveBeenCalledWith(0);
    });

    it('renders the left nav with tag for SIT extension request without a number tag', () => {
      expect(wrapper.find('nav').exists()).toBe(true);
      const navLinks = wrapper.find('nav a');
      // We should get just the shipment text in the nav link
      expect(navLinks.at(1).text()).toEqual('HHG shipment ');
    });
  });

  describe('permission dependent rendering', () => {
    beforeEach(() => {
      useMoveTaskOrderQueries.mockReturnValue(someShipmentsApprovedMTOQuery);
    });

    const testProps = {
      ...requiredProps,
      setUnapprovedShipmentCount,
      setUnapprovedServiceItemCount,
      setExcessWeightRiskCount,
      setUnapprovedSITExtensionCount,
    };

    it('renders the financial review flag button when user has permission', async () => {
      render(
        <MockProviders permissions={[permissionTypes.updateFinancialReviewFlag, permissionTypes.updateMTOPage]}>
          <MoveTaskOrder {...testProps} />
        </MockProviders>,
      );

      expect(await screen.getByText('Flag move for financial review')).toBeInTheDocument();
    });

    it('does not show the financial review flag button if user does not have permission', () => {
      render(
        <MockProviders>
          <MoveTaskOrder {...testProps} />
        </MockProviders>,
      );

      expect(screen.queryByText('Flag move for financial review')).not.toBeInTheDocument();
    });

    it('does not show the financial review flag button if user does not have updateMTOPage permission', () => {
      render(
        <MockProviders permissions={[permissionTypes.updateFinancialReviewFlag]}>
          <MoveTaskOrder {...testProps} />
        </MockProviders>,
      );

      expect(screen.queryByText('Flag move for financial review')).not.toBeInTheDocument();
    });
  });

  describe('estimated weight breakdown', () => {
    it('should show UB estimated weight', () => {
      const testShipments = [
        {
          id: '1',
          moveTaskOrderID: '2',
          shipmentType: SHIPMENT_OPTIONS.UNACCOMPANIED_BAGGAGE,
          scheduledPickupDate: '2020-03-16',
          requestedPickupDate: '2020-03-15',
          pickupAddress: {
            streetAddress1: '932 Baltic Avenue',
            city: 'Chicago',
            state: 'IL',
            postalCode: '60601',
          },
          destinationAddress: {
            streetAddress1: '10 Park Place',
            city: 'Atlantic City',
            state: 'NJ',
            postalCode: '08401',
          },
          status: 'APPROVED',
          eTag: '1234',
          primeEstimatedWeight: 1850,
          primeActualWeight: 1841,
          sitExtensions: [],
          sitStatus: SITStatusOrigin,
        },
      ];
      useMoveTaskOrderQueries.mockReturnValue({
        ...riskOfExcessWeightQueryExternalUBShipment,
        mtoShipments: testShipments,
      });

      render(
        <MockProviders>
          <MoveTaskOrder
            {...requiredProps}
            setUnapprovedShipmentCount={setUnapprovedShipmentCount}
            setUnapprovedServiceItemCount={setUnapprovedServiceItemCount}
            setExcessWeightRiskCount={setExcessWeightRiskCount}
            setUnapprovedSITExtensionCount={setUnapprovedSITExtensionCount}
          />
        </MockProviders>,
      );

      const breakdownToggle = screen.queryByText('Show Breakdown');
      expect(breakdownToggle).toBeInTheDocument();

      breakdownToggle.click();
      expect(breakdownToggle).toHaveTextContent('Hide Breakdown');
      expect(screen.queryByText('110% Estimated UB')).toBeInTheDocument();

      const ubEstimatedWeightValue = screen.getByTestId('breakdownUBEstimatedWeight');
      expect(ubEstimatedWeightValue).toBeInTheDocument();
      expect(ubEstimatedWeightValue).toHaveTextContent(`${formatWeight(testShipments[0].primeEstimatedWeight * 1.1)}`);
    });
  });
});

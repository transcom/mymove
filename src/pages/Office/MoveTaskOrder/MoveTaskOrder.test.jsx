import React from 'react';
import { mount } from 'enzyme';
import { render, screen } from '@testing-library/react';

import {
  unapprovedMTOQuery,
  approvedMTOWithCancelledShipmentQuery,
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
} from './moveTaskOrderUnitTestData';

import { MoveTaskOrder } from 'pages/Office/MoveTaskOrder/MoveTaskOrder';
import { useMoveTaskOrderQueries } from 'hooks/queries';
import { MockProviders } from 'testUtils';
import SERVICE_ITEM_STATUS from 'constants/serviceItems';

jest.mock('hooks/queries', () => ({
  useMoveTaskOrderQueries: jest.fn(),
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

const moveCode = 'WE31AZ';
const requiredProps = {
  match: { params: { moveCode } },
  history: { push: jest.fn() },
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
        <MockProviders initialEntries={['moves/1000/allowances']}>
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
    });

    it('displays the max billable weight', async () => {
      render(
        <MockProviders initialEntries={['moves/1000/allowances']}>
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

      render(
        <MockProviders initialEntries={['moves/1000/allowances']}>
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
        <MockProviders initialEntries={['moves/1000/allowances']}>
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
        <MockProviders initialEntries={['moves/1000/allowances']}>
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
        <MockProviders initialEntries={['moves/1000/allowances']}>
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
        <MockProviders initialEntries={['moves/1000/allowances']}>
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
        <MockProviders initialEntries={['moves/1000/allowances']}>
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
        <MockProviders initialEntries={['moves/1000/allowances']}>
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
      expect(weightSummaries[1]).toHaveTextContent('101');
    });

    it('displays the move weight total with some sent and some not sent', async () => {
      useMoveTaskOrderQueries.mockReturnValue(someWeightNotReturned);

      render(
        <MockProviders initialEntries={['moves/1000/allowances']}>
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
        <MockProviders initialEntries={['moves/1000/allowances']}>
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

    it('displays risk of excess alert', async () => {
      useMoveTaskOrderQueries.mockReturnValue(riskOfExcessWeightQuery);

      render(
        <MockProviders initialEntries={['moves/1000/allowances']}>
          <MoveTaskOrder
            {...requiredProps}
            setUnapprovedShipmentCount={setUnapprovedShipmentCount}
            setUnapprovedServiceItemCount={setUnapprovedServiceItemCount}
            setExcessWeightRiskCount={setExcessWeightRiskCount}
            setUnapprovedSITExtensionCount={setUnapprovedSITExtensionCount}
          />
        </MockProviders>,
      );

      expect(setExcessWeightRiskCount).toHaveBeenCalledWith(1);

      const riskOfExcessAlert = await screen.getByText(/This move is at risk for excess weight./);
      expect(riskOfExcessAlert).toBeInTheDocument();
    });

    it('displays the estimated total weight', async () => {
      useMoveTaskOrderQueries.mockReturnValue(allApprovedMTOQuery);

      render(
        <MockProviders initialEntries={['moves/1000/allowances']}>
          <MoveTaskOrder
            {...requiredProps}
            setUnapprovedShipmentCount={setUnapprovedShipmentCount}
            setUnapprovedServiceItemCount={setUnapprovedServiceItemCount}
            setExcessWeightRiskCount={setExcessWeightRiskCount}
            setUnapprovedSITExtensionCount={setUnapprovedSITExtensionCount}
          />
        </MockProviders>,
      );

      const estimatedWeightTotal = await screen.getByText(/400 lbs/);
      expect(estimatedWeightTotal).toBeInTheDocument();
    });

    it('displays the move weight total', async () => {
      useMoveTaskOrderQueries.mockReturnValue(allApprovedMTOQuery);

      render(
        <MockProviders initialEntries={['moves/1000/allowances']}>
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

    it('displays the move weight total using lower reweighs', async () => {
      useMoveTaskOrderQueries.mockReturnValue(lowerReweighsMTOQuery);

      render(
        <MockProviders initialEntries={['moves/1000/allowances']}>
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
        <MockProviders initialEntries={['moves/1000/allowances']}>
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
  });

  describe('check loading and error component states', () => {
    it('renders the Loading Placeholder when the query is still loading', async () => {
      useMoveTaskOrderQueries.mockReturnValue(loadingReturnValue);

      render(
        <MockProviders initialEntries={['moves/1000/allowances']}>
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
        <MockProviders initialEntries={['moves/1000/allowances']}>
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
      expect(wrapper.find('h1').text()).toBe('Move task order');
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

  describe('approved mto with both submitted and approved shipments', () => {
    useMoveTaskOrderQueries.mockReturnValue(someShipmentsApprovedMTOQuery);
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
      expect(wrapper.find('h1').text()).toBe('Move task order');
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
      expect(wrapper.find('[data-testid="button"]').exists()).toBe(true);
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
      expect(requestedServiceItemsTable.length).toBe(3);
      expect(requestedServiceItemsTable.at(0).prop('statusForTableType')).toBe(SERVICE_ITEM_STATUS.SUBMITTED);
      expect(requestedServiceItemsTable.at(1).prop('statusForTableType')).toBe(SERVICE_ITEM_STATUS.APPROVED);
      expect(requestedServiceItemsTable.at(2).prop('statusForTableType')).toBe(SERVICE_ITEM_STATUS.REJECTED);
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
      expect(wrapper.find('h1').text()).toBe('Move task order');
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
      expect(wrapper.find('ShipmentContainer').length).toBe(5);
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

  describe('approved mto with cancelled shipment', () => {
    useMoveTaskOrderQueries.mockReturnValue(approvedMTOWithCancelledShipmentQuery);
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
      expect(wrapper.find('h1').text()).toBe('Move task order');
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
      expect(wrapper.find('span[data-testid="tag"]').at(0).text()).toEqual('cancelled');
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
});

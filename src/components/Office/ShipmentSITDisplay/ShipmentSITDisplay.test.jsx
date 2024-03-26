import React from 'react';
import { render, screen, waitFor, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import ShipmentSITDisplay from './ShipmentSITDisplay';
import {
  futureSITShipment,
  futureSITStatus,
  SITExtensions,
  SITStatusOrigin,
  SITStatusOriginAuthorized,
  SITStatusDestination,
  SITStatusDestinationWithoutCustomerDeliveryInfo,
  SITStatusOriginWithoutCustomerDeliveryInfo,
  SITShipment,
  SITStatusWithPastSITOriginServiceItem,
  SITStatusWithPastSITServiceItems,
  SITStatusWithPastSITServiceItemsDeparted,
  SITExtensionsWithComments,
  SITExtensionPending,
  SITExtensionDenied,
  SITStatusExpired,
  SITStatusShowConvert,
  SITStatusDontShowConvert,
} from './ShipmentSITDisplayTestParams';

import { MockProviders } from 'testUtils';
import { permissionTypes } from 'constants/permissions';

describe('ShipmentSITDisplay', () => {
  it('renders the Shipment SIT Extensions', async () => {
    render(
      <MockProviders>
        <ShipmentSITDisplay sitExtensions={SITExtensions} sitStatus={SITStatusOrigin} shipment={SITShipment} />
      </MockProviders>,
    );
    expect(screen.getByText('SIT (STORAGE IN TRANSIT)')).toBeTruthy();
    const sitStatusTable = await screen.findByTestId('sitStatusTable');
    expect(sitStatusTable).toBeInTheDocument();
    expect(within(sitStatusTable).getByText('Total days of SIT approved')).toBeInTheDocument();
    expect(within(sitStatusTable).getByText('270')).toBeInTheDocument();
    expect(within(sitStatusTable).getByText('Total days used')).toBeInTheDocument();
    expect(within(sitStatusTable).getByText('45')).toBeInTheDocument();
    expect(within(sitStatusTable).getByText('Total days remaining')).toBeInTheDocument();
    expect(within(sitStatusTable).getByText('60')).toBeInTheDocument();

    expect(screen.getByText('Current location: origin SIT')).toBeInTheDocument();

    expect(screen.getByText('Total days in origin SIT')).toBeInTheDocument();
    expect(screen.getByText(`13 Aug 2021`)).toBeInTheDocument();

    expect(await screen.queryByText('Office remarks:')).toBeFalsy();
  });

  it('renders the Shipment SIT at Destination, no previous SIT', async () => {
    render(
      <MockProviders>
        <ShipmentSITDisplay sitStatus={SITStatusDestination} shipment={SITShipment} />
      </MockProviders>,
    );

    expect(screen.getByText('Current location: destination SIT')).toBeInTheDocument();
    expect(screen.getByText('Total days in destination SIT')).toBeInTheDocument();
    expect(screen.getByText('15')).toBeInTheDocument();
    const sitStartAndEndTable = await screen.findByTestId('sitStartAndEndTable');
    expect(sitStartAndEndTable).toBeInTheDocument();
    expect(within(sitStartAndEndTable).getByText('Calculated total SIT days')).toBeInTheDocument();
    expect(within(sitStartAndEndTable).getByText('45')).toBeInTheDocument();
  });

  it('renders the Shipment SIT at Origin, with customer delivery info', async () => {
    render(
      <MockProviders>
        <ShipmentSITDisplay sitStatus={SITStatusOrigin} shipment={SITShipment} />
      </MockProviders>,
    );

    expect(screen.getByText('Customer delivery request')).toBeInTheDocument();
    expect(screen.getByText('Customer contact date')).toBeInTheDocument();
    expect(screen.getByText('26 Aug 2021')).toBeInTheDocument();
    expect(screen.getByText('Requested delivery date')).toBeInTheDocument();
    expect(screen.getByText('30 Aug 2021')).toBeInTheDocument();
  });

  it('renders the Shipment SIT at Destination, with customer delivery info', async () => {
    render(
      <MockProviders>
        <ShipmentSITDisplay sitStatus={SITStatusDestination} shipment={SITShipment} />
      </MockProviders>,
    );

    expect(screen.getByText('Customer delivery request')).toBeInTheDocument();
    expect(screen.getByText('Customer contact date')).toBeInTheDocument();
    expect(screen.getByText('26 Aug 2021')).toBeInTheDocument();
    expect(screen.getByText('Requested delivery date')).toBeInTheDocument();
    expect(screen.getByText('30 Aug 2021')).toBeInTheDocument();
    expect(screen.getByText('SIT departure date')).toBeInTheDocument();
  });

  it('renders the Shipment SIT at Destination, without customer delivery info', async () => {
    render(
      <MockProviders>
        <ShipmentSITDisplay sitStatus={SITStatusDestinationWithoutCustomerDeliveryInfo} shipment={SITShipment} />
      </MockProviders>,
    );

    expect(screen.getByText('Customer delivery request')).toBeInTheDocument();
    expect(screen.getByText('Customer contact date')).toBeInTheDocument();
    expect(screen.getByText('Requested delivery date')).toBeInTheDocument();
    expect(screen.getAllByText('—')).toHaveLength(3);
  });

  it('renders the Shipment SIT at Origin, without customer delivery info', async () => {
    render(
      <MockProviders>
        <ShipmentSITDisplay sitStatus={SITStatusOriginWithoutCustomerDeliveryInfo} shipment={SITShipment} />
      </MockProviders>,
    );

    expect(screen.getByText('Customer delivery request')).toBeInTheDocument();
    expect(screen.getByText('Customer contact date')).toBeInTheDocument();
    expect(screen.getByText('Requested delivery date')).toBeInTheDocument();
    expect(screen.getAllByText('—')).toHaveLength(3);
  });

  it('renders the Shipment SIT at Destination, previous destination SIT', async () => {
    render(
      <MockProviders>
        <ShipmentSITDisplay sitStatus={SITStatusWithPastSITOriginServiceItem} shipment={SITShipment} />
      </MockProviders>,
    );

    expect(screen.getByText('Previously used SIT')).toBeInTheDocument();
    expect(screen.getByText(`30 days at origin (24 Jul 2021 - 23 Aug 2021)`)).toBeInTheDocument();
  });

  it('renders the Shipment SIT at Destination, multiple previous SIT', async () => {
    render(
      <MockProviders>
        <ShipmentSITDisplay sitStatus={SITStatusWithPastSITServiceItems} shipment={SITShipment} />
      </MockProviders>,
    );
    expect(screen.getByText('Previously used SIT')).toBeInTheDocument();
    expect(screen.getByText(`30 days at origin (24 Jul 2021 - 23 Aug 2021)`)).toBeInTheDocument();
    expect(screen.getByText(`21 days at destination (03 Sep 2021 - 24 Sep 2021)`)).toBeInTheDocument();
  });

  it('renders with no current or future sit and multiple departed SIT', async () => {
    render(
      <MockProviders>
        <ShipmentSITDisplay sitStatus={SITStatusWithPastSITServiceItemsDeparted} shipment={SITShipment} />
      </MockProviders>,
    );
    expect(screen.getByText('Total days of SIT approved')).toBeInTheDocument();
    expect(screen.getByText('Total days used')).toBeInTheDocument();
    expect(screen.getByText('Total days remaining')).toBeInTheDocument();

    expect(screen.queryByText('Current location')).not.toBeInTheDocument();
    expect(screen.queryByText('SIT start date')).not.toBeInTheDocument();
    expect(screen.queryByText('SIT authorized end date')).not.toBeInTheDocument();
    expect(screen.queryByText('Calculated total SIT days')).not.toBeInTheDocument();

    expect(screen.getByText('Previously used SIT')).toBeInTheDocument();
    expect(screen.getByText(`30 days at origin (24 Jul 2021 - 23 Aug 2021)`)).toBeInTheDocument();
    expect(screen.getByText(`21 days at destination (03 Sep 2021 - 24 Sep 2021)`)).toBeInTheDocument();
  });

  it('renders the approved Shipment SIT Extensions', async () => {
    render(
      <MockProviders>
        <ShipmentSITDisplay sitExtensions={SITExtensions} sitStatus={SITStatusDestination} shipment={SITShipment} />
      </MockProviders>,
    );
    expect(screen.getByText('SIT history')).toBeInTheDocument();
    expect(screen.getByText('Total days of SIT approved: 270')).toBeInTheDocument();
    expect(screen.getByText('updated on 13 Sep 2021')).toBeInTheDocument();
    expect(screen.getByText('Serious illness of the member')).toBeInTheDocument();
  });

  it('renders the approved Shipment SIT Extensions with comments', async () => {
    render(
      <MockProviders>
        <ShipmentSITDisplay
          sitExtensions={SITExtensionsWithComments}
          sitStatus={SITStatusDestination}
          shipment={SITShipment}
        />
        ,
      </MockProviders>,
    );
    expect(screen.getByText('SIT (STORAGE IN TRANSIT)')).toBeInTheDocument();

    expect(screen.getByText('Office remarks:')).toBeInTheDocument();
    expect(screen.getByText('The customer requested an extension.')).toBeInTheDocument();
    expect(screen.getByText('Contractor remarks:')).toBeInTheDocument();
    expect(
      screen.getByText('The service member is unable to move into their new home at the expected time.'),
    ).toBeInTheDocument();
  });

  it('renders the denied Shipment SIT Extensions', async () => {
    render(
      <MockProviders>
        <ShipmentSITDisplay
          sitExtensions={SITExtensionDenied}
          sitStatus={SITStatusDestination}
          shipment={SITShipment}
        />
      </MockProviders>,
    );
    expect(screen.getByText('SIT history')).toBeInTheDocument();
    expect(screen.getByText('Total days of SIT approved: 270')).toBeInTheDocument();
    expect(screen.getByText('updated on 13 Sep 2021')).toBeInTheDocument();
    expect(screen.getByText('Serious illness of the member')).toBeInTheDocument();
  });

  it('omits SIT Extension history when there is only a pending SIT Extension', async () => {
    render(
      <MockProviders>
        <ShipmentSITDisplay
          sitExtensions={SITExtensionPending}
          sitStatus={SITStatusDestination}
          shipment={SITShipment}
        />
      </MockProviders>,
    );

    expect(screen.queryByText('SIT extensions')).not.toBeInTheDocument();
  });

  it('renders the future SIT', async () => {
    render(
      <MockProviders>
        <ShipmentSITDisplay shipment={futureSITShipment} sitStatus={futureSITStatus} />
      </MockProviders>,
    );
    const sitStatusTable = await screen.findByTestId('sitStatusTable');
    expect(sitStatusTable).toBeInTheDocument();
    expect(within(sitStatusTable).getByText('Total days of SIT approved')).toBeInTheDocument();
    expect(within(sitStatusTable).getByText('Total days remaining')).toBeInTheDocument();
    const daysApprovedAndRemaining = within(sitStatusTable).getAllByText('365');
    expect(daysApprovedAndRemaining).toHaveLength(1);
    const sitStartAndEndTable = await screen.findByTestId('sitStartAndEndTable');
    expect(sitStartAndEndTable).toBeInTheDocument();
    expect(within(sitStartAndEndTable).queryByText('Current location')).not.toBeInTheDocument();
    expect(within(sitStartAndEndTable).getByText('SIT start date')).toBeInTheDocument();
    expect(within(sitStartAndEndTable).getByText('SIT authorized end date')).toBeInTheDocument();
    expect(within(sitStartAndEndTable).getByText('Calculated total SIT days')).toBeInTheDocument();
    expect(within(sitStartAndEndTable).getByText('0')).toBeInTheDocument();
    const sitDaysAtCurrentLocation = await screen.findByTestId('sitDaysAtCurrentLocation');
    expect(sitDaysAtCurrentLocation).toBeInTheDocument();
    expect(within(sitDaysAtCurrentLocation).getByText('Total days in origin SIT')).toBeInTheDocument();
    expect(within(sitDaysAtCurrentLocation).getByText('0')).toBeInTheDocument();
  });

  it('calls SIT extension callback when button clicked', async () => {
    const onClick = jest.fn();
    const OpenModalButton = (
      <button type="button" onClick={() => onClick()}>
        Edit
      </button>
    );
    render(
      <MockProviders permissions={[permissionTypes.updateSITExtension]}>
        <ShipmentSITDisplay
          sitExtensions={SITExtensions}
          sitStatus={SITStatusDestination}
          shipment={SITShipment}
          openModalButton={OpenModalButton}
        />
      </MockProviders>,
    );

    const editButton = screen.getByRole('button', { name: 'Edit' });

    await userEvent.click(editButton);

    await waitFor(() => {
      expect(onClick).toHaveBeenCalledTimes(1);
    });
  });

  it('show convert SIT To Customer Expense callback when show convert is true', async () => {
    const onClick = jest.fn();
    const OpenConvertModalButton = (
      <button type="button" onClick={() => onClick()}>
        Convert to customer expense
      </button>
    );
    render(
      <MockProviders permissions={[permissionTypes.updateSITExtension]}>
        <ShipmentSITDisplay
          sitStatus={SITStatusShowConvert}
          shipment={SITShipment}
          openConvertModalButton={OpenConvertModalButton}
        />
      </MockProviders>,
    );

    const convertButton = screen.getByRole('button', { name: 'Convert to customer expense' });

    await userEvent.click(convertButton);

    await waitFor(() => {
      expect(onClick).toHaveBeenCalledTimes(1);
    });
  });

  it('hide convert SIT To Customer Expense button when show button is false', async () => {
    const onClick = jest.fn();
    const OpenConvertModalButton = (
      <button type="button" onClick={() => onClick()}>
        Convert to customer expense
      </button>
    );
    render(
      <MockProviders permissions={[permissionTypes.updateSITExtension]}>
        <ShipmentSITDisplay
          sitStatus={SITStatusDontShowConvert}
          shipment={SITShipment}
          openConvertModalButton={OpenConvertModalButton}
        />
      </MockProviders>,
    );

    expect(screen.queryByRole('button', { name: 'Convert to customer expense' })).not.toBeInTheDocument();
  });

  it('hides review pending SIT Extension button when hide prop is true', async () => {
    render(
      <MockProviders permissions={[permissionTypes.createSITExtension]}>
        <ShipmentSITDisplay
          sitExtensions={SITExtensionPending}
          sitStatus={SITStatusDestination}
          shipment={SITShipment}
        />
      </MockProviders>,
    );

    expect(screen.queryByRole('button', { name: 'View request' })).not.toBeInTheDocument();
  });

  it('hides submit new SIT Extension button when hide prop is true', async () => {
    render(
      <MockProviders permissions={[permissionTypes.updateSITExtension]}>
        <ShipmentSITDisplay sitExtensions={SITExtensions} sitStatus={SITStatusDestination} shipment={SITShipment} />
      </MockProviders>,
    );

    expect(screen.queryByRole('button', { name: 'Edit' })).not.toBeInTheDocument();
  });

  it('View request button is hidden when user does not have permissions', async () => {
    render(
      <MockProviders>
        <ShipmentSITDisplay sitExtensions={SITExtensions} sitStatus={SITStatusDestination} shipment={SITShipment} />
      </MockProviders>,
    );

    expect(screen.queryByRole('button', { name: 'View request' })).not.toBeInTheDocument();
  });

  it('Edit button is hidden when user does not have permissions', async () => {
    render(
      <MockProviders>
        <ShipmentSITDisplay sitExtensions={SITExtensions} sitStatus={SITStatusDestination} shipment={SITShipment} />
      </MockProviders>,
    );

    expect(screen.queryByRole('button', { name: 'Edit' })).not.toBeInTheDocument();
  });
  it('shows Expired when the used days is greater than the approved days', async () => {
    render(
      <MockProviders>
        <ShipmentSITDisplay sitExtensions={SITExtensions} sitStatus={SITStatusExpired} shipment={SITShipment} />
      </MockProviders>,
    );
    expect(screen.getByText('Expired')).toBeInTheDocument();
  });
  it('renders the Shipment SIT at Origin, with current SIT authorized end date', async () => {
    render(
      <MockProviders>
        <ShipmentSITDisplay sitStatus={SITStatusOriginAuthorized} shipment={SITShipment} />
      </MockProviders>,
    );

    const sitStartAndEndTable = await screen.findByTestId('sitStartAndEndTable');
    expect(sitStartAndEndTable).toBeInTheDocument();
    expect(within(sitStartAndEndTable).getByText('SIT authorized end date')).toBeInTheDocument();
    expect(within(sitStartAndEndTable).getByText('28 Aug 2021')).toBeInTheDocument();
  });
});

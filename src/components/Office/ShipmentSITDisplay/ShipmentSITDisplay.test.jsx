import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import ShipmentSITDisplay from './ShipmentSITDisplay';
import {
  SITExtensions,
  SITStatusOrigin,
  SITStatusDestination,
  SITShipment,
  SITStatusWithPastSITOriginServiceItem,
  SITStatusWithPastSITServiceItems,
  SITExtensionsWithComments,
  SITExtensionPending,
  SITExtensionDenied,
} from './ShipmentSITDisplayTestParams';

import { MockProviders } from 'testUtils';
import { permissionTypes } from 'constants/permissions';

describe('ShipmentSITDisplay', () => {
  it('renders the Shipment SIT Extensions', async () => {
    render(<ShipmentSITDisplay sitExtensions={SITExtensions} sitStatus={SITStatusOrigin} shipment={SITShipment} />);
    expect(screen.getByText('SIT (STORAGE IN TRANSIT)')).toBeTruthy();

    expect(screen.getByText('270 authorized')).toBeInTheDocument();
    expect(screen.getByText('45 used')).toBeInTheDocument();
    expect(screen.getByText('60 remaining')).toBeInTheDocument();

    expect(screen.getByText('Current location: origin')).toBeInTheDocument();

    expect(screen.getByText('Days in origin SIT')).toBeInTheDocument();
    expect(screen.getByText('45')).toBeInTheDocument();
    expect(screen.getByText(`13 Aug 2021`)).toBeInTheDocument();

    expect(await screen.queryByText('Office remarks:')).toBeFalsy();
  });

  it('renders the Shipment SIT at Destination, no previous SIT', async () => {
    render(<ShipmentSITDisplay sitStatus={SITStatusDestination} shipment={SITShipment} />);

    expect(screen.getByText('Current location: destination')).toBeInTheDocument();
    expect(screen.getByText('Days in destination SIT')).toBeInTheDocument();
  });

  it('renders the Shipment SIT at Destination, previous origin SIT', async () => {
    render(<ShipmentSITDisplay sitStatus={SITStatusWithPastSITOriginServiceItem} shipment={SITShipment} />);

    expect(screen.getByText('Previously used SIT')).toBeInTheDocument();
    expect(await screen.getByText(`30 days at origin (24 Jul 2021 - 23 Aug 2021)`)).toBeInTheDocument();
  });

  it('renders the Shipment SIT at Destination, multiple previous SIT', async () => {
    render(<ShipmentSITDisplay sitStatus={SITStatusWithPastSITServiceItems} shipment={SITShipment} />);
    expect(screen.getByText('Previously used SIT')).toBeInTheDocument();
    expect(await screen.getByText(`30 days at origin (24 Jul 2021 - 23 Aug 2021)`)).toBeInTheDocument();
    expect(await screen.getByText(`21 days at destination (03 Sep 2021 - 24 Sep 2021)`)).toBeInTheDocument();
  });

  it('renders the approved Shipment SIT Extensions', async () => {
    render(
      <ShipmentSITDisplay sitExtensions={SITExtensions} sitStatus={SITStatusDestination} shipment={SITShipment} />,
    );
    expect(screen.getByText('SIT extensions')).toBeInTheDocument();
    expect(screen.getByText('30 days added')).toBeInTheDocument();
    expect(screen.getByText('on 13 Sep 2021')).toBeInTheDocument();
    expect(screen.getByText('Serious illness of the member')).toBeInTheDocument();
  });

  it('renders the approved Shipment SIT Extensions with comments', async () => {
    render(
      <ShipmentSITDisplay
        sitExtensions={SITExtensionsWithComments}
        sitStatus={SITStatusDestination}
        shipment={SITShipment}
      />,
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
      <ShipmentSITDisplay sitExtensions={SITExtensionDenied} sitStatus={SITStatusDestination} shipment={SITShipment} />,
    );
    expect(screen.getByText('SIT extensions')).toBeInTheDocument();
    expect(screen.getByText('0 days added')).toBeInTheDocument();
    expect(screen.getByText('on 13 Sep 2021 — request rejected')).toBeInTheDocument();
    expect(screen.getByText('Serious illness of the member')).toBeInTheDocument();
  });

  it('omits SIT Extension history when there is only a pending SIT Extension', async () => {
    render(
      <ShipmentSITDisplay
        sitExtensions={SITExtensionPending}
        sitStatus={SITStatusDestination}
        shipment={SITShipment}
      />,
    );

    expect(await screen.queryByText('SIT extensions')).not.toBeInTheDocument();
  });

  it('calls review SIT extension callback when button is clicked', async () => {
    const showReviewSITExtension = jest.fn();
    render(
      <MockProviders permissions={[permissionTypes.createSITExtension]}>
        <ShipmentSITDisplay
          sitExtensions={SITExtensionPending}
          sitStatus={SITStatusDestination}
          shipment={SITShipment}
          showReviewSITExtension={showReviewSITExtension}
        />
      </MockProviders>,
    );

    const reviewButton = screen.getByRole('button', { name: 'View request' });

    userEvent.click(reviewButton);

    await waitFor(() => {
      expect(showReviewSITExtension).toHaveBeenCalledWith(true);
    });
  });

  it('calls submit SIT extension callback when button is clicked', async () => {
    const showSubmitSITExtension = jest.fn();
    render(
      <MockProviders permissions={[permissionTypes.updateSITExtension]}>
        <ShipmentSITDisplay
          sitExtensions={SITExtensions}
          sitStatus={SITStatusDestination}
          shipment={SITShipment}
          showSubmitSITExtension={showSubmitSITExtension}
        />
      </MockProviders>,
    );

    const reviewButton = screen.getByRole('button', { name: 'Edit' });

    userEvent.click(reviewButton);

    await waitFor(() => {
      expect(showSubmitSITExtension).toHaveBeenCalledWith(true);
    });
  });

  it('hides review pending SIT Extension button when hide prop is true', async () => {
    render(
      <MockProviders permissions={[permissionTypes.createSITExtension]}>
        <ShipmentSITDisplay
          sitExtensions={SITExtensionPending}
          sitStatus={SITStatusDestination}
          shipment={SITShipment}
          hideSITExtensionAction
        />
      </MockProviders>,
    );

    expect(await screen.queryByRole('button', { name: 'View request' })).not.toBeInTheDocument();
  });

  it('hides submit new SIT Extension button when hide prop is true', async () => {
    render(
      <MockProviders permissions={[permissionTypes.updateSITExtension]}>
        <ShipmentSITDisplay
          sitExtensions={SITExtensions}
          sitStatus={SITStatusDestination}
          shipment={SITShipment}
          hideSITExtensionAction
        />
      </MockProviders>,
    );

    expect(await screen.queryByRole('button', { name: 'Edit' })).not.toBeInTheDocument();
  });

  it('View request button is hidden when user does not have permissions', async () => {
    render(
      <ShipmentSITDisplay sitExtensions={SITExtensions} sitStatus={SITStatusDestination} shipment={SITShipment} />,
    );

    expect(await screen.queryByRole('button', { name: 'View request' })).not.toBeInTheDocument();
  });

  it('Edit button is hidden when user does not have permissions', async () => {
    render(
      <ShipmentSITDisplay sitExtensions={SITExtensions} sitStatus={SITStatusDestination} shipment={SITShipment} />,
    );

    expect(await screen.queryByRole('button', { name: 'Edit' })).not.toBeInTheDocument();
  });
});

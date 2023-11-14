import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import { IncompleteShipmentModal } from 'components/Customer/Review/IncompleteShipmentModal/IncompleteShipmentModal';

let onClose;
const shipmentLabel = 'PPM_1';
const moveCodeLabel = 'D889F48D';
const shipmentType = 'PPM';

beforeEach(() => {
  onClose = jest.fn();
});

describe('IncompleteShipmentModal', () => {
  it('verify incompleteShipmentModal display', async () => {
    render(
      <IncompleteShipmentModal
        closeModal={onClose}
        shipmentLabel={shipmentLabel}
        shipmentMoveCode={moveCodeLabel}
        shipmentType={shipmentType}
      />,
    );

    expect(await screen.findByRole('heading', { level: 3, name: 'INCOMPLETE SHIPMENT' })).toBeInTheDocument();

    const keepButton = await screen.findByRole('button', { name: 'OK' });
    await userEvent.click(keepButton);
    expect(onClose).toHaveBeenCalledTimes(1);

    expect(screen.getByText(/PPM_1: #D889F48D/, { selector: 'b' })).toBeInTheDocument();

    // make sure shipment type display is not hardcoded
    expect(screen.getByText(/PPM shipment/, { selector: 'p' })).toBeInTheDocument();
    expect(screen.getByText(/PPM information/, { selector: 'p' })).toBeInTheDocument();
  });
});

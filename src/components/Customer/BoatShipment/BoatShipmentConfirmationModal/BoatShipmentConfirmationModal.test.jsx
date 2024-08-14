import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import { BoatShipmentConfirmationModal } from './BoatShipmentConfirmationModal';

let closeModal;
let handleConfirmationContinue;
let handleConfirmationRedirect;
let handleConfirmationDeleteAndRedirect;

beforeEach(() => {
  closeModal = jest.fn();
  handleConfirmationContinue = jest.fn();
  handleConfirmationRedirect = jest.fn();
  handleConfirmationDeleteAndRedirect = jest.fn();
});

describe('BoatShipmentConfirmationModal', () => {
  it('renders the component with default settings', async () => {
    render(
      <BoatShipmentConfirmationModal
        isDimensionsMeetReq
        boatShipmentType="TOW_AWAY"
        closeModal={closeModal}
        handleConfirmationContinue={handleConfirmationContinue}
        handleConfirmationRedirect={handleConfirmationRedirect}
        handleConfirmationDeleteAndRedirect={handleConfirmationDeleteAndRedirect}
        isEditPage={false}
      />,
    );

    expect(await screen.findByRole('heading', { level: 3, name: 'Boat Tow-Away (BTA)' })).toBeInTheDocument();
  });

  it('closes the modal when the back button is clicked', async () => {
    render(
      <BoatShipmentConfirmationModal
        isDimensionsMeetReq
        boatShipmentType="TOW_AWAY"
        closeModal={closeModal}
        handleConfirmationContinue={handleConfirmationContinue}
        handleConfirmationRedirect={handleConfirmationRedirect}
        handleConfirmationDeleteAndRedirect={handleConfirmationDeleteAndRedirect}
        isEditPage={false}
      />,
    );

    const backButton = await screen.findByTestId('boatConfirmationBack');

    await userEvent.click(backButton);

    expect(closeModal).toHaveBeenCalledTimes(1);
  });

  it('calls handleConfirmationContinue when the continue button is clicked and dimensions meet the requirements', async () => {
    render(
      <BoatShipmentConfirmationModal
        isDimensionsMeetReq
        boatShipmentType="HAUL_AWAY"
        closeModal={closeModal}
        handleConfirmationContinue={handleConfirmationContinue}
        handleConfirmationRedirect={handleConfirmationRedirect}
        handleConfirmationDeleteAndRedirect={handleConfirmationDeleteAndRedirect}
        isEditPage={false}
      />,
    );

    const continueButton = await screen.findByTestId('boatConfirmationContinue');

    await userEvent.click(continueButton);

    expect(handleConfirmationContinue).toHaveBeenCalledTimes(1);
  });

  it('calls handleConfirmationRedirect when the continue button is clicked and dimensions do not meet the requirements', async () => {
    render(
      <BoatShipmentConfirmationModal
        isDimensionsMeetReq={false}
        boatShipmentType="HAUL_AWAY"
        closeModal={closeModal}
        handleConfirmationContinue={handleConfirmationContinue}
        handleConfirmationRedirect={handleConfirmationRedirect}
        handleConfirmationDeleteAndRedirect={handleConfirmationDeleteAndRedirect}
        isEditPage={false}
      />,
    );

    const continueButton = await screen.findByTestId('boatConfirmationContinue');

    await userEvent.click(continueButton);

    expect(handleConfirmationRedirect).toHaveBeenCalledTimes(1);
  });

  it('calls handleConfirmationDeleteAndRedirect when delete & continue button is clicked on the edit page', async () => {
    render(
      <BoatShipmentConfirmationModal
        isDimensionsMeetReq={false}
        boatShipmentType=""
        closeModal={closeModal}
        handleConfirmationContinue={handleConfirmationContinue}
        handleConfirmationRedirect={handleConfirmationRedirect}
        handleConfirmationDeleteAndRedirect={handleConfirmationDeleteAndRedirect}
        isEditPage
      />,
    );

    const deleteContinueButton = await screen.findByTestId('boatConfirmationContinue');

    await userEvent.click(deleteContinueButton);

    expect(handleConfirmationDeleteAndRedirect).toHaveBeenCalledTimes(1);
  });
});

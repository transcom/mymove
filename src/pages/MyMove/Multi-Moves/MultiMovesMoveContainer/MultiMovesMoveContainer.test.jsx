import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import '@testing-library/jest-dom/extend-expect'; // For expect assertions

import { mockMovesPCS, mockMovesPPMWithAdvanceOptions } from '../MultiMovesTestData';

import MultiMovesMoveContainer from './MultiMovesMoveContainer';

import { MockProviders } from 'testUtils';

describe('MultiMovesMoveContainer', () => {
  const mockCurrentMoves = mockMovesPCS.currentMove;
  const mockPreviousMoves = mockMovesPCS.previousMoves;

  it('renders current move list correctly', () => {
    render(
      <MockProviders>
        <MultiMovesMoveContainer moves={mockCurrentMoves} />
      </MockProviders>,
    );

    expect(screen.getByTestId('move-info-container')).toBeInTheDocument();
    expect(screen.getByText('#MOVECO')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Go to Move' })).toBeInTheDocument();
  });

  it('renders previous move list correctly', () => {
    render(
      <MockProviders>
        <MultiMovesMoveContainer moves={mockPreviousMoves} />
      </MockProviders>,
    );

    expect(screen.queryByText('#SAMPLE')).toBeInTheDocument();
    expect(screen.queryByText('#EXAMPL')).toBeInTheDocument();
  });

  it('expands and collapses moves correctly', () => {
    render(
      <MockProviders>
        <MultiMovesMoveContainer moves={mockCurrentMoves} />
      </MockProviders>,
    );

    // Initially, the move details should not be visible
    expect(screen.queryByText('Shipment')).not.toBeInTheDocument();
    fireEvent.click(screen.getByTestId('expand-icon'));

    // Now, the move details should be visible
    expect(screen.getByText('Shipments')).toBeInTheDocument();

    fireEvent.click(screen.getByTestId('expand-icon'));

    // The move details should be hidden again
    expect(screen.queryByText('Shipments')).not.toBeInTheDocument();
  });

  it('renders Go to Move & Download buttons for current move', () => {
    render(
      <MockProviders>
        <MultiMovesMoveContainer moves={mockCurrentMoves} />
      </MockProviders>,
    );

    expect(screen.getByTestId('headerBtns')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Go to Move' })).toBeInTheDocument();
  });

  it('renders Go to Move & Download buttons for previous moves exceeding one', () => {
    render(
      <MockProviders>
        <MultiMovesMoveContainer moves={mockPreviousMoves} />
      </MockProviders>,
    );

    // Check for the container that holds the buttons - there should be 2
    const headerBtnsElements = screen.getAllByTestId('headerBtns');
    expect(headerBtnsElements).toHaveLength(2);

    // Check for Go to Move buttons - there should be 2
    const goToMoveButtons = screen.getAllByRole('button', { name: 'Go to Move' });
    expect(goToMoveButtons).toHaveLength(2);
  });

  it('renders correct shipment type headings for PPMs', async () => {
    render(
      <MockProviders>
        <MultiMovesMoveContainer moves={mockCurrentMoves} />
      </MockProviders>,
    );

    // expand the shipments
    fireEvent.click(screen.getByTestId('expand-icon'));
    await waitFor(() => expect(screen.getByText('Shipments')).toBeInTheDocument());

    const headings = screen.queryAllByRole('heading', { level: 4 });
    expect(headings).toHaveLength(8);
    expect(headings[0]).toHaveTextContent('Household Goods');
    expect(headings[1]).toHaveTextContent('Personally Procured Move');
    expect(headings[2]).toHaveTextContent('Household Goods NTS');
    expect(headings[3]).toHaveTextContent('Household Goods NTSR');
    expect(headings[4]).toHaveTextContent('Mobile Home');
    expect(headings[5]).toHaveTextContent('Boat Haul Away');
    expect(headings[6]).toHaveTextContent('Boat Tow Away');
    expect(headings[7]).toHaveTextContent('Unaccompanied Baggage');
  });

  it('renders appropriate ppm download options', async () => {
    const { currentMove } = mockMovesPPMWithAdvanceOptions;
    render(
      <MockProviders>
        <MultiMovesMoveContainer moves={currentMove} />
      </MockProviders>,
    );

    // expand the shipments
    fireEvent.click(screen.getByTestId('expand-icon'));
    await waitFor(() => expect(screen.getByText('Shipments')).toBeInTheDocument());

    // Should be four shipments
    const shipments = screen.getAllByTestId('shipment-container');
    expect(shipments).toHaveLength(4);

    // Check that there are three download buttons (one shipment should not have a download btn)
    const downloadButtons = screen.getAllByRole('button', { name: 'Download' });
    expect(downloadButtons).toHaveLength(3);

    // the first shipment should have both options
    fireEvent.click(downloadButtons[0]);
    await waitFor(() => expect(screen.queryByText('AOA Packet')).toBeInTheDocument());
    await waitFor(() => expect(screen.queryByText('PPM Packet')).toBeInTheDocument());
    // close it & verify the dropdown is closed
    fireEvent.click(downloadButtons[0]);
    expect(screen.queryByText('AOA Packet')).not.toBeInTheDocument();
    expect(screen.queryByText('PPM Packet')).not.toBeInTheDocument();

    // the second shipment should only have the PPM packet option
    fireEvent.click(downloadButtons[1]);
    await waitFor(() => expect(screen.queryByText('AOA Packet')).not.toBeInTheDocument());
    await waitFor(() => expect(screen.queryByText('PPM Packet')).toBeInTheDocument());
    // close it & verify the dropdown is closed
    fireEvent.click(downloadButtons[1]);
    expect(screen.queryByText('AOA Packet')).not.toBeInTheDocument();
    expect(screen.queryByText('PPM Packet')).not.toBeInTheDocument();

    // the third shipment should only have the AOA packet option
    fireEvent.click(downloadButtons[2]);
    await waitFor(() => expect(screen.queryByText('AOA Packet')).toBeInTheDocument());
    await waitFor(() => expect(screen.queryByText('PPM Packet')).not.toBeInTheDocument());
    // close it & verify the dropdown is closed
    fireEvent.click(downloadButtons[2]);
    expect(screen.queryByText('AOA Packet')).not.toBeInTheDocument();
    expect(screen.queryByText('PPM Packet')).not.toBeInTheDocument();
  });

  it('renders correct shipment type headings for PPMs', async () => {
    const { currentMove } = mockMovesPPMWithAdvanceOptions;
    render(
      <MockProviders>
        <MultiMovesMoveContainer moves={currentMove} />
      </MockProviders>,
    );

    // expand the shipments
    fireEvent.click(screen.getByTestId('expand-icon'));
    await waitFor(() => expect(screen.getByText('Shipments')).toBeInTheDocument());

    const headings = screen.queryAllByRole('heading', { level: 4 });
    expect(headings).toHaveLength(4);
    expect(headings[0]).toHaveTextContent('Personally Procured Move');
    expect(headings[1]).toHaveTextContent('Personally Procured Move');
    expect(headings[2]).toHaveTextContent('Personally Procured Move');
    expect(headings[3]).toHaveTextContent('Personally Procured Move');
  });

  it('renders Canceled move label', () => {
    render(
      <MockProviders>
        <MultiMovesMoveContainer moves={mockMovesPCS.canceledMove} />
      </MockProviders>,
    );

    expect(screen.getByText('Canceled')).toBeInTheDocument();
  });
});

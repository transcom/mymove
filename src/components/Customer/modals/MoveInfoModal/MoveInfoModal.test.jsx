import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';

import '@testing-library/jest-dom';
import { MoveInfoModal } from './MoveInfoModal';

describe('MoveInfoModal', () => {
  const mockCloseModal = jest.fn();

  // renders the MoveInfoModal with default props
  test('renders MoveInfoModal with default props', () => {
    render(<MoveInfoModal closeModal={mockCloseModal} />);
    // checks if MoveInfoModal heading is present
    expect(screen.getByTestId('moveInfoModalHeading')).toHaveTextContent('More info about shipments');
    // checks if HHG section is rendered
    expect(screen.getByText(/Professional movers pack and ship your things/)).toBeInTheDocument();
    expect(screen.getByTestId('hhgSubHeading')).toBeInTheDocument();
    expect(screen.getByTestId('hhgProsList')).toBeInTheDocument();
    expect(screen.getByTestId('hhgConsList')).toBeInTheDocument();
  });

  // checks if PPM section is rendered when enablePPM is true
  test('renders PPM section when enablePPM is true', () => {
    render(<MoveInfoModal closeModal={mockCloseModal} enablePPM />);
    expect(screen.getByText(/PPM: You get your things packed and moved/)).toBeInTheDocument();
    expect(screen.getByTestId('ppmSubHeading')).toBeInTheDocument();
    expect(screen.getByTestId('ppmProsList')).toBeInTheDocument();
    expect(screen.getByTestId('ppmConsList')).toBeInTheDocument();
  });

  test('does not render PPM section when enablePPM is false', () => {
    render(<MoveInfoModal closeModal={mockCloseModal} enablePPM={false} />);
    expect(screen.queryByText(/PPM: You get your things packed and moved/)).toBeNull();
    expect(screen.queryByTestId('ppmSubHeading')).not.toBeInTheDocument();
    expect(screen.queryByTestId('ppmProsList')).not.toBeInTheDocument();
    expect(screen.queryByTestId('ppmConsList')).not.toBeInTheDocument();
  });

  // tests UB section based on enableUB and hasOconusDutyLocation
  test('renders UB section when enableUB and hasOconusDutyLocation are true', () => {
    render(<MoveInfoModal closeModal={mockCloseModal} enableUB hasOconusDutyLocation />);
    expect(
      screen.getByText(/UB: Professional movers pack and ship your more essential personal property/),
    ).toBeInTheDocument();
    expect(screen.getByTestId('ubSubHeading')).toBeInTheDocument();
    expect(screen.getByTestId('ubProsList')).toBeInTheDocument();
    expect(screen.getByTestId('ubConsList')).toBeInTheDocument();
  });

  test('does not render UB section when enableUB is false', () => {
    render(<MoveInfoModal closeModal={mockCloseModal} enableUB={false} hasOconusDutyLocation />);

    expect(
      screen.queryByText(/UB: Professional movers pack and ship your more essential personal property/),
    ).toBeNull();
    expect(screen.queryByTestId('ubSubHeading')).not.toBeInTheDocument();
    expect(screen.queryByTestId('ubProsList')).not.toBeInTheDocument();
    expect(screen.queryByTestId('ubConsList')).not.toBeInTheDocument();
  });

  test('does not render UB section when hasOconusDutyLocation is false', () => {
    render(<MoveInfoModal closeModal={mockCloseModal} enableUB hasOconusDutyLocation={false} />);

    expect(
      screen.queryByText(/UB: Professional movers pack and ship your more essential personal property/),
    ).toBeNull();
    expect(screen.queryByTestId('ubSubHeading')).not.toBeInTheDocument();
    expect(screen.queryByTestId('ubProsList')).not.toBeInTheDocument();
    expect(screen.queryByTestId('ubConsList')).not.toBeInTheDocument();
  });

  // tests "Got it" button functionality
  test('calls closeModal function when "Got it" button is clicked', () => {
    render(<MoveInfoModal closeModal={mockCloseModal} />);

    const button = screen.getByRole('button', { name: /Got it/ });
    fireEvent.click(button);

    // Verify the closeModal function was called
    expect(mockCloseModal).toHaveBeenCalled();
  });
});

import React from 'react';
import { render } from '@testing-library/react';
import { renderHook, act } from '@testing-library/react-hooks';

import { Modal, connectModal, useModal } from './MigratedModal';

/** This is a straightforward port of the Modal component from React-USWDS 1.17
 *  into the MilMove project, as the component is being deprecated in USWDS 2.x. */

describe('Modal component', () => {
  it('renders without errors', () => {
    const { queryByText } = render(<Modal>My Modal</Modal>);
    expect(queryByText('My Modal')).toBeInTheDocument();
  });
});

const TestModal = () => <div>My Modal</div>;

describe('connectModal', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });
  const mockClose = jest.fn();

  describe('if isOpen is false', () => {
    it('does not render its children', () => {
      const TestConnectedModal = connectModal(TestModal);
      const { queryByText } = render(<TestConnectedModal isOpen={false} onClose={mockClose} />);
      expect(queryByText('My Modal')).not.toBeInTheDocument();
    });
  });

  describe('if isOpen is true', () => {
    it('renders its children', () => {
      const TestConnectedModal = connectModal(TestModal);
      const { queryByText } = render(<TestConnectedModal isOpen onClose={mockClose} />);
      expect(queryByText('My Modal')).toBeInTheDocument();
    });
  });
});

describe('useModal', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('provides state and functions for opening/closing a modal', () => {
    const { result } = renderHook(() => useModal());
    expect(result.current.isOpen).toEqual(false);
    expect(typeof result.current.openModal).toBe('function');
    expect(typeof result.current.closeModal).toBe('function');

    act(() => {
      result.current.openModal();
    });

    expect(result.current.isOpen).toEqual(true);

    act(() => {
      result.current.closeModal();
    });

    expect(result.current.isOpen).toEqual(false);
  });
});

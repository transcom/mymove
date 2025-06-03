import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';

import CertificationText from './CertificationText';

// Helper to mock a scrollable div
const mockScrollToBottom = (element) => {
  Object.defineProperty(element, 'scrollHeight', {
    value: 1000,
    writable: true,
  });
  Object.defineProperty(element, 'clientHeight', {
    value: 500,
    writable: true,
  });
  Object.defineProperty(element, 'scrollTop', {
    value: 500,
    writable: true,
  });
};

describe('CertificationText component', () => {
  it('renders markdown text', () => {
    const text = '# Hello World\nThis is **bold** text.';
    render(<CertificationText certificationText={text} />);

    expect(screen.getByText('Hello World')).toBeInTheDocument();
    expect(screen.getByText('bold')).toBeInTheDocument();
  });

  it('calls onScrollToBottom when scrolled to bottom', () => {
    const onScrollToBottomMock = jest.fn();
    render(<CertificationText certificationText="Sample certification" onScrollToBottom={onScrollToBottomMock} />);

    const box = screen.getByTestId('certificationTextBox');
    mockScrollToBottom(box);

    fireEvent.scroll(box, {
      target: {
        scrollTop: 500,
        scrollHeight: 1000,
        clientHeight: 500,
      },
    });

    expect(onScrollToBottomMock).toHaveBeenCalledWith(true);
  });

  it('does not call onScrollToBottom if not at bottom', () => {
    const onScrollToBottomMock = jest.fn();
    render(<CertificationText certificationText="Sample certification" onScrollToBottom={onScrollToBottomMock} />);

    const box = screen.getByTestId('certificationTextBox');

    // Mock scroll values directly on the box element
    Object.defineProperty(box, 'scrollTop', { value: 100, configurable: true });
    Object.defineProperty(box, 'scrollHeight', { value: 1000, configurable: true });
    Object.defineProperty(box, 'clientHeight', { value: 500, configurable: true });

    fireEvent.scroll(box); // No target needed now, it uses the box element's mocked props

    expect(onScrollToBottomMock).not.toHaveBeenCalled();
  });

  it('calls onScrollToBottom only once when scrolled to bottom multiple times', () => {
    const onScrollToBottomMock = jest.fn();
    render(<CertificationText certificationText="Sample certification" onScrollToBottom={onScrollToBottomMock} />);

    const box = screen.getByTestId('certificationTextBox');
    mockScrollToBottom(box);

    // First scroll - should trigger callback
    fireEvent.scroll(box, {
      target: {
        scrollTop: 500,
        scrollHeight: 1000,
        clientHeight: 500,
      },
    });

    // Scroll again - shouldn't trigger callback again
    fireEvent.scroll(box, {
      target: {
        scrollTop: 500,
        scrollHeight: 1000,
        clientHeight: 500,
      },
    });

    expect(onScrollToBottomMock).toHaveBeenCalledTimes(1);
  });
  it('renders nothing when certificationText is undefined', () => {
    render(<CertificationText />);

    const box = screen.getByTestId('certificationTextBox');
    expect(box).toBeInTheDocument();
    expect(box).toBeEmptyDOMElement();
  });
  it('does not call onScrollToBottom again if already scrolled to bottom', () => {
    const onScrollToBottomMock = jest.fn();
    render(<CertificationText certificationText="Sample" onScrollToBottom={onScrollToBottomMock} />);

    const box = screen.getByTestId('certificationTextBox');

    // First scroll to bottom
    Object.defineProperty(box, 'scrollTop', { value: 500, configurable: true });
    Object.defineProperty(box, 'scrollHeight', { value: 1000, configurable: true });
    Object.defineProperty(box, 'clientHeight', { value: 500, configurable: true });

    fireEvent.scroll(box);

    // Scroll again - should not trigger again
    fireEvent.scroll(box);

    expect(onScrollToBottomMock).toHaveBeenCalledTimes(1);
  });
  it('does not throw if onScrollToBottom is not provided', () => {
    render(<CertificationText certificationText="Just text" />);

    const box = screen.getByTestId('certificationTextBox');

    Object.defineProperty(box, 'scrollTop', { value: 500, configurable: true });
    Object.defineProperty(box, 'scrollHeight', { value: 1000, configurable: true });
    Object.defineProperty(box, 'clientHeight', { value: 500, configurable: true });

    expect(() => {
      fireEvent.scroll(box);
    }).not.toThrow();
  });
});

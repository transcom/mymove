import React from 'react';
import { fireEvent, render, screen } from '@testing-library/react';
import { Button } from '@trussworks/react-uswds';

import DocumentViewerSidebar from './DocumentViewerSidebar';

describe('DocumentViewerSidebar', () => {
  it('renders header, content and footer', () => {
    const title = 'Review weights';
    const subtitle = 'Shipment weights';
    const supertitle = 'Weight 1 of 2';
    const description = 'Shipment 1 of 2';
    const content = 'Some content';
    const buttonText = 'Review billable weights';
    const mockOnClose = jest.fn();
    render(
      <DocumentViewerSidebar
        title={title}
        subtitle={subtitle}
        description={description}
        onClose={mockOnClose}
        supertitle={supertitle}
      >
        <DocumentViewerSidebar.Content>{content}</DocumentViewerSidebar.Content>
        <DocumentViewerSidebar.Footer>
          <Button>{buttonText}</Button>
        </DocumentViewerSidebar.Footer>
      </DocumentViewerSidebar>,
    );

    fireEvent.click(screen.getByTestId('closeSidebar'));
    expect(mockOnClose).toHaveBeenCalledTimes(1);
    expect(screen.getByText(title)).toBeInTheDocument();
    expect(screen.getByText(subtitle)).toBeInTheDocument();
    expect(screen.getByText(supertitle)).toBeInTheDocument();
    expect(screen.getByText(description)).toBeInTheDocument();
    expect(screen.getByText(content)).toBeInTheDocument();
    expect(screen.getByText(buttonText)).toBeInTheDocument();
  });
});

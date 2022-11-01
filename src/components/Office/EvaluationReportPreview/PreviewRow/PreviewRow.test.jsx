import { render, screen } from '@testing-library/react';
import React from 'react';

import PreviewRow from './PreviewRow';

describe('Preview Row', () => {
  it('renders a basic preview row', async () => {
    render(<PreviewRow label="Label" data="Data" />);

    expect(screen.getByText('Label')).toBeInTheDocument();
    expect(screen.getByText('Data')).toBeInTheDocument();
  });

  it('does not render if isShown is false', async () => {
    render(<PreviewRow label="Label" data="Data" isShown={false} />);

    expect(screen.queryByText('Label')).not.toBeInTheDocument();
    expect(screen.queryByText('Data')).not.toBeInTheDocument();
  });
});

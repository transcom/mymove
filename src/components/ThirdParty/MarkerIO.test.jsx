import React from 'react';
import { render, screen } from '@testing-library/react';

import MarkerIO, { setMarkerIOMilmoveTraceID } from './MarkerIO';

describe('components/ThirdParty/MarkerIO', () => {
  it('renders a script tag within the page', async () => {
    render(<MarkerIO />);

    expect(await screen.findByTestId('markerio script tag')).toBeInTheDocument();
  });
  it('adds trace ID to the customData', () => {
    render(<MarkerIO />);
    setMarkerIOMilmoveTraceID('test');
    expect(window.markerConfig.customData.milmoveTraceId).toEqual('test');
  });
});

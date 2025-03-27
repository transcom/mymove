import { render, screen } from '@testing-library/react';

import * as formLabels from './formLabels';

describe('form labels', () => {
  it(`redAsterisk renders a red asterisk`, () => {
    render(formLabels.redAsterisk);

    const asterisk = screen.getByText('*');
    expect(asterisk).toBeInTheDocument();
    expect(asterisk).toHaveStyle({ color: 'red' });
  });

  it(`requiredAsteriskMessage renders correct message with red asterisk`, () => {
    render(formLabels.requiredAsteriskMessage);

    const message = screen.getByTestId('reqAsteriskMsg');
    expect(message).toBeInTheDocument();
    expect(message).toHaveTextContent('Fields marked with * are required.');

    const asterisk = screen.getByText('*');
    expect(asterisk).toHaveStyle({ color: 'red' });
  });

  it(`getLabelWithAsterisk renders label with red asterisk`, () => {
    const label = 'Orders type';

    render(formLabels.getLabelWithAsterisk(label));

    const element = screen.getByTestId('labelWithAsterisk');

    expect(element).toBeInTheDocument();
    expect(element).toHaveTextContent('Orders type *');
  });
});

import React from 'react';
import { render, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import CurrentDutyLocationForm from './CurrentDutyLocationForm';

import dutyLocationFactory from 'utils/test/factories/dutyLocation';

const mockDutyLocation = dutyLocationFactory();
const mockDutyLocation2 = dutyLocationFactory();

describe('CurrentDutyLocationForm component', () => {
  it('renders the form input', async () => {
    const { getByLabelText } = render(
      <CurrentDutyLocationForm
        onSubmit={jest.fn()}
        onBack={jest.fn()}
        initialValues={{ current_location: {} }}
        newDutyLocation={{}}
      />,
    );
    await waitFor(() => {
      expect(getByLabelText('What is your current duty location?')).toBeInstanceOf(HTMLInputElement);
    });
  });

  it('keeps the next button disabled if the form is not filled out', async () => {
    const { getByRole } = render(
      <CurrentDutyLocationForm
        onSubmit={jest.fn()}
        onBack={jest.fn()}
        initialValues={{ current_location: null }}
        newDutyLocation={{}}
      />,
    );
    const submitBtn = getByRole('button', { name: 'Next' });

    await waitFor(() => {
      expect(submitBtn).toHaveAttribute('disabled');
    });
  });

  it('does not disable submit if current and new duty locations are the same', async () => {
    const onSubmit = jest.fn();
    const { getByRole } = render(
      <CurrentDutyLocationForm
        onSubmit={onSubmit}
        onBack={jest.fn()}
        initialValues={{
          current_location: mockDutyLocation,
        }}
        newDutyLocation={mockDutyLocation}
      />,
    );
    const submitBtn = getByRole('button', { name: 'Next' });

    await waitFor(() => {
      expect(submitBtn).not.toHaveAttribute('disabled');
    });
  });

  it('submits the form when its valid', async () => {
    const onSubmit = jest.fn();
    const { getByRole } = render(
      <CurrentDutyLocationForm
        onSubmit={onSubmit}
        onBack={jest.fn()}
        initialValues={{
          current_location: mockDutyLocation,
        }}
        newDutyLocation={mockDutyLocation2}
      />,
    );
    const submitBtn = getByRole('button', { name: 'Next' });

    await userEvent.click(submitBtn);

    await waitFor(() => {
      expect(onSubmit).toHaveBeenCalled();
    });
  });

  it('uses the onBack handler when the back button is clicked', async () => {
    const onBack = jest.fn();
    const { getByRole } = render(
      <CurrentDutyLocationForm
        onSubmit={jest.fn()}
        onBack={onBack}
        initialValues={{
          current_location: {},
        }}
      />,
    );
    const backBtn = getByRole('button', { name: 'Back' });

    await userEvent.click(backBtn);

    await waitFor(() => {
      expect(onBack).toHaveBeenCalled();
    });
  });

  afterEach(jest.resetAllMocks);
});

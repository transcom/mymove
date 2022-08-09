import React from 'react';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import selectEvent from 'react-select-event';

import ServiceInfoForm from './ServiceInfoForm';

jest.mock('components/DutyLocationSearchBox/api', () => ({
  ShowAddress: jest.fn().mockImplementation(() =>
    Promise.resolve({
      city: 'Test City',
      id: '25be4d12-fe93-47f1-bbec-1db386dfa67f',
      postalCode: '12345',
      state: 'NY',
      streetAddress1: '123 Main St',
    }),
  ),
  SearchDutyLocations: jest.fn().mockImplementation(() =>
    Promise.resolve([
      {
        address: {
          city: '',
          id: '00000000-0000-0000-0000-000000000000',
          postalCode: '',
          state: '',
          streetAddress1: '',
        },
        address_id: '46c4640b-c35e-4293-a2f1-36c7b629f903',
        affiliation: 'AIR_FORCE',
        created_at: '2021-02-11T16:48:04.117Z',
        id: '93f0755f-6f35-478b-9a75-35a69211da1c',
        name: 'Altus AFB',
        updated_at: '2021-02-11T16:48:04.117Z',
      },
      {
        address: {
          city: '',
          id: '00000000-0000-0000-0000-000000000000',
          postalCode: '',
          state: '',
          streetAddress1: '',
        },
        address_id: '2d7e17f6-1b8a-4727-8949-007c80961a62',
        affiliation: 'AIR_FORCE',
        created_at: '2021-02-11T16:48:04.117Z',
        id: '7d123884-7c1b-4611-92ae-e8d43ca03ad9',
        name: 'Hill AFB',
        updated_at: '2021-02-11T16:48:04.117Z',
      },
      {
        address: {
          city: 'Test City',
          id: '00000000-0000-0000-0000-000000000010',
          postalCode: '12345',
          state: 'NY',
          streetAddress1: '123 Main St',
        },
        address_id: '25be4d12-fe93-47f1-bbec-1db386dfa67f',
        affiliation: 'AIR_FORCE',
        created_at: '2021-02-11T16:48:04.117Z',
        id: 'a8d6b33c-8370-4e92-8df2-356b8c9d0c1a',
        name: 'Luke AFB',
        updated_at: '2021-02-11T16:48:04.117Z',
      },
      {
        address: {
          city: '',
          id: '00000000-0000-0000-0000-000000000000',
          postalCode: '',
          state: '',
          streetAddress1: '',
        },
        address_id: '3dbf1fc7-3289-4c6e-90aa-01b530a7c3c3',
        affiliation: 'AIR_FORCE',
        created_at: '2021-02-11T16:48:20.225Z',
        id: 'd01bd2a4-6695-4d69-8f2f-69e88dff58f8',
        name: 'Shaw AFB',
        updated_at: '2021-02-11T16:48:20.225Z',
      },
      {
        address: {
          city: '',
          id: '00000000-0000-0000-0000-000000000000',
          postalCode: '',
          state: '',
          streetAddress1: '',
        },
        address_id: '1af8f0f3-f75f-46d3-8dc8-c67c2feeb9f0',
        affiliation: 'AIR_FORCE',
        created_at: '2021-02-11T16:49:14.322Z',
        id: 'b1f9a535-96d4-4cc3-adf1-b76505ce0765',
        name: 'Yuma AFB',
        updated_at: '2021-02-11T16:49:14.322Z',
      },
      {
        address: {
          city: '',
          id: '00000000-0000-0000-0000-000000000000',
          postalCode: '',
          state: '',
          streetAddress1: '',
        },
        address_id: 'f2adfebc-7703-4d06-9b49-c6ca8f7968f1',
        affiliation: 'AIR_FORCE',
        created_at: '2021-02-11T16:48:20.225Z',
        id: 'a268b48f-0ad1-4a58-b9d6-6de10fd63d96',
        name: 'Los Angeles AFB',
        updated_at: '2021-02-11T16:48:20.225Z',
      },
      {
        address: {
          city: '',
          id: '00000000-0000-0000-0000-000000000000',
          postalCode: '',
          state: '',
          streetAddress1: '',
        },
        address_id: '13eb2cab-cd68-4f43-9532-7a71996d3296',
        affiliation: 'AIR_FORCE',
        created_at: '2021-02-11T16:48:20.225Z',
        id: 'a48fda70-8124-4e90-be0d-bf8119a98717',
        name: 'Wright-Patterson AFB',
        updated_at: '2021-02-11T16:48:20.225Z',
      },
    ]),
  ),
}));

describe('ServiceInfoForm', () => {
  const testProps = {
    onSubmit: jest.fn(),
    onCancel: jest.fn(),
    initialValues: {
      first_name: '',
      middle_name: '',
      last_name: '',
      suffix: '',
      affiliation: '',
      edipi: '',
      rank: '',
      current_location: {},
    },
    newDutyLocation: {},
  };

  it('renders the form inputs', async () => {
    render(<ServiceInfoForm {...testProps} />);

    const firstNameInput = await screen.findByLabelText('First name');
    expect(firstNameInput).toBeInstanceOf(HTMLInputElement);
    expect(firstNameInput).toBeRequired();

    expect(await screen.findByLabelText(/Middle name/)).toBeInstanceOf(HTMLInputElement);

    const lastNameInput = await screen.findByLabelText('Last name');
    expect(lastNameInput).toBeInstanceOf(HTMLInputElement);
    expect(lastNameInput).toBeRequired();

    expect(await screen.findByLabelText(/Suffix/)).toBeInstanceOf(HTMLInputElement);

    const branchInput = await screen.findByLabelText('Branch of service');
    expect(branchInput).toBeInstanceOf(HTMLSelectElement);
    expect(branchInput).toBeRequired();

    const dodInput = await screen.findByLabelText('DoD ID number');
    expect(dodInput).toBeInstanceOf(HTMLInputElement);
    expect(dodInput).toBeRequired();

    const rankInput = await screen.findByLabelText('Rank');
    expect(rankInput).toBeInstanceOf(HTMLSelectElement);
    expect(rankInput).toBeRequired();

    expect(await screen.findByLabelText('Current duty location')).toBeInstanceOf(HTMLInputElement);
  });

  it('validates the DOD ID number on blur', async () => {
    render(<ServiceInfoForm {...testProps} />);

    const dodInput = await screen.findByLabelText('DoD ID number');
    await userEvent.type(dodInput, 'not a valid ID number');
    await userEvent.tab();

    expect(dodInput).not.toBeValid();
    expect(await screen.findByText('Enter a 10-digit DOD ID number')).toBeInTheDocument();
  });

  it('validates the new duty location against the current duty location', async () => {
    render(
      <ServiceInfoForm
        {...testProps}
        newDutyLocation={{ name: 'Luke AFB', id: 'a8d6b33c-8370-4e92-8df2-356b8c9d0c1a' }}
      />,
    );

    // Test Duty Location Search Box interaction
    const dutyLocationInput = await screen.getByLabelText('Current duty location');
    fireEvent.change(dutyLocationInput, { target: { value: 'AFB' } });
    await selectEvent.select(dutyLocationInput, /Luke/);

    expect(await screen.findByRole('form')).toHaveFormValues({
      current_location: 'Luke AFB',
    });

    expect(await screen.findByRole('button', { name: 'Save' })).toHaveAttribute('disabled');
    expect(
      await screen.findByText(
        'You entered the same duty location for your origin and destination. Please change one of them.',
      ),
    ).toBeInTheDocument();
  });

  it('shows an error message if trying to submit an invalid form', async () => {
    render(<ServiceInfoForm {...testProps} />);

    // Touch required fields to show validation errors
    await userEvent.click(screen.getByLabelText('First name'));
    await userEvent.click(screen.getByLabelText('Last name'));
    await userEvent.click(screen.getByLabelText('Branch of service'));
    await userEvent.click(screen.getByLabelText('DoD ID number'));
    await userEvent.click(screen.getByLabelText('Rank'));

    const submitBtn = screen.getByRole('button', { name: 'Save' });
    await userEvent.click(submitBtn);

    await waitFor(() => {
      expect(screen.getAllByText('Required').length).toBe(5);
    });
    expect(testProps.onSubmit).not.toHaveBeenCalled();
  });

  it('submits the form when its valid', async () => {
    render(<ServiceInfoForm {...testProps} />);
    const submitBtn = screen.getByRole('button', { name: 'Save' });

    await userEvent.type(screen.getByLabelText('First name'), 'Leo');
    await userEvent.type(screen.getByLabelText('Last name'), 'Spaceman');
    await userEvent.selectOptions(screen.getByLabelText('Branch of service'), ['NAVY']);
    await userEvent.type(screen.getByLabelText('DoD ID number'), '1234567890');
    await userEvent.selectOptions(screen.getByLabelText('Rank'), ['E_5']);
    fireEvent.change(screen.getByLabelText('Current duty location'), { target: { value: 'AFB' } });
    await selectEvent.select(screen.getByLabelText('Current duty location'), /Luke/);

    expect(screen.getByRole('form')).toHaveFormValues({
      current_location: 'Luke AFB',
    });

    await userEvent.click(submitBtn);

    await waitFor(() => {
      expect(testProps.onSubmit).toHaveBeenCalledWith(
        expect.objectContaining({
          first_name: 'Leo',
          last_name: 'Spaceman',
          affiliation: 'NAVY',
          edipi: '1234567890',
          rank: 'E_5',
          current_location: {
            address: {
              city: 'Test City',
              id: '25be4d12-fe93-47f1-bbec-1db386dfa67f',
              postalCode: '12345',
              state: 'NY',
              streetAddress1: '123 Main St',
            },
            address_id: '25be4d12-fe93-47f1-bbec-1db386dfa67f',
            affiliation: 'AIR_FORCE',
            created_at: '2021-02-11T16:48:04.117Z',
            id: 'a8d6b33c-8370-4e92-8df2-356b8c9d0c1a',
            name: 'Luke AFB',
            updated_at: '2021-02-11T16:48:04.117Z',
          },
        }),
        expect.anything(),
      );
    });
  });

  it('uses the onCancel handler when the cancel button is clicked', async () => {
    const onCancel = jest.fn();
    render(<ServiceInfoForm {...testProps} onCancel={onCancel} />);
    const cancelBtn = screen.getByRole('button', { name: 'Cancel' });

    await userEvent.click(cancelBtn);

    await waitFor(() => {
      expect(onCancel).toHaveBeenCalled();
    });
  });

  afterEach(jest.restoreAllMocks);
});

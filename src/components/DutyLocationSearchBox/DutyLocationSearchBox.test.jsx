import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import DutyLocationSearchBox from './DutyLocationSearchBox';

const testAddress = {
  city: 'Glendale Luke AFB',
  country: 'United States',
  id: 'fa51dab0-4553-4732-b843-1f33407f77bc',
  postalCode: '85309',
  state: 'AZ',
  streetAddress1: 'n/a',
};

const testDutyLocations = [
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
      city: '',
      id: '00000000-0000-0000-0000-000000000000',
      postalCode: '',
      state: '',
      streetAddress1: '',
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
];

jest.mock('./api.js', () => ({
  SearchDutyLocations: async (search) => {
    if (search === 'empty') {
      return [];
    }

    if (search === 'broken') {
      throw new Error('Server returned an error');
    }

    return testDutyLocations;
  },
  ShowAddress: async () => {
    return testAddress;
  },
}));

describe('DutyLocationSearchBoxContainer', () => {
  describe('basic rendering', () => {
    it('renders with minimal props', () => {
      render(<DutyLocationSearchBox input={{ name: 'test_component' }} name="test_component" />);
      expect(screen.getByLabelText('Name of Duty Location:')).toBeInTheDocument();
    });

    it('renders the title', () => {
      render(<DutyLocationSearchBox input={{ name: 'test_component' }} name="test_component" title="Test Component" />);
      expect(screen.getByLabelText('Test Component')).toBeInTheDocument();
    });

    it('renders the default placeholder text', () => {
      render(<DutyLocationSearchBox input={{ name: 'test_component' }} name="test_component" />);
      expect(screen.getByText('Start typing a duty location...')).toBeInTheDocument();
    });

    it('renders an error message', () => {
      render(
        <DutyLocationSearchBox
          input={{ name: 'test_component' }}
          name="test_component"
          errorMsg="Test Error Message"
        />,
      );
      expect(screen.getByText('Test Error Message')).toBeInTheDocument();
    });

    it('renders a value passed in via prop', () => {
      render(
        <DutyLocationSearchBox
          name="test_component"
          input={{
            name: 'test_component',
            value: {
              ...testDutyLocations[2],
              address: testAddress,
            },
          }}
        />,
      );
      expect(screen.getByText('Luke AFB')).toBeInTheDocument();
      expect(screen.getByText('Glendale Luke AFB, AZ 85309')).toBeInTheDocument();
    });

    it('can render without the address', () => {
      render(
        <DutyLocationSearchBox
          name="test_component"
          displayAddress={false}
          input={{
            name: 'test_component',
            value: {
              ...testDutyLocations[2],
              address: testAddress,
            },
          }}
        />,
      );
      expect(screen.getByText('Luke AFB')).toBeInTheDocument();
      expect(screen.queryByText('Glendale Luke AFB, AZ 85309')).not.toBeInTheDocument();
    });

    it('can show placeholder text based on prop', () => {
      const testPlaceholderText = 'Test Placeholder Text';
      render(
        <DutyLocationSearchBox
          input={{ name: 'test_component' }}
          name="test_component"
          placeholder={testPlaceholderText}
        />,
      );
      expect(screen.getByText(testPlaceholderText)).toBeInTheDocument();
    });
  });

  describe('updating options based on text', () => {
    it('searches user input and renders options', async () => {
      render(<DutyLocationSearchBox input={{ name: 'test_component' }} title="Test Component" name="test_component" />);
      await userEvent.type(screen.getByLabelText('Test Component'), 'AFB');

      const option = await screen.findByText('Luke');
      expect(option).toBeInTheDocument();

      const optionsContainer = option.closest('div').parentElement;
      expect(optionsContainer.children.length).toEqual(7);
    });

    it('searches user input and renders a message if empty', async () => {
      render(<DutyLocationSearchBox input={{ name: 'test_component' }} title="Test Component" name="test_component" />);
      await userEvent.type(screen.getByLabelText('Test Component'), 'empty');

      expect(await screen.findByText('No Options')).toBeInTheDocument();
    });

    it("doesnt search if user input isn't 2+ characters in length", async () => {
      render(<DutyLocationSearchBox input={{ name: 'test_component' }} title="Test Component" name="test_component" />);
      await userEvent.type(screen.getByLabelText('Test Component'), '1');

      expect(await screen.findByText('No Options')).toBeInTheDocument();
    });

    it('handles server errors', async () => {
      render(<DutyLocationSearchBox input={{ name: 'test_component' }} title="Test Component" name="test_component" />);
      await userEvent.type(screen.getByLabelText('Test Component'), 'broken');

      expect(await screen.findByText('No Options')).toBeInTheDocument();
    });
  });

  describe('selecting options', () => {
    it('selects an option, calls the onChange callback prop', async () => {
      const onChange = jest.fn();
      render(
        <DutyLocationSearchBox
          input={{ name: 'test_component', onChange }}
          title="Test Component"
          name="test_component"
        />,
      );
      await userEvent.type(screen.getByLabelText('Test Component'), 'AFB');
      await userEvent.click(await screen.findByText('Luke'));

      await waitFor(() =>
        expect(onChange).toHaveBeenCalledWith({
          ...testDutyLocations[2],
          address: testAddress,
        }),
      );
    });
  });
});

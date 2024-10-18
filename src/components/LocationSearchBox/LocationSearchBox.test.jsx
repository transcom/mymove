import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import LocationSearchBox from './LocationSearchBox';

import dutyLocationFactory from 'utils/test/factories/dutyLocation';

const mockLocations = [
  dutyLocationFactory(),
  dutyLocationFactory(),
  dutyLocationFactory(),
  dutyLocationFactory(),
  dutyLocationFactory(),
  dutyLocationFactory(),
  dutyLocationFactory(),
];

const selectedLocations = mockLocations[2];
const mockAddress = selectedLocations.address;
const displayAddress = `${mockAddress.city}, ${mockAddress.state} ${mockAddress.postalCode}`;
const optionName = selectedLocations.name.split(' AFB')[0];

jest.mock('./api.js', () => ({
  ShowAddress: async () => {
    return mockAddress;
  },
}));

const mockLocationSearch = jest.fn(async (search) => {
  if (search === 'empty') {
    return [];
  }

  if (search === 'broken') {
    throw new Error('Server returned an error');
  }

  return mockLocations;
});

describe('LocationSearchBoxContainer', () => {
  describe('basic rendering', () => {
    it('renders with minimal props', () => {
      render(
        <LocationSearchBox
          input={{ name: 'test_component' }}
          name="test_component"
          searchLocations={mockLocationSearch}
        />,
      );
      expect(screen.getByLabelText('Name of Duty Location:')).toBeInTheDocument();
    });

    it('renders the title', () => {
      render(
        <LocationSearchBox
          input={{ name: 'test_component' }}
          name="test_component"
          title="Test Component"
          searchLocations={mockLocationSearch}
        />,
      );
      expect(screen.getByLabelText('Test Component')).toBeInTheDocument();
    });

    it('renders the default placeholder text', () => {
      render(
        <LocationSearchBox
          input={{ name: 'test_component' }}
          name="test_component"
          searchLocations={mockLocationSearch}
        />,
      );
      expect(screen.getByText('Start typing a duty location...')).toBeInTheDocument();
    });

    it('renders an error message', () => {
      render(
        <LocationSearchBox
          input={{ name: 'test_component' }}
          name="test_component"
          errorMsg="Test Error Message"
          searchLocations={mockLocationSearch}
        />,
      );
      expect(screen.getByText('Test Error Message')).toBeInTheDocument();
    });

    it('renders a value passed in via prop', () => {
      render(
        <LocationSearchBox
          name="test_component"
          input={{
            name: 'test_component',
            value: {
              ...selectedLocations,
              address: mockAddress,
            },
          }}
        />,
      );
      expect(screen.getByText(selectedLocations.name)).toBeInTheDocument();
      expect(screen.getByText(displayAddress)).toBeInTheDocument();
    });

    it('can render without the address', () => {
      render(
        <LocationSearchBox
          name="test_component"
          displayAddress={false}
          searchLocations={mockLocationSearch}
          input={{
            name: 'test_component',
            value: {
              ...selectedLocations,
              address: mockAddress,
            },
          }}
        />,
      );
      expect(screen.getByText(selectedLocations.name)).toBeInTheDocument();
      expect(screen.queryByText(displayAddress)).not.toBeInTheDocument();
    });

    it('can show placeholder text based on prop', () => {
      const testPlaceholderText = 'Test Placeholder Text';
      render(
        <LocationSearchBox
          input={{ name: 'test_component' }}
          name="test_component"
          placeholder={testPlaceholderText}
          searchLocations={mockLocationSearch}
        />,
      );
      expect(screen.getByText(testPlaceholderText)).toBeInTheDocument();
    });
  });

  describe('updating options based on text', () => {
    it('searches user input and renders options', async () => {
      render(
        <LocationSearchBox
          input={{ name: 'test_component' }}
          title="Test Component"
          name="test_component"
          searchLocations={mockLocationSearch}
        />,
      );
      await userEvent.type(screen.getByLabelText('Test Component'), 'AFB');

      const option = await screen.findByText(optionName);
      expect(option).toBeInTheDocument();

      const optionsContainer = option.closest('div').parentElement;
      expect(optionsContainer.children.length).toEqual(7);
    });

    it('searches user input and renders a message if empty', async () => {
      render(
        <LocationSearchBox
          input={{ name: 'test_component' }}
          title="Test Component"
          name="test_component"
          searchLocations={mockLocationSearch}
        />,
      );
      await userEvent.type(screen.getByLabelText('Test Component'), 'empty');

      expect(await screen.findByText('No Options')).toBeInTheDocument();
    });

    it("doesnt search if user input isn't 2+ characters in length", async () => {
      render(
        <LocationSearchBox
          input={{ name: 'test_component' }}
          title="Test Component"
          name="test_component"
          searchLocations={mockLocationSearch}
        />,
      );
      await userEvent.type(screen.getByLabelText('Test Component'), '1');

      expect(await screen.findByText('No Options')).toBeInTheDocument();
    });

    it('handles server errors', async () => {
      render(
        <LocationSearchBox
          input={{ name: 'test_component' }}
          title="Test Component"
          name="test_component"
          searchLocations={mockLocationSearch}
        />,
      );
      await userEvent.type(screen.getByLabelText('Test Component'), 'broken');

      expect(await screen.findByText('No Options')).toBeInTheDocument();
    });
  });

  describe('selecting options', () => {
    it('selects an option, calls the onChange callback prop', async () => {
      const onChange = jest.fn();
      render(
        <LocationSearchBox
          input={{ name: 'test_component', onChange }}
          title="Test Component"
          name="test_component"
          searchLocations={mockLocationSearch}
        />,
      );
      await userEvent.type(screen.getByLabelText('Test Component'), 'AFB');
      await userEvent.click(await screen.findByText(optionName));

      await waitFor(() =>
        expect(onChange).toHaveBeenCalledWith({
          ...selectedLocations,
          address: mockAddress,
        }),
      );
    });
  });
});

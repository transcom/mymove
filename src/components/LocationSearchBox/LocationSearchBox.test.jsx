import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import LocationSearchBox from './LocationSearchBox';

import dutyLocationFactory from 'utils/test/factories/dutyLocation';
// import transportationOfficeFactory from 'utils/test/factories/transportationOffice';

// const mockCloseoutLocations = [
//   transportationOfficeFactory(),
//   transportationOfficeFactory(),
//   transportationOfficeFactory(),
//   transportationOfficeFactory(),
//   transportationOfficeFactory(),
//   transportationOfficeFactory(),
//   transportationOfficeFactory(),
// ];

// const selectedCloseoutLocation = mockCloseoutLocations[2];
// const mockCloseoutAddress = selectedCloseoutLocation.address;
// const displayCloseoutAddress = `${mockCloseoutAddress.city}, ${mockCloseoutAddress.state} ${mockCloseoutAddress.postalCode}`;
// const closeoutOptionName = selectedCloseoutLocation.name.split(' AFB'[0]);

const mockDutyLocations = [
  dutyLocationFactory(),
  dutyLocationFactory(),
  dutyLocationFactory(),
  dutyLocationFactory(),
  dutyLocationFactory(),
  dutyLocationFactory(),
  dutyLocationFactory(),
];

const selectedDutyLocation = mockDutyLocations[2];
const mockAddress = selectedDutyLocation.address;
const displayAddress = `${mockAddress.city}, ${mockAddress.state} ${mockAddress.postalCode}`;
const optionName = selectedDutyLocation.name.split(' AFB')[0];

jest.mock('./api.js', () => ({
  SearchDutyLocations: async (search) => {
    if (search === 'empty') {
      return [];
    }

    if (search === 'broken') {
      throw new Error('Server returned an error');
    }

    return mockDutyLocations;
  },
  ShowAddress: async () => {
    return mockAddress;
  },
}));

describe('LocationSearchBoxContainer', () => {
  describe('basic rendering', () => {
    it('renders with minimal props', () => {
      render(<LocationSearchBox input={{ name: 'test_component' }} name="test_component" />);
      expect(screen.getByLabelText('Name of Duty Location:')).toBeInTheDocument();
    });

    it('renders the title', () => {
      render(<LocationSearchBox input={{ name: 'test_component' }} name="test_component" title="Test Component" />);
      expect(screen.getByLabelText('Test Component')).toBeInTheDocument();
    });

    it('renders the default placeholder text', () => {
      render(<LocationSearchBox input={{ name: 'test_component' }} name="test_component" />);
      expect(screen.getByText('Start typing a duty location...')).toBeInTheDocument();
    });

    it('renders an error message', () => {
      render(
        <LocationSearchBox input={{ name: 'test_component' }} name="test_component" errorMsg="Test Error Message" />,
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
              ...selectedDutyLocation,
              address: mockAddress,
            },
          }}
        />,
      );
      expect(screen.getByText(selectedDutyLocation.name)).toBeInTheDocument();
      expect(screen.getByText(displayAddress)).toBeInTheDocument();
    });

    it('can render without the address', () => {
      render(
        <LocationSearchBox
          name="test_component"
          displayAddress={false}
          input={{
            name: 'test_component',
            value: {
              ...selectedDutyLocation,
              address: mockAddress,
            },
          }}
        />,
      );
      expect(screen.getByText(selectedDutyLocation.name)).toBeInTheDocument();
      expect(screen.queryByText(displayAddress)).not.toBeInTheDocument();
    });

    it('can show placeholder text based on prop', () => {
      const testPlaceholderText = 'Test Placeholder Text';
      render(
        <LocationSearchBox
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
      render(<LocationSearchBox input={{ name: 'test_component' }} title="Test Component" name="test_component" />);
      await userEvent.type(screen.getByLabelText('Test Component'), 'AFB');

      const option = await screen.findByText(optionName);
      expect(option).toBeInTheDocument();

      const optionsContainer = option.closest('div').parentElement;
      expect(optionsContainer.children.length).toEqual(7);
    });

    it('searches user input and renders a message if empty', async () => {
      render(<LocationSearchBox input={{ name: 'test_component' }} title="Test Component" name="test_component" />);
      await userEvent.type(screen.getByLabelText('Test Component'), 'empty');

      expect(await screen.findByText('No Options')).toBeInTheDocument();
    });

    it("doesnt search if user input isn't 2+ characters in length", async () => {
      render(<LocationSearchBox input={{ name: 'test_component' }} title="Test Component" name="test_component" />);
      await userEvent.type(screen.getByLabelText('Test Component'), '1');

      expect(await screen.findByText('No Options')).toBeInTheDocument();
    });

    it('handles server errors', async () => {
      render(<LocationSearchBox input={{ name: 'test_component' }} title="Test Component" name="test_component" />);
      await userEvent.type(screen.getByLabelText('Test Component'), 'broken');

      expect(await screen.findByText('No Options')).toBeInTheDocument();
    });
  });

  describe('selecting options', () => {
    it('selects an option, calls the onChange callback prop', async () => {
      const onChange = jest.fn();
      render(
        <LocationSearchBox input={{ name: 'test_component', onChange }} title="Test Component" name="test_component" />,
      );
      await userEvent.type(screen.getByLabelText('Test Component'), 'AFB');
      await userEvent.click(await screen.findByText(optionName));

      await waitFor(() =>
        expect(onChange).toHaveBeenCalledWith({
          ...selectedDutyLocation,
          address: mockAddress,
        }),
      );
    });
  });
});

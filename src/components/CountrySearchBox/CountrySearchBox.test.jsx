import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import CountrySearchBox from './CountrySearchBox';

const mockCountries = [
  {
    code: 'US',
    name: 'UNITED STATES',
    id: '791899e6-cd77-46f2-981b-176ecb8d7098',
  },
];

const selectedCountries = mockCountries[0];

const mockCountrySearch = jest.fn(async (search) => {
  if (search === 'empty') {
    return [];
  }

  if (search === 'broken') {
    throw new Error('Server returned an error');
  }

  return mockCountries;
});

describe('CountrySearchBoxContainer', () => {
  describe('basic rendering', () => {
    it('renders with minimal props', () => {
      render(
        <CountrySearchBox
          input={{ name: 'test_component' }}
          name="test_component"
          searchCountries={mockCountrySearch}
        />,
      );
      expect(screen.getByText('Start typing a country name, code')).toBeInTheDocument();
    });

    it('renders the title', () => {
      render(
        <CountrySearchBox
          input={{ name: 'test_component' }}
          name="test_component"
          title="Test Component"
          searchCountries={mockCountrySearch}
        />,
      );
      expect(screen.getByLabelText('Test Component')).toBeInTheDocument();
    });

    it('renders the default placeholder text', () => {
      render(
        <CountrySearchBox
          input={{ name: 'test_component' }}
          name="test_component"
          searchCountries={mockCountrySearch}
        />,
      );
      expect(screen.getByText('Start typing a country name, code')).toBeInTheDocument();
    });

    it('renders the required asterisk when prop is provided', () => {
      render(
        <CountrySearchBox
          input={{ name: 'test_component' }}
          name="test_component"
          searchCountries={mockCountrySearch}
          showRequiredAsterisk
        />,
      );
      expect(screen.getByTestId('requiredAsterisk')).toBeInTheDocument();
    });

    it('renders an error message', () => {
      const onChange = jest.fn();
      const locationState = jest.fn();
      render(
        <CountrySearchBox
          input={{ name: 'test_component', onChange, locationState }}
          name="test_component"
          errorMsg="Test Error Message"
          searchCountries={mockCountrySearch}
        />,
      );
      expect(screen.getByText('Test Error Message')).toBeInTheDocument();
    });

    it('renders a value passed in via prop', async () => {
      render(
        <CountrySearchBox
          name="test_component"
          input={{
            name: 'test_component',
            value: {
              ...mockCountries,
              country: selectedCountries,
            },
          }}
          searchCountries={mockCountrySearch}
        />,
      );
      await waitFor(() => {
        const elements = screen.queryAllByText(/(US)/);
        expect(elements).toHaveLength(1);
      });
    });

    it('can show placeholder text based on prop', () => {
      const testPlaceholderText = 'Test Placeholder Text';
      render(
        <CountrySearchBox
          input={{ name: 'test_component' }}
          name="test_component"
          placeholder={testPlaceholderText}
          searchCountries={mockCountrySearch}
        />,
      );
      expect(screen.getByText(testPlaceholderText)).toBeInTheDocument();
    });

    it('renders a required asterisk', () => {
      const onChange = jest.fn();
      const locationState = jest.fn();
      render(
        <CountrySearchBox
          input={{ name: 'test_component', onChange, locationState }}
          name="test_component"
          searchCountries={mockCountrySearch}
          showRequiredAsterisk
        />,
      );
      expect(screen.getByTestId('requiredAsterisk')).toBeInTheDocument();
    });
  });

  describe('updating options based on text', () => {
    it('searches user input and renders options', async () => {
      const onChange = jest.fn();
      const locationState = jest.fn();
      render(
        <CountrySearchBox
          input={{ name: 'test_component', onChange, locationState }}
          title="Test Component"
          name="test_component"
          searchCountries={mockCountrySearch}
        />,
      );
      await userEvent.type(screen.getByLabelText('Test Component'), 'UNITED STATES');

      const option = await screen.findByText('(US)');
      expect(option).toBeInTheDocument();
    });

    it('searches user input and renders a message if empty', async () => {
      const onChange = jest.fn();
      const locationState = jest.fn();
      render(
        <CountrySearchBox
          input={{ name: 'test_component', onChange, locationState }}
          title="Test Component"
          name="test_component"
          searchCountries={mockCountrySearch}
        />,
      );
      await userEvent.type(screen.getByLabelText('Test Component'), 'empty');

      expect(await screen.findByText('No Options')).toBeInTheDocument();
    });

    it("doesnt search if user input isn't 2+ characters in length", async () => {
      const onChange = jest.fn();
      const locationState = jest.fn();
      render(
        <CountrySearchBox
          input={{ name: 'test_component', onChange, locationState }}
          title="Test Component"
          name="test_component"
          searchCountries={mockCountrySearch}
        />,
      );
      await userEvent.type(screen.getByLabelText('Test Component'), '1');

      expect(await screen.findByText('No Options')).toBeInTheDocument();
    });

    it('handles server errors', async () => {
      const onChange = jest.fn();
      const locationState = jest.fn();
      render(
        <CountrySearchBox
          input={{ name: 'test_component', onChange, locationState }}
          title="Test Component"
          name="test_component"
          searchCountries={mockCountrySearch}
        />,
      );
      await userEvent.type(screen.getByLabelText('Test Component'), 'broken');

      expect(await screen.findByText('No Options')).toBeInTheDocument();
    });
  });

  describe('selecting options', () => {
    it('selects an option, calls the onChange callback prop', async () => {
      const onChange = jest.fn();
      const locationState = jest.fn();
      render(
        <CountrySearchBox
          input={{ name: 'test_component', onChange, locationState }}
          title="Test Component"
          name="test_component"
          searchLocations={mockCountrySearch}
        />,
      );

      await userEvent.type(screen.getByLabelText('Test Component'), 'UNITED STATES');
      await userEvent.click(await screen.findByText('(US)'));

      await waitFor(() =>
        expect(onChange).toHaveBeenCalledWith({
          ...mockCountries,
          country: selectedCountries,
        }),
      );
    });
  });
});

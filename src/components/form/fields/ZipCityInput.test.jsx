import React from 'react';
import { shallow } from 'enzyme';
import AsyncSelect from 'react-select/async';

import { ZipCityInput } from './ZipCityInput';

import { LocationSearchBoxComponent, LocationSearchBoxContainer } from 'components/LocationSearchBox/LocationSearchBox';
import { searchLocationByZipCity, showCounselingOffices } from 'services/internalApi';

const mockOnChange = jest.fn();
const mockSetValue = jest.fn();
// mock out formik hook as we are not testing formik
// needs to be before first describe
jest.mock('formik', () => {
  return {
    ...jest.requireActual('formik'),
    useField: () => [
      {
        onChange: mockOnChange,
      },
      { touched: true, error: 'sample error' },
      { setValue: mockSetValue },
    ],
  };
});

jest.mock('services/internalApi', () => ({
  ...jest.requireActual('services/internalApi'),
  showCounselingOffices: jest.fn().mockImplementation(() =>
    Promise.resolve({
      body: [
        {
          id: '3e937c1f-5539-4919-954d-017989130584',
          name: 'Albuquerque AFB',
        },
        {
          id: 'fa51dab0-4553-4732-b843-1f33407f77bc',
          name: 'Glendale Luke AFB',
        },
      ],
    }),
  ),
  searchLocationByZipCity: jest.fn(),
}));

const handleZipCityChange = jest.fn();

describe('ZipCityInput', () => {
  describe('with all required props', () => {
    const wrapper = shallow(
      <ZipCityInput
        name="zipCity"
        placeholder="Start typing a Zip Code or City..."
        label="Zip/City Lookup"
        displayAddress={false}
        handleZipCityChange={handleZipCityChange}
      />,
    );

    it('renders a Zip City search input', () => {
      const input = wrapper.find(LocationSearchBoxContainer);
      expect(input.length).toBe(1);
    });

    it('triggers onChange properly', async () => {
      const cityName = 'El Paso';
      showCounselingOffices.mockImplementation(() => Promise.resolve({}));
      searchLocationByZipCity.mockImplementation(() => Promise.resolve(cityName));
      const container = wrapper.find(LocationSearchBoxContainer).dive();
      const component = container.find(LocationSearchBoxComponent).dive();
      const select = component.find(AsyncSelect);
      await select.simulate('change', { city: cityName });
      expect(mockSetValue).toHaveBeenCalledWith({ city: cityName });
    });
  });

  afterEach(jest.resetAllMocks);
});

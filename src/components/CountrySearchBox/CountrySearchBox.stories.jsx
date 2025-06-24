import React, { useState } from 'react';

import { CountrySearchBoxComponent } from './CountrySearchBox';

export default {
  title: 'Components/Country Search Box',
  component: CountrySearchBoxComponent,
};

const testCountries = [
  {
    code: 'US',
    name: 'UNITED STATES',
    id: '791899e6-cd77-46f2-981b-176ecb8d7098',
  },
];

const mockSearchCountries = async () => {
  return testCountries;
};

const baseValue = {
  ...testCountries[0],
  country: { ...testCountries[0] },
};

export const CountrySearchBasic = () => {
  const [value, setValue] = useState();

  const onCountryChange = (newValue) => {
    if (newValue !== null) {
      setValue(newValue.name);
    }
  };

  return (
    <CountrySearchBoxComponent
      input={{
        value,
        onChange: onCountryChange,
        countryState: () => {},
        name: 'test_component',
      }}
      title="Search Country Test Component"
      name="test_component"
      searchCountries={mockSearchCountries}
      handleCountryOnChange={onCountryChange}
    />
  );
};

export const CountrySearchWithDefaultValue = () => {
  const [value, setValue] = useState(baseValue);

  const onChange = (newValue) => {
    setValue(newValue);
  };

  return (
    <CountrySearchBoxComponent
      input={{ name: 'test_component', onChange, value }}
      title="Search Country Test Component With Default Value"
      searchCountries={mockSearchCountries}
    />
  );
};

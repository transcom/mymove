import React, { useState } from 'react';

import { LocationSearchBoxComponent } from './LocationSearchBox';

export default {
  title: 'Components/Location Search Box',
  component: LocationSearchBoxComponent,
};

const testZipCity = [
  {
    city: 'Glendale Luke AFB',
    county: 'Maricopa',
    postalCode: '85309',
    state: 'AZ',
  },
  {
    city: 'El Paso',
    county: 'El Paso',
    postalCode: '79912',
    state: 'TX',
  },
];

const testAddress = {
  city: 'Glendale Luke AFB',
  country: 'United States',
  id: 'fa51dab0-4553-4732-b843-1f33407f77bc',
  postalCode: '85309',
  state: 'AZ',
  streetAddress1: 'n/a',
};

const testLocations = [
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

const baseValue = {
  ...testLocations[2],
  address: { ...testAddress },
};

const mockSearchLocations = async () => {
  return testLocations;
};

const mockShowAddress = async () => {
  return testAddress;
};

const mockZipCity = async () => {
  return testZipCity;
};

export const DutyStationBasic = () => {
  const [value, setValue] = useState();

  const onChange = (newValue) => {
    setValue(newValue);
  };

  return (
    <LocationSearchBoxComponent
      input={{ name: 'test_component', onChange, value }}
      title="Duty Station Test Component"
      name="test_component"
      searchLocations={mockSearchLocations}
      showAddress={mockShowAddress}
    />
  );
};

export const DutyStationWithValue = () => {
  const [value, setValue] = useState(baseValue);

  const onChange = (newValue) => {
    setValue(newValue);
  };

  return (
    <LocationSearchBoxComponent
      input={{ name: 'test_component', onChange, value }}
      title="Duty Station Test Component"
      displayAddress={false}
      searchLocations={mockSearchLocations}
      showAddress={mockShowAddress}
    />
  );
};

export const DutyStationWithValueAndAddress = () => {
  const [value, setValue] = useState(baseValue);

  const onChange = (newValue) => {
    setValue(newValue);
  };

  return (
    <LocationSearchBoxComponent
      input={{ name: 'test_component', onChange, value }}
      title="Duty Station Test Component"
      searchLocations={mockSearchLocations}
      showAddress={mockShowAddress}
    />
  );
};

export const DutyStationWithErrorMessage = () => {
  const [value, setValue] = useState();

  const onChange = (newValue) => {
    setValue(newValue);
  };

  return (
    <LocationSearchBoxComponent
      input={{ name: 'test_component', onChange, value }}
      title="Duty Station Test Component"
      errorMsg="Something went wrong"
      searchLocations={mockSearchLocations}
      showAddress={mockShowAddress}
    />
  );
};

export const DutyStationWithLocalError = () => {
  const [value, setValue] = useState();

  const onChange = (newValue) => {
    setValue(newValue);
  };

  const brokenSearchDutyLocations = async () => {
    throw new Error('Artificial error message text');
  };

  return (
    <LocationSearchBoxComponent
      input={{ name: 'test_component', onChange, value }}
      title="Duty Station Test Component"
      searchLocations={brokenSearchDutyLocations}
      showAddress={mockShowAddress}
    />
  );
};

export const TransportationLocationBasic = () => {
  const [value, setValue] = useState();

  const onChange = (newValue) => {
    setValue(newValue);
  };

  return (
    <LocationSearchBoxComponent
      input={{ name: 'test_component', onChange, value }}
      placeholder="Start typing a closeout office..."
      title="Transportation Office Test Component"
      name="test_component"
      searchLocations={mockSearchLocations}
      showAddress={mockShowAddress}
    />
  );
};

export const TransportationLocationWithValue = () => {
  const [value, setValue] = useState(baseValue);

  const onChange = (newValue) => {
    setValue(newValue);
  };

  return (
    <LocationSearchBoxComponent
      input={{ name: 'test_component', onChange, value }}
      placeholder="Start typing a closeout office..."
      title="Transportation Office Test Component"
      displayAddress={false}
      searchLocations={mockSearchLocations}
      showAddress={mockShowAddress}
    />
  );
};

export const ZipCityLocationBasic = () => {
  const [value, setValue] = useState();

  const handleZipCityOnChange = (newValue) => {
    setValue(newValue);
  };

  return (
    <LocationSearchBoxComponent
      input={{
        value,
        onChange: handleZipCityOnChange,
        locationState: () => {},
        name: 'test_component',
      }}
      placeholder="Start typing a Zip Code or City..."
      title="Zip/City Lookup"
      name="test_component"
      searchLocations={mockZipCity}
      displayAddress={false}
      handleLocationOnChange={handleZipCityOnChange}
    />
  );
};

export const DutyStationWithRequiredAsterisk = () => {
  const [value, setValue] = useState();

  const onChange = (newValue) => {
    setValue(newValue);
  };

  return (
    <LocationSearchBoxComponent
      input={{ name: 'test_component', onChange, value }}
      title="Duty Station Test Component"
      displayAddress={false}
      searchLocations={mockSearchLocations}
      showAddress={mockShowAddress}
      showRequiredAsterisk
    />
  );
};

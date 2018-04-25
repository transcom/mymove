import React from 'react';
import { render } from 'react-dom';
// import _ from 'lodash';
import { RetrieveDutyStations } from './api.js';

// Import React Table
import ReactTable from 'react-table';
import 'react-table/react-table.css';

// To see data pulled from the API, swap rawData in wherever hardcodedData is used.
// We used hardcoded data to avoid the problem of auth in the VM. It worked :)
const rawData = RetrieveDutyStations();

const hardcodedData = [
  {
    address: {
      city: 'Altus AFB',
      country: 'United States',
      postal_code: '73523',
      state: 'OK',
      street_address_1: 'n/a',
    },
    branch: 'AIRFORCE',
    created_at: '2018-04-20T21:50:42.037Z',
    id: 'a1f25199-543b-43a2-9728-563a5aa9e460',
    name: 'Altus AFB',
    updated_at: '2018-04-20T21:50:42.037Z',
  },
  {
    address: {
      city: 'Barksdale AFB',
      country: 'United States',
      postal_code: '71110',
      state: 'LA',
      street_address_1: 'n/a',
    },
    branch: 'AIRFORCE',
    created_at: '2018-04-20T21:50:42.037Z',
    id: 'e881cb42-88f3-4238-9ffa-0f6f5c72f7fd',
    name: 'Barksdale AFB',
    updated_at: '2018-04-20T21:50:42.037Z',
  },
  {
    address: {
      city: 'Columbus',
      country: 'United States',
      postal_code: '39710',
      state: 'MS',
      street_address_1: 'n/a',
    },
    branch: 'AIRFORCE',
    created_at: '2018-04-20T21:50:42.037Z',
    id: '183decc4-c035-4a0b-a464-0fbb73e36eb9',
    name: 'Columbus AFB',
    updated_at: '2018-04-20T21:50:42.037Z',
  },
  {
    address: {
      city: 'Dyess AFB',
      country: 'United States',
      postal_code: '79607',
      state: 'TX',
      street_address_1: 'n/a',
    },
    branch: 'AIRFORCE',
    created_at: '2018-04-20T21:50:42.037Z',
    id: '8a049f2a-4a7b-4062-aefc-3288f8b93207',
    name: 'Dyess AFB',
    updated_at: '2018-04-20T21:50:42.037Z',
  },
  {
    address: {
      city: 'Eglin AFB',
      country: 'United States',
      postal_code: '32542',
      state: 'FL',
      street_address_1: 'n/a',
    },
    branch: 'AIRFORCE',
    created_at: '2018-04-20T21:50:42.037Z',
    id: '9a460f99-882b-4f38-b30e-b20e132bfd77',
    name: 'Eglin AFB',
    updated_at: '2018-04-20T21:50:42.037Z',
  },
  {
    address: {
      city: 'Goodfellow AFB',
      country: 'United States',
      postal_code: '76908',
      state: 'TX',
      street_address_1: 'n/a',
    },
    branch: 'AIRFORCE',
    created_at: '2018-04-20T21:50:42.037Z',
    id: 'c3d1a219-fa0e-4cd6-a572-b9bf88e88ee5',
    name: 'Goodfellow AFB',
    updated_at: '2018-04-20T21:50:42.037Z',
  },
  {
    address: {
      city: 'Hurlburt Field',
      country: 'United States',
      postal_code: '32544',
      state: 'FL',
      street_address_1: 'n/a',
    },
    branch: 'AIRFORCE',
    created_at: '2018-04-20T21:50:42.037Z',
    id: 'f06c5dd9-e99b-45f4-8f1b-4bd3b030fdf0',
    name: 'Hurlburt Field AFB',
    updated_at: '2018-04-20T21:50:42.037Z',
  },
  {
    address: {
      city: 'Biloxi',
      country: 'United States',
      postal_code: '39534',
      state: 'MS',
      street_address_1: 'n/a',
    },
    branch: 'AIRFORCE',
    created_at: '2018-04-20T21:50:42.037Z',
    id: '732605c9-f688-44da-9009-2c504e51f86d',
    name: 'Keesler AFB',
    updated_at: '2018-04-20T21:50:42.037Z',
  },
  {
    address: {
      city: 'Lackland AFB',
      country: 'United States',
      postal_code: '78236',
      state: 'TX',
      street_address_1: 'n/a',
    },
    branch: 'AIRFORCE',
    created_at: '2018-04-20T21:50:42.037Z',
    id: '7e42d0ec-61ef-4419-b7f2-3c72de774374',
    name: 'Lackland AFB',
    updated_at: '2018-04-20T21:50:42.037Z',
  },
  {
    address: {
      city: 'Laughlin AFB',
      country: 'United States',
      postal_code: '78843',
      state: 'TX',
      street_address_1: 'n/a',
    },
    branch: 'AIRFORCE',
    created_at: '2018-04-20T21:50:42.037Z',
    id: '3cf37423-7f99-40e1-af05-02d320ca646b',
    name: 'Laughlin AFB',
    updated_at: '2018-04-20T21:50:42.037Z',
  },
  {
    address: {
      city: 'Little Rock AFB',
      country: 'United States',
      postal_code: '72099',
      state: 'AR',
      street_address_1: 'n/a',
    },
    branch: 'AIRFORCE',
    created_at: '2018-04-20T21:50:42.037Z',
    id: '6305ff63-1f7d-43b7-b438-bedb8ad044de',
    name: 'Little Rock AFB',
    updated_at: '2018-04-20T21:50:42.037Z',
  },
  {
    address: {
      city: 'Tampa',
      country: 'United States',
      postal_code: '33621',
      state: 'FL',
      street_address_1: 'n/a',
    },
    branch: 'AIRFORCE',
    created_at: '2018-04-20T21:50:42.037Z',
    id: 'ad78e596-8657-4bc7-be0f-645581f2f6fe',
    name: 'MacDill AFB',
    updated_at: '2018-04-20T21:50:42.037Z',
  },
  {
    address: {
      city: 'Montgomery',
      country: 'United States',
      postal_code: '36112',
      state: 'AL',
      street_address_1: 'n/a',
    },
    branch: 'AIRFORCE',
    created_at: '2018-04-20T21:50:42.037Z',
    id: '499bdfb5-9701-4aed-9478-1ce142e826bb',
    name: 'Maxwell AFB',
    updated_at: '2018-04-20T21:50:42.037Z',
  },
  {
    address: {
      city: 'Moody AFB',
      country: 'United States',
      postal_code: '31699',
      state: 'GA',
      street_address_1: 'n/a',
    },
    branch: 'AIRFORCE',
    created_at: '2018-04-20T21:50:42.037Z',
    id: '871038c1-2037-40f1-beb5-aa56cc88690a',
    name: 'Moody AFB',
    updated_at: '2018-04-20T21:50:42.037Z',
  },
  {
    address: {
      city: 'Patrick AFB',
      country: 'United States',
      postal_code: '32925',
      state: 'FL',
      street_address_1: 'n/a',
    },
    branch: 'AIRFORCE',
    created_at: '2018-04-20T21:50:42.037Z',
    id: '522376bd-3eda-4171-a2f2-ba312ccbf722',
    name: 'Patrick AFB',
    updated_at: '2018-04-20T21:50:42.037Z',
  },
  {
    address: {
      city: 'Randolph AFB',
      country: 'United States',
      postal_code: '78150',
      state: 'TX',
      street_address_1: 'n/a',
    },
    branch: 'AIRFORCE',
    created_at: '2018-04-20T21:50:42.037Z',
    id: '5bc53f20-b675-4b1f-a6da-30e987c931c2',
    name: 'Randolph AFB',
    updated_at: '2018-04-20T21:50:42.037Z',
  },
  {
    address: {
      city: 'Warner Robins',
      country: 'United States',
      postal_code: '31098',
      state: 'GA',
      street_address_1: 'n/a',
    },
    branch: 'AIRFORCE',
    created_at: '2018-04-20T21:50:42.037Z',
    id: '36af5a76-4bbd-49b1-924a-63ac52af4c6c',
    name: 'Robins AFB',
    updated_at: '2018-04-20T21:50:42.037Z',
  },
  {
    address: {
      city: 'Sheppard AFB',
      country: 'United States',
      postal_code: '76311',
      state: 'TX',
      street_address_1: 'n/a',
    },
    branch: 'AIRFORCE',
    created_at: '2018-04-20T21:50:42.037Z',
    id: '1957468e-6ffd-4ff1-8e83-1857bc412c7c',
    name: 'Sheppard AFB',
    updated_at: '2018-04-20T21:50:42.037Z',
  },
  {
    address: {
      city: 'Oklahoma City',
      country: 'United States',
      postal_code: '73145',
      state: 'OK',
      street_address_1: 'n/a',
    },
    branch: 'AIRFORCE',
    created_at: '2018-04-20T21:50:42.037Z',
    id: '9e67fd93-22d5-4b97-b3c7-6ca1cd759c38',
    name: 'Tinker AFB',
    updated_at: '2018-04-20T21:50:42.037Z',
  },
  {
    address: {
      city: 'Panama City',
      country: 'United States',
      postal_code: '32403',
      state: 'FL',
      street_address_1: 'n/a',
    },
    branch: 'AIRFORCE',
    created_at: '2018-04-20T21:50:42.037Z',
    id: 'eef7d1e2-8e31-410e-9d0c-2b1c1d570145',
    name: 'Tyndall AFB',
    updated_at: '2018-04-20T21:50:42.037Z',
  },
  {
    address: {
      city: 'Enid',
      country: 'United States',
      postal_code: '73705',
      state: 'OK',
      street_address_1: 'n/a',
    },
    branch: 'AIRFORCE',
    created_at: '2018-04-20T21:50:42.037Z',
    id: '3961723e-4599-4db5-ae8f-eff56ca5f480',
    name: 'Vance AFB',
    updated_at: '2018-04-20T21:50:42.037Z',
  },
  {
    address: {
      city: 'Beale AFB',
      country: 'United States',
      postal_code: '95903',
      state: 'CA',
      street_address_1: 'n/a',
    },
    branch: 'AIRFORCE',
    created_at: '2018-04-20T21:50:42.037Z',
    id: '99cd5892-f4a5-438c-a72e-789cb07c9622',
    name: 'Beale AFB',
    updated_at: '2018-04-20T21:50:42.037Z',
  },
  {
    address: {
      city: 'Aurora',
      country: 'United States',
      postal_code: '80011',
      state: 'CO',
      street_address_1: 'n/a',
    },
    branch: 'AIRFORCE',
    created_at: '2018-04-20T21:50:42.037Z',
    id: '2fb078cb-d2c7-40b5-8f55-d96b62ddd9be',
    name: 'Buckley AFB',
    updated_at: '2018-04-20T21:50:42.037Z',
  },
  {
    address: {
      city: 'Cannon AFB',
      country: 'United States',
      postal_code: '88103',
      state: 'NM',
      street_address_1: 'n/a',
    },
    branch: 'AIRFORCE',
    created_at: '2018-04-20T21:50:42.037Z',
    id: 'ce3a3dbb-4589-47d5-a92a-8b4b62e1c251',
    name: 'Cannon AFB',
    updated_at: '2018-04-20T21:50:42.037Z',
  },
  {
    address: {
      city: 'Indian Springs',
      country: 'United States',
      postal_code: '89018',
      state: 'NV',
      street_address_1: 'n/a',
    },
    branch: 'AIRFORCE',
    created_at: '2018-04-20T21:50:42.037Z',
    id: 'a5b57512-ea80-4a65-851f-57fd09df7db1',
    name: 'Creech AFB',
    updated_at: '2018-04-20T21:50:42.037Z',
  },
  {
    address: {
      city: 'Tucson',
      country: 'United States',
      postal_code: '85707',
      state: 'AZ',
      street_address_1: 'n/a',
    },
    branch: 'AIRFORCE',
    created_at: '2018-04-20T21:50:42.037Z',
    id: 'a68dc050-d17b-443b-8d30-36bfd8a9599a',
    name: 'Davis Monthan AFB',
    updated_at: '2018-04-20T21:50:42.037Z',
  },
  {
    address: {
      city: 'Edwards',
      country: 'United States',
      postal_code: '93524',
      state: 'CA',
      street_address_1: 'n/a',
    },
    branch: 'AIRFORCE',
    created_at: '2018-04-20T21:50:42.037Z',
    id: 'd49f4281-dbef-45bd-b258-365d50db4f07',
    name: 'Edwards AFB',
    updated_at: '2018-04-20T21:50:42.037Z',
  },
  {
    address: {
      city: 'Ellsworth AFB',
      country: 'United States',
      postal_code: '57706',
      state: 'SD',
      street_address_1: 'n/a',
    },
    branch: 'AIRFORCE',
    created_at: '2018-04-20T21:50:42.037Z',
    id: 'd19ef524-924a-4671-b178-9de4b7205208',
    name: 'Ellsworth AFB',
    updated_at: '2018-04-20T21:50:42.037Z',
  },
  {
    address: {
      city: 'F.E. Warren AFB',
      country: 'United States',
      postal_code: '82005',
      state: 'WY',
      street_address_1: 'n/a',
    },
    branch: 'AIRFORCE',
    created_at: '2018-04-20T21:50:42.037Z',
    id: '0684325d-3c64-4d3a-8fc7-9a809daaed89',
    name: 'F.E. Warren AFB',
    updated_at: '2018-04-20T21:50:42.037Z',
  },
  {
    address: {
      city: 'Spokane',
      country: 'United States',
      postal_code: '99208',
      state: 'WA',
      street_address_1: 'n/a',
    },
    branch: 'AIRFORCE',
    created_at: '2018-04-20T21:50:42.037Z',
    id: '39a863bb-dcfd-498a-8827-e20d201363a7',
    name: 'Fairchild AFB',
    updated_at: '2018-04-20T21:50:42.037Z',
  },
  {
    address: {
      city: 'Grand Forks AFB',
      country: 'United States',
      postal_code: '58205',
      state: 'ND',
      street_address_1: 'n/a',
    },
    branch: 'AIRFORCE',
    created_at: '2018-04-20T21:50:42.037Z',
    id: 'f9afc55b-7499-41ea-8c73-6073ef464b26',
    name: 'Grand Forks AFB',
    updated_at: '2018-04-20T21:50:42.037Z',
  },
  {
    address: {
      city: 'Hill AFB',
      country: 'United States',
      postal_code: '84056',
      state: 'UT',
      street_address_1: 'n/a',
    },
    branch: 'AIRFORCE',
    created_at: '2018-04-20T21:50:42.037Z',
    id: 'e2eb09e5-7bf5-4358-b176-f7abf4940e81',
    name: 'Hill AFB',
    updated_at: '2018-04-20T21:50:42.037Z',
  },
  {
    address: {
      city: 'Holloman AFB',
      country: 'United States',
      postal_code: '88330',
      state: 'NM',
      street_address_1: 'n/a',
    },
    branch: 'AIRFORCE',
    created_at: '2018-04-20T21:50:42.037Z',
    id: 'b7927c63-3f4a-4c35-ae69-03303ac46384',
    name: 'Holloman AFB',
    updated_at: '2018-04-20T21:50:42.037Z',
  },
  {
    address: {
      city: 'Kortland AFB',
      country: 'United States',
      postal_code: '87117',
      state: 'NM',
      street_address_1: 'n/a',
    },
    branch: 'AIRFORCE',
    created_at: '2018-04-20T21:50:42.037Z',
    id: '1600e625-9f74-4ac8-8d5d-07ff8f035d68',
    name: 'Kirtland AFB',
    updated_at: '2018-04-20T21:50:42.037Z',
  },
  {
    address: {
      city: 'Glendale Luke AFB',
      country: 'United States',
      postal_code: '85309',
      state: 'AZ',
      street_address_1: 'n/a',
    },
    branch: 'AIRFORCE',
    created_at: '2018-04-20T21:50:42.037Z',
    id: 'e57ab6f2-f75a-4d36-83e1-3eedd88c9a16',
    name: 'Luke AFB',
    updated_at: '2018-04-20T21:50:42.037Z',
  },
  {
    address: {
      city: 'Malmstrom AFB',
      country: 'United States',
      postal_code: '59402',
      state: 'MT',
      street_address_1: 'n/a',
    },
    branch: 'AIRFORCE',
    created_at: '2018-04-20T21:50:42.037Z',
    id: 'fc7203ff-60d1-44e3-9105-29f0f24a90d5',
    name: 'Malmstrom AFB',
    updated_at: '2018-04-20T21:50:42.037Z',
  },
  {
    address: {
      city: 'McConnell AFB',
      country: 'United States',
      postal_code: '67221',
      state: 'KS',
      street_address_1: 'n/a',
    },
    branch: 'AIRFORCE',
    created_at: '2018-04-20T21:50:42.037Z',
    id: '34de74ea-5886-4be8-bd6e-cce7ccf7093d',
    name: 'McConnell AFB',
    updated_at: '2018-04-20T21:50:42.037Z',
  },
  {
    address: {
      city: 'Minot AFB',
      country: 'United States',
      postal_code: '58705',
      state: 'ND',
      street_address_1: 'n/a',
    },
    branch: 'AIRFORCE',
    created_at: '2018-04-20T21:50:42.037Z',
    id: 'bbdee24a-20f4-4e46-95d9-f6effead55f9',
    name: 'Minot AFB',
    updated_at: '2018-04-20T21:50:42.037Z',
  },
  {
    address: {
      city: 'Mountain Home AFB',
      country: 'United States',
      postal_code: '83648',
      state: 'ID',
      street_address_1: 'n/a',
    },
    branch: 'AIRFORCE',
    created_at: '2018-04-20T21:50:42.037Z',
    id: 'aa2ae5b1-1c2d-4906-9243-c0ebc291eaca',
    name: 'Mountain Home AFB',
    updated_at: '2018-04-20T21:50:42.037Z',
  },
  {
    address: {
      city: 'Nellis AFB',
      country: 'United States',
      postal_code: '89191',
      state: 'NV',
      street_address_1: 'n/a',
    },
    branch: 'AIRFORCE',
    created_at: '2018-04-20T21:50:42.037Z',
    id: 'ab0c7b90-5519-4872-871d-cf1afd72afdd',
    name: 'Nellis AFB',
    updated_at: '2018-04-20T21:50:42.037Z',
  },
  {
    address: {
      city: 'Offutt AFB',
      country: 'United States',
      postal_code: '68113',
      state: 'NE',
      street_address_1: 'n/a',
    },
    branch: 'AIRFORCE',
    created_at: '2018-04-20T21:50:42.037Z',
    id: '549f42db-bc7e-464a-a771-e00ca6ed2830',
    name: 'Offutt AFB',
    updated_at: '2018-04-20T21:50:42.037Z',
  },
  {
    address: {
      city: 'Colorado Springs',
      country: 'United States',
      postal_code: '80916',
      state: 'CO',
      street_address_1: 'n/a',
    },
    branch: 'AIRFORCE',
    created_at: '2018-04-20T21:50:42.037Z',
    id: 'b290f232-9de4-4fbb-8924-918fcf29184f',
    name: 'Peterson AFB',
    updated_at: '2018-04-20T21:50:42.037Z',
  },
  {
    address: {
      city: 'Colorado Springs',
      country: 'United States',
      postal_code: '80912',
      state: 'CO',
      street_address_1: 'n/a',
    },
    branch: 'AIRFORCE',
    created_at: '2018-04-20T21:50:42.037Z',
    id: '435dd49e-befb-4d09-9be6-321e7232580a',
    name: 'Schiever AFB',
    updated_at: '2018-04-20T21:50:42.037Z',
  },
  {
    address: {
      city: 'Travis AFB',
      country: 'United States',
      postal_code: '94535',
      state: 'CA',
      street_address_1: 'n/a',
    },
    branch: 'AIRFORCE',
    created_at: '2018-04-20T21:50:42.037Z',
    id: '155ae0e3-90bb-4033-9215-722f3b590983',
    name: 'Travis AFB',
    updated_at: '2018-04-20T21:50:42.037Z',
  },
  {
    address: {
      city: 'Lompoc',
      country: 'United States',
      postal_code: '93437',
      state: 'CA',
      street_address_1: 'n/a',
    },
    branch: 'AIRFORCE',
    created_at: '2018-04-20T21:50:42.037Z',
    id: 'df3b9549-9d70-4668-b832-096c3f5d28e1',
    name: 'Vandenberg AFB',
    updated_at: '2018-04-20T21:50:42.037Z',
  },
  {
    address: {
      city: 'Whiteman AFB',
      country: 'United States',
      postal_code: '65305',
      state: 'MO',
      street_address_1: 'n/a',
    },
    branch: 'AIRFORCE',
    created_at: '2018-04-20T21:50:42.037Z',
    id: '7a88b71c-b5aa-4305-b810-41cf4c0dcde3',
    name: 'Whiteman AFB',
    updated_at: '2018-04-20T21:50:42.037Z',
  },
];

// Example based on the one here: https://react-table.js.org/#/story/server-side-data
const requestData = (pageSize, page, sorted, filtered) => {
  return new Promise((resolve, reject) => {
    // You can retrieve your data however you want, in this case, we will just use some local data.
    // let filteredData = rawData;

    // You can use the filters in your request, but you are responsible for applying them.
    // if (filtered.length) {
    //   filteredData = filtered.reduce((filteredSoFar, nextFilter) => {
    //     return filteredSoFar.filter(row => {
    //       return (row[nextFilter.id] + '').includes(nextFilter.value);
    //     });
    //   }, filteredData);
    // }
    // You can also use the sorting in your request, but again, you are responsible for applying it.
    // const sortedData = _.orderBy(
    //   filteredData,
    //   sorted.map(sort => {
    //     return row => {
    //       if (row[sort.id] === null || row[sort.id] === undefined) {
    //         return -Infinity;
    //       }
    //       return typeof row[sort.id] === 'string'
    //         ? row[sort.id].toLowerCase()
    //         : row[sort.id];
    //     };
    //   }),
    //   sorted.map(d => (d.desc ? 'desc' : 'asc')),
    // );

    // You must return an object containing the rows of the current page, and optionally the total pages number.
    const res = {
      rows: hardcodedData.slice(pageSize * page, pageSize * page + pageSize),
      pages: Math.ceil(hardcodedData.length / pageSize),
    };

    // Here we'll simulate a server response with 500ms of delay.
    setTimeout(() => resolve(res), 500);
  });
};

class Admin extends React.Component {
  constructor() {
    super();
    this.state = {
      data: [],
      pages: null,
      loading: true,
    };
    this.fetchData = this.fetchData.bind(this);
  }
  fetchData(state, instance) {
    // Whenever the table model changes, or the user sorts or changes pages, this method gets called and passed the current table model.
    // You can set the `loading` prop of the table to true to use the built-in one or show you're own loading bar if you want.
    this.setState({ loading: true });
    // Request the data however you want.  Here, we'll use our mocked service we created earlier
    requestData(state.pageSize, state.page, state.sorted, state.filtered).then(
      res => {
        // Now just get the rows of data to your React Table (and update anything else like total pages or loading)
        this.setState({
          data: res.rows,
          pages: res.pages,
          loading: false,
        });
      },
    );
  }
  render() {
    const { data, pages, loading } = this.state;
    return (
      <div>
        <ReactTable
          columns={[
            {
              Header: 'Flair',
              Cell: row => (
                <div>
                  <img src="https://upload.wikimedia.org/wikipedia/commons/thumb/f/f1/Heart_coraz%C3%B3n.svg/130px-Heart_coraz%C3%B3n.svg.png" />
                </div>
              ),
            },
            {
              Header: 'ID',
              accessor: 'id',
            },
            {
              Header: 'Name',
              accessor: 'name',
            },
            {
              Header: 'Branch',
              accessor: 'branch',
            },
            {
              Header: 'Created At',
              accessor: 'created_at',
            },
            {
              Header: 'Updated At',
              accessor: 'updated_at',
            },
          ]}
          // manual // Forces table not to paginate or sort automatically, so we can handle it server-side
          data={data}
          pages={pages} // Display the total number of pages
          loading={loading} // Display the loading overlay when we need it
          onFetchData={this.fetchData} // Request new data when things change
          filterable
          defaultPageSize={50}
          className="-striped -highlight"
        />
        <br />
      </div>
    );
  }
}

export default Admin;

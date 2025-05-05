import React from 'react';
import { object } from '@storybook/addon-knobs';

import CustomerInfoList from '../DefinitionLists/CustomerInfoList';

import DetailsPanel from './DetailsPanel';

import { ReviewButton } from 'components/form/IconButtons';
import ButtonDropdown from 'components/ButtonDropdown/ButtonDropdown';

export default {
  title: 'Office Components/DetailsPanel',
  component: DetailsPanel,
};

export const Basic = () => (
  <div className="officeApp">
    <DetailsPanel title="Details panel">
      <p>Child content!</p>
    </DetailsPanel>
  </div>
);

export const WithEditButton = () => (
  <div className="officeApp">
    <DetailsPanel
      title="Details panel with edit button"
      editButton={
        <a href="#" className="usa-button usa-button--secondary">
          Edit
        </a>
      }
    >
      <p>Child content!</p>
    </DetailsPanel>
  </div>
);

export const WithDropdownButton = () => (
  <div className="officeApp">
    <DetailsPanel
      title="With Dropdown Button"
      editButton={
        <ButtonDropdown>
          <option value="">Dropdown Button</option>
          <option>Option 1</option>
          <option>Option 2</option>
          <option>Option 3</option>
          <option>Option 4</option>
        </ButtonDropdown>
      }
    >
      <p>Child Content!</p>
    </DetailsPanel>
  </div>
);

export const WithReviewButton = () => (
  <div className="officeApp">
    <DetailsPanel title="With Review Button" reviewButton={<ReviewButton label="Review Button" secondary />}>
      <p>Child Content!</p>
    </DetailsPanel>
  </div>
);

export const WithEditAndReviewButton = () => (
  <div className="officeApp">
    <DetailsPanel
      title="With Edit and Review Button"
      editButton={
        <a href="#" className="usa-button usa-button--secondary">
          Edit
        </a>
      }
      reviewButton={<ReviewButton label="Review Button" secondary />}
    >
      <p>Child Content!</p>
    </DetailsPanel>
  </div>
);

export const WithDropdownAndReviewButton = () => (
  <div className="officeApp">
    <DetailsPanel
      title="With Dropdown And Review Button"
      reviewButton={<ReviewButton label="Review Button" secondary />}
      editButton={
        <ButtonDropdown>
          <option value="">Dropdown Button</option>
          <option>Option 1</option>
          <option>Option 2</option>
          <option>Option 3</option>
          <option>Option 4</option>
        </ButtonDropdown>
      }
    >
      <p>Child Content!</p>
    </DetailsPanel>
  </div>
);

export const WithClassname = () => (
  <div className="officeApp">
    <DetailsPanel title="Details panel with added CSS class" className="border-2px">
      <p>I have a border class added via props!</p>
    </DetailsPanel>
  </div>
);

export const WithTag = () => (
  <div className="officeApp">
    <DetailsPanel title="Details panel with tag" tag="NEW">
      <p>I have a tag added via props!</p>
    </DetailsPanel>
  </div>
);

const info = {
  name: 'Smith, Kerry',
  dodId: '9999999999',
  phone: '+1 999-999-9999',
  email: 'ksmith@email.com',
  currentAddress: {
    streetAddress1: '812 S 129th St',
    city: 'San Antonio',
    state: 'TX',
    postalCode: '78234',
  },
  backupContact: {
    name: 'Quinn Ocampo',
    email: 'quinnocampo@myemail.com',
    phone: '999-999-9999',
  },
};

export const WithDefinitionListComponent = () => (
  <div className="officeApp">
    <DetailsPanel
      title="Customer info"
      editButton={
        <a href="#" className="usa-button usa-button--secondary">
          Edit customer info
        </a>
      }
    >
      <CustomerInfoList customerInfo={object('customerInfo', info)} />
    </DetailsPanel>
  </div>
);

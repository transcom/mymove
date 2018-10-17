import React, { Fragment } from 'react';
import PropTypes from 'prop-types';

import { PanelField } from 'shared/EditablePanel';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';

export const AddressElementDisplay = ({ address, title }) => (
  <PanelField title={title}>
    {address.street_address_1 && (
      <span>
        {address.street_address_1}
        <br />
      </span>
    )}
    {address.street_address_2 && (
      <span>
        {address.street_address_2}
        <br />
      </span>
    )}
    {address.street_address_3 && (
      <span>
        {address.street_address_3}
        <br />
      </span>
    )}
    {address.city}, {address.state} {address.postal_code}
  </PanelField>
);

AddressElementDisplay.propTypes = {
  address: PropTypes.shape({
    street_address_1: PropTypes.string,
    street_address_2: PropTypes.string,
    street_address_3: PropTypes.string,
    city: PropTypes.string.isRequired,
    state: PropTypes.string.isRequired,
    postal_code: PropTypes.string.isRequired,
  }),
  title: PropTypes.string.isRequired,
};

export const AddressElementEdit = ({ addressProps, title }) => (
  <Fragment>
    <div className="panel-subhead">{title}</div>
    <SwaggerField fieldName="street_address_1" {...addressProps} required />
    <SwaggerField fieldName="street_address_2" {...addressProps} />
    <SwaggerField fieldName="street_address_3" {...addressProps} />
    <SwaggerField fieldName="city" {...addressProps} required />
    <SwaggerField fieldName="state" {...addressProps} required />
    <SwaggerField fieldName="postal_code" {...addressProps} required />
  </Fragment>
);

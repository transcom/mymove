import React from 'react';

import { PanelField } from 'shared/EditablePanel';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';

export function addressElementDisplay(address, title) {
  return (
    <React.Fragment>
      <PanelField title={title}>
        {address.street_address_1}
        <br />
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
    </React.Fragment>
  );
}

export function addressElementEdit(addressProps, title) {
  return (
    <React.Fragment>
      <div className="panel-subhead">{title}</div>
      <SwaggerField fieldName="street_address_1" {...addressProps} required />
      <SwaggerField fieldName="street_address_2" {...addressProps} />
      <SwaggerField fieldName="city" {...addressProps} required />
      <SwaggerField fieldName="state" {...addressProps} required />
      <SwaggerField fieldName="postal_code" {...addressProps} required />
    </React.Fragment>
  );
}

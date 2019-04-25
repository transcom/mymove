import React, { Fragment } from 'react';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import { AddressElementEdit } from '../Address';
import { get } from 'lodash';

export const Code125Form = props => {
  const { ship_line_item_schema } = props;
  return (
    <Fragment>
      <SwaggerField
        className="textarea-half"
        title="Reason for service"
        fieldName="reason"
        swagger={ship_line_item_schema}
        required
      />
      <SwaggerField title="Date of service" fieldName="date" swagger={ship_line_item_schema} required />
      <SwaggerField title="Time of service" fieldName="time" swagger={ship_line_item_schema} />
      <AddressElementEdit
        fieldName="address"
        schema={get(ship_line_item_schema, 'properties.address')}
        title="Truck-to-truck transfer location"
        zipPattern="USA"
      />
      <div className="bq-explanation">
        <p>Enter amount after service is completed</p>
      </div>
    </Fragment>
  );
};

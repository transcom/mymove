import React, { Fragment } from 'react';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';

export const Code226Form = props => {
  const { ship_line_item_schema } = props;
  return (
    <Fragment>
      <SwaggerField
        className="textarea-half"
        title="Description of charge"
        fieldName="description"
        swagger={ship_line_item_schema}
        required
      />
      <SwaggerField
        className="textarea-half"
        title="Reason for charge"
        fieldName="reason"
        swagger={ship_line_item_schema}
        required
      />
      <SwaggerField title="Amount" fieldName="actual_amount_cents" swagger={ship_line_item_schema} required />
    </Fragment>
  );
};

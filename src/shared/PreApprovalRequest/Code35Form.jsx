import React, { Fragment } from 'react';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import { Code35FormAlert } from './Code35FormAlert';

export const Code35Form = props => {
  const { ship_line_item_schema } = props;
  return (
    <Fragment>
      <SwaggerField
        className="textarea-half"
        title="Description of service"
        fieldName="description"
        swagger={ship_line_item_schema}
        required
      />
      <SwaggerField
        className="textarea-half"
        title="Reason for service"
        fieldName="reason"
        swagger={ship_line_item_schema}
        required
      />
      <SwaggerField
        title="Estimate, not to exceed"
        fieldName="estimate_amount_cents"
        swagger={ship_line_item_schema}
        required
      />
      <SwaggerField title="Actual amount of service" fieldName="actual_amount_cents" swagger={ship_line_item_schema} />
      <div className="bq-explanation">
        <p>Enter amount after service is completed</p>
      </div>
      <Code35FormAlert showAlert={props.showAlert} />
    </Fragment>
  );
};

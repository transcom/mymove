import React, { Fragment } from 'react';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import { DimensionsField } from '../JsonSchemaForm/DimensionsField';

export const Code105Details = props => {
  const { ship_line_item_schema } = props;
  return (
    <Fragment>
      <SwaggerField fieldName="description" swagger={ship_line_item_schema} required />
      <DimensionsField
        fieldName="item_dimensions"
        swagger={ship_line_item_schema}
        labelText="Item Dimensions (inches)"
        isRequired={true}
      />
      <DimensionsField
        fieldName="crate_dimensions"
        swagger={ship_line_item_schema}
        labelText="Crate Dimensions (inches)"
        isRequired={true}
      />
      <div className="bq-explanation">
        <p>Crate can only exceed item size by:</p>
        <ul>
          <li>
            <em>Internal crate</em>: Up to 3" larger
          </li>
          <li>
            <em>External crate</em>: Up to 5" larger
          </li>
        </ul>
      </div>
    </Fragment>
  );
};

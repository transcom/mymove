import React, { Fragment } from 'react';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';

export const DefaultForm = props => {
  const { ship_line_item_schema } = props;
  return (
    <Fragment>
      <SwaggerField fieldName="quantity_1" className="half-width" swagger={ship_line_item_schema} required />
      <div className="bq-explanation">
        <p>
          Enter numbers only, no symbols or units. <em>Examples:</em>
        </p>
        <ul>
          <li>
            Crating: enter "<strong>47.4</strong>" for crate size of 47.4 cu. ft.
          </li>
          <li>
            {' '}
            3rd-party service: enter "<strong>1299.99</strong>" for cost of $1,299.99.
          </li>
          <li>
            Bulky item: enter "<strong>1</strong>" for a single item.
          </li>
        </ul>
      </div>
    </Fragment>
  );
};

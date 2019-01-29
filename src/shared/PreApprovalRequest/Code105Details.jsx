import React, { Fragment } from 'react';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';

export const Code105Details = props => {
  const { ship_line_item_schema } = props;
  return (
    <Fragment>
      <SwaggerField fieldName="description" swagger={ship_line_item_schema} required />
    </Fragment>
  );
};

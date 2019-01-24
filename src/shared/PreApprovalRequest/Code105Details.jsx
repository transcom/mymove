import React, { Fragment } from 'react';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';


export const Code105Details = props => {
  const { swagger } = props;
  return (
    <Fragment>
      <div>More to come!</div>
      <SwaggerField
        fieldName="description"
        className="three-quarter-width"
        swagger={swagger}
        required
      />
    </Fragment>
  );
};

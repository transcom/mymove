import React from 'react';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';

const AddressForm = ({ schema }) => (
  <>
    <div className="grid-row">
      <div className="grid-col-12">
        <SwaggerField fieldName="street_address_1" swagger={schema} required />
        <SwaggerField fieldName="street_address_2" swagger={schema} />
      </div>
    </div>
    <div className="grid-row grid-gap">
      <div className="tablet:grid-col-4 grid-col-8">
        <SwaggerField fieldName="city" swagger={schema} required />
      </div>
      <div className="tablet:grid-col-2 grid-col-4">
        <SwaggerField fieldName="state" swagger={schema} required />
      </div>
    </div>
    <div className="grid-row grid-gap">
      <div className="tablet:grid-col-2 grid-col-4">
        <SwaggerField fieldName="postal_code" swagger={schema} required />
      </div>
    </div>
  </>
);

export default AddressForm;

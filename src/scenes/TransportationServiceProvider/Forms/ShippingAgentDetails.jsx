import PropTypes from 'prop-types';
import React, { Component } from 'react';

import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';

class ShippingAgentDetails extends Component {
  render() {
    return (
      <form onSubmit={acceptFormShipment}>
        <h3 className="smheading">Origin Shipping Agent</h3>
        <SwaggerField fieldName="origin_agent_name" swagger={schema} required />
        <SwaggerField
          fieldName="origin_agent_phone_number"
          swagger={schema}
          required
        />
        <SwaggerField fieldName="origin_agent_email" swagger={schema} required />

        <h3 className="smheading">Destination Shipping Agent</h3>
        <SwaggerField
          fieldName="destination_agent_name"
          swagger={schema}
          required
        />
        <SwaggerField
          fieldName="destination_agent_phone_number"
          swagger={schema}
          required
        />
        <SwaggerField
          fieldName="destination_agent_email"
          swagger={schema}
          required
        />

        <div className="usa-grid">
          <div className="usa-width-one-whole extras options">
            <a onClick={setButtonState}>Never mind</a>
          </div>
          <div className="usa-width-one-whole extras options">
            <button type="submit">Submit</button>
          </div>
        </div>
      </form>
    );
  }
}

ShippingAgentDetails.propTypes = {
  schema: PropTypes.object.isRequired,
  error: PropTypes.object,
  formValues: PropTypes.object,
};

export default ShippingAgentDetails;

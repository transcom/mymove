import React, { Component } from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { get, capitalize } from 'lodash';

import { getFormValues, reduxForm } from 'redux-form';
import { NavLink } from 'react-router-dom';

import Alert from 'shared/Alert'; // eslint-disable-line
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import { withContext } from 'shared/AppContext';

import { loadShipmentDependencies, acceptShipment } from './ducks';
import { formatDate } from 'shared/formatters';

const shipmentAcceptFormName = 'shipment_accept';

let ShipmentAcceptForm = props => {
  const { schema, setButtonState, acceptFormShipment } = props;

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
};

ShipmentAcceptForm = reduxForm({
  form: shipmentAcceptFormName,
})(ShipmentAcceptForm);
